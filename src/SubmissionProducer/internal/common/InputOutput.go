package common

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
)

// Input
// -----------------------

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// CheckArgs should be used to ensure the right command line arguments are
// passed before executing an example.
func CheckArgs(arg ...string) {
	if len(os.Args) < len(arg)+1 {
		Warning("Usage: %s %s", os.Args[0], strings.Join(arg, " "))
		os.Exit(1)
	}
}

// Console Out
// -----------------------

func Error(mess string) {
	fmt.Println(mess)
	os.Exit(1)
}

// CheckIfError should be used to naively panics if an error is not nil.
func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n\n", fmt.Sprintf("error: %s", err))
	debug.PrintStack()

	os.Exit(1)
}

// CheckIfErrorWithMessage should be used to naively panics if an error is not nil with a message.
func CheckIfErrorWithMessage(err error, message string) {
	if err == nil {
		return
	}

	fmt.Println(message)
	fmt.Printf("\x1b[31;1m%s\x1b[0m\n\n", fmt.Sprintf("error: %s", err))
	debug.PrintStack()

	os.Exit(1)
}

// Info should be used to describe the example commands that are about to run.
func Info(format string, args ...interface{}) {
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

// Warning should be used to display a warning
func Warning(format string, args ...interface{}) {
	fmt.Printf("\x1b[36;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

// File IO
// -----------------------

func MakeDir(path string) {
	if err := os.MkdirAll(path, OwnerPermRw); err != nil {
		if !os.IsExist(err) {
			CheckIfErrorWithMessage(err, "Unable to create directory")
		}
		err = os.Chmod(path+PathSeparator()+"submission", OwnerPermRw)
		if err != nil {
			CheckIfErrorWithMessage(err, "Unable to update directory mode")
		}
	}
}

func HandleTarStream(reader io.ReadCloser, destination string) {
	tr := tar.NewReader(reader)
	if tr != nil {
		err := processTarStream(tr, destination)
		if err != nil {
			log.Print(err)
		}
	} else {
		log.Printf("Unable to create image tar reader")
	}
}

func processTarStream(tr *tar.Reader, destination string) error {

	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("Unable to extract container: %v\n", err)
		}

		hdrInfo := hdr.FileInfo()

		temp := strings.TrimPrefix(hdr.Name, DockerTarPrefix)

		dstpath := path.Join(destination, temp)

		if runtime.GOOS == "windows" {
			dstpath = strings.Replace(dstpath, "/", "\\", -1)
		}

		// Overriding permissions to allow writing content
		mode := hdrInfo.Mode() | OwnerPermRw

		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(dstpath, mode); err != nil {
				if !os.IsExist(err) {
					return fmt.Errorf("Unable to create directory: %v", err)
				}
				err = os.Chmod(dstpath, mode)
				if err != nil {
					return fmt.Errorf("Unable to update directory mode: %v", err)
				}
			}
		case tar.TypeReg, tar.TypeRegA:
			file, err := os.OpenFile(dstpath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mode)
			if err != nil {
				return fmt.Errorf("Unable to create file: %v", err)
			}
			if _, err := io.Copy(file, tr); err != nil {
				err := file.Close()
				CheckIfError(err)

				return fmt.Errorf("Unable to write into file: %v", err)
			}
			err = file.Close()
			CheckIfError(err)

		case tar.TypeSymlink:
			if err := os.Symlink(hdr.Linkname, dstpath); err != nil {
				return fmt.Errorf("Unable to create symlink: %v\n", err)
			}
		case tar.TypeLink:
			target := path.Join(destination, strings.TrimPrefix(hdr.Linkname, DockerTarPrefix))
			if err := os.Link(target, dstpath); err != nil {
				return fmt.Errorf("Unable to create link: %v\n", err)
			}

		default:
			// For now we're skipping anything else. Special device files and
			// symlinks are not needed or anyway probably incorrect.
		}

		// maintaining access and modification time in best effort fashion
		err = os.Chtimes(dstpath, hdr.AccessTime, hdr.ModTime)
		CheckIfError(err)
	}
}

func Tar(source, target string) error {
	filename := filepath.Base(source)
	target = filepath.Join(target, fmt.Sprintf("%s.tar", filename))
	tarfile, err := os.Create(target)
	CheckIfError(err)

	defer func(tarfile *os.File) {
		err := tarfile.Close()
		CheckIfError(err)
	}(tarfile)

	tarball := tar.NewWriter(tarfile)
	defer func(tarball *tar.Writer) {
		err := tarball.Close()
		CheckIfError(err)
	}(tarball)

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	return filepath.Walk(source,
		func(path string, info os.FileInfo, err error) error {
			CheckIfError(err)

			header, err := tar.FileInfoHeader(info, info.Name())
			CheckIfError(err)

			if baseDir != "" {
				header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
			}

			if err := tarball.WriteHeader(header); err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			file, err := os.Open(path)
			CheckIfError(err)

			defer func(file *os.File) {
				err := file.Close()
				CheckIfError(err)
			}(file)

			_, err = io.Copy(tarball, file)
			return err
		})
}

func Untar(tarball, target string) error {
	reader, err := os.Open(tarball)
	CheckIfError(err)

	defer func(reader *os.File) {
		err := reader.Close()
		CheckIfError(err)
	}(reader)
	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		path := filepath.Join(target, header.Name)
		info := header.FileInfo()
		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			continue
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		CheckIfError(err)

		defer func(file *os.File) {
			err := file.Close()
			CheckIfError(err)
		}(file)

		_, err = io.Copy(file, tarReader)
		CheckIfError(err)
	}
	return nil
}

func Gzip(source, target string) error {
	reader, err := os.Open(source)
	CheckIfError(err)

	filename := filepath.Base(source)
	target = filepath.Join(target, fmt.Sprintf("%s.gz", filename))
	writer, err := os.Create(target)
	CheckIfError(err)

	defer func(writer *os.File) {
		err := writer.Close()
		CheckIfError(err)
	}(writer)

	archiver := gzip.NewWriter(writer)
	archiver.Name = filename
	defer func(archiver *gzip.Writer) {
		err := archiver.Close()
		CheckIfError(err)
	}(archiver)

	_, err = io.Copy(archiver, reader)
	return err
}

func UnGzip(source, target string) error {
	reader, err := os.Open(source)
	if err != nil {
		return err
	}
	defer func(reader *os.File) {
		err := reader.Close()
		CheckIfError(err)
	}(reader)

	archive, err := gzip.NewReader(reader)
	CheckIfError(err)

	defer func(archive *gzip.Reader) {
		err := archive.Close()
		CheckIfError(err)
	}(archive)

	target = filepath.Join(target, archive.Name)
	writer, err := os.Create(target)
	CheckIfError(err)

	defer func(writer *os.File) {
		err := writer.Close()
		CheckIfError(err)
	}(writer)

	_, err = io.Copy(writer, archive)
	return err
}

func ZipSource(source, target string) error {
	// 1. Create a ZIP file and zip.Writer
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		CheckIfError(err)
	}(f)

	writer := zip.NewWriter(f)
	defer func(writer *zip.Writer) {
		err := writer.Close()
		CheckIfError(err)
	}(writer)

	// 2. Go through all the files of the source
	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 3. Create a local file header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// set compression
		header.Method = zip.Deflate

		// 4. Set relative path of a file as the header name
		header.Name, err = filepath.Rel(filepath.Dir(source), path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			header.Name += "/"
		}

		// 5. Create writer for the file header and save content of the file
		headerWriter, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(headerWriter, f)
		return err
	})
}

// ZipFiles compresses one or many files into a single zip archive file.
// Param 1: filename is the output zip file's name.
// Param 2: files is a list of files to add to the zip.
func ZipFiles(filename string, files []string) error {

	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// Add files to zip
	for _, file := range files {
		if err = AddFileToZip(zipWriter, file); err != nil {
			return err
		}
	}
	return nil
}

func AddFileToZip(zipWriter *zip.Writer, filename string) error {

	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	// Using FileInfoHeader() above only uses the basename of the file. If we want
	// to preserve the folder structure we can overwrite this with the full path.
	header.Name = filename

	// Change to deflate to gain better compression
	// see http://golang.org/pkg/archive/zip/#pkg-constants
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}

func Zipit(source, target string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	//gw := gzip.NewWriter(zipfile)
	//defer gw.Close()
	//archive := tar.NewWriter(gw)
	//defer archive.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = "" //filepath.Base(source)
	}

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info) //tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}

		//if baseDir != "" {
		header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		//}
		if strings.HasPrefix(header.Name, "\\") {
			header.Name = header.Name[1:]
		}

		if info.IsDir() {
			header.Name += "/" //PathSeparator()
		} // else {
		//	header.Method = zip.Deflate
		//}

		header.Name = strings.Replace(header.Name, "\\", "/", -1)

		writer, err := archive.CreateHeader(header)
		//err = archive.WriteHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		//_, err = io.Copy(archive, file)
		_, err = io.Copy(writer, file)
		return err
	})

	return err
}

func GetAndCreateTGZArchive(source, target string) {

	zipfile, err := os.Create(target)
	CheckIfError(err)

	defer func(zipfile *os.File) {
		err := zipfile.Close()
		CheckIfError(err)
	}(zipfile)

	var archiveFiles []string

	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		archiveFiles = append(archiveFiles, path) //strings.TrimPrefix(path, source))

		return nil
	})
	CheckIfError(err)

	err = createTGZArchive(archiveFiles, zipfile)
	CheckIfError(err)

}

func createTGZArchive(files []string, buf io.Writer) error {
	// Create new Writers for gzip and tar
	// These writers are chained. Writing to the tar writer will
	// write to the gzip writer which in turn will write to
	// the "buf" writer
	gw := gzip.NewWriter(buf)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	// Iterate over files and add them to the tar archive
	for _, file := range files {
		err := addToTGZArchive(tw, file)
		if err != nil {
			return err
		}
	}

	return nil
}

func addToTGZArchive(tw *tar.Writer, filename string) error {
	// Open the file which will be written into the archive
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get FileInfo about our file providing file size, mode, etc.
	info, err := file.Stat()
	if err != nil {
		return err
	}

	// Create a tar Header from the FileInfo data
	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}

	// Use full path as name (FileInfoHeader only takes the basename)
	// If we don't do this the directory strucuture would
	// not be preserved
	// https://golang.org/src/archive/tar/common.go?#L626
	header.Name = filename

	// Write file header to the tar archive
	err = tw.WriteHeader(header)
	if err != nil {
		return err
	}

	// Copy file content to tar archive
	_, err = io.Copy(tw, file)
	if err != nil {
		return err
	}

	return nil
}

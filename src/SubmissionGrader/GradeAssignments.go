package main

import (
	// "context"
	"fmt"
	"github.com/go-git/go-git/v5"
	//"github.com/go-git/go-git/storage/memory"
	memory "github.com/go-git/go-git/v5/storage/memory"
	"os"
	"os/exec"
	// "runtime"
)

func GetEnvVar(envName string) string {
	if envName != "" {
		ret := os.Getenv("envName")
		return ret
	}
	fmt.Printf(envName + "environment variable was not found.")
	return "error"
}

var (
	//secret token
	githubToken = GetEnvVar("GITHUBTOKEN")

	//env variables
	language                  = GetEnvVar("LANGUAGE")
	repoName                  = GetEnvVar("REPOFULLNAME")
	teacherUnitTestsEnabled   = GetEnvVar("TEACHER_UNIT_TESTS")
	courseType                = GetEnvVar("COURSE_TYPE")
	assignmentName            = GetEnvVar("ASSIGNMENT_NAME")
	gradeDocsEnabled          = GetEnvVar("GRADE_DOCS")
	studentTestsEnabled       = GetEnvVar("STUDENT_TESTS_ENABLED")
	nonCodeSubmissionsEnabled = GetEnvVar("NONCODE_SUBMISSIONS")
	courseID                  = GetEnvVar("COURSE_ID")
	semesterID                = GetEnvVar("SEMESTER_ID")
	studentUserName           = GetEnvVar("SUDENT_USERNAME")

	// copyToPath = "/tmp/AutoGrader/$CourseID-$SemesterID/$AssignmentName/$StudentUserName"
	copyToPath = "/tmp/AutoGrader/" + courseID + "-" + semesterID + "/" + assignmentName + "/" + studentUserName
)

func main() {

	// done: start with environment variables
	// done: switch statement
	// TODO: one language at a time
	// TODO: java>c++>duplicate
	// TODO: sigterm handling on grader level

	switch language {
	case "java":
		{
			cmd := exec.Command("git", "clone", fmt.Sprintf("https://%v@github.com/%v.git", common.GithubAPIKey["token"], *repo.FullName))
			cmd.Dir = "/opt/gradle"
			err := cmd.Run()
			return err

			if teacherUnitTestsEnabled == "TRUE" {

				cmd := exec.Command("cp", "--recursive", oldDir, newDir)
				cmd.Run()
				//git.Clone(memory.NewStorage(), nil, &git.CloneOptions{URL: "https://" + githubToken + "@github.com/IUS-CS/AutoGraderConfig.git"})
				//common.RunDockerCommand(cli, containerName,
				//	[]string{"git", "clone", fmt.Sprintf("https://%v@github.com/%v.git", common.GithubAPIKey["token"], "IUS-CS/AutoGraderConfig")},
				//	"/tmp", verbose)
			}
			if gradeDocsEnabled == "TRUE" {
				fmt.Println("Grading java")
			}

		}
	case "c++":
		{
			git.Clone(memory.NewStorage(), nil, &git.CloneOptions{URL: "https://" + githubToken + "@github.com/IUS-CS/" + repoName + ".git"})
			if teacherUnitTestsEnabled == "TRUE" {
				git.Clone(memory.NewStorage(), nil, &git.CloneOptions{URL: "https://" + githubToken + "@github.com/IUS-CS/AutoGraderConfig.git"})
			}
			if gradeDocsEnabled == "TRUE" {
				fmt.Println("Grading c++")
			}
		}

	case "python":
		{

		}
	default:
		{
			//log and deal with unsupported language
			fmt.Println(language + " is not supported. Nothing has been graded.")
		}
	}
	//	ctx := context.Background()

	// 	verbose := false

	// 	containerName := *repo.Name + "-Container"

	// 	cli, err := client.NewClientWithOpts(client.FromEnv)
	// 	common.CheckIfError(err)

	// 	// Create and Start Container while destroying an existing container if it exists.
	// 	common.ForceStartContainer(cli, containerName, verbose)

	// 	var CopyToPath string
	// 	if runtime.GOOS == "windows" {
	// 		CopyToPath = fmt.Sprintf("%s\\AutoGrader\\%s\\%s\\%s", common.TempPath(), assignment.CourseIDSemID, assignment.AssignmentName, common.ParseUsername(*repo.HTMLURL, assignment))
	// 	} else {
	// 		CopyToPath = fmt.Sprintf("%s/AutoGrader/%s/%s/%s", common.TempPath(), assignment.CourseIDSemID, assignment.AssignmentName, common.ParseUsername(*repo.HTMLURL, assignment))
	// 	}

	// 	// See if previous zip folder exists, if so delete
	// 	if _, err := os.Stat(CopyToPath); !os.IsNotExist(err) {
	// 		err := os.RemoveAll(CopyToPath)
	// 		common.CheckIfError(err)
	// 	}

	// 	if assignment.AssignmentConfig.Language == "java" {

	// 		common.RunDockerCommand(cli, containerName,
	// 			[]string{"git", "clone", fmt.Sprintf("https://%v@github.com/%v.git", common.GithubAPIKey["token"], *repo.FullName)},
	// 			"/opt/gradle", verbose)

	// 		if assignment.AssignmentConfig.TeacherUnitTests {
	// 			//Clone Config Repo
	// 			// Could also use //err := cli.CopyToContainer(ctx, containerName,constants.TempPath(), , types.CopyToContainerOptions{AllowOverwriteDirWithFile: false,CopyUIDGID: false} )
	// 			// but I did not know how to use a tar reader to pass it in to the container
	// 			// This would also prevent parallelization in the future so we when with each one downloading itself
	// 			common.RunDockerCommand(cli, containerName,
	// 				[]string{"git", "clone", fmt.Sprintf("https://%v@github.com/%v.git", common.GithubAPIKey["token"], "IUS-CS/AutoGraderConfig")},
	// 				"/tmp", verbose)

	// 			// Copy tests to folder
	// 			common.RunDockerCommand(cli, containerName,
	// 				[]string{"cp", "-r", ".", fmt.Sprintf("/opt/gradle/%v/src/test/java", *repo.Name)},
	// 				fmt.Sprintf("/tmp/AutoGraderConfig/src/%v/%v/current/tests", assignment.CourseType, assignment.AssignmentName), verbose)
	// 		}

	// 		if assignment.AssignmentConfig.GradeDocs {
	// 			common.RunDockerCommand(cli, containerName,
	// 				[]string{"gradle", "javadoc"},
	// 				fmt.Sprintf("/opt/gradle/%v", *repo.Name), verbose)

	// 			tarStream, _, err := cli.CopyFromContainer(ctx, containerName, fmt.Sprintf("/opt/gradle/%v/build/docs/javadoc", *repo.Name))
	// 			//common.CheckIfError(err)

	// 			if err == nil {
	// 				common.HandleTarStream(tarStream, CopyToPath)

	// 				err = os.Rename(CopyToPath+common.PathSeparator()+"javadoc", CopyToPath+common.PathSeparator()+"docs")
	// 				common.CheckIfError(err)
	// 			}
	// 		}

	// 		if assignment.AssignmentConfig.StudentTestsEnabled || assignment.AssignmentConfig.TeacherUnitTests {

	// 			common.RunDockerCommand(cli, containerName,
	// 				[]string{"gradle", "test"},
	// 				fmt.Sprintf("/opt/gradle/%v", *repo.Name), verbose)

	// 			tarStream, _, err := cli.CopyFromContainer(ctx, containerName, fmt.Sprintf("/opt/gradle/%v/build/test-results/test", *repo.Name))
	// 			//common.CheckIfError(err)

	// 			if err == nil {
	// 				common.HandleTarStream(tarStream, CopyToPath)

	// 				err = os.Rename(CopyToPath+common.PathSeparator()+"test", CopyToPath+common.PathSeparator()+"tests")
	// 				common.CheckIfError(err)
	// 			}

	// 		}

	// 		if assignment.AssignmentConfig.NonCodeSubmissions {

	// 			tarStream, _, err := cli.CopyFromContainer(ctx, containerName, fmt.Sprintf("/opt/gradle/%v/submission", *repo.Name))

	// 			if err == nil {
	// 				common.HandleTarStream(tarStream, CopyToPath)
	// 			}
	// 		}

	// 	} else if assignment.AssignmentConfig.Language == "c++" {

	// 		common.RunDockerCommand(cli, containerName,
	// 			[]string{"git", "clone", fmt.Sprintf("https://%v@github.com/%v.git", common.GithubAPIKey["token"], *repo.FullName)},
	// 			"/opt/cmake", verbose)

	// 		if assignment.AssignmentConfig.TeacherUnitTests {
	// 			//Clone Config Repo
	// 			// Could also use //err := cli.CopyToContainer(ctx, containerName,constants.TempPath(), , types.CopyToContainerOptions{AllowOverwriteDirWithFile: false,CopyUIDGID: false} )
	// 			// but I did not know how to use a tar reader to pass it in to the container
	// 			// This would also prevent parallelization in the future so we when with each one downloading itself
	// 			common.RunDockerCommand(cli, containerName,
	// 				[]string{"git", "clone", fmt.Sprintf("https://%v@github.com/%v.git", common.GithubAPIKey["token"], "IUS-CS/AutoGraderConfig")},
	// 				"/tmp", verbose)

	// 			// Copy tests to folder
	// 			common.RunDockerCommand(cli, containerName,
	// 				[]string{"cp", "-r", ".", fmt.Sprintf("/opt/cmake/%v/test", *repo.Name)},
	// 				fmt.Sprintf("/tmp/AutoGraderConfig/src/%v/%v/current/tests", assignment.CourseType, assignment.AssignmentName), verbose)
	// 		}

	// 		common.RunDockerCommand(cli, containerName,
	// 			[]string{"cmake", "-S", fmt.Sprintf("../%v", *repo.Name), "-B", "."},
	// 			"/opt/cmake/cmake-build-debug", verbose)

	// 		common.RunDockerCommand(cli, containerName,
	// 			[]string{"cmake", "--build", ".", "--target", "all"},
	// 			"/opt/cmake/cmake-build-debug", verbose)

	// 		// Build has to be ran and it covers doc generation too
	// 		// Only copy over if we are grading it
	// 		if assignment.AssignmentConfig.GradeDocs {
	// 			tarStream, _, err := cli.CopyFromContainer(ctx, containerName, "/opt/cmake/cmake-build-debug/docs/doc_doxygen/html")

	// 			if err == nil {
	// 				common.HandleTarStream(tarStream, CopyToPath)

	// 				err = os.Rename(CopyToPath+common.PathSeparator()+"html", CopyToPath+common.PathSeparator()+"docs")
	// 				common.CheckIfError(err)
	// 			}

	// 		}

	// 		if assignment.AssignmentConfig.StudentTestsEnabled || assignment.AssignmentConfig.TeacherUnitTests {

	// 			common.RunDockerCommand(cli, containerName,
	// 				[]string{"ctest", ".", "--output-junit", "./test/reports/test_output.xml"},
	// 				"/opt/cmake/cmake-build-debug", verbose)

	// 			tarStream, _, err := cli.CopyFromContainer(ctx, containerName, "/opt/cmake/cmake-build-debug/test/reports")

	// 			if err == nil {
	// 				common.HandleTarStream(tarStream, CopyToPath)

	// 				err = os.Rename(CopyToPath+common.PathSeparator()+"reports", CopyToPath+common.PathSeparator()+"tests")
	// 				common.CheckIfError(err)
	// 			}
	// 		}

	// 		if assignment.AssignmentConfig.NonCodeSubmissions {

	// 			tarStream, _, err := cli.CopyFromContainer(ctx, containerName, fmt.Sprintf("/opt/gradle/%v/submission", *repo.Name))
	// 			if err == nil {
	// 				common.HandleTarStream(tarStream, CopyToPath)
	// 			}
	// 		}

	// 	} else {
	// 		// TODO log unsupported language used and
	// 		common.Warning("%v is not supported. Nothing has been graded.", assignment.AssignmentConfig.Language)
	// 	}

	// 	common.CheckIfError(common.StopAndRemoveContainer(cli, containerName))

	//SIGTERM STUFF COPIED FROM MY LAST PROJECT
	// fmt.Println("Bot is now running. Press CTRL-C to exit.")
	// sc := make(chan os.Signal, 1)
	// signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	// <-sc
}

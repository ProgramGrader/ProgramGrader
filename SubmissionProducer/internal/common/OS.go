package common

import (
	"fmt"
	"os/exec"
	"syscall"
	"unsafe"
)

// Needs work this does not currently work
func shutdownLinux() {

	cmd := exec.Command("shutdown now")

	err := cmd.Run()

	if err != nil {
		Warning("unable to shutdown server, err: %v", err)
	}

}

func rebootLinux() {

	cmd := exec.Command("reboot")

	err := cmd.Run()

	if err != nil {
		Warning("unable to shutdown server, err: %v", err)
	}

}

func rebootWindows() error {

	user32 := syscall.MustLoadDLL("user32")
	defer user32.Release()

	kernel32 := syscall.MustLoadDLL("kernel32")
	defer user32.Release()

	advapi32 := syscall.MustLoadDLL("advapi32")
	defer advapi32.Release()

	ExitWindowsEx := user32.MustFindProc("ExitWindowsEx")
	GetCurrentProcess := kernel32.MustFindProc("GetCurrentProcess")
	GetLastError := kernel32.MustFindProc("GetLastError")
	OpenProdcessToken := advapi32.MustFindProc("OpenProcessToken")
	LookupPrivilegeValue := advapi32.MustFindProc("LookupPrivilegeValueW")
	AdjustTokenPrivileges := advapi32.MustFindProc("AdjustTokenPrivileges")

	currentProcess, _, _ := GetCurrentProcess.Call()

	const tokenAdjustPrivileges = 0x0020
	const tokenQuery = 0x0008
	var hToken uintptr

	result, _, err := OpenProdcessToken.Call(currentProcess, tokenAdjustPrivileges|tokenQuery, uintptr(unsafe.Pointer(&hToken)))
	if result != 1 {
		fmt.Println("OpenProcessToken(): ", result, " err: ", err)
		return err
	}
	//fmt.Println("hToken: ", hToken)

	const SeShutdownName = "SeShutdownPrivilege"

	type Luid struct {
		lowPart  uint32 // DWORD
		highPart int32  // long
	}
	type LuidAndAttributes struct {
		luid       Luid   // LUID
		attributes uint32 // DWORD
	}

	type TokenPrivileges struct {
		privilegeCount uint32 // DWORD
		privileges     [1]LuidAndAttributes
	}

	var tkp TokenPrivileges

	result, _, err = LookupPrivilegeValue.Call(uintptr(0), uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(SeShutdownName))), uintptr(unsafe.Pointer(&(tkp.privileges[0].luid))))
	if result != 1 {
		fmt.Println("LookupPrivilegeValue(): ", result, " err: ", err)
		return err
	}
	//fmt.Println("LookupPrivilegeValue luid: ", tkp.privileges[0].luid)

	const SePrivilegeEnabled uint32 = 0x00000002

	tkp.privilegeCount = 1
	tkp.privileges[0].attributes = SePrivilegeEnabled

	result, _, err = AdjustTokenPrivileges.Call(hToken, 0, uintptr(unsafe.Pointer(&tkp)), 0, uintptr(0), 0)
	if result != 1 {
		fmt.Println("AdjustTokenPrivileges() ", result, " err: ", err)
		return err
	}

	result, _, _ = GetLastError.Call()
	if result != 0 {
		fmt.Println("GetLastError() ", result)
		return err
	}

	const ewxForceIfHung = 0x00000010
	const ewxReboot = 0x00000002
	const shutdownReasonMajorSoftware = 0x00030000

	result, _, err = ExitWindowsEx.Call(ewxReboot|ewxForceIfHung, shutdownReasonMajorSoftware)
	if result != 1 {
		fmt.Println("Failed to initiate reboot:", err)
		return err
	}

	return nil
}

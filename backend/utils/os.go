package utils

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"syscall"
	"unsafe"

	"../models"
)

func GetFPRegs(pid int, data *models.FPRegs) error {
	_, _, errno := syscall.RawSyscall6(uintptr(syscall.SYS_PTRACE),
		uintptr(syscall.PTRACE_GETFPREGS),
		uintptr(pid),
		uintptr(0),
		uintptr(unsafe.Pointer(data)),
		0,
		0)

	var err error
	if errno != 0 {
		err = errno
		return err
	}
	return nil
}

func prLimit(pid int, limit uintptr, rlimit *syscall.Rlimit) error {
	_, _, errno := syscall.RawSyscall6(syscall.SYS_PRLIMIT64,
		uintptr(pid),
		limit,
		uintptr(unsafe.Pointer(rlimit)),
		0, 0, 0)
	var err error
	if errno != 0 {
		err = errno
		return err
	}
	return nil
}

func LimitFileSize(pid int, maxSize uint64) {
	var rlimit syscall.Rlimit

	rlimit.Cur = maxSize
	rlimit.Max = maxSize
	prLimit(pid, syscall.RLIMIT_FSIZE, &rlimit)
}

func LimitCPUTime(pid int, maxTime uint64) {
	var rlimit syscall.Rlimit

	rlimit.Cur = maxTime
	rlimit.Max = maxTime
	prLimit(pid, syscall.RLIMIT_CPU, &rlimit)
}

func KillProcess(pid int, err string) models.ResponseObj {
	fmt.Println("Killing Process")
	var ws syscall.WaitStatus
	_, _, killErr := syscall.RawSyscall6(syscall.SYS_KILL,
		uintptr(pid),
		uintptr(syscall.SIGKILL),
		0, 0, 0, 0)
	syscall.Wait4(pid, &ws, syscall.WALL, nil)
	if PidExists(pid) {
		return models.ResponseObj{ConsoleOut: err + "\nCould not kill process: " + strconv.Itoa(pid) + "\nError: " + killErr.Error()}

	}
	fmt.Println("Process killed succesfully.")

	return models.ResponseObj{ConsoleOut: err + "\nProcess killed succesfully."}
}

func PidExists(pid int) bool {
	_, err := ioutil.ReadFile("/proc/" + strconv.Itoa(pid) + "/status")
	return err == nil
}

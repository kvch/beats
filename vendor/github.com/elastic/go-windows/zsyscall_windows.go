// MACHINE GENERATED BY 'go generate' COMMAND; DO NOT EDIT

package windows

import (
	"syscall"
	"unsafe"
)

var _ unsafe.Pointer

// Do the interface allocations only once for common
// Errno values.
const (
	errnoERROR_IO_PENDING = 997
)

var (
	errERROR_IO_PENDING error = syscall.Errno(errnoERROR_IO_PENDING)
)

// errnoErr returns common boxed Errno values, to prevent
// allocations at runtime.
func errnoErr(e syscall.Errno) error {
	switch e {
	case 0:
		return nil
	case errnoERROR_IO_PENDING:
		return errERROR_IO_PENDING
	}
	// TODO: add more here, after collecting data on the common
	// error values see on Windows. (perhaps when running
	// all.bat?)
	return e
}

var (
	modkernel32 = syscall.NewLazyDLL("kernel32.dll")
	modversion  = syscall.NewLazyDLL("version.dll")
	modpsapi    = syscall.NewLazyDLL("psapi.dll")

	procGetNativeSystemInfo     = modkernel32.NewProc("GetNativeSystemInfo")
	procGetTickCount64          = modkernel32.NewProc("GetTickCount64")
	procGetSystemTimes          = modkernel32.NewProc("GetSystemTimes")
	procGlobalMemoryStatusEx    = modkernel32.NewProc("GlobalMemoryStatusEx")
	procGetFileVersionInfoW     = modversion.NewProc("GetFileVersionInfoW")
	procGetFileVersionInfoSizeW = modversion.NewProc("GetFileVersionInfoSizeW")
	procVerQueryValueW          = modversion.NewProc("VerQueryValueW")
	procGetProcessMemoryInfo    = modpsapi.NewProc("GetProcessMemoryInfo")
)

func _GetNativeSystemInfo(systemInfo *SystemInfo) {
	syscall.Syscall(procGetNativeSystemInfo.Addr(), 1, uintptr(unsafe.Pointer(systemInfo)), 0, 0)
	return
}

func _GetTickCount64() (millis uint64, err error) {
	r0, _, e1 := syscall.Syscall(procGetTickCount64.Addr(), 0, 0, 0, 0)
	millis = uint64(r0)
	if millis == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func _GetSystemTimes(idleTime *syscall.Filetime, kernelTime *syscall.Filetime, userTime *syscall.Filetime) (err error) {
	r1, _, e1 := syscall.Syscall(procGetSystemTimes.Addr(), 3, uintptr(unsafe.Pointer(idleTime)), uintptr(unsafe.Pointer(kernelTime)), uintptr(unsafe.Pointer(userTime)))
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func _GlobalMemoryStatusEx(buffer *MemoryStatusEx) (err error) {
	r1, _, e1 := syscall.Syscall(procGlobalMemoryStatusEx.Addr(), 1, uintptr(unsafe.Pointer(buffer)), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func _GetProcessHandleCount(handle syscall.Handle, pdwHandleCount *uint32) (err error) {
	r1, _, e1 := syscall.Syscall(procGetProcessHandleCount.Addr(), 2, uintptr(handle), uintptr(unsafe.Pointer(pdwHandleCount)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func _GetFileVersionInfo(filename string, reserved uint32, dataLen uint32, data *byte) (success bool, err error) {
	var _p0 *uint16
	_p0, err = syscall.UTF16PtrFromString(filename)
	if err != nil {
		return
	}
	return __GetFileVersionInfo(_p0, reserved, dataLen, data)
}

func __GetFileVersionInfo(filename *uint16, reserved uint32, dataLen uint32, data *byte) (success bool, err error) {
	r0, _, e1 := syscall.Syscall6(procGetFileVersionInfoW.Addr(), 4, uintptr(unsafe.Pointer(filename)), uintptr(reserved), uintptr(dataLen), uintptr(unsafe.Pointer(data)), 0, 0)
	success = r0 != 0
	if !success {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func _GetFileVersionInfoSize(filename string, handle uintptr) (size uint32, err error) {
	var _p0 *uint16
	_p0, err = syscall.UTF16PtrFromString(filename)
	if err != nil {
		return
	}
	return __GetFileVersionInfoSize(_p0, handle)
}

func __GetFileVersionInfoSize(filename *uint16, handle uintptr) (size uint32, err error) {
	r0, _, e1 := syscall.Syscall(procGetFileVersionInfoSizeW.Addr(), 2, uintptr(unsafe.Pointer(filename)), uintptr(handle), 0)
	size = uint32(r0)
	if size == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func _VerQueryValueW(data *byte, subBlock string, pBuffer *uintptr, len *uint32) (success bool, err error) {
	var _p0 *uint16
	_p0, err = syscall.UTF16PtrFromString(subBlock)
	if err != nil {
		return
	}
	return __VerQueryValueW(data, _p0, pBuffer, len)
}

func __VerQueryValueW(data *byte, subBlock *uint16, pBuffer *uintptr, len *uint32) (success bool, err error) {
	r0, _, e1 := syscall.Syscall6(procVerQueryValueW.Addr(), 4, uintptr(unsafe.Pointer(data)), uintptr(unsafe.Pointer(subBlock)), uintptr(unsafe.Pointer(pBuffer)), uintptr(unsafe.Pointer(len)), 0, 0)
	success = r0 != 0
	if !success {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func _GetProcessMemoryInfo(handle syscall.Handle, psmemCounters *ProcessMemoryCountersEx, cb uint32) (err error) {
	r1, _, e1 := syscall.Syscall(procGetProcessMemoryInfo.Addr(), 3, uintptr(handle), uintptr(unsafe.Pointer(psmemCounters)), uintptr(cb))
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

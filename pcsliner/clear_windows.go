package pcsliner

import (
	"syscall"
	"unsafe"
)

const (
	std_output_handle = uint32(-11 & 0xFFFFFFFF)
)

var (
	kernel32 = syscall.NewLazyDLL("kernel32.dll")

	procGetStdHandle               = kernel32.NewProc("GetStdHandle")
	procSetConsoleCursorPosition   = kernel32.NewProc("SetConsoleCursorPosition")
	procGetConsoleScreenBufferInfo = kernel32.NewProc("GetConsoleScreenBufferInfo")
	procFillConsoleOutputCharacter = kernel32.NewProc("FillConsoleOutputCharacterW")
)

type (
	coord struct {
		x, y int16
	}
	smallRect struct {
		left, top, right, bottom int16
	}
	consoleScreenBufferInfo struct {
		dwSize              coord
		dwCursorPosition    coord
		wAttributes         int16
		srWindow            smallRect
		dwMaximumWindowSize coord
	}
)

// ClearScreen 清空屏幕
func (pl *PCSLiner) ClearScreen() {
	ClearScreen()
}

// ClearScreen 清空屏幕
func ClearScreen() {
	out, _, _ := procGetStdHandle.Call(uintptr(std_output_handle))
	hOut := syscall.Handle(out)

	var sbi consoleScreenBufferInfo
	procGetConsoleScreenBufferInfo.Call(uintptr(hOut), uintptr(unsafe.Pointer(&sbi)))

	var numWritten uint32
	procFillConsoleOutputCharacter.Call(uintptr(hOut), uintptr(' '),
		uintptr(sbi.dwSize.x)*uintptr(sbi.dwSize.y),
		0,
		uintptr(unsafe.Pointer(&numWritten)))
	procSetConsoleCursorPosition.Call(uintptr(hOut), 0)
}

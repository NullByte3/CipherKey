package main

import (
	"github.com/gonutz/w32/v2"
	"github.com/lxn/win"
	"strings"
	"syscall"
	"unsafe"
)

// https://stackoverflow.com/questions/42500570/how-to-hide-command-prompt-window-when-using-exec-in-golang
func hideTerminal() {
	console := w32.GetConsoleWindow()
	if console == 0 {
		return
	}
	_, consoleProcID := w32.GetWindowThreadProcessId(console)
	if w32.GetCurrentProcessId() == consoleProcID {
		w32.ShowWindowAsync(console, w32.SW_HIDE)
	}
}
func getActiveWindowTitle() string {
	hwnd, _, _ := procGetForegroundWindow.Call()
	if hwnd == 0 {
		return "None"
	}

	textLength, _, _ := procGetWindowTextLength.Call(hwnd)
	buff := make([]uint16, textLength+1)
	procGetWindowText.Call(hwnd, uintptr(unsafe.Pointer(&buff[0])), textLength+1)

	var windowTitle = syscall.UTF16ToString(buff)

	if strings.Contains(windowTitle, " — Mozilla Firefox") {
		return "SITE: " + strings.Replace(windowTitle, " — Mozilla Firefox", "", 1)
	} else if strings.Contains(windowTitle, " - Google Chrome") {
		return "SITE: " + strings.Replace(windowTitle, " - Google Chrome", "", 1)
	} else {
		return windowTitle
	}
}

func showMessage(title, message string) {
	win.MessageBox(0, syscall.StringToUTF16Ptr(message), syscall.StringToUTF16Ptr(title), win.MB_ICONINFORMATION)
}

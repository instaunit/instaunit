package term

import (
  "unsafe"
  "syscall"
)

// IsTerminal returns true if the given file descriptor is a terminal. Lifted from https://github.com/golang/crypto/blob/master/ssh/terminal/util.go#L30
func IsTerminal(fd int) bool {
  var termios syscall.Termios
  _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), ioctlReadTermios, uintptr(unsafe.Pointer(&termios)), 0, 0, 0)
  return err == 0
}

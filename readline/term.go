package readline

// copy pasted from https://mzunino.com.uy/til/2025/03/building-a-terminal-raw-mode-input-reader-in-go/
// i'll study it later

import (
	"syscall"
	"unsafe"
)

// termios holds the terminal attributes
type termios struct {
	Iflag  uint32
	Oflag  uint32
	Cflag  uint32
	Lflag  uint32
	Cc     [20]byte
	Ispeed uint32
	Ospeed uint32
}

// enableRawMode switches the terminal to raw mode and returns the original state
func enableRawMode() (*termios, error) {
	fd := int(syscall.Stdin)
	var oldState termios

	// Get the current terminal attributes
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(fd), uintptr(syscall.TCGETS),
		uintptr(unsafe.Pointer(&oldState)))
	if errno != 0 {
		return nil, errno
	}

	// Modify the attributes to enable raw mode
	newState := oldState
	// Disable canonical mode (ICANON) and echo (ECHO)
	newState.Lflag &^= syscall.ICANON | syscall.ECHO

	// Set the new terminal attributes
	_, _, errno = syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(fd), uintptr(syscall.TCSETS),
		uintptr(unsafe.Pointer(&newState)))
	if errno != 0 {
		return nil, errno
	}

	return &oldState, nil
}

// disableRawMode restores the terminal to its original state
func disableRawMode(oldState *termios) error {
	fd := int(syscall.Stdin)
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(fd), uintptr(syscall.TCSETS),
		uintptr(unsafe.Pointer(oldState)))
	if errno != 0 {
		return errno
	}
	return nil
}

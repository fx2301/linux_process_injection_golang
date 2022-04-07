package main

import (
	"fmt"
	"os"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

const fakeProcessName = "injectedprocess"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Sleeping for 100 seconds. Use --inject to process inject.")
		time.Sleep(time.Second * 100)
		os.Exit(0)
	}

	if os.Args[1] != "--inject" {
		fmt.Println("Use --inject to process inject.")
		os.Exit(1)
	}

	fmt.Println("Injecting into process...")

	const fileDescriptorName = ""

	// not using close-on-exec because it didn't work
	fd, err := unix.MemfdCreate(fileDescriptorName, 0)

	if err != nil {
		panic(fmt.Sprintf("MemfdCreate failed with: %q", err))
	}

	fmt.Printf("File descriptor is %d\n", fd)

	payload, err := os.ReadFile("./memfd_launch")

	if err != nil {
		panic(err)
	}

	written, err := unix.Write(fd, payload)

	if err != nil {
		panic(fmt.Sprintf("Write failed with: %q", err))
	}

	fdPath := fmt.Sprintf("/proc/self/fd/%d", fd)

	fmt.Printf("Bytes written to memfd %d at path %s\n", written, fdPath)
	// no x/sys/unix fork wrapper available
	ret, _, _ := syscall.Syscall(syscall.SYS_FORK, 0, 0, 0)
	switch ret {
	case 0:
		break
	default:
		os.Exit(0)
	}
	_, err = unix.Setsid()
	if err != nil {
		panic(err)
	}

	os.Chdir("/")
	f, e := os.OpenFile("/dev/null", os.O_RDWR, 0)
	if e == nil {
		fd := int(f.Fd())
		unix.Dup2(fd, int(os.Stdin.Fd()))
		unix.Dup2(fd, int(os.Stdout.Fd()))
		unix.Dup2(fd, int(os.Stderr.Fd()))
	}

	_ = syscall.Exec(fdPath, []string{fakeProcessName}, nil)
}

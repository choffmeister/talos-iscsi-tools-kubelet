package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"

	"golang.org/x/sys/unix"
)

const CLONE_NEWNS = 0x20000

func main() {
	pids, _ := getPids("/usr/local/bin/kubelet")
	fmt.Printf("+%v", pids)
	mountns, _ := getMountNamespaceHandleFromPid(os.Getpid())

	for _, pid := range pids {
		fmt.Printf("Patching kubelet with pid %d", pid)
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

		unix.setns(getMountNamespaceHandleFromPid(pid), CLONE_NEWNS)

		err := writeFile()
		if err != nil {
			panic(err)
		}

		unix.setns(mountns, CLONE_NEWNS)
	}

	return

}

func writeFile() error {
	content := `#!/bin/sh
	set -e
		
	iscsid_pid=$(pidof iscsid)
	nsenter --mount="/proc/${iscsid_pid}/ns/mnt" --net="/proc/${iscsid_pid}/ns/net" -- /usr/local/bin/iscsiadm "$@"
	`

	f, err := os.OpenFile("/sbin/iscsiadm", os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	ioutil.WriteFile("/sbin/iscsiadm", []byte(content), 0755)

	return nil
}

func getPids(path string) ([]int, error) {
	pids := make([]int, 0)
	re := regexp.MustCompile("/proc/([0-9]+)/exe")

	matches, err := filepath.Glob("/proc/[0-9][0-9]*/exe")
	if err != nil {
		return nil, err
	}

	for _, file := range matches {
		fmt.Println(file)
		target, err := os.Readlink(file)

		if err != nil {
			continue
		}

		if target == path {
			match := re.FindStringSubmatch(file)
			pid, err := strconv.Atoi(match[1])
			if err != nil {
				continue
			}

			pids = append(pids, pid)
		}
	}

	return pids, nil
}

func setMountNamespace(namespace int) error {
	return unix.Setns(namespace, CLONE_NEWNS)
}

func getMountNamespaceHandleFromPid(pid int) (int, error) {
	return getNamespaceHandleFromPath(fmt.Sprintf("/proc/$d/ns/mnt", pid))
}

func getNamespaceHandleFromPath(path string) (int, error) {
	fd, err := unix.Open(path, unix.O_RDONLY|unix.O_CLOEXEC, 0)
	if err != nil {
		return -1, err
	}
	return fd, nil
}


package main

import (
	"errors"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

// const (
// 	CLONE_NEWNS = 0x20000 /* New mount namespace */
// )

func main() {
	sigs := make(chan os.Signal, 1)
	logger, _ := zap.NewProduction()

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	wrapperPath := "/host/run/containerd/io.containerd.runtime.v2.task/system/kubelet/rootfs/sbin/iscsiadm"

	for {
		if _, err := os.Stat(wrapperPath); err == nil {
			logger.Info("Wrapper is already in place")
		} else if errors.Is(err, os.ErrNotExist) {
			logger.Info("Deploying wrapper")
			err = writeWrapper(wrapperPath)
			if err != nil {
				logger.Error("Failed to deploy wrapper", zap.Error(err))
			}
		} else {
			logger.Error("Unable to determine state of wrapper installation", zap.Error(err))
		}

		logger.Info("Sleeping")

		select {
		case <-sigs:
			logger.Info("Terminating")
			os.Exit(0)
		case <-time.After(15 * time.Second):
		}
	}
}

func writeWrapper(path string) error {
	content := `#!/bin/sh
set -e
	
iscsid_pid=$(pidof iscsid)
nsenter --mount="/proc/${iscsid_pid}/ns/mnt" --net="/proc/${iscsid_pid}/ns/net" -- /usr/local/sbin/iscsiadm "$@"
`

	return ioutil.WriteFile(path, []byte(content), 0755)
}

// func main() {
// 	for {
// 		pids, _ := getPids("/usr/local/bin/kubelet")

// 		if len(pids) == 0 {
// 			fmt.Println("No kubelets found. Sleeping...")
// 		} else {
// 			fmt.Printf("Found %d kubelet instances", len(pids))
// 			fmt.Printf("+%v", pids)
// 			mountns, _ := getNamespaceFromPid(os.Getpid())

// 			for _, pid := range pids {
// 				fmt.Printf("Patching kubelet with pid %d", pid)
// 				runtime.LockOSThread()
// 				defer runtime.UnlockOSThread()

// 				pid, err := getNamespaceFromPid(pid)
// 				if err != nil {
// 					panic(err)
// 				}
// 				unix.Setns(pid, CLONE_NEWNS)

// 				err = writeFile()
// 				if err != nil {
// 					panic(err)
// 				}

// 				unix.Setns(mountns, CLONE_NEWNS)
// 			}
// 		}

// 		time.Sleep(15 * time.Second)
// 	}
// }

// func getPids(path string) ([]int, error) {
// 	filepath.Walk("/proc", func(path string, info fs.FileInfo, err error) error {
// 		fmt.Println(path)
// 		return err
// 	})

// 	pids := make([]int, 0)
// 	re := regexp.MustCompile("/proc/([0-9]+)/exe")

// 	matches, err := filepath.Glob("/proc/[0-9][0-9]*/exe")
// 	if err != nil {
// 		return nil, err
// 	}

// 	fmt.Printf("Found %d matches\n", len(matches))

// 	for _, file := range matches {
// 		fmt.Println(file)
// 		target, err := os.Readlink(file)

// 		if err != nil {
// 			continue
// 		}

// 		fmt.Println(target)

// 		if target == path {
// 			match := re.FindStringSubmatch(file)
// 			pid, err := strconv.Atoi(match[1])
// 			if err != nil {
// 				continue
// 			}

// 			pids = append(pids, pid)
// 		}
// 	}

// 	return pids, nil
// }

// func setNamespace(namespace int) error {
// 	return unix.Setns(namespace, CLONE_NEWNS)
// }

// func getNamespaceFromPid(pid int) (int, error) {
// 	return getNamespaceFromPath(fmt.Sprintf("/proc/%d/ns/mnt", pid))
// }

// func getNamespaceFromPath(path string) (int, error) {
// 	fd, err := unix.Open(path, unix.O_RDONLY|unix.O_CLOEXEC, 0)
// 	if err != nil {
// 		return -1, err
// 	}
// 	return fd, nil
// }

package main

import (
	"os"
	"syscall"
	"time"
)

const (
	// timeBetweenPidMonitoringChecks는 프로세스가 실행 중인지 확인하는 데 사용되는 시간 간격입니다.
	timeBetweenPidMonitoringChecks = 2 * time.Second
)

// isProcessRunning은 주어진 pid를 가진 프로세스가 실행 중이면 true를 반환합니다
func isProcessRunning(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	err = process.Signal(syscall.Signal(0))
	return err == nil
}

// monitorPid는 서버 수명을 다른(클라이언트 앱) 프로세스와 동기화하는 데 사용됩니다.
func monitorPid(pid int) {
	go func() {
		for {
			// 프로세스가 실행 중이 아니면 프로그램을 종료합니다.
			if !isProcessRunning(pid) {
				println("Monitored process not found, exiting.")
				os.Exit(1)
			}
			// 프로세스가 실행 중이면 2초 후에 다시 확인합니다.
			time.Sleep(timeBetweenPidMonitoringChecks)
		}
	}()
}

func main() {
	println("Hello, World!")
}

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"time"
)

const (
	// Dead critical state at which reboot is recommended
	Dead = "dead"
	// ConnectionDead state at which commuication is lost
	ConnectionDead = "connection_dead"
	// NoAccessToPollingServer state at which polling is failing
	NoAccessToPollingServer = "no_access"
)

var (
	logger = log.New(os.Stdout, "polling", 0)
)

// PollingStatus object to return on polling
type PollingStatus struct {
	Status string `json:"status"`
}

func main() {
	chExecPolling := make(chan struct{})
	pollingTimer(chExecPolling)

	for {
		select {
		case <-chExecPolling:
			polling()
		}
	}
}

func pollingTimer(ch chan struct{}) {
	go func() {
		ticker := time.NewTicker(1 * time.Minute)

		for {
			select {
			case <-ticker.C:
				if isTimeToCheck(time.Now()) {
					ch <- struct{}{}
				}
			}
		}
	}()
}

func isTimeToCheck(time time.Time) bool {
	return true
}

func getScriptsDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Join(path.Dir(filename), "scripts")
}

func polling() {
	logger.Printf("exec polling %v", time.Now())

	status := getStatus()
	logger.Printf("polling status %s", status)

	if status == NoAccessToPollingServer || status == Dead {
		// reboot
		execPath := path.Join(getScriptsDir(), "reboot.sh")
		logger.Printf("exec reboot, path: %s", execPath)
		if err := exec.Command(execPath).Run(); err != nil {
			logger.Printf("script %s exec, %s", execPath, err)
		}

	} else if status == ConnectionDead {
		// restart ssh
		execPath := path.Join(getScriptsDir(), "restart_ssh_tunnel.sh")
		logger.Printf("exec %s", execPath)
		if err := exec.Command(execPath).Run(); err != nil {
			logger.Printf("script %s exec, %s", execPath, err)
		}
	}
}

func getStatus() string {
	res, err := http.Get("http://35.200.74.27/polling/status")
	if err != nil {
		logger.Print(err)
		return NoAccessToPollingServer
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Print(err)
		return NoAccessToPollingServer
	}

	s, err := parseResponse([]byte(body))
	if err != nil {
		return NoAccessToPollingServer
	}
	return s
}

func parseResponse(body []byte) (string, error) {
	var s PollingStatus
	err := json.Unmarshal(body, &s)
	if err != nil {
		logger.Print(err)
		return "", err
	}
	return s.Status, nil
}

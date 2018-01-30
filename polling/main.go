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
	Dead                    = "dead"
	ConnectionDead          = "connection_dead"
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
		ticker := time.NewTicker(10 * time.Second) // 1秒間隔のTicker

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

func polling() {
	logger.Printf("exec polling %v", time.Now())

	status := getStatus()
	logger.Printf("polling status %s", status)

	if status == NoAccessToPollingServer || status == Dead {
		// reboot
		logger.Print("exec reboot")

		_, filename, _, _ := runtime.Caller(1)
		shellPath := path.Join(path.Dir(filename), "scripts", "reboot.sh")
		logger.Printf("exec %s", shellPath)

		if err := exec.Command(shellPath).Run(); err != nil {
			logger.Printf("script exec, %s", err)
		}

	} else if status == ConnectionDead {
		logger.Printf("connection is dead")
		// restart ssh
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

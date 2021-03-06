package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo"
)

// StatusNormal status returned when normal
const StatusNormal = "normal"

// PollingStatus object to return on polling
type PollingStatus struct {
	Status string `json:"status"`
}

var (
	logger = log.New(os.Stdout, "[polling server]", 0)
)

func main() {
	e := echo.New()
	e.GET("/status", statusGetHandler)
	e.GET("/status-check", statusCheckHandler)
	e.POST("/status", statusSetHandler)
	e.Logger.Fatal(e.Start(":9000"))
}

func statusGetHandler(c echo.Context) error {
	status, err := readStatus()
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	// reset to normal status
	if err = writeStatus(StatusNormal); err != nil {
		return c.JSON(http.StatusInternalServerError, &PollingStatus{Status: status})
	}
	return c.JSON(http.StatusOK, &PollingStatus{Status: status})
}

// statusCheckHandler only returns current status but do not reset to StatusNormal
func statusCheckHandler(c echo.Context) error {
	status, err := readStatus()
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, &PollingStatus{Status: status})
}

func statusSetHandler(c echo.Context) error {
	var status PollingStatus
	if err := c.Bind(&status); err != nil {
		logger.Printf("on status set, bind status json to structure failed, %s", err)
		return c.String(http.StatusInternalServerError, fmt.Sprint(err))
	}

	if err := writeStatus(status.Status); err != nil {
		logger.Print(err)
		return c.String(http.StatusInternalServerError, fmt.Sprint(err))
	}
	return c.NoContent(http.StatusOK)
}

func readStatus() (string, error) {
	return redisGetStatus()
}

func writeStatus(status string) error {
	return redisSetStatus(status)
}

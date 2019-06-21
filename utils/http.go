package utils

import (
	"fmt"
	"net/http"

	"github.com/walk1ng/cf-push-stress/model"
)

func doHTTPVerifyApp(app model.App) (verified bool, err error) {
	route := fmt.Sprintf("http://%s.%s", app.Name, app.Domain)
	resp, err := http.Get(route)

	if err != nil {
		logger.Printf("[%s]:Send http request failed\n", app.Name)
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		logger.Printf("[%s]: Verify pass with code: %d\n", app.Name, resp.StatusCode)
		return true, nil
	}
	logger.Printf("[%s]: Verify failed with code: %d\n", app.Name, resp.StatusCode)
	return false, nil
}

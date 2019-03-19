package main

import (
	"flag"
	"fmt"
	"net/http"
	"os/exec"
	"time"
)

type testResult struct {
	id          string
	push        bool
	pushElapsed int
	verify      bool
}

func pushApp(round int, index int, cfHost string) (pushed bool, elapsed int, url string) {
	start := time.Now().Unix()
	appName := fmt.Sprintf("app-%d-%d", round, index)
	fmt.Printf("[%s]:Pushing\n", appName)
	cmd := exec.Command("cf", "push", appName)
	_, err := cmd.Output()
	end := time.Now().Unix()
	ela := end - start
	fmt.Printf("[%s]:Push elapsed: %v seconds\n", appName, ela)
	if err != nil {
		fmt.Printf("[%s]:Push failed with error:\n%s\n", appName, err.Error())
		return false, 0, ""
	}
	fmt.Printf("[%s]:Push successfully\n", appName)
	url = fmt.Sprintf("http://%s.%s", appName, cfHost)
	return true, int(ela), url
}

func httpVerifyApp(round int, index int, url string) (verified bool) {
	appName := fmt.Sprintf("app-%d-%d", round, index)
	fmt.Printf("[%s]:Sending http request to verify (URL: %s)\n", appName, url)
	time.Sleep(6 * time.Second)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("[%s]:Send http request failed\n", appName)
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		fmt.Printf("[%s]:http request successfully with code: %d\n", appName, resp.StatusCode)
		return true
	}
	fmt.Printf("[%s]:http request failed with code: %d\n", appName, resp.StatusCode)
	return false
}

func summaryTest(testCount int, testResults []testResult) {
	var pass, failed int
	var totalPushElapsed int
	allAvailablePushElapsed := make([]int, testCount)

	for i, result := range testResults {
		fmt.Printf("[DEBUG]:%+v\n", result)
		if result.push && result.verify {
			pass++
		} else {
			failed++
		}
		allAvailablePushElapsed[i] = result.pushElapsed
		totalPushElapsed += result.pushElapsed
	}

	// get the fastest push
	fastest := allAvailablePushElapsed[0]
	for i := 1; i < testCount; i++ {
		if allAvailablePushElapsed[i] < fastest {
			fastest = allAvailablePushElapsed[i]
		}
	}

	// get the slowest push
	slowest := allAvailablePushElapsed[0]
	for j := 1; j < testCount; j++ {
		if allAvailablePushElapsed[j] > slowest {
			slowest = allAvailablePushElapsed[j]
		}
	}

	successRate := (float64(pass) / float64(testCount)) * 100
	failureRate := (float64(failed) / float64(testCount)) * 100

	fmt.Printf("Summary\n push app %d times, pass %d times, failed %d times\n success rate: %.2f%%, failure rate: %.2f%%\n push elapsed\n  average (success): %d seconds\n  fastest: %v seconds\n  slowest: %v seconds\n",
		testCount, pass, failed, successRate, failureRate, totalPushElapsed/testCount, fastest, slowest)
}

func serialRun(count int, cfHost string) (testResults []testResult) {
	fmt.Printf("Serial Test Run, push ----> %d <---- apps sequentially\n", count)
	testResults = make([]testResult, count)
	// serial
	for i := 1; i <= count; i++ {
		pushed, elapsed, url := pushApp(1, i, cfHost)
		var verified bool
		if pushed {
			verified = httpVerifyApp(1, i, url)
		} else {
			verified = false
		}
		testResults[i-1] = testResult{
			id:          fmt.Sprintf("%d-%d", 1, i),
			push:        pushed,
			pushElapsed: elapsed,
			verify:      verified,
		}
	}
	return
}

func concurrencyRun(round int, concurrency int, cfHost string) (testResults []testResult) {
	fmt.Printf("Concurrency Test Run, push ----> %d <---- apps concurrently in round ----> %d <----\n", concurrency, round)
	testResults = make([]testResult, concurrency)
	chResults := make(chan testResult)

	for i := 1; i <= concurrency; i++ {
		go func(index int) {
			pushed, elapsed, url := pushApp(round, index, cfHost)
			var verified bool
			if pushed {
				verified = httpVerifyApp(round, index, url)
			} else {
				verified = false
			}
			chResults <- testResult{
				id:          fmt.Sprintf("%d-%d", round, index),
				push:        pushed,
				pushElapsed: elapsed,
				verify:      verified,
			}
		}(i)
	}

	for i := 1; i <= concurrency; i++ {
		select {
		case rst := <-chResults:
			testResults[i-1] = rst
		}
	}

	return
}

func main() {

	// define flags
	cfHost := flag.String("host", "", "specify the host of cf")
	runConc := flag.Bool("runConcurrently", false, "specify it when run tests concurrently")
	testsCount := flag.Int("testCount", 5, "specify the amount of apps will be pushed")
	conc := flag.Int("concurrency", 5, "specify the concurrency when run test concurrently")

	flag.Parse()

	var testResults []testResult
	// serial run
	if !*runConc {
		// serial
		testResults = serialRun(*testsCount, *cfHost)
	} else {
		// concurrent
		rounds := *testsCount / *conc
		for rd := 1; rd <= rounds; rd++ {
			tr := concurrencyRun(rd, *conc, *cfHost)
			testResults = append(testResults, tr...)
		}
	}

	// summary
	summaryTest(*testsCount, testResults)
}

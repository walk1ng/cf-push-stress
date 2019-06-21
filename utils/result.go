package utils

import (
	"github.com/walk1ng/cf-push-stress/model"
)

func SummaryTest(testCount int, testResults []model.PushAppResult) {
	var pass, failed int
	var totalPushElapsed int
	allAvailablePushElapsed := make([]int, testCount)

	for i, result := range testResults {
		logger.Printf("%+v\n", result)
		if result.PushSucced && result.HTTPVerificationSucced {
			pass++
		} else {
			failed++
		}
		allAvailablePushElapsed[i] = result.PushElapsed
		totalPushElapsed += result.PushElapsed
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

	logger.Printf("Summary\n push app %d times, pass %d times, failed %d times\n success rate: %.2f%%, failure rate: %.2f%%\n push elapsed\n  average (success): %d seconds\n  fastest: %v seconds\n  slowest: %v seconds\n",
		testCount, pass, failed, successRate, failureRate, totalPushElapsed/testCount, fastest, slowest)
}

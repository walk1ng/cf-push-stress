package utils

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/walk1ng/cf-push-stress/model"
)

func pushApp(app model.App) (pushSucced bool, elapsed int) {

	log.Printf("[%s]:Pushing\n", app.Name)
	start := time.Now().Unix()
	cmd := exec.Command("cf", "push", app.Name)
	_, err := cmd.Output()
	end := time.Now().Unix()
	ela := end - start

	if err != nil {
		log.Fatalf("[%s]:Push failed with error:\n%s\n", app.Name, err.Error())
		return false, 0
	}

	log.Printf("[%s]:Push successfully and elapsed: %v seconds\n", app.Name, ela)

	return true, int(ela)
}

// SerialPush func
func SerialPush(count int, cfHost string) (testResults []model.PushAppResult) {

	log.Printf("Serial Test Run, push ----> %d <---- apps sequentially\n", count)
	testResults = make([]model.PushAppResult, count)
	// serial
	for i := 1; i <= count; i++ {

		app := model.App{
			Name:   fmt.Sprintf("app-%d-%d", 1, i),
			Domain: cfHost,
		}

		pushSucced, elapsed := pushApp(app)
		testResults[i-1] = model.PushAppResult{
			App:         app,
			PushSucced:  pushSucced,
			PushElapsed: elapsed,
		}
	}

	time.Sleep(5 * time.Second)
	testResults = finalVerify(testResults)

	return
}

// ConcurrencyPush func
func ConcurrencyPush(round int, concurrency int, cfHost string) (testResults []model.PushAppResult) {

	log.Printf("Concurrency Test Run, push ----> %d <---- apps concurrently in round ----> %d <----\n", concurrency, round)
	testResults = make([]model.PushAppResult, concurrency)
	chResults := make(chan model.PushAppResult)

	for i := 1; i <= concurrency; i++ {
		go func(index int) {
			app := model.App{
				Name:   fmt.Sprintf("app-%d-%d", round, index),
				Domain: cfHost,
			}
			pushSucced, elapsed := pushApp(app)

			chResults <- model.PushAppResult{
				App:         app,
				PushSucced:  pushSucced,
				PushElapsed: elapsed,
			}
		}(i)
	}

	for i := 1; i <= concurrency; i++ {
		select {
		case rst := <-chResults:
			testResults[i-1] = rst
		}
	}

	time.Sleep(5 * time.Second)
	testResults = finalVerify(testResults)

	return
}

func finalVerify(rsts []model.PushAppResult) []model.PushAppResult {
	for i := range rsts {
		if rsts[i].PushSucced {
			if httpVerificationSucced, err := doHTTPVerifyApp(rsts[i].App); err != nil {
				// try once
				rsts[i].HTTPVerificationSucced, _ = doHTTPVerifyApp(rsts[i].App)
			} else {
				rsts[i].HTTPVerificationSucced = httpVerificationSucced
			}
		} else {
			rsts[i].HTTPVerificationSucced = false
		}
	}

	return rsts
}

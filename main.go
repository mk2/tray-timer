package main

import (
	"fmt"
	"time"

	"github.com/getlantern/systray"
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTitle("TrayTimer")
	m5Min := systray.AddMenuItem("5 Minutes", "Start 5 minutes timer")
	m10Min := systray.AddMenuItem("10 Minutes", "Start 10 minutes timer")
	m15Min := systray.AddMenuItem("15 Minutes", "Start 15 minutes timer")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quite the app")

	finishCh := make(chan struct{})

	go func() {
		<-mQuit.ClickedCh
		close(finishCh)
		systray.Quit()
	}()

	for {
		select {
		case <-m5Min.ClickedCh:
			go startTimer(5, finishCh)
		case <-m10Min.ClickedCh:
			go startTimer(10, finishCh)
		case <-m15Min.ClickedCh:
			go startTimer(15, finishCh)
		}

		<-finishCh

		fmt.Println("timer finished")
	}
}

func onExit() {

}

func startTimer(min int, finishCh chan<- struct{}) {
	startTime := time.Now()
	timer := time.After(time.Duration(min) * time.Second)
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for now := range ticker.C {
			diff := now.Sub(startTime)
			diff = diff.Round(time.Second)
			systray.SetTitle(fmt.Sprintf("%02.f:%02.f", diff.Minutes(), diff.Seconds()))
		}
	}()

	<-timer
	ticker.Stop()
	<-time.After(1 * time.Second)
	systray.SetTitle("timer finished!!")
	finishCh <- struct{}{}
}

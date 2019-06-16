package main

import (
	"sync"
	"testing"
	"time"

	"github.com/getlantern/systray"
	"github.com/promignis/cwitch/menu"
	"github.com/promignis/cwitch/timer"
)

func setupSystray() {
	InitMenu()
	systray.Run(onReady, onExit)
}

func TestSequentialItemClicks(t *testing.T) {
	go setupSystray()
	time.Sleep(time.Second * 3)
	for _, m := range menu.MenuMap {
		if m.MenuItem != nil {
			multiClick(m.MenuItem.ClickedCh, 3)
			time.Sleep(time.Second * 3)
		}
	}
}

func TestParallelItemClicks(t *testing.T) {
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(menu.MenuMap))
	for _, m := range menu.MenuMap {
		// closure
		go func(m *timer.Timer) {
			if m.MenuItem != nil {
				multiClick(m.MenuItem.ClickedCh, 3)
				time.Sleep(time.Second * 3)
				waitGroup.Done()
			}
		}(m)
	}
	waitGroup.Wait()
}

func TestSameItemClick(t *testing.T) {
	ticker := time.NewTicker(time.Second)

	// 5 seconds
	count := 5
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(menu.MenuMap))
	for _, m := range menu.MenuMap {
		go func(m *timer.Timer) {
			for t := range ticker.C {
				_ = t
				count -= 1
				if count == 0 {
					break
				}
				click(m.MenuItem.ClickedCh)
				waitGroup.Done()
			}
		}(m)
	}
	waitGroup.Wait()
}

func click(clickedCh chan struct{}) {
	clickedCh <- struct{}{}
}

func multiClick(clickedCh chan struct{}, n int) {
	for i := 0; i < n; i++ {
		clickedCh <- struct{}{}
	}
}

package utils

import (
	"fmt"
	"time"
)

type DotProgressBar struct {
	closer chan struct{}
}

func NewDotProgressBar() *DotProgressBar {
	return &DotProgressBar{
		closer: make(chan struct{}),
	}
}

func (d *DotProgressBar) Start() {
	go func() {
		for {
			select {
			case <-time.Tick(1000 * time.Millisecond):
				fmt.Print(".")
			case <-d.closer:
				return
			}
		}
	}()
}

func (d *DotProgressBar) Stop() {
	d.closer <- struct{}{}
}

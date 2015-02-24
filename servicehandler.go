package main

import (
	"sync"
)

type ServiceHandler struct {
	ch chan bool
	wg *sync.WaitGroup
}

func NewServiceHandler() *ServiceHandler {
	s := &ServiceHandler {
		ch: make(chan bool),
		wg: &sync.WaitGroup{},
    }
	s.wg.Add(1) // Add main wait. Closed on Stop().
    return s
}

func (s *ServiceHandler) addWatchedProcess() {
	s.wg.Add(1)
}

func (s *ServiceHandler) finishProcess() {
	s.wg.Done()
}

// Stop the service by closing the service's channel.  Block until the service
// is really stopped.
func (s *ServiceHandler) Stop() {
	close(s.ch)
	s.wg.Done()
	s.wg.Wait()
}

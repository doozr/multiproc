package multiproc

import (
	"fmt"
	"sync"
)

// MessageProcessor to process incoming messages
type MessageProcessor func(interface{}) error

// processMessages until one returns an error or the shutdown channel is closed
func processMessages(messages <-chan interface{}, processor MessageProcessor, shutdown <-chan interface{}) (err error) {
	for {
		select {
		case msg, ok := <-messages:
			if !ok {
				err = fmt.Errorf("Delivery channel closed unexpectedly")
				return
			}

			err = processor(msg)
			if err != nil {
				return
			}

		case <-shutdown:
			return
		}
	}
}

// ProcessConcurrent messages from the messages channel until an error occurs or done is closed
func ProcessConcurrent(messages <-chan interface{}, processor MessageProcessor, concurrency int, done <-chan interface{}) (err error) {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(concurrency)

	waitChan := make(chan interface{})
	errorChan := make(chan error)
	shutdown := make(chan interface{})

	defer close(errorChan)

	// Run a number of processors in goroutines
	for x := 0; x < concurrency; x++ {
		go func() {
			defer waitGroup.Done()

			err := processMessages(messages, processor, shutdown)
			if err != nil {
				errorChan <- err
				return
			}
		}()
	}

	go func() {
		// Wait until all children are stopped then close waitChan
		waitGroup.Wait()
		close(waitChan)
	}()

	for {
		select {
		case err = <-errorChan:
			// If one of the goroutines has an error, set err and close the shutdown channel
			close(shutdown)
		case _, ok := <-done:
			// If the done channel is closed, close the shutdown channel and wait for goroutines to stop
			if !ok {
				close(shutdown)
			}
		case _, ok := <-waitChan:
			// If waitchan is closed, we're done
			if !ok {
				return
			}
		}
	}
}

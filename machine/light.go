package machine

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type Light struct {
	control         uint
	shutdownChannel chan struct{}
	logger          *log.Logger
}

const (
	Yellow uint = 0
	Green  uint = 1
	Red    uint = 2
	off    uint = 4
)

func (l *Light) init() error {
	return nil
}

func (l *Light) open() error {
	return nil
}

func (l *Light) run(commChan chan struct{}, state uint, lightFinishWg *sync.WaitGroup) error {

	runDoneChannel := make(chan struct{})

	defer close(runDoneChannel)

	go monitorFunction(runDoneChannel, l.shutdownChannel, l.abort, l.logger) // Running the monitor function

	select {
	case <-commChan:

	case <-l.shutdownChannel:

		l.abort()
		lightFinishWg.Done()
		return errors.New("Shutdown Called")
	}

	l.control = state

	switch {
	case l.control == Yellow: // Yellow State, to signify a wait state/test run
		log.Println("Light is in Yellow State")
	case l.control == Green: //Green State, to signify the on running state
		log.Println("Light is in Green State")
	case l.control == Red: //Red State, to signify Error state
		log.Println("Light is in Red State")
	case l.control == off: // Off State, only used to simulate blinking during the estop process
		log.Println("Light is in off State")
	}

	lightFinishWg.Done()

	return nil
}

func (l *Light) close() error {
	os.Exit(0)
	return nil
}

func (l *Light) abort() error {
	fmt.Println("The problem is in light") // For debugging Purposes
	os.Exit(1)
	return nil
}
func (l *Light) estop(length int) error {
	// Simulating blinking of a lightbulb during the estop function
	for i := 0; i < length; i++ {
		l.control = Red
		log.Println("Red")
		time.Sleep(100 * time.Millisecond)
		l.control = off
		log.Println("off")
	}

	return nil
}

func newLight(logger *log.Logger, shutdownChannel chan struct{}) *Light {
	return &Light{
		shutdownChannel: shutdownChannel,
		logger:          logger,
	}
}

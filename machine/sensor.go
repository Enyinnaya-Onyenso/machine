package machine

import (
	"fmt"
	"log"
	"os"
	"sync"
)

type Sensor struct {
	isDone          bool
	shutdownChannel chan struct{}
	logger          *log.Logger
}

func (s *Sensor) init() error {
	return nil
}

func (s *Sensor) open() error {
	return nil
}

func (s *Sensor) run(commChan chan struct{}, commFinishWg *sync.WaitGroup) error {

	runDoneChannel := make(chan struct{})
	defer close(runDoneChannel)

	go monitorFunction(runDoneChannel, s.shutdownChannel, s.abort, s.logger) // Running the monitor function

	s.isDone = true
	log.Println("Sensor reads true")
	select {
	case <-commChan:
		// communication channel closed
	default:
		close(commChan)
	}

	commFinishWg.Wait() // wait

	return nil
}

func (s *Sensor) close() error {
	os.Exit(0)
	return nil
}

func (s *Sensor) abort() error {
	fmt.Println("The problem is in sensor") // For debugging purposes
	os.Exit(1)
	return nil
}

func newSensor(logger *log.Logger, shutdownChannel chan struct{}) *Sensor {
	return &Sensor{
		isDone:          false,
		shutdownChannel: shutdownChannel,
		logger:          logger,
	}
}

package machine

import (
	"errors"
	"log"
	"math/rand"
	"os"
	"sync"
)

type BarScanner struct {
	reader          int
	isDoneReading   bool
	shutdownChannel chan struct{}
	logger          *log.Logger
}

func (b *BarScanner) init() error {
	b.reader = 0
	return nil
}

func (b *BarScanner) open() error {
	return nil
}

func (b *BarScanner) run(commChan chan struct{}, commFinishWg *sync.WaitGroup) error {
	runDoneChannel := make(chan struct{})
	defer close(runDoneChannel)

	go monitorFunction(runDoneChannel, b.shutdownChannel, b.abort, b.logger) // Running the monitor function
	select {
	case <-commChan:

	case <-b.shutdownChannel:

		b.abort()
		return errors.New("Shutdown Called")
	}

	b.isDoneReading = true

	b.reader = rand.Intn(10) // Setting a value for the reader to simulate scanning a barcode
	log.Println("Barcode has read: ", b.reader)

	return nil
}

func (b *BarScanner) close() error {
	os.Exit(0)
	return nil
}

func (b *BarScanner) abort() error {
	os.Exit(1)
	return nil
}

func newBarScanner(logger *log.Logger, shutdownChannel chan struct{}) *BarScanner {
	return &BarScanner{
		logger:          logger,
		shutdownChannel: shutdownChannel,
	}
}

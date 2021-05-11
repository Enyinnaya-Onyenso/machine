package machine

import (
	"fmt"
	"log"
	"os"
)

type Conveyor struct {
	// The Conveyor is going to edit values in the other components.
	belt            int
	shutdownChannel chan struct{}
	logger          *log.Logger
	barcode         *BarScanner
	lights          []*Light
	sensors         []*Sensor
}

func (c *Conveyor) init() error {
	c.belt = 1 // Setting the intial value of the belt to 1.
	return nil
}

func (c *Conveyor) open() error {
	return nil
}

func (c *Conveyor) run() error {
	runDoneChannel := make(chan struct{})
	defer close(runDoneChannel)

	go monitorFunction(runDoneChannel, c.shutdownChannel, c.abort, c.logger) // Running the monitor function

	if c.barcode.isDoneReading {
		for i := 1; i < len(c.lights); i++ {
			switch {
			case c.lights[i].control == 0:
				fmt.Println("Test mode, light ", i, " is in Yellow State, Please Wait")
			case c.lights[i].control == 1:
				c.lights[i].control = 0 // Conveyor editing the control value for the lights
				c.logger.Println("Light", i, " is now in Yellow State")
				for i := 0; i < len(c.sensors); i++ {
					c.sensors[i].isDone = false // Conveyor editing the isDone value for the sensors
				}
			case c.lights[i].control == 2:
				fmt.Println("Error found ", i, " is in Red State, Please Wait")
			case c.lights[i].control == 3:
				fmt.Println("Light ", i, " is off")
			}
		}

		c.belt++
		log.Println("Now on belt ", c.belt)
	}

	return nil
}

func (c *Conveyor) close() error {
	os.Exit(0)
	return nil
}

func (c *Conveyor) abort() error {
	os.Exit(1)
	return nil
}

func newConveyor(logger *log.Logger, shutdownChannel chan struct{}, sensors []*Sensor, lights []*Light, barcode *BarScanner) *Conveyor {
	return &Conveyor{
		logger:          logger,
		shutdownChannel: shutdownChannel,
		sensors:         sensors,
		lights:          lights,
		barcode:         barcode,
	}
}

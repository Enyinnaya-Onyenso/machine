package machine

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"
	"unicode"

	//	"fmt"
	"log"
	"sync"
)

type System struct {
	sensors         []*Sensor
	lights          []*Light
	barscanner      *BarScanner
	conveyor        *Conveyor
	shutDownChannel chan struct{}
	logger          *log.Logger
}

func (s *System) Init() error {
	// Initializing Logger and functions

	s.logger = log.New(new(bytes.Buffer), "Log", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

	err := s.barscanner.init()
	if err != nil {
		s.logger.Println(err)
	}

	err = s.conveyor.init()
	if err != nil {
		s.logger.Println(err)
	}

	for i := 0; i < len(s.sensors); i++ {
		err = s.sensors[i].init()
		if err != nil {
			s.logger.Println(err)
		}
	}

	for i := 0; i < len(s.lights); i++ {
		err = s.lights[i].init()
		if err != nil {
			s.logger.Println(err)
		}
	}

	return nil
}

func (s *System) Open() error {
	err := s.barscanner.open()
	if err != nil {
		s.logger.Println(err)
	}

	err = s.conveyor.open()
	if err != nil {
		s.logger.Println(err)
	}

	for i := 0; i < len(s.sensors); i++ {
		err = s.sensors[i].open()
		if err != nil {
			s.logger.Println(err)
		}
	}

	for i := 0; i < len(s.lights); i++ {
		err = s.lights[i].open()
		if err != nil {
			s.logger.Println(err)
		}
	}
	return nil
}

func (s *System) Run() error {

	s.shutDownChannel = make(chan struct{}) // Ensuring shutdown channel
	commChan := make(chan struct{})         // Initializing communication channel
	defer close(commChan)                   // Ensuring communication channel closes after successful run
	var lightFinish sync.WaitGroup

	var wg sync.WaitGroup

	errArr := make([]error, 0)

	done := make(chan struct{})
	defer close(done) // Ensuring done closes after successful run

	fmt.Println("Enter Command: ")
	commandInput := bufio.NewScanner(os.Stdin)
	commandInput.Scan()
	command := commandInput.Text()     // Converting Scan to text
	cmdrune := []rune(command)         // Converting to rune for use in unicode package commands
	cmdlow := strings.ToLower(command) // Converting capital case letters to lower case letters
	var (
		// Variables for use in unicode conversions
		cmdArr   []bool
		isnumber bool
		isletter bool
		num      bool
		state    uint
	)
	runNum := 1
	for i := 0; i < len(cmdrune); i++ {
		if unicode.IsNumber(cmdrune[i]) {
			cmdArr = append(cmdArr, true)
		} else if unicode.IsLetter(cmdrune[i]) {
			cmdArr = append(cmdArr, false)
		}
	}
	for _, num = range cmdArr {
		// Setting booleans in the array to number or letter
		if !num {
			isletter = true
		} else {
			isnumber = true
		}
	}
	// Forcing the numbers/letters to a single type
	if isnumber {
		if isletter {
			isletter = false
		} else {
			isnumber = true
		}
	}

	if isletter && !isnumber {
		isletter = true
	}
	time.Sleep(2 * time.Second)
	if len(cmdArr) > 1 {
		switch { // Switch to determine the state of the machine when the user inputs multiple characters
		case isnumber: // Estop State
			s.logger.Println("ESTOP Initiated")
			fmt.Println("Estop Initiated")
			err := s.Estop()
			if err != nil {
				s.logger.Println(err)
			}
		case isletter: // Test Run State
			for i := 0; i < len(cmdArr); i++ {
				fmt.Println("Test mode Run")
				for i := 0; i < len(s.sensors); i++ {
					s.sensors[i].shutdownChannel = s.shutDownChannel
					wg.Add(1)
					go func(index int) {
						defer wg.Done()
						errArr = append(errArr, s.sensors[index].run(commChan, &lightFinish))
					}(i)
				}

				wg.Wait()

				var aggregateErr error
				for _, err := range errArr {
					if err != nil {
						aggregateErr = fmt.Errorf("%v : %v", aggregateErr, err)
					}
				}
				if aggregateErr != nil {
					fmt.Println(aggregateErr)
					return aggregateErr
				}

				time.Sleep(2 * time.Second)

				s.barscanner.shutdownChannel = s.shutDownChannel
				wg.Add(1)
				go func() {
					defer wg.Done()
					errArr = append(errArr, s.barscanner.run(commChan, &lightFinish))
				}()

				for i := 0; i < len(s.lights); i++ {
					s.lights[i].shutdownChannel = s.shutDownChannel
					wg.Add(1)
					lightFinish.Add(1)
					go func(index int) {
						defer wg.Done()
						state = 0
						errArr = append(errArr, s.lights[index].run(commChan, state, &lightFinish))
					}(i)
				}
				wg.Wait()

				for _, err := range errArr {
					if err != nil {
						aggregateErr = fmt.Errorf("%v : %v", aggregateErr, err)
					}
				}
				if aggregateErr != nil {
					fmt.Println(aggregateErr)
					return aggregateErr
				}

				time.Sleep(2 * time.Second)

				s.conveyor.shutdownChannel = s.shutDownChannel
				err := s.conveyor.run()
				if err != nil {
					fmt.Println(err)
					return err
				}
				time.Sleep(2 * time.Second)
			}
		}

	} else if len(cmdArr) == 1 {
		switch {
		// Switch to determine the state of the machine when the user inputs a single character
		case isnumber: // ESTOP State
			s.logger.Println("ESTOP Initiated")
			err := s.Estop()
			if err != nil {
				s.logger.Println(err)
			}
		case isletter && cmdlow != "y": // Idle State
			s.logger.Println("Part is still in idle")
			fmt.Println("Part is still in idle")
		case isletter && cmdlow == "y": // Normal run state
			s.logger.Println("Starting Machine")
			fmt.Println("Starting Machine")
			for i := 0; i < len(s.sensors); i++ {
				s.sensors[i].shutdownChannel = s.shutDownChannel
				wg.Add(1)
				go func(index int) {
					defer wg.Done()
					errArr = append(errArr, s.sensors[index].run(commChan, &lightFinish))
				}(i)
			}

			wg.Wait()

			var aggregateErr error
			for _, err := range errArr {
				if err != nil {
					aggregateErr = fmt.Errorf("%v : %v", aggregateErr, err)
				}
			}
			if aggregateErr != nil {
				fmt.Println(aggregateErr)
				return aggregateErr
			}

			time.Sleep(2 * time.Second)

			s.barscanner.shutdownChannel = s.shutDownChannel
			wg.Add(1)
			go func() {
				defer wg.Done()
				errArr = append(errArr, s.barscanner.run(commChan, &lightFinish))
			}()

			for i := 0; i < len(s.lights); i++ {
				s.lights[i].shutdownChannel = s.shutDownChannel
				wg.Add(1)
				lightFinish.Add(1)
				go func(index int) {
					defer wg.Done()
					state = 1
					errArr = append(errArr, s.lights[index].run(commChan, state, &lightFinish))
				}(i)
			}
			wg.Wait()

			for _, err := range errArr {
				if err != nil {
					aggregateErr = fmt.Errorf("%v : %v", aggregateErr, err)
				}
			}
			if aggregateErr != nil {
				fmt.Println(aggregateErr)
				return aggregateErr
			}

			time.Sleep(2 * time.Second)

			s.conveyor.shutdownChannel = s.shutDownChannel
			err := s.conveyor.run()
			if err != nil {
				fmt.Println(err)
				return err
			}
			s.logger.Println("You have run ", runNum, "times")
			runNum++
			time.Sleep(2 * time.Second)
			s.Run()
		}
	}

	return nil
}

func (s *System) Close() error {

	err := s.barscanner.close()
	if err != nil {
		s.logger.Println(err)
	}

	err = s.conveyor.close()
	if err != nil {
		s.logger.Println(err)
	}

	for i := 0; i < len(s.sensors); i++ {
		err = s.sensors[i].close()
		if err != nil {
			s.logger.Println(err)
		}
	}

	for i := 0; i < len(s.lights); i++ {
		err = s.lights[i].close()
		if err != nil {
			s.logger.Println(err)
		}
	}
	return nil
}

func (s *System) Abort() error {
	select {
	case <-s.shutDownChannel:

	default:
		close(s.shutDownChannel)
	}
	return nil
}
func (s *System) Estop() error {
	var wg sync.WaitGroup
	for i := 0; i < len(s.lights); i++ {
		time.Sleep(1 * time.Second)
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			fmt.Println("light", index)
			s.lights[index].estop(len(s.lights))

		}(i)
		time.Sleep(1 * time.Second)
	}
	wg.Wait()
	os.Exit(21) // Status 21 meaning Emergency Stop
	return nil
}

func NewSystem(sensorNum int, lightNum int) *System {

	logger := log.New(new(bytes.Buffer), "Log", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

	systemChan := make(chan struct{})

	systemSensors := make([]*Sensor, 0)
	for i := 0; i < sensorNum; i++ {
		systemSensors = append(systemSensors, newSensor(logger, systemChan))
	}

	systemLights := make([]*Light, 0)
	for i := 0; i < lightNum; i++ {
		systemLights = append(systemLights, newLight(logger, systemChan))
	}

	systemBarScanner := newBarScanner(logger, systemChan)

	systemConveyor := newConveyor(logger, systemChan, systemSensors, systemLights, systemBarScanner)

	return &System{
		sensors:         systemSensors,
		barscanner:      systemBarScanner,
		shutDownChannel: systemChan,
		conveyor:        systemConveyor,
		lights:          systemLights,
	}
}

func (s *System) Logger() *log.Logger {
	// Exposing logger to log in main.go
	return s.logger
}

## Problem Statement:
Simulate a machine with components: Sensors, Lights, Barcode Scanner and a Conveyor. Take user input to run the machine.
If the user enters "y" or "Y", run the machine in a normal mode
If the user enters a single alphabet that is not "y" or "Y" notify the user that the machine is still idle
If the user enters multiple alphabets, run the machine in test mode. Test mode runs the code as many times as the lenght of the input characters
If there is a number among the input characters, notify the user and send the machine into an Emergency stop

If the machine is in normal mode, turn the lights on after the sensor reads and turn it to a wait state before the conveyor moves

If the machine is in test mode, the lights should always be in the wait state

When an error is found, the lights enter an alert state and blink

### Run  Instructions:
Navigate to cmd folder and run:
`go run main.go`
package main

import (
	"machine"
)

func main() {
	wantedSensorNum := 10
	wantedLightNum := 10

	thisMachine := machine.NewSystem(wantedSensorNum, wantedLightNum)

	err := thisMachine.Init()
	if err != nil {
		thisMachine.Logger().Println(err)
		return
	}

	err = thisMachine.Open()
	if err != nil {
		thisMachine.Logger().Println(err)
		return
	}

	defer thisMachine.Close()

	err = thisMachine.Run()
	if err != nil {
		thisMachine.Logger().Println(err)
		return
	}
}

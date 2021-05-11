package machine

import "log"

func monitorFunction(runChannel chan struct{}, shutdownChannel chan struct{}, abort func() error, logger *log.Logger) {
	// Monitor function for all the component run functions
	select {
	case <-runChannel:
		return
	case <-shutdownChannel:
		err := abort()
		if err != nil {
			logger.Println(err)
		}
		return
	}
}

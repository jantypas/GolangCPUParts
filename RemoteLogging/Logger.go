package RemoteLogging

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

var logFile *os.File
var logChannel chan LogEventStruct
var logApp string
var loggingActive bool = false

type LogEventStruct struct {
	EventTime   string `json:"event_time"`
	EventApp    string `json:"event_app"`
	EventLevel  string `json:"event_level"`
	EventSource string `json:"event_source"`
	EventMsg    string `json:"event_msg"`
}

func LogInit(appname string) error {
	f, err := os.OpenFile("/tmp/"+appname+".log", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Error opening log file")
		return err
	} else {
		logFile = f
		logApp = appname
		logChannel = make(chan LogEventStruct, 5)
		go func() {
			for {
				msg := <-logChannel
				msg.EventApp = logApp
				msg.EventTime = time.Now().Format("2006-01-02 15:04:05")
				s, _ := json.Marshal(msg)
				logFile.Write(s)
			}
		}()
		time.Sleep(1)
		return nil
	}
}

func SetLogginActive(state bool) {
	loggingActive = state
}

func LogEvent(msg LogEventStruct) {
	if loggingActive {
		logChannel <- msg
	}
}

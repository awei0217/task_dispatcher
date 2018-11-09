package common

import (
	"github.com/sirupsen/logrus"
	"os"
)

var log *logrus.Logger

func init()  {
	log = logrus.New()
	log.SetLevel(logrus.InfoLevel)
	log.Formatter = &logrus.TextFormatter{
				TimestampFormat:YYYY_MM_DD_HH_MM_SS,
			}
	/*file, err := os.OpenFile("F:\\export\\Logs\\task.dispatcher\\task_dispatcher.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		//log.Out = file
	} else {
		log.Info("Failed to log to file, using default stderr")
	}*/
	log.SetOutput(os.Stdout)

}

func GetLog() (*logrus.Logger)  {

	return log
}
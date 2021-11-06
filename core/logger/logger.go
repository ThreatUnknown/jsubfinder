package logger

import (
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	Log *logrus.Logger
	//zip    = flag.String("zipkin", os.Getenv("ZIPKIN"), "Zipkin address")
	//	port        = flag.String("port", "8080", "Port number on which the service should run")
	//	ip          = flag.String("ip", "0.0.0.0", "Preferred IP address to run the service on")
	//serviceName = "user"
)

// This initiates a new Logger and defines the format for logs
func InitDetailedLogger(f *os.File) {

	Log = logrus.New()
	Log.SetReportCaller(true)

	Log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "",
		PrettyPrint:     true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			s := strings.Split(f.Function, ".")
			funcname := s[len(s)-1]
			_, filename := path.Split(f.File)
			return funcname, filename
		},
	})

	// Set output of logs to Stdout
	// Change to f for redirecting to file
	Log.SetOutput(os.Stdout)

}

package llog

import (
	"io"
	"os"
	"strings"
)

var (
	Debug       = StdDebug
	Error       = StdError
	Warn        = StdWarn
	Info        = StdInfo
	CalledDebug = StdCalledDebug
	CalledError = StdCalledError
	CalledWarn  = StdCalledWarn
	CalledInfo  = StdCalledInfo
)

var (
	Writter io.Writer = os.Stdout
)

func SetLevel(in string) {
	switch strings.ToUpper(in) {
	case "INFO":
		Level = L_Info
	case "DEBUG":
		Level = L_Debug
	case "WARN":
		Level = L_Warn
	case "FATAL":
		Level = L_Fatal
	case "ERROR":
		Level = L_Error
	case "OFF":
		Level = L_Off
	default:
		Level = L_All
	}
}

//func SetType() {
//	viper.AddConfigPath("conf")
//	viper.SetConfigName("conf")
//	viper.SetConfigType("yaml")
//	if err := viper.ReadInConfig(); err != nil {
//		panic("Config init failed:" + err.Error())
//	}
//	l_type := viper.GetString("log.writers")
//	level := viper.GetString("log.logger_level")
//	if l_type == "file" {
//		fn := viper.GetString("log.logger_file")
//		filename := path.Join("/log", fn)
//		SetTypeFile(filename)
//	} else if l_type == "default" {
//		SetTypeStream()
//	}
//	SetLevel(level)
//}

func SetTypeFile(filename string) {
	var flogger = LogFile{FileName: filename}
	Debug = flogger.Debug
	Error = flogger.Error
	Warn = flogger.Warn
	Info = flogger.Info
	CalledDebug = flogger.CalledDebug
	CalledError = flogger.CalledError
	CalledWarn = flogger.CalledWarn
	CalledInfo = flogger.CalledInfo
	Writter = &flogger
}

func SetTypeStream() {
	Debug = StdDebug
	Error = StdError
	Warn = StdWarn
	Info = StdInfo
	CalledDebug = StdCalledDebug
	CalledError = StdCalledError
	CalledWarn = StdCalledWarn
	CalledInfo = StdCalledInfo
	Writter = os.Stdout
}

package llog

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	SplitDay int   = 1
	KeepDay  int64 = 90
)

const (
	DATEFORMAT = "2006-01-02"
)

var TimeFormatRexMap = map[string]string{
	`^[+|-]\d{4}\s\d{4}-\d{2}-\d{2}\s\d{2}:\d{2}:\d{2}$`: "-0700 2006-01-02 15:04:05",
	`^\d{4}-\d{2}-\d{2}\s\d{2}:\d{2}:\d{2}$`:             "2006-01-02 15:04:05",
	`^\d{4}/\d{2}/\d{2}\s\d{2}:\d{2}:\d{2}[+|-]\d{4}$`:   "2006/01/02 15:04:05-0700",
	`^\d{4}/\d{2}/\d{2}T\d{2}:\d{2}:\d{2}[+|-]\d{4}$`:    "2006/01/02T15:04:05-0700",
	`^\d{4}/\d{2}/\d{2}\s\d{2}:\d{2}:\d{2}$`:             "2006/01/02 15:04:05",
	`^\d{4}/\d{2}/\d{2}T\d{2}:\d{2}:\d{2}$`:              "2006/01/02T15:04:05",
} // 字符串转时间戳匹配模式

// 根据指定时间格式返回时间戳
func ToTimeStamp(in string) int64 {
	var timeFormat = ""
	for r, v := range TimeFormatRexMap {
		if matched, _ := regexp.Match(r, []byte(in)); matched {
			timeFormat = v
			break
		}
	}
	if timeFormat == "" {
		return 0
	}
	ret, _ := time.ParseInLocation(timeFormat, in, time.Local)
	return ret.Unix()
}

type LogFile struct {
	FileName string
	mu       *sync.RWMutex
	date     *time.Time
	sinit    sync.Once
	name     string
	path     string
	io.Writer
}

func (slog *LogFile) Write(p []byte) (nn int, err error) {
	slog.mu.Lock()
	defer slog.mu.Unlock()
	f, err := slog.GetFile()
	if err != nil {
		StdError("fail to create log file!")
	}
	defer f.Close()
	return f.Write(p)

}

func (slog *LogFile) init() {
	slog.mu = new(sync.RWMutex)
	//t := time.Now()
	var t time.Time
	_, err := os.Stat(slog.FileName)
	if os.IsNotExist(err) {
		t, _ = time.Parse(DATEFORMAT, time.Now().Format(DATEFORMAT))
	} else {
		f, _ := os.OpenFile(slog.FileName, os.O_RDONLY, 0666)
		defer f.Close()
		d := make([]byte, 50)
		_, _ = f.Read(d)
		v := strings.Split(string(d), " ")
		t, _ = time.Parse("2006/01/02", v[1])
	}
	slog.date = &t
	path, name := filepath.Split(slog.FileName)
	slog.name = name
	slog.path, _ = filepath.Abs(path)
}

func (slog *LogFile) Clean() {
	path := filepath.Dir(slog.FileName)
	_, filename := filepath.Split(slog.FileName)
	rd, _ := ioutil.ReadDir(path)
	for _, fi := range rd {
		if fi.IsDir() {
			continue
		}
		name := fi.Name()
		_, _name := filepath.Split(name)
		if ok, _ := regexp.MatchString(filename+"\\."+".*", _name); ok {
			date := strings.Split(name, ".")[2]
			t, err := time.Parse(DATEFORMAT, date)
			if err != nil {
				continue
			}
			t = t.Add(time.Duration(KeepDay*24) * time.Hour)
			u := time.Now()
			if u.After(t) {
				os.Remove(filepath.Join(path, name))
			}
		}
	}
}
func (slog *LogFile) Split() {
	t, _ := time.Parse(DATEFORMAT, time.Now().Format(DATEFORMAT))
	u := slog.date.Add(time.Duration(SplitDay*24) * time.Hour)
	logFile := slog.FileName
	logFileBak := slog.FileName + "." + slog.date.Format(DATEFORMAT)
	if t.Before(u) {
		return
	}
	slog.Clean()
	_, err_file := os.Stat(logFileBak)
	if os.IsNotExist(err_file) {
		err := os.Rename(logFile, logFileBak)
		if err != nil {
			StdError(err)
		}
	}
	slog.date = &t
}
func (slog *LogFile) GetFile() (f *os.File, err error) {
	slog.Split()
	return os.OpenFile(slog.FileName, os.O_SYNC|os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
}
func (slog *LogFile) Debug(message ...interface{}) {
	if Level&r_debug == 0 {
		return
	}
	slog.sinit.Do(slog.init)
	slog.mu.Lock()
	defer slog.mu.Unlock()
	f, err := slog.GetFile()
	if err != nil {
		StdError("fail to create log file!")
	}
	defer f.Close()
	logger := log.New(f, "Debug: ", logflag)
	_ = logger.Output(2, fmt.Sprintln(message...))
}

func (slog *LogFile) Info(message ...interface{}) {
	if Level&r_debug == 0 {
		return
	}
	slog.sinit.Do(slog.init)
	slog.mu.Lock()
	defer slog.mu.Unlock()
	f, err := slog.GetFile()
	if err != nil {
		StdError("fail to create log file!")
	}
	defer f.Close()
	logger := log.New(f, "Info: ", logflag)
	_ = logger.Output(2, fmt.Sprintln(message...))
	//logger.Println(message...)
}

func (slog *LogFile) Warn(message ...interface{}) {
	if Level&r_debug == 0 {
		return
	}
	slog.sinit.Do(slog.init)
	slog.mu.Lock()
	defer slog.mu.Unlock()
	f, err := slog.GetFile()
	if err != nil {
		StdError("fail to create log file!")
	}
	defer f.Close()
	logger := log.New(f, "Warn: ", logflag)
	_ = logger.Output(2, fmt.Sprintln(message...))
}

func (slog *LogFile) Error(message ...interface{}) {
	if Level&r_debug == 0 {
		return
	}
	slog.sinit.Do(slog.init)
	slog.mu.Lock()
	defer slog.mu.Unlock()
	f, err := slog.GetFile()
	if err != nil {
		StdError("fail to create log file!")
	}
	defer f.Close()
	logger := log.New(f, "Error: ", logflag)
	_ = logger.Output(2, fmt.Sprintln(message...))
}

func (slog *LogFile) CalledDebug(calldepth int, message ...interface{}) {
	if Level&r_debug == 0 {
		return
	}
	slog.sinit.Do(slog.init)
	slog.mu.Lock()
	defer slog.mu.Unlock()
	f, err := slog.GetFile()
	if err != nil {
		StdError("fail to create log file!")
	}
	defer f.Close()
	logger := log.New(f, "Debug: ", logflag)
	_ = logger.Output(calldepth, fmt.Sprintln(message...))
}

func (slog *LogFile) CalledInfo(calldepth int, message ...interface{}) {
	if Level&r_debug == 0 {
		return
	}
	slog.sinit.Do(slog.init)
	slog.mu.Lock()
	defer slog.mu.Unlock()
	f, err := slog.GetFile()
	if err != nil {
		StdError("fail to create log file!")
	}
	defer f.Close()
	logger := log.New(f, "Info: ", logflag)
	_ = logger.Output(calldepth, fmt.Sprintln(message...))
	//logger.Println(message...)
}

func (slog *LogFile) CalledWarn(calldepth int, message ...interface{}) {
	if Level&r_debug == 0 {
		return
	}
	slog.sinit.Do(slog.init)
	slog.mu.Lock()
	defer slog.mu.Unlock()
	f, err := slog.GetFile()
	if err != nil {
		StdError("fail to create log file!")
	}
	defer f.Close()
	logger := log.New(f, "Warn: ", logflag)
	_ = logger.Output(calldepth, fmt.Sprintln(message...))
}

func (slog *LogFile) CalledError(calldepth int, message ...interface{}) {
	if Level&r_debug == 0 {
		return
	}
	slog.sinit.Do(slog.init)
	slog.mu.Lock()
	defer slog.mu.Unlock()
	f, err := slog.GetFile()
	if err != nil {
		StdError("fail to create log file!")
	}
	defer f.Close()
	logger := log.New(f, "Error: ", logflag)
	_ = logger.Output(calldepth, fmt.Sprintln(message...))
}

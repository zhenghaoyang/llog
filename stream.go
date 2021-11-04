package llog

import (
	"bytes"
	"fmt"
	"log"
	"runtime"
	"strconv"
)

//OFF、FATAL、ERROR、WARN、INFO、DEBUG、ALL
const (
	r_debug = 1 << iota
	r_info
	r_warn
	r_error
	r_fatal
)
const (
	L_All   = r_debug | r_info | r_warn | r_error | r_fatal
	L_Debug = r_debug | r_info | r_warn | r_error | r_fatal
	L_Info  = r_info | r_warn | r_error | r_fatal
	L_Warn  = r_warn | r_error | r_fatal
	L_Error = r_error | r_fatal
	L_Fatal = r_fatal
	L_Off   = 0
)

const logflag = log.Ldate | log.Ltime | log.Llongfile
const goroutineid = false

var Level int = L_All

func GetGoroutineID() uint64 {
	b := make([]byte, 64)
	runtime.Stack(b, false)
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

func clolorprint(messages interface{}, color int) {
	b := 40
	f := color //30 黑 31 红 32 绿 33 黄 34 蓝 35 紫红 36 青蓝  37 白
	fmt.Printf("%c[%d;%d;%dm%s%c[0m", 27, 0, b, f, messages, 27)
}

func StdDebug(message ...interface{}) {
	var buf bytes.Buffer
	if Level&r_debug == 0 {
		return
	}
	logger := log.New(&buf, "Debug: ", logflag)
	if goroutineid {
		pm := []interface{}{"goroutine id (", GetGoroutineID(), ")"}
		pm = append(pm, message...)
		_ = logger.Output(2, fmt.Sprintln(pm...))
	} else {
		_ = logger.Output(2, fmt.Sprintln(message...))
	}
	clolorprint(buf.String(), 35)
}
func StdInfo(message ...interface{}) {
	var buf bytes.Buffer
	if Level&r_info == 0 {
		return
	}
	logger := log.New(&buf, "Info: ", logflag)
	_ = logger.Output(2, fmt.Sprintln(message...))
	clolorprint(buf.String(), 32)
}
func StdWarn(message ...interface{}) {
	var buf bytes.Buffer
	if Level&r_info == 0 {
		return
	}
	logger := log.New(&buf, "Warn: ", logflag)
	_ = logger.Output(2, fmt.Sprintln(message...))
	clolorprint(buf.String(), 33)
}
func StdError(message ...interface{}) {
	var buf bytes.Buffer
	if Level&r_error == 0 {
		return
	}
	logger := log.New(&buf, "Error: ", logflag)
	if goroutineid {
		pm := []interface{}{"goroutine id (", GetGoroutineID(), ")"}
		pm = append(pm, message...)
		_ = logger.Output(2, fmt.Sprintln(pm...))
	} else {
		_ = logger.Output(2, fmt.Sprintln(message...))
	}
	clolorprint(buf.String(), 31)
}
func StdCalledDebug(calldepth int, message ...interface{}) {
	var buf bytes.Buffer
	if Level&r_debug == 0 {
		return
	}
	logger := log.New(&buf, "Debug: ", logflag)
	_ = logger.Output(calldepth, fmt.Sprintln(message...))
	clolorprint(buf.String(), 35)
}
func StdCalledInfo(calldepth int, message ...interface{}) {
	var buf bytes.Buffer
	if Level&r_info == 0 {
		return
	}
	logger := log.New(&buf, "Info: ", logflag)
	_ = logger.Output(calldepth, fmt.Sprintln(message...))
	clolorprint(buf.String(), 32)
}
func StdCalledWarn(calldepth int, message ...interface{}) {
	var buf bytes.Buffer
	if Level&r_info == 0 {
		return
	}
	logger := log.New(&buf, "Warn: ", logflag)
	_ = logger.Output(calldepth, fmt.Sprintln(message...))
	clolorprint(buf.String(), 33)
}
func StdCalledError(calldepth int, message ...interface{}) {
	var buf bytes.Buffer
	if Level&r_error == 0 {
		return
	}
	logger := log.New(&buf, "Error: ", logflag)
	_ = logger.Output(calldepth, fmt.Sprintln(message...))
	clolorprint(buf.String(), 31)
}

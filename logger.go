package logger

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type outputType int

const (
	NotPrint  outputType = 0 // not to print to console
	OnlyPrint outputType = 1 // only print to console
	WriteLog  outputType = 2 // print to console and write to log file
)

var (
	Trace *log.Logger
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger

	SeparateLogByDay bool = true
)

type logSetting struct {
	opType outputType
	prefix string
	flags  int
}

var (
	file               *os.File
	path               string
	filename           string = "now.log"
	separateLogByDayOn bool   = false

	trace logSetting = logSetting{
		opType: NotPrint,
		prefix: "TRACE ",
		flags:  log.Ldate | log.Ltime,
	}
	info logSetting = logSetting{
		opType: OnlyPrint,
		prefix: "INFO ",
		flags:  log.Ldate | log.Ltime,
	}
	warn logSetting = logSetting{
		opType: OnlyPrint,
		prefix: "WARN ",
		flags:  log.Ldate | log.Ltime | log.Llongfile,
	}
	err logSetting = logSetting{
		opType: WriteLog,
		prefix: "ERROR ",
		flags:  log.Ldate | log.Ltime | log.Llongfile,
	}
)

// Init log file folder path
func Init(folderPath string) error {
	if len(folderPath) == 0 || folderPath[len(folderPath)-1:] != "/" {
		folderPath = folderPath + "/"
	}

	err := setFolder(folderPath)
	if err != nil {
		return err
	}

	file, err = openFile(folderPath)
	if err != nil {
		return err
	}

	path = folderPath

	setLogSetting()

	if !separateLogByDayOn {
		go separateLog()
		separateLogByDayOn = true
	}

	return nil
}

// SetTraceFlags
func SetTraceFlags(typ outputType, prefix string, flags int) {
	trace = logSetting{
		opType: typ,
		prefix: prefix,
		flags:  flags,
	}
	Trace = setFlags(trace)
}

// SetInfoFlags
func SetInfoFlags(typ outputType, prefix string, flags int) {
	info = logSetting{
		opType: typ,
		prefix: prefix,
		flags:  flags,
	}
	Info = setFlags(info)
}

// SetWarnFlags
func SetWarnFlags(typ outputType, prefix string, flags int) {
	warn = logSetting{
		opType: typ,
		prefix: prefix,
		flags:  flags,
	}
	Warn = setFlags(warn)
}

// SetErrorFlags
func SetErrorFlags(typ outputType, prefix string, flags int) {
	err = logSetting{
		opType: typ,
		prefix: prefix,
		flags:  flags,
	}
	Error = setFlags(err)
}

func openFile(folderPath string) (*os.File, error) {
	return os.OpenFile(folderPath+filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
}

func setFolder(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.MkdirAll(path, 0666)
		return err
	}

	return nil
}

func setFlags(setting logSetting) *log.Logger {
	if setting.opType == 0 {
		return log.New(ioutil.Discard, setting.prefix, setting.flags)

	} else if setting.opType == 1 {
		return log.New(os.Stdout, setting.prefix, setting.flags)

	} else if setting.opType == 2 {
		if file != nil {
			return log.New(io.MultiWriter(file, os.Stderr), setting.prefix, setting.flags)
		}
	}

	return log.New(os.Stdout, setting.prefix, setting.flags)
}

func setLogSetting() {
	Trace = setFlags(trace)
	Info = setFlags(info)
	Warn = setFlags(warn)
	Error = setFlags(err)
}

func separateLog() {
	now := time.Now()
	tomorrow := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.Local)

	time.Sleep(time.Duration(tomorrow.Unix()-now.Unix()) * time.Second)

	if SeparateLogByDay {
		file.Close()
		date := fmt.Sprintf("%d-%02d-%02d.log", now.Year(), now.Month(), now.Day())
		err := os.Rename(path+filename, path+date)
		if err == nil {
			fmt.Println(err.Error())
			file, err = openFile(path)
			if err != nil {
				fmt.Println(err.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
	}

	separateLog()
}

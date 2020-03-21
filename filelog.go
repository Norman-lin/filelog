//This package allows to write logs to your self-defining file.use function NewFileLogger to add log level(which log level you want to see), your path and filename,
//The interface DEBUG,INFO,WARNING,ERROR,FATAL allow you to defining your own log contents.

package filelog

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
)

type LogLevel uint16

const(
	DEBUG LogLevel = iota
	TRACE
	INFO
	WARNING
	ERROR
	FATAL
)

type FileLogger struct{
	Level LogLevel
	filePath string //日志文件保存的路径
	fileName string //日志文件保存的文件名
	fileObj *os.File
}

func parseLogLevel(s string)LogLevel{
	s = strings.ToLower(s)
	switch s{
	case "debug":
		return DEBUG
	case "trace":
		return TRACE
	case "info":
		return INFO
	case "warning":
		return WARNING
	case "error":
		return ERROR
	case "fatal":
		return FATAL
	default:
		return DEBUG
	}
}

func unParseLogLevel(level LogLevel)string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case TRACE:
		return "TRACE"
	case INFO:
		return "INFO"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "DEBUG"
	}
}

func NewFileLogger(levelStr,fp,fn string)*FileLogger{
	logLevel := parseLogLevel(levelStr)
	file := &FileLogger{
		Level:    logLevel,
		filePath: fp,
		fileName: fn,
	}
	err := file.initFile() //按照文件路径和文件名将文件打开
	if err != nil{
		panic(err)
	}
	return file
}

func (f *FileLogger) initFile()(error){
	fullFileName := path.Join(f.filePath,f.fileName)
	fileObj,err := os.OpenFile(fullFileName,os.O_APPEND|os.O_CREATE|os.O_WRONLY,0644)
    if err != nil{
    	fmt.Println("open file failed,err:",err)
    	return err
	}
	//日志文件都打开
	f.fileObj = fileObj
	return nil
}

func (f *FileLogger)Close(){
	f.fileObj.Close()
}

// getinfo get name,path,line
func getInfo(skip int)(funcName,filePath string,lineNum int){
	pc,filePath,lineNum,ok := runtime.Caller(skip)
	if !ok {
		fmt.Printf("runtime.Caller()failed\n")
		return
	}
	funcName = runtime.FuncForPC(pc).Name()
	return

}

func (f *FileLogger) enable (logLevel LogLevel) bool{
	return f.Level <= logLevel
}

func (f *FileLogger) log(level LogLevel,msg string){
	if f.enable(level) {
		now := time.Now().Format("2006-01-02 15:04:05")
		funcName, fileName, lineNum := getInfo(3)
		fmt.Fprintf(f.fileObj,"[%s][%s][%s:%s]:%d:%s\n", now, unParseLogLevel(level), funcName, fileName, lineNum, msg)
	}
}

func (f *FileLogger) Debug (s string) {
		f.log(DEBUG,s)
}

func (f *FileLogger) Info (s string) {
		f.log(INFO,s)
}

func (f *FileLogger) Warning (s string) {
		f.log(WARNING,s)
}

func (f *FileLogger) Error (s string) {
		f.log(ERROR,s)
}

func (f *FileLogger) Fatal (s string) {
		f.log(FATAL,s)
}
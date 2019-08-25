package writer

import (
	"fmt"
	"os"
	"strings"
	"time"
	"util/filesystem"
	"util/logger/writers"
)

type FileWriter struct {
	file *os.File
}

var DefaultPath = "../../../storage/logs"

func (writer *FileWriter) write(level string, data interface{}) error {
	timestamp := time.Now().Unix()
	if writer.file == nil {
		if err := writer.SetFile(fmt.Sprintf("%s/%s/%s.log", strings.TrimRight(DefaultPath, "/"), level, time.Unix(timestamp, 0).Format("2006-01-02"))); err != nil {
			panic(err)
		}
	}

	defer writer.file.Close()

	now := time.Unix(timestamp, 0).Format("2006-01-02 03:04:05 PM")
	str := fmt.Sprintf("%s [%s] %v\n", now, level, data)

	_, err := writer.file.Write([]byte(str))
	return err
}

func (writer *FileWriter) SetFile(filename string) (err error) {
	if writer.file != nil {
		writer.file.Close()
	}

	dir := filesystem.Dir(filename)
	if exist, err := filesystem.PathExists(dir); err != nil {
		return err
	} else if !exist {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}

	writer.file, err = os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModePerm)
	return err
}

func (writer FileWriter) Emergency(message interface{}, context interface{}) error {
	return writer.write(writers.EmergencyLevel, message)
}

func (writer FileWriter) Alert(message interface{}, context interface{}) error {
	return writer.write(writers.AlertLevel, message)
}

func (writer FileWriter) Info(message interface{}, context interface{}) error {
	return writer.write(writers.InfoLevel, message)
}

func (writer FileWriter) Debug(message interface{}, context interface{}) error {
	return  writer.write(writers.DebugLevel, message)
}

func (writer FileWriter) Error(message interface{}, context interface{}) error {
	return  writer.write(writers.ErrorLevel, message)
}

func (writer FileWriter) Critical(message interface{}, context interface{}) error {
	return writer.write(writers.CriticalLevel, message)
}

func (writer FileWriter) Warning(message interface{}, context interface{}) error {
	return writer.write(writers.WarningLevel, message)
}

func (writer FileWriter) Notice(message interface{}, context interface{}) error {
	return writer.write(writers.NoticeLevel, message)
}
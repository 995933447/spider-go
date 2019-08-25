package writer

import (
	"fmt"
	"time"
	 "util/logger/writers"
)

type ConsoleWriter struct {

}

func write(level string, message interface{})  {
	message = fmt.Sprintf("%s [%s] %v\n", time.Unix(time.Now().Unix(), 0).Format("2006-01-02 03:04:05"), level, message)
	var b, d, f int
	switch level {
		case writers.EmergencyLevel:
			b, d, f = 4, 40,31
		case writers.AlertLevel:
			b, d, f = 1, 40, 32
	    case writers.DebugLevel:
	    	b, d, f = 1, 40, 36
	    case writers.InfoLevel:
	    	b, d, f = 1, 40, 34
	    case writers.ErrorLevel:
	    	b, d, f = 1, 40, 35
	    case writers.WarningLevel:
	    	b, d, f = 1, 40, 33
	    case writers.NoticeLevel:
			b, d, f = 1, 40, 37
		case writers.CriticalLevel:
			b, d, f = 1, 40,31
	}

	fmt.Printf(" %c[%d;%d;%dm %s%c[0m\n", 0x1B, d, b, f, message, 0x1B)
}

func (writer ConsoleWriter) Emergency(message interface{}, context interface{}) error {
	write(writers.EmergencyLevel, message)
	return nil
}

func (writer ConsoleWriter) Alert(message interface{}, context interface{}) error {
	write(writers.AlertLevel, message)
	return nil
}

func (writer ConsoleWriter) Info(message interface{}, context interface{}) error {
	write(writers.InfoLevel, message)
	return nil
}

func (writer ConsoleWriter) Error(message interface{}, context interface{}) error {
	write(writers.ErrorLevel, message)
	return nil
}

func (writer ConsoleWriter) Critical(message interface{}, context interface{}) error {
	write(writers.CriticalLevel, message)
	return nil
}

func (writer ConsoleWriter) Warning(message interface{}, context interface{}) error {
	write(writers.WarningLevel, message)
	return nil
}

func (writer ConsoleWriter) Notice(message interface{}, context interface{}) error {
	write(writers.NoticeLevel, message)
	return nil
}

func (writer ConsoleWriter) Debug(message interface{}, context interface{}) error {
	write(writers.DebugLevel, message)
	return nil
}

//颜色表
func colorList() {
	for b := 40; b <= 47; b++ { // 背景色彩 = 40-47
		for f := 30; f <= 37; f++ { // 前景色彩 = 30-37
			for d := range []int{0, 1, 4, 5, 7, 8} { // 显示方式 = 0,1,4,5,7,8
				fmt.Printf(" %c[%d;%d;%dm(f=%d,b=%d,d=%d)%c[0m ", 0x1B, d, b, f, f, b, d, 0x1B)
			}
			fmt.Println("")
		}
		fmt.Println("")
	}
}

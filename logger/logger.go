package logger

import (
	"fmt"
	"log"
	"runtime"
	"os"
	"math"
	"path"
	"time"
	"JustServer/utils/simplejson"
	"errors"
)

/**
	logger builder const
 */
const (
	chunk_size int = 4000;
	top_left_conner string = "╔";
	button_left_conner string = "╚";
	middle_conner string = "╟";
	horizontal_double_line string = "║";
	double_divider string = "════════════════════════════════════════════";
	single_divider string = "────────────────────────────────────────────";
	top_border string = top_left_conner + double_divider + double_divider;
	bottom_border string = button_left_conner + double_divider + double_divider;
	middle_border string = middle_conner + single_divider + single_divider;
)

/**
	logger level
 */
const (
	debug string = "[ debug ]"
	release string = "[release]"
	err string = "[ error ]"
	wtf string = "[  wtf  ]"
)

var D_level int = 0
var R_level int = 1
var E_level int = 2
var W_level int = 3

var Global_Logger, Global_Error = New("", "global", false)

type Logger struct {
	maxLevel      int
	moreMsg       bool
	console       bool
	logFile       *os.File
	logger        *log.Logger
	consoleLogger *log.Logger
}

func (logger *Logger)Console(console bool) {
	logger.console = console
}

func New(pathname string, module string, moreMsg bool) (*Logger, error) {
	logger := new(Logger)
	var Blogger *log.Logger
	var Clogger *log.Logger
	if pathname != "" {
		now := time.Now()
		filename := fmt.Sprintf("%d%02d%02d_%02d_%02d_%02d.log",
			now.Year(),
			now.Month(),
			now.Day(),
			now.Hour(),
			now.Minute(),
			now.Second())
		f, err := os.Create(path.Join(pathname, module + "-" + filename))
		if err != nil {
			return nil, err
		}
		logger.logFile = f
		Blogger = log.New(f, "", log.LstdFlags)
		Clogger = log.New(os.Stdout, "", log.LstdFlags)
	}else {
		Blogger = log.New(os.Stdout, "", log.LstdFlags)
		Clogger = Blogger
	}

	logger.moreMsg = moreMsg
	logger.logger = Blogger
	logger.consoleLogger = Clogger
	logger.console = true
	return logger, nil
}

func (logger *Logger)printer(level string, format string, a ...interface{}) {
	logger.logTop()

	levelString := getLevel(level)

	logger.logHeader(levelString)

	msg := fmt.Sprintf(format, a...)
	// msg
	length := len(msg)

	if length < chunk_size {
		logger.logContent(msg)
		logger.logBottom()
		return
	}

	for i := 0; i < length / chunk_size; i++ {
		count := math.Min(float64(length - 1), float64(chunk_size))
		logger.logContent(msg[i:int(count)])
	}
	logger.logBottom()
}

func (logger *Logger)logPrinter(msg string) {
	if logger.console && logger.consoleLogger != logger.logger {
		logger.consoleLogger.Println(msg)
	}
	logger.logger.Println(msg);
}

func (logger *Logger)logTop() {
	logger.logPrinter(top_border)
}

func (logger *Logger)logBottom() {
	logger.logPrinter(bottom_border)
}

func (logger *Logger)logHeader(levelString string) {
	// 4 floor
	pc, _, line, _ := runtime.Caller(4)

	//	checkError(ok)
	logger.logPrinter(horizontal_double_line)
	// left msg
	left := levelString + ". " + runtime.FuncForPC(pc).Name()

	len := len(left)

	if len < chunk_size {
		logger.logPrinter("║ " + left)
	}

	for i := 0; i < len / chunk_size; i++ {
		count := math.Min(float64(len - 1), float64(chunk_size))
		logger.logPrinter("║ " + left[i:int(count)])
	}

	logger.logger.Println("║ " + "line :", line)
}

func (logger *Logger)logContent(msg string) {
	logger.logPrinter(horizontal_double_line + "  " + msg)

}

func (logger *Logger)logDivider() {
	logger.logPrinter(middle_border)
}

func (logger *Logger)R(format string, a ...interface{}) {
	logger.log("R", format, a...)
}

func (logger *Logger)D(format string, a ...interface{}) {
	logger.log("D", format, a...)
}

func (logger *Logger)E(format string, a ...interface{}) {
	logger.log("E", format, a...)
}

func (logger *Logger)WTF(format string, a ...interface{}) {
	logger.log("WTF", format, a...)
}

func (logger *Logger)json(json string) error {
	js, jsErr := simplejson.NewJson([]byte(json))

	if jsErr != nil {
		return errors.New("json error")
	}
	jsMap, err := js.Map()

	if err != nil {
		return errors.New("json conver to map error")
	}

	for k, v := range jsMap {
		logger.logPrinter("key :" + k + " " + "val :" + fmt.Sprint(v))
	}
	return nil
}

func (logger *Logger)log(level string, format string, a ...interface{}) {
	if getLevelNum(level) > logger.maxLevel {
		if logger.moreMsg {
			logger.printer(level, format, a...)
		}else {
			msg := fmt.Sprintf(format, a...)
			logger.logPrinter(getLevel(level) + " : " + msg)
		}
	}
}

func getLevel(level string) string {
	levelString := ""
	switch level {
	case "D":
		levelString = debug
	case "R":
		levelString = release
	case "E":
		levelString = err
	case "WTF":
		levelString = wtf
	default:
		levelString = debug
	}
	return levelString
}

func getLevelNum(level string) int {
	switch level {
	case "D":
		return D_level
	case "WTF":
		return W_level
	case "R":
		return R_level
	case "E":
		return E_level
	default:
		return W_level
	}
	return W_level
}

func (logger *Logger)SetMaxLevel(Level int) {
	if Level <= 3 && Level >= 0 {
		logger.maxLevel = Level
	}
}

func (logger *Logger)Close() error {
	return logger.logFile.Close()
}



package gologs

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// LogLevel 定义日志等级
type LogLevel int

var (
	// config 实例在程序运行期间保持唯一
	Log *Logger = NewLogger(Warn)
)

const (
	Debug     LogLevel = 10
	Info      LogLevel = 20
	Hint      LogLevel = 22
	Important LogLevel = 24
	Warn      LogLevel = 30
	Error     LogLevel = 40
)

var Levels = map[LogLevel]string{
	Debug:     "Debug",
	Info:      "Info",
	Error:     "Error",
	Warn:      "Warn",
	Hint:      "Hint",
	Important: "Important",
}

// 颜色部分
var defaultColor = func(s string) string { return s }
var DefaultColorMap = map[LogLevel]func(string) string{
	Debug:     Yellow,
	Error:     RedBold,
	Info:      Cyan,
	Hint:      CyanBold,
	Warn:      YellowBold,
	Important: PurpleBold,
}

var DefaultFormatterMap = map[LogLevel]string{
	Debug:     "[Debug] %s \n",
	Warn:      "[Warn] %s \n",
	Error:     "[Error] %s \n",
	Info:      "[-] %s {{suffix}}\n",
	Hint:      "[+] %s {{suffix}}\n",
	Important: "[*] %s {{suffix}}\n",
}

func (l LogLevel) Name() string {
	if name, ok := Levels[l]; ok {
		return name
	} else {
		return strconv.Itoa(int(l))
	}
}

func (l LogLevel) Formatter() string {
	if formatter, ok := DefaultFormatterMap[l]; ok {
		return formatter
	} else {
		return "[" + l.Name() + "] %s"
	}
}

func (l LogLevel) Color() func(string) string {
	if f, ok := DefaultColorMap[l]; ok {
		return f
	} else {
		return defaultColor
	}
}

func AddLevel(level LogLevel, name string, opts ...interface{}) {
	Levels[level] = name
	for _, opt := range opts {
		switch opt := opt.(type) {
		case string:
			DefaultFormatterMap[level] = opt
		case func(string) string:
			DefaultColorMap[level] = opt
		}
	}
}

func NewLogger(level LogLevel) *Logger {
	log := &Logger{
		Level:     level,
		Color:     false,
		LogToFile: false,
		writer:    os.Stdout,
		levels:    Levels,
		formatter: DefaultFormatterMap,
		colorMap:  DefaultColorMap,
		SuffixFunc: func() string {
			return ", " + getCurtime()
		},
		PrefixFunc: func() string {
			return ""
		},
	}

	return log
}

// NewFileLogger create a pure file logger
func NewFileLogger(filename string) (*Logger, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	log := &Logger{
		Level:     Warn,
		writer:    file,
		formatter: DefaultFormatterMap,
		levels:    Levels,
		SuffixFunc: func() string {
			return ", " + getCurtime()
		},
		PrefixFunc: func() string {
			return ""
		},
	}
	return log, nil
}

// Logger 结构体，包含日志级别、日志文件等配置
type Logger struct {
	Quiet       bool // is enable Print
	Clean       bool // is enable Console()
	Level       LogLevel
	Color       bool
	LogToFile   bool
	LogFileName string
	SuffixFunc  func() string
	PrefixFunc  func() string

	mu        sync.RWMutex
	muf       sync.Mutex
	logFile   *os.File
	writer    io.Writer
	levels    map[LogLevel]string
	formatter map[LogLevel]string
	colorMap  map[LogLevel]func(string) string
}

func (log *Logger) SetQuiet(q bool) {
	log.Quiet = q
}

func (log *Logger) SetClean(c bool) {
	log.Clean = c
}

func (log *Logger) SetColor(c bool) {
	log.Color = c
}
func (log *Logger) SetIsLogToFile(l bool) {
	log.LogToFile = l
	if log.LogToFile {
		log.InitLogFile()
	} else {
		// 关闭原日志文件
		if log.logFile != nil {
			log.logFile.Close()
			log.logFile = nil
		}
	}
}

func (log *Logger) InitLogFile() {
	// 关闭原日志文件
	if log.logFile != nil {
		log.logFile.Close()
		log.logFile = nil
	}

	file, err := os.OpenFile(log.LogFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error: Set log file is error")
	}
	log.logFile = file
}

func (log *Logger) SetColorMap(cm map[LogLevel]func(string) string) {
	log.colorMap = cm
}

func (log *Logger) SetLevel(l LogLevel) {
	log.Level = l
}

func (log *Logger) SetOutput(w io.Writer) {
	log.writer = w
}

func (log *Logger) SetFile(filename string) {
	log.LogFileName = filename
}

func (log *Logger) SetFormatter(formatter map[LogLevel]string) {
	log.formatter = formatter
}

func (log *Logger) Console(s string) {
	if !log.Clean {
		fmt.Fprint(log.writer, s)
	}
}

func (log *Logger) Consolef(format string, s ...interface{}) {
	if !log.Clean {
		fmt.Fprintf(log.writer, format, s...)
	}
}

func (log *Logger) FConsolef(writer io.Writer, format string, s ...interface{}) {
	if !log.Clean {
		fmt.Fprintf(writer, format, s...)
	}
}

func (log *Logger) logInterface(writer io.Writer, level LogLevel, s interface{}) {
	log.mu.RLock()
	defer log.mu.RUnlock()
	if !log.Quiet && level >= log.Level {
		line := log.Format(level, s)
		if log.Color {
			fmt.Fprint(writer, log.SetLevelColor(level, line))
		} else {
			fmt.Fprint(writer, line)
		}

		// 写入到日志文件
		if log.LogToFile {
			log.writeToFile(line)
		}
	}
}

func (log *Logger) logInterfacef(writer io.Writer, level LogLevel, format string, s ...interface{}) {
	log.mu.RLock()
	defer log.mu.RUnlock()
	if !log.Quiet && level >= log.Level {
		line := log.Format(level, fmt.Sprintf(format, s...))
		if log.Color {
			fmt.Fprint(writer, log.SetLevelColor(level, line))
		} else {
			fmt.Fprint(writer, line)
		}

		// 写入到日志文件
		if log.LogToFile {
			log.writeToFile(line)
		}
	}
}

func (log *Logger) Log(level LogLevel, s interface{}) {
	log.logInterface(log.writer, level, s)
}

func (log *Logger) Logf(level LogLevel, format string, s ...interface{}) {
	log.logInterfacef(log.writer, level, format, s...)
}

func (log *Logger) FLogf(writer io.Writer, level LogLevel, s ...interface{}) {
	log.logInterface(writer, level, fmt.Sprintln(s...))
}

func (log *Logger) Important(s interface{}) {
	log.logInterface(log.writer, Important, s)
}

func (log *Logger) Importantf(format string, s ...interface{}) {
	log.logInterfacef(log.writer, Important, format, s...)
}

func (log *Logger) FImportantf(writer io.Writer, format string, s ...interface{}) {
	log.logInterfacef(writer, Important, format, s...)
}

func (log *Logger) Info(s interface{}) {
	log.logInterface(log.writer, Info, s)
}

func (log *Logger) Infof(format string, s ...interface{}) {
	log.logInterfacef(log.writer, Info, format, s...)
}

func (log *Logger) Hint(s interface{}) {
	log.logInterface(log.writer, Hint, s)
}

func (log *Logger) Hintf(format string, s ...interface{}) {
	log.logInterfacef(log.writer, Hint, format, s...)
}
func (log *Logger) FInfof(writer io.Writer, format string, s ...interface{}) {
	log.logInterfacef(writer, Info, format, s...)
}

func (log *Logger) Error(s interface{}) {
	log.logInterface(log.writer, Error, s)
}

func (log *Logger) Errorf(format string, s ...interface{}) {
	log.logInterfacef(log.writer, Error, format, s...)
}

func (log *Logger) FErrorf(writer io.Writer, format string, s ...interface{}) {
	log.logInterfacef(writer, Error, format, s...)
}

func (log *Logger) Warn(s interface{}) {
	log.logInterface(log.writer, Warn, s)
}

func (log *Logger) Warnf(format string, s ...interface{}) {
	log.logInterfacef(log.writer, Warn, format, s...)
}

func (log *Logger) FWarnf(writer io.Writer, format string, s ...interface{}) {
	log.logInterfacef(writer, Warn, format, s...)
}

func (log *Logger) Debug(s interface{}) {
	log.logInterface(log.writer, Debug, s)
}

func (log *Logger) Debugf(format string, s ...interface{}) {
	log.logInterfacef(log.writer, Debug, format, s...)
}

func (log *Logger) FDebugf(writer io.Writer, format string, s ...interface{}) {
	log.logInterfacef(writer, Debug, format, s...)
}

func (log *Logger) SetLevelColor(level LogLevel, line string) string {
	if c, ok := log.colorMap[level]; ok {
		return c(line)
	} else if c, ok := DefaultColorMap[level]; ok {
		return c(line)
	} else {
		return line
	}
}

func (log *Logger) Format(level LogLevel, s ...interface{}) string {
	var line string
	if f, ok := log.formatter[level]; ok {
		line = fmt.Sprintf(f, s...)
	} else if f, ok := DefaultFormatterMap[level]; ok {
		line = fmt.Sprintf(f, s...)
	} else {
		line = fmt.Sprintf("[%s] %s ", append([]interface{}{level.Name()}, s...)...)
	}
	line = strings.Replace(line, "{{suffix}}", log.SuffixFunc(), -1)
	line = strings.Replace(line, "{{prefix}}", log.PrefixFunc(), -1)
	return line
}

func (log *Logger) writeToFile(line string) {
	log.muf.Lock()
	defer log.muf.Unlock()
	// 检查 logFile 是否已初始化
	if log.logFile == nil {
		// 如果日志文件未初始化，打印警告
		fmt.Println("Error: Log file is not initialized.")
		return
	}

	// 写入日志到文件
	_, err := log.logFile.WriteString(line)
	if err != nil {
		// 写入失败时的错误处理
		fmt.Printf("Error writing to logfile: %s\n", err.Error())
		return
	}
}
func (log *Logger) Close(remove bool) {

	log.mu.Lock()
	defer log.mu.Unlock()

	// 关闭日志文件
	if log.logFile != nil {
		log.logFile.Close()
		log.logFile = nil
	}

	// 删除日志文件
	if remove {
		err := os.Remove(log.LogFileName)
		if err != nil {
			fmt.Printf("Error removing logfile: " + err.Error())
		}
	}
}

// 获取当前时间
func getCurtime() string {
	curtime := time.Now().Format("2006-01-02 15:04.05")
	return curtime
}

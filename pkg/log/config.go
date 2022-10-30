package log

import "go.uber.org/zap/zapcore"

// Levels zapcore level
var Levels = map[string]zapcore.Level{
	"":      zapcore.DebugLevel,
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"warn":  zapcore.WarnLevel,
	"error": zapcore.ErrorLevel,
	"fatal": zapcore.FatalLevel,
}

const (
	// OutputConsole console output
	OutputConsole = "console"
	// OutputFile file output
	OutputFile = "file"

	// EncoderTypeConsole log format console encoder
	EncoderTypeConsole = "console"
	// EncoderTypeJSON log format json encoder
	EncoderTypeJSON = "json"
)

var defaultConfig = &OutputConfig{
	Writer:     OutputConsole,
	Level:      "debug",
	Formatter:  EncoderTypeConsole,
	CallerSkip: 1,
}

// OutputConfig log output: console file remote
type OutputConfig struct {
	// Writer 日志输出端 (console, file)
	Writer string
	// FileConfig 日志文件配置，如果 Writer 为 file 则该配置不能为空
	FileConfig FileConfig `yaml:"file_config"`

	// Formatter 日志输出格式 (console, json)
	Formatter string

	// Level 日志级别 debug info error
	Level string

	// CallerSkip 控制 log 函数嵌套深度
	CallerSkip int `yaml:"caller_skip"`
}

// FileConfig 日志文件的配置
type FileConfig struct {
	// LogPath 日志路径
	LogPath string `yaml:"log_path"`
	// Filename 日志文件名
	Filename string `yaml:"filename"`
	// WriteMode 日志写入模式，1-同步，2-异步，3-极速(异步丢弃)
	WriteMode int `yaml:"write_mode"`
	// RollType 文件滚动类型，size-按大小分割文件，time-按时间分割文件，默认按大小分割
	RollType string `yaml:"roll_type"`
	// MaxAge 日志最大保留时间, 天
	MaxAge int `yaml:"max_age"`
	// MaxBackups 日志最大文件数
	MaxBackups int `yaml:"max_backups"`
	// Compress 日志文件是否压缩
	Compress bool `yaml:"compress"`
	// MaxSize 日志文件最大大小（单位MB）
	MaxSize int `yaml:"max_size"`

	// 以下参数按时间分割时才有效
	// TimeUnit 按时间分割文件的时间单位
	// 支持year/month/day/hour/minute, 默认为day
	TimeUnit TimeUnit `yaml:"time_unit"`
}

// WriteMode 日志写入模式，支持：1/2/3
type WriteMode int

const (
	// WriteSync 同步写
	WriteSync = 1
	// WriteAsync 异步写
	WriteAsync = 2
	// WriteFast 极速写(异步丢弃)
	WriteFast = 3
)

// 文件滚动类型配置字段
const (
	// RollBySize 按大小分割文件
	RollBySize = "size"
	// RollByTime 按时间分割文件
	RollByTime = "time"
)

// 常用时间格式
const (
	// TimeFormatMinute 分钟
	TimeFormatMinute = "%Y%m%d%H%M"
	// TimeFormatHour 小时
	TimeFormatHour = "%Y%m%d%H"
	// TimeFormatDay 天
	TimeFormatDay = "%Y%m%d"
	// TimeFormatMonth 月
	TimeFormatMonth = "%Y%m"
	// TimeFormatYear 年
	TimeFormatYear = "%Y"
)

// TimeUnit 文件按时间分割的时间单位，支持：minute/hour/day/month/year
type TimeUnit string

const (
	// Minute 按分钟分割
	Minute = "minute"
	// Hour 按小时分割
	Hour = "hour"
	// Day 按天分割
	Day = "day"
	// Month 按月分割
	Month = "month"
	// Year 按年分割
	Year = "year"
)

// Format 返回时间单位的格式字符串（c风格），默认返回day的格式字符串
func (t TimeUnit) Format() string {
	var timeFmt string
	switch t {
	case Minute:
		timeFmt = TimeFormatMinute
	case Hour:
		timeFmt = TimeFormatHour
	case Day:
		timeFmt = TimeFormatDay
	case Month:
		timeFmt = TimeFormatMonth
	case Year:
		timeFmt = TimeFormatYear
	default:
		timeFmt = TimeFormatDay
	}
	return "." + timeFmt
}

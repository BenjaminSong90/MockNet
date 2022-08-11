package logger

import (
	"fmt"
	"time"
)

var vLogger logFunc
var dLogger logFunc
var iLogger logFunc
var wLogger logFunc
var eLogger logFunc

func init() {

	vLogger = newLogFunc("V")
	dLogger = newLogFunc("D")
	iLogger = newLogFunc("I")
	wLogger = newLogFunc("W")
	eLogger = newLogFunc("E")

}

var logColor = map[string]string{
	"log_color_error": "bright_red",
	"log_color_debug": "cyan",

	"log_color_V": "cyan",
	"log_color_D": "green",
	"log_color_I": "magenta",
	"log_color_W": "yellow",
	"log_color_E": "red",
}

var colors = map[string]string{
	"reset":          "\033[0m",
	"black":          "\033[30m",
	"red":            "\033[31m",
	"green":          "\033[32m",
	"yellow":         "\033[33m",
	"blue":           "\033[34m",
	"magenta":        "\033[35m",
	"cyan":           "\033[36m",
	"white":          "\033[37m",
	"bold_black":     "\033[30;1m",
	"bold_red":       "\033[31;1m",
	"bold_green":     "\033[32;1m",
	"bold_yellow":    "\033[33;1m",
	"bold_blue":      "\033[34;1m",
	"bold_magenta":   "\033[35;1m",
	"bold_cyan":      "\033[36;1m",
	"bold_white":     "\033[37;1m",
	"bright_black":   "\033[30;2m",
	"bright_red":     "\033[31;2m",
	"bright_green":   "\033[32;2m",
	"bright_yellow":  "\033[33;2m",
	"bright_blue":    "\033[34;2m",
	"bright_magenta": "\033[35;2m",
	"bright_cyan":    "\033[36;2m",
	"bright_white":   "\033[37;2m",
}

type logFunc func(string)

func newLogFunc(prefix string) func(string) {
	color, clear := "", ""
	color = getLogColor(prefix)
	clear = colors["reset"]
	prefix = fmt.Sprintf("%-11s", prefix)

	return func(msg string) {
		timeString := time.Now().Format("2006/01/02 - 15:04:05")
		fmt.Printf("%s%s |%s |%s %s\n", color, timeString, prefix, msg, clear)
	}
}

func getLogColor(logName string) string {
	logColorName := fmt.Sprintf("log_color_%s", logName)
	colorName := logColor[logColorName]

	return colors[colorName]
}

func FormatPanicString(err error, info string) string {
	color := fmt.Sprintf("\033[0;%sm", "40;31")
	clear := colors["reset"]
	return fmt.Sprintf("%s %v %s %s \n", color, err, info, clear)
}

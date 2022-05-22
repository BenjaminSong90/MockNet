package utils

import (
	"fmt"
	"time"
)


var ErrorLogger logFunc
var DebugLogger logFunc

func init() {
	ErrorLogger = newLogFunc("error")
	DebugLogger = newLogFunc("debug")
}



var logColor = map[string]string{
	"log_color_error": "bright_red",
	"log_color_debug": "cyan",
}

var colors = map[string]string{
	"reset":          "0",
	"black":          "30",
	"red":            "31",
	"green":          "32",
	"yellow":         "33",
	"blue":           "34",
	"magenta":        "35",
	"cyan":           "36",
	"white":          "37",
	"bold_black":     "30;1",
	"bold_red":       "31;1",
	"bold_green":     "32;1",
	"bold_yellow":    "33;1",
	"bold_blue":      "34;1",
	"bold_magenta":   "35;1",
	"bold_cyan":      "36;1",
	"bold_white":     "37;1",
	"bright_black":   "30;2",
	"bright_red":     "31;2",
	"bright_green":   "32;2",
	"bright_yellow":  "33;2",
	"bright_blue":    "34;2",
	"bright_magenta": "35;2",
	"bright_cyan":    "36;2",
	"bright_white":   "37;2",
}


type logFunc func(string, ...interface{})

func newLogFunc(prefix string) func(string, ...interface{}) {
	color, clear := "", ""
	color = fmt.Sprintf("\033[%sm", getLogColor(prefix))
	clear = fmt.Sprintf("\033[%sm", colors["reset"])
	prefix = fmt.Sprintf("%-11s", prefix)

	return func(format string, v ...interface{}) {
		now := time.Now()
		timeString := fmt.Sprintf("%d:%d:%02d", now.Hour(), now.Minute(), now.Second())
		format = fmt.Sprintf("%s%s %s |%s %s", color, timeString, prefix, format, clear)
		fmt.Printf(format, v...)
	}
}


func getLogColor(logName string) string {
	logColorName := fmt.Sprintf("log_color_%s", logName)
	colorName := logColor[logColorName]

	return colors[colorName]
}



package helper

import (
	"time"

	"github.com/fatih/color"
	"github.com/briandowns/spinner"
)

var (
	SuccessLog = color.New(color.FgHiGreen).PrintfFunc()
	ErrorLog = color.New(color.FgHiRed).PrintfFunc()
	GreenString = color.HiGreenString
	RedString = color.HiRedString
	CyanString = color.HiCyanString
	Loader = spinner.New(spinner.CharSets[41], 100*time.Millisecond)
)

func RedPanic(err string) {
	panic(color.HiRedString(err))
}

func StartLoader() {
	Loader.Start()
}

func StopLoader() {
	Loader.Stop()
}
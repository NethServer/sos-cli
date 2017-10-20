/*
 * Copyright (C) 2017 Nethesis S.r.l.
 * http://www.nethesis.it - nethserver@nethesis.it
 *
 * This script is part of NethServer.
 *
 * NethServer is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License,
 * or any later version.
 *
 * NethServer is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with NethServer.  If not, see COPYING.
 *
 * author: Edoardo Spadoni <edoardo.spadoni@nethesis.it>
 */

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
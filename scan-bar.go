/*
 * Minio Client (C) 2014, 2015 Minio, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"fmt"
	"strings"

	"github.com/minio/mc/pkg/console"
	"github.com/minio/pb"
)

// Set up global cursor channel.
var cursorCh = cursorAnimate()

// Set up a constant width to follow.
const (
	messageWidth = 40
)

/******************************** Scan Bar ************************************/
// fixateScanBar truncates long text to fit within the terminal size.
func fixateScanBar(text string, width int) string {
	if len([]rune(text)) > width {
		// Trim text to fit within the screen
		trimSize := len([]rune(text)) - width + 3 //"..."
		if trimSize < len([]rune(text)) {
			text = "..." + text[trimSize:]
		}
	}
	return text
}

// TODO: create global totalTest count var.

// Progress bar function report objects being scaned.
type scanBarFunc func(string)

// scanBarFactory returns a progress bar function to report URL scanning.
func scanBarFactory() scanBarFunc {
	prevLineSize := 0
	prevMessage := ""
	//	testCount := 0
	termWidth, e := pb.GetTerminalWidth()
	if e != nil {
		console.Fatalln("Unable to get terminal size. Please use --quiet option.")
	}

	return func(message string) {
		scanPrefix := fmt.Sprintf("%s", message)
		//cmnPrefix := commonPrefix(message, prevMessage)
		eraseLen := prevLineSize // - len([]rune(cmnPrefix))
		if eraseLen < 1 {
			eraseLen = 0
		}
		if prevLineSize != 0 { // erase previous line
			console.PrintC("\r" + strings.Repeat(" ", eraseLen))
		}
		// TODO: need to find length of message and space accordingly.
		padding := messageWidth - len([]rune(scanPrefix))
		message = fixateScanBar(message, termWidth-len([]rune(scanPrefix))-1)
		barText := scanPrefix + strings.Repeat(" ", padding) + string(<-cursorCh)
		console.PrintC("\r" + barText)
		prevMessage = message
		prevLineSize = len([]rune(barText))
		//	testCount++
	}
}
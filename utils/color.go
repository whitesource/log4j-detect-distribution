package utils

import (
	"github.com/gookit/color"
)

var ColorsDisabled = false

var (
	// RGB colors does not work well in ssh terminals
	// using 256 colors instead
	red    = color.C256(160)
	orange = color.C256(214)
	yellow = color.C256(190)
	green  = color.C256(46)
	blue   = color.C256(33)
)

/* *** Private *** */

func colorText(c color.Color256, text string) string {
	if ColorsDisabled {
		return text
	}
	return c.Sprintf(text)
}

/* *** Public *** */

func MakeRedText(text string) string {
	return colorText(red, text)
}

func MakeOrangeText(text string) string {
	return colorText(orange, text)
}

func MakeYellowText(text string) string {
	return colorText(yellow, text)
}

func MakeGreenText(text string) string {
	return colorText(green, text)
}

func MakeBlueText(text string) string {
	return colorText(blue, text)
}

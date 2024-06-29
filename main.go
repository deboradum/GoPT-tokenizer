package main

import (
	"fmt"

	gogpt "github.com/deboradum/GoPT-tokenizer/tokenizer"
)

func main() {
	text := "Ｕｎｉｃｏｄｅ! 🅤🅝🅘🅒🅞🅓🅔‽ 🇺‌🇳‌🇮‌🇨‌🇴‌🇩‌🇪! 😄 The very name strikes fear and awe into the hearts of programmers worldwide. We all know we ought to “support Unicode” in our software (whatever that means—like using wchar_t for all the strings, right?). But Unicode can be abstruse, and diving into the thousand-page Unicode Standard plus its dozens of supplementary annexes, reports, and notes can be more than a little intimidating. I don’t blame programmers for still finding the whole thing mysterious, even 30 years after Unicode’s inception."

	tokens := gogpt.EncodeConversion(text)
	fmt.Println(tokens)

	stats := gogpt.GetStats(tokens)
	mostCommmon, count := gogpt.GetMostCommonBytePair(&stats)
	fmt.Println(mostCommmon, count)
}

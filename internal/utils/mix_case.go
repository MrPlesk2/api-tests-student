package utils

import "unicode"

// MixCase чередует верхний и нижний регистр букв, начиная с верхнего.
// Символы, не являющиеся буквами, остаются без изменений.
func MixCase(text string) string {
	letters := []rune(text)
	upper := true
	for i, letter := range letters {
		if !unicode.IsLetter(letter) {
			continue
		}
		if upper {
			letters[i] = unicode.ToUpper(letter)
		} else {
			letters[i] = unicode.ToLower(letter)
		}
		upper = !upper
	}
	return string(letters)
}

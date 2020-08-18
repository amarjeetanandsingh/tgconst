package text

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

type TransformFormat string

const (
	CamelCase  = "camelcase"
	LispCase   = "lispcase"
	PascalCase = "pascalcase"
	SnakeCase  = "snakecase"
	Mirror     = "mirror"

	LineLength = 60
)

// Consider 60 char width
// outputs "xyz" as "----  xyz  -------"
func CenterAlignedPadded(str, padWith string) string {
	space := "  "

	numOfPadEachSide := (LineLength - len(str) - len(space)*2) / 2
	if numOfPadEachSide <= 0 {
		return str
	}

	dashedString := strings.Repeat(padWith, numOfPadEachSide)
	return dashedString + space + str + space + dashedString
}

// NOTE: This code is copied from https://github.com/fatih/gomodifytags/blob/master/main.go#L353
func Transform(txt string, format TransformFormat) string {
	if len(txt) == 0 {
		return ""
	}

	splitted := Split(txt)

	switch format {
	case SnakeCase:
		var lowerSplitted []string
		for _, s := range splitted {
			lowerSplitted = append(lowerSplitted, strings.ToLower(s))
		}
		return strings.Join(lowerSplitted, "_")

	case LispCase:
		var lowerSplitted []string
		for _, s := range splitted {
			lowerSplitted = append(lowerSplitted, strings.ToLower(s))
		}
		return strings.Join(lowerSplitted, "-")

	case CamelCase:
		var titled []string
		for _, s := range splitted {
			titled = append(titled, strings.Title(s))
		}
		titled[0] = strings.ToLower(titled[0])
		return strings.Join(titled, "")

	case PascalCase:
		var titled []string
		for _, s := range splitted {
			titled = append(titled, strings.Title(s))
		}
		return strings.Join(titled, "")
	}

	return txt
}

// NOTE: This code is copied from https://github.com/fatih/camelcase/blob/master/camelcase.go
// Splitting rules
//  1) If string is not valid UTF-8, return it without splitting as
//     single item array.
//  2) Assign all unicode characters into one of 4 sets: lower case
//     letters, upper case letters, numbers, and all other characters.
//  3) Iterate through characters of string, introducing splits
//     between adjacent characters that belong to different sets.
//  4) Iterate through array of split strings, and if a given string
//     is upper case:
//       if subsequent string is lower case:
//         move last character of upper case string to beginning of
//         lower case string
func Split(src string) (entries []string) {
	// don't split invalid utf8
	if !utf8.ValidString(src) {
		return []string{src}
	}
	entries = []string{}
	var runes [][]rune
	lastClass := 0
	class := 0
	// split into fields based on class of unicode character
	for _, r := range src {
		switch true {
		case unicode.IsLower(r):
			class = 1
		case unicode.IsUpper(r):
			class = 2
		case unicode.IsDigit(r):
			class = 3
		default:
			class = 4
		}
		if class == lastClass {
			runes[len(runes)-1] = append(runes[len(runes)-1], r)
		} else {
			runes = append(runes, []rune{r})
		}
		lastClass = class
	}
	// handle upper case -> lower case sequences, e.g.
	// "PDFL", "oader" -> "PDF", "Loader"
	for i := 0; i < len(runes)-1; i++ {
		if unicode.IsUpper(runes[i][0]) && unicode.IsLower(runes[i+1][0]) {
			runes[i+1] = append([]rune{runes[i][len(runes[i])-1]}, runes[i+1]...)
			runes[i] = runes[i][:len(runes[i])-1]
		}
	}
	// construct []string from results
	for _, s := range runes {
		if len(s) > 0 {
			entries = append(entries, string(s))
		}
	}
	return
}

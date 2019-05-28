package StemmerRu

import (
	"runtime"
	"strings"
)

var (
	perfectiveGerund = [][]string{{`в`, `вши`, `вшись`}, {`ив`, `ивши`, `ившись`, `ыв`, `ывши`, `ывшись`}}
	adjective        = []string{`ее`, `ие`, `ые`, `ое`, `ими`, `ыми`, `ей`, `ий`, `ый`, `ой`, `ем`, `им`, `ым`, `ом`, `его`, `ого`, `ему`, `ому`, `их`, `ых`, `ую`, `юю`, `ая`, `яя`, `ою`, `ею`}
	participle       = [][]string{{`ем`, `нн`, `вш`, `ющ`, `щ`}, {`ивш`, `ывш`, `ующ`}}
	reflexive        = []string{`ся`, `сь`}
	verb             = [][]string{{`ла`, `на`, `ете`, `йте`, `ли`, `й`, `л`, `ем`, `н`, `ло`, `но`, `ет`, `ют`, `ны`, `ть`, `ешь`, `нно`}, {`ила`, `ыла`, `ена`, `ейте`, `уйте`, `ите`, `или`, `ыли`, `ей`, `уй`, `ил`, `ыл`, `им`, `ым`, `ен`, `ило`, `ыло`, `ено`, `ят`, `ует`, `уют`, `ит`, `ыт`, `ены`, `ить`, `ыть`, `ишь`, `ую`, `ю`}}
	noun             = []string{`а`, `ев`, `ов`, `ие`, `ье`, `е`, `иями`, `ями`, `ами`, `еи`, `ии`, `и`, `ией`, `ей`, `ой`, `ий`, `й`, `иям`, `ям`, `ием`, `ем`, `ам`, `ом`, `о`, `у`, `ах`, `иях`, `ях`, `ы`, `ь`, `ию`, `ью`, `ю`, `ия`, `ья`, `я`}
	superlative      = []string{`ейш`, `ейше`}
	derivational     = []string{`ост`, `ость`}

	vowels = `аеиоуыэюя`
)

func Stem(word string) string {

	word = strings.Replace(word, `ё`, `е`, -1)

	RVpos := getRVPart(word)

	if RVpos == -1 {
		return word
	}

	R1pos := getRNPart(word, 0)
	R2pos := getRNPart(word, R1pos)
	if R2pos < RVpos {
		R2pos = 0
	} else {
		R2pos -= RVpos
	}

	suffix := string([]rune(word)[RVpos:])
	prefix := string([]rune(word)[:RVpos])

	// Step 1
	suffix, isTrimmed := trimSuffix(suffix, perfectiveGerund[1], perfectiveGerund[0])
	if !isTrimmed {
		suffix, isTrimmed = trimSuffix(suffix, reflexive, nil)
		suffix, isTrimmed = trimAdjectival(suffix)
		if !isTrimmed {
			suffix, isTrimmed = trimSuffix(suffix, verb[1], verb[0])
			if !isTrimmed {
				suffix, _ = trimSuffix(suffix, noun, nil)
			}
		}
	}

	// Step 2
	suffix = strings.TrimSuffix(suffix, `и`)

	// Step 3
	if R2pos < len([]rune(suffix)) {
		R2suffix := string([]rune(suffix)[R2pos:])
		R2prefix := string([]rune(suffix)[:R2pos])
		R2suffix, _ = trimSuffix(R2suffix, derivational, nil)
		suffix = R2prefix + R2suffix
	}

	// Step 4
	suffix, isTrimmed = trimNN(suffix)
	if !isTrimmed {
		suffix, isTrimmed = trimSuffix(suffix, superlative, nil)
		if isTrimmed {
			suffix, _ = trimNN(suffix)
		} else {
			suffix = strings.TrimSuffix(suffix, `ь`)
		}
	}

	return prefix + suffix
}

func trimNN(word string) (string, bool) {
	if strings.HasSuffix(word, `нн`) {
		return strings.TrimSuffix(word, `нн`) + `н`, true
	}

	return word, false
}

func trimAdjectival(word string) (string, bool) {
	isTrimmedParticiple := false
	word, isTrimmedAdjective := trimSuffix(word, adjective, nil)
	if isTrimmedAdjective {
		word, isTrimmedParticiple = trimSuffix(word, participle[1], participle[0])
	}

	return word, isTrimmedAdjective || isTrimmedParticiple
}

func trimSuffix(word string, suffixes []string, suffixes2 []string) (string, bool) {
	for _, suffix := range suffixes {
		if strings.HasSuffix(word, suffix) {
			return strings.TrimSuffix(word, suffix), true
		}
	}
	if suffixes2 != nil {
		for _, suffix := range suffixes2 {
			if strings.HasSuffix(word, `а`+suffix) ||
				strings.HasSuffix(word, `я`+suffix) {
				return strings.TrimSuffix(word, suffix), true
			}
		}
	}

	return word, false
}

func isVowel(char rune) bool {
	return strings.Contains(vowels, string(char))
}

func getRVPart(word string) int {
	chars := []rune(word)
	for idx, char := range chars {
		if isVowel(char) {
			return idx + 1
		}
	}

	return -1
}

func getRNPart(word string, startPos int) int {
	chars := []rune(word)[startPos:]
	for idx, char := range chars {
		if idx+2 < len(chars) {
			if isVowel(char) && !isVowel(chars[idx+1]) {
				return startPos + idx + 2
			}
		}
	}

	return startPos
}

// Code from https://github.com/caneroj1/stemmer

// StemMultiple accepts a slice of strings and stems each of them.
func StemMultiple(words []string) (output []string) {
	output = make([]string, len(words))
	for idx, word := range words {
		output[idx] = Stem(word)
	}

	return
}

// StemMultipleMutate accepts a pointer to a slice of strings and stems them in place.
// It modifies the original slice.
func StemMultipleMutate(words *[]string) {
	for idx, word := range *words {
		(*words)[idx] = Stem(word)
	}
}

// StemConcurrent accepts a pointer to a slice of strings and stems them in place.
// It tries to offload the work into multiple threads. It makes no guarantees about
// the order of the stems in the modified slice.
func StemConcurrent(words *[]string) {
	CPUs := runtime.NumCPU()
	length := len(*words)
	output := make(chan string)
	partition := length / CPUs

	var CPU int
	for CPU = 0; CPU < CPUs; CPU++ {
		go func(strs []string) {
			for _, word := range strs {
				output <- Stem(word)
			}
		}((*words)[CPU*partition : (CPU+1)*partition])
	}

	// if there are leftover words, stem them now
	if length-(CPU)*partition > 0 {
		go func(strs []string) {
			for _, word := range strs {
				output <- Stem(word)
			}
		}((*words)[(CPU)*partition : length])
	}

	for idx := 0; idx < length; idx++ {
		(*words)[idx] = <-output
	}
}

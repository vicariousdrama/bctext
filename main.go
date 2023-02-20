package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
)

var (
	characterWidthMap map[int]int
	debug             bool
	align             bool
	nopad             bool
	blockclockip      string
	texttoshow        string
	eInkWidth         int
)

func init() {
	characterWidthMap = map[int]int{}
	eInkWidth = 128
	flag.StringVar(&blockclockip, "blockclockip", "21.21.21.21", "Blockclock IP Address")
	flag.StringVar(&texttoshow, "texttoshow", "This is sample output created with bctext by vicariousdrama", "Text to Show")
	flag.BoolVar(&align, "wordalign", false, "Avoid breaking word over panels, starting with next panel for each word")
	flag.BoolVar(&nopad, "nopadding", false, "Controls whether periods should be added to pad out the text in panels to edges. With no padding, data in panels is centered")
	flag.BoolVar(&debug, "debugmode", false, "Show data, dont send to blockclock")
	initCharacterWidthMap()
}
func getStringWidthInPixels(s string) (string, int) {
	width := 0
	s2 := ""
	l := len(s)
	for i := 0; i < l; i++ {
		ch := s[i]
		chi := int(ch)
		chs, chw := getCharacterWidthInPixels(chi)
		s2 += chs
		width += chw
	}
	return s2, width
}
func getCharacterWidthInPixels(chr int) (string, int) {
	// lookup in the map
	val, ok := characterWidthMap[chr]
	if ok {
		return string(chr), val
	}
	// all others will be defaulted to periods instead
	dVal := '.'
	val, ok = characterWidthMap[int(dVal)]
	// but convert spaces to 2 periods
	if chr == 32 {
		return strings.Repeat(string(dVal), 2), val * 2
	}
	return string(dVal), val
}
func fillAvailableSpaceWithPeriods(s string, maxWidth int) string {
	_, periodWidth := getCharacterWidthInPixels(int('.'))
	for true {
		_, currentWidth := getStringWidthInPixels(s)
		availableWidth := maxWidth - currentWidth
		if availableWidth > periodWidth {
			s += "."
		} else {
			break
		}
	}
	return s
}
func initCharacterWidthMap() {
	// character map setup (based on eink width of 128px)
	// letters A..Z
	for c := 'A'; c <= 'Z'; c++ {
		characterWidthMap[int(c)] = 25
	}
	// letter exceptions
	characterWidthMap['I'] = 13
	characterWidthMap['J'] = 20
	characterWidthMap['M'] = 32
	characterWidthMap['W'] = 32
	// numbers 0..9
	for c := '0'; c <= '9'; c++ {
		characterWidthMap[int(c)] = 25
	}
	// symbols
	characterWidthMap['+'] = 40
	characterWidthMap[','] = 13 // note: additional space follows first in series
	characterWidthMap['-'] = 18
	characterWidthMap['.'] = 12 // note: additional space follows first in series
	// no support for !, $, &, ', (, ), *, /, :, ;, <, =, >, ?, @, [, \, ], ^, _
	// some others get urlencoded so we can't use: ", <, >, [, ], %
}
func renderDebugOutput(slots *[]string, panelCount int) {
	slotMax := panelCount * 2
	slotCount := len(*slots)
	if slotCount < slotMax {
		slotMax = slotCount
	}
	fmt.Println("Debug results for this text string")
	fmt.Println("---------------------------------------------------------------------------------------")
	fmt.Println(" slot      over            under       url")
	for i := 0; i < slotMax && i < panelCount; i++ {
		otext := (*slots)[i]
		utext := ""
		if (i + panelCount) < slotMax {
			utext = (*slots)[i+panelCount]
		}
		url := fmt.Sprintf("http://%s/api/ou_text/%d/%s/%s", blockclockip, i, otext, utext)
		fmt.Println(fmt.Sprintf("%5d  %-14s  %-14s  %s", i, otext, utext, url))
	}
}
func renderToBlockclock(slots *[]string, panelCount int) {
	slotMax := panelCount * 2
	slotCount := len(*slots)
	if slotCount < slotMax {
		slotMax = slotCount
	}
	for i := 0; i < slotMax && i < panelCount; i++ {
		otext := (*slots)[i]
		utext := ""
		if (i + panelCount) < slotMax {
			utext = (*slots)[i+panelCount]
		}
		url := fmt.Sprintf("http://%s/api/ou_text/%d/%s/%s", blockclockip, i, otext, utext)
		// TODO: Graceful handling of errors here instead of fatal
		_, err := http.Get(url)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
func main() {
	// parse command line flags
	flag.Parse()
	pad := !(nopad)

	// default to debugmode if the blockclock ip was not given
	if blockclockip == "21.21.21.21" || blockclockip == "" {
		debug = true
	}

	// panels - info about the eInk displays
	panelCount := 7   //
	panelWidth := 128 // the width of the eInk display in pixels

	// slot setup
	var slotMax = panelCount * 2 // two rows for over / under text
	var slots []string           // start empty and add as needed
	var slotIndex int = 0
	slots = append(slots, "")

	// for local convenience
	convertedSpace, convertedSpaceWidth := getCharacterWidthInPixels(int(' '))

	// Text string to words preparation
	textToShow_UpperCased := strings.ToUpper(texttoshow) // lowercase not supported for over under text
	textToShow_Words := strings.Split(textToShow_UpperCased, " ")
	textToShow_WordsCount := len(textToShow_Words)

	// Process words of the string
	for textToShowWordIndex, textToShow_UpperCasedword := range textToShow_Words {
		wordConverted, wordWidth := getStringWidthInPixels(textToShow_UpperCasedword)
		wordLength := len(wordConverted)
		// RULE: Advance slot if using alignment mode
		if align && textToShowWordIndex > 0 {
			slotIndex += 1
			slots = append(slots, "")
		}
		// determine current slot contents
		currentSlotContent := slots[slotIndex]
		_, currentSlotWidth := getStringWidthInPixels(currentSlotContent)
		// determine available space of slot and row
		slotWidthAvailable := panelWidth - currentSlotWidth
		rowWidthAvailable := ((panelCount - (slotIndex % panelCount) - 1) * panelWidth) + slotWidthAvailable
		// track for two letters to avoid islanding single letter at end of slot
		twoletterwidth := wordWidth
		if wordLength >= 2 {
			_, twoletterwidth = getStringWidthInPixels(textToShow_UpperCasedword[0:2])
		}
		// RULE: if small word and not enough space, or any word not enough space in the row.
		if ((wordLength < 4) && (slotWidthAvailable < wordWidth) && (currentSlotWidth > 0)) ||
			(slotWidthAvailable < twoletterwidth) ||
			(rowWidthAvailable < wordWidth) {
			// advance to next slot
			if pad {
				slots[slotIndex] = fillAvailableSpaceWithPeriods(slots[slotIndex], panelWidth)
			}
			slotIndex += 1
			slots = append(slots, "")
			// RULE: dont break words across rows
			if (rowWidthAvailable < wordWidth) && (slotIndex%panelCount != 0) {
				// advance slots as needed for next row
				for {
					if pad {
						slots[slotIndex] = fillAvailableSpaceWithPeriods(slots[slotIndex], panelWidth)
					}
					slotIndex += 1
					slots = append(slots, "")
					if slotIndex%panelCount == 0 {
						break
					} else if pad {
						slots[slotIndex] = fillAvailableSpaceWithPeriods(slots[slotIndex], panelWidth)
					}
				}
			}
		} else {
			// if at end of row and already has content, insert a space
			if (((slotIndex + 1) % panelCount) == 0) && (currentSlotWidth > 0) && (rowWidthAvailable > convertedSpaceWidth) {
				fmt.Println("adding a space for end of row slot that has content")
				slots[slotIndex] = slots[slotIndex] + convertedSpace
			}
		}
		// process characters of this word into slots
		for charIndex := 0; charIndex < wordLength; charIndex++ {
			charCurrent := textToShow_UpperCasedword[charIndex]
			charCurrentInt := int(charCurrent)
			charString, charWidth := getCharacterWidthInPixels(charCurrentInt)
			currentSlotContent := slots[slotIndex]
			_, currentSlotWidth := getStringWidthInPixels(currentSlotContent)
			slotWidthAvailable := panelWidth - currentSlotWidth
			// if we do not fit
			if charWidth > slotWidthAvailable {
				// advance to next slot
				slotIndex += 1
				slots = append(slots, charString)
			} else {
				// append to current slot
				currentSlotContent = currentSlotContent + charString
				slots[slotIndex] = currentSlotContent
			}
		}
		// add a space if not at end of row and not last word
		if pad {
			slotIsNotLastPanel := ((slotIndex+1)%panelCount != 0)
			notLastWord := (textToShowWordIndex < textToShow_WordsCount-1)
			if slotIsNotLastPanel && notLastWord {
				// add whatever we can fit of a converted space
				for si := 0; si < len(convertedSpace); si++ {
					currentSlotContent = slots[slotIndex]
					_, currentSlotWidth = getStringWidthInPixels(currentSlotContent)
					slotWidthAvailable = panelWidth - currentSlotWidth
					convertedCharacter, convertedCharacterWidth := getCharacterWidthInPixels(int(convertedSpace[si]))
					if slotWidthAvailable > convertedCharacterWidth {
						slots[slotIndex] = slots[slotIndex] + convertedCharacter
					} else {
						break
					}
				}
			}
		}
	} // Done processing words

	// Pad slots as needed
	if pad {
		// fill periods for remaining space in last slot with content
		slots[slotIndex] = fillAvailableSpaceWithPeriods(slots[slotIndex], panelWidth)
		// adds periods to fill end of existing slots that trail with spaces
		for fillslot := 0; fillslot < slotIndex; fillslot++ {
			if strings.HasSuffix(slots[fillslot], convertedSpace) {
				slots[fillslot] = fillAvailableSpaceWithPeriods(slots[fillslot], panelWidth)
			}
		}
		// create additional slots to fill out slotMax
		for fillslot := slotIndex; fillslot < slotMax; fillslot++ {
			slots = append(slots, "............")
		}
	}

	// display output
	if debug {
		renderDebugOutput(&slots, panelCount)
	} else {
		renderToBlockclock(&slots, panelCount)
	}
}

package main

import (
        "flag"
        "fmt"
        "log"
        "net/http"
        "strings"
)

var (
        charWidthMap map[int]int
        debug bool
        align bool
        nopad bool
        blockclockip string
        texttoshow string
)
func init() {
        charWidthMap = map[int]int{}
        flag.StringVar(&blockclockip, "blockclockip", "21.21.21.21", "Blockclock IP Address")
        flag.StringVar(&texttoshow, "texttoshow", "This is sample output created with bctext by vicariousdrama", "Text to Show")
        flag.BoolVar(&align, "wordalign", false, "Avoid breaking word over panels, starting with next panel for each word")
        flag.BoolVar(&nopad, "nopadding", false, "Controls whether periods should be added to pad out the text in panels to edges. With no padding, data in panels is centered")
        flag.BoolVar(&debug, "debugmode", false, "Show data, dont send to blockclock")
}
func getstringwidth(s string) (string, int) {
        width := 0
        s2 := ""
        l := len(s)
        for i := 0; i < l; i++ {
                ch := s[i]
                chi := int(ch)
                chs, chw := getcharacterwidth(chi)
                s2 += chs
                width += chw
        }
        return s2, width
}
func getcharacterwidth(chr int) (string, int) {
        dVal := '.' // convert unmapped to a period
        val, ok := charWidthMap[chr]
        if ok {
                return string(chr), val
        }
        val, ok = charWidthMap[int(dVal)]
        if chr == 32 {
                return strings.Repeat(string(dVal),2), val*2
        }
        return string(dVal), val
}
func addperiodstoend(s string, m int) (string) {
        for true {
                _, w := getstringwidth(s)
                a := m - w
                if a > 12 { // period width
                        s = fmt.Sprintf("%s.",s)
                } else {
                        break
                }
        }
        return s
}
func main() {
        // parse command line flags
        flag.Parse()
        pad := !(nopad)

        // default to debugmode if the blockclock ip was not given
        if (blockclockip == "21.21.21.21" || blockclockip == "") {
                debug = true
        }

        // character map setup (based on eink width of 128px)
        // letters A..Z
        for c := 'A'; c <= 'Z'; c++ {
                charWidthMap[int(c)] = 25
        }
        // letter exceptions
        charWidthMap['I'] = 13
        charWidthMap['J'] = 20
        charWidthMap['M'] = 32
        charWidthMap['W'] = 32
        // numbers 0..9
        for c := '0'; c <= '9'; c++ {
                charWidthMap[int(c)] = 25
        }
        // symbols
        periodwidth := 12
        charWidthMap['+'] = 40
        charWidthMap[','] = 13          // note: additional space follows first in series
        charWidthMap['-'] = 18
        charWidthMap['.'] = periodwidth // note: additional space follows first in series
        // no support for !, $, &, ', (, ), *, /, :, ;, <, =, >, ?, @, [, \, ], ^, _
        // some others get urlencoded so we can't use: ", <, >, [, ], %

        // panels - the number of eink displays
        panelcount := 7

        // slot setup
        var slotmax = panelcount * 2    // two rows
        var slotwidthmax = 128          // pixels
        var slots []string              // our empty slot array, will add as needed and fill later
        var slotidx int = 0
        slots = append(slots, "",)

        tts := strings.ToUpper(texttoshow)      // the blockclock mini doesnt support lowercase for over under text
        ttswords := strings.Split(tts, " ")
        ttswordscount := len(ttswords)
        // process words of the string
        for ttswidx, ttsword := range ttswords {
                wordconverted, wordwidth := getstringwidth(ttsword)
                wordlen := len(wordconverted)
                if align && ttswidx > 0 {
                        slotidx += 1
                        slots = append(slots, "", )
                }
                currentslotcontent := slots[slotidx]
                _, currentslotwidth := getstringwidth(currentslotcontent)
                // track for two letters to avoid islanding single letter at end of slot
                twoletterwidth := wordwidth
                if wordlen >= 2 {
                        _, twoletterwidth = getstringwidth(ttsword[0:2])
                }
                slotavail := slotwidthmax - currentslotwidth
                rowavail := ((panelcount - (slotidx%panelcount) - 1) * slotwidthmax) + slotavail
                // if small word and not enough space, or any word not enough space in the row.
                if ((wordlen < 4) && (slotavail < wordwidth) && (currentslotwidth > 0)) ||
                (slotavail < twoletterwidth) ||
                (rowavail < wordwidth) {
                        // advance to next slot
                        if pad {
                                slots[slotidx] = addperiodstoend(slots[slotidx],slotwidthmax)
                        }
                        slotidx += 1
                        slots = append(slots, "", )
                        if (rowavail < wordwidth) && (slotidx%panelcount != 0) {
                                // advance slots as needed for next row
                                for {
                                        if pad {
                                                slots[slotidx] = addperiodstoend(slots[slotidx],slotwidthmax)
                                        }
                                        slotidx += 1
                                        slots = append(slots, "", )
                                        if slotidx%panelcount == 0 {
                                                break
                                        } else if pad {
                                                slots[slotidx] = addperiodstoend(slots[slotidx],slotwidthmax)
                                        }
                                }
                        }
                } else {
                        // if at end of row and already has content, insert a space
                        if (((slotidx+1)%panelcount)==0) && (currentslotwidth > 0) && (rowavail > periodwidth*2) {
                                slots[slotidx] = fmt.Sprintf("%s%s",slots[slotidx],"..")
                        }
                }
                // process characters of this word into slots
                for charidx := 0; charidx < wordlen; charidx ++ {
                        charcur := ttsword[charidx]
                        charcuri := int(charcur)
                        charstring, charwidth := getcharacterwidth(charcuri)
                        currentslotcontent := slots[slotidx]
                        _, currentslotwidth := getstringwidth(currentslotcontent)
                        slotavail := slotwidthmax - currentslotwidth
                        // if we aren't gonna fit
                        if charwidth > slotavail {
                                // advance to next slot
                                slotidx +=1
                                slots = append(slots, charstring, )
                        } else {
                                // append it
                                currentslotcontent = fmt.Sprintf("%s%s",currentslotcontent,charstring)
                                slots[slotidx] = currentslotcontent
                        }
                }
                // add a space if not at end of row and not last word
                if pad {
                        if ((slotidx+1)%panelcount != 0) && (ttswidx < ttswordscount -1) {
                                for si := 0; si < 2; si ++ {
                                        currentslotcontent = slots[slotidx]
                                        _, currentslotwidth = getstringwidth(currentslotcontent)
                                        slotavail = slotwidthmax - currentslotwidth
                                        if slotavail > periodwidth {
                                                slots[slotidx] = fmt.Sprintf("%s%s", slots[slotidx], ".")
                                        }
                                }
                        }
                }
        }
        if pad {
                // fill periods for remaining space in last slot with content
                slots[slotidx] = addperiodstoend(slots[slotidx], slotwidthmax)
                // adds periods to fill end of existing slots that trail with two periods
                for fillslot := 0; fillslot < slotidx; fillslot ++ {
                        if strings.HasSuffix(slots[fillslot], "..") {
                                slots[fillslot] = addperiodstoend(slots[fillslot], slotwidthmax)
                        }
                }
                // create additional slots to fill out slotmax
                for fillslot := slotidx; fillslot < slotmax; fillslot ++ {
                        slots = append(slots, "............", )
                }
        }


        // render to blockclock =======================================================================================
        slotcount := len(slots)
        if slotcount < slotmax {
                slotmax = slotcount
        }
        if debug {
                fmt.Println("Debug results for this text string")
                fmt.Println("---------------------------------------------------------------------------------------")
                fmt.Println(" slot      over            under       url")
        }
        for i := 0; i < slotmax && i < panelcount; i++ {
                otext := slots[i]
                utext := ""
                if (i+panelcount) < slotmax {
                        utext = slots[i+panelcount]
                }
                url := fmt.Sprintf("http://%s/api/ou_text/%d/%s/%s", blockclockip, i, otext, utext)
                if debug {
                        fmt.Println(fmt.Sprintf("%5d  %-14s  %-14s  %s", i, otext, utext, url))
                } else {
                        _, err := http.Get(url)
                        if err != nil {
                                log.Fatalln(err)
                        }
                }
        }
}

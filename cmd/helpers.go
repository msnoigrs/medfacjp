package cmd

import (
	"fmt"
	"os"
)

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(-1)
}

func VerifyDate(text string) (string, bool) {
	maxlen := 20
	rval := make([]rune, 0, maxlen)
	fStart := 0
	fEnd := 0
	fNen := 0
	fTuki := 0
	fHi := 0
	ai := 0
	for _, c := range text {
		if fStart == 0 {
			if c == '[' {
				fStart++
			}
		} else {
			if c == ']' {
				fEnd++
				break
			}
			if c != ' ' && c != '　' && c != '現' && c != '在' {
				rval = append(rval, c)
				ai++
				if ai == maxlen {
					break
				}
			}
			if c == '年' {
				fNen++
			} else if c == '月' {
				fTuki++
			} else if c == '日' {
				fHi++
			}
		}
	}
	rstr := string(rval)
	if fStart == 1 && fEnd == 1 && fNen == 1 && fTuki == 1 && fHi == 1 {
		return rstr, true
	}
	return rstr, false
}

// rbufを使いまわししているので、注意
func RemoveSpace(rbuf []rune, maxlen int, text string) string {
	ai := 0
	for _, c := range text {
		if c != ' ' && c != '　' {
			rbuf = append(rbuf, c)
			ai++
			if ai == maxlen {
				break
			}
		}
	}
	return string(rbuf)
}

package cmd

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"strconv"
)

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(-1)
}

var ByosyoTable = map[string]int{
	"一般":     1,
	"一般（感染）": 1,
	"療養":     1,
	"精神":     1,
	"介護":     1,
	"結核":     1,
	"特定":     1,
	"その他":    1,
}

var HyobouTable = map[string]int{
	"内":  1,
	"外":  1,
	"胃":  1,
	"消":  1,
	"小":  1,
	"小外": 1,
	"整外": 1,
	"皮ひ": 1,
	"こう": 1,
	"美外": 1,
	"ひ":  1,
	"皮":  1,
	"耳い": 1,
	"眼":  1,
	"麻":  1,
	"産婦": 1,
	"産":  1,
	"婦":  1,
	"精":  1,
	"歯":  1,
	"矯歯": 1,
	"小歯": 1,
	"歯外": 1,
	"リハ": 1,
	"放":  1,
	"病理": 1,
	"脳外": 1,
	"形外": 1,
	"循":  1,
	"神内": 1,
	"リウ": 1,
	"アレ": 1,
	"救命": 1,
	"心外": 1,
	"神":  1,
	"呼":  1,
	"呼内": 1,
	"呼外": 1,
	"透析": 1,
}
var PrefCodeTable = map[string]string{
	"北海道":  "01",
	"青森県":  "02",
	"岩手県":  "03",
	"宮城県":  "04",
	"秋田県":  "05",
	"山形県":  "06",
	"福島県":  "07",
	"茨城県":  "08",
	"栃木県":  "09",
	"群馬県":  "10",
	"埼玉県":  "11",
	"千葉県":  "12",
	"東京都":  "13",
	"神奈川県": "14",
	"新潟県":  "15",
	"富山県":  "16",
	"石川県":  "17",
	"福井県":  "18",
	"山梨県":  "19",
	"長野県":  "20",
	"岐阜県":  "21",
	"静岡県":  "22",
	"愛知県":  "23",
	"三重県":  "24",
	"滋賀県":  "25",
	"京都府":  "26",
	"大阪府":  "27",
	"兵庫県":  "28",
	"奈良県":  "29",
	"和歌山県": "30",
	"鳥取県":  "31",
	"島根県":  "32",
	"岡山県":  "33",
	"広島県":  "34",
	"山口県":  "35",
	"徳島県":  "36",
	"香川県":  "37",
	"愛媛県":  "38",
	"高知県":  "39",
	"福岡県":  "40",
	"佐賀県":  "41",
	"長崎県":  "42",
	"熊本県":  "43",
	"大分県":  "44",
	"宮崎県":  "45",
	"鹿児島県": "46",
	"沖縄県":  "47",
}

func GetTown(text string, wbuff []rune) string {
	wbuff = wbuff[:0]
	textrune := []rune(text)
	for _, c := range textrune {
		wbuff = append(wbuff, c)
		if c == '市' || c == '町' || c == '村' || c == '区' {
			break
		} else if c == '郡' {
			wbuff = wbuff[:0]
		}
	}
	return string(wbuff)
}

func GetPrefCode(postal string, address string, db *sql.DB, wbuff []rune) (string, bool) {
	pcode := SelectNumbers(postal, wbuff)
	if len(pcode) < 7 {
		fmt.Println("invalid postal?: " + postal + " " + pcode)
		return "", false
	}

	rows, err := db.Query(
		`SELECT PREF FROM KENALL WHERE ZIP = ?`, pcode,
	)
	if err != nil {
		er(err)
	}
	defer rows.Close()
	count := 0
	var multiple []string
	one := ""
	for rows.Next() {
		var pref string
		err := rows.Scan(&pref)
		if err != nil {
			er(err)
		}
		if count == 0 {
			one = pref
			count++
		} else if count == 1 {
			if one != pref {
				multiple = make([]string, 2)
				multiple[0] = one
				multiple[1] = pref
				one = pref
				count++
			}
		} else {
			if one != pref {
				multiple = append(multiple, pref)
				one = pref
				count++
			}
		}
	}

	if count == 0 {
		rows2, err := db.Query(
			`SELECT PREF FROM KENALL WHERE OLDZIP = ?`, pcode[:5],
		)
		if err != nil {
			er(err)
		}
		defer rows2.Close()
		for rows2.Next() {
			var pref string
			err := rows2.Scan(&pref)
			if err != nil {
				er(err)
			}
			if count == 0 {
				one = pref
				count++
			} else if count == 1 {
				if one != pref {
					multiple = make([]string, 2)
					multiple[0] = one
					multiple[1] = pref
					one = pref
					count++
				}
			} else {
				if one != pref {
					multiple = append(multiple, pref)
					one = pref
					count++
				}
			}
		}
	}

	if count == 0 {
		rows3, err := db.Query(
			`SELECT PREF FROM KENALL WHERE OLDZIP = ?`, pcode[:3]+"  ",
		)
		if err != nil {
			er(err)
		}
		defer rows3.Close()
		for rows3.Next() {
			var pref string
			err := rows3.Scan(&pref)
			if err != nil {
				er(err)
			}
			if count == 0 {
				one = pref
				count++
			} else if count == 1 {
				if one != pref {
					multiple = make([]string, 2)
					multiple[0] = one
					multiple[1] = pref
					one = pref
					count++
				}
			} else {
				if one != pref {
					multiple = append(multiple, pref)
					one = pref
					count++
				}
			}
		}
	}

	if count == 1 {
		c, ok := PrefCodeTable[one]
		if !ok {
			er("invalid pref " + one)
		}
		return c, true
	} else {
		town := GetTown(address, wbuff)
		trows, err := db.Query(
			`SELECT PREF FROM TOWNCODE WHERE TOWN = ?`, town,
		)
		if err != nil {
			er(err)
		}
		defer trows.Close()
		mtowns := make([]string, 0)
		tcount := 0
		for trows.Next() {
			var pref string
			err := trows.Scan(&pref)
			if err != nil {
				er(err)
			}
			mtowns = append(mtowns, pref)
			tcount++
		}
		if count == 0 {
			if tcount > 1 || tcount == 0 {
				return "", false
			}
			c, ok := PrefCodeTable[mtowns[0]]
			if !ok {
				er("invalid pref " + one)
			}
			return c, true
		}
		if tcount > 0 {
			for _, t := range mtowns {
				for _, pt := range multiple {
					if t == pt {
						c, ok := PrefCodeTable[t]
						if !ok {
							er("invalid pref " + one)
						}
						return c, true
					}
				}
			}
		}
		return "", false
	}

}

type Column8Type int

const (
	NotByosyo Column8Type = iota
	ByosyoWithNumber
	ByosyoSingle
	ByosyoNumberOnly
)

func Byosyo(text string, wbuff []rune) (Column8Type, string, int) {
	wbuff = wbuff[:0]
	fToken := false
	fNumber := false
	tCount := 0
	token := ""
	for _, c := range text {
		if fToken {
			if c == ' ' || c == '　' || c == '、' {
				token = string(wbuff)
				wbuff = wbuff[:0]
				fToken = false
				tCount++
			} else {
				wbuff = append(wbuff, c)
			}
		} else {
			if c == ' ' || c == '　' {
				if fNumber {
					break
				}
			} else {
				if tCount == 0 {
					wbuff = append(wbuff, c)
					if c == '0' || c == '1' || c == '2' || c == '3' || c == '4' || c == '5' || c == '6' || c == '7' || c == '8' || c == '9' {
						if !fNumber {
							fNumber = true
						}
					} else {
						fToken = true
					}
				} else {
					if c == '0' || c == '1' || c == '2' || c == '3' || c == '4' || c == '5' || c == '6' || c == '7' || c == '8' || c == '9' {
						wbuff = append(wbuff, c)
						if !fNumber {
							fNumber = true
						}
					} else {
						tCount++
						break
					}
				}
			}
		}
	}
	if tCount > 1 {
		return NotByosyo, token, 0
	} else if tCount == 0 && fToken {
		return ByosyoSingle, string(wbuff), 0
	} else if tCount == 1 {
		if fNumber {
			c, err := strconv.Atoi(string(wbuff))
			if err != nil {
				er(err.Error() + ": " + string(wbuff))
			}
			return ByosyoWithNumber, token, c
		} else {
			return ByosyoSingle, token, 0
		}
	} else {
		if fNumber {
			c, err := strconv.Atoi(string(wbuff))
			if err != nil {
				er(err.Error() + ": " + string(wbuff))
			}
			return ByosyoNumberOnly, "", c
		}
	}
	return NotByosyo, "", 0
}

func SplitPostal(text string, wbuff []rune) (string, string, bool) {
	wbuff = wbuff[:0]
	textrune := []rune(text)
	if textrune[0] != '〒' {
		return "", "", false
	}
	yMark := 0
	hMark := 0
	for _, c := range textrune {
		if c == '0' || c == '1' || c == '2' || c == '3' || c == '4' || c == '5' || c == '6' || c == '7' || c == '8' || c == '9' || c == '-' || c == '〒' {
			if c == '〒' {
				yMark++
			} else if c == '-' {
				hMark++
			}
			wbuff = append(wbuff, c)
		}
	}
	if yMark != 1 && hMark != 1 {
		return "", "", false
	}
	return string(wbuff), string(textrune[len(wbuff):]), true
}

func SplitGenzonOrKyushi(text string) (string, string, bool) {
	textrune := []rune(text)
	boundary := len(textrune) - 2
	if boundary <= 0 {
		return "", "", false
	}
	genOrKyu := string(textrune[boundary:])
	if genOrKyu != "現存" && genOrKyu != "休止" && genOrKyu != "辞退" && genOrKyu != "廃止" {
		return "", "", false
	}
	return string(textrune[:boundary]), genOrKyu, true
}

func SelectNumbers(text string, wbuff []rune) string {
	wbuff = wbuff[:0]
	for _, c := range text {
		if c == '0' || c == '1' || c == '2' || c == '3' || c == '4' || c == '5' || c == '6' || c == '7' || c == '8' || c == '9' {
			wbuff = append(wbuff, c)
		}
	}
	return string(wbuff)
}

func NormalizeTel(text string, wbuff []rune) string {
	wbuff = wbuff[:0]
	for _, c := range text {
		if c == '0' || c == '1' || c == '2' || c == '3' || c == '4' || c == '5' || c == '6' || c == '7' || c == '8' || c == '9' || c == '-' {
			wbuff = append(wbuff, c)
		} else if c == '(' || c == ')' || c == '（' || c == '）' {
			wbuff = append(wbuff, '-')
		}
	}
	return string(wbuff)
}

func VerifyPageHeader(text []rune) (string, string, bool) {
	maxlen := 11
	rval1 := make([]rune, 0, maxlen)
	rval2 := make([]rune, 0, 2)
	len1 := 0
	len2 := 0
	fStart := 0
	fEnd := 0
	fNen := 0
	fTuki := 0
	fHi := 0
	for _, c := range text {
		if fStart == 0 {
			if c == '［' {
				fStart++
			}
		} else {
			if c == '医' || c == '科' || c == '歯' || c == '薬' || c == '局' {
				len2++
				if len2 > 2 {
					return "", "", false
				}
				rval2 = append(rval2, c)
			} else if c == '0' || c == '1' || c == '2' || c == '3' || c == '4' || c == '5' || c == '6' || c == '7' || c == '8' || c == '9' {
				len1++
				if len1 > maxlen {
					return "", "", false
				}
				rval1 = append(rval1, c)
			} else if c == '平' || c == '成' {
				len1++
				if len1 > maxlen {
					return "", "", false
				}
				rval1 = append(rval1, c)
			} else if c == '年' {
				fNen++
				len1++
				if len1 > maxlen {
					return "", "", false
				}
				rval1 = append(rval1, c)
			} else if c == '月' {
				fTuki++
				len1++
				if len1 > maxlen {
					return "", "", false
				}
				rval1 = append(rval1, c)
			} else if c == '日' {
				fHi++
				len1++
				if len1 > maxlen {
					return "", "", false
				}
				rval1 = append(rval1, c)
			} else if c == '］' {
				fEnd++
			}
		}
	}
	if fStart == 1 && fEnd == 1 && fNen == 1 && fTuki == 1 && fHi == 1 && len2 == 2 {
		return string(rval1), string(rval2), true
	}
	return "", "", false
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

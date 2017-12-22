package cmd

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func init() {
	RootCmd.AddCommand(codeCmd)
}

var codeCmd = &cobra.Command{
	Use:   "code csvfile",
	Short: "get facilities list",
	Long:  "get facilities list",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		code(args)
	},
}

func code(args []string) {
	mkFacilitiesListFile(args[0], args[1])
}

type Data1 struct {
	No       string
	PrefCode string
	Code     string
	Name     string
	Postal   string
	Address  string
	Tel      string
	GenOrKyu string
}

func mkFacilitiesListFile(input string, outdir string) {
	var writer1 *csv.Writer
	var writer2 *csv.Writer
	var file1 *os.File
	var file2 *os.File
	data1line := make([]string, 11, 11)
	data2line := make([]string, 6, 6)

	ifile, err := os.Open(input)
	if err != nil {
		er(input + ": " + err.Error())
	}
	defer ifile.Close()

	reader := csv.NewReader(ifile)
	reader.Comma = '\t'

	db, err := sql.Open("sqlite3", "prefdb.db")
	if err != nil {
		er(input + ": " + err.Error())
	}
	defer db.Close()

	var data1 Data1

	wbuff := make([]rune, 0, 20)
	lastdate := ""
	lastkubun := ""
	inRecord := false
	var lastc8 Column8Type = ByosyoWithNumber
	lastbyosyo := ""
	//prefcode := ""
	lastno := 0
	lastpref := ""
	var byosyomap map[string]int
	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				// ここでも出力TODO
				for key, value := range byosyomap {
					data2line[0] = data1.No
					data2line[1] = lastdate
					data2line[2] = lastkubun
					data2line[3] = data1.PrefCode + data1.Code
					data2line[4] = key
					data2line[5] = strconv.Itoa(value)
					err = writer2.Write(data2line)
					if err != nil {
						er(err)
					}
					//fmt.Printf("%s %s %s %s %s %d\n", data1.No, lastdate, lastkubun, data1.PrefCode+data1.Code, key, value)
				}
				break
			} else {
				er(input + ": " + err.Error())
			}
		}
		if !inRecord {
			r0 := strings.TrimSpace(record[0])
			if len(r0) == 0 {
				continue
			}
			r0rune := []rune(r0)
			if r0rune[0] == '［' {
				d, k, ok := VerifyPageHeader(r0rune)
				if !ok {
					er(input + ": invalid page header")
				}
				if lastdate != d {
					lastdate = d
				}
				if lastkubun != k {
					lastkubun = k
				}
				continue
			} else if r0rune[0] == '0' || r0rune[0] == '1' || r0rune[0] == '2' || r0rune[0] == '3' || r0rune[0] == '4' || r0rune[0] == '5' || r0rune[0] == '6' || r0rune[0] == '7' || r0rune[0] == '8' || r0rune[0] == '9' {
				if lastdate == "" || lastkubun == "" {
					er(input + ": miss page header")
				}
				no, err := strconv.Atoi(record[0])
				if err != nil {
					er(input + ": " + err.Error() + ": " + record[0])
				}
				if lastno != no-1 {
					fmt.Println(input + ": not sequential number: " + record[0])
					//break
				}
				lastno = no
				data1.No = record[0]
				data1.Code = SelectNumbers(record[1], wbuff)
				data1.Name = record[2]
				data1.Address = record[3]
				data1.Tel = NormalizeTel(record[4], wbuff)
				data1.GenOrKyu = record[9]
				t, g, ok := SplitGenzonOrKyushi(record[9])
				if !ok {
					er(input + ": invalid value of Genzon or Kyushi column: " + record[2] + " " + record[3] + " " + record[9])
				}
				p, a, ok := SplitPostal(record[3], wbuff)
				if !ok {
					if record[3] == "恵庭市黄金中央２丁目１番地２" {
						p = "〒061-1449"
						a = record[3]
						data1.Address = p + a
					} else {
						er(input + ": invalid value of Postal code or Address: " + record[2] + " " + record[3])
					}
				}
				pc, ok := GetPrefCode(p, a, db, wbuff)
				if !ok {
					if record[2] == "みぞはた眼科" && p == "〒940-0" {
						p = "〒940-2108"
						pc = "15"
						data1.Address = p + a
					} else {
						er(input + ": Unknown City: " + record[2] + " " + record[3])
					}
				}
				if pc == "15" && data1.Code == "0214099" {
					data1.Name = "辻本皮膚科"
				}
				if lastpref == "" {
					// ここで都道府県ディレクトリを作る
					pcodedir := filepath.Join(outdir, pc)
					stat, err := os.Stat(pcodedir)
					if err != nil {
						if os.IsNotExist(err) {
							err := os.MkdirAll(pcodedir, 0755)
							if err != nil {
								er(err)
							}
						} else {
							er(err)
						}
					} else {
						if !stat.IsDir() {
							er(pcodedir + " file already exists")
						}
					}
					// ここで区分ディレクトリを作る
					var ekubun string
					switch lastkubun {
					case "医科":
						ekubun = "ika"
					case "歯科":
						ekubun = "sika"
					case "薬局":
						ekubun = "yaku"
					default:
						er("Unknown kubun: " + lastkubun)
					}
					kubundir := filepath.Join(outdir, pc, ekubun)
					stat, err = os.Stat(kubundir)
					if err != nil {
						if os.IsNotExist(err) {
							err := os.Mkdir(kubundir, 0755)
							if err != nil {
								er(err)
							}
						} else {
							er(err)
						}
					} else {
						if !stat.IsDir() {
							er(kubundir + " file already exists")
						}
					}
					// data1
					data1path := filepath.Join(outdir, pc, ekubun, "data1.txt")
					file1, err = os.OpenFile(data1path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
					if err != nil {
						er(err)
					}
					defer file1.Close()
					// data2
					data2path := filepath.Join(outdir, pc, ekubun, "data2.txt")
					file2, err = os.OpenFile(data2path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
					if err != nil {
						er(err)
					}
					defer file2.Close()
					writer1 = csv.NewWriter(file1)
					writer1.Comma = '\t'
					defer writer1.Flush()
					writer2 = csv.NewWriter(file2)
					writer2.Comma = '\t'
					defer writer2.Flush()
				}
				if lastpref != "" && lastpref != pc {
					fmt.Println(input + ": change pref code: " + lastpref + " to " + pc + ": " + record[2] + " " + record[3])
					if record[2] == "若葉歯科医院" && p == "〒891-3206" {
						pc = "43"
						p = "〒861-3206"
						data1.Address = p + a
						fmt.Println("fix it")
					}
				}
				lastpref = pc
				data1.PrefCode = pc
				//fmt.Printf("%s %s %s %s %s %s %s %s %s %s %s\n", record[0], lastdate, lastkubun, pc+data1.Code, record[2], p, a, record[4], record[9], t, g)
				data1line[0] = data1.No
				data1line[1] = lastdate
				data1line[2] = lastkubun
				data1line[3] = data1.PrefCode + data1.Code
				data1line[4] = data1.Name
				data1line[5] = p
				data1line[6] = a
				data1line[7] = data1.Tel
				data1line[8] = data1.GenOrKyu
				data1line[9] = t
				data1line[10] = g
				err = writer1.Write(data1line)
				if err != nil {
					er(err)
				}

				byosyomap = map[string]int{}
				c8, byosyo, bnum := Byosyo(record[8], wbuff)
				if c8 == ByosyoWithNumber {
					byosyomap[byosyo] = bnum
					//fmt.Printf("%s %s %s %s %s %d\n", record[0], lastdate, lastkubun, pc+data1.Code, byosyo, bnum)
					lastbyosyo = ""
				} else if c8 == ByosyoSingle {
					lastbyosyo = byosyo
				} else {
					lastbyosyo = ""
				}
				lastc8 = c8
				inRecord = true
			}
		} else {
			// 以下の条件のどれかが現れたとき、一つの施設情報が終端していると判断する
			// ----- の行
			// 空行
			// 1つめのカラムが、数字の行
			// EOF
			r0 := strings.TrimSpace(record[0])
			if len(r0) > 0 {
				r0rune := []rune(r0)
				if r0rune[0] == '-' {
					// ここで出力 TODO
					for key, value := range byosyomap {
						data2line[0] = data1.No
						data2line[1] = lastdate
						data2line[2] = lastkubun
						data2line[3] = data1.PrefCode + data1.Code
						data2line[4] = key
						data2line[5] = strconv.Itoa(value)
						err = writer2.Write(data2line)
						if err != nil {
							er(err)
						}
						// fmt.Printf("%s %s %s %s %s %d\n", data1.No, lastdate, lastkubun, data1.PrefCode+data1.Code, key, value)
					}
					inRecord = false
				} else if r0rune[0] == '0' || r0rune[0] == '1' || r0rune[0] == '2' || r0rune[0] == '3' || r0rune[0] == '4' || r0rune[0] == '5' || r0rune[0] == '6' || r0rune[0] == '7' || r0rune[0] == '8' || r0rune[0] == '9' {
					// 万が一いきなり項番の場合、出力 TODO
					for key, value := range byosyomap {
						data2line[0] = data1.No
						data2line[1] = lastdate
						data2line[2] = lastkubun
						data2line[3] = data1.PrefCode + data1.Code
						data2line[4] = key
						data2line[5] = strconv.Itoa(value)
						err = writer2.Write(data2line)
						if err != nil {
							er(err)
						}
						// fmt.Printf("%s %s %s %s %s %d\n", data1.No, lastdate, lastkubun, data1.PrefCode+data1.Code, key, value)
					}
					no, err := strconv.Atoi(record[0])
					if err != nil {
						er(input + ": " + err.Error() + ": " + record[0])
					}
					if lastno != no-1 {
						fmt.Println(input + ": not sequential number: " + record[0])
						//break
					}
					lastno = no
					data1.No = record[0]
					data1.Code = SelectNumbers(record[1], wbuff)
					data1.Name = record[2]
					data1.Address = record[3]
					data1.Tel = NormalizeTel(record[4], wbuff)
					data1.GenOrKyu = record[9]
					t, g, ok := SplitGenzonOrKyushi(record[9])
					if !ok {
						er(input + ": invalid value of Genzon or Kyushi column: " + record[2] + " " + record[3] + " " + record[9])
					}
					p, a, ok := SplitPostal(record[3], wbuff)
					if !ok {
						if record[3] == "恵庭市黄金中央２丁目１番地２" {
							p = "〒061-1449"
							a = record[3]
							data1.Address = p + a
						} else {
							er(input + ": invalid value of Postal code or Address: " + record[3])
						}
					}
					pc, ok := GetPrefCode(p, a, db, wbuff)
					if !ok {
						if record[2] == "みぞはた眼科" && p == "〒940-0" {
							p = "〒940-2108"
							pc = "15"
							data1.Address = p + a
						} else {
							er(input + ": Unknown City: " + record[2] + " " + record[3])
						}
					}
					if lastpref != pc {
						fmt.Println(input + ": change pref code: " + lastpref + " to " + pc + ": " + record[2] + " " + record[3])
						if record[2] == "若葉歯科医院" && p == "〒891-3206" {
							pc = "43"
							p = "〒861-3206"
							data1.Address = p + a
							fmt.Println("fix it")
						}
					}
					lastpref = pc
					data1.PrefCode = pc
					data1line[0] = data1.No
					data1line[1] = lastdate
					data1line[2] = lastkubun
					data1line[3] = data1.PrefCode + data1.Code
					data1line[4] = data1.Name
					data1line[5] = p
					data1line[6] = a
					data1line[7] = data1.Tel
					data1line[8] = data1.GenOrKyu
					data1line[9] = t
					data1line[10] = g
					err = writer1.Write(data1line)
					if err != nil {
						er(err)
					}

					byosyomap = map[string]int{}
					//fmt.Printf("%s %s %s %s %s %s %s %s %s %s %s\n", record[0], lastdate, lastkubun, pc+data1.Code, record[2], p, a, record[4], record[9], t, g)

					c8, byosyo, bnum := Byosyo(record[8], wbuff)
					if c8 == ByosyoWithNumber {
						byosyomap[byosyo] = bnum
						//fmt.Printf("%s %s %s %s %s %d\n", record[0], lastdate, lastkubun, pc+data1.Code, byosyo, bnum)
						lastbyosyo = ""
					} else if c8 == ByosyoSingle {
						lastbyosyo = byosyo
					} else {
						lastbyosyo = ""
					}
					lastc8 = c8
				}
			} else {
				// 空行かどうかの確認
				var i int
				var r string
				nstrcount := 0
				for i, r = range record {
					if len(strings.TrimSpace(r)) > 0 {
						// ここから下病床数取得コード
						if i == 8 {
							c8, byosyo, bnum := Byosyo(r, wbuff)
							if c8 == ByosyoWithNumber {
								c, ok := byosyomap[byosyo]
								if !ok {
									byosyomap[byosyo] = bnum
								} else {
									byosyomap[byosyo] = c + bnum
								}
								//fmt.Printf("%s %s %s %s %s %d\n", data1.No, lastdate, lastkubun, data1.PrefCode+data1.Code, byosyo, bnum)
							} else if c8 == ByosyoNumberOnly && lastc8 == ByosyoSingle {
								c, ok := byosyomap[lastbyosyo]
								if !ok {
									byosyomap[lastbyosyo] = bnum
								} else {
									byosyomap[lastbyosyo] = c + bnum
								}
								//fmt.Printf("%s %s %s %s %s %d\n", data1.No, lastdate, lastkubun, data1.PrefCode+data1.Code, lastbyosyo, bnum)
								lastbyosyo = ""
							} else if c8 == ByosyoSingle {
								lastbyosyo = byosyo
							} else {
								lastbyosyo = ""
							}
							lastc8 = c8
						}
					} else {
						nstrcount++
					}
				}
				if nstrcount == len(record) {
					// 空行
					// ここで出力 TODO
					for key, value := range byosyomap {
						data2line[0] = data1.No
						data2line[1] = lastdate
						data2line[2] = lastkubun
						data2line[3] = data1.PrefCode + data1.Code
						data2line[4] = key
						data2line[5] = strconv.Itoa(value)
						err = writer2.Write(data2line)
						if err != nil {
							er(err)
						}
						//fmt.Printf("%s %s %s %s %s %d\n", data1.No, lastdate, lastkubun, data1.PrefCode+data1.Code, key, value)
					}
					inRecord = false
				}
			}
		}
	}
}

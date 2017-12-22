package cmd

import (
	"encoding/csv"
	"github.com/spf13/cobra"
	"github.com/tealeg/xlsx"
	//"github.com/360EntSecGroup-Skylar/excelize"
	//"golang.org/x/text/encoding/japanese"
	//"golang.org/x/text/transform"
	"fmt"
	//"io"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	RootCmd.AddCommand(dataCmd)
}

var dataCmd = &cobra.Command{
	Use:   "data file.xlsx outputdir",
	Short: "get csv data",
	Long:  "get csv data",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		data(args)
	},
}

func data(args []string) {
	mkDataFile(args[0], args[1])
}

func mkDataFile(input string, outdir string) {
	xfile, err := xlsx.OpenFile(input)
	if err != nil {
		er(err)
	}

	// xlsx, err := excelize.OpenFile(input)
	// if err != nil {
	// 	er(err)
	// }

	// 都道府県コード
	// 医科/歯科/薬局
	// data1.txt
	// data2.txt
	// data3.txt

	var kubun string
	var pcode string
	var code string
	var writer1 *csv.Writer
	var writer2 *csv.Writer
	var writer3 *csv.Writer
	var file1 *os.File
	var file2 *os.File
	var file3 *os.File
	var opened = false

	count := 0
	lastLine1 := make([]string, 9, 9)
	data2 := make([]string, 6, 6)
	data3 := make([]string, 6, 6)
	maxlen := 20
	rbuf := make([]rune, 0, maxlen)

	jiten := ""
	fDate := false
	fStart := false
	for sheetidx, sheet := range xfile.Sheets {
		rowslen := len(sheet.Rows)
		for rowidx, row := range sheet.Rows {
			cellslen := len(row.Cells)
			if cellslen == 0 {
				continue
			}
			// 項番
			columnA := strings.TrimSpace(row.Cells[0].String())
			if columnA == "7000" {
				fmt.Println(jiten)
			}
			if sheetidx == 0 {
				if !fDate {
					var ok bool
					jiten, ok = VerifyDate(columnA)
					if !ok {
						er("maybe invalid xlsx: no date")
					}
					fDate = true
					fmt.Printf("%d: %s %s\n", rowslen, columnA, input)
					continue
				}
				if !fStart {
					if columnA == "" {
						continue
					}
					if columnA == "項番" {
						fStart = true
						continue
					} else {
						er("maybe invalid xlsx: no header")
					}
				}
			} else {
				if rowidx == 0 {
					jitennext, ok := VerifyDate(columnA)
					if ok {
						jiten = jitennext
						fStart = false
						continue
					}
				}
				if !fStart {
					if columnA == "" {
						continue
					}
					if columnA == "項番" {
						fStart = true
						continue
					} else {
						er("maybe invalid xlsx: no header")
					}
				}
			}

			count++
			if count == 1 {
				// 都道府県コード2桁化
				// pcode = row[1]
				pcode = row.Cells[1].String()
				if len(pcode) == 1 {
					pcode = "0" + pcode
				}
				// ここで都道府県ディレクトリを作る
				pcodedir := filepath.Join(outdir, pcode)
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
				// kubun = row[3]
				kubun = row.Cells[3].String()
				var ekubun string
				switch kubun {
				case "医科":
					ekubun = "ika"
				case "歯科":
					ekubun = "sika"
				case "薬局":
					ekubun = "yaku"
				default:
					er("Unknown kubun: " + kubun)
				}
				kubundir := filepath.Join(outdir, pcode, ekubun)
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
				data1path := filepath.Join(outdir, pcode, ekubun, "data1.txt")
				file1, err = os.OpenFile(data1path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
				if err != nil {
					er(err)
				}
				// defer file1.Close()
				// data2
				data2path := filepath.Join(outdir, pcode, ekubun, "data2.txt")
				file2, err = os.OpenFile(data2path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
				if err != nil {
					er(err)
				}
				// defer file2.Close()
				// data3
				data3path := filepath.Join(outdir, pcode, ekubun, "data3.txt")
				file3, err = os.OpenFile(data3path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
				if err != nil {
					er(err)
				}
				// defer file3.Close()
				writer1 = csv.NewWriter(file1)
				writer1.Comma = '\t'
				// defer writer1.Flush()
				writer2 = csv.NewWriter(file2)
				writer2.Comma = '\t'
				// defer writer2.Flush()
				writer3 = csv.NewWriter(file3)
				writer3.Comma = '\t'
				// defer writer3.Flush()

				opened = true

				// code = pcode + row[4]
				code = pcode + row.Cells[4].String()

				// 項番
				lastLine1[0] = columnA
				// 時点
				lastLine1[1] = jiten
				// 区分
				lastLine1[2] = kubun
				// 都道府県コード+医療機関番号
				lastLine1[3] = code
				// 医療機関名称
				// lastLine1[4] = row[7]
				lastLine1[4] = row.Cells[7].String()
				// 郵便番号
				// lastLine1[5] = "〒" + row[8]
				lastLine1[5] = "〒" + row.Cells[8].String()
				// 住所
				// lastLine1[6] = row[9]
				lastLine1[6] = row.Cells[9].String()
				// 電話番号
				// lastLine1[7] = row[10]
				lastLine1[7] = row.Cells[10].String()
				// FAX番号
				// lastLine1[8] = row[11]
				lastLine1[8] = ""
				if cellslen > 11 {
					lastLine1[8] = row.Cells[11].String()
				}

				// 書き出し
				err = writer1.Write(lastLine1)
				if err != nil {
					er(err)
				}
				//fmt.Println(lastLine1)

				// 項番
				data2[0] = columnA
				// 時点
				data2[1] = jiten
				// 区分
				data2[2] = kubun
				// 都道府県コード+医療機関番号
				data2[3] = code
				// if row[12] != "" {
				// 	SplitByByosyo1(writer2, data2, row[12])
				// }
				if cellslen > 12 {
					byosyo := row.Cells[12].String()
					if byosyo != "" {
						SplitByByosyo1(writer2, data2, byosyo)
					}
				}

				if cellslen > 16 {
					// 項番
					data3[0] = columnA
					// 時点
					data3[1] = jiten
					// 区分
					data3[2] = kubun
					// 都道府県コード+医療機関番号
					data3[3] = code
					// 受理記号
					// data3[4] = "（" + row[14] + "）" + row.Cells[15].String()
					data3[4] = "（" + row.Cells[14].String() + "）" + row.Cells[15].String()
					// 算定開始日
					// 初期化trick
					rbuf = rbuf[:0]
					// data3[5] = RemoveSpace(rbuf, maxlen, row[16])
					data3[5] = RemoveSpace(rbuf, maxlen, row.Cells[16].String())
					err = writer3.Write(data3)
					if err != nil {
						er(err)
					}
					//fmt.Println(data3)
				}
				continue
			}

			var ekubun string
			changed := false
			// 都道府県コード2桁化
			// npcode := row[1]
			npcode := row.Cells[1].String()
			if len(npcode) == 1 {
				npcode = "0" + npcode
			}
			if npcode != pcode {
				pcode = npcode

				// ここで都道府県ディレクトリを作る
				pcodedir := filepath.Join(outdir, pcode)
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
				changed = true
			}

			kubun = row.Cells[3].String()
			if lastLine1[2] != kubun {
				switch kubun {
				case "医科":
					ekubun = "ika"
				case "歯科":
					ekubun = "sika"
				case "薬局":
					ekubun = "yaku"
				default:
					er("Unknown kubun: " + kubun)
				}
				kubundir := filepath.Join(outdir, pcode, ekubun)
				stat, err := os.Stat(kubundir)
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
				changed = true
			}

			if changed {
				fmt.Println("Flush")
				writer1.Flush()
				writer2.Flush()
				writer3.Flush()
				_ = file1.Close()
				_ = file2.Close()
				_ = file3.Close()
				// data1
				data1path := filepath.Join(outdir, pcode, ekubun, "data1.txt")
				file1, err = os.OpenFile(data1path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
				if err != nil {
					er(err)
				}
				// data2
				data2path := filepath.Join(outdir, pcode, ekubun, "data2.txt")
				file2, err = os.OpenFile(data2path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
				if err != nil {
					er(err)
				}
				// data3
				data3path := filepath.Join(outdir, pcode, ekubun, "data3.txt")
				file3, err = os.OpenFile(data3path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
				if err != nil {
					er(err)
				}
				writer1 = csv.NewWriter(file1)
				writer1.Comma = '\t'
				writer2 = csv.NewWriter(file2)
				writer2.Comma = '\t'
				writer3 = csv.NewWriter(file3)
				writer3.Comma = '\t'
			}

			// code = pcode + row[4]
			code = pcode + row.Cells[4].String()

			row7 := row.Cells[7].String()
			row8 := "〒" + row.Cells[8].String()
			row9 := row.Cells[9].String()
			row10 := row.Cells[10].String()
			row11 := ""
			if cellslen > 11 {
				row11 = row.Cells[11].String()
			}
			// if lastLine1[0] != columnA ||
			// 	lastLine1[3] != code ||
			// 	lastLine1[4] != row[7] ||
			// 	lastLine1[5] != "〒"+row[8] ||
			// 	lastLine1[6] != row[9] ||
			// 	lastLine1[7] != row[10] ||
			// 	lastLine1[8] != row[11] {
			if lastLine1[0] != columnA ||
				lastLine1[2] != kubun ||
				lastLine1[3] != code ||
				lastLine1[4] != row7 ||
				lastLine1[5] != row8 ||
				lastLine1[6] != row9 ||
				lastLine1[7] != row10 ||
				lastLine1[8] != row11 {

				// if lastLine1[0] == columnA {
				// 	er("which the different value???")
				// }

				// if lastLine1[2] != kubun {
				// 	er("Oops! kubun is changed")
				// }

				// 項番
				lastLine1[0] = columnA
				// 時点
				lastLine1[1] = jiten
				// 区分
				lastLine1[2] = kubun
				// 都道府県コード+医療機関番号
				lastLine1[3] = code
				// 医療機関名称
				// lastLine1[4] = row[7]
				lastLine1[4] = row7
				// 郵便番号
				// lastLine1[5] = "〒" + row[8]
				lastLine1[5] = row8
				// 住所
				// lastLine1[6] = row[9]
				lastLine1[6] = row9
				// 電話番号
				// lastLine1[7] = row[10]
				lastLine1[7] = row10
				// FAX番号
				// lastLine1[8] = row[11]
				lastLine1[8] = row11
				err = writer1.Write(lastLine1)
				if err != nil {
					er(err)
				}
				//fmt.Println(lastLine1)

				// 項番
				data2[0] = columnA
				// 時点
				data2[1] = jiten
				// 区分
				data2[2] = kubun
				// 都道府県コード+医療機関番号
				data2[3] = code
				// if row[12] != "" {
				// 	SplitByByosyo1(writer2, data2, row[12])
				// }
				if cellslen > 12 {
					byosyo := row.Cells[12].String()
					if byosyo != "" {
						SplitByByosyo1(writer2, data2, byosyo)
					}
				}

				if cellslen > 16 {
					// 項番
					data3[0] = columnA
					// 時点
					data3[1] = jiten
					// 区分
					data3[2] = kubun
					// 都道府県コード+医療機関番号
					data3[3] = code
					// 受理記号
					// data3[4] = "（" + row[14] + "）" + row.Cells[15].String()
					data3[4] = "（" + row.Cells[14].String() + "）" + row.Cells[15].String()
					// 算定開始日
					// 初期化trick
					rbuf = rbuf[:0]
					// data3[5] = RemoveSpace(rbuf, maxlen, row[16])
					data3[5] = RemoveSpace(rbuf, maxlen, row.Cells[16].String())
					err = writer3.Write(data3)
					if err != nil {
						er(err)
					}
					//fmt.Println(data3)
				}
			} else {
				if cellslen > 16 {
					// 受理記号
					// data3[4] = "(" + row[14] + ")"
					data3[4] = "（" + row.Cells[14].String() + "）" + row.Cells[15].String()
					// 算定開始日
					// 初期化trick
					rbuf = rbuf[:0]
					// data3[5] = RemoveSpace(rbuf, maxlen, row[16])
					data3[5] = RemoveSpace(rbuf, maxlen, row.Cells[16].String())
					err := writer3.Write(data3)
					if err != nil {
						er(err)
					}
					//fmt.Println(data3)
				}
			}
		}
	}

	if opened {
		fmt.Println("Flush")
		writer1.Flush()
		writer2.Flush()
		writer3.Flush()
		_ = file1.Close()
		_ = file2.Close()
		_ = file3.Close()
	}
}

/*
 * 「一般　　12／介護　　5／療養　　2」 のような入力
 * この場合3行つくる
 * 病床数が未記入の場合出力しない
 */
func SplitByByosyo1(writer2 *csv.Writer, data2 []string, text string) {
	kubuns := strings.Split(text, "／")
	for _, kubun := range kubuns {
		kubun = strings.TrimSpace(kubun)
		if kubun != "" {
			byosyo := strings.Split(kubun, "　　")
			if len(byosyo) == 2 {
				byosyokubun := strings.TrimSpace(byosyo[0])
				byosyosu := strings.TrimSpace(byosyo[1])
				if byosyokubun != "" && byosyosu != "" {
					data2[4] = byosyokubun
					data2[5] = byosyosu
					err := writer2.Write(data2)
					if err != nil {
						er(err)
					}
					//fmt.Println(data2)
				}
			}
		}
	}
}

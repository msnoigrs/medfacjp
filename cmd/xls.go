package cmd

import (
	//"encoding/csv"
	"github.com/spf13/cobra"
	//"github.com/tealeg/xlsx"
	//"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/extrame/xls"
	//"golang.org/x/text/encoding/japanese"
	//"golang.org/x/text/transform"
	"fmt"
	//"io"
	//"os"
	"strings"
)

func init() {
	RootCmd.AddCommand(dataFromXlsCmd)
}

var dataFromXlsCmd = &cobra.Command{
	Use:   "xls file.xls",
	Short: "get csv data",
	Long:  "get csv data",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dataFromXls(args)
	},
}

func dataFromXls(args []string) {
	mkDataFileFromXls(args[0], "")
}

func mkDataFileFromXls(input string, outdir string) {
	xlsfile, err := xls.Open(input, "CP932")
	if err != nil {
		er(err)
	}
	fDate := false
	jiten := ""
	fStart := false
	count := 0
	lastLine1 := make([]string, 9, 9)
	data2 := make([]string, 6, 6)
	data3 := make([]string, 6, 6)
	maxlen := 20
	rbuf := make([]rune, 0, maxlen)

	sheet := xlsfile.GetSheet(0)
	if sheet == nil {
		er("sheet not found")
	}
	fmt.Println(sheet.MaxRow)
	for i := 0; i <= (int(sheet.MaxRow)); i++ {
		row := sheet.Row(i)
		// 項番
		columnA := strings.TrimSpace(row.Col(0))
		if !fDate {
			var ok bool
			jiten, ok = VerifyDate(columnA)
			if !ok {
				er("maybe invalid xls: no date")
			}
			fDate = true
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
				er("maybe invalid xls: no header")
			}
		}

		// 都道府県コード2桁化
		pcode := row.Col(1)
		if len(pcode) == 1 {
			pcode = "0" + pcode
		}
		code := pcode + row.Col(4)

		count++
		if count == 1 {
			// 項番
			lastLine1[0] = columnA
			// 時点
			lastLine1[1] = jiten
			// 区分
			lastLine1[2] = row.Col(3)
			// 都道府県コード+医療機関番号
			lastLine1[3] = code
			// 医療機関名称
			lastLine1[4] = row.Col(7)
			// 郵便番号
			lastLine1[5] = "〒" + row.Col(8)
			// 住所
			lastLine1[6] = row.Col(9)
			// 電話番号
			lastLine1[7] = row.Col(10)
			// FAX番号
			lastLine1[8] = row.Col(11)
			// 書き出し
			fmt.Println(lastLine1)

			// 項番
			data2[0] = columnA
			// 時点
			data2[1] = jiten
			// 区分
			data2[2] = row.Col(3)
			// 都道府県コード+医療機関番号
			data2[3] = code
			if row.Col(12) != "" {
				SplitByByosyo2(data2, row.Col(12))
			}

			// 項番
			data3[0] = columnA
			// 時点
			data3[1] = jiten
			// 区分
			data3[2] = row.Col(3)
			// 都道府県コード+医療機関番号
			data3[3] = code
			// 受理記号
			data3[4] = "(" + row.Col(14) + ")"
			// 算定開始日
			data3[5] = RemoveSpace(rbuf, maxlen, row.Col(16))
			fmt.Println(data3)
			continue
		}

		if lastLine1[0] != columnA ||
			lastLine1[2] != row.Col(3) ||
			lastLine1[3] != code ||
			lastLine1[4] != row.Col(7) ||
			lastLine1[5] != "〒"+row.Col(8) ||
			lastLine1[6] != row.Col(9) ||
			lastLine1[7] != row.Col(10) ||
			lastLine1[8] != row.Col(11) {

			if lastLine1[0] == columnA {
				er("which the different value???")
			}

			// 書き出し
			// 初期化trick
			//lastLine1 = lastLine1[:0]

			// 項番
			lastLine1[0] = columnA
			// 時点
			lastLine1[1] = jiten
			// 区分
			lastLine1[2] = row.Col(3)
			// 都道府県コード+医療機関番号
			lastLine1[3] = code
			// 医療機関名称
			lastLine1[4] = row.Col(7)
			// 郵便番号
			lastLine1[5] = "〒" + row.Col(8)
			// 住所
			lastLine1[6] = row.Col(9)
			// 電話番号
			lastLine1[7] = row.Col(10)
			// FAX番号
			lastLine1[8] = row.Col(11)
			fmt.Println(lastLine1)

			// 項番
			data2[0] = columnA
			// 時点
			data2[1] = jiten
			// 区分
			data2[2] = row.Col(3)
			// 都道府県コード+医療機関番号
			data2[3] = code
			if row.Col(12) != "" {
				SplitByByosyo2(data2, row.Col(12))
			}

			// 項番
			data3[0] = columnA
			// 時点
			data3[1] = jiten
			// 区分
			data3[2] = row.Col(3)
			// 都道府県コード+医療機関番号
			data3[3] = code
			// 受理記号
			data3[4] = "(" + row.Col(14) + ")"
			// 算定開始日
			// 初期化trick
			rbuf = rbuf[:0]
			data3[5] = RemoveSpace(rbuf, maxlen, row.Col(16))
			fmt.Println(data3)
		} else {
			// 受理記号
			data3[4] = "(" + row.Col(14) + ")"
			// 算定開始日
			// 初期化trick
			rbuf = rbuf[:0]
			data3[5] = RemoveSpace(rbuf, maxlen, row.Col(16))
			fmt.Println(data3)
		}
	}
}

/*
 * 「一般　　12／介護　　5／療養　　2」 のような入力
 * この場合3行つくる
 * 病床数が未記入の場合出力しない
 */
func SplitByByosyo2(data2 []string, text string) {
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
					fmt.Println(data2)
				}
			}
		}
	}
}

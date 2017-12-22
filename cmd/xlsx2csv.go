package cmd

import (
	"encoding/csv"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tealeg/xlsx"
	"os"
	"path/filepath"
)

func init() {
	RootCmd.AddCommand(xlsx2csvCmd)
}

var xlsx2csvCmd = &cobra.Command{
	Use:   "xlsx2csv xlsxfile",
	Short: "convert a xlsx file to a csv file",
	Long:  "convert a xlsx file to a csv file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		xlsx2csv(args)
	},
}

func xlsx2csv(args []string) {
	absp, err := filepath.Abs(args[0])
	if err != nil {
		er(err)
	}
	dir, filename := filepath.Split(absp)
	basename := filepath.Base(filename)
	ext := filepath.Ext(filename)
	if ext != ".xlsx" {
		er("Please specify a xlsx file")
	}
	basename = basename[:len(basename)-len(ext)]
	csvbase := filepath.Join(dir, basename)
	convertXlsxToCsv(args[0], csvbase)
}

func convertXlsxToCsv(input string, output string) {
	xfile, err := xlsx.OpenFile(input)
	if err != nil {
		er(err)
	}
	var ofile *os.File
	var writer *csv.Writer
	cellslen := -1
	var line []string
	for sheetidx, sheet := range xfile.Sheets {
		if sheetidx > 0 {
			writer.Flush()
			ofile.Close()
		}
		csvfile := fmt.Sprintf("%s_%d.txt", output, sheetidx)
		_, err := os.Stat(csvfile)
		if !os.IsNotExist(err) {
			er(csvfile + " is already exist")
		}
		ofile, err = os.OpenFile(csvfile, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			er(err)
		}
		writer = csv.NewWriter(ofile)
		writer.Comma = '\t'
		for _, row := range sheet.Rows {
			if cellslen == -1 {
				cellslen = len(row.Cells)
				line = make([]string, 0, cellslen)
			} else {
				line = line[:0]
			}
			for _, cell := range row.Cells {
				switch cell.Type() {
				case xlsx.CellTypeString:
					line = append(line, cell.String())
				case xlsx.CellTypeStringFormula:
					line = append(line, cell.String())
				case xlsx.CellTypeNumeric:
					line = append(line, cell.String())
				case xlsx.CellTypeBool:
					line = append(line, cell.String())
				case xlsx.CellTypeDate:
					t, err := cell.GetTime(false)
					if err != nil {
						er(err)
					}
					line = append(line, t.Format("2006-01-02T15:04:05"))
				default:
					er(fmt.Sprintf("a Unknonw type of cell: %d", cell.Type()))
				}
			}
			err = writer.Write(line)
			if err != nil {
				er(err)
			}
		}
	}
	writer.Flush()
	ofile.Close()
}

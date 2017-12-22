package cmd

import (
	//"fmt"
	"github.com/spf13/cobra"
	"github.com/tealeg/xlsx"
	//"os"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"strings"
)

func init() {
	RootCmd.AddCommand(soumuCmd)
}

var soumuCmd = &cobra.Command{
	Use:   "soumu xlsxfile",
	Short: "convert a xlsx file to a go source",
	Long:  "convert a xlsx file to a go source",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		soumu(args)
	},
}

// http://www.soumu.go.jp/denshijiti/code.html

func soumu(args []string) {
	// absp, err := filepath.Abs(args[0])
	// if err != nil {
	// 	er(err)
	// }
	// dir, filename := filepath.Split(absp)
	// basename := filepath.Base(filename)
	// ext := filepath.Ext(filename)
	// if ext != ".xlsx" {
	// 	er("Please specify a xlsx file")
	// }
	// basename = basename[:len(basename)-len(ext)]
	// csvbase := filepath.Join(dir, basename)
	soumuToGo(args[0], "prefdb.db")
}

func soumuToGo(input string, output string) {
	db, err := sql.Open("sqlite3", output)
	if err != nil {
		er(err)
	}
	defer db.Close()

	_, err = db.Exec(
		`DROP TABLE IF EXISTS "TOWNCODE"`,
	)
	if err != nil {
		er(err)
	}
	_, err = db.Exec(
		`CREATE TABLE "TOWNCODE" ("CODE" TEXT, "PREF" TEXT, "PREFRUBY" TEXT, "TOWN" TEXT, "TOWNRUBY" TEXT)`,
	)
	if err != nil {
		er(err)
	}

	xfile, err := xlsx.OpenFile(input)
	if err != nil {
		er(err)
	}

	tx, err := db.Begin()
	if err != nil {
		er(err)
	}
	stmt, err := tx.Prepare(`INSERT INTO TOWNCODE (CODE, PREF, PREFRUBY, TOWN, TOWNRUBY) VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		er(err)
	}
	defer stmt.Close()

	sheet := xfile.Sheets[0]

	for rowidx, row := range sheet.Rows {
		if rowidx == 0 {
			continue
		}
		c2 := strings.TrimSpace(row.Cells[2].String())
		if c2 == "" {
			continue
		}
		_, err := stmt.Exec(row.Cells[0].String(), row.Cells[1].String(), row.Cells[3].String(), c2, row.Cells[4].String())
		if err != nil {
			tx.Rollback()
			er(err)
		}

	}
	tx.Commit()
}

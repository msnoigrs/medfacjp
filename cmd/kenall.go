package cmd

import (
	//"fmt"
	"github.com/spf13/cobra"
	"lufia.org/pkg/japanese/zipcode"
	"os"
	//"strings"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	RootCmd.AddCommand(kenallCmd)
}

var kenallCmd = &cobra.Command{
	Use:   "kenall x-ken-all-utf8.csv",
	Short: "import postal map from x-ken-all-utf8.csv to sqldb3",
	Long:  "import postal map from x-ken-all-utf8.csv to sqldb3",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		kenall(args)
	},
}

// http://www.soumu.go.jp/denshijiti/code.html

func kenall(args []string) {
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
	kenallTo(args[0], "prefdb.db")
}

func kenallTo(input string, output string) {
	// ofile, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE, 0666)
	// if err != nil {
	// 	er(err)
	// }
	// defer ofile.Close()

	// _, err = ofile.WriteString("package cmd\n\nvar PrefTable = map[string]string{\n")
	// if err != nil {
	// 	er(err)
	// }

	db, err := sql.Open("sqlite3", output)
	if err != nil {
		er(err)
	}
	defer db.Close()

	_, err = db.Exec(
		`DROP TABLE IF EXISTS "KENALL"`,
	)
	if err != nil {
		er(err)
	}
	// _, err = db.Exec(
	// 	`CREATE TABLE IF NOT EXISTS "KENALL" ("CODE" TEXT, "OLDZIP" TEXT, "ZIP" TEXT, "PREF" TEXT, "PREFRUBY" TEXT, "REGION" TEXT, "REGIONRUBY" TEXT, "TOWN" TEXT, "TOWNRUBY" TEXT, "ISPARTIALTOWN" INTEGER, "ISLARGETOWN" INTEGER, "ISBLOCKEDSCHEME" INTEGER, "ISOVERLAPPEDZIP" INTEGER, "STATUS" INTEGER, "REASON" INTEGER)`,
	// )
	_, err = db.Exec(
		`CREATE TABLE "KENALL" ("CODE" TEXT, "OLDZIP" TEXT, "ZIP" TEXT, "PREF" TEXT, "PREFRUBY" TEXT, "REGION" TEXT, "REGIONRUBY" TEXT, "TOWN" TEXT, "TOWNRUBY" TEXT, "ISPARTIALTOWN" INTEGER, "ISLARGETOWN" INTEGER, "ISBLOCKEDSCHEME" INTEGER, "ISOVERLAPPEDZIP" INTEGER, "STATUS" INTEGER, "REASON" INTEGER)`,
	)
	if err != nil {
		er(err)
	}

	ifile, err := os.Open(input)
	if err != nil {
		er(err)
	}

	tx, err := db.Begin()
	if err != nil {
		er(err)
	}
	stmt, err := tx.Prepare(`INSERT INTO KENALL (CODE, OLDZIP, ZIP, PREF, PREFRUBY, REGION, REGIONRUBY, TOWN, TOWNRUBY, ISPARTIALTOWN, ISLARGETOWN, ISBLOCKEDSCHEME, ISOVERLAPPEDZIP, STATUS, REASON) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		er(err)
	}
	defer stmt.Close()

	var parser zipcode.Parser
	c := parser.Parse(ifile)
	for e := range c {
		//fmt.Println(e)
		isPartialTown := 0
		isLargeTown := 0
		isBlockedScheme := 0
		isOverlappedZip := 0
		if e.IsPartialTown {
			isPartialTown = 1
		}
		if e.IsLargeTown {
			isLargeTown = 1
		}
		if e.IsBlockedScheme {
			isBlockedScheme = 1
		}
		if e.IsOverlappedZip {
			isOverlappedZip = 1
		}
		_, err := stmt.Exec(e.Code, e.OldZip, e.Zip, e.Pref.Text, e.Pref.Ruby, e.Region.Text, e.Region.Ruby, e.Town.Text, e.Town.Ruby, isPartialTown, isLargeTown, isBlockedScheme, isOverlappedZip, e.Status, e.Reason)
		if err != nil {
			tx.Rollback()
			er(err)
		}
	}
	if parser.Error != nil {
		tx.Rollback()
		er(parser.Error)
	}

	tx.Commit()
}

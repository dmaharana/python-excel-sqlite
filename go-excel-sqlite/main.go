package main

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tealeg/xlsx"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
)

type User struct {
	ID    int    `bun:"id,pk"`
	Eid   string `bun:"eid"`
	Fname string `bun:"fname"`
	Lname string `bun:"lname"`
	Addr  string `bun:"addr"`
}

func readExcelToSQLite(sheetName string, excelFile string, dbFile string) error {
	ctx := context.Background()

	sqldb, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return err
	}
	defer sqldb.Close()

	db := bun.NewDB(sqldb, sqlitedialect.New())

	err = db.Ping()
	if err != nil {
		return err
	}

	_, err = db.NewCreateTable().Model((*User)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		return err
	}

	xlFile, err := xlsx.OpenFile(excelFile)
	if err != nil {
		return err
	}

	sheet := xlFile.Sheet[sheetName]
	rowsUpdated := 0
	rowsInserted := 0

	for i, row := range sheet.Rows {
		// skip header
		if i == 0 {
			continue
		}
		log.Printf("Row %d: %s", i, row.Cells[0].String())
		if row == nil {
			continue
		}

		eid := row.Cells[0].String()
		fname := row.Cells[1].String()
		lname := row.Cells[2].String()
		addr := row.Cells[3].String()

		var user User
		err = db.NewSelect().Model(&user).Where("id = ?", i).Scan(ctx)
		if err != nil {
			if err == sql.ErrNoRows {
				user = User{
					ID:    i,
					Eid:   eid,
					Fname: fname,
					Lname: lname,
					Addr:  addr,
				}
				_, err = db.NewInsert().Model(&user).Exec(ctx)
				if err != nil {
					return err
				}
				rowsInserted++
			} else {
				return err
			}
		} else {
			if user.Fname != fname || user.Lname != lname || user.Addr != addr {
				user.Eid = eid
				user.Fname = fname
				user.Lname = lname
				user.Addr = addr
				_, err = db.NewUpdate().Model(&user).Where("id = ?", user.ID).Exec(ctx)
				if err != nil {
					return err
				}
				rowsUpdated++
			}
		}
	}

	if rowsInserted == 0 && rowsUpdated == 0 {
		log.Println("No rows were inserted or updated.")
	} else {
		log.Printf("Rows inserted: %d\n", rowsInserted)
		log.Printf("Rows updated: %d\n", rowsUpdated)
	}

	return nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	excelFile := "users.xlsx"
	sheetName := "Sheet1_2"
	dbFile := "users.db"
	err := readExcelToSQLite(sheetName, excelFile, dbFile)
	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}
}

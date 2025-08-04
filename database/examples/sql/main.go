package main

import (
	"fmt"
	"log"
	"os"

	"github.com/grokify/mogo/config"
	"github.com/grokify/mogo/database/sqlutil"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/strconv/strconvutil"
)

func main() {
	files, err := config.LoadDotEnv(
		[]string{".env", os.Getenv("ENV_PATH")}, -1)
	if err != nil {
		log.Fatal(err)
	}
	fmtutil.MustPrintJSON(files)

	col, err := strconvutil.Atou32(os.Getenv("SQL_CSV_FILE_COL"))
	if err != nil {
		log.Fatal(err)
	}

	sqls, values, err := sqlutil.ReadFileCSVToSQLs(
		os.Getenv("SQL_FORMAT"),
		os.Getenv("SQL_CSV_FILE"),
		',',
		true,
		true,
		col)
	if err != nil {
		log.Fatal(err)
	}

	fmtutil.MustPrintJSON(sqls)

	fmt.Printf("SQL_ITEMS [%v]\n", len(values))
	fmt.Printf("SQL_STATEMENTS [%v]\n", len(sqls))

	fmt.Println("DONE")
}

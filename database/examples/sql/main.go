package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/grokify/mogo/config"
	"github.com/grokify/mogo/database/sql"
	"github.com/grokify/mogo/fmt/fmtutil"
)

func main() {
	files, err := config.LoadDotEnv(
		[]string{".env", os.Getenv("ENV_PATH")}, -1)
	if err != nil {
		log.Fatal(err)
	}
	fmtutil.PrintJSON(files)

	col, err := strconv.Atoi(os.Getenv("SQL_CSV_FILE_COL"))
	if err != nil {
		log.Fatal(err)
	}

	sqls, values, err := sql.ReadFileCSVToSQLs(
		os.Getenv("SQL_FORMAT"),
		os.Getenv("SQL_CSV_FILE"),
		',',
		true,
		true,
		uint(col))
	if err != nil {
		log.Fatal(err)
	}

	fmtutil.PrintJSON(sqls)

	fmt.Printf("SQL_ITEMS [%v]\n", len(values))
	fmt.Printf("SQL_STATEMENTS [%v]\n", len(sqls))

	fmt.Println("DONE")
}

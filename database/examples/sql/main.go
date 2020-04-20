package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/grokify/gotilla/config"
	"github.com/grokify/gotilla/database/sql"
	"github.com/grokify/gotilla/fmt/fmtutil"
)

func main() {
	files, err := config.LoadDotEnv(
		".env", os.Getenv("ENV_PATH"))
	if err != nil {
		log.Fatal(err)
	}
	fmtutil.PrintJSON(files)

	col, err := strconv.Atoi(os.Getenv("SQL_CSV_FILE_COL"))
	if err != nil {
		log.Fatal(err)
	}

	sqls, err := sql.ReadFileCSVToSQLs(
		os.Getenv("SQL_FORMAT"),
		os.Getenv("SQL_CSV_FILE"),
		",",
		true,
		true,
		uint(col))
	if err != nil {
		log.Fatal(err)
	}

	fmtutil.PrintJSON(sqls)

	fmt.Println("DONE")
}

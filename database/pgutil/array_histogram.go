package pgutil

import (
	"fmt"

	"github.com/go-pg/pg"
)

const (
	histogramSql = `SELECT %s AS name, COUNT(*) FROM %s, UNNEST(%s) AS t(%s) GROUP BY %s ORDER BY count DESC;`
	querySql     = `SELECT * FROM %s WHERE ? = ANY (%s);`
)

type TagInfo struct {
	Name  string
	Count int
}

func TagHistogramSQL(table, tags, tag string) string {
	return fmt.Sprintf(histogramSql, tag, table, tags, tag, tag)
}

func SelectTagHistogramResults(db *pg.DB, table, tags, tag string, res interface{}) error {
	stmt := TagHistogramSQL(table, tags, tag)
	_, err := db.Query(res, stmt)
	return err
}

func SelectTagHistogram(db *pg.DB, table, tags, tag string) ([]TagInfo, error) {
	var res []TagInfo
	return res, SelectTagHistogramResults(db, table, tags, tag, &res)
}

func SelectArrayAny(db *pg.DB, table, arrayCol, arrayVar string, res interface{}) error {
	stmt := fmt.Sprintf(querySql, table, arrayCol)
	_, err := db.Query(res, stmt, arrayVar)
	return err
}

package ynab

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"fmt"
	"strconv"
	"time"
)

type Entry struct {
	Date    time.Time
	Payee   string
	Memo    string
	Outflow sql.NullInt64 // in cents
	Inflow  sql.NullInt64 // in cents
}

func money(cents int64) string {
	d := cents / 100
	c := cents % 100
	return fmt.Sprintf("%s.%02s", strconv.FormatInt(d, 10), strconv.FormatInt(c, 10))
}

func record(e *Entry) []string {
	var in, out string
	if e.Outflow.Valid {
		out = money(e.Outflow.Int64)
	}
	if e.Inflow.Valid {
		in = money(e.Inflow.Int64)
	}
	return []string{
		e.Date.Format("2006-01-02"),
		e.Payee,
		e.Memo,
		out,
		in,
	}
}

// Make sure entries are sorted by date before calling Encode
func Encode(entries []Entry) ([]byte, error) {
	var out []byte
	buf := bytes.NewBuffer(out)
	w := csv.NewWriter(buf)
	w.Write([]string{
		"Date",
		"Payee",
		"Memo",
		"Outflow",
		"Inflow",
	})
	for _, entry := range entries {
		if err := w.Write(record(&entry)); err != nil {
			return nil, err
		}
	}

	// Write any buffered data to the underlying writer (standard output).
	w.Flush()

	if err := w.Error(); err != nil {
		return nil, err
	} else {
		return buf.Bytes(), nil
	}
}

package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kyleconroy/ynabnk/bnkdev"
	"github.com/kyleconroy/ynabnk/ynab"
)

func abs(n int) int64 {
	if n < 0 {
		return int64(-n)
	}
	return int64(n)
}

func main() {
	client := bnkdev.NewClient(os.Getenv("BNKDEV_API_KEY"))
	account := os.Getenv("BNKDEV_ACCOUNT_ID")
	ctx := context.Background()

	resp, err := client.ListTransactions(ctx, &bnkdev.ListTransactionsRequest{
		AccountID: account,
	})
	if err != nil {
		log.Fatal(err)
	}

	var entries []ynab.Entry
	for _, tx := range resp.Data {
		date, err := time.Parse("2006-01-02", tx.Date)
		if err != nil {
			log.Fatal(err)
		}
		entry := ynab.Entry{
			Date:  date,
			Payee: tx.AccountID,
			Memo:  tx.Description,
		}
		if tx.Amount >= 0 {
			entry.Inflow = sql.NullInt64{
				Int64: abs(tx.Amount),
				Valid: true,
			}
		} else {
			entry.Outflow = sql.NullInt64{
				Int64: abs(tx.Amount),
				Valid: true,
			}
		}
		entries = append(entries, entry)
	}

	out, err := ynab.Encode(entries)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(string(out))
}

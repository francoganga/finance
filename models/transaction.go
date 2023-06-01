package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Transaction struct {
	bun.BaseModel `bun:"transactions,alias:t"`

    ID int `bun:"id,pk,autoincrement" json:"-"`

	Date        time.Time
	Code        string
	Description string
	Amount      int
	Balance     int
}

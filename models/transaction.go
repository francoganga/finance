package models

import (
	"fmt"
	"time"

	"strings"

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
	Category    Category `bun:"rel:belongs-to"`
}

func (t *Transaction) AmountStr() string {
	var value float32

	value = float32(t.Amount) / 100

	strVal := fmt.Sprintf("%.2f", value)

	if value < 0 {
		return strings.ReplaceAll(strVal, "-", "-$")
	}

	return fmt.Sprintf("$%s", strVal)
}

func (t *Transaction) DateStr() string {

	return t.Date.Format("02-01-2006")
}

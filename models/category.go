package models

import "github.com/uptrace/bun"

type Category struct {
	bun.BaseModel `bun:"category,alias:c"`

	ID int `bun:"id,pk,autoincrement" json:"-"`

	Name string
}

package models

import "github.com/uptrace/bun"

type Category struct {
	bun.BaseModel `bun:"category,alias:c"`

	Name string
}

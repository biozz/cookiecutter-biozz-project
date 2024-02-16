package models

type Item struct {
	ID   string `bun:"id,pk"`
	Name string `bun:"name"`
}

package book

import "time"

type Book struct {
	Id            int
	Title         string
	Author        string
	ISBN          string
	Publisher     string
	PublishedYear int
	Pages         int
	Description   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

package model

import "time"

type File struct {
	Id        int64     `db:"id"`
	Name      string    `db:"name"`
	Servers   string    `db:"servers"`
	CreatedAt time.Time `db:"created_at"`
}

func (f File) Table() string {
	return "files"
}

func (f File) InsertColumns() []string {
	return []string{
		"name",
		"servers",
	}
}

func (f File) SelectColumns() []string {
	return []string{
		"id",
		"name",
		"servers",
		"created_at",
	}
}

func (f File) InsertValues() []any {
	return []any{
		f.Name,
		f.Servers,
	}
}

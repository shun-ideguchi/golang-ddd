package gorm

type Circle struct {
	ID      string
	Name    string
	Owner   string
	Members []User
}

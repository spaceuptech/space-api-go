package model

// Meta contains the meta information required to make a request
type Meta struct {
	DB, Col string
	Token   string
	Project string
}

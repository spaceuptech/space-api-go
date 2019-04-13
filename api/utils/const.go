package utils

const (
	// All is used when all the records needs to be worked on
	All string = "all"

	// One is used when oly a single record needs to be worked on
	One string = "one"

	// Count is used to count the number of documents returned
	Count string = "count"

	// Distinct is used to get the distinct values
	Distinct string = "distinct"

	// Upsert is used to upsert documents
	Upsert string = "upsert"
)

const (
	// Mongo is the constant for selecting MongoDB
	Mongo string = "mongo"

	// MySQL is the constant for selected MySQL
	MySQL string = "sql-mysql"

	// Postgres is the constant for selected Postgres
	Postgres string = "sql-postgres"
)

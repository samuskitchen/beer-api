package beer

const (
	// selectAllBeersTest is a query that selects all rows in the beers table
	selectAllBeersTest = "SELECT id, \"name\", brewery, country, price, currency, created_at FROM beers;"

	// selectBeerByIdTest is a query that selects a row from the beers table based off of the given id.
	// You must escape the code and to escape the code use
	// https://regex-escape.com/preg_quote-online.php
	selectBeerByIdTest = "SELECT id, \"name\", brewery, country, price, currency, created_at FROM beers WHERE id \\= \\$1;"

	// insertBeerWithIdTest is a query that inserts a new row in the beers table with a given id and using the values
	// given in order for id, "name", brewery, country, price, currency, created_at.
	// You must escape the code and to escape the code use
	// https://regex-escape.com/preg_quote-online.php
	insertBeerWithIdTest = "INSERT INTO beers \\(id, \"name\", brewery, country, price, currency, created_at\\) VALUES \\(\\$1, \\$2, \\$3, \\$4, \\$5, \\$6, \\$7\\) RETURNING id;"
)


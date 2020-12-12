package persistence

const (
	// selectAllBeers is a query that selects all rows in the beers table
	selectAllBeers = "SELECT id, \"name\", brewery, country, price, currency, created_at FROM beers;"

	// selectBeerById is a query that selects a row from the beers table based off of the given id.
	selectBeerById = "SELECT id, \"name\", brewery, country, price, currency, created_at FROM beers WHERE id = $1;"

	// insertBeerWithId is a query that inserts a new row in the beers table with a given id and using the values
	// given in order for id, "name", brewery, country, price, currency, created_at.
	insertBeerWithId = "INSERT INTO beers (id, \"name\", brewery, country, price, currency, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;"
)


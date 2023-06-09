package repository

const (
	save = `INSERT INTO account (first_name, last_name) 
					VALUES ($1, $2) 
					RETURNING *`

	update = `UPDATE account 
					SET first_name = COALESCE(NULLIF($1, ''), first_name),
						last_name = COALESCE(NULLIF($2, ''), last_name) 
					WHERE id = $3
					RETURNING *`

	deleteById = `DELETE FROM account WHERE id = $1`

	findById = `SELECT * 
				FROM account 
				WHERE id = $1`
)

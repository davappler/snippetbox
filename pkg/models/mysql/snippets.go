package mysql

import (
	"database/sql"

	"davappler/snippetbox/pkg/models"
)

// Define a SnippetModel type which wraps a sql.DB connection pool.
type SnippetModel struct { 
	DB *sql.DB
}
// This will insert a new snippet into the database.
func (m *SnippetModel) Insert(title, content, expires string) (int, error) { 



	stmt := `INSERT INTO snippets (title, content, created, expires)
			  VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`


	result, err := m.DB.Exec(stmt, title, content, expires) 
	if err != nil {
		return 0, err 
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err 
	}

	return int(id), nil
}
// This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (*models.Snippet, error) { 

	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id = ?`
	// Use the QueryRow() method on the connection pool to execute our
	// SQL statement, passing in the untrusted id variable as the value for the 
	// placeholder parameter. This returns a pointer to a sql.Row object which 
	// holds the result from the database.
	row := m.DB.QueryRow(stmt, id)

	// Initialize a pointer to a new zeroed Snippet struct.
	s := &models.Snippet{}
	// Use row.Scan() to copy the values from each field in sql.Row to the
	// corresponding field in the Snippet struct. Notice that the arguments
	// to row.Scan are *pointers* to the place you want to copy the data into,
	// and the number of arguments must be exactly the same as the number of
	// columns returned by your statement. If the query returns no rows, then
	// row.Scan() will return a sql.ErrNoRows error. We check for that and return 
	// our own models.ErrNoRecord error instead of a Snippet object.
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err == sql.ErrNoRows {
	return nil, models.ErrNoRecord } else if err != nil {
	return nil, err }
	// If everything went OK then return the Snippet object.
	return s, nil
}

// This will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*models.Snippet, error) { 
	return nil, nil
}
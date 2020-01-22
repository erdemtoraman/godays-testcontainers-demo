package demo1

import "github.com/jmoiron/sqlx"

type User struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type userRepo struct {
	conn *sqlx.DB
}

func NewRepo(conn *sqlx.DB) *userRepo {
	return &userRepo{conn: conn}
}

//noinspection ALL
func (r *userRepo) CreateUser(name string) (user User, err error) {
	err = r.conn.Get(&user, "INSERT INTO users(name) VALUES ($1) RETURNING *", name)
	return
}

//noinspection ALL
func (r *userRepo) GetUserByID(id int) (user User, err error) {
	err = r.conn.Get(&user, "SELECT * FROM  users WHERE  id = $1", id)
	return
}

//noinspection ALL
func (r *userRepo) GetAllUsers() (users []User, err error) {
	err = r.conn.Select(&users, "SELECT * FROM  users") // should be: .Select instead of .Get
	return
}

//noinspection ALL
func runMigrations(conn *sqlx.DB) error {
	_, err := conn.Exec(`
		CREATE TABLE IF NOT EXISTS users 
		(
		    id   SERIAL PRIMARY KEY, 
		    name TEXT NOT NULL UNIQUE
		)
	`)
	return err
}

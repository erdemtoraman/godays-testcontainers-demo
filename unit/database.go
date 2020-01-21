package unit

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
func (r *userRepo) CreateUser(name string) (User, error) {
	var user User
	err := r.conn.Get(&user, "INSERT INTO users(name) VALUES ($1) RETURNING *", name)
	return user, err
}

//noinspection ALL
func (r *userRepo) GetUserByID(id int) (User, error) {
	var user User
	err := r.conn.Get(&user, "SELECT * FROM  users WHERE  id = $1", id)
	return user, err
}

//noinspection ALL
func (r *userRepo) GetAllUsers() ([]User, error) {
	var users []User
	err := r.conn.Select(&users, "SELECT * FROM  users") // should be: .Select instead of .Get
	return users, err
}

func runMigrations(conn *sqlx.DB) error {
	_, err := conn.Exec("CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY,  name TEXT NOT NULL UNIQUE)")
	return err
}

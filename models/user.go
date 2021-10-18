package models

import (
	"database/sql"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type User struct {
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"-"`

	Groups []uuid.UUID `json:"groups"`

	ID        uuid.UUID `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// Model Functions
func (u *User) CheckPasswordHash(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) Create(c *Client) error {
	// encode password
	var (
		err          error
		passwordHash string
	)
	passwordHash, err = hashPassword(u.Password)
	if err != nil {
		return err
	}

	// add to database
	if u.ID == uuid.Nil {
		// id doesn't exist
		err = c.db.
			QueryRowx(`INSERT INTO "public"."users"("username", "password")
			VALUES ($1, $2) RETURNING id, created_at, updated_at;`, u.Username, passwordHash).
			Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
	} else {
		// id exists
		err = c.db.
			QueryRowx(`INSERT INTO "public"."users"("id", "username", "password")
			VALUES ($1, $2, $3) RETURNING created_at, updated_at;`, u.ID, u.Username, passwordHash).
			Scan(&u.CreatedAt, &u.UpdatedAt)
	}

	return err
}

func (u *User) SetPassword(password string) error {
	password, err := hashPassword(password)
	if err != nil {
		return err
	}

	u.Password = password
	return nil
}

// Client Functions
func (c *Client) ReadUser(id uuid.UUID) (*User, error) {
	var user User
	err := c.db.
		Get(&user, `SELECT id ,username, password, created_at, updated_at 
		FROM public.users WHERE id = $1 AND deleted_at IS NULL;`, id)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var groups []uuid.UUID
	err = c.db.
		Select(&groups, `SELECT group_id 
		FROM public.group_membership WHERE user_id = $1 AND deleted_at IS NULL;`, user.ID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	user.Groups = groups

	return &user, nil
}

func (c *Client) ReadUserByUsername(username string) (*User, error) {
	var user User
	err := c.db.
		Get(&user, `SELECT id ,username, password, created_at, updated_at 
		FROM public.users WHERE username = $1 AND deleted_at IS NULL;`, username)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var groups []uuid.UUID
	err = c.db.
		Select(&groups, `SELECT group_id 
		FROM public.group_membership WHERE user_id = $1 AND deleted_at IS NULL;`, user.ID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	user.Groups = groups

	return &user, nil
}

func (c *Client) ReadUsersPage(index, count int, orderBy string, asc bool) (*[]User, error) {
	var userList []User

	// build query
	query := "SELECT id, username, password, created_at, updated_at FROM users WHERE deleted_at IS NULL ORDER BY "

	switch strings.ToLower(orderBy) {
	case "created_at":
		query = query + "created_at "
	case "username":
		query = query + "username "
	default:
		return nil, ErrUnknownAttribute
	}

	if asc {
		query = query + "ASC "
	} else {
		query = query + "DESC "
	}

	query = query + "OFFSET $1 LIMIT $2;"

	// run query
	offset := index * count
	err := c.db.Select(&userList, query, offset, count)
	if err != nil && err != sql.ErrNoRows {
		logger.Errorf("cant get user page: %s")
		return nil, err
	}

	return &userList, nil
}

// Private Functions
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

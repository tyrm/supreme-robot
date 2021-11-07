package models

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/tyrm/supreme-robot/util"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

// User is used to login and keep authentication information
type User struct {
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"-"`

	Groups []uuid.UUID `json:"groups"`

	ID        uuid.UUID `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

// Model Functions

// AddGroup will add a group to the user and update update the database
func (u *User) AddGroup(c *Client, groups ...uuid.UUID) error {
	// start transaction
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}

	for _, group := range groups {
		logger.Tracef("adding group %s to %s", GroupTitle[group], u.Username)

		// add
		_, err = tx.
			Exec(`INSERT INTO "public"."group_membership"("user_id", "group_id")
			VALUES ($1, $2)`, u.ID, group)

		// rollback on error
		if err != nil {
			logger.Errorf("tx error: %s", err.Error())
			rberr := tx.Rollback()
			if rberr != nil {
				logger.Errorf("rollback error: %s", rberr.Error())
				// something went REALLY wrong
				return rberr
			}
			return err
		}

		u.Groups = append(u.Groups, group)
	}

	// commit transaction
	logger.Tracef("committing group memberships")
	err = tx.Commit()
	if err != nil {
		logger.Errorf("commit transaction: %s", err.Error())
		return err
	}

	return nil
}

// CheckPasswordHash is used to validate that a given password matches the stored hash
func (u *User) CheckPasswordHash(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) create(c *Client) error {
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
	err = c.db.
		QueryRowx(`INSERT INTO "public"."users"("username", "password")
			VALUES ($1, $2) RETURNING id, created_at, updated_at;`, u.Username, passwordHash).
		Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)

	return err
}

// IsMemberOfGroup checks if a user is in a given set of groups
func (u *User) IsMemberOfGroup(groups ...uuid.UUID) bool {
	return util.ContainsOneOfUUIDs(&u.Groups, &groups)
}

// SetPassword updates the user object's password hash
func (u *User) SetPassword(password string) error {
	password, err := hashPassword(password)
	if err != nil {
		return err
	}

	u.Password = password
	return nil
}

// Client Functions

// ReadUser will retrieve a user by their uuid from the database
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

// ReadUserByUsername will read a user by username from the database
func (c *Client) ReadUserByUsername(username string) (*User, error) {
	var user User
	err := c.db.
		Get(&user, `SELECT id ,username, password, created_at, updated_at 
		FROM public.users WHERE lower(username) = lower($1) AND deleted_at IS NULL;`, username)
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

// ReadUsersPage can retrieve a sorted and paginated list of users from the database
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
		return nil, errUnknownAttribute
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

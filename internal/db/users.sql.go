// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: users.sql

package q

import (
	"context"
	"time"
)

const createNewPasswordResetToken = `-- name: CreateNewPasswordResetToken :exec
INSERT INTO
    user_auth (
        user_id,
        username,
        password_reset_token
    )
VALUES ($1, $2, $3)
`

type CreateNewPasswordResetTokenParams struct {
	UserID             int64  `json:"userId"`
	Username           string `json:"username"`
	PasswordResetToken string `json:"passwordResetToken"`
}

func (q *Queries) CreateNewPasswordResetToken(ctx context.Context, arg *CreateNewPasswordResetTokenParams) error {
	_, err := q.db.Exec(ctx, createNewPasswordResetToken, arg.UserID, arg.Username, arg.PasswordResetToken)
	return err
}

const createOrUpdateUser = `-- name: CreateOrUpdateUser :one
INSERT INTO
    users (
        email,
        username,
        name,
        is_email_verified,
        provider,
        avatar,
        password,
        created_at,
        updated_at
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8,
        $9
    )
ON CONFLICT (email) DO
UPDATE
SET
    avatar = $6,
    password = $7,
    updated_at = $9
RETURNING
    id,
    username,
    name,
    email,
    alias_count,
    destination_count,
    is_premium,
    provider,
    avatar,
    password_changed_at,
    is_active,
    password,
    is_email_verified,
    created_at,
    updated_at
`

type CreateOrUpdateUserParams struct {
	Email           string    `json:"email"`
	Username        string    `json:"username"`
	Name            string    `json:"name"`
	IsEmailVerified bool      `json:"isEmailVerified"`
	Provider        string    `json:"provider"`
	Avatar          string    `json:"avatar"`
	Password        string    `json:"password"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type CreateOrUpdateUserRow struct {
	ID                int64     `json:"id"`
	Username          string    `json:"username"`
	Name              string    `json:"name"`
	Email             string    `json:"email"`
	AliasCount        int64     `json:"aliasCount"`
	DestinationCount  int64     `json:"destinationCount"`
	IsPremium         bool      `json:"isPremium"`
	Provider          string    `json:"provider"`
	Avatar            string    `json:"avatar"`
	PasswordChangedAt time.Time `json:"passwordChangedAt"`
	IsActive          bool      `json:"isActive"`
	Password          string    `json:"password"`
	IsEmailVerified   bool      `json:"isEmailVerified"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

func (q *Queries) CreateOrUpdateUser(ctx context.Context, arg *CreateOrUpdateUserParams) (*CreateOrUpdateUserRow, error) {
	row := q.db.QueryRow(ctx, createOrUpdateUser,
		arg.Email,
		arg.Username,
		arg.Name,
		arg.IsEmailVerified,
		arg.Provider,
		arg.Avatar,
		arg.Password,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i CreateOrUpdateUserRow
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Name,
		&i.Email,
		&i.AliasCount,
		&i.DestinationCount,
		&i.IsPremium,
		&i.Provider,
		&i.Avatar,
		&i.PasswordChangedAt,
		&i.IsActive,
		&i.Password,
		&i.IsEmailVerified,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const createUser = `-- name: CreateUser :exec

INSERT INTO
    users (
        username,
        name,
        email,
        password,
        avatar
    )
VALUES ($1, $2, $3, $4, $5)
`

type CreateUserParams struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Avatar   string `json:"avatar"`
}

// queries/users.sql
func (q *Queries) CreateUser(ctx context.Context, arg *CreateUserParams) error {
	_, err := q.db.Exec(ctx, createUser,
		arg.Username,
		arg.Name,
		arg.Email,
		arg.Password,
		arg.Avatar,
	)
	return err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteUser, id)
	return err
}

const findPasswordById = `-- name: FindPasswordById :one
SELECT password FROM users WHERE id = $1
`

func (q *Queries) FindPasswordById(ctx context.Context, id int64) (string, error) {
	row := q.db.QueryRow(ctx, findPasswordById, id)
	var password string
	err := row.Scan(&password)
	return password, err
}

const findUserByID = `-- name: FindUserByID :one
SELECT
    id,
    username,
    name,
    email,
    alias_count,
    destination_count,
    is_premium,
    provider,
    avatar,
    password_changed_at,
    is_active,
    is_email_verified,
    created_at,
    updated_at
FROM users
WHERE
    id = $1
`

type FindUserByIDRow struct {
	ID                int64     `json:"id"`
	Username          string    `json:"username"`
	Name              string    `json:"name"`
	Email             string    `json:"email"`
	AliasCount        int64     `json:"aliasCount"`
	DestinationCount  int64     `json:"destinationCount"`
	IsPremium         bool      `json:"isPremium"`
	Provider          string    `json:"provider"`
	Avatar            string    `json:"avatar"`
	PasswordChangedAt time.Time `json:"passwordChangedAt"`
	IsActive          bool      `json:"isActive"`
	IsEmailVerified   bool      `json:"isEmailVerified"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

func (q *Queries) FindUserByID(ctx context.Context, id int64) (*FindUserByIDRow, error) {
	row := q.db.QueryRow(ctx, findUserByID, id)
	var i FindUserByIDRow
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Name,
		&i.Email,
		&i.AliasCount,
		&i.DestinationCount,
		&i.IsPremium,
		&i.Provider,
		&i.Avatar,
		&i.PasswordChangedAt,
		&i.IsActive,
		&i.IsEmailVerified,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const findUserByUsername = `-- name: FindUserByUsername :one
SELECT
    id,
    username,
    name,
    email,
    alias_count,
    destination_count,
    is_premium,
    provider,
    avatar,
    password_changed_at,
    is_active,
    is_email_verified,
    created_at,
    updated_at
FROM users
WHERE
    username = $1
`

type FindUserByUsernameRow struct {
	ID                int64     `json:"id"`
	Username          string    `json:"username"`
	Name              string    `json:"name"`
	Email             string    `json:"email"`
	AliasCount        int64     `json:"aliasCount"`
	DestinationCount  int64     `json:"destinationCount"`
	IsPremium         bool      `json:"isPremium"`
	Provider          string    `json:"provider"`
	Avatar            string    `json:"avatar"`
	PasswordChangedAt time.Time `json:"passwordChangedAt"`
	IsActive          bool      `json:"isActive"`
	IsEmailVerified   bool      `json:"isEmailVerified"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

func (q *Queries) FindUserByUsername(ctx context.Context, username string) (*FindUserByUsernameRow, error) {
	row := q.db.QueryRow(ctx, findUserByUsername, username)
	var i FindUserByUsernameRow
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Name,
		&i.Email,
		&i.AliasCount,
		&i.DestinationCount,
		&i.IsPremium,
		&i.Provider,
		&i.Avatar,
		&i.PasswordChangedAt,
		&i.IsActive,
		&i.IsEmailVerified,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const findUserByUsernameOrEmail = `-- name: FindUserByUsernameOrEmail :one
SELECT
    id,
    username,
    name,
    email,
    alias_count,
    destination_count,
    is_premium,
    provider,
    avatar,
    password_changed_at,
    is_active,
    password,
    is_email_verified,
    created_at,
    updated_at
FROM users
WHERE
    username = $1
    OR email = $2
`

type FindUserByUsernameOrEmailParams struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type FindUserByUsernameOrEmailRow struct {
	ID                int64     `json:"id"`
	Username          string    `json:"username"`
	Name              string    `json:"name"`
	Email             string    `json:"email"`
	AliasCount        int64     `json:"aliasCount"`
	DestinationCount  int64     `json:"destinationCount"`
	IsPremium         bool      `json:"isPremium"`
	Provider          string    `json:"provider"`
	Avatar            string    `json:"avatar"`
	PasswordChangedAt time.Time `json:"passwordChangedAt"`
	IsActive          bool      `json:"isActive"`
	Password          string    `json:"password"`
	IsEmailVerified   bool      `json:"isEmailVerified"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

func (q *Queries) FindUserByUsernameOrEmail(ctx context.Context, arg *FindUserByUsernameOrEmailParams) (*FindUserByUsernameOrEmailRow, error) {
	row := q.db.QueryRow(ctx, findUserByUsernameOrEmail, arg.Username, arg.Email)
	var i FindUserByUsernameOrEmailRow
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Name,
		&i.Email,
		&i.AliasCount,
		&i.DestinationCount,
		&i.IsPremium,
		&i.Provider,
		&i.Avatar,
		&i.PasswordChangedAt,
		&i.IsActive,
		&i.Password,
		&i.IsEmailVerified,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const findUserByValidResetToken = `-- name: FindUserByValidResetToken :one
SELECT user_id
FROM user_auth
WHERE
    password_reset_token = $1
    AND password_reset_expires > now()
`

func (q *Queries) FindUserByValidResetToken(ctx context.Context, passwordResetToken string) (int64, error) {
	row := q.db.QueryRow(ctx, findUserByValidResetToken, passwordResetToken)
	var user_id int64
	err := row.Scan(&user_id)
	return user_id, err
}

const hasNoActiveResetToken = `-- name: HasNoActiveResetToken :one
SELECT id
FROM user_auth
WHERE
    user_id = $1
    AND password_reset_expires > now()
`

func (q *Queries) HasNoActiveResetToken(ctx context.Context, userID int64) (int64, error) {
	row := q.db.QueryRow(ctx, hasNoActiveResetToken, userID)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const updatePasswordAuth = `-- name: UpdatePasswordAuth :exec
UPDATE user_auth
SET
    password_reset_token = NULL,
    password_reset_expires = NULL
WHERE
    user_id = $1
`

func (q *Queries) UpdatePasswordAuth(ctx context.Context, userID int64) error {
	_, err := q.db.Exec(ctx, updatePasswordAuth, userID)
	return err
}

const updatePasswordUser = `-- name: UpdatePasswordUser :exec
UPDATE users
SET
    password = $1,
    password_changed_at = $2
WHERE
    id = $3
`

type UpdatePasswordUserParams struct {
	Password          string    `json:"password"`
	PasswordChangedAt time.Time `json:"passwordChangedAt"`
	ID                int64     `json:"id"`
}

func (q *Queries) UpdatePasswordUser(ctx context.Context, arg *UpdatePasswordUserParams) error {
	_, err := q.db.Exec(ctx, updatePasswordUser, arg.Password, arg.PasswordChangedAt, arg.ID)
	return err
}

const updateProviderByID = `-- name: UpdateProviderByID :exec
UPDATE users SET provider = $1 WHERE id = $2
`

type UpdateProviderByIDParams struct {
	Provider string `json:"provider"`
	ID       int64  `json:"id"`
}

func (q *Queries) UpdateProviderByID(ctx context.Context, arg *UpdateProviderByIDParams) error {
	_, err := q.db.Exec(ctx, updateProviderByID, arg.Provider, arg.ID)
	return err
}

const updateUser = `-- name: UpdateUser :exec
UPDATE users
SET
    name = $2,
    email = $3,
    avatar = $4,
    username = $5,
    provider = $6
WHERE
    id = $1
`

type UpdateUserParams struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	Username string `json:"username"`
	Provider string `json:"provider"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg *UpdateUserParams) error {
	_, err := q.db.Exec(ctx, updateUser,
		arg.ID,
		arg.Name,
		arg.Email,
		arg.Avatar,
		arg.Username,
		arg.Provider,
	)
	return err
}

const verifyEmailByID = `-- name: VerifyEmailByID :exec
UPDATE users SET is_email_verified = true WHERE id = $1
`

func (q *Queries) VerifyEmailByID(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, verifyEmailByID, id)
	return err
}

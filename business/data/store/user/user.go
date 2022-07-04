// Package user provides an example of a core business API. Right now these
// calls are just wrapping the data/data layer. But at some point you will
// want auditing or something that isn't specific to the data/store layer.
package user

import (
	"context"
	"fmt"
	"time"

	"github.com/Joggz/services/business/sys/database"
	"github.com/Joggz/services/business/sys/validate"
	"github.com/Joggz/services/business/web/auth"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// Store manages the set of APIs for user access.
type Store struct {
	log          *zap.SugaredLogger
	db           *sqlx.DB
}

// NewStore constructs a data for api access.
func NewStore(log *zap.SugaredLogger, db *sqlx.DB) Store {
	return Store{
		log: log,
		db:  db,
	}
}


// Create inserts a new user into the database.
func (s Store) Create(ctx context.Context, nu NewUser, now time.Time) (User, error) {
	if err := validate.Check(nu); err != nil {
		return User{}, fmt.Errorf("validating data: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("generating password hash: %w", err)
	}

	usr := User{
		ID:           validate.GenerateID(),
		Name:         nu.Name,
		Email:        nu.Email,
		PasswordHash: hash,
		Roles:        nu.Roles,
		DateCreated:  now,
		DateUpdated:  now,
	}

	const q = `
	INSERT INTO users
		(user_id, name, email, password_hash, roles, date_created, date_updated)
	VALUES
		(:user_id, :name, :email, :password_hash, :roles, :date_created, :date_updated)`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, usr); err != nil {
		return User{}, fmt.Errorf("inserting user: %w", err)
	}

	return usr, nil
}



// Update replaces a user document in the database.
func (s Store) Update(ctx context.Context, claims auth.Claims,  userID string, uu UpdateUser, now time.Time) error {
	if err :=	validate.CheckID(userID); err != nil {
		return database.ErrInvalidID
	}
	if err := validate.Check(uu); err != nil {
		return fmt.Errorf("Validating data : %w",err)
	}

	usr, err := s.QueryByID(ctx, claims, userID)
	if err != nil {
		return fmt.Errorf("updating user userID[%q]: %w", userID, err)
	}

	if uu.Name != nil {
		usr.Name = *uu.Name
	}
	if uu.Email != nil {
		usr.Email = *uu.Email
	}
	if uu.Roles != nil {
		usr.Roles = uu.Roles
	}
	if uu.Password != nil {
		pw, err := bcrypt.GenerateFromPassword([]byte(*uu.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("generating password hash: %w", err)
		}
		usr.PasswordHash = pw
	}
	usr.DateUpdated = now

	const q = `
	UPDATE
		users
	SET 
		"name" = :name,
		"email" = :email,
		"roles" = :roles,
		"password_hash" = :password_hash,
		"date_updated" = :date_updated
	WHERE
		user_id = :user_id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, usr); err != nil {
		return fmt.Errorf("updating userID[%s]: %w", usr.ID, err)
	}

	return nil
}



// Delete removes a user from the database.
func (s Store) Delete(ctx context.Context, claims auth.Claims, userID string) error {
	if err := validate.Check(userID); err != nil {
		return database.ErrInvalidID
	}

	// if you an not an admin and looking to delete someone other than yourself
	if !claims.Authorized(auth.RoleAdmin) && claims.Subject != userID {
		return database.ErrForbidden
	}

	data := struct {
		UserID string `db:"user_id"`
	}{
		UserID: userID,
	}

	const q = `
	DELETE FROM
		users
	WHERE
		user_id = :user_id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("deleting userID[%s]: %w", userID, err)
	}

	return nil
}


func (s Store) Query(ctx context.Context, pageNumber int, rowsPerPage int) ([]User, error) {
	data := struct{
		Offset int `db:"offset"`
		RowsPerPage int`db:"rows_per_page"`
	}{
		Offset: (pageNumber - 1) * rowsPerPage,
		RowsPerPage: rowsPerPage,
	}

	const q =`
	SELECT * users
	ORDER BY user_id
	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY
	`

	var users []User
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &users ); err != nil {
		if err == database.ErrNotFound {
			return nil, database.ErrNotFound
		}
	}

	return users, nil;
}



func (s Store) QueryByID(ctx context.Context, claims auth.Claims, userID string) (User, error) {
	if err := validate.Check(userID); err != nil {
		return User{}, database.ErrInvalidID
	}

	// if you an not an admin and looking to delete someone other than yourself
	if !claims.Authorized(auth.RoleAdmin) && claims.Subject != userID {
		return User{}, database.ErrForbidden
	}
	data := struct{
		UserID string `db:"user_id"`
	}{
		UserID: userID,
	}

	const q = `
		SELECT * users WHERE user_id = :user_id
	`

	var usr User
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &usr ); err != nil {
		if err == database.ErrNotFound {
			return User{}, database.ErrNotFound
		}
	}
	return usr, nil;
}




func (s Store) QueryByEmail(ctx context.Context, claims auth.Claims, email string) (User, error) {
	// if err := validate.Check(userID); err != nil {
	// 	return User{}, database.ErrInvalidID
	// }

	
	data := struct{
		Email string `db:"email"`
	}{
		Email: email,
	}

	const q = `
		SELECT * users WHERE email = :email
	`

	var usr User
	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &usr ); err != nil {
		if err == database.ErrNotFound {
			return User{}, database.ErrNotFound
		}
		return User{}, fmt.Errorf("selecting email[%q]: %w", email, err)
	}
	// if you an not an admin and looking to delete someone other than yourself
	if !claims.Authorized(auth.RoleAdmin) && claims.Subject != usr.ID {
		return User{}, database.ErrForbidden
	}
	return usr, nil;
}


func (s Store) Authenticate(ctx context.Context, now time.Time, email string, password string) (auth.Claims, error){
	data := struct{
		Email string `db:"email"`
	}{
		Email: email,
	}

	const q = ` SELECT * user WHERE email= :email`
	var usr User 

	if err := database.NamedQueryStruct(ctx,s.log,s.db, q, data, &usr); err != nil{
		if err == database.ErrNotFound{
			return auth.Claims{}, err
		}
		return auth.Claims{}, fmt.Errorf("selecting user[%q]: %w", email, err)
	}

	// compare password and usr.hash
	if err := bcrypt.CompareHashAndPassword(usr.PasswordHash, []byte(password));err != nil {
		return auth.Claims{}, database.ErrAuthenticationFailure
	}	

	// generate Claim for user

	claims := auth.Claims{
		StandardClaims: jwt.StandardClaims {
			Issuer: "service project",
			Subject: usr.ID,
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
			IssuedAt: time.Now().UTC().Unix(),
		},
		Role: usr.Roles,
	}
	return claims, nil
}
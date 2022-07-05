package user_test

import (
	"context"
	"testing"
	"time"

	"github.com/Joggz/services/business/data/dbtest"
	"github.com/Joggz/services/business/data/store/user"
	"github.com/Joggz/services/business/web/auth"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/go-cmp/cmp"
)

//  fix the error on the file 
var dbc = dbtest.DBContainer{
	Image: "postgres:14-alpine",
	Port:  "5432",
	Args:  []string{"-e", "POSTGRES_PASSWORD=password"},
}


func TestUser(t *testing.T) {
	log, db, teardown := dbtest.NewUnit(t, dbc)
	t.Cleanup(teardown)


	store := user.NewStore(log, db)

	t.Log("Given the need to work with user") 
	{
		testID := 0
		t.Logf("\tTest %d:\tWhen handling a single user.", testID)
		{
			ctx := context.Background()
			now := time.Date(2018, time.October, 1,0,0,0,0, time.UTC)

			nu := user.NewUser{
				Name:            "joggz swizz",
				Email:           "joggzswizz@gmail.com",
				Roles:           []string{auth.RoleAdmin},
				Password:        "gophers",
				PasswordConfirm: "gophers",
			}

		usr, err :=	store.Create(ctx, nu, now)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tshould be able to create user : %s.", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\t Test %d\tShould be able to create user: %s ", dbtest.Success, testID, usr)
			
			claims := auth.Claims{
				StandardClaims: jwt.StandardClaims {
					Issuer: "service project",
					Subject: usr.ID,
					ExpiresAt: time.Now().Add(time.Hour).Unix(),
					IssuedAt: time.Now().UTC().Unix(),
				},
				Role:[]string{auth.RoleUser},
			}
			}
			savedusr, err := store.QueryByID(ctx, claims, usr.ID)
			if err != nil {
				t.Fatalf("\t%s\tTest %d:\tshould be able to query user by id : %s.", dbtest.Failed, testID, err)
			}
			t.Fatalf("\t%s\tTest %d:\tshould be able to query user by id : %s.", dbtest.Success, testID, savedusr)

		if diff :=	cmp.Diff(usr, savedusr); diff !=  "" {
			t.Fatalf("\t%s\tTest %d:\tShould get back the same user. Diff:\n%s", dbtest.Failed, testID, diff)
		}
		}
	}
}


package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Joggz/services/app/services/sales-api/handlers"
	"github.com/Joggz/services/business/data/dbtest"
	"github.com/Joggz/services/business/web/auth"
)

// UserTests holds methods for each user subtest. This type allows passing
// dependencies for tests while still providing a convenient syntax when
// subtests are registered.
type UserTest  struct {
	app http.Handler
	userToken string
	adminToken string
	auth  *auth.Auth
}


// Test_Users is the entry point for testing user management functions.
func Test_Users(t *testing.T) {

	test := dbtest.NewIntegration(t, dbtest.DBContainer{
		Image: "postgres:14-alpine",
		Port:  "5432",
		Args:  []string{"-e", "POSTGRES_PASSWORD=password"},
	} )


	t.Cleanup(test.Teardown)


	shutdown := make(chan os.Signal, 1)

	tests := UserTest{
		app:        handlers.APIMux(handlers.APIMuxConfig{
			Shutdown: shutdown,
			Log:      test.Log,
			Auth:     test.Auth,
			DB:       test.DB,
		}),
		userToken:  test.Token("user@example.com", "gophers"),
		adminToken: test.Token("admin@example.com", "gophers"),
		auth: test.Auth,
	}


	t.Run("getToken200", tests.getToken200)
	t.Run("getToken404", tests.getToken404)
}

func (ut *UserTest) getToken200 (t *testing.T){
	r := httptest.NewRequest(http.MethodGet, "/v1/users/token", nil)
	w := httptest.NewRecorder()

	r.SetBasicAuth("admin@example.com", "gophers")
	ut.app.ServeHTTP(w, r)

	t.Log("Given the need to issues tokens to known users.", w.Result())


	{
		testID := 0
		t.Logf("\tTest %d:\tWhen fetching a token with valid credentials.", testID)
		{
			if w.Code != http.StatusOK {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 200 for the response : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 200 for the response.", dbtest.Success, testID)

			var got struct {
				Token string `json:"token"`
			}
		if	err := json.NewDecoder(w.Body).Decode(&got); err != nil {
			t.Fatalf("\t%s\tTest %d:\tShould be able to unmarshal the response : %v", dbtest.Failed, testID, err)
			}
			t.Logf("\t%s\tTest %d:\tShould be able to unmarshal the response.", dbtest.Success, testID)

		_, err :=	ut.auth.ValidateToken(got.Token)
		if err != nil {
			t.Fatalf("\t%s\tTest %d:\tShould be able to validate token from the response : %v", dbtest.Failed, testID, err)
		}
		t.Logf("\t%s\tTest %d:\tShould be able to validate token from the response.", dbtest.Success, testID)


		}
	}
}


func (ut *UserTest) getToken404(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/users/token", nil)
	w := httptest.NewRecorder()

	t.Log("Given the need to issues tokens to known users.", w.Result())


	r.SetBasicAuth("unknown@example.com", "gophamint")
	ut.app.ServeHTTP(w, r)

	{
		testID := 0
		t.Logf("\tTest %d:\tWhen fetching a token with an unrecognized email.", testID, )
		{
			if w.Code != http.StatusNotFound {
				t.Fatalf("\t%s\tTest %d:\tShould receive a status code of 404 for the response : %v", dbtest.Failed, testID, w.Code)
			}
			t.Logf("\t%s\tTest %d:\tShould receive a status code of 404 for the response.", dbtest.Success, testID)
		}

	}
}
package users

import (
	"net/http"

	"github.com/pressly/warpdrive/services"
	"github.com/pressly/warpdrive/web"
)

type createUser struct {
	Name     *string `json:"name,required"`
	Email    *string `json:"email,required"`
	Password *string `json:"password,required"`
}

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	search := query.Get("q")

	users := services.QueryUsersByEmail(search)

	web.Respond(w, http.StatusOK, users)
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := web.ParamAsInt64(r, "userId")

	if err != nil {
		web.Respond(w, http.StatusBadRequest, err)
		return
	}

	user := services.FindUserByID(userID)

	if user == nil {
		web.Respond(w, http.StatusNotFound, nil)
		return
	}

	web.Respond(w, http.StatusOK, user)
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	body := ctx.Value("parsed:body").(*createUser)

	name := *body.Name
	email := *body.Email
	password := *body.Password

	user, err := services.CreateUser(name, email, password)
	if err != nil {
		web.Respond(w, http.StatusBadRequest, err)
		return
	}

	web.Respond(w, http.StatusOK, user)
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`¯\_(ツ)_/¯`))
}

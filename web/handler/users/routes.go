package users

import (
	"github.com/pressly/chi"
	"github.com/pressly/warpdrive/web/middleware"
)

func Routes() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.JwtHandler())

	r.Post("/", middleware.BodyParser(createUserRequestBuilder, 512), createUserHandler)
	r.Delete("/:userId", deleteUserHandler)

	return r
}

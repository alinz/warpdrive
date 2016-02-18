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

	r.Patch("/", middleware.BodyParser(updateUserRequestBuilder, 512), updateUserHandler)
	r.Patch("/:userId", middleware.BodyParser(updateUserRequestBuilder, 512), updateUserHandler)

	return r
}

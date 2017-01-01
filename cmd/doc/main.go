package main

import (
	"github.com/pressly/chi/docgen"
	"github.com/pressly/warpdrive/web/routes"
)

func main() {
	r := routes.New(true)

	//docgen.JSONRoutesDoc(r)
	docgen.PrintRoutes(r)
}

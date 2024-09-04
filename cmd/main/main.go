package main

import (
	_ "JwtWithGo/docs"
	"JwtWithGo/pkg/config"
	"JwtWithGo/pkg/routes"
	"github.com/swaggo/http-swagger"
	"log"
	"net/http"
)

// @title Your Project API
// @version 1.0
// @description This is a sample server for a bookstore application.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8000
func main() {
	db := config.InitDB()
	r := routes.SetupRoutes(db)
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	log.Fatal(http.ListenAndServe(":8000", r))
}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.

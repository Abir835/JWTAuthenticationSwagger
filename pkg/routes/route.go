package routes

import (
	"JwtWithGo/pkg/controllers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
)

// HealthCheck godoc
// @Summary Health Check
// @Description Check the health of the application.
// @Tags Health
// @Success 200 {string} string "OK"
// @Router / [get]
func HealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Health is OK!"))
	}
}

func SetupRoutes(db *gorm.DB) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", HealthCheck()).Methods("GET")
	r.HandleFunc("/register", controllers.Register(db)).Methods("POST")
	r.HandleFunc("/login", controllers.Login(db)).Methods("POST")
	r.HandleFunc("/refresh", controllers.RefreshToken(db)).Methods("POST")
	//r.HandleFunc("/logout", controllers.Logout(db)).Methods("POST")

	return r
}

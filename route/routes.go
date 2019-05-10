package route

import (
	"fmt"
	"net/http"

	"github.com/yusufpapurcu/Library/models"

	"github.com/gorilla/mux"
)

var routes = models.Routes{ //Route List
	models.Route{
		Name:        "Create User API",
		Path:        "/api/user/create",
		Method:      "Post",
		HandlerFunc: CreateUser,
	}, // Duzenlemek gerekiyor Test Basarili
	models.Route{
		Name:        "Auth API",
		Path:        "/api/user/auth",
		Method:      "Post",
		HandlerFunc: Auth,
	}, // Duzenleme Yapildi. *Ufak Tefek sorunlari var*
	models.Route{
		Name:        "Find One User API",
		Path:        "/api/user/getone",
		Method:      "Post",
		HandlerFunc: FindUser,
	}, // Test edildi. Problem Yok
	models.Route{
		Name:        "Find All Users API",
		Path:        "/api/user/getall",
		Method:      "Post",
		HandlerFunc: FindAllUser,
	}, // Test Edildi Problem yok

	models.Route{
		Name:        "Create Book API",
		Path:        "/api/book/create",
		Method:      "Post",
		HandlerFunc: CreateBook,
	}, // Duzenleme Yapildi. *Ufak Tefek sorunlari var*
	models.Route{
		Name:        "Find One Book API",
		Path:        "/api/book/getone",
		Method:      "Post",
		HandlerFunc: FindBook,
	}, // Duzenleme Yapildi. *Ufak Tefek sorunlari var*
	models.Route{
		Name:        "Find All Books API",
		Path:        "/api/book/getall",
		Method:      "Post",
		HandlerFunc: FindAllBook,
	}, // Duzenleme Yapildi. *Ufak Tefek sorunlari var*
	models.Route{
		Name:        "Find One Book API",
		Path:        "/api/book/borrow",
		Method:      "Post",
		HandlerFunc: Borrow,
	},
	models.Route{
		Name:        "Find All Books API",
		Path:        "/api/book/deliver",
		Method:      "Post",
		HandlerFunc: Deliver,
	},

	models.Route{
		Name:        "Index",
		Path:        "/",
		Method:      "Get",
		HandlerFunc: Empty,
	},
}

// SetRouter Function for create Api End Points.
func SetRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true) // Create router
	router.Use(JwtAuthentication)
	for _, route := range routes { //Adding route list's property
		router.
			Methods(route.Method).
			Path(route.Path).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	return router // Return router
}

func Empty(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Yakinda Acilacak")
}

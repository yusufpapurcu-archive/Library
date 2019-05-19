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
	}, //Test Basarili
	models.Route{
		Name:        "Create User API",
		Path:        "/api/create/admin",
		Method:      "Post",
		HandlerFunc: CreateAdmin,
	}, //Test
	models.Route{
		Name:        "Find One User API",
		Path:        "/api/user/getone",
		Method:      "Post",
		HandlerFunc: FindUser,
	}, // Test Basarili
	models.Route{
		Name:        "Find All Users API",
		Path:        "/api/user/getall",
		Method:      "Post",
		HandlerFunc: FindAllUser,
	}, // Test Basarili

	models.Route{
		Name:        "Create Book API",
		Path:        "/api/book/create",
		Method:      "Post",
		HandlerFunc: CreateBook,
	}, // Test Basarili
	models.Route{
		Name:        "Find One Book API",
		Path:        "/api/book/getone",
		Method:      "Post",
		HandlerFunc: FindBook,
	}, // Test Basarili
	models.Route{
		Name:        "Find All Books API",
		Path:        "/api/book/getall",
		Method:      "Post",
		HandlerFunc: FindAllBook,
	}, // Test Basarili
	models.Route{
		Name:        "Find One Book API",
		Path:        "/api/book/borrow",
		Method:      "Post",
		HandlerFunc: Borrow,
	}, // Test Basarili
	models.Route{
		Name:        "Find All Books API",
		Path:        "/api/book/deliver",
		Method:      "Post",
		HandlerFunc: Deliver,
	}, // Test Basarili
	models.Route{
		Name:        "Index",
		Path:        "/",
		Method:      "Get",
		HandlerFunc: Index,
	},
	models.Route{
		Name:        "Login",
		Path:        "/login",
		Method:      "Get",
		HandlerFunc: Login,
	},
	models.Route{
		Name:        "Auth API",
		Path:        "/api/user/auth",
		Method:      "Post",
		HandlerFunc: Auth,
	}, // Test Basarili
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

package models

import "net/http"

// Route struct. This Route for Routes list.
type Route struct {
	Name        string           // Route Name  		        Example : "UserRegister"
	Path        string           // Route Path  		        Example : "/register"
	Method      string           // Route Method 		        Example : "Get"
	HandlerFunc http.HandlerFunc // Route HandlerFunc       For Route Method
}

// Routes Type for Route array
type Routes = []Route

package route

import (
	"fmt"
	"html/template"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("pages/index.html")
	if err != nil {
		fmt.Fprint(w, "hata")
	}
	t.Execute(w, nil)
}
func Login(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("pages/login.html")
	if err != nil {
		fmt.Fprint(w, "hata")
	}
	t.Execute(w, nil)
}

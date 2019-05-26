package route

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/yusufpapurcu/Library/models"
	u "github.com/yusufpapurcu/Library/utils"
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
func MainPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("pages/main.html")
	if err != nil {
		fmt.Fprint(w, "hata")
	}
	t.Execute(w, nil)
}

func MainLoader(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&req) //decode the request body into struct and failed if any error occur
	if err != nil {
		fmt.Println(err)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	tokenHeader := req["token"]
	if tokenHeader == "" {
		t, err := template.ParseFiles("pages/main.admin.html")
		if err != nil {
			fmt.Fprint(w, "hata")
		}
		var temp bytes.Buffer
		t.Execute(&temp, nil)
		a := temp.String()
		out, err := json.Marshal(&a) // Convert the Json
		if err != nil {
			fmt.Println(err)
		}
		fmt.Fprint(w, string(out))
		return
	}
	tk := &models.Token{}

	token, err := jwt.ParseWithClaims(tokenHeader.(string), tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("token_password")), nil
	})
	if err != nil { //Malformed token, returns with http code 403 as usual
		t, err := template.ParseFiles("pages/main.admin.html")
		if err != nil {
			fmt.Fprint(w, "hata")
		}
		var temp bytes.Buffer
		t.Execute(&temp, nil)
		a := temp.String()
		out, err := json.Marshal(&a) // Convert the Json
		if err != nil {
			fmt.Println(err)
		}
		fmt.Fprint(w, string(out))
		return
	}
	if !token.Valid { //Token is invalid, maybe not signed on this server
		t, err := template.ParseFiles("pages/cf.html")
		if err != nil {
			fmt.Fprint(w, "hata")
		}
		var temp bytes.Buffer
		t.Execute(&temp, nil)
		a := temp.String()
		out, err := json.Marshal(&a) // Convert the Json
		if err != nil {
			fmt.Println(err)
		}
		fmt.Fprint(w, string(out))
		return
	}
	user, err := models.GetUser(tk.UserId)
	if err != nil {
		response := u.Message(false, err.Error())
		w.WriteHeader(http.StatusForbidden)
		w.Header().Add("Content-Type", "application/json")
		u.Respond(w, response)
		return
	}
	if user.Admin {
		t, err := template.ParseFiles("pages/main.admin.html")
		if err != nil {
			fmt.Fprint(w, "hata")
		}
		var temp bytes.Buffer
		t.Execute(&temp, nil)
		a := temp.String()
		out, err := json.Marshal(&a) // Convert the Json
		if err != nil {
			fmt.Println(err)
		}
		fmt.Fprint(w, string(out))
		return
	} else {
		t, err := template.ParseFiles("pages/main.user.html")
		if err != nil {
			fmt.Fprint(w, "hata")
		}
		var temp bytes.Buffer
		t.Execute(&temp, nil)
		a := temp.String()
		out, err := json.Marshal(&a) // Convert the Json
		if err != nil {
			fmt.Println(err)
		}
		fmt.Fprint(w, string(out))
		return
	}
}

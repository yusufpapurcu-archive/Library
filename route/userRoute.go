package route

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yusufpapurcu/Library/models"
	u "github.com/yusufpapurcu/Library/utils"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {

	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	id := r.Context().Value("user")
	resp := user.Create(id) //Create account
	u.Respond(w, resp)
}
func CreateAdmin(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	token := r.Header.Get("Authorization")
	resp := user.CreateAdmin(token) //Create account
	u.Respond(w, resp)
}

func Auth(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	c := user.SchoolTag
	fmt.Println(user)
	resp := models.Login(user.Email, user.Password, c)
	u.Respond(w, resp)
}

func FindUser(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	if user.ID.Hex() == "000000000000000000000000" {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	usra, err := models.GetUser(user.ID.Hex())
	if err != nil {
		u.Respond(w, u.Message(false, "Kullanici Bulunamadi"))
		return
	}
	out, err := json.Marshal(&usra) // Convert the Json
	if err != nil {
		u.Respond(w, u.Message(false, "Json Convert Error."))
		return
	}
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, string(out))
}

func FindAllUser(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	res := user.FindAllUser()
	out, err := json.Marshal(res) // Convert the Json
	if err != nil {
		u.Respond(w, u.Message(false, "Json Convert Error."))
		return
	}
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, string(out))
}

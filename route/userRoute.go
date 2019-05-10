package route

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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
	fmt.Println(id)
	resp := user.Create(id) //Create account
	u.Respond(w, resp)
}

func Auth(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	c := r.Header.Get("SchoolTag")
	if c == "" { //Token is missing, returns with error code 403 Unauthorized
		response := u.Message(false, "Missing auth token")
		w.WriteHeader(http.StatusForbidden)
		w.Header().Add("Content-Type", "application/json")
		u.Respond(w, response)
		return
	}

	splitted := strings.Split(c, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
	if len(splitted) != 2 {
		response := u.Message(false, "Invalid/Malformed auth token")
		w.WriteHeader(http.StatusForbidden)
		w.Header().Add("Content-Type", "application/json")
		u.Respond(w, response)
		return
	}
	st := []string{splitted[0], splitted[1]}
	resp := models.Login(user.Email, user.Password, st)
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
	out, err := json.Marshal(&usra) // Convert the Json
	if err != nil {
		u.Respond(w, u.Message(false, "Json Convert Error."))
		return
	}
	fmt.Fprint(w, string(out))
}

func FindAllUser(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	out := user.FindAllUser()

	fmt.Fprint(w, out)
}

package route

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yusufpapurcu/Library/models"
	u "github.com/yusufpapurcu/Library/utils"
)

func CreateBook(w http.ResponseWriter, r *http.Request) {
	book := &models.Book{}
	err := json.NewDecoder(r.Body).Decode(book) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	id := r.Context().Value("user")
	fmt.Println(id)
	out, ok := book.CreateBook(id.(string))
	if !ok {
		fmt.Fprint(w, out)
		return
	}
	fmt.Fprint(w, out)
	return
}

func FindBook(w http.ResponseWriter, r *http.Request) {
	book := &models.Book{}
	err := json.NewDecoder(r.Body).Decode(book) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	if book.ID.Hex() == "000000000000000000000000" {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	booka := models.GetBook(book.ID)
	if booka == nil {
		fmt.Fprint(w, "Can't Find Anything")
		return
	}
	out, err := json.Marshal(&booka) // Convert the Json
	if err != nil {
		u.Respond(w, u.Message(false, "Json Convert Error."))
		return
	}
	fmt.Fprint(w, string(out))
}

func FindAllBook(w http.ResponseWriter, r *http.Request) {
	book := &models.Book{}
	err := json.NewDecoder(r.Body).Decode(book) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	out := book.FindAllBook()

	fmt.Fprint(w, out)
}

func Borrow(w http.ResponseWriter, r *http.Request) {
	ubid := &models.UBdecoder{}
	err := json.NewDecoder(r.Body).Decode(ubid) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	book := models.GetBook(ubid.BID)
	if book == nil {
		fmt.Fprint(w, "Can't Find Book")
		return
	}
	user, err := models.GetUser(ubid.UID.Hex())
	if user == nil {
		fmt.Fprint(w, "Can't Find Book. Error : "+err.Error())
		return
	}
	id := r.Context().Value("user")
	out := user.Borrow(book, id)
	fmt.Fprint(w, out)
}

func Deliver(w http.ResponseWriter, r *http.Request) {
	ubid := &models.UBdecoder{}
	err := json.NewDecoder(r.Body).Decode(ubid) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	book := models.GetBook(ubid.BID)
	if book == nil {
		fmt.Fprint(w, "Can't Find Book")
		return
	}
	user, err := models.GetUser(ubid.UID.Hex())
	if user == nil {
		fmt.Fprint(w, "Can't Find Book. Error : "+err.Error())
		return
	}
	id := r.Context().Value("user")
	out := user.Deliver(book, id)
	fmt.Fprint(w, out)
}

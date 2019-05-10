package models

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/yusufpapurcu/Library/database"
	u "github.com/yusufpapurcu/Library/utils"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`                // User ID
	Name       string             `bson:"name,omitempty" json:"name,omitempty"`             // User Name
	Email      string             `bson:"email,omitempty" json:"email,omitempty"`           // User Email
	Password   string             `bson:"password,omitempty" json:"password,omitempty"`     // User Password
	SchoolTag  []string           `bson:"schooltag,omitempty" json:"schooltag,omitempty"`   // User School Tag
	Read       []Book             `bson:"read,omitempty" json:"read,omitempty"`             // User Readed Books
	Delivery   []Book             `bson:"delivery,omitempty" json:"delivery,omitempty"`     // User Delivery Books
	Token      string             `json:"token" bson:"-"`                                   // User Token
	Admin      bool               `json:"admin,omitempty" bson:"admin,omitempty"`           // User Admin
	CreatedAt  time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`   // CreatedAt User
	ModifiedAt time.Time          `bson:"modifiedAt,omitempty" json:"modifiedAt,omitempty"` // CreatedAt User
}

type UBdecoder struct {
	BID primitive.ObjectID `json:"bid"`
	UID primitive.ObjectID `json:"uid`
}

func (user *User) Create(c interface{}) map[string]interface{} {
	usra := &User{}
	id, err := primitive.ObjectIDFromHex(c.(string))
	if err != nil {
		fmt.Println("Line 45 : " + err.Error())
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second) // Context for Serach
	err = database.GetDB("user").FindOne(ctx, bson.M{"_id": id}).Decode(&usra)
	fmt.Println(usra)
	if err != nil && err.Error() != "mongo: no documents in result" {
		return u.Message(false, "Failed to create account, connection error.")
	}
	if !usra.Admin {
		return u.Message(false, "Failed to create account, Please Be Admin.")
	}
	if resp, ok := user.Validate(); !ok {
		return resp
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	user.SchoolTag = usra.SchoolTag
	user.Read = []Book{}
	user.Delivery = []Book{}
	user.CreatedAt = time.Now()
	user.ModifiedAt = time.Now()
	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second) // Context for Create function
	res, err := database.GetDB("user").InsertOne(ctx, user)           // Create Function
	if err != nil {
		return u.Message(false, "Failed to create account, connection error.")
	}
	out, err := json.Marshal(&res.InsertedID) // Convert the Json
	if err != nil {
		return u.Message(false, "Json Transform Problem.")
	}
	tk := &Token{UserId: string(out)}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	user.Token = tokenString

	user.Password = "" //delete password

	response := u.Message(true, "Account has been created")
	response["user"] = user
	return response
}

func (user *User) Validate() (map[string]interface{}, bool) {

	if !strings.Contains(user.Email, "@") {
		return u.Message(false, "Email address is required"), false
	}

	if len(user.Password) < 6 {
		return u.Message(false, "Password is required"), false
	}

	//Email must be unique
	temp := &User{}
	filter := bson.M{"email": user.Email}
	//check for errors and duplicate emails
	err := database.GetDB("user").FindOne(context.TODO(), filter).Decode(&temp)
	if err != nil && err.Error() != "mongo: no documents in result" {
		return u.Message(false, "Connection error. Please retry"), false
	}
	if temp.Email != "" {
		return u.Message(false, "Email address already in use by another user."), false
	}

	return u.Message(false, "Requirement passed"), true
}

func Login(email, password string, c []string) map[string]interface{} {

	user := &User{}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second) // Context for Serach
	err := database.GetDB("user").FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil && err.Error() != "mongo: no documents in result" {
		fmt.Println(err)
		return u.Message(false, "Connection error. Please retry")
	}
	fmt.Println(user)
	if user.Email == "" {
		fmt.Println(user)
		return u.Message(false, "Email address not found")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return u.Message(false, "Invalid login credentials. Please try again")
	}
	//Worked! Logged In
	user.Password = ""
	tk := &Token{UserId: user.ID.Hex()}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	user.Token = tokenString

	user.Password = "" //delete password

	resp := u.Message(true, "Logged In")
	resp["user"] = user
	return resp
}

func GetUser(u string) (*User, error) {

	id, err := primitive.ObjectIDFromHex(u)
	if err != nil {
		return nil, err
	}
	user := &User{}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second) // Context for Serach
	err = database.GetDB("user").FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	user.Password = ""
	return user, nil
}

func (user *User) Borrow(b *Book, c interface{}) map[string]interface{} {
	usra := &User{}
	id, err := primitive.ObjectIDFromHex(c.(string))
	if err != nil {
		return u.Message(false, "Invalid request")
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second) // Context for Serach
	err = database.GetDB("user").FindOne(ctx, bson.M{"_id": id}).Decode(&usra)
	if err != nil {
		return u.Message(false, err.Error())
	}
	if !usra.Admin {
		return u.Message(false, "Failed to borrow book, Please Be Admin.")
	}
	if len(user.Delivery) >= 3 {
		return u.Message(false, "Failed to borrow book, You can Borrow 3 Books Maksimum.")
	}
	b.User = *user
	err = b.Update()
	if err != nil {
		return u.Message(false, err.Error())
	}
	user.Delivery = append(user.Delivery, *b)
	err = user.Update()
	if err != nil {
		return u.Message(false, err.Error())
	}
	return u.Message(true, "Borrow Complate")
}

func (user *User) Deliver(book *Book, c interface{}) map[string]interface{} {
	usra := &User{}
	id, err := primitive.ObjectIDFromHex(c.(string))
	if err != nil {
		return u.Message(false, "Invalid request")
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second) // Context for Serach
	err = database.GetDB("user").FindOne(ctx, bson.M{"_id": id}).Decode(&usra)
	if err != nil {
		return nil
	}
	if !usra.Admin {
		return u.Message(false, "Failed to borrow book, Please Be Admin.")
	}
	for i, b := range user.Delivery {
		if b.ID == book.ID {
			user.Delivery = append(user.Delivery[:i], user.Delivery[i+1:]...)
			book.User = User{}
			err = user.Update()
			if err != nil {
				return u.Message(false, err.Error())
			}
			err = book.Update()
			if err != nil {
				return u.Message(false, err.Error())
			}
			return u.Message(true, "Deliver is Succesful.")
		}
	}
	return u.Message(false, "Failed to borrow book, This User Haven't This book.")
}

func (user User) FindAllUser() map[string]interface{} {
	var filtre bson.M                // Filtre variable created
	bytes, err := bson.Marshal(user) // User convert the []byte
	if err != nil {
		return u.Message(false, "Failed to Serach, filter create error.")
	}
	bson.Unmarshal(bytes, &filtre)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second) // Context for Find function
	cur, err := database.GetDB("user").Find(ctx, filtre)                // Empty Find function
	if err != nil && err.Error() != "mongo: no documents in result" {
		return u.Message(false, "Failed to Serach, connection error.")
	}
	defer cur.Close(ctx) // Close cursor
	var result User      // Create Models.User for Decode result
	var list []string    // Decoded and Transformed Json's String
	for cur.Next(ctx) {  // Loop all cursors
		err := cur.Decode(&result) // Decode
		if err != nil {
			return u.Message(false, "Failed to Serach, Decode Error.")
		}
		out, err := json.Marshal(result) // Transform JSON
		if err != nil {
			return u.Message(false, "Failed to Serach, Convert JSON Error.")
		}
		list = append(list, string(out)) // Storage the string array
	}
	if err := cur.Err(); err != nil {
		return u.Message(false, "Failed to Serach, Cursor error.")
	}
	msg := u.Message(true, "Serach Succesful")
	msg["users"] = list
	return msg
}

func (user User) Update() error {
	update := bson.M{"$set": user}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)              // Context for Update function
	_, err := database.GetDB("user").UpdateOne(ctx, bson.M{"_id": user.ID}, update) // Update Document
	if err != nil {
		return err
	}
	return nil
}

// TODO : Update Yazilacak

/*
	usra := &User{}
	filtre := bson.D{{"_id", c}}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second) // Context for Serach
	err := database.GetDB("user").FindOne(ctx, filtre).Decode(&usra)
	if err != nil && err.Error() != "mongo: no documents in result" {
		return u.Message(false, "Failed to create account, connection error.")
	}
	if !usra.Admin {
		return u.Message(false, "Failed to create account, Please Be Admin.")
	}
*/

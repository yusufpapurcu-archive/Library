package models

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	u "github.com/yusufpapurcu/Library/utils"

	"github.com/yusufpapurcu/Library/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Book Struct for Manage and Store Books
type Book struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`                  // Book ID
	Name            string             `bson:"name,omitempty" json:"name,omitempty"`               // Book Name
	PrintedYear     int                `bson:"printedyear,omitempty" json:"printedyear,omitempty"` // Printed Year of Book
	Author          string             `bson:"author,omitempty" json:"author,omitempty"`           // Author Book
	PublishingHouse string             `bson:"pubhouse,omitempty" json:"pubhouse,omitempty"`       // Publishing House Book
	TimeLeft        time.Time          `bson:"timeleft,omitempty" json:"timeleft,omitempty"`       // Timeleft Book
	User            User               `bson:"user,omitempty" json:"user,omitempty"`
	No              int                `bson:"no,omitempty" json:"no,omitempty"`
	SchoolTag       []string           `bson:"schooltag,omitempty" json:"schooltag,omitempty"`   // User School Tag
	CreatedAt       time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`   // CreatedAt Book
	ModifiedAt      time.Time          `bson:"modifiedAt,omitempty" json:"modifiedAt,omitempty"` // ModifiedAt Book
}

func (book *Book) CreateBook(c string) (map[string]interface{}, bool) {
	user := &User{}
	id, err := primitive.ObjectIDFromHex(c)
	if err != nil {
		fmt.Println("Line 45 : " + err.Error())
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second) // Context for Serach
	err = database.GetDB("user").FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil && err.Error() != "mongo: no documents in result" {
		return u.Message(false, "Failed to create account, connection error."), false
	}
	book.SchoolTag = user.SchoolTag
	temp := &Book{}
	filter := bson.M{"no": book.No, "schooltag": book.SchoolTag}
	//check for errors and duplicate emails
	err = database.GetDB("book").FindOne(context.TODO(), filter).Decode(&temp)
	fmt.Println(temp.ID.Hex())
	if err != nil && err.Error() != "mongo: no documents in result" {
		return u.Message(false, "Connection error. Please retry"), false
	}
	if temp.ID.Hex() == "000000000000000000000000" {
		if user.Admin {
			book.CreatedAt = time.Now()
			book.CreatedAt = time.Now()
			book.TimeLeft = time.Now()
			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second) // Context for Create function
			_, err := database.GetDB("book").InsertOne(ctx, book)              // Create Function
			if err != nil {
				return u.Message(false, "Failed to book account, connection error."), false
			}
			return u.Message(true, "Book Succesfully Created."), true
		}
		return u.Message(false, "Be Admin Please."), false
	}
	return u.Message(false, "No already in use by another book."), false
}

func (b Book) FindAllBook() map[string]interface{} {
	var filtre bson.M             // Filtre variable created
	bytes, err := bson.Marshal(b) // User convert the []byte
	if err != nil {
		return u.Message(false, "Failed to Serach, filter create error.")
	}
	bson.Unmarshal(bytes, &filtre)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second) // Context for Find function
	cur, err := database.GetDB("book").Find(ctx, filtre)                // Empty Find function
	if err != nil && err.Error() != "mongo: no documents in result" {
		return u.Message(false, "Failed to Serach, connection error.")
	}
	defer cur.Close(ctx) // Close cursor
	var result Book      // Create Models.User for Decode result
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
	msg["books"] = list
	return msg
}

func GetBook(b interface{}) *Book {

	book := &Book{}
	filtre := bson.D{{"_id", b}}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second) // Context for Serach
	err := database.GetDB("book").FindOne(ctx, filtre).Decode(&book)
	if err != nil {
		return nil
	}
	return book
}
func (book Book) Update() error {
	update := bson.M{"$set": book}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)              // Context for Update function
	_, err := database.GetDB("book").UpdateOne(ctx, bson.M{"_id": book.ID}, update) // Update Document
	if err != nil {
		return err
	}
	return nil
}

// TODO : Update Yazilacak

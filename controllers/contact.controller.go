package controllers

import (
	"contact/handlers"
	models "contact/models"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ContactHandler struct {
	Collection   *mongo.Collection
	errorHandler *handlers.ErrorHandler
}

func NewContactCollection(collection *mongo.Collection, errorHandler *handlers.ErrorHandler) *ContactHandler {
	return &ContactHandler{
		Collection:   collection,
		errorHandler: errorHandler,
	}
}

// creates a new contact
func (h *ContactHandler) CreateContact(c *gin.Context) {
	var contact models.Contact

	if err := c.BindJSON(&contact); err != nil {
		h.errorHandler.HandleBadRequest(c) // Use errorHandler to handle bad request
		return
	}

	_, err := h.Collection.InsertOne(c, contact)
	if err != nil {
		h.errorHandler.HandleInternalServerError(c) // Use ErrorHandler to handle internal server error
		return
	}
	c.JSON(201, contact)
}

// gets all contacts

func (h *ContactHandler) GetContacts(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := h.Collection.Find(ctx, bson.M{})

	if err != nil {
		h.errorHandler.HandleInternalServerError(c) // Use ErrorHandler to handle internal server error
		return
	}
	defer cursor.Close(ctx)

	var contacts []models.Contact

	if err := cursor.All(ctx, &contacts); err != nil {
		h.errorHandler.HandleInternalServerError(c) // Use ErrorHandler to handle internal server error
		return
	}
	c.JSON(200, contacts)
}

// get contact by id
func (h *ContactHandler) GetContactByID(c *gin.Context) {
	id := c.Param("id")

	var contact models.Contact
	cxt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := h.Collection.FindOne(cxt, bson.M{"_id": id}).Decode(&contact)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			h.errorHandler.HandleNotFound(c) // Use ErrorHandler to handle not found error
			return
		}
		h.errorHandler.HandleInternalServerError(c) // Use ErrorHandler to handle internal server error
		return
	}
	c.JSON(200, contact)
}

//upadate contact updates a contact by id

func (h *ContactHandler) UpdateContact(c *gin.Context) {
	id := c.Params.ByName("id")
	var updateContact models.Contact

	if err := c.BindJSON(&updateContact); err != nil {
		h.errorHandler.HandleBadRequest(c) // Use ErrorHandler to handle bad request
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.M{"_id": id}
	update := bson.M{"$set": updateContact}
	_, err := h.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		h.errorHandler.HandleInternalServerError(c) // Use ErrorHandler to handle internal server error
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "contact updated"})
}

//delete contact deletes a contact by id

func (h *ContactHandler) DeleteContact(c *gin.Context) {
	id := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.M{"_id": id}
	_, err := h.Collection.DeleteOne(ctx, filter)
	if err != nil {
		h.errorHandler.HandleInternalServerError(c) // Use ErrorHandler to handle internal server error
		return
	}
	c.JSON(200, gin.H{"message": "contact deleted"})
}

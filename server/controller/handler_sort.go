package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/harry713j/minurly/config"
	"github.com/harry713j/minurly/helper"
	"github.com/harry713j/minurly/models"
	"github.com/harry713j/minurly/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ShortUrlType struct {
	ShortenCode string `json:"shortCode"`
}

type OriginalUrlType struct {
	OriginalUrl string `json:"originalUrl"`
}

func HandleCreateUrl(w http.ResponseWriter, r *http.Request) {
	session, err := config.SessionStore.Get(r, "sessionId")

	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userId, ok := session.Values["userId"].(string)

	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid or expired session")
		return
	}

	var incomingUrl OriginalUrlType
	err = json.NewDecoder(r.Body).Decode(&incomingUrl)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	randomChars, err := utils.GenerateRandomStrings(8)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Unable to create a short url")
		return
	}

	hexUserId, err := bson.ObjectIDFromHex(userId)

	if err != nil {
		log.Println("Unable to convert to mongo objectId")
		utils.RespondWithError(w, http.StatusInternalServerError, "Unable to create url")
		return
	}

	// create the short url
	url := &models.ShortUrl{
		OriginalUrl: incomingUrl.OriginalUrl,
		Visits:      0,
		ShortCode:   randomChars, // created by some crypto algo
		UserId:      hexUserId,
		CreatedAt:   time.Now(),
		LastVisited: time.Now(),
	}

	urlId, err := helper.InsertOneUrl(*url)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Unable to create a short url")
		return
	}

	fmt.Println("New Url id: ", urlId)
	utils.RespondWithJSON(w, http.StatusCreated, ShortUrlType{
		ShortenCode: randomChars,
	})

}

func HandleGetUrl(w http.ResponseWriter, r *http.Request) {
	// extract the path and run a database query and fetch the original url and send back to front-end
	vars := mux.Vars(r)
	shortCode := vars["short-code"]

	if shortCode == "" || len(shortCode) != 8 {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid parameters provided")
		return
	}

	shortUrlDoc, err := helper.FindOneUrlByShort(shortCode)

	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "No original url found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, OriginalUrlType{
		OriginalUrl: shortUrlDoc.OriginalUrl,
	})

}

func HandleDeleteUrl(w http.ResponseWriter, r *http.Request) {
	session, err := config.SessionStore.Get(r, "sessionId")

	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userId, ok := session.Values["userId"].(string)

	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid or expired session")
		return
	}

	vars := mux.Vars(r)
	shortCode := vars["short-code"]

	if shortCode == "" || len(shortCode) != 8 {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid parameters provided")
		return
	}

	hexUserId, err := bson.ObjectIDFromHex(userId)

	if err != nil {
		log.Println("Unable to convert to mongo objectId")
		utils.RespondWithError(w, http.StatusInternalServerError, "Unable to remove url")
		return
	}

	_, err = helper.DeleteUrlByShort(shortCode, hexUserId)

	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "No original url found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, "Successfully Deleted")
}

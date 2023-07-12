package controller

import (
	"context"
	"encoding/json"
	"main/model"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func AddAddress() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(status{"error": "user_id query is empty"})
			return
		}

		var address model.Address
		err := json.NewDecoder(r.Body).Decode(&address)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(status{"error": "invalid json format:"})
			return
		}

		err = validate.Struct(address)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(status{"error": "failed to validate json:"})
			return
		}
		address.AddressID = uuid.NewString()

		var foundUser model.User
		_, err = userCollection.ReadDocument(context.TODO(), userID, &foundUser)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(status{"error": "failed to fetch user"})
			return
		}

		addresses := foundUser.AddressDetails
		if len(addresses) > 2 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(status{"error": "can't add more than two addresses"})
			return
		}
		addresses = append(addresses, address)

		updateObject := make(map[string]interface{})
		updateObject["address_details"] = addresses

		_, err = userCollection.UpdateDocument(context.TODO(), userID, updateObject)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(status{"error": "failed to update user item"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("successfully added address")
	}
}

func EditHomeAddress() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(status{"error": "user_id query is empty"})
			return
		}

		var address model.Address
		err := json.NewDecoder(r.Body).Decode(&address)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(status{"error": "invalid json format:"})
			return
		}

		var foundUser model.User
		_, err = userCollection.ReadDocument(context.TODO(), userID, &foundUser)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(status{"error": "failed to fetch user"})
			return
		}

		if len(foundUser.AddressDetails) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(status{"error": "address_details empty"})
			return
		}

		address.AddressID = foundUser.AddressDetails[0].AddressID
		if address.City == nil {
			address.City = foundUser.AddressDetails[0].City
		}
		if address.House == nil {
			address.House = foundUser.AddressDetails[0].House
		}
		if address.Street == nil {
			address.Street = foundUser.AddressDetails[0].Street
		}
		if address.PostCode == nil {
			address.Street = foundUser.AddressDetails[0].PostCode
		}

		addresses := foundUser.AddressDetails
		addresses[0] = address
		updateObject := make(map[string]interface{})
		updatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		updateObject["address_details"] = addresses
		updateObject["updated_at"] = updatedAt

		_, err = userCollection.UpdateDocument(context.TODO(), userID, updateObject)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(status{"error": "failed to update user item"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("successfully edited home address")
	}
}

func EditWorkAddress() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(status{"error": "user_id query is empty"})
			return
		}

		var address model.Address
		err := json.NewDecoder(r.Body).Decode(&address)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(status{"error": "invalid json format:"})
			return
		}

		var foundUser model.User
		_, err = userCollection.ReadDocument(context.TODO(), userID, &foundUser)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(status{"error": "failed to fetch user"})
			return
		}

		if len(foundUser.AddressDetails) != 2 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(status{"error": "work address is empty"})
			return
		}

		address.AddressID = foundUser.AddressDetails[1].AddressID
		if address.City == nil {
			address.City = foundUser.AddressDetails[1].City
		}
		if address.House == nil {
			address.House = foundUser.AddressDetails[1].House
		}
		if address.Street == nil {
			address.Street = foundUser.AddressDetails[1].Street
		}
		if address.PostCode == nil {
			address.Street = foundUser.AddressDetails[1].PostCode
		}

		addresses := foundUser.AddressDetails
		addresses[1] = address
		updateObject := make(map[string]interface{})
		updatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		updateObject["address_details"] = addresses
		updateObject["updated_at"] = updatedAt

		_, err = userCollection.UpdateDocument(context.TODO(), userID, updateObject)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(status{"error": "failed to update user item"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("successfully edited work address")
	}
}

func DeleteAddresses() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(status{"error": "user_id query is empty"})
			return
		}

		addresses := make([]model.Address, 0)
		updatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		updateObject := make(map[string]interface{})
		updateObject["address_details"] = addresses
		updateObject["updated_at"] = updatedAt

		_, err := userCollection.UpdateDocument(context.TODO(), userID, updateObject)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(status{"error": "failed to update user item"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("successfully deleted address_details")
	}
}

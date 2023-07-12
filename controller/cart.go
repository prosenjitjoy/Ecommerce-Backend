package controller

import (
	"context"
	"encoding/json"
	"main/database"
	"main/model"
	"net/http"
)

func AddToCart() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productQueryID := r.URL.Query().Get("product_id")
		userQueryID := r.URL.Query().Get("user_id")

		if productQueryID == "" || userQueryID == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(status{"error": "product_id or user_id query is empty"})
			return
		}

		err := database.AddProductToCart(prodCollection, userCollection, productQueryID, userQueryID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(status{"error": "failed to add product to the cart"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Successfully added to the cart")
	}
}

func RemoveItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productQueryID := r.URL.Query().Get("product_id")
		userQueryID := r.URL.Query().Get("user_id")

		if productQueryID == "" || userQueryID == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(status{"error": "product_id or user_id query is empty"})
			return
		}

		err := database.RemoveCartItem(prodCollection, userCollection, productQueryID, userQueryID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(status{"error": "failed to remove product from the cart"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Successfully removed item from cart")
	}
}

func GetItemFromCart() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(status{"error": "user_id query is empty"})
			return
		}

		var foundUser model.User
		_, err := userCollection.ReadDocument(context.TODO(), userID, &foundUser)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(status{"error": "failed to fetch user"})
			return
		}

		totalPrice := 0
		for _, v := range foundUser.UserCart {
			totalPrice += v.Price
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(status{"payment_due": totalPrice, "user_cart": foundUser.UserCart})
	}
}

func BuyFromCart() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userQueryID := r.URL.Query().Get("user_id")

		if userQueryID == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(status{"error": "product_id or user_id query is empty"})
			return
		}

		err := database.BuyItemFromCart(prodCollection, userCollection, userQueryID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(status{"error": "failed to placed the order"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Successfully placed the order")
	}
}

func InstantBuy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productQueryID := r.URL.Query().Get("product_id")
		userQueryID := r.URL.Query().Get("user_id")

		if productQueryID == "" || userQueryID == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(status{"error": "product_id or user_id query is empty"})
			return
		}

		err := database.InstantBuyer(prodCollection, userCollection, productQueryID, userQueryID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(status{"error": "failed to placed order"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Successfully placed the order")
	}
}

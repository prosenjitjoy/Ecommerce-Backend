package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"main/model"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func AddProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := CheckUserType(r.Context(), "ADMIN")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(status{"error": "unauthorized user"})
			return
		}

		var product model.Product
		err = json.NewDecoder(r.Body).Decode(&product)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(status{"error": "invalid json format:"})
			return
		}

		err = validate.Struct(product)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(status{"error": "failed to validate json:"})
			return
		}

		product.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		product.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		product.ProductID = uuid.NewString()

		meta, err := prodCollection.CreateDocument(context.TODO(), product)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(status{"error": "failed to create product item"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(meta.Key)
	}
}

func EditProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := CheckUserType(r.Context(), "ADMIN")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(status{"error": "unauthorized user"})
			return
		}

		productID := r.URL.Query().Get("product_id")
		if productID == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(status{"error": "product_id query is empty"})
			return
		}

		var product model.Product
		err = json.NewDecoder(r.Body).Decode(&product)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(status{"error": "invalid json format:"})
			return
		}

		updateObject := make(map[string]interface{})

		if product.ProductName != nil {
			updateObject["product_name"] = product.ProductName
		}
		if product.Price != nil {
			updateObject["price"] = product.Price
		}
		if product.Rating != nil {
			updateObject["product_image"] = product.Rating
		}
		if product.Image != nil {
			updateObject["image"] = product.Image
		}

		product.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObject["updated_at"] = product.UpdatedAt

		meta, err := prodCollection.UpdateDocument(context.TODO(), productID, updateObject)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(status{"error": "failed to update product item"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(meta.Key)
	}
}

func DeleteProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := CheckUserType(r.Context(), "ADMIN")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(status{"error": "unauthorized user"})
			return
		}

		productID := r.URL.Query().Get("product_id")
		if productID == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(status{"error": "product_id query is empty"})
			return
		}

		meta, err := prodCollection.RemoveDocument(context.TODO(), productID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(status{"error": "failed to delete product item"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(meta.Key)
	}
}

func CheckUserType(c context.Context, role string) error {
	userType := c.Value("user_type")
	if userType != role {
		return fmt.Errorf("unauthorized to access this resource")
	}
	return nil
}

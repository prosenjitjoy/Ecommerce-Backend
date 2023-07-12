package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"main/database"
	"main/model"
	"net/http"

	"github.com/arangodb/go-driver"
	"github.com/go-playground/validator/v10"
)

var (
	validate       *validator.Validate = validator.New()
	db             driver.Database     = database.DBinstance()
	userCollection driver.Collection   = database.OpenCollection(db, "users")
	prodCollection driver.Collection   = database.OpenCollection(db, "products")
)

type status map[string]interface{}

func GetProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := "FOR product IN products RETURN product"
		cursor, err := db.Query(context.TODO(), query, nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(status{"error": "failed to execute query"})
			return
		}
		defer cursor.Close()

		var productList []model.Product
		for {
			var product model.Product
			_, err := cursor.ReadDocument(context.TODO(), &product)

			if driver.IsNoMoreDocuments(err) {
				break
			} else if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(status{"error": "failed to read product item"})
				return
			}

			productList = append(productList, product)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(productList)
	}
}

func GetProductByQuery() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productName := r.URL.Query().Get("product_name")
		if productName == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(status{"error": "product_name is empty"})
			return
		}

		query := fmt.Sprintf(`FOR product IN products FILTER product.product_name == "%s" RETURN product`, productName)
		cursor, err := db.Query(context.TODO(), query, nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(status{"error": err.Error()})
			return
		}
		defer cursor.Close()

		var searchProducts []model.Product
		for {
			var product model.Product
			_, err := cursor.ReadDocument(context.TODO(), &product)

			if driver.IsNoMoreDocuments(err) {
				break
			} else if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(status{"error": "failed to read product items"})
				return
			}

			searchProducts = append(searchProducts, product)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(searchProducts)
	}
}

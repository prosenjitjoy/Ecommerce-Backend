package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"main/model"
	"main/token"
	"net/http"
	"time"

	"github.com/arangodb/go-driver"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(status{"error": "invalid json format:"})
			return
		}

		err = validate.Struct(user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(status{"error": "failed to validate json:"})
			return
		}

		query := fmt.Sprintf(`FOR user IN users FILTER user.email == "%s" RETURN user`, *user.Email)
		cursor, err := db.Query(context.TODO(), query, nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(status{"error": "failed to count email:"})
			return
		}
		defer cursor.Close()
		emailCount := cursor.Count()

		query = fmt.Sprintf(`FOR user IN users FILTER user.phone == "%s" RETURN user`, *user.Phone)
		cursor, err = db.Query(context.TODO(), query, nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(status{"error": "failed to count phone:"})
			return
		}
		defer cursor.Close()
		phoneCount := cursor.Count()

		if emailCount > 0 || phoneCount > 0 {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(status{"error": "duplicate email or phone:"})
			return
		}

		password := HashPassword(*user.Password)
		user.Password = &password
		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UserID = uuid.NewString()

		token, refreshToken, err := jwttoken.GenerateAllTokens(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(status{"error": "failed to generate token:"})
			return
		}

		user.Token = &token
		user.RefreshToken = &refreshToken
		user.UserCart = []model.ProductUser{}
		user.AddressDetails = []model.Address{}
		user.OrderStatus = []model.Order{}

		meta, err := userCollection.CreateDocument(context.TODO(), user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(status{"error": "failed to insert data:"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(meta.Key)
	}
}

func Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(status{"error": "invalid json format:"})
			return
		}

		if user.Email == nil || user.Password == nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(status{"error": "enter email and password"})
			return
		}

		query := fmt.Sprintf(`FOR user IN users FILTER user.email == "%s" RETURN user`, *user.Email)

		cursor, err := db.Query(context.TODO(), query, nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(status{"error": "failed to find and map user:"})
			return
		}
		defer cursor.Close()

		users := []model.User{}
		for {
			var user model.User
			_, err := cursor.ReadDocument(context.TODO(), &user)

			if driver.IsNoMoreDocuments(err) {
				break
			} else if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(status{"error": "failed to read menu items"})
				return
			}

			users = append(users, user)
		}

		// verify password
		err = bcrypt.CompareHashAndPassword([]byte(*users[0].Password), []byte(*user.Password))
		if err != nil {

			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(status{"error": "failed to verify password:"})
			return
		}

		token, refreshToken, err := jwttoken.GenerateAllTokens(users[0])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(status{"error": "failed to generate tokens"})
			return
		}

		err = jwttoken.UpdateAllTokens(token, refreshToken, users[0].UserID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(status{"error": "failed to update tokens" + err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(users[0])
	}
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

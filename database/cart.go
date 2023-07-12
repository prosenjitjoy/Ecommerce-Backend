package database

import (
	"context"
	"errors"
	"main/model"
	"time"

	"github.com/arangodb/go-driver"
	"github.com/google/uuid"
)

func AddProductToCart(prodCollection, userCollection driver.Collection, productID, userID string) (err error) {
	var product model.ProductUser
	_, err = prodCollection.ReadDocument(context.TODO(), productID, &product)
	if err != nil {
		err = errors.New("can't find the product")
		return
	}

	var foundUser model.User
	_, err = userCollection.ReadDocument(context.TODO(), userID, &foundUser)
	if err != nil {
		err = errors.New("can't find user")
		return
	}

	foundUser.UserCart = append(foundUser.UserCart, product)

	updateObject := make(map[string]interface{})
	updateObject["user_cart"] = foundUser.UserCart

	_, err = userCollection.UpdateDocument(context.TODO(), userID, updateObject)
	if err != nil {
		err = errors.New("failed to update user cart")
		return
	}

	return nil
}

func RemoveCartItem(prodCollection, userCollection driver.Collection, productID, userID string) (err error) {
	var product model.ProductUser
	_, err = prodCollection.ReadDocument(context.TODO(), productID, &product)
	if err != nil {
		err = errors.New("can't find the product")
		return
	}

	var foundUser model.User
	_, err = userCollection.ReadDocument(context.TODO(), userID, &foundUser)
	if err != nil {
		err = errors.New("can't find user")
		return
	}

	var productCart []model.ProductUser
	for _, v := range foundUser.UserCart {
		if v.ProductID == productID {
			continue
		}
		productCart = append(productCart, v)
	}

	updateObject := make(map[string]interface{})
	updateObject["user_cart"] = productCart

	_, err = userCollection.UpdateDocument(context.TODO(), userID, updateObject)
	if err != nil {
		err = errors.New("failed to update user cart")
		return
	}

	return nil
}

func BuyItemFromCart(prodCollection, userCollection driver.Collection, userID string) (err error) {
	var foundUser model.User
	_, err = userCollection.ReadDocument(context.TODO(), userID, &foundUser)
	if err != nil {
		err = errors.New("can't find user")
		return
	}

	totalPrice := 0
	for _, v := range foundUser.UserCart {
		totalPrice += v.Price
	}

	var order model.Order
	order.OrderID = uuid.NewString()
	order.OrderedAt = time.Now()
	order.OrderCart = foundUser.UserCart
	order.PaymentMethod.COD = true
	order.TotalPrice = totalPrice

	foundUser.OrderStatus = append(foundUser.OrderStatus, order)
	emptyUserCart := make([]model.ProductUser, 0)

	updateObject := make(map[string]interface{})
	updateObject["order_status"] = foundUser.OrderStatus
	updateObject["user_cart"] = emptyUserCart

	_, err = userCollection.UpdateDocument(context.TODO(), userID, updateObject)
	if err != nil {
		err = errors.New("failed to update order cart")
		return
	}

	return nil
}

func InstantBuyer(prodCollection, userCollection driver.Collection, productID, userID string) (err error) {
	var foundProduct model.ProductUser
	_, err = prodCollection.ReadDocument(context.TODO(), productID, &foundProduct)
	if err != nil {
		err = errors.New("can't find product")
		return
	}

	var order model.Order
	order.OrderID = uuid.NewString()
	order.OrderedAt = time.Now()
	order.OrderCart = []model.ProductUser{foundProduct}
	order.PaymentMethod.COD = true
	order.TotalPrice = foundProduct.Price

	var foundUser model.User
	_, err = userCollection.ReadDocument(context.TODO(), userID, &foundUser)
	if err != nil {
		err = errors.New("can't find user")
		return
	}

	foundUser.OrderStatus = append(foundUser.OrderStatus, order)

	updateObject := make(map[string]interface{})
	updateObject["order_status"] = foundUser.OrderStatus

	_, err = userCollection.UpdateDocument(context.TODO(), userID, updateObject)
	if err != nil {
		err = errors.New("failed to update order cart")
		return
	}

	return nil
}

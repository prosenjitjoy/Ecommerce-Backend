package model

import (
	"time"
)

// user model
type User struct {
	UserID         string        `json:"_key"`
	FirstName      *string       `json:"first_name" validate:"required,min=3,max=30"`
	LastName       *string       `json:"last_name" validate:"required,min=3,max=30"`
	Email          *string       `json:"email" validate:"required,email"`
	Password       *string       `json:"password" validate:"required,min=3"`
	Phone          *string       `json:"phone" validate:"required,min=11,max=11"`
	UserType       *string       `json:"user_type" validate:"required,eq=ADMIN|eq=USER"`
	Token          *string       `json:"token"`
	RefreshToken   *string       `json:"refresh_token"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
	UserCart       []ProductUser `json:"user_cart"`
	AddressDetails []Address     `json:"address_details"`
	OrderStatus    []Order       `json:"order_status"`
}

// product model
type Product struct {
	ProductID   string    `json:"_key"`
	ProductName *string   `json:"product_name" validate:"required"`
	Price       *uint64   `json:"price" validate:"required"`
	Rating      *uint8    `json:"rating"`
	Image       *string   `json:"image" validate:"required"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// productuser model
type ProductUser struct {
	ProductID   string  `json:"_key"`
	ProductName *string `json:"product_name"`
	Price       int     `json:"price"`
	Rating      *uint8  `json:"rating"`
	Image       *string `json:"image"`
}

// address model
type Address struct {
	AddressID string  `json:"_key"`
	House     *string `json:"house" validate:"required"`
	Street    *string `json:"street" validate:"required"`
	City      *string `json:"city" validate:"required"`
	PostCode  *string `json:"post_code" validate:"required"`
}

// order model
type Order struct {
	OrderID       string        `json:"_key"`
	OrderCart     []ProductUser `json:"order_cart"`
	OrderedAt     time.Time     `json:"ordered_at"`
	TotalPrice    int           `json:"total_price"`
	Discount      *int          `json:"discount"`
	PaymentMethod Payment       `json:"payment_method"`
}

type Payment struct {
	Digital bool
	COD     bool
}

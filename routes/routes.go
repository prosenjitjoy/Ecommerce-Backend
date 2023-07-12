package routes

import (
	"main/controller"
	"main/middleware"

	"github.com/go-chi/chi/v5"
)

func Use(router *chi.Mux) {
	// public routes
	router.Group(func(r chi.Router) {
		r.Post("/user/signup", controller.Register())
		r.Post("/user/signin", controller.Login())
		r.Get("/user/products", controller.GetProducts())
		r.Get("/user/search", controller.GetProductByQuery())
	})

	// protected routes
	router.Group(func(r chi.Router) {
		r.Use(middleware.Authenticator)

		r.Post("/admin/addproduct", controller.AddProduct())
		r.Patch("/admin/editproduct", controller.EditProduct())
		r.Delete("/admin/deleteproduct", controller.DeleteProduct())
		r.Get("/addtocart", controller.AddToCart())
		r.Get("/removeitem", controller.RemoveItem())
		r.Get("/listcart", controller.GetItemFromCart())
		r.Post("/addaddress", controller.AddAddress())
		r.Patch("/edithomeaddress", controller.EditHomeAddress())
		r.Patch("/editworkaddress", controller.EditWorkAddress())
		r.Delete("/deleteaddresses", controller.DeleteAddresses())
		r.Get("/cartcheckout", controller.BuyFromCart())
		r.Get("/instantbuy", controller.InstantBuy())
	})
}

package api

import (
	"a21hc3NpZ25tZW50/model"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func (api *API) AddCart(w http.ResponseWriter, r *http.Request) {
	// Get username context to struct model.Cart.
	// username := "" // TODO: replace this
	username := fmt.Sprintf("%s", r.Context().Value("username"))

	r.ParseForm()

	// Check r.Form with key product, if not found then return response code 400 and message "Request Product Not Found".
	// TODO: answer here
	if r.FormValue("product") == "" {
		var x model.ErrorResponse
		x.Error = "Request Product Not Found"
		output, _ := json.Marshal(x)
		w.WriteHeader(400)
		w.Write(output)
		return
	}

	var list []model.Product
	var float float64
	for _, formList := range r.Form {
		for _, v := range formList {
			item := strings.Split(v, ",")
			p, _ := strconv.ParseFloat(item[2], 64)
			q, _ := strconv.ParseFloat(item[3], 64)
			total := p * q
			float += total
			list = append(list, model.Product{
				Id:       item[0],
				Name:     item[1],
				Price:    item[2],
				Quantity: item[3],
				Total:    total,
			})
		}
	}

	// Add data field Name, Cart and TotalPrice with struct model.Cart.
	carts := model.Cart{
		Name:       username,
		Cart:       list,
		TotalPrice: float64(float),
	} // TODO: replace this

	_, err := api.cartsRepo.CartUserExist(carts.Name)
	if err != nil {
		api.cartsRepo.AddCart(carts)
	} else {
		api.cartsRepo.UpdateCart(carts)
	}
	api.dashboardView(w, r)

}

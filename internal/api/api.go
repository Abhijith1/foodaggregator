package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Abhijith01/foodaggregator/internal/item"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// NewRouter creates all the handler functions and returns gorilla mux's router
func NewRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/buy-item/{item}", BuyItem).Methods("GET")
	r.HandleFunc("/buy-item-qty/{item}", BuyItemQty).Methods("GET")
	r.HandleFunc("/buy-item-qty-price/{item}", BuyItemQtyPrice).Methods("GET")
	r.HandleFunc("/show-summary", ShowSummary).Methods("GET")
	r.HandleFunc("/fast-buy-item/{item}", FastBuyItem).Methods("GET")

	return r
}

// BuyItem is the handler for /buy-item end-point and returns the item if available
func BuyItem(w http.ResponseWriter, r *http.Request) {
	requestedItem := mux.Vars(r)["item"]
	if requestedItem == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("item not specified")))
		return
	}

	order := item.Order{ItemName: requestedItem}

	availableItem := order.BuyItem()
	if len(availableItem) == 0 {
		w.Write([]byte(fmt.Sprintf("NOT_FOUND")))
		return
	}

	response, err := json.Marshal(availableItem)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func BuyItemQty(w http.ResponseWriter, r *http.Request) {
	requestedItem := mux.Vars(r)["item"]
	var err error
	if requestedItem == "" {
		err = errors.New("item not specified")
	}
	requestedQty := r.URL.Query()["quantity"]
	if len(requestedQty) != 1 {
		err = errors.New("quantity not specified")
	}
	quantity, err := strconv.Atoi(requestedQty[0])
	if  err != nil {
		err = errors.New("incorrect quantity")
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(err.Error())))
		return
	}

	order := item.Order{ItemName: requestedItem, Quantity: &quantity}

	availableItem := order.BuyItem()
	if len(availableItem) == 0 {
		w.Write([]byte(fmt.Sprintf("NOT_FOUND")))
		return
	}

	response, err := json.Marshal(availableItem)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func BuyItemQtyPrice(w http.ResponseWriter, r *http.Request) {
	requestedItem := mux.Vars(r)["item"]
	var err error
	if requestedItem == "" {
		err = errors.New("item not specified")
	}
	requestedQty := r.URL.Query()["quantity"]
	if len(requestedQty) != 1 {
		err = errors.New("quantity not specified")
	}
	quantity, err := strconv.Atoi(requestedQty[0])
	if  err != nil {
		err = errors.New("incorrect quantity")
	}
	requestedPrice := r.URL.Query()["price"]
	if len(requestedPrice) != 1 {
		err = errors.New("price not specified")
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(err.Error())))
		return
	}
	order := item.Order{ItemName: requestedItem, Quantity: &quantity, Price: &requestedPrice[0], UseCache: true}
	availableItem := order.BuyItem()
	if len(availableItem) == 0 {
		w.Write([]byte(fmt.Sprintf("NOT_FOUND")))
		return
	}

	response, err := json.Marshal(availableItem)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func ShowSummary(w http.ResponseWriter, r *http.Request) {
	summary := item.ShowSummary()
	response, err := json.Marshal(summary)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func FastBuyItem(w http.ResponseWriter, r *http.Request) {
	requestedItem := mux.Vars(r)["item"]
	if requestedItem == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("item not specified")))
		return
	}

	order := item.Order{ItemName: requestedItem}

	availableItem := order.FastBuyItem()
	if len(availableItem) == 0 {
		w.Write([]byte(fmt.Sprintf("NOT_FOUND")))
		return
	}

	response, err := json.Marshal(availableItem)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
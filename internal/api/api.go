package api

import (
	"encoding/json"
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

// BuyItem is the handler for /buy-item end-point which returns the item if available
func BuyItem(w http.ResponseWriter, r *http.Request) {
	// Validation
	requestedItem := mux.Vars(r)["item"]
	if requestedItem == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("item not specified")))
		return
	}

	order := item.Order{ItemName: requestedItem}

	availableItem := order.BuyItem()
	// If item is not present, return NOT_FOUND
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

// BuyItemQty is the handler for /buy-item-qty end-point which returns the item if the specified quantity is available
func BuyItemQty(w http.ResponseWriter, r *http.Request) {
	requestedItem := mux.Vars(r)["item"]
	var err error
	if requestedItem == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("item not specified")))
		return
	}
	requestedQty := r.URL.Query()["quantity"]
	if len(requestedQty) != 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("quantity not specified")))
		return
	}
	quantity, err := strconv.Atoi(requestedQty[0])
	if err != nil || quantity < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("incorrect quantity")))
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

// BuyItemQtyPrice is the handler for /buy-item-qty-price end-point which returns the item if the specified quantity is available
// within the specified price
func BuyItemQtyPrice(w http.ResponseWriter, r *http.Request) {
	requestedItem := mux.Vars(r)["item"]
	var err error
	if requestedItem == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("item not specified")))
		return
	}
	requestedQty := r.URL.Query()["quantity"]
	if len(requestedQty) != 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("quantity not specified")))
		return
	}
	quantity, err := strconv.Atoi(requestedQty[0])
	if err != nil || quantity < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("incorrect quantity")))
		return
	}
	requestedPrice := r.URL.Query()["price"]
	if len(requestedPrice) != 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("price not specified")))
		return
	}

	// cache is to be used for buy-item-qty-price end point
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

// ShowSummary is the handler for /show-summary end-point which returns the stock available in cache
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

// FastBuyItem is the handler for /fast-buy-item end-point which returns the item if available by making calls to suppliers parallelly
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

package item

import (
	"encoding/json"
	"fmt"
	"github.com/Abhijith01/foodaggregator/internal/config"
	"strconv"
	"strings"
)

type Item []struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Price    string `json:"price"`
	Quantity int    `json:"quantity"`
}

type Fruit []struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Price    string `json:"price"`
	Quantity int    `json:"quantity"`
}

type Vegetable []struct {
	Id       string `json:"productId"`
	Name     string `json:"productName"`
	Price    string `json:"price"`
	Quantity int    `json:"quantity"`
}

type Grain []struct {
	Id       string `json:"itemId"`
	Name     string `json:"itemName"`
	Price    string `json:"price"`
	Quantity int    `json:"quantity"`
}

type Order struct {
	ItemName string
	Quantity *int
	Price    *string
	UseCache bool
}

type fruitSupplier struct{}
type vegetableSupplier struct{}
type grainSupplier struct{}

type supplier interface {
	getStock() (Item, error)
}

// BuyItem is the function used to check if the requested item is present in the stock of all the 3 suppliers
// quantity and price conditions are also checked if sent
func (o *Order) BuyItem() Item {
	var availableItem Item
	var err error
	// Cache is used only for buy-item-qty-price end-point
	if o.UseCache {
		cacheData := getCacheData()
		for _, eachStoreData := range cacheData {
			availableItem = o.checkAvailability(eachStoreData)
			if len(availableItem) > 0 {
				return availableItem
			}
		}
	}
	availableItem, err = o.checkFruits()
	// If there is an error or if requested item is not present, call next supplier and repeat the process
	if err == nil && len(availableItem) > 0 {
		return availableItem
	}

	availableItem, err = o.checkVegetables()
	if err == nil && len(availableItem) > 0 {
		return availableItem
	}

	availableItem, err = o.checkGrains()
	if err == nil && len(availableItem) > 0 {
		return availableItem
	}
	return availableItem
}

// checkFruits makes an API call to fruit supplier and checks if the requested item is present in the stock or not
func (o *Order) checkFruits() (Item, error) {
	var availableItem Item
	fruitSupplier := fruitSupplier{}
	allFruits, err := fruitSupplier.getStock()
	if err == nil && allFruits != nil {
		// The supplier data is saved in cache only for buy-item-qty-price API
		if o.UseCache {
			go storeInCache("fruitSupplier", allFruits)
		}
		availableItem = o.checkAvailability(allFruits)
		if len(availableItem) > 0 {
			return availableItem, nil
		}
	}
	return availableItem, err
}

// checkVegetables makes an API call to vegetable supplier and checks if the requested item is present in the stock or not
func (o *Order) checkVegetables() (Item, error) {
	var availableItem Item
	vegetableSupplier := vegetableSupplier{}
	allVegetables, err := vegetableSupplier.getStock()
	if err == nil && allVegetables != nil {
		if o.UseCache {
			go storeInCache("vegetableSupplier", allVegetables)
		}
		availableItem = o.checkAvailability(allVegetables)
		if len(availableItem) > 0 {
			return availableItem, nil
		}
	}
	return availableItem, err
}

// checkGrains makes an API call to grain supplier and checks if the requested item is present in the stock or not
func (o *Order) checkGrains() (Item, error) {
	var availableItem Item
	grainSupplier := grainSupplier{}
	allGrains, err := grainSupplier.getStock()
	if err == nil && allGrains != nil {
		if o.UseCache {
			go storeInCache("grainSupplier", allGrains)
		}
		availableItem = o.checkAvailability(allGrains)
		if len(availableItem) > 0 {
			return availableItem, nil
		}
	}
	return availableItem, err
}

// checkAvailability loops over the items passed and checks if the requested item is present in the items passed
func (o *Order) checkAvailability(allItems Item) Item {
	var availableItem Item
	for _, eachItem := range allItems {
		if strings.ToLower(o.ItemName) == strings.ToLower(eachItem.Name) {
			found := true
			// If quantity is passed, check if the item's quantity in stock >= requested quantity
			if o.Quantity != nil {
				if *o.Quantity > eachItem.Quantity {
					found = false
					continue
				}
			}
			// If price is passed, check if the item's price in stock <= requested price
			if o.Price != nil {
				if !validatePrice(*o.Price, eachItem.Price) {
					found = false
					continue
				}
			}
			if found {
				availableItem = append(availableItem, eachItem)
			}
		}
	}
	return availableItem
}

// validatePrice checks if the requestedPrice >= price of available stock.
func validatePrice(requestedPrice, availablePrice string) bool {
	// Remove the dollar symbol from the price string got from supplier data  for comparing
	availablePrice = string([]rune(availablePrice)[1:])

	requested, err := strconv.ParseFloat(requestedPrice, 64)
	if err != nil {
		fmt.Println("Failed to parse float", err)
		return false
	}
	available, err := strconv.ParseFloat(availablePrice, 64)
	if err != nil {
		fmt.Println("Failed to parse float", err)
		return false
	}
	if requested >= available {
		return true
	} else {
		return false
	}
}

// FastBuyItem calls the 3 supplier end-points asynchronously by making use of go routines
// If the requested item is present in any of the supplier data,
// the item is returned from the respective go routine functions by making use of channels
func (o *Order) FastBuyItem() Item {
	// 3 separate channels for each of the suppliers
	itemFruits, itemVegetables, itemGrains := make(chan *Item), make(chan *Item), make(chan *Item)
	go o.checkFruitsForFastBuy(itemFruits)
	go o.checkVegesForFastBuy(itemVegetables)
	go o.checkGrainsForFastBuy(itemGrains)

	// A for loop which iterates 3 times is created since there are 3 suppliers
	// If channel holds non-nil value, return the item then and there without waiting for all the go routines to complete
	for i := 0; i < 3; i++ {
		select {
		case availableItem := <-itemFruits:
			if availableItem != nil {
				return *availableItem
			}
		case availableItem := <-itemVegetables:
			if availableItem != nil {
				return *availableItem
			}
		case availableItem := <-itemGrains:
			if availableItem != nil {
				return *availableItem
			}
		}
	}
	return Item{}
}

// checkFruitsForFastBuy makes a call to get the stock from fruit supplier
// If the requested item is found, data is passed on with the help of channel
func (o *Order) checkFruitsForFastBuy(itemFruits chan *Item) {
	availableItem, err := o.checkFruits()
	if err == nil && len(availableItem) > 0 {
		itemFruits <- &availableItem
	}
	itemFruits <- nil
}

// checkVegesForFastBuy makes a call to get the stock from vegetable supplier
// If the requested item is found, data is passed on with the help of channel
func (o *Order) checkVegesForFastBuy(itemVegetables chan *Item) {
	availableItem, err := o.checkVegetables()
	if err == nil && len(availableItem) > 0 {
		itemVegetables <- &availableItem
	}
	itemVegetables <- nil
}

// checkGrainsForFastBuy makes a call to get the stock from grain supplier
// If the requested item is found, data is passed on with the help of channel
func (o *Order) checkGrainsForFastBuy(itemGrains chan *Item) {
	availableItem, err := o.checkGrains()
	if err == nil && len(availableItem) > 0 {
		itemGrains <- &availableItem
	}
	itemGrains <- nil
}

// fruitSupplier.getStock makes a GET call to the fruit supplier end point and returns the data in Item struct
func (fs fruitSupplier) getStock() (Item, error) {
	var availableFruits Item
	url := config.FruitSupplierUrl
	rawData, err := makeGetCall(url)
	if err != nil {
		return availableFruits, err
	}
	// Unmarshal the response to Fruit struct
	var fruit Fruit
	err = json.Unmarshal(rawData, &fruit)
	if err != nil {
		fmt.Println("Error Unmarshalling Fruit supplier data", err)
		return nil, err
	}
	// Casting fruit struct to Item struct so that checking the data in the parent struct is simplified
	// Casting is achieved since both the structs have the same fields
	availableFruits = Item(fruit)
	return availableFruits, nil
}

// vegetableSupplier.getStock makes a GET call to the vegetable supplier end point and returns the data in Item struct
func (vs vegetableSupplier) getStock() (Item, error) {
	var availableVegetables Item
	url := config.VegetableSupplierUrl
	rawData, err := makeGetCall(url)
	if err != nil {
		return availableVegetables, err
	}
	// Unmarshal the response to Vegetable struct
	var vegetable Vegetable
	err = json.Unmarshal(rawData, &vegetable)
	if err != nil {
		fmt.Println("Error Unmarshalling Vegetable supplier data", err)
		return availableVegetables, err
	}
	// Casting vegetable struct to Item struct so that checking the data in the parent struct is simplified
	// Casting is achieved since both the structs have the same fields
	availableVegetables = Item(vegetable)
	return availableVegetables, nil
}

// grainSupplier.getStock makes a GET call to the grain supplier end point and returns the data in Item struct
func (gs grainSupplier) getStock() (Item, error) {
	var availableGrains Item
	url := config.GrainSupplierUrl
	rawData, err := makeGetCall(url)
	if err != nil {
		return availableGrains, err
	}
	// Unmarshal the response to Grain struct
	var grain Grain
	err = json.Unmarshal(rawData, &grain)
	if err != nil {
		fmt.Println("Error Unmarshalling Grain supplier data", err)
		return nil, err
	}
	// Casting grain struct to Item struct so that checking the data in the parent struct is simplified
	// Casting is achieved since both the structs have the same fields
	availableGrains = Item(grain)
	return availableGrains, nil
}

// ShowSummary returns the stock data that is held in cachedItems
func ShowSummary() map[string]Item {
	return cachedItems
}

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

func (o *Order) BuyItem() Item {
	var availableItem Item
	var err error
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

func (o *Order) checkFruits() (Item, error) {
	var availableItem Item
	fruitSupplier := fruitSupplier{}
	allFruits, err := fruitSupplier.getStock()
	if err == nil && allFruits != nil {
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

func (o *Order) checkAvailability(allItems Item) Item {
	var availableItem Item
	for _, eachItem := range allItems {
		if strings.ToLower(o.ItemName) == strings.ToLower(eachItem.Name) {
			found := true
			if o.Quantity != nil {
				if *o.Quantity > eachItem.Quantity {
					found = false
					continue
				}
			}
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

func validatePrice(requestedPrice, availablePrice string) bool {
	// Remove the dollar symbol from the string for comparing
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

func (o *Order) FastBuyItem() Item {
	itemFruits, itemVegetables, itemGrains := make(chan *Item), make(chan *Item), make(chan *Item)
	go o.checkFruitsForFastBuy(itemFruits)
	go o.checkVegesForFastBuy(itemFruits)
	go o.checkGrainsForFastBuy(itemFruits)

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

func (o *Order) checkFruitsForFastBuy(itemFruits chan *Item) {
	availableItem, err := o.checkFruits()
	if err == nil && len(availableItem) > 0 {
		itemFruits <- &availableItem
	}
	itemFruits <- nil
}

func (o *Order) checkVegesForFastBuy(itemVegetables chan *Item) {
	availableItem, err := o.checkVegetables()
	if err == nil && len(availableItem) > 0 {
		itemVegetables <- &availableItem
	}
	itemVegetables <- nil
}

func (o *Order) checkGrainsForFastBuy(itemGrains chan *Item) {
	availableItem, err := o.checkGrains()
	if err == nil && len(availableItem) > 0 {
		itemGrains <- &availableItem
	}
	itemGrains <- nil
}

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
	availableFruits = Item(fruit)
	return availableFruits, nil
}

func (vs vegetableSupplier) getStock() (Item, error) {
	var availableVegetables Item
	url := config.VegetableSupplierUrl
	rawData, err := makeGetCall(url)
	if err != nil {
		return availableVegetables, err
	}
	// Unmarshal the response to Fruit struct
	var vegetable Vegetable
	err = json.Unmarshal(rawData, &vegetable)
	if err != nil {
		fmt.Println("Error Unmarshalling Vegetable supplier data", err)
		return availableVegetables, err
	}
	availableVegetables = Item(vegetable)
	return availableVegetables, nil
}

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
	availableGrains = Item(grain)
	return availableGrains, nil
}

func ShowSummary() map[string]Item {
	return cachedItems
}

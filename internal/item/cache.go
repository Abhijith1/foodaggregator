package item

// cachedItems is the cache data structure used to store the supplier stock data
var cachedItems = map[string]Item{}

// storeInCache stores the supplier's stock data into cachedItems.
// Old data of the supplier will be replaced
func storeInCache(supplier string, items Item) {
	cachedItems[supplier] = items
}

// getCacheData returns the stock data present in the cache data structure
func getCacheData() map[string]Item {
	return cachedItems
}

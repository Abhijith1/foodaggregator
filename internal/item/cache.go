package item

var cachedItems = map[string]Item{}

func storeInCache(supplier string, items Item) {
	cachedItems[supplier] = items
}

func getCacheData() map[string]Item {
	return cachedItems
}

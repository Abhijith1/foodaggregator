package item

var cachedItems = map[string]Item{}

func storeInCache(supplier string, items Item) {
	cachedItems[supplier] = items
}



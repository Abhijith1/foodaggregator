# foodaggregator

foodaggregator is a simple REST API service which returns available stock of items from external mock APIs

### How to run

Assuming you have Go installed on your system, run the following commands in your terminal:

```bash
$ cd $GOPATH/src
$ mkdir -p github.com/Abhijith01
$ cd github.com/Abhijith01
$ git clone https://github.com/Abhijith1/foodaggregator.git
$ cd foodaggregator
$ go build
$ ./foodaggregator

Started server at port :8080
```

End points exposed:

1.   localhost:8080/buy-item/:item
Eg.: localhost:8080/buy-item/apple

2.   localhost:8080/buy-item-qty/:item
Eg.: localhost:8080/buy-item-qty/apple?quantity=9

3.   localhost:8080/buy-item-qty-price/:item
Eg.: localhost:8080/buy-item-qty-price/apple?quantity=10&price=100

4.   localhost:8080/show-summary

5.   localhost:8080/fast-buy-item/:item
Eg.: localhost:8080/fast-buy-item/apple

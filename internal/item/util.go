package item

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func makeGetCall(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error in making http get call", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error in reading get API response body", err)
		return nil, err
	}
	return body, nil
}

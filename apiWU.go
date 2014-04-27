package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

func main() {
	URL := "http://api.wunderground.com/api/06ae7ac7474e21f1/geolookup/conditions/q/24356.json"
	resp, err := http.Get(URL)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	var dat map[string]interface{}

	if err := json.Unmarshal(b, &dat); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(dat["current_observation"][1])
}

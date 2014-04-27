package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

func main() {
	baseURL := "http://api.nytimes.com/svc/mostpopular/v2/mostviewed/arts/7.json?api-key=98278497df4cbc1264ca529a8f412acf:4:69326077"
	resp, err := http.Get(baseURL)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)

	var dat map[string]interface{}

	if err := json.Unmarshal(b, &dat); err != nil {
		panic(err)
	}

	resultArray := dat["results"]
	
}

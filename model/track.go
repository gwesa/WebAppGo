package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func checkErr(error error){
	if error != nil {
		fmt.Println(error)
	}
}

func IsValueSlice(slice []int, value int) (result bool) {
	for _, number := range slice {
		if number == value {
			return true
		}

	}
	return false

}

type Item struct {
	ItemName string `json:"itemName"`
	ItemCost float64 `json:"itemCost"`
}


type AllItems struct{
	Id int `json:"id"`
	AllItem []*Item 
}

type TotalItems struct {
	TotalItem []*AllItems
}

func ShowTotalItems() (au *TotalItems) {
	file, error := os.OpenFile("spendtrack.json", os.O_RDWR|os.O_APPEND, 0666)
	checkError(error)
	b, err := ioutil.ReadAll(file)
	var totalItems TotalItems
	json.Unmarshal(b, &totalItems.TotalItem)
	checkError(err)
	return &totalItems
}
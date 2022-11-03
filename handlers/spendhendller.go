package handlers

import (
	"WebAppGo/model"
	"encoding/json"
	// "errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

func checkValue(w http.ResponseWriter, r *http.Request, forms ...string) (res bool, errStr string) {
	for _, form := range forms {
		m, _ := regexp.MatchString("^[a-zA-Z]+$", r.FormValue(form))
		if r.FormValue(form) == "" {
			return false, "All forms must be completed"
		}
		if m == false {
			return false, "Use only english letters to fill up the Booking Form"
		}

	}
	return true, ""
}

func ShowSpending(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("templates/showSpendPage.html")
		t.Execute(w, nil)
	} else {
		id, error := strconv.Atoi(r.FormValue("id"))
		checkError(error)
		var allTotal model.TotalItems
		file, error := os.OpenFile("spendtrack.json", os.O_RDONLY, 0666)
		checkError(error)
		data, error := ioutil.ReadAll(file)
		checkError(error)
		json.Unmarshal(data, &allTotal.TotalItem)

		var allID []int
		for _, items := range allTotal.TotalItem {
			allID = append(allID, items.Id)
		}
		for _, items := range allTotal.TotalItem {
			if model.IsValueSlice(allID, id) != true {
				http.Redirect(w, r, "/showuser/notsuccededshow/", 302)
				return
			}
			if items.Id != id {
				continue
			} else {
				t, err := template.ParseFiles("templates/showSpendPage.html")
				checkError(err)
				t.Execute(w, items)
			}

		}

	}
}

func checkErr(error error){
	if error != nil{
		fmt.Println(error)
	}
}


func AddNewSpend(w http.ResponseWriter, r *http.Request){

	newSpend := &model.Item{}
	groupSpend := &model.AllItems{}
	if r.Method != "GET" {
		template, _ := template.ParseFiles("templates/spend_track.html")
		template.Execute(w, nil)
	} else{
		resBool, errStr := checkFormValue(w, r, "itemName")
		if resBool == false {
			t, err := template.ParseFiles("templates/notSucceded.html")
			checkError(err)
			t.Execute(w, errStr)

			return
		}

		newSpend.ItemName = r.FormValue("itemName")
		var error error
		newSpend.ItemCost, error = strconv.ParseFloat(r.FormValue("itemCost"), 64)
		checkErr(error)

		//open file
		file, error := os.OpenFile("spendtrack.json", os.O_RDWR, 0644)
		checkError(error)
		defer file.Close()

		//read file and unmarshall json file to slice of users
		data, error := ioutil.ReadAll(file)
		var AllItems model.TotalItems
		error = json.Unmarshal(data, &AllItems.TotalItem)
		checkError(error)
		max := 0

		//generation of id(last id at the json file+1)
		for _, items := range AllItems.TotalItem {
			if items.Id > max {
				max = items.Id
			}
		}
		id := max + 1
		groupSpend.Id = id

		//appending newUser to slice of all Users and rewrite json file
		AllItems.TotalItem = append(AllItems.TotalItem, groupSpend)
		newUserBytes, err := json.MarshalIndent(&AllItems.TotalItem, "", " ")
		checkError(err)
		ioutil.WriteFile("spendtrack.json", newUserBytes, 0666)
		http.Redirect(w, r, "/", 301)

	}
}
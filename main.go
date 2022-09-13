 package main

import (
    "encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"errors"
	"github.com/gorilla/mux"
)

type shoppingItem struct {
	ID string
	Item string 
	Quantity int
}


func getLocalList() ([]shoppingItem,error){
	content, err := ioutil.ReadFile("./resources/shopList.json")
    if err != nil {
        return nil, errors.New("ShopList not found")
    }
    
    var payload = []shoppingItem{}
    err = json.Unmarshal(content, &payload)
    if err != nil {
        return nil, errors.New("ShopList parse error")
    }

    return payload,nil
}

func setLocalList(list []shoppingItem) (error){

    body, err := json.Marshal(&list)
    if err != nil {
        return errors.New("ShopList WriteFile parse error")
    }

	ioutil.WriteFile("./resources/shopList.json",body,0644)
    return nil
}


func getHttpShoppingList(w http.ResponseWriter, r *http.Request){
	list,error := getLocalList()
	if error != nil {
        fmt.Fprintf(w, "getHttpShoppingList error")
    }

    json.NewEncoder(w).Encode(list)
}

func postHttpShoppingList(w http.ResponseWriter, r *http.Request){
	//get local list...
	list,listError := getLocalList()
	if listError != nil {
		fmt.Fprintf(w, "postHttpShoppingList error")
		return
	}


	//get new item from post body...
	var newItem shoppingItem
	reqBody, bodyErr := ioutil.ReadAll(r.Body)
	if bodyErr != nil {
		fmt.Fprintf(w, "Body format not correct")
		return
	}

	//convert data to shopping item
	jsonErr := json.Unmarshal(reqBody, &newItem)
    if jsonErr != nil {
    	fmt.Fprintf(w, "Body format not correct")
        return 
    }

    //update local list
    list = append(list,newItem)
    err := setLocalList(list)
    if err != nil {
    	fmt.Fprintf(w, "setLocalList error")
        return 
    }

    //return result
    w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newItem)
}


func updateShoppingItem(w http.ResponseWriter, r *http.Request){
	eventID := mux.Vars(r)["id"]

	list,listError := getLocalList()
	if listError != nil {
		fmt.Fprintf(w, "updateShoppingItem error")
		return
	}

	reqBody, bodyErr := ioutil.ReadAll(r.Body)
	if bodyErr != nil {
		fmt.Fprintf(w, "Body format not correct")
		return
	}


	//convert data to shopping item
	var newItem shoppingItem
	jsonErr := json.Unmarshal(reqBody, &newItem)
    if jsonErr != nil {
    	fmt.Fprintf(w, "Body format not correct")
        return 
    }


	for i,item := range list {
		if item.ID == eventID {
			item.Quantity = newItem.Quantity
			item.Item = newItem.Item
			list[i] = newItem
			//list = append(list[:i], item)
		}
	}

	err := setLocalList(list)
    if err != nil {
    	fmt.Fprintf(w, "setLocalList error")
        return 
    }

    fmt.Fprintf(w,"updateShoppingItem SUCCESS")

}

func welcome(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "GOLANG ShopItem Example\n\nGET localhost:8080/shoppingList\nPOST localhost:8080/shoppingList (with body)\nPATCH localhost:8080/shoppingList/id\n")
}

func main(){
    //start router
    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/", welcome).Methods("GET")
    router.HandleFunc("/shoppingList", getHttpShoppingList).Methods("GET")
    router.HandleFunc("/shoppingList", postHttpShoppingList).Methods("POST")
    router.HandleFunc("/shoppingList/{id}", updateShoppingItem).Methods("PATCH")

    fmt.Printf("Starting server at port 8080...\n\n")


	log.Fatal(http.ListenAndServe(":8080", router))
}
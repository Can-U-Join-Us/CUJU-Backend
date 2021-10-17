package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"strconv"

	"github.com/gorilla/mux"
)

const port = ":5500"

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", rootPage)
	router.HandleFunc("/competitions/{Count}", competitions).Methods("GET")
	fmt.Println("Serving @ http://127.0.0.1" + port)
	log.Fatal(http.ListenAndServe(port, router))

}

func rootPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is root page"))
}

func competitions(w http.ResponseWriter, r *http.Request) {
	Count, errInput := strconv.ParseFloat(mux.Vars(r)["Count"], 64)
	log.Println("Competitions Get Request len : " + mux.Vars(r)["Count"])

	fetchCount := 0

	if errInput != nil {
		fmt.Println(errInput.Error())
	} else {
		fetchCount = int(Count)
		if fetchCount > len(competitionList) {
			log.Println("Competition len over")
			fetchCount = len(competitionList)
		}
	}
	fmt.Println(competitionList[0:fetchCount])
	jsonList, err := json.Marshal(competitionList[0:fetchCount])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.Header().Set("content-type", "application/json")
		w.Write(jsonList)
	}
}

type competition struct {
	Title       string
	Description string
	DueDate     string
	ImageUrl    string
}

var competitionList = []competition{
	competition{"p1", "test1", "1998-12-16", "localhost/image"},
	competition{"p2", "test2", "1999-12-16", "localhost/image"},
	competition{"p3", "test3", "1918-03-16", "localhost/image"},
	competition{"p4", "test4", "1928-01-16", "localhost/image"},
	competition{"p5", "test5", "1938-05-16", "localhost/image"},
	competition{"p6", "test6", "1948-02-16", "localhost/image"},
	competition{"p7", "test7", "1958-04-16", "localhost/image"},
	competition{"p8", "test8", "1968-10-16", "localhost/image"},
	competition{"p9", "test9", "1978-11-16", "localhost/image"},
}

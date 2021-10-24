package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"net"

	"github.com/gorilla/mux"
)

const port = ":5500"

// Get preferred outbound ip of this machine
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
func main() {
	log.Println(competitionList[1].Title)

	router := mux.NewRouter()
	router.HandleFunc("/", rootPage)
	router.HandleFunc("/competitions/{Count}", competitions).Methods("GET")
	fmt.Println("Serving @ " + GetLocalIP() + port)

	log.Fatal(http.ListenAndServe(port, router))
}

func rootPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is root page"))
}
func competitions(w http.ResponseWriter, r *http.Request) {
	Count, errInput := strconv.ParseFloat(mux.Vars(r)["Count"], 64)
	if errInput != nil {
		fmt.Println(errInput.Error())
	} else {
		data := competitionList[0:int(Count)]
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(data)

	}

	log.Println("Competitions Get Request len : " + mux.Vars(r)["Count"])

	// fetchCount := 0

	// if errInput != nil {
	// 	fmt.Println(errInput.Error())
	// } else {
	// 	fetchCount = int(Count)
	// 	if fetchCount > len(competitionList) {
	// 		log.Println("Competition len over")
	// 		fetchCount = len(competitionList)
	// 	}
	// }
	// _, err := json.Marshal(competitionList[0:fetchCount])

	// w.Header().Set("content-type", "application/json")

	// json.NewEncoder(w).Encode(cps)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// } else {
	// 	w.Header().Set("content-type", "application/json")

	// }

}

type competition struct {
	Title       string `json : "title"`
	Description string `json : "description"`
	DueDate     string `json : "dueDate"`
	ImageUrl    string `json : "imageUrl"`
}

var competitionList = []competition{
	competition{Title: "p1", Description: "test1", DueDate: "1998-12-16", ImageUrl: "localhost/image"},
	competition{Title: "p2", Description: "test2", DueDate: "1999-12-16", ImageUrl: "localhost/image"},
	competition{Title: "p3", Description: "test3", DueDate: "1918-03-16", ImageUrl: "localhost/image"},
	competition{Title: "p4", Description: "test4", DueDate: "1928-01-16", ImageUrl: "localhost/image"},
	competition{Title: "p5", Description: "test5", DueDate: "1938-05-16", ImageUrl: "localhost/image"},
	competition{Title: "p6", Description: "test6", DueDate: "1948-02-16", ImageUrl: "localhost/image"},
	competition{Title: "p7", Description: "test7", DueDate: "1958-04-16", ImageUrl: "localhost/image"},
	competition{Title: "p8", Description: "test8", DueDate: "1968-10-16", ImageUrl: "localhost/image"},
	competition{Title: "p9", Description: "test9", DueDate: "1978-11-16", ImageUrl: "localhost/image"},
}

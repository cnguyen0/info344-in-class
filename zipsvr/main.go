package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/cnguyen0/info344-in-class/handlers"
	"github.com/cnguyen0/info344-in-class/zipsvr/models"
)

const zipPath = "/zips/"

func helloHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	w.Header().Add("Content-Type", "text/plain")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, "Hello %s!", name)
	//w.Write([]byte("Hello, World!"))
}

func memoryHandler(w http.ResponseWriter, r *http.Request) {
	runtime.GC()
	stats := &runtime.MemStats{}
	runtime.ReadMemStats(stats)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func main() {
	// fmt.Println("Hello, World")
	addr := os.Getenv("ADDR") // convention, its in all capitalized letters)
	if len(addr) == 0 {
		addr = ":443"
	}

	tlskey := os.Getenv("TLSKEY")
	tlscert := os.Getenv("TLSCERT")
	if len(tlskey) == 0 || len(tlscert) == 0 {
		log.Fatal("please set key and cert")
	}

	zips, err := models.LoadZips("zips.csv")
	if err != nil {
		log.Fatalf("error loading zips: %v", err)
	}

	log.Printf("loaded %d zips", len(zips))

	cityIndex := models.ZipIndex{}

	for _, z := range zips {
		// use .ToLower() to lowercase the name
		cityLower := strings.ToLower(z.City)

		// the city index is a zip index. what u get back is a slice of zip
		cityIndex[cityLower] = append(cityIndex[cityLower], z)
	}

	// mux is able to handle multiple resource path
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/memory", memoryHandler)

	cityHandler := &handlers.CityHandler{
		Index:      cityIndex,
		PathPrefix: zipPath, // when you declare a struct u need a comma in every single line
	}

	// .handle vs .handlefunc
	mux.Handle(zipPath, cityHandler)

	fmt.Printf("server is listening at https://%s\n", addr)

	// mux parameter is the traffic handler
	log.Fatal(http.ListenAndServeTLS(addr, tlscert, tlskey, mux))

}

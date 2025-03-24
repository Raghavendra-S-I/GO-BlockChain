package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Block struct {
	Id     int
	Title  string
	Author string
}

type BookCheckout struct {
}

type Book struct {
}

type BlockChain struct {
	block []*Block
}

var BlockChain *BlockChain

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", getBlockChain).Methods("GET")
	r.HandleFunc("/", writeblock).Methods("POST")
	r.HandleFunc("/new", newBook).Methods("POST")

	log.Println("Listening on port 3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}

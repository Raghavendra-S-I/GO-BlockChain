package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Block struct {
	Pos       int          `json:"pos"`
	Data      BookCheckout `json:"data"`
	TimeStamp string       `json:"timestamp"`
	Hash      string       `json:"hash"`
	PrevHash  string       `json:"prevhash"`
}

type BookCheckout struct {
	BookID       string `json:"bookid"`
	User         string `json:"user"`
	CheckoutDate string `json:"checkoutdata"`
	IsGenesis    bool   `json:"isgenesis"`
}

type Book struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	PublishDate string `json:"publishdate"`
	ISBN        string `json:"isbn"`
}

type BlockChain struct {
	block []*Block
}

var BlockChain *BlockChain

func (b *Block) generateHash() {

	bytes, _ := json.Marshal(b.Data)

	data := string(b.Pos) + b.TimeStamp + string(bytes) + b.PrevHash

	hash := sha256.New()
	hash.Write([]byte(data))
	b.Hash = hex.EncodeToString(hash.Sum(nil))
}

func CreateBlock(prevBlock *Block, data BlockChain) *Block {
	block := &Block{}
	block.Pos = prevBlock.Pos + 1
	block.PrevHash = prevBlock.Hash
	block.generateHash()

	return block
}

func (bc *BlockChain) AddBlock(data BookCheckout) {
	prevBlock := bc.blocks[len(bc.blocks)-1]

	block := CreateBlock(prevBlock, data)

	if validBlock(block, prevBlock) {
		bc.blocks = append(bc.blocks, block)
	}
}

func newBook(w http.ResponseWriter, r *http.Request) {
	var book Book

	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("could not create book : %v", err)
		w.Write([]byte("could not create book"))
		return
	}

	h := md5.New()
	io.WriteString(h, book.ISBN+book.PublishDate)
	book.ID = fmt.Sprintf("%x", h.Sum(nil))

	resp, err := json.MarshalIndent(book, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("could not create book : %v", err)
		w.Write([]byte("Could not create book"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)

}

func validBlock(block, prevBlock *Block) bool {
	if prevBlock.Hash != block.PrevHash {
		return false
	}
	if !block.validateHash(block.Hash) {
		return false

	}

	if prevBlock.Pos+1 != block.Pos {
		return false
	}

	return true
}

func (b *Block) validateHash(hash string) bool {
	b.generateHash()
	if b.Hash != hash {
		return false
	}
	return true
}

func writeblock(w http.ResponseWriter, r *http.Request) {
	var checkoutitem BookCheckout

	if err := json.NewDecoder(r.Body).Decode(&checkoutitem); err != nil {
		r.WriteHeader(http.StatusInternalServerError)
		log.Printf("could not create block : %v", err)
		w.Write([]byte("could not create block"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp)

	BlockChain.AddBlock(checkoutitem)
}

func GenesisBlock() *BlockChain {
	return CreateBlock(&Block{}, BookCheckout{IsGenesis: true})
}

func NewBlockChain() *BlockChain {
	retrun & BlockChain{[]*BlockChain{GenesisBlock()}}
}

func main() {

	BlockChain = NewBlockChain()

	r := mux.NewRouter()
	r.HandleFunc("/", getBlockchain).Methods("GET")
	r.HandleFunc("/", writeblock).Methods("POST")
	r.HandleFunc("/new", newBook).Methods("POST")

	go func() {

		for _, block := range BlockChain.blocks {
			fmt.Printf("PrevHash: %s\n", block.PrevHash)
			json.MarshalIndent(block.Data, "", " ")
			fmt.Printf("Hash: %s\n", block.Hash)
			fmt.Printf("Data: %s\n", string(bytes))

		}
	}()

	log.Println("Listening on port 3000")
	log.Fatal(http.ListenAndServe(":3000", r))

}

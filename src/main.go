package main

import (
	"net/http"
	"fmt"
	"log"
)

func print_binary(s []byte) {
	fmt.Printf("Received b:");
	for n := 0;n < len(s);n++ {
		fmt.Printf("%d,",s[n]);
	}
	fmt.Printf("\n");
}



func main() {
	LoggerInit()

	HubsInitialize()

	router := NewRouter()
	log.Fatal(http.ListenAndServe("127.0.0.1:8000", router))
}
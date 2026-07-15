package main

import (
	"fmt"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Servidor Ativo")
}

func main() {
	mux := http.NewServeMux()
	
	mux.HandleFunc("GET /", hello)

	println("Servidor rodando em http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
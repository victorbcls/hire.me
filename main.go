package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/victorbcls/hire.me/controllers/db"
	"github.com/victorbcls/hire.me/controllers/shorter"
	"go.mongodb.org/mongo-driver/mongo"
)

var _client *mongo.Client

func main() {
	client, err := db.Connect()
	if err != nil {
		log.Fatalf("Erro na conexão com o banco de dados: %v", err)
	}
	_client = client

	router := mux.NewRouter()
	router.Use(corsMiddleware)
	router.HandleFunc("/create", createHandler).Methods(http.MethodPut)
	router.HandleFunc("/{alias}", retrieveHandler).Methods(http.MethodGet)

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalf("Erro ao iniciar o servidor HTTP: %v", err)
	}
}

// createHandler cria uma nova URL encurtada a partir de uma URL original.
func createHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	url := queryParams.Get("url")
	customAlias := queryParams.Get("CUSTOM_ALIAS")

	responseJson, err := shorter.GenerateShortURL(url, customAlias, _client)
	if err != nil {
		http.Error(w, toJSONError(err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJson)
}

// retrieveHandler redireciona o usuário para a URL original associada a um alias.
func retrieveHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	alias, ok := vars["alias"]
	if !ok {
		http.Error(w, toJSONError(errors.New("Alias não encontrado")), http.StatusBadRequest)
		return
	}

	responseJson, err := shorter.RetrieveUrl(alias, _client)
	if err != nil {
		http.Error(w, toJSONError(err), http.StatusInternalServerError)
		return
	}
	var m map[string]interface{}
	json.Unmarshal([]byte(responseJson), &m)
	if m["url"] != nil {
		http.Redirect(w, r, m["url"].(string), http.StatusSeeOther)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(responseJson)
	}

}

// toJSONError converte um erro em uma string JSON no formato {"error": "mensagem de erro"}.
func toJSONError(err error) string {
	jsonErr, _ := json.Marshal(struct{ Error string }{Error: err.Error()})
	return string(jsonErr)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}

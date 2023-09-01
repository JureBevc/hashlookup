package api

import (
	"encoding/json"
	"fmt"
	"hashlookup/util"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Service is running\n")
}

func HandleLookupHash(w http.ResponseWriter, r *http.Request) {
	algorithmName := r.URL.Query().Get("algorithm")
	hash := r.URL.Query().Get("hash")

	log.Printf("Hash lookup request %s %s\n", algorithmName, hash)
	if algorithmName == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid algorithm"))
		return
	}

	if hash == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid hash"))
		return
	}

	message, err := util.GetMessageByHash(algorithmName, hash)

	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": nil,
		})
	} else {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": message,
		})
	}

}

func HandleLookupMessage(w http.ResponseWriter, r *http.Request) {
	algorithmName := r.URL.Query().Get("algorithm")
	message := r.URL.Query().Get("message")

	log.Printf("Message lookup request %s %s\n", algorithmName, message)
	if algorithmName == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid algorithm"))
		return
	}

	if message == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid message"))
		return
	}

	hash, err := util.GetHashByMessage(algorithmName, message)

	if err != nil {
		log.Println(err)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"hash": nil,
		})
	} else {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"hash": hash,
		})
	}
}

func StartAPI(port int) {
	log.Println("Running API...")
	godotenv.Load()
	router := mux.NewRouter()

	router.HandleFunc("/lookup/hash", HandleLookupHash).Methods("GET")
	router.HandleFunc("/lookup/message", HandleLookupMessage).Methods("GET")
	router.HandleFunc("/", HandleHome).Methods("GET")

	log.Printf("Listening on port %d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}

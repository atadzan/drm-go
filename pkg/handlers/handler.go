package handlers

import (
	"encoding/hex"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/md4"
)

type Handler interface {
	CasServerHandler(w http.ResponseWriter, r *http.Request)
	UserKeyHandler(w http.ResponseWriter, r *http.Request)
}

type handler struct {
}

func New() Handler {
	return &handler{}
}

func (h *handler) CasServerHandler(w http.ResponseWriter, r *http.Request) {
	resource := r.URL.Query().Get("file")
	number := r.URL.Query().Get("number")

	if resource == "" || number == "" {
		http.Error(w, "Missing file or number parameter", http.StatusBadRequest)
		return
	}

	log.Printf("Server is requesting key %s for %s from %s", number, resource, r.RemoteAddr)

	// Generate key using MD4 hash
	input := fmt.Sprintf("%s.%s", resource, number)
	hash := md4.New()
	hash.Write([]byte(input))
	key := hex.EncodeToString(hash.Sum(nil))

	// Respond with X-Key-Url header pointing to user-key server
	keyURL := fmt.Sprintf("http://%s/user-key?name=%s&number=%s", r.Host, resource, number)
	w.Header().Set("X-Key-Url", keyURL)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(key)))
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(key))
	return

}

func (h *handler) UserKeyHandler(w http.ResponseWriter, r *http.Request) {
	resource := r.URL.Query().Get("name")
	number := r.URL.Query().Get("number")

	if resource == "" || number == "" {
		http.Error(w, "Missing name or number parameter", http.StatusBadRequest)
		return
	}

	log.Printf("User is requesting key %s for %s from %s", number, resource, r.RemoteAddr)

	// Generate key using MD4 hash
	input := fmt.Sprintf("%s.%s", resource, number)
	hash := md4.New()
	hash.Write([]byte(input))
	key := hash.Sum(nil)

	// Set headers and return key
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(key)))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(key)
}

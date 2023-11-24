package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"product-cmp-api/storage"
	"product-cmp-api/types"
)

type APIServer struct {
	listenAddr string
	store      storage.Storage
}

func NewAPIServer(listenAddr string, store storage.Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	//todo
	router.HandleFunc("/product", makeHTTPHandleFunc(s.handleCreateProduct))
	router.HandleFunc("/product/{brand}/{model}", makeHTTPHandleFunc(s.handleGetProduct))
	//router.HandleFunc("/product/{brand}/{model}", makeHTTPHandleFunc(s.handleDeleteProduct))

	log.Println("JSON API server running on port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func getModelBrand(r *http.Request) (string, string, error) {
	model := mux.Vars(r)["model"]
	brand := mux.Vars(r)["brand"]
	if model == "" || brand == "" {
		return "", "", fmt.Errorf("invalid model/brankd")
	}
	return brand, model, nil
}

func (s *APIServer) handleGetProduct(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodDelete {
		return s.handleDeleteProduct(w, r)
	}
	if r.Method != http.MethodGet {
		return fmt.Errorf("invalid method")
	}
	brand, model, err := getModelBrand(r)
	if err != nil {
		return err
	}
	product, err := s.store.GetProduct(brand, model)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, product)
}

func (s *APIServer) handleCreateProduct(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return fmt.Errorf("invalid method")
	}
	prod := new(types.Product)
	if err := json.NewDecoder(r.Body).Decode(prod); err != nil {
		return err
	}

	if err := s.store.CreateProduct(prod); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, prod)
}
func (s *APIServer) handleDeleteProduct(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodDelete {
		return fmt.Errorf("invalid delete method")
	}
	brand, model, err := getModelBrand(r)
	if err != nil {
		return err
	}
	err = s.store.DeleteProduct(brand, model)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, "deleted successfully")

}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

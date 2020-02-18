package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nikvas0/dc-homework/objects"
	"github.com/nikvas0/dc-homework/storage"
)

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Create product request error: Error while reading request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	product := objects.Product{}
	err = json.Unmarshal(reqBody, &product)
	if err != nil {
		log.Println("Create product request error: Got broken JSON.")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Broken JSON.")
		return
	}

	err = storage.CreateProduct(&product)
	if err != nil {
		log.Println("Create product request error: Database error.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json;")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(product)
	if err != nil {
		log.Println("Create product request error: Encoded broken JSON.")
		return
	}
	log.Println("Create product request: success.")
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		log.Println("Delete product request error: id must be an integer.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = storage.DeleteProductById(uint32(id))
	if storage.IsNotFoundError(err) {
		log.Printf("Get product request: not found (id=%d).", id)
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		log.Println("Delete product request error: Database error.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Println("Delete product request: success.")
}

func GetProducts(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 32)
	if err != nil {
		log.Println("Get product request error: Failed to parse id.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	product := objects.Product{}
	err = storage.GetProductById(&product, uint32(id))
	if storage.IsNotFoundError(err) {
		log.Printf("Get product request: not found (id=%d).", id)
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		log.Println("Get product request error: Failed to parse id.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json;")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(product)
	if err != nil {
		log.Println("Get product request error: Encoded broken JSON.")
		return
	}
	log.Println("Get product request: success.")
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Update product request error: Error while reading request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	product := objects.Product{}
	err = json.Unmarshal(reqBody, &product)
	if err != nil {
		log.Println("Update product request error: Got broken JSON.")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Broken JSON.")
		return
	}

	err = storage.UpdateProduct(&product)
	if storage.IsNotFoundError(err) {
		log.Printf("Get product request: not found (id=%d).", product.ID)
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		log.Println("Update product request error: Database error.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json;")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(product)
	if err != nil {
		log.Println("Update product request error: Encoded broken JSON.")
		return
	}
	log.Println("Update product request: success.")
}

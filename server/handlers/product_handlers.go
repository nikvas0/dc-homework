package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"server/objects"
	"server/storage"

	"github.com/gorilla/mux"
)

const productsPageLimit = 100

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value("role").(uint32) != 1 {
		w.WriteHeader(http.StatusForbidden)
		return
	}

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
		return
	}
	err = objects.FixProduct(&product)
	if err != nil {
		log.Println("Create product request error: Got broken product.", err)
		w.WriteHeader(http.StatusBadRequest)
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
	log.Printf("Create product request: success (id=%d).", product.ID)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value("role").(uint32) != 1 {
		w.WriteHeader(http.StatusForbidden)
		return
	}

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
	log.Printf("Delete product request: success (id=%d).", id)
}

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	var products []objects.Product
	err := storage.GetAllProducts(&products)
	if err != nil && !storage.IsNotFoundError(err) {
		log.Println("Get product request error: Failed to parse id.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json;")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(products)
	if err != nil {
		log.Println("Get products request error: Encoded broken JSON.")
		return
	}

	log.Printf("Get products request: success (got=%d).", len(products))
}

func GetProductsPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var offset uint32
	var limit uint32

	parseUint := func(value *uint32, name string) error {
		value_parsed, err := strconv.ParseUint(vars[name], 10, 32)
		if err != nil {
			log.Println("Get products request error: %s must be an integer.", name)
			return err
		}
		*value = uint32(value_parsed)
		return nil
	}
	if parseUint(&offset, "offset") != nil || parseUint(&limit, "limit") != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if limit > productsPageLimit {
		limit = productsPageLimit
	}

	var products []objects.Product
	err := storage.GetProductsPage(&products, offset, limit)
	if err != nil && !storage.IsNotFoundError(err) {
		log.Println("Get product request error: Failed to parse id.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	total, err := storage.GetProductsCount()
	if err != nil {
		log.Println("Get product request error: Failed to get total.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json;")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"total": total,
		"items": products,
	})
	if err != nil {
		log.Println("Get products request error: Encoded broken JSON.")
		return
	}

	log.Printf("Get products request: success (offset=%d, limit=%d, got=%d).", offset, limit, len(products))
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
	log.Printf("Get product request: success (id=%d).", id)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value("role").(uint32) != 1 {
		w.WriteHeader(http.StatusForbidden)
		return
	}

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
		return
	}
	err = objects.FixProduct(&product)
	if err != nil {
		log.Println("Update product request error: Got broken product.", err)
		w.WriteHeader(http.StatusBadRequest)
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
	log.Printf("Update product request: success (id=%d).", product.ID)
}

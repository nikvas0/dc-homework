package handlers

import (
	"bufio"
	"io"
	"log"
	"net/http"

	"product_upload/objects"
	"product_upload/queues"
	"product_upload/reader"
)

func Upload(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value("role").(uint32) != 1 {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	r.ParseMultipartForm(10 << 30)
	file, header, err := r.FormFile("products_file")
	if err != nil {
		log.Printf("bad file: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("Upload file: name=%s, size=%d", header.Filename, header.Size)

	defer file.Close()
	rr := bufio.NewReader(file)

	var products []objects.Product
	for {
		line, err := rr.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("bad line: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		product, err := reader.ReadProduct(line)
		if err != nil {
			log.Printf("bad line: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = objects.FixProduct(&product)
		if err != nil {
			log.Printf("bad product: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		products = append(products, product)

		if len(products) == 100 {
			queues.SheduleProductBatch(products)
			products = nil
		}
	}
	if len(products) != 0 {
		queues.SheduleProductBatch(products)
		products = nil
	}
	log.Println("file upload: ok")

	w.WriteHeader(http.StatusOK)
}

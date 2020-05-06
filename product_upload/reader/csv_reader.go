package reader

import (
	"errors"
	"strconv"
	"strings"

	"product_upload/objects"
)

func ReadProduct(s string) (objects.Product, error) {
	spl := strings.Split(s, ",")
	if len(spl) != 3 {
		return objects.Product{}, errors.New("wrong row splt=" + strconv.Itoa(len(spl)))
	}
	res := objects.Product{}
	var err error
	ans, err := strconv.ParseUint(strings.Trim(spl[0], " \t\n"), 10, 64)
	if err != nil {
		return res, err
	}
	res.ID = uint32(ans)
	res.Name = strings.Trim(spl[1], " \t\n")
	ans, err = strconv.ParseUint(strings.Trim(spl[2], " \t\n"), 10, 64)
	if err != nil {
		return res, err
	}
	res.Category = uint32(ans)
	return res, nil
}

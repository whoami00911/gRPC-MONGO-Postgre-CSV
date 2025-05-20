package domain

import "errors"

var ErrProductExists = errors.New("productExists")
var ErrProductNotFound = errors.New("productNotFound")

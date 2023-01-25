package handlers

import (
	"context"
	"net/http"

	"github.com/maan19/go-coffeshop-microservices/product-api/data"
)

func (p *Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := &data.Product{}
		err := data.FromJSON(prod, r.Body)
		if err != nil {
			p.l.Println("[ERROR] deserializing product", err)

			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericError{Message: err.Error()}, rw)
			return
		}

		//validate product
		errs := p.v.Validate(prod)
		if len(errs) != 0 {
			p.l.Println("[ERROR] validating product", errs)
		}

		//set product in request context
		ctx := context.WithValue(r.Context(), KeyProduct, prod)
		req := r.WithContext(ctx)
		next.ServeHTTP(rw, req)
	})
}

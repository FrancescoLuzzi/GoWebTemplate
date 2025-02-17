package utils

import (
	"encoding/json"
	"net/http"

	"github.com/FrancescoLuzzi/GoWebTemplate/app/interfaces"
)

func ParseUrlEncoded[T any](r *http.Request, decoder interfaces.Decoder, validator interfaces.Validator) (T, error) {
	var result T
	err := r.ParseForm()
	if err != nil {
		return result, err
	}
	if err = decoder.Decode(&result, r.PostForm); err != nil {
		return result, err
	}
	if err = validator.StructCtx(r.Context(), &result); err != nil {
		return result, err
	}
	return result, err
}

func ParseJson[T any](r *http.Request, validator interfaces.Validator) (T, error) {
	var result T
	err := r.ParseForm()
	if err != nil {
		return result, err
	}
	err = json.NewDecoder(r.Body).Decode(&result)
	if err = validator.StructCtx(r.Context(), &result); err != nil {
		return result, err
	}
	return result, err
}

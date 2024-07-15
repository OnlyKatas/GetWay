package main

import (
	"getway/getway/internal/models"
)

//go:generate swagger generate spec -o ./swagger.json --scan-models

// swagger:meta
// Package classification GeocodService.
//
// Documentation of your project API.
//
//	Schemes:
//	- http
//	- https
//	BasePath: /
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//	- multipart/form-data
//
//	Produces:
//	- application/json
//
//	Security:
//	- basic
//
//	SecurityDefinitions:
//	  Bearer:
//	    type: apiKey
//	    name: Authorization
//	    in: header

// swagger:route POST /api/register address registerRequest
// Регистрация нового пользователя.
// responses:
//   201: registerResponse

// swagger:parameters registerRequest
type registerRequest struct {
	// in:body
	Body models.User
}

// swagger:response registerResponse
type registerResponse struct {
	Str string `json:"str"`
}

// swagger:route POST /api/login address loginRequest
// Авторизация пользователя.
// responses:
//   200: loginResponse

// swagger:parameters loginRequest
type loginRequest struct {
	// in:body
	Body models.User
}

// swagger:response loginResponse
type loginResponse struct {
	// in:body
	Body models.Token
}

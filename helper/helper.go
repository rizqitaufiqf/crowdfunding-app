package helper

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"math/rand"
	"strings"
	"time"
)

type Response struct {
	Meta Meta        `json:"meta"` // when "Meta" is used/converted to json, use "meta" instead of "Meta" as object key
	Data interface{} `json:"data"`
}

type Meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

func APIResponse(message string, code int, status string, data interface{}) Response {
	meta := Meta{
		Message: message,
		Code:    code,
		Status:  status,
	}

	jsonResponse := Response{
		Meta: meta,
		Data: data,
	}

	return jsonResponse
}

func FormatValidationError(err error) []string {
	// store errors to slice
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		var errorsSlice []string
		for _, e := range ve {
			errorsSlice = append(errorsSlice, e.Error())
		}
		return errorsSlice
	}

	return []string{err.Error()}
}

func SanitizePerksSplitString(input string) string {
	parts := strings.FieldsFunc(input, func(r rune) bool {
		return r == ';'
	})
	return strings.Join(parts, ";")
}

func GenerateTransactionCode() string {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	currentDate := time.Now().Format("200601021504")
	randomNumber := rng.Intn(1000000)
	transactionCode := fmt.Sprintf("TRX%s%06d", currentDate, randomNumber)

	return transactionCode
}

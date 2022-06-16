package response

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/google/jsonapi"
)

func ErrorResponse(w http.ResponseWriter, status int, e error) {
	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.Header().Set("Cache-Control", "private, max-age=0")
	w.WriteHeader(status)
	debugging := true
	var eMessage ErrResponse
	//only show real errors when debugging
	if debugging {
		eMessage = ErrResponse{
			Message:    e.Error(),
			StatusCode: status,
		}
	}
	eMessage.Message = friendlyError(eMessage.Message)
	if err := json.NewEncoder(w).Encode(eMessage); err != nil {
		log.Println(err)
	}
}

func Standard(w http.ResponseWriter, statusCode int, response *Response) {
	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.Header().Set("Cache-Control", "private, max-age=0")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println(err)
	}
}

func Action(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.Header().Set("Cache-Control", "private, max-age=0")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(message); err != nil {
		log.Println(err)
	}
}

func friendlyError(message string) string {
	duplicateConstraint := "pq: duplicate key value violates unique constraint"
	if strings.Contains(message, duplicateConstraint) {
		return "Duplicate entry " + strings.Replace(message, duplicateConstraint, "", -1)
	}
	return message
}

type Response struct {
	Data interface{} `json:"data"`
	Meta MetaData    `json:"meta"`
}

type MetaData struct {
	Pagination Pagination     `json:"pagination"`
	Sort       Sort           `json:"sort"`
	Response   ResponseAction `json:"response"`
}

type ResponseAction struct {
	Message string `json:"message"`
}
type Pagination struct {
	PageOffset  int `json:"pageOffset"`
	PageSize    int `json:"pageSize"`
	ResultTotal int `json:"resultTotal"`
	Total       int `json:"total"`
}
type Sort struct {
	Sort string `json:"sort"`
}

type ErrResponse struct {
	StatusCode int    `json:"Code,omitempty"`
	Title      string `json:"Title,omitempty"`
	Message    string `json:"Message,omitempty"`
}

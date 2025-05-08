package mainhandler

import (
	"encoding/json"
	"net/http"

	msgresponse "github.com/uVazzi/go-otus/hw12_13_14_15_calendar/internal/infra/http/main/dto/response"
)

func MainHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r
		response := msgresponse.MessageResponse{
			Message: "Hello World!",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
}

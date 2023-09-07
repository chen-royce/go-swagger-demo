package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"strings"
	"unicode"
)

// A 500 error returned in the response
type errorResponse500 struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
}

// EchoHandlerResponse is the response from the echo handler endpoint
type EchoHandlerResponse struct {
	ResponseString string
}

type echoCaseType string

const (
	upperCase echoCaseType = "upper"
	lowerCase echoCaseType = "lower"
	spongeBob echoCaseType = "spongebob"
)

// EchoHandler echoes and formats text content based on provided query parameters
func EchoHandler(w http.ResponseWriter, r *http.Request) {
	// get input from query param
	data := r.URL.Query().Get("data")

	if len(data) == 0 {
		resp, _ := json.Marshal(errorResponse500{
			Error:  "bad input",
			Status: 500,
		})
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	}

	// get formatting from query param and then format the input
	letterCase := strings.ToLower(r.URL.Query().Get("case"))
	switch echoCaseType(letterCase) {
	case upperCase:
		// POTATO
		data = strings.ToUpper(data)
	case lowerCase:
		// potato
		data = strings.ToLower(data)
	case spongeBob:
		// PoTaTo
		tmp := make([]byte, len(data))
		for idx, char := range data {
			if idx%2 == 0 {
				tmp[idx] = byte(unicode.ToUpper(char))
			} else {
				tmp[idx] = byte(unicode.ToLower(char))
			}
		}
		data = string(tmp)
	}

	// get # of times to repeat/echo the input
	repetitions := r.URL.Query().Get("repetitions")
	numRepetitions, err := strconv.Atoi(repetitions)
	if err != nil {
		numRepetitions = 1
	}
	responseString := data
	for i := 0; i < numRepetitions-1; i++ {
		responseString += " " + data
	}

	// return response as JSON
	resp, _ := json.Marshal(EchoHandlerResponse{ResponseString: responseString})
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

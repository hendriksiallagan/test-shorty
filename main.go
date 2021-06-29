package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"

	"github.com/thanhpk/randstr"
)

// define validate for validation req body
var validate *validator.Validate

// struct to req body
type ShortenParam struct {
	URL       string `json:"url" validate:"required"`
	Shortcode string `json:"shortcode" validate:"omitempty,alphanum,len=6"`
}

// struct for error response
type ErrorResponse struct {
	Message string `json:"message"`
}

// struct for add shortcode response
type ShortenAddResponse struct {
	Shortcode string `json:"shortcode"`
}

// struct for get shorten url location response
type GetShortenResponse struct {
	Location string `json:"location"`
}

// struct for data saved to memory
type DataSaved struct {
	Shortcode     string
	URL           string
	StartDate     time.Time
	LastSeenDate  time.Time
	RedirectCount int
}

// struct for get stats shorten data
type StatsData struct {
	StartDate     time.Time `json:"startDate" `
	LastSeenDate  time.Time `json:"lastSeenDate"`
	RedirectCount int       `json:"redirectCount"`
}

// declare mapping for saved data
var ShortenData map[string][]*DataSaved

func main() {
	// route
	r := mux.NewRouter()
	r.HandleFunc("/shorten", createShorten).Methods(http.MethodPost)
	r.HandleFunc("/shortcode/{shorten}", getDataByShortCode).Methods(http.MethodGet)
	r.HandleFunc("/shortcode/stats/{shorten}", getStatsByShortcode).Methods(http.MethodGet)

	ShortenData = make(map[string][]*DataSaved)

	port := ":7777"
	fmt.Println("Service run on port " + port)
	log.Fatal(http.ListenAndServe(port, r))

}

/**
* Func for generate random shortcode
 */
func randomStr(n int, ShortenData map[string][]*DataSaved) string {
	result := randstr.String(n)
	if _, ok := ShortenData[result]; ok {
		result = randomStr(n, ShortenData)
	}

	return result
}

/**
* Func for response func
 */
func Response(w http.ResponseWriter, httpStatus int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(data)
}

/**
* Func for create shorten code
 */
func createShorten(w http.ResponseWriter, r *http.Request) {
	var params ShortenParam

	// decode req body
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		Response(w, http.StatusInternalServerError, &ErrorResponse{
			Message: "Error Decoded Request",
		})
		return
	}

	// validate req body
	validate = validator.New()
	err = validate.Struct(params)
	if err != nil {
		Response(w, http.StatusBadRequest, &ErrorResponse{
			Message: "Error Validation: " + err.Error(),
		})
		return
	}

	// check whether shorten code exist or not
	if _, ok := ShortenData[params.Shortcode]; ok {
		Response(w, http.StatusConflict, &ErrorResponse{
			Message: "Shortcode already exist",
		})
		return
	}

	// validation if req body of shotcode is empty
	if params.Shortcode == "" {
		params.Shortcode = randomStr(6, ShortenData)
	}

	// Save data to map
	var shortenData DataSaved
	shortenData.Shortcode = params.Shortcode
	shortenData.URL = params.URL
	shortenData.StartDate = time.Now()

	ShortenData[params.Shortcode] = append(ShortenData[params.Shortcode], &shortenData)

	// success response
	Response(w, http.StatusCreated, &ShortenAddResponse{
		Shortcode: params.Shortcode,
	})
	return
}

/**
* Func for get data by shorten code
 */
func getDataByShortCode(w http.ResponseWriter, r *http.Request) {
	shorten := mux.Vars(r)["shorten"]

	// check if shortencode is empty or not found
	_, ok := ShortenData[shorten]
	if shorten == "" || !ok {
		Response(w, http.StatusNotFound, &ErrorResponse{
			Message: "Shorten Not Found",
		})
		return
	}

	// add count data if func called
	var url string
	for _, val := range ShortenData[shorten] {
		val.RedirectCount += 1
		val.LastSeenDate = time.Now()

		url = val.URL
	}

	// success response
	Response(w, http.StatusFound, &GetShortenResponse{
		Location: url,
	})
	return

}

/**
* Func for get stats data by shorten code
 */
func getStatsByShortcode(w http.ResponseWriter, r *http.Request) {
	shorten := mux.Vars(r)["shorten"]

	// check if shortencode is empty or not found
	_, ok := ShortenData[shorten]
	if shorten == "" || !ok {
		Response(w, http.StatusNotFound, &ErrorResponse{
			Message: "Not Found",
		})
		return
	}

	// get stats data
	var data StatsData
	for _, val := range ShortenData[shorten] {
		data = StatsData{
			StartDate:     val.StartDate,
			LastSeenDate:  val.LastSeenDate,
			RedirectCount: val.RedirectCount,
		}
	}

	// success response
	Response(w, http.StatusOK, data)
	return

}

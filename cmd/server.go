package cmd

import (
	"fmt"
	mux "github.com/gorilla/mux"
	"net/http"
	"os"
	"time"
)

type tempData struct {
	data        []byte
	prevRefresh int64
}

var temporaryData = make(map[string]tempData)

func ServerUp() {

	router := mux.NewRouter()
	router.HandleFunc("/", head)
	router.HandleFunc("/api", handleConnection)
	router.HandleFunc("/api/", handleCountry)

	server := &http.Server{
		Addr:    c.serverHost,
		Handler: router,
	}

	fmt.Printf("server work on: %s\n", c.serverHost)
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("local server down: %s", err)
		os.Exit(1)
	}
}

func head(w http.ResponseWriter, r *http.Request) {

	_, _ = w.Write([]byte("data only on /api"))
}

func handleConnection(w http.ResponseWriter, r *http.Request) {

	if notSupportedMethod(r.Method, http.MethodGet, w) {
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	country := r.URL.Query().Get("c")
	key := "getResultT"

	if country != "" {
		var err error
		country, err = getCountryISO3166_1(country, 69)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		key += country
	}

	if temporaryData[key].prevRefresh <
		time.Now().Unix()-c.refreshTimeSec {

		data, err := getResultTEmail(country)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		temporaryData[key] = tempData{
			data:        data,
			prevRefresh: time.Now().Unix()}
	}

	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "%s\n", temporaryData[key].data)
}

func handleCountry(w http.ResponseWriter, r *http.Request) {

	if notSupportedMethod(r.Method, http.MethodGet, w) {
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	country := r.URL.Query().Get("c")
	key := "getResultT"

	if country != "" {
		var err error
		country, err = getCountryISO3166_1(country, 69)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		key += country
	}

	if temporaryData[key].prevRefresh <
		time.Now().Unix()-c.refreshTimeSec {

		data, err := getResultTEmail(country)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		temporaryData[key] = tempData{
			data:        data,
			prevRefresh: time.Now().Unix()}
	}

	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "%s\n", temporaryData[key].data)
}

func notSupportedMethod(request string, support string,
	w http.ResponseWriter) bool {

	if request != support {
		w.Header().Set("Support", support)
		err := fmt.Sprintf("method \"%s\" not allowed [%s]", request, support)
		http.Error(w, err, http.StatusMethodNotAllowed)
		return true
	}
	return false
}

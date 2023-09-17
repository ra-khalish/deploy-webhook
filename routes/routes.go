package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/ra-khalish/deploy-webhook/models"
	"github.com/ra-khalish/deploy-webhook/services"
)

type malformedRequest struct {
	status int
	msg    string
}

func (mr *malformedRequest) Error() string {
	return mr.msg
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	if r.Header.Get("Content-Type") != "" {
		value := r.Header.Get("Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			log.Println(msg)
			return &malformedRequest{status: http.StatusUnsupportedMediaType, msg: msg}
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	//dec.DisallowUnknownFields()
	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			log.Println(msg)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			log.Println(msg)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Offset)
			log.Println(msg)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			log.Println(msg)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			log.Println(msg)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case err.Error() == "http: request body too large":
			msg := "Request body must not be largen than 1MB"
			log.Println(msg)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		msg := "Request body must only contain a single JSON object"
		log.Println(msg)
		return &malformedRequest{status: http.StatusBadRequest, msg: msg}
	}
	return nil
}

func getAll(w http.ResponseWriter, r *http.Request) {
	var me models.MinioEvent

	err := decodeJSONBody(w, r, &me)
	if err != nil {
		var mr malformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)
		} else {
			log.Print(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	log.Println(me.Key)
	fmt.Fprintf(w, "OK")

	//var me models.MinioEvent
	//err := json.NewDecoder(r.Body).Decode(&me)

	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusBadRequest)
	//	log.Fatal(err.Error())
	//	return
	//}

	//log.Println(me)

	//header := r.Header
	//defer r.Body.Close()
	//body, err := httputil.DumpRequest(r, true)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Printf("Request:\n%s\n", string(body))
	//log.Printf("Header:\n%s\n", header)
	//result := services.All("kubisa-id")

	//if result != nil {
	//	log.Println(result)
	//}
	//fmt.Fprintln(w, me)
}

func apply(w http.ResponseWriter, r *http.Request) {
	var me models.MinioEvent

	err := decodeJSONBody(w, r, &me)
	if err != nil {
		var mr malformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)
		} else {
			log.Print(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	ns := r.URL.Query().Get("ns")
	result := services.Apply(&ns, &me.Key)

	//if result != nil {
	//	log.Println(result)
	//}
	log.Println(result)
	fmt.Fprintf(w, "OK")
}

func Init() {
	mux := http.NewServeMux()

	//mux.HandleFunc("/all", getAll)
	mux.HandleFunc("/apply", apply)

	err := http.ListenAndServe(":8090", mux)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("listening...")
}

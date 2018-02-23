package martialarts

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"regexp"
	"fmt"
)

type TrainingUnitService struct {
	Repository TrainingUnitRepository
}

func (s TrainingUnitService) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.get(rw, r)
	case http.MethodPost:
		s.post(rw, r)
	}
}

func (s TrainingUnitService) get(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	trainingUnitRegex := regexp.MustCompile("^/training-unit/([^/]+)/([^/]+)$")

	trainingSeriesName := trainingUnitRegex.ReplaceAll([]byte(r.URL.Path), []byte("$1"))

	trainingUnitIndex := trainingUnitRegex.ReplaceAll([]byte(r.URL.Path), []byte("$2"))

	trainingUnit, err := s.Repository.Read(string(trainingSeriesName), string(trainingUnitIndex))
	if err != nil {
		fmt.Println(err)
		switch err {
		case NotFound:
			rw.WriteHeader(http.StatusNotFound)
		case Invalid:
			rw.WriteHeader(http.StatusInternalServerError)
		default:
			rw.WriteHeader(http.StatusInternalServerError)
		}
	}

	trainingUnitJSON, err := json.Marshal(trainingUnit)
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
	}

	rw.Write(trainingUnitJSON)
}

func (s TrainingUnitService) post(rw http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
	}

	trainingUnit := TrainingUnit{}
	err = json.Unmarshal(content, &trainingUnit)
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
	}

	index, err := s.Repository.Save(trainingUnit)
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
	}

	rw.Write([]byte(index))
}

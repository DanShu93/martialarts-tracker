package service

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"regexp"
	"fmt"
	"github.com/DanShu93/martialarts-tracker/storage"
	"github.com/DanShu93/martialarts-tracker/entity"
)

type TrainingUnitService struct {
	Repository    Repository
	UUIDGenerator UUIDGenerator
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

	trainingUnitRegex := regexp.MustCompile("^/training-unit/([^/]+)$")

	trainingUnitIndex := trainingUnitRegex.ReplaceAll([]byte(r.URL.Path), []byte("$1"))

	trainingUnit := entity.TrainingUnit{}
	err := s.Repository.Read(string(trainingUnitIndex), &trainingUnit)

	if err != nil {
		fmt.Println(err)
		switch err {
		case storage.NotFound:
			rw.WriteHeader(http.StatusNotFound)
		case storage.Invalid:
			rw.WriteHeader(http.StatusInternalServerError)
		default:
			rw.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	trainingUnitJSON, err := json.Marshal(trainingUnit)
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)

		return
	}

	rw.Write(trainingUnitJSON)
}

func (s TrainingUnitService) post(rw http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
	}

	trainingUnit := entity.TrainingUnit{}
	err = json.Unmarshal(content, &trainingUnit)
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusBadRequest)

		return
	}

	trainingUnit.ID = s.UUIDGenerator.Generate()

	err = s.Repository.Save(trainingUnit)
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)

		return
	}

	rw.Write([]byte(trainingUnit.ID))
}
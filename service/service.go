package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ilja2209/batch-manager-ui-client/db"
)

type Service struct {
	processRepository *db.ProcessRepository
}

func NewService(processRepository *db.ProcessRepository) *Service {
	return &Service{
		processRepository: processRepository,
	}
}

func (service *Service) GetProcessesHandler(writer http.ResponseWriter, request *http.Request) {
	processes, err := service.processRepository.GetProcesses()
	if err != nil {
		//todo: log error
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = writer.Write([]byte(err.Error()))
		return
	}

	writer.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(writer).Encode(processes)

	if err != nil {
		// todo: log error
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = writer.Write([]byte(err.Error()))
		return
	}
}

func (service *Service) GetProcessesByIdHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	idStr := vars["id"]

	id, _ := strconv.ParseInt(idStr, 10, 64)

	process, err := service.processRepository.GetProcessById(id)

	if err != nil {
		fmt.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = writer.Write([]byte(err.Error()))
		return
	}

	writer.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(writer).Encode(process)

	if err != nil {
		// todo: log error
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = writer.Write([]byte(err.Error()))
		return
	}
}

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ilja2209/batch-manager-ui-client/db"
	pb "github.com/ilja2209/batch-manager-ui-client/go-grpc/proto"
)

type Service struct {
	processRepository     *db.ProcessRepository
	stopProcessAuthorized bool
	batchManager          pb.BatchManagerServiceClient
}

func NewService(
	processRepository *db.ProcessRepository,
	stopProcessAuthorized bool,
	batchManager pb.BatchManagerServiceClient,
) *Service {
	return &Service{
		processRepository:     processRepository,
		stopProcessAuthorized: stopProcessAuthorized,
		batchManager:          batchManager,
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
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = writer.Write([]byte(err.Error()))
		return
	}
}

func (service *Service) StopProcessHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["id"]
	fmt.Println(id)

	if !service.stopProcessAuthorized {
		writer.WriteHeader(http.StatusUnauthorized)
		_, _ = writer.Write([]byte("Unauthorized to stop process " + id))
		return
	}

	clientReq := &pb.Get{
		Id: id,
	}

	_, err := service.batchManager.StopProcess(context.Background(), clientReq)
	if err != nil {
		writer.WriteHeader(http.StatusServiceUnavailable)
		_, _ = writer.Write([]byte(fmt.Sprintf("Failed to call stopProcess: %v", err)))
		return
	}

	writer.WriteHeader(http.StatusAccepted)
	//_, _ = writer.Write([]byte(""))
}

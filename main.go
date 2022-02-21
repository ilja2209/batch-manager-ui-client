package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/ilja2209/batch-manager-ui-client/db"
	pb "github.com/ilja2209/batch-manager-ui-client/go-grpc/proto"
	"github.com/ilja2209/batch-manager-ui-client/service"
	"github.com/ilja2209/batch-manager-ui-client/utils"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello world!")
	dbCon, err := db.NewDatabase(
		utils.GetEnvOrPanic("DB_HOST"), //format: <host>:<port>
		utils.GetEnvOrPanic("DB_USER"),
		utils.GetEnvOrPanic("DB_PASSWORD"),
		utils.GetEnvOrPanic("DB_NAME"),
	)

	if err != nil {
		panic(err)
	}
	defer dbCon.Close()

	procRepository := db.NewProcessRepository(dbCon)

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	bmConn, err := grpc.Dial(utils.GetEnvOrPanic("BATCH_MANAGER_URL"), opts...)

	if err != nil {
		log.Fatalf("Fail to dial batch-manager-service: %v", err)
		return
	}
	defer bmConn.Close()
	batchManagerClient := pb.NewBatchManagerServiceClient(bmConn)

	procService := service.NewService(procRepository, utils.GetEnvAsBool("STOP_PROCESS_AUTHORIZED", false), batchManagerClient)

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/api/v1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	}).Methods("GET", "POST")

	router.HandleFunc("/api/v1/processes", procService.GetProcessesHandler).Methods("GET")
	router.HandleFunc("/api/v1/processes/{id}", procService.GetProcessesByIdHandler).Methods("GET")
	router.HandleFunc("/api/v1/processes/{id}", procService.StopProcessHandler).Methods("DELETE")

	staticFilesPath := utils.GetEnv("STATIC_PATH", "./static/")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(staticFilesPath)))

	srv := &http.Server{
		Handler:      router,
		Addr:         ":" + utils.GetEnvOrPanic("SERVICE_PORT"),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Printf("Server is running")
	log.Fatal(srv.ListenAndServe())
}

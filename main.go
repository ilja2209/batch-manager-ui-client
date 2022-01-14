package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/ilja2209/batch-manager-ui-client/db"
	"github.com/ilja2209/batch-manager-ui-client/service"
	"github.com/ilja2209/batch-manager-ui-client/utils"
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
	procService := service.NewService(procRepository)

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/api/v1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	}).Methods("GET", "POST")

	router.HandleFunc("/api/v1/processes", procService.GetProcessesHandler).Methods("GET")
	router.HandleFunc("/api/v1/processes/{id}", procService.GetProcessesByIdHandler).Methods("GET")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	srv := &http.Server{
		Handler:      router,
		Addr:         ":" + utils.GetEnvOrPanic("SERVICE_PORT"),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Printf("Server is running")
	log.Fatal(srv.ListenAndServe())
}

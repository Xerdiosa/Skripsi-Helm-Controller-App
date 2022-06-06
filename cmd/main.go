package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gudangada/data-warehouse/warehouse-controller/internal/configs"
	"github.com/gudangada/data-warehouse/warehouse-controller/internal/routes"
)

func main() {
	var rt routes.Route

	appConfigs, err := configs.InitAppConfigs()
	if err != nil {
		panic(err)
	}

	router := rt.Init(*appConfigs)

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		sig := <-gracefulStop
		fmt.Printf("caught sig: %+v", sig)
		fmt.Println("Wait for 2 second to finish processing")
		time.Sleep(2 * time.Second)
		os.Exit(0)
	}()

	headersOK := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "NAME", "MODULE_NAME", "VERSION", "ON_FAILURE"})
	originsOK := handlers.AllowedOrigins([]string{"*"})
	methodsOK := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "DELETE", "PUT"})
	host := appConfigs.Server.Host
	port := strconv.Itoa(appConfigs.Server.Port)
	fmt.Println("Server served at port " + port)
	if err := http.ListenAndServe(host+":"+port, handlers.CORS(originsOK, headersOK, methodsOK)(router)); err != nil {
		log.Fatal("Unable to start service: " + err.Error())
	}
}

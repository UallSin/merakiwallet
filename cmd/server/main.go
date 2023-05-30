package main

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"meraki/finalwallet/config"
	"meraki/finalwallet/internal/handler"
	"meraki/finalwallet/internal/migration"
	"meraki/finalwallet/internal/repo"
	"meraki/finalwallet/internal/services"
	"meraki/finalwallet/pkg/pg"
	"net/http"
)

func main() {
	config.SetEnv()
	db, err := pg.ConnectDB(config.AppConfig{
		Host:     config.LoadEnv().Host,
		Port:     config.LoadEnv().Port,
		Username: config.LoadEnv().Username,
		Password: config.LoadEnv().Password,
		Dbname:   config.LoadEnv().Dbname,
	})
	if err != nil {
		logrus.Errorf("Failed to connect db: %v", err.Error())
		return
	}
	logrus.Infof("Connect to db successfull. Db name: %s", db.Name())
	logrus.Infof("Start http sever at 8080")

	userRepo := repo.NewUserRepo(db)
	userSrv := services.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSrv)
	r := mux.NewRouter()

	r.HandleFunc("/register", userHandler.Register).Methods("POST")
	r.HandleFunc("/migrate", migration.Migrate).Methods("GET")

	logrus.Infof("API register: %s", "/register POST")
	logrus.Infof("API migrate: %s", "/migrate GET")

	if err := http.ListenAndServe(":8080", r); err != nil {
		logrus.Errorf("Failed to get start server, err: %v", err)
		return
	}
}

package migration

import (
	"github.com/sirupsen/logrus"
	"meraki/finalwallet/config"
	"meraki/finalwallet/internal/model"
	"meraki/finalwallet/pkg/pg"
	"net/http"
)

func Migrate(w http.ResponseWriter, r *http.Request) {
	config.SetEnv()

	db, err := pg.ConnectDB(config.AppConfig{
		Host:     config.LoadEnv().Host,
		Port:     config.LoadEnv().Port,
		Username: config.LoadEnv().Username,
		Password: config.LoadEnv().Password,
		Dbname:   config.LoadEnv().Dbname,
	})
	if err != nil {
		logrus.Errorf("Failed to migrate the db: %v", err.Error())
		return
	}
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		logrus.Errorf("Failed to automigrate the db: %v", err.Error())
		return
	}
	logrus.Infof("Migrate successful")
}

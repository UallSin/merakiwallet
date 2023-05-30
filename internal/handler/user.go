package handler

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"meraki/finalwallet/internal/model"
	"meraki/finalwallet/internal/services"
	"net/http"
)

type UserHandler struct {
	userSrv services.UserService
}

func NewUserHandler(userSrv services.UserService) UserHandler {
	return UserHandler{
		userSrv: userSrv,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	requestUser := model.UserRequest{}
	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&requestUser)
	if err != nil {
		logrus.Errorf("Failed to get request body : %v", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		if err != nil {
			logrus.Errorf(err.Error())
			return
		}
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestUser.Password), bcrypt.DefaultCost)
	if err != nil {
		logrus.Errorf("Failed hashed password :%v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		err := json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		if err != nil {
			logrus.Errorf(err.Error())
			return
		}
	}
	if err := h.userSrv.Register(requestUser.Email, string(hashedPassword)); err != nil {
		logrus.Errorf("Failed to create user : %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		err := json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		if err != nil {
			logrus.Errorf(err.Error())
			return
		}
	}
	if err = json.NewEncoder(w).Encode(requestUser); err != nil {
		logrus.Errorf("Failed to encode respone: %v", err.Error())
		return
	}
}

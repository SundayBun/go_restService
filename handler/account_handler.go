package handler

import (
	"errors"
	"github.com/labstack/echo/v4"
	"goWebService/config"
	"goWebService/models"
	"goWebService/repository"
	"log"
	"net/http"
)

type Handlers interface {
	Save() echo.HandlerFunc
	Update() echo.HandlerFunc
	GetByID() echo.HandlerFunc
	Delete() echo.HandlerFunc
}

type accountHandler struct {
	cfg        *config.Config
	repository repository.Repository
}

func NewAccountHandler(cfg *config.Config, repository repository.Repository) Handlers {
	return &accountHandler{cfg: cfg, repository: repository}
}

func (a accountHandler) Save() echo.HandlerFunc {
	return func(c echo.Context) error {
		account := &models.AccountModel{}
		if err := c.Bind(account); err != nil {
			log.Printf("AccountHandler.Save: %v", err)
			return c.JSON(http.StatusBadRequest, BadRequest)
		}

		createdAccount, err := a.repository.Save(c.Request().Context(), account)
		if err != nil {
			log.Printf("AccountHandler.Save: %v", err)
			return c.JSON(http.StatusBadRequest, BadRequest)
		}

		return c.JSON(http.StatusCreated, createdAccount)
	}
}

func (a accountHandler) Update() echo.HandlerFunc {
	//TODO implement me
	panic("implement me")
}

func (a accountHandler) GetByID() echo.HandlerFunc {
	//TODO implement me
	panic("implement me")
}

func (a accountHandler) Delete() echo.HandlerFunc {
	//TODO implement me
	panic("implement me")
}

var (
	BadRequest = errors.New("Bad request")
)

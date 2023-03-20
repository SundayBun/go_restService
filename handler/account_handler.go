package handler

import (
	"errors"
	"github.com/labstack/echo/v4"
	"goWebService/config"
	"goWebService/models"
	"goWebService/repository"
	"log"
	"net/http"
	"strconv"
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
	return func(c echo.Context) error {
		account := &models.AccountModel{}
		if err := c.Bind(account); err != nil {
			log.Printf("AccountHandler.Save: %v", err)
			return c.JSON(http.StatusBadRequest, BadRequest)
		}

		updatedAccount, err := a.repository.Update(c.Request().Context(), account)
		if err != nil {
			log.Printf("AccountHandler.Update: %v", err)
			return c.JSON(http.StatusBadRequest, BadRequest)
		}

		return c.JSON(http.StatusCreated, updatedAccount)
	}
}

func (a accountHandler) GetByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Printf("AccountHandler.GetByID: %v", err)
			return c.JSON(http.StatusBadRequest, BadRequest)
		}

		account, err := a.repository.GetById(c.Request().Context(), id)
		return c.JSON(http.StatusOK, account)
	}
}

func (a accountHandler) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Printf("AccountHandler.Delete: %v", err)
			return c.JSON(http.StatusBadRequest, BadRequest)
		}

		err = a.repository.DeleteById(c.Request().Context(), id)
		if err != nil {
			log.Printf("AccountHandler.Delete: %v", err)
			return c.JSON(http.StatusBadRequest, BadRequest)
		}
		return c.JSON(http.StatusOK, Success)
	}
}

type SuccessResponse struct {
	s string
}

var (
	BadRequest = errors.New("Bad request")
	Success    = &SuccessResponse{"Success"}
)

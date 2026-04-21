package controller

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/laurencefluciano/content-api/internal/exception"
)

func (h *UserController) handleAppError(c *gin.Context, err error) {
	var appErr *exception.AppError

	if errors.As(err, &appErr) {
		status := http.StatusBadRequest

		log.Print("[API ERROR]", err.Error())

		switch appErr.Code {
		case exception.EntityNotFoundCode:
			status = http.StatusNotFound
		case exception.EntityConflictCode, exception.RoleAssignedCode:
			status = http.StatusConflict
		case exception.ValidationFailedCode:
			status = http.StatusUnprocessableEntity
		case exception.EmptyUpdateCode:
			c.Status(http.StatusNoContent)
			return
		}

		c.JSON(status, appErr)

		return
	}

	log.Printf("[INTERNAL ERROR]: %v", err)

	c.JSON(http.StatusInternalServerError, gin.H{
		"status":  500,
		"message": "Ocorreu um erro inesperado em nossos servidores.",
	})
}

package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	models "medidor_enerbit/models"
	"medidor_enerbit/utils"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func CreateMedidor(c *gin.Context) {
	var medidor models.Medidor

	request_id := c.GetString("x-request-id")

	// Bind request payload with our model
	if binderr := c.ShouldBindJSON(&medidor); binderr != nil {

		log.Error().Err(binderr).Str("request_id", request_id).
			Msg("Error occurred while binding request data")

		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": binderr.Error(),
		})
		return
	}

	medidor.SetUUID()

	db, conErr := utils.GetDatabaseConnection()
	if conErr != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"message": "Service is unavailable",
		})
		return
	}

	result := db.Create(&medidor)
	if result.Error != nil && result.RowsAffected != 1 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error occurred while creating a new medidor",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Medidor created successfully",
		"id":      medidor.ID,
	})
}

func UpdateMedidor(c *gin.Context) {
	var medidor models.Medidor

	request_id := c.GetString("x-request-id")

	if binderr := c.ShouldBindJSON(&medidor); binderr != nil {

		log.Error().Err(binderr).Str("request_id", request_id).
			Msg("Error occurred while binding request data")

		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": binderr.Error(),
		})
		return
	}
	medidor.SetUUID()

	db, conErr := utils.GetDatabaseConnection()
	if conErr != nil {
		log.Err(conErr).Str("request_id", request_id).Msg("Error occurred while getting a DB connection from the connection pool")
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"message": "Service is unavailable",
		})
		return
	}

	var value models.Medidor

	result := db.First(&value, "id = ?", medidor.ID)
	if result.Error != nil {
		log.Err(result.Error).Str("request_id", request_id).Msg("Error occurred while updating the medidor")
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Record not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error occurred while updating medidor",
			})
		}
		return
	}

	value.Address = medidor.Address
	value.RetirementDate = medidor.RetirementDate
	value.Lines = medidor.Lines
	value.IsActive = medidor.IsActive

	tx := db.Save(&value)
	if tx.RowsAffected != 1 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error occurred while updating medidor",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Medidor updated successfully",
		"result":  value,
	})

}

func GetMedidor(c *gin.Context) {
	var medidorId models.MedidorID

	request_id := c.GetString("x-request-id")

	// Bind request payload with our model
	if binderr := c.ShouldBindUri(&medidorId); binderr != nil {

		log.Error().Err(binderr).Str("request_id", request_id).
			Msg("Error occurred while binding request data")

		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": binderr.Error(),
		})
		return
	}
	// Get a connection
	db, conErr := utils.GetDatabaseConnection()
	if conErr != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"message": "Service is unavailable",
		})
		return
	}

	var medidor models.Medidor
	result := db.First(&medidor, "id = ?", medidorId.ID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Record not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error occurred while fetching medidor",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": medidor,
	})
}

func GetMedidores(c *gin.Context) {
	var medidors []models.Medidor

	request_id := c.GetString("x-request-id")

	earliest := c.DefaultQuery("earliest", "0")
	latest := c.DefaultQuery("latest", fmt.Sprint(time.Now()))

	// Get a connection
	db, conErr := utils.GetDatabaseConnection()
	if conErr != nil {
		log.Err(conErr).Str("request_id", request_id).Msg("Error occurred while getting a DB connection from the connection pool")
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"message": "Service is unavailable",
		})
		return
	}

	tx := db.Find(&medidors)

	if tx.RowsAffected == 0 {
		log.Info().Msg("Read medidores returned with empty results")
	}
	c.JSON(http.StatusOK, gin.H{
		"results": medidors,
	})
}

func DeleteMedidor(c *gin.Context) {
	var medidorId models.MedidorID

	request_id := c.GetString("x-request-id")

	if binderr := c.ShouldBindUri(&medidorId); binderr != nil {

		log.Error().Err(binderr).Str("request_id", request_id).
			Msg("Error occurred while binding request data")

		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": binderr.Error(),
		})
		return
	}

	// Get a connection
	db, conErr := utils.GetDatabaseConnection()
	if conErr != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"message": "Service is unavailable",
		})
		return
	}

	var medidor models.Medidor
	result := db.First(&medidor, "id = ?", medidorId.ID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Record not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error occurred while deleting medidor",
			})
		}
		return
	}

	tx := db.Delete(&medidor)
	if tx.RowsAffected != 1 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error occurred while deleting medidor",
		})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{})
}

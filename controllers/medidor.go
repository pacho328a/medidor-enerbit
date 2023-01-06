package controllers

import (
	"errors"
	"net/http"

	models "medidor_enerbit/models"
	redis "medidor_enerbit/stream"
	"medidor_enerbit/utils"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// PostMedidor	 godoc
// @Summary      Create a new medidor
// @Description  Takes a Medidor JSON and store in DB postgres.
// @Tags         Medidor
// @Produce      json
// @Param        Medidor  body      models.MedidorSwCreate  true  "Medidor JSON"
// @Success      200   {object}  models.MedidorResponse
// @Router       /medidor [post]
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

	client := redis.GetRedis()
	err := redis.SendStreamMedidor(medidor, client)

	if err != nil {
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

// UpdateMedidorById		 godoc
// @Summary      Update single Medidor by id
// @Description  Updates and returns a single Medidor whose Id value matches the id. New data must be passed in the body.
// @Tags         Medidor
// @Produce      json
// @Param        Medidor  body      models.MedidorSwUpdate  true  "update Medidor by id"
// @Success      200  {object}  models.MedidorSwUpdateResponse
// @Router       /medidor [PATCH]
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

// GetMedidorById		 godoc
// @Summary      Get single Medidor by Id
// @Description  Returns the Medidor whose Id value matches the Id.
// @Tags         Medidor
// @Produce      json
// @Param        id  path      string  true  "search Medidor by Id"
// @Success      200  {object}  models.Medidor
// @Router       /medidor/{id} [get]
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

// GetMedidores	 godoc
// @Summary      Get Medidores array
// @Description  Responds with the list of all Medidores as JSON.
// @Tags         Medidores
// @Produce      json
// @Success      200  {array}  models.Medidor
// @Router       /medidores [get]
func GetMedidores(c *gin.Context) {
	var medidors []models.Medidor

	request_id := c.GetString("x-request-id")

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

// DeleteMedidorById		 godoc
// @Summary      Remove single Medidor by id
// @Description  Delete a single entry from the database based on id.
// @Tags         Medidor
// @Produce      json
// @Param        id  path      string  true  "delete Medidor by id"
// @Success      204
// @Router       /medidor/{id} [delete]
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

package models

import (
	"time"

	"github.com/google/uuid"
)

type MedidorID struct {
	ID string `uri:"id" binding:"required"`
}

type Medidor struct {
	ID               string    `gorm:"primaryKey" json:"id"`
	Brand            string    `json:"brand"`
	Address          string    `json:"address" binding:"required"`
	InstallationDate time.Time `json:"installationdate"`
	RetirementDate   time.Time `json:"retirementdate"`
	Serial           string    `json:"serial"`
	Lines            uint64    `json:"lines" binding:"required"`
	IsActive         bool      `json:"isactive"`
	CreatedAt        time.Time `gorm:"autoCreateTime:milli" json:"created_at"`
}

type MedidorSwCreate struct {
	Brand            string    `json:"brand" example:"Marca x"`
	Address          string    `json:"address" example:"Calle-street x"`
	InstallationDate time.Time `json:"installationdate" example:"2022-05-25T00:53:16.535668Z"`
	RetirementDate   time.Time `json:"retirementdate" example:"2022-05-25T00:53:16.535668Z"`
	Serial           string    `json:"serial" example:"Serial x"`
	// min: 1
	// example: 1
	Lines    uint64 `json:"lines" binding:"required"`
	IsActive bool   `json:"isactive"`
}

type MedidorSwUpdate struct {
	ID             string    `json:"id" example:"6cd9f3c8-7bc8-40e7-8a4b-b575e63f0..."`
	Address        string    `json:"address" example:"Calle-street x"`
	RetirementDate time.Time `json:"retirementdate" example:"2022-05-25T00:53:16.535668Z"`
	// min: 1
	// example: 1
	Lines    uint64 `json:"lines" binding:"required"`
	IsActive bool   `json:"isactive"`
}

type MedidorSwUpdateResponse struct {
	Message string `json:"Message" example:"Medidor ..."`
	Result  Medidor
}

type MedidorResponse struct {
	Id      string `json:"Id" example:"6cd9f3c8-7bc8-40e7-8a4b-b575e63f0..."`
	Message string `json:"Message" example:"Medidor ..."`
}

func (x *Medidor) SetUUID() {
	if x.ID == "" {
		x.ID = uuid.New().String()
	}
}

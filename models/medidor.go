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
	DeletedAt        time.Time `json:"deleted_at"`
}

func (x *Medidor) SetUUID() {
	if x.ID == "" {
		x.ID = uuid.New().String()
	}
}

package gRPC

import (
	"context"
	"errors"
	"fmt"
	"log"
	"medidor_enerbit/gRPC/medidorgRPC"
	models "medidor_enerbit/models"
	redis "medidor_enerbit/stream"
	"medidor_enerbit/utils"
	"net"

	"gorm.io/gorm"

	"google.golang.org/grpc"
)

const (
	gRpcPort = "50001"
)

type medidorService struct {
	medidorgRPC.UnimplementedMedidorServiceServer
}

func GRPCListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	var s = grpc.NewServer()
	medidorgRPC.RegisterMedidorServiceServer(s, &medidorService{})

	log.Printf("gRPC Server started on port %s", gRpcPort)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed 2 to listen for gRPC: %v", err)
	}
}

func (a *medidorService) WriteMedidor(ctx context.Context, req *medidorgRPC.MedidorRequest) (*medidorgRPC.MedidorCreateResponse, error) {

	err := req.Validate()
	if err != nil {
		res := &medidorgRPC.MedidorCreateResponse{Result: "Error"}
		return res, err
	}

	input := req.GetMedidorEntry()

	var medidor models.Medidor

	medidor.Brand = input.Brand
	medidor.Address = input.Address
	if input.InstallationDate != nil {
		medidor.InstallationDate = input.InstallationDate.AsTime()
	}
	if input.RetirementDate != nil {
		medidor.RetirementDate = input.RetirementDate.AsTime()
	}
	medidor.Serial = input.Serial
	medidor.Lines = input.Lines
	medidor.IsActive = input.IsActive

	medidor.SetUUID()

	db, conErr := utils.GetDatabaseConnection()
	if conErr != nil {
		res := &medidorgRPC.MedidorCreateResponse{Result: "error db"}
		return res, conErr
	}

	var medidors []models.Medidor
	db.Find(&medidors)

	flag := false
	flagd := false
	for i := 0; i < len(medidors); i++ {
		if medidors[i].Serial == input.Serial {
			if medidors[i].Brand == input.Brand {
				flag = true
			}
		}
		if medidors[i].Address == input.Address {
			flagd = true
		}
	}
	if flagd {
		res := &medidorgRPC.MedidorCreateResponse{Result: "Solo Puede existir 1 medidor por direccion"}
		return res, conErr
	}

	if flag {
		res := &medidorgRPC.MedidorCreateResponse{Result: "Serial y Marca ya existen"}
		return res, conErr
	} else {
		result := db.Create(&medidor)
		if result.Error != nil && result.RowsAffected != 1 {
			res := &medidorgRPC.MedidorCreateResponse{Result: "error creating medidor"}
			return res, conErr
		}

		client := redis.GetRedis()
		err := redis.SendStreamMedidor(medidor, client)
		if err != nil {
			res := &medidorgRPC.MedidorCreateResponse{Result: "error creating stream"}
			return res, conErr
		}
		res := &medidorgRPC.MedidorCreateResponse{
			Id:     medidor.ID,
			Result: "Medidor Created"}
		return res, conErr
	}
}

func (a *medidorService) GetMedidorInstalled(ctx context.Context, req *medidorgRPC.MedidorIsActive) (*medidorgRPC.MedidorIsActiveResponse, error) {

	input := req.GetIsActive()

	db, conErr := utils.GetDatabaseConnection()
	if conErr != nil {
		return nil, conErr
	}

	var medidors []models.Medidor
	db.Where("is_active = ?", input).Find(&medidors)

	medidorsg := make([]*medidorgRPC.MedidorGet, 0)

	for _, u := range medidors {
		medidorsg = append(medidorsg, &medidorgRPC.MedidorGet{
			Brand:            u.Brand,
			Address:          u.Address,
			Installationdate: u.InstallationDate.String(),
			Retirementdate:   u.RetirementDate.String(),
			Serial:           u.Serial,
			Lines:            int32(u.Lines),
			Isactive:         u.IsActive,
		})
	}

	res := &medidorgRPC.MedidorIsActiveResponse{
		Medidores: medidorsg}
	return res, conErr

}

func (a *medidorService) DeleteMedidor(ctx context.Context, req *medidorgRPC.MedidorUUID) (*medidorgRPC.MedidorResponse, error) {
	input := req.GetId()
	var medidor models.Medidor

	db, conErr := utils.GetDatabaseConnection()
	if conErr != nil {
		res := &medidorgRPC.MedidorResponse{Result: "error db"}
		return res, conErr
	}

	result := db.First(&medidor, "id = ?", input)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			res := &medidorgRPC.MedidorResponse{Result: "Record not found"}
			return res, conErr
		} else {
			res := &medidorgRPC.MedidorResponse{Result: "Error occurred while deleting medidor"}
			return res, conErr
		}
	}

	if medidor.IsActive {
		res := &medidorgRPC.MedidorResponse{
			Result: "Medidor not installed"}
		return res, conErr
	} else {
		tx := db.Delete(&medidor)
		if tx.RowsAffected != 1 {
			res := &medidorgRPC.MedidorResponse{
				Result: "Error occurred while deleting medidor"}
			return res, conErr
		}

		res := &medidorgRPC.MedidorResponse{
			Result: "Medidor Deleted"}
		return res, conErr
	}

}

func (a *medidorService) UpdateMedidor(ctx context.Context, req *medidorgRPC.MedidorUpdate) (*medidorgRPC.MedidorResponse, error) {

	var medidor models.Medidor

	db, conErr := utils.GetDatabaseConnection()
	if conErr != nil {
		res := &medidorgRPC.MedidorResponse{Result: "error db"}
		return res, conErr
	}

	result := db.First(&medidor, "id = ?", req.GetId())
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			res := &medidorgRPC.MedidorResponse{Result: "Record not found"}
			return res, conErr
		} else {
			res := &medidorgRPC.MedidorResponse{Result: "Error occurred while deleting medidor"}
			return res, conErr
		}
	}

	medidor.Address = req.GetAddress()
	medidor.RetirementDate = req.GetRetirementDate().AsTime()
	medidor.Lines = uint64(req.GetLines())
	medidor.IsActive = req.GetIsActive()

	tx := db.Save(&medidor)
	if tx.RowsAffected != 1 {
		res := &medidorgRPC.MedidorResponse{
			Result: "Error occurred while updating medidor"}
		return res, conErr
	}

	res := &medidorgRPC.MedidorResponse{
		Result: "Medidor Updated"}
	return res, conErr

}

func (a *medidorService) GetMedidor(ctx context.Context, req *medidorgRPC.MedidorUUID) (*medidorgRPC.MedidorGet, error) {
	input := req.GetId()
	var medidor models.Medidor

	db, conErr := utils.GetDatabaseConnection()
	if conErr != nil {
		return nil, conErr
	}

	result := db.First(&medidor, "id = ?", input)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, conErr
		} else {
			return nil, conErr
		}
	}

	res := &medidorgRPC.MedidorGet{
		Brand:            medidor.Brand,
		Address:          medidor.Address,
		Serial:           medidor.Serial,
		Installationdate: medidor.InstallationDate.String(),
		Retirementdate:   medidor.RetirementDate.String(),
		Lines:            int32(medidor.Lines),
		Isactive:         medidor.IsActive,
	}
	return res, conErr

}

func (a *medidorService) RecentInstallationMarca(ctx context.Context, req *medidorgRPC.MedidorMarca) (*medidorgRPC.MedidorGet, error) {
	input := req.GetMarca()

	db, conErr := utils.GetDatabaseConnection()
	if conErr != nil {
		return nil, conErr
	}

	var medidors []models.Medidor
	db.Order("installation_date desc").Where("brand = ?", input).Find(&medidors)
	medidor := medidors[0]
	res := &medidorgRPC.MedidorGet{
		Brand:            medidor.Brand,
		Address:          medidor.Address,
		Serial:           medidor.Serial,
		Installationdate: medidor.InstallationDate.String(),
		Retirementdate:   medidor.RetirementDate.String(),
		Lines:            int32(medidor.Lines),
		Isactive:         medidor.IsActive,
	}
	return res, conErr

}

func (a *medidorService) RecentInstallationSerial(ctx context.Context, req *medidorgRPC.MedidorSerial) (*medidorgRPC.MedidorGet, error) {
	input := req.GetSerial()

	db, conErr := utils.GetDatabaseConnection()
	if conErr != nil {
		return nil, conErr
	}

	var medidors []models.Medidor
	db.Order("installation_date desc").Where("serial = ?", input).Find(&medidors)
	medidor := medidors[0]
	res := &medidorgRPC.MedidorGet{
		Brand:            medidor.Brand,
		Address:          medidor.Address,
		Serial:           medidor.Serial,
		Installationdate: medidor.InstallationDate.String(),
		Retirementdate:   medidor.RetirementDate.String(),
		Lines:            int32(medidor.Lines),
		Isactive:         medidor.IsActive,
	}
	return res, conErr

}

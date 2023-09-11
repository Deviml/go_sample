package main

import (
	"fmt"

	"github.com/Equiphunter-com/equipment-hunter-api/pkg/models"

	//"github.com/Equiphunter-com/equipment-hunter-api/pkg/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DBConfig struct {
	Host       string `env:"DB_HOST,required"`
	Port       string `env:"DB_PORT,required"`
	DBName     string `env:"DB_NAME,required"`
	DBUser     string `env:"DB_USER,required"`
	DBPassword string `env:"DB_PASSWORD,required"`
}

func (db DBConfig) buildDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", db.DBUser, db.DBPassword, db.Host, db.Port, db.DBName)
}

func buildDB() (*gorm.DB, error) {
	dbConfig := &DBConfig{
		Host:       "equiphunterdb.chtwvgzu6d0b.us-east-2.rds.amazonaws.com",
		Port:       "3306",
		DBName:     "equiphunterdb",
		DBUser:     "admin",
		DBPassword: "equiphuntersecret",
	}
	return gorm.Open(mysql.Open(dbConfig.buildDSN()), &gorm.Config{})
}

func main() {
	db, err := buildDB()
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}
	err = db.AutoMigrate(
		&models.VerificationCode{},
		&models.WebUser{},
		&models.Subcontractor{},
		&models.GeneralContractor{},
		&models.VendorRental{},
		&models.Company{},
		&models.Location{},
		&models.EquipmentCategory{},
		&models.EquipmentSubcategory{},
		&models.Equipment{},
		&models.EquipmentRequest{},
		&models.SupplyCategory{},
		&models.Supply{},
		&models.SupplyRequest{},
		&models.Quote{},
		&models.Sublist{},
		&models.City{},
		&models.SublistsCompany{},
		&models.County{},
		&models.State{},
		&models.Payment{},
		&models.PaymentDetail{},
		&models.InvitationCode{},
		&models.Proposal{},
	)
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
}

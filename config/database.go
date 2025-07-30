package config

import (
	"ayo-indonesia-api/app/models"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var DBS = make(map[string]*gorm.DB)

// Database Initialization
func Database() *gorm.DB {
	host := LoadConfig().DatabaseHost
	user := LoadConfig().DatabaseUsername
	password := LoadConfig().DatabasePassword
	name := LoadConfig().DatabaseName
	port := LoadConfig().DatabasePort

	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", host, user, password, name, port)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	if LoadConfig().EnableDatabaseAutomigration {
		err = DB.AutoMigrate(
			&models.GlobalUser{},
			&models.Team{},
			&models.GlobalFile{},
			&models.Player{},
			&models.Goal{},
			&models.Match{},
			&models.GlobalFile{},
		)
		if err != nil {
			log.Fatalf("Auto migration failed: %v", err)
		} else {
			fmt.Println("Auto migration success ...")

		}
	}

	fmt.Println("Connected to Database:", name)

	return DB
}

func GetRespectiveID(db *gorm.DB, tablename string, synchronizeSequence bool) (respectiveID uint, err error) {
	err = db.Unscoped().Table(tablename).Order("id desc").Limit(1).Pluck("id", &respectiveID).Error
	if err != nil {
		return
	}

	if respectiveID == 0 {
		err = fmt.Errorf("no data found in table %s", tablename)

		return
	}

	if respectiveID%2 == 0 {
		respectiveID += 2
	} else {
		respectiveID += 1
	}
	if synchronizeSequence {
		err = SynchronizeSequence(db, tablename, respectiveID)
	}

	return
}

func SynchronizeSequence(db *gorm.DB, table string, lastID uint) error {
	sequenceQuery := fmt.Sprintf(
		"SELECT setval(pg_get_serial_sequence('%s', 'id'), %v)",
		table, lastID,
	)
	err := DB.Exec(sequenceQuery).Error
	if err != nil {
		return fmt.Errorf("failed to synchronize sequence for table %s: %w", table, err)
	}

	return nil
}

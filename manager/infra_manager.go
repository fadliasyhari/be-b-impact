package manager

import (
	"context"
	"errors"
	"fmt"
	"os"

	"be-b-impact.com/csr/config"
	"be-b-impact.com/csr/model"
	"be-b-impact.com/csr/utils/constants"
	firebase "firebase.google.com/go/v4"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type InfraManager interface {
	Conn() *gorm.DB
	Migrate(model ...any) error
	Log() *logrus.Logger
	LogFilePath() string
	UploadPath() string
	FirebaseApp() *firebase.App
}

type infraManager struct {
	db          *gorm.DB
	cfg         *config.Config
	log         *logrus.Logger
	firebaseApp *firebase.App
}

// FirebaseApp implements InfraManager
func (i *infraManager) FirebaseApp() *firebase.App {
	return i.firebaseApp
}

func (i *infraManager) UploadPath() string {
	return i.cfg.UploadPath
}

func (i *infraManager) Log() *logrus.Logger {
	return logrus.New()
}

func (i *infraManager) LogFilePath() string {
	return i.cfg.LogPath
}

func (i *infraManager) Conn() *gorm.DB {
	return i.db
}

func insertSeedData(db *gorm.DB, data interface{}) error {
	return db.Clauses(clause.OnConflict{DoNothing: true}).Create(data).Error
}

func (i *infraManager) Migrate(models ...any) error {
	db := i.Conn()
	err := db.AutoMigrate(models...)
	if err != nil {
		return err
	}
	if db.Migrator().HasTable(&model.User{}) {
		if err := db.First(&model.User{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			seedUser := constants.UserSeed
			if err := insertSeedData(db, seedUser); err != nil {
				return err
			}
		}
	}
	if db.Migrator().HasTable(&model.Progress{}) {
		if err := db.First(&model.Progress{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			seedProgress := constants.ProgressSeed
			if err := insertSeedData(db, seedProgress); err != nil {
				return err
			}
		}
	}
	if db.Migrator().HasTable(&model.Category{}) {
		var count int64
		db.Model(&model.Category{}).Count(&count)
		if count == 0 {
			categories := constants.CategorySeed
			if err := insertSeedData(db, categories); err != nil {
				return err
			}
		}
	}
	if db.Migrator().HasTable(&model.Tag{}) {
		var count int64
		db.Model(&model.Tag{}).Count(&count)
		if count == 0 {
			categories := constants.TagSeed
			if err := insertSeedData(db, categories); err != nil {
				return err
			}
		}
	}
	return nil
}

func (i *infraManager) initDb() error {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		i.cfg.Host, i.cfg.Port, i.cfg.User, i.cfg.Password, i.cfg.Name)
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	i.db = conn

	if i.cfg.FileConfig.Env == "DEV" {
		i.db = conn.Debug()
	} else {
		// production / release
	}
	return nil
}

func NewInfraManager(cfg *config.Config) (InfraManager, error) {
	conn := &infraManager{cfg: cfg}

	credentialsFilePath := os.Getenv("FIREBASE_CRED_JSON")
	if credentialsFilePath == "" {
		return nil, errors.New("missing required environment variables")
	}

	// Initialize Firebase app
	ctx := context.Background()
	opt := option.WithCredentialsFile(credentialsFilePath)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}

	// Set the Firebase app instance in the infraManager
	conn.firebaseApp = app

	err = conn.initDb()
	if err != nil {
		return nil, err
	}
	return conn, nil
}

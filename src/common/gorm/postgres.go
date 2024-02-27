package appgorm

import (
	"context"
	"fmt"
	"go-kit/src/common/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBProvider struct {
	db *gorm.DB
}

func (p *DBProvider) Start(ctx context.Context) error {
	return nil
}

func (p *DBProvider) Stop(ctx context.Context) error {
	sqlDB, err := p.db.DB()
	if err != nil {
		return fmt.Errorf("[Postgres] failed to disconnect %w", err)
	}
	err = sqlDB.Close()
	if err != nil {
		return fmt.Errorf("[Postgres] failed to disconnect %w", err)
	}
	return nil
}

func (p *DBProvider) DB() *gorm.DB {
	return p.db
}

func NewPostgres(config *configs.Config) (*DBProvider, error) {
	cf := config.Postgres
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s", cf.Host,
		cf.Port, cf.User, cf.DbName, cf.SslMode, cf.Password)

	gormConfig := &gorm.Config{Logger: logger.Default.LogMode(logger.Info)}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("[Postgres] failed to connect to DB %w", err)
	}

	return &DBProvider{db: db}, nil
}

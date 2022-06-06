package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gudangada/data-warehouse/warehouse-controller/internal/configs"
	"github.com/gudangada/data-warehouse/warehouse-controller/internal/models"
	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var database *gorm.DB

func GetDB(config configs.DBConfig) (*gorm.DB, error) {
	if database != nil {
		return database, nil
	}
	var err error
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=allow TimeZone=UTC",
		config.Host,
		config.User,
		config.Password,
		config.Name,
		config.Port,
	)
	database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger(),
	})
	if err != nil {
		return nil, err
	}

	err = database.AutoMigrate(&models.ChartRelease{})
	if err != nil {
		return nil, err
	}

	err = database.AutoMigrate(&models.Module{})
	if err != nil {
		return nil, err
	}

	err = database.AutoMigrate(&models.ModuleRelease{})
	if err != nil {
		return nil, err
	}

	err = database.AutoMigrate(&models.Kinesis{})
	if err != nil {
		return nil, err
	}

	return database, nil
}

type zerologAdapter struct {
	log zerolog.Logger
}

func (z *zerologAdapter) LogMode(level logger.LogLevel) logger.Interface {
	switch level {
	case logger.Error:
		z.log = z.log.With().Logger().Level(zerolog.ErrorLevel)
	case logger.Warn:
		z.log = z.log.With().Logger().Level(zerolog.WarnLevel)
	case logger.Info:
		z.log = z.log.With().Logger().Level(zerolog.InfoLevel)
	case logger.Silent:
		z.log = z.log.With().Logger().Level(zerolog.Disabled)
	}
	return z
}

func (z *zerologAdapter) Info(ctx context.Context, s string, i ...interface{}) {
	z.log.Info().Msgf(s, i...)
}

func (z *zerologAdapter) Warn(ctx context.Context, s string, i ...interface{}) {
	z.log.Warn().Msgf(s, i...)
}

func (z *zerologAdapter) Error(ctx context.Context, s string, i ...interface{}) {
	z.log.Error().Msgf(s, i...)
}

func (z *zerologAdapter) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	stmt, row := fc()
	z.log.Trace().Time("begin", begin).Str("stmt", stmt).Int64("row_affected", row).Err(err).Send()
}

func newLogger() logger.Interface {
	return &zerologAdapter{
		log: zerolog.New(os.Stderr).With().Timestamp().Logger(),
	}
}

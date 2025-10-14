// Package config
package config

import (
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DatabaseType string

const (
	MySQL      DatabaseType = "mysql"
	PostgreSQL DatabaseType = "postgres"
	SQLite     DatabaseType = "sqlite3"
)

var allowedDatabaseType = []DatabaseType{MySQL, PostgreSQL, SQLite}

type DatabaseConfig struct {
	Type                 string        `json:"type"`
	DBType               DatabaseType  `json:"-"`
	Database             string        `json:"database"`
	Host                 string        `json:"host"`
	Port                 int           `json:"port"`
	Username             string        `json:"username"`
	Password             string        `json:"password"`
	EnableSSL            bool          `json:"enable_ssl"`
	ConnectIdleTimeout   string        `json:"connect_idle_timeout"` // 连接空闲超时时间
	ConnectIdleDuration  time.Duration `json:"-"`
	QueryTimeout         string        `json:"connect_timeout"` // 每次查询超时时间
	QueryDuration        time.Duration `json:"-"`
	ServerMaxConnections int           `json:"server_max_connections"` // 最大连接池大小
}

func defaultDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Type:                 "sqlite3",
		Database:             "database.db",
		Host:                 "",
		Port:                 0,
		Username:             "",
		Password:             "",
		EnableSSL:            false,
		ConnectIdleTimeout:   "1h",
		QueryTimeout:         "5s",
		ServerMaxConnections: 32,
	}
}

func (config *DatabaseConfig) checkValid(_ log.LoggerInterface) *ValidResult {
	config.DBType = DatabaseType(config.Type)
	if !slices.Contains(allowedDatabaseType, config.DBType) {
		return ValidFail(fmt.Errorf("database type %s is not allowed, support database is %v, please check the configuration file", config.DBType, allowedDatabaseType))
	}

	if duration, err := time.ParseDuration(config.ConnectIdleTimeout); err != nil {
		return ValidFailWith(errors.New("invalid json field connect_idel_timeout"), err)
	} else {
		config.ConnectIdleDuration = duration
	}

	if duration, err := time.ParseDuration(config.QueryTimeout); err != nil {
		return ValidFailWith(errors.New("invalid json field query_timeout"), err)
	} else {
		config.QueryDuration = duration
	}
	return ValidPass()
}

func (config *DatabaseConfig) GetConnection(logger log.LoggerInterface) gorm.Dialector {
	switch config.DBType {
	case MySQL:
		return mySQLConnection(logger, config)
	case PostgreSQL:
		return postgreSQLConnection(logger, config)
	case SQLite:
		return sqliteConnection(logger, config)
	default:
		return nil
	}
}

func mySQLConnection(logger log.LoggerInterface, db *DatabaseConfig) gorm.Dialector {
	var enableSSL string
	if db.EnableSSL {
		enableSSL = "true"
	} else {
		enableSSL = "false"
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&tls=%s",
		db.Username,
		db.Password,
		db.Host,
		db.Port,
		db.Database,
		enableSSL,
	)
	logger.DebugF("Mysql Connection DSN %s", dsn)
	return mysql.Open(dsn)
}

func postgreSQLConnection(logger log.LoggerInterface, db *DatabaseConfig) gorm.Dialector {
	var enableSSL string
	if db.EnableSSL {
		enableSSL = "enable"
	} else {
		enableSSL = "disable"
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Shanghai",
		db.Host,
		db.Username,
		db.Password,
		db.Database,
		db.Port,
		enableSSL,
	)
	logger.DebugF("PostgreSQL Connection DSN %s", dsn)
	return postgres.Open(dsn)
}

func sqliteConnection(_ log.LoggerInterface, db *DatabaseConfig) gorm.Dialector {
	return sqlite.Open(db.Database)
}

package database

import (
	"context"
	. "fmt"
	"github.com/half-nothing/simple-fsd/internal/interfaces/config"
	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
	. "github.com/half-nothing/simple-fsd/internal/interfaces/operation"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type DBCloseCallback struct {
	logger log.LoggerInterface
	db     *gorm.DB
}

func NewDBCloseCallback(logger log.LoggerInterface, db *gorm.DB) *DBCloseCallback {
	return &DBCloseCallback{
		logger: logger,
		db:     db,
	}
}

func (dc *DBCloseCallback) Invoke(ctx context.Context) error {
	dc.logger.InfoF("Closing database connection")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	db, err := dc.db.DB()
	if err != nil {
		return err
	}
	err = db.Close()
	return err
}

func ConnectDatabase(lg log.LoggerInterface, config *config.Config, debug bool) (*DBCloseCallback, *DatabaseOperations, error) {
	queryTimeout := config.Database.QueryDuration

	connection := config.Database.GetConnection(lg)

	gormConfig := gorm.Config{}
	gormConfig.DefaultTransactionTimeout = 5 * time.Second
	gormConfig.PrepareStmt = true
	gormConfig.TranslateError = true

	if debug {
		gormConfig.Logger = logger.Default.LogMode(logger.Error)
	} else {
		gormConfig.Logger = logger.Default.LogMode(logger.Silent)
	}

	db, err := gorm.Open(connection, &gormConfig)
	if err != nil {
		return nil, nil, Errorf("error occured while connecting to operation: %v", err)
	}

	if err = db.Migrator().AutoMigrate(&User{}, &FlightPlan{}, &History{}, &Activity{}, &ActivityATC{},
		&ActivityPilot{}, &ActivityFacility{}, &AuditLog{}, &ControllerRecord{}, &Ticket{}, &ControllerApplication{}); err != nil {
		return nil, nil, Errorf("error occured while migrating operation: %v", err)
	}

	dbPool, err := db.DB()
	if err != nil {
		return nil, nil, Errorf("error occured while creating operation pool: %v", err)
	}

	maxOpenConnections := config.Database.ServerMaxConnections * 4 / 5 // 不超过数据库最大连接的80%
	maxIdleConnections := maxOpenConnections / 5                       // 空闲连接约为最大连接的20%

	dbPool.SetMaxIdleConns(maxIdleConnections)
	dbPool.SetMaxOpenConns(maxOpenConnections)
	dbPool.SetConnMaxLifetime(config.Database.ConnectIdleDuration)

	err = dbPool.Ping()
	if err != nil {
		return nil, nil, Errorf("error occured while pinging operation: %v", err)
	}
	lg.Info("Database initialized and connection established")

	userOperation := NewUserOperation(lg, db, queryTimeout, config.Server.General)
	flightPlanOperation := NewFlightPlanOperation(lg, db, queryTimeout, config.Server.General)
	historyOperation := NewHistoryOperation(lg, db, queryTimeout)
	activityOperation := NewActivityOperation(lg, db, queryTimeout)
	auditLogOperation := NewAuditLogOperation(lg, db, queryTimeout)
	controllerOperation := NewControllerOperation(lg, db, queryTimeout)
	controllerRecordOperation := NewControllerRecordOperation(lg, db, queryTimeout)
	controllerApplicationOperation := NewControllerApplicationOperation(lg, db, queryTimeout)
	ticketOperation := NewTicketOperation(lg, db, queryTimeout)

	return NewDBCloseCallback(lg, db),
		NewDatabaseOperations(
			userOperation,
			flightPlanOperation,
			historyOperation,
			activityOperation,
			auditLogOperation,
			controllerOperation,
			controllerRecordOperation,
			controllerApplicationOperation,
			ticketOperation,
		),
		nil
}

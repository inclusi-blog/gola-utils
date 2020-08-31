package tracing

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gola-glitch/gola-utils/model"
	"github.com/stretchr/testify/suite"
	"testing"
)

type TracingTestSuite struct {
	suite.Suite
}

func TestTracingTestSuite(t *testing.T) {
	suite.Run(t, new(TracingTestSuite))
}

func (suite *TracingTestSuite) TearDownTest() {
	RegisterOracleDriverWithInstrumentationFunc = RegisterOracleDriverWithInstrumentation
	OpenDbConnectionFunc = OpenDbConnection
	SqlOpenFunc = sql.Open
}

func (suite *TracingTestSuite) TestOpenConnectionSuccessfully() {
	db, mock, _ := sqlmock.NewWithDSN("sqlmock_db_0", sqlmock.MonitorPingsOption(true))
	SqlOpenFunc = func(string2 string, string3 string) (*sql.DB, error) {
		return db, nil
	}
	mock.ExpectPing()
	db1, dbErr := OpenDbConnection("mock", "sqlmock_db_0")
	suite.Nil(dbErr)
	suite.NotNil(db1)
	db.Close()
}

func (suite *TracingTestSuite) TestOpenConnectionShouldReturnErrorWhenConnectionFails() {
	SqlOpenFunc = func(string2 string, string3 string) (*sql.DB, error) {
		return nil, errors.New("unable to open")
	}
	db1, dbErr := OpenDbConnection("mock", "sqlmock_db_0")
	suite.Nil(db1)
	suite.Equal(errors.New("unable to open"), dbErr)
}

func (suite *TracingTestSuite) TestOpenConnectionShouldReturnErrorWhenPingFails() {
	db, mock, _ := sqlmock.NewWithDSN("sqlmock_db_0", sqlmock.MonitorPingsOption(true))
	SqlOpenFunc = func(string2 string, string3 string) (*sql.DB, error) {
		return db, nil
	}
	mock.ExpectPing().WillReturnError(errors.New("unable to ping"))
	db1, dbErr := OpenDbConnection("mock", "sqlmock_db_0")
	suite.Nil(db1)
	suite.Equal(errors.New("unable to ping"), dbErr)
	db.Close()
}

func (suite *TracingTestSuite) TestRegisterOracleDriverWithInstrumentation() {
	db, mock, _ := sqlmock.NewWithDSN("sqlmock_db_0", sqlmock.MonitorPingsOption(true))
	SqlOpenFunc = func(string2 string, string3 string) (*sql.DB, error) {
		return db, nil
	}
	mock.ExpectClose()
	dbErr := RegisterOracleDriverWithInstrumentation("mock", "sqlmock_db_0")
	suite.Nil(dbErr)
}

func (suite *TracingTestSuite) TestRegisterOracleDriverWithInstrumentationShouldFailWhenUnableToOpenConnection() {
	err := errors.New("unable to open connection")
	SqlOpenFunc = func(string2 string, string3 string) (*sql.DB, error) {
		return nil, err
	}
	dbErr := RegisterOracleDriverWithInstrumentation("mock", "sqlmock_db_0")
	suite.Equal(err, dbErr)
}

func (suite *TracingTestSuite) TestRegisterOracleDriverWithInstrumentationShouldFailWhenUnableToCloseConnection() {
	db, mock, _ := sqlmock.NewWithDSN("sqlmock_db_0", sqlmock.MonitorPingsOption(true))
	SqlOpenFunc = func(string2 string, string3 string) (*sql.DB, error) {
		return db, nil
	}
	err := errors.New("unable to close connection")
	mock.ExpectClose().WillReturnError(err)
	dbErr := RegisterOracleDriverWithInstrumentation("mock", "sqlmock_db_0")
	suite.Equal(err, dbErr)
	db.Close()
}

func (suite *TracingTestSuite) TestInitSqlxOracleDBWithInstrumentationAndConnectionConfigSuccessfully() {
	db, _, _ := sqlmock.NewWithDSN("sqlmock_db_0", sqlmock.MonitorPingsOption(true))
	RegisterOracleDriverWithInstrumentationFunc = func(string2 string, string3 string) error {
		return nil
	}
	OpenDbConnectionFunc = func(string2 string, string3 string) (*sql.DB, error) {
		return db, nil
	}
	dbConn, err := InitSqlxOracleDBWithInstrumentationAndConnectionConfig("mock",
		"sqlmock_db_0", model.DBConnectionPoolConfig{
			MaxOpenConnections:             1,
			MaxIdleConnections:             2,
			MaxConnectionLifetimeInMinutes: 3,
		})
	suite.Nil(err)
	suite.NotNil(dbConn)
	db.Close()
	dbConn.Close()
}

func (suite *TracingTestSuite) TestInitSqlxOracleDBWithInstrumentationAndConnectionConfigShouldGiveErrorWhenUnableToRegisterOracleDriverWithInstrumentation() {
	dbErr := errors.New("unable to register")
	RegisterOracleDriverWithInstrumentationFunc = func(string2 string, string3 string) error {
		return dbErr
	}
	dbConn, err := InitSqlxOracleDBWithInstrumentationAndConnectionConfig("mock",
		"sqlmock_db_0", model.DBConnectionPoolConfig{})
	suite.Nil(dbConn)
	suite.Equal(dbErr, err)
}

func (suite *TracingTestSuite) TestInitSqlxOracleDBWithInstrumentationAndConnectionConfigShouldGiveErrorWhenUnableToOpenDBConnection() {
	dbErr := errors.New("unable to open db connection")
	RegisterOracleDriverWithInstrumentationFunc = func(string2 string, string3 string) error {
		return nil
	}
	OpenDbConnectionFunc = func(string2 string, string3 string) (*sql.DB, error) {
		return nil, dbErr
	}
	dbConn, err := InitSqlxOracleDBWithInstrumentationAndConnectionConfig("mock",
		"sqlmock_db_0", model.DBConnectionPoolConfig{})
	suite.Equal(dbErr, err)
	suite.Nil(dbConn)
}

func (suite *TracingTestSuite) TestInitSqlxOracleDBWithInstrumentationSuccessfully() {
	db, _, _ := sqlmock.NewWithDSN("sqlmock_db_0", sqlmock.MonitorPingsOption(true))
	RegisterOracleDriverWithInstrumentationFunc = func(string2 string, string3 string) error {
		return nil
	}
	OpenDbConnectionFunc = func(string2 string, string3 string) (*sql.DB, error) {
		return db, nil
	}
	dbConn, err := InitSqlxOracleDBWithInstrumentation("mock", "sqlmock_db_0")
	suite.Nil(err)
	suite.NotNil(dbConn)
	db.Close()
	dbConn.Close()
}

func (suite *TracingTestSuite) TestInitSqlxOracleDBWithInstrumentationShouldGiveErrorWhenUnableToRegisterOracleDriverWithInstrumentation() {
	dbErr := errors.New("unable to register")
	RegisterOracleDriverWithInstrumentationFunc = func(string2 string, string3 string) error {
		return dbErr
	}
	dbConn, err := InitSqlxOracleDBWithInstrumentation("mock", "sqlmock_db_0")
	suite.Nil(dbConn)
	suite.Equal(dbErr, err)
}

func (suite *TracingTestSuite) TestInitSqlxOracleDBWithInstrumentationShouldGiveErrorWhenUnableToOpenDBConnection() {
	dbErr := errors.New("unable to open db connection")
	RegisterOracleDriverWithInstrumentationFunc = func(string2 string, string3 string) error {
		return nil
	}
	OpenDbConnectionFunc = func(string2 string, string3 string) (*sql.DB, error) {
		return nil, dbErr
	}
	dbConn, err := InitSqlxOracleDBWithInstrumentation("mock", "sqlmock_db_0")
	suite.Equal(dbErr, err)
	suite.Nil(dbConn)
}

func (suite *TracingTestSuite) TestInitSqlOracleDBWithInstrumentationShouldOpenConnectionSuccessfully() {
	db, _, _ := sqlmock.NewWithDSN("sqlmock_db_0", sqlmock.MonitorPingsOption(true))
	RegisterOracleDriverWithInstrumentationFunc = func(string2 string, string3 string) error {
		return nil
	}
	OpenDbConnectionFunc = func(string2 string, string3 string) (*sql.DB, error) {
		return db, nil
	}
	dbConn, err := InitSqlOracleDBWithInstrumentation("mock", "sqlmock_db_0")
	suite.Nil(err)
	suite.NotNil(dbConn)
	db.Close()
	dbConn.Close()
}

func (suite *TracingTestSuite) TestInitSqlOracleDBWithInstrumentationShouldReturnErrorWhenRegisterTracingFails() {
	dbErr := errors.New("unable to register")
	RegisterOracleDriverWithInstrumentationFunc = func(string2 string, string3 string) error {
		return dbErr
	}
	dbConn, err := InitSqlOracleDBWithInstrumentation("mock", "sqlmock_db_0")
	suite.Equal(dbErr, err)
	suite.Nil(dbConn)
}

func (suite *TracingTestSuite) TestInitSqlOracleDBWithInstrumentationShouldReturnErrorWhenOpenDBConnectionFails() {
	dbErr := errors.New("unable to open db connection")
	RegisterOracleDriverWithInstrumentationFunc = func(string2 string, string3 string) error {
		return nil
	}
	OpenDbConnectionFunc = func(string2 string, string3 string) (*sql.DB, error) {
		return nil, dbErr
	}
	dbConn, err := InitSqlOracleDBWithInstrumentation("mock", "sqlmock_db_0")
	suite.Equal(dbErr, err)
	suite.Nil(dbConn)
}

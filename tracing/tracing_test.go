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
	RegisterPostgresDriverWithInstrumentationFunc = RegisterPostgresDriverWithInstrumentation
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

func (suite *TracingTestSuite) TestRegisterPostgresDriverWithInstrumentation() {
	db, mock, _ := sqlmock.NewWithDSN("sqlmock_db_0", sqlmock.MonitorPingsOption(true))
	SqlOpenFunc = func(string2 string, string3 string) (*sql.DB, error) {
		return db, nil
	}
	mock.ExpectClose()
	dbErr := RegisterPostgresDriverWithInstrumentation("mock", "sqlmock_db_0")
	suite.Nil(dbErr)
}

func (suite *TracingTestSuite) TestRegisterPostgresDriverWithInstrumentationShouldFailWhenUnableToOpenConnection() {
	err := errors.New("unable to open connection")
	SqlOpenFunc = func(string2 string, string3 string) (*sql.DB, error) {
		return nil, err
	}
	dbErr := RegisterPostgresDriverWithInstrumentation("mock", "sqlmock_db_0")
	suite.Equal(err, dbErr)
}

func (suite *TracingTestSuite) TestRegisterPostgresDriverWithInstrumentationShouldFailWhenUnableToCloseConnection() {
	db, mock, _ := sqlmock.NewWithDSN("sqlmock_db_0", sqlmock.MonitorPingsOption(true))
	SqlOpenFunc = func(string2 string, string3 string) (*sql.DB, error) {
		return db, nil
	}
	err := errors.New("unable to close connection")
	mock.ExpectClose().WillReturnError(err)
	dbErr := RegisterPostgresDriverWithInstrumentation("mock", "sqlmock_db_0")
	suite.Equal(err, dbErr)
	db.Close()
}

func (suite *TracingTestSuite) TestInitPostgresDBWithInstrumentationAndConnectionConfigSuccessfully() {
	db, _, _ := sqlmock.NewWithDSN("sqlmock_db_0", sqlmock.MonitorPingsOption(true))
	RegisterPostgresDriverWithInstrumentationFunc = func(string2 string, string3 string) error {
		return nil
	}
	OpenDbConnectionFunc = func(string2 string, string3 string) (*sql.DB, error) {
		return db, nil
	}
	dbConn, err := InitPostgresDBWithInstrumentationAndConnectionConfig("mock",
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

func (suite *TracingTestSuite) TestInitPostgresDBWithInstrumentationAndConnectionConfigShouldGiveErrorWhenUnableToRegisterPostgresDriverWithInstrumentation() {
	dbErr := errors.New("unable to register")
	RegisterPostgresDriverWithInstrumentationFunc = func(string2 string, string3 string) error {
		return dbErr
	}
	dbConn, err := InitPostgresDBWithInstrumentationAndConnectionConfig("mock",
		"sqlmock_db_0", model.DBConnectionPoolConfig{})
	suite.Nil(dbConn)
	suite.Equal(dbErr, err)
}

func (suite *TracingTestSuite) TestInitPostgresDBWithInstrumentationAndConnectionConfigShouldGiveErrorWhenUnableToOpenDBConnection() {
	dbErr := errors.New("unable to open db connection")
	RegisterPostgresDriverWithInstrumentationFunc = func(string2 string, string3 string) error {
		return nil
	}
	OpenDbConnectionFunc = func(string2 string, string3 string) (*sql.DB, error) {
		return nil, dbErr
	}
	dbConn, err := InitPostgresDBWithInstrumentationAndConnectionConfig("mock",
		"sqlmock_db_0", model.DBConnectionPoolConfig{})
	suite.Equal(dbErr, err)
	suite.Nil(dbConn)
}

func (suite *TracingTestSuite) TestInitPostgresDBWithInstrumentationSuccessfully() {
	db, _, _ := sqlmock.NewWithDSN("sqlmock_db_0", sqlmock.MonitorPingsOption(true))
	RegisterPostgresDriverWithInstrumentationFunc = func(string2 string, string3 string) error {
		return nil
	}
	OpenDbConnectionFunc = func(string2 string, string3 string) (*sql.DB, error) {
		return db, nil
	}
	dbConn, err := InitPostgresDBWithInstrumentation("mock", "sqlmock_db_0")
	suite.Nil(err)
	suite.NotNil(dbConn)
	db.Close()
	dbConn.Close()
}

func (suite *TracingTestSuite) TestInitPostgresDBWithInstrumentationShouldGiveErrorWhenUnableToRegisterPostgresDriverWithInstrumentation() {
	dbErr := errors.New("unable to register")
	RegisterPostgresDriverWithInstrumentationFunc = func(string2 string, string3 string) error {
		return dbErr
	}
	dbConn, err := InitPostgresDBWithInstrumentation("mock", "sqlmock_db_0")
	suite.Nil(dbConn)
	suite.Equal(dbErr, err)
}

func (suite *TracingTestSuite) TestInitPostgresDBWithInstrumentationShouldGiveErrorWhenUnableToOpenDBConnection() {
	dbErr := errors.New("unable to open db connection")
	RegisterPostgresDriverWithInstrumentationFunc = func(string2 string, string3 string) error {
		return nil
	}
	OpenDbConnectionFunc = func(string2 string, string3 string) (*sql.DB, error) {
		return nil, dbErr
	}
	dbConn, err := InitPostgresDBWithInstrumentation("mock", "sqlmock_db_0")
	suite.Equal(dbErr, err)
	suite.Nil(dbConn)
}

func (suite *TracingTestSuite) TestInitSqlPostgresDBWithInstrumentationShouldOpenConnectionSuccessfully() {
	db, _, _ := sqlmock.NewWithDSN("sqlmock_db_0", sqlmock.MonitorPingsOption(true))
	RegisterPostgresDriverWithInstrumentationFunc = func(string2 string, string3 string) error {
		return nil
	}
	OpenDbConnectionFunc = func(string2 string, string3 string) (*sql.DB, error) {
		return db, nil
	}
	dbConn, err := InitSqlPostgresDBWithInstrumentation("mock", "sqlmock_db_0")
	suite.Nil(err)
	suite.NotNil(dbConn)
	db.Close()
	dbConn.Close()
}

func (suite *TracingTestSuite) TestInitSqlPostgresDBWithInstrumentationShouldReturnErrorWhenRegisterTracingFails() {
	dbErr := errors.New("unable to register")
	RegisterPostgresDriverWithInstrumentationFunc = func(string2 string, string3 string) error {
		return dbErr
	}
	dbConn, err := InitSqlPostgresDBWithInstrumentation("mock", "sqlmock_db_0")
	suite.Equal(dbErr, err)
	suite.Nil(dbConn)
}

func (suite *TracingTestSuite) TestInitSqlPostgresDBWithInstrumentationShouldReturnErrorWhenOpenDBConnectionFails() {
	dbErr := errors.New("unable to open db connection")
	RegisterPostgresDriverWithInstrumentationFunc = func(string2 string, string3 string) error {
		return nil
	}
	OpenDbConnectionFunc = func(string2 string, string3 string) (*sql.DB, error) {
		return nil, dbErr
	}
	dbConn, err := InitSqlPostgresDBWithInstrumentation("mock", "sqlmock_db_0")
	suite.Equal(dbErr, err)
	suite.Nil(dbConn)
}

package tracing

import (
	"database/sql"
	"fmt"
	"github.com/gola-glitch/gola-utils/constants"
	"github.com/gola-glitch/gola-utils/logging"
	"github.com/gola-glitch/gola-utils/model"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"

	"contrib.go.opencensus.io/exporter/ocagent"
	"contrib.go.opencensus.io/integrations/ocsql"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
)

var (
	RegisterOracleDriverWithInstrumentationFunc = RegisterOracleDriverWithInstrumentation
	OpenDbConnectionFunc                        = OpenDbConnection
	SqlOpenFunc                                 = sql.Open
)

func Init(serviceName string, ocAgent string) *ocagent.Exporter {
	oce, _ := ocagent.NewExporter(
		ocagent.WithInsecure(),
		ocagent.WithReconnectionPeriod(1*time.Second),
		ocagent.WithAddress(ocAgent),
		ocagent.WithServiceName(serviceName),
		ocagent.WithSpanConfig(ocagent.SpanConfig{
			AnnotationEventsPerSpan: constants.TRACE_CONFIG_MAX_ANNOTATIONS,
		}))
	trace.RegisterExporter(oce)
	trace.ApplyConfig(trace.Config{
		DefaultSampler:             trace.AlwaysSample(),
		MaxAnnotationEventsPerSpan: constants.TRACE_CONFIG_MAX_ANNOTATIONS,
	})

	return oce
}

func WithTracing(app http.Handler, healthz string) http.Handler {
	return &ochttp.Handler{
		Handler: app,
		GetStartOptions: func(r *http.Request) trace.StartOptions {
			startOptions := trace.StartOptions{}

			if r.URL.Path == healthz {
				startOptions.Sampler = trace.NeverSample()
			}
			return startOptions
		}}
}

func InitSqlOracleDBWithInstrumentation(driverName, conn string) (*sql.DB, error) {
	logger := logging.NewLoggerEntry()
	if err := RegisterOracleDriverWithInstrumentationFunc(driverName, conn); err != nil {
		logger.Error("error while registering driver with tracing instrumentation ", err.Error())
		return nil, err
	}
	db, err := OpenDbConnectionFunc(driverName, conn)
	if err != nil {
		logger.Error("error while opening DB connection", err.Error())
		return nil, err
	}
	return db, nil
}

func RegisterOracleDriverWithInstrumentation(driverName, conn string) error {
	logger := logging.NewLoggerEntry()
	db, dbErr := SqlOpenFunc(driverName, conn)
	if dbErr != nil {
		logger.Error("Unable connect to db to fetch driver", dbErr.Error())
		return dbErr
	}
	addTracingOptions(db, GetInstrumentedDriverName(driverName))
	dbErr = db.Close()
	if dbErr != nil {
		logger.Error("Unable to close db connection ", dbErr.Error())
		return dbErr
	}
	return nil
}

func GetInstrumentedDriverName(driverName string) string {
	return fmt.Sprintf("%s%s", driverName, "-with-oc")
}

func OpenDbConnection(driverName, conn string) (*sql.DB, error) {
	logger := logging.NewLoggerEntry()
	db, dbErr := SqlOpenFunc(GetInstrumentedDriverName(driverName), conn)
	if dbErr != nil {
		logger.Error("Unable to connect to Gola DB ", dbErr.Error())
		return nil, dbErr
	}
	if dbErr = db.Ping(); dbErr != nil {
		logger.Error("Unable to ping ", dbErr.Error())
		return nil, dbErr
	}
	return db, nil
}

func addTracingOptions(db *sql.DB, driverName string) {
	var found = false
	for _, name := range sql.Drivers() {
		if name == driverName {
			found = true
		}
	}
	if !found {
		traceOptions := ocsql.TraceOptions{AllowRoot: false, Ping: true, RowsNext: false,
			RowsClose: false, RowsAffected: true, LastInsertID: false, Query: true, QueryParams: false}
		driver := ocsql.Wrap(db.Driver(), ocsql.WithOptions(traceOptions))
		sql.Register(driverName, driver)
	}
}

// can be used for default custom config of oracle
func InitSqlxOracleDBWithInstrumentation(driverName, conn string) (*sqlx.DB, error) {
	err := RegisterOracleDriverWithInstrumentationFunc(driverName, conn)
	if err != nil {
		return nil, err
	}
	db, dbErr := OpenDbConnectionFunc(driverName, conn)
	if dbErr != nil {
		return nil, dbErr
	}
	return sqlx.NewDb(db, driverName), nil
}

func InitSqlxOracleDBWithInstrumentationAndConnectionConfig(driverName, conn string, dbConnectionPoolConfig model.DBConnectionPoolConfig) (*sqlx.DB, error) {
	err := RegisterOracleDriverWithInstrumentationFunc(driverName, conn)
	if err != nil {
		return nil, err
	}
	db, dbErr := OpenDbConnectionFunc(driverName, conn)
	if dbErr != nil {
		return nil, dbErr
	}
	newDb := sqlx.NewDb(db, driverName)
	setConnectionPoolConfig(dbConnectionPoolConfig, newDb)
	return newDb, nil
}

func setConnectionPoolConfig(config model.DBConnectionPoolConfig, newDb *sqlx.DB) {
	if config.MaxIdleConnections > 0 {
		newDb.SetMaxIdleConns(config.MaxIdleConnections)
	}
	if config.MaxConnectionLifetimeInMinutes > 0 {
		newDb.SetConnMaxLifetime(time.Duration(config.MaxConnectionLifetimeInMinutes) * time.Minute)
	}
	if config.MaxOpenConnections > 0 {
		newDb.SetMaxOpenConns(config.MaxOpenConnections)
	}
}

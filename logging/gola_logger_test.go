package logging

import (
	"context"
	"errors"
	"fmt"
	"github.com/inclusi-blog/gola-utils/constants"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	logrusTest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/suite"
	"go.opencensus.io/trace"
)

type GolaLoggerTestSuite struct {
	suite.Suite
	context  *gin.Context
	recorder *httptest.ResponseRecorder
}

func TestGolaLoggerTestSuite(t *testing.T) {
	suite.Run(t, new(GolaLoggerTestSuite))
}

func (suite *GolaLoggerTestSuite) SetupTest() {
	suite.recorder = httptest.NewRecorder()
	suite.context, _ = gin.CreateTestContext(suite.recorder)
}

func (suite *GolaLoggerTestSuite) TearDownTest() {
	suite.recorder = nil
	suite.context = nil
}

func (suite *GolaLoggerTestSuite) TestShouldCreatePrivateNewLoggerEntrySuccessfully() {
	logrusEntry := logrus.StandardLogger().WithContext(context.Background())
	golaLoggerEntry := newLoggerEntry(context.Background(), logrusEntry)
	suite.Equal(context.Background(), golaLoggerEntry.context)
	suite.Equal(logrusEntry, golaLoggerEntry.stdEntry)
	suite.Equal(0, len(golaLoggerEntry.Data))
	suite.Equal("", golaLoggerEntry.err)
}

func (suite *GolaLoggerTestSuite) TestShouldCreatePublicNewLoggerEntrySuccessfully() {
	golaLoggerEntry := NewLoggerEntry()
	suite.Equal(context.TODO(), golaLoggerEntry.context)
	suite.Equal(logrus.StandardLogger().WithContext(context.TODO()), golaLoggerEntry.stdEntry)
	suite.Equal(0, len(golaLoggerEntry.Data))
	suite.Equal("", golaLoggerEntry.err)
}

func (suite *GolaLoggerTestSuite) TestShouldSetWithErrorSuccessfully() {
	err := errors.New("dummyError")
	golaLoggerEntry := NewLoggerEntry().WithError(err)
	suite.Equal(context.TODO(), golaLoggerEntry.context)
	suite.Equal(logrus.StandardLogger().WithContext(context.TODO()).WithError(err), golaLoggerEntry.stdEntry)
	suite.Equal(1, len(golaLoggerEntry.Data))
	suite.Equal("", golaLoggerEntry.err)
	suite.Equal(err, golaLoggerEntry.Data["error"])
}

func (suite *GolaLoggerTestSuite) TestShouldSetWithFieldSuccessfully() {
	golaLoggerEntry := NewLoggerEntry().WithField("key1", "value1")
	suite.Equal(context.TODO(), golaLoggerEntry.context)
	suite.Equal(logrus.StandardLogger().WithContext(context.TODO()).WithField("key1", "value1"), golaLoggerEntry.stdEntry)
	suite.Equal(1, len(golaLoggerEntry.Data))
	suite.Equal("", golaLoggerEntry.err)
	suite.Equal("value1", golaLoggerEntry.Data["key1"])
}

func (suite *GolaLoggerTestSuite) TestShouldSetWithFieldsSuccessfully() {
	expectedLogrusEntry := logrus.StandardLogger().WithContext(context.TODO()).WithFields(logrus.Fields{"key1": "value1", "key2": "value2"})
	golaLoggerEntry := NewLoggerEntry().WithFields(GolaFields{"key1": "value1", "key2": "value2"})
	suite.Equal(context.TODO(), golaLoggerEntry.context)
	suite.Equal(expectedLogrusEntry, golaLoggerEntry.stdEntry)
	suite.Equal(2, len(golaLoggerEntry.Data))
	suite.Equal("", golaLoggerEntry.err)
	suite.Equal("value1", golaLoggerEntry.Data["key1"])
	suite.Equal("value2", golaLoggerEntry.Data["key2"])
}

//TODO: Add rest cases for WithFields

func (suite *GolaLoggerTestSuite) TestShouldSetWithContextSuccessfully() {
	expectedLogrusEntry := logrus.StandardLogger().WithContext(context.Background())
	golaLoggerEntry := newLoggerEntry(context.Background(), expectedLogrusEntry)
	suite.Equal(context.Background(), golaLoggerEntry.context)
	golaLoggerEntry = golaLoggerEntry.WithContext(context.TODO())
	suite.Equal(context.TODO(), golaLoggerEntry.context)
	suite.Equal(expectedLogrusEntry, golaLoggerEntry.stdEntry)
	suite.Equal(0, len(golaLoggerEntry.Data))
	suite.Equal("", golaLoggerEntry.err)
}

func (suite *GolaLoggerTestSuite) TestShouldSetFormatterSuccessfully() {
	golaLoggerEntry := newLoggerEntry(context.TODO(), logrus.StandardLogger().WithContext(context.TODO()))
	golaLoggerEntry.SetFormatter("json")
	suite.Equal(&logrus.JSONFormatter{}, golaLoggerEntry.stdEntry.Logger.Formatter)

	suite.Equal(context.TODO(), golaLoggerEntry.context)
	suite.Equal(0, len(golaLoggerEntry.Data))
	suite.Equal("", golaLoggerEntry.err)
}

//TODO: testSetFormatter

func (suite *GolaLoggerTestSuite) TestShouldSetPanicLevelSuccessfully() {
	expectedLogrusEntry := logrus.StandardLogger().WithContext(context.Background())
	golaLoggerEntry := newLoggerEntry(context.Background(), expectedLogrusEntry)
	err := golaLoggerEntry.SetLevel("panic")

	suite.Nil(err)
	suite.Equal(context.TODO(), golaLoggerEntry.context)
	suite.Equal(logrus.PanicLevel, golaLoggerEntry.stdEntry.Logger.Level)
	suite.Equal(expectedLogrusEntry, golaLoggerEntry.stdEntry)
	suite.Equal(0, len(golaLoggerEntry.Data))
	suite.Equal("", golaLoggerEntry.err)
}

func (suite *GolaLoggerTestSuite) TestShouldSetFatalLevelSuccessfully() {
	expectedLogrusEntry := logrus.StandardLogger().WithContext(context.Background())
	golaLoggerEntry := newLoggerEntry(context.Background(), expectedLogrusEntry)
	err := golaLoggerEntry.SetLevel("Fatal")

	suite.Nil(err)
	suite.Equal(context.TODO(), golaLoggerEntry.context)
	suite.Equal(logrus.FatalLevel, golaLoggerEntry.stdEntry.Logger.Level)
	suite.Equal(expectedLogrusEntry, golaLoggerEntry.stdEntry)
	suite.Equal(0, len(golaLoggerEntry.Data))
	suite.Equal("", golaLoggerEntry.err)
}

func (suite *GolaLoggerTestSuite) TestShouldSetErrorLevelSuccessfully() {
	expectedLogrusEntry := logrus.StandardLogger().WithContext(context.TODO())
	golaLoggerEntry := newLoggerEntry(context.Background(), expectedLogrusEntry)
	err := golaLoggerEntry.SetLevel("ERROR")

	suite.Nil(err)
	suite.Equal(context.TODO(), golaLoggerEntry.context)
	suite.Equal(logrus.ErrorLevel, golaLoggerEntry.stdEntry.Logger.Level)
	suite.Equal(expectedLogrusEntry, golaLoggerEntry.stdEntry)
	suite.Equal(0, len(golaLoggerEntry.Data))
	suite.Equal("", golaLoggerEntry.err)
}

func (suite *GolaLoggerTestSuite) TestShouldSetWarnLevelSuccessfully() {
	expectedLogrusEntry := logrus.StandardLogger().WithContext(context.Background())
	golaLoggerEntry := newLoggerEntry(context.Background(), expectedLogrusEntry)
	err := golaLoggerEntry.SetLevel("warn")

	suite.Nil(err)
	suite.Equal(context.TODO(), golaLoggerEntry.context)
	suite.Equal(logrus.WarnLevel, golaLoggerEntry.stdEntry.Logger.Level)
	suite.Equal(expectedLogrusEntry, golaLoggerEntry.stdEntry)
	suite.Equal(0, len(golaLoggerEntry.Data))
	suite.Equal("", golaLoggerEntry.err)
}

func (suite *GolaLoggerTestSuite) TestShouldSetWarningLevelSuccessfully() {
	expectedLogrusEntry := logrus.StandardLogger().WithContext(context.Background())
	golaLoggerEntry := newLoggerEntry(context.Background(), expectedLogrusEntry)
	err := golaLoggerEntry.SetLevel("warning")

	suite.Nil(err)
	suite.Equal(context.TODO(), golaLoggerEntry.context)
	suite.Equal(logrus.WarnLevel, golaLoggerEntry.stdEntry.Logger.Level)
	suite.Equal(expectedLogrusEntry, golaLoggerEntry.stdEntry)
	suite.Equal(0, len(golaLoggerEntry.Data))
	suite.Equal("", golaLoggerEntry.err)
}

func (suite *GolaLoggerTestSuite) TestShouldSetInfoLevelSuccessfully() {
	expectedLogrusEntry := logrus.StandardLogger().WithContext(context.Background())
	golaLoggerEntry := newLoggerEntry(context.Background(), expectedLogrusEntry)
	err := golaLoggerEntry.SetLevel("info")

	suite.Nil(err)
	suite.Equal(context.TODO(), golaLoggerEntry.context)
	suite.Equal(logrus.InfoLevel, golaLoggerEntry.stdEntry.Logger.Level)
	suite.Equal(expectedLogrusEntry, golaLoggerEntry.stdEntry)
	suite.Equal(0, len(golaLoggerEntry.Data))
	suite.Equal("", golaLoggerEntry.err)
}

func (suite *GolaLoggerTestSuite) TestShouldSetPrintLevelSuccessfully() {
	expectedLogrusEntry := logrus.StandardLogger().WithContext(context.Background())
	golaLoggerEntry := newLoggerEntry(context.Background(), expectedLogrusEntry)
	err := golaLoggerEntry.SetLevel("print")

	suite.Nil(err)
	suite.Equal(context.TODO(), golaLoggerEntry.context)
	suite.Equal(logrus.InfoLevel, golaLoggerEntry.stdEntry.Logger.Level)
	suite.Equal(expectedLogrusEntry, golaLoggerEntry.stdEntry)
	suite.Equal(0, len(golaLoggerEntry.Data))
	suite.Equal("", golaLoggerEntry.err)
}

func (suite *GolaLoggerTestSuite) TestShouldSetDebugLevelSuccessfully() {
	expectedLogrusEntry := logrus.StandardLogger().WithContext(context.Background())
	golaLoggerEntry := newLoggerEntry(context.Background(), expectedLogrusEntry)
	err := golaLoggerEntry.SetLevel("DEBUG")

	suite.Nil(err)
	suite.Equal(context.TODO(), golaLoggerEntry.context)
	suite.Equal(logrus.DebugLevel, golaLoggerEntry.stdEntry.Logger.Level)
	suite.Equal(expectedLogrusEntry, golaLoggerEntry.stdEntry)
	suite.Equal(0, len(golaLoggerEntry.Data))
	suite.Equal("", golaLoggerEntry.err)
}

func (suite *GolaLoggerTestSuite) TestShouldSetTraceLevelSuccessfully() {
	expectedLogrusEntry := logrus.StandardLogger().WithContext(context.Background())
	golaLoggerEntry := newLoggerEntry(context.Background(), expectedLogrusEntry)
	err := golaLoggerEntry.SetLevel("traCe")

	suite.Nil(err)
	suite.Equal(context.TODO(), golaLoggerEntry.context)
	suite.Equal(logrus.TraceLevel, golaLoggerEntry.stdEntry.Logger.Level)
	suite.Equal(expectedLogrusEntry, golaLoggerEntry.stdEntry)
	suite.Equal(0, len(golaLoggerEntry.Data))
	suite.Equal("", golaLoggerEntry.err)
}

func (suite *GolaLoggerTestSuite) TestSetLevelShouldThrowErrorWhenInvalidLevelStringProvided() {
	expectedLogrusEntry := logrus.StandardLogger().WithContext(context.Background())
	golaLoggerEntry := newLoggerEntry(context.Background(), expectedLogrusEntry)
	err := golaLoggerEntry.SetLevel("dummy")

	suite.NotNil(err)
	suite.Equal("not a valid golaLogger Level: \"dummy\"", err.Error())
	suite.Equal(logrus.InfoLevel, golaLoggerEntry.stdEntry.Logger.Level)
}

type testExporter struct {
	SpanData *trace.SpanData
}

func (t *testExporter) ExportSpan(s *trace.SpanData) {
	t.SpanData = s
}

func (suite *GolaLoggerTestSuite) TestShouldPrintTracefLogsInLogrusAndSpanSuccessfully() {
	logger, hook := logrusTest.NewNullLogger()
	t := testExporter{}
	trace.RegisterExporter(&t)
	ctx, s := trace.StartSpan(suite.context, "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))

	golaLoggerEntry := newLoggerEntry(ctx, logger.WithContext(context.Background()))
	_ = golaLoggerEntry.SetLevel("trace")
	golaLoggerEntry.Tracef("%d", 100)
	s.End()

	suite.Equal("100", hook.LastEntry().Message)
	suite.Equal(logrus.TraceLevel, hook.LastEntry().Level)
	suite.Equal("trace", t.SpanData.Annotations[0].Attributes["level"])
	suite.Equal("100", t.SpanData.Annotations[0].Message)
	hook.Reset()
}

func (suite *GolaLoggerTestSuite) TestShouldPrintDebugfLogsInLogrusAndSpanSuccessfully() {
	logger, hook := logrusTest.NewNullLogger()
	t := testExporter{}
	trace.RegisterExporter(&t)
	ctx, s := trace.StartSpan(suite.context, "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))

	golaLoggerEntry := newLoggerEntry(ctx, logger.WithContext(context.Background()))
	golaLoggerEntry.SetLevel("debug")
	golaLoggerEntry.Debugf("%d", 100)
	s.End()

	suite.Equal("100", hook.LastEntry().Message)
	suite.Equal(logrus.DebugLevel, hook.LastEntry().Level)
	suite.Equal("debug", t.SpanData.Annotations[0].Attributes["level"])
	suite.Equal("100", t.SpanData.Annotations[0].Message)
	hook.Reset()
}

func (suite *GolaLoggerTestSuite) TestShouldPrintInfofLogsInLogrusAndSpanSuccessfully() {
	logger, hook := logrusTest.NewNullLogger()
	t := testExporter{}
	trace.RegisterExporter(&t)
	ctx, s := trace.StartSpan(suite.context, "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))

	golaLoggerEntry := newLoggerEntry(ctx, logger.WithContext(context.Background()))
	golaLoggerEntry.Infof("%t", true)
	s.End()

	suite.Equal("true", hook.LastEntry().Message)
	suite.Equal(logrus.InfoLevel, hook.LastEntry().Level)
	suite.Equal("info", t.SpanData.Annotations[0].Attributes["level"])
	suite.Equal("true", t.SpanData.Annotations[0].Message)
	hook.Reset()
}

func (suite *GolaLoggerTestSuite) TestShouldPrintPrintfLogsInLogrusAndSpanSuccessfully() {
	logger, hook := logrusTest.NewNullLogger()
	t := testExporter{}
	trace.RegisterExporter(&t)
	ctx, s := trace.StartSpan(suite.context, "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))

	golaLoggerEntry := newLoggerEntry(ctx, logger.WithContext(context.Background()))
	golaLoggerEntry.Printf("%s", "hi")
	s.End()

	suite.Equal("hi", hook.LastEntry().Message)
	suite.Equal(logrus.InfoLevel, hook.LastEntry().Level)
	suite.Equal("info", t.SpanData.Annotations[0].Attributes["level"])
	suite.Equal("hi", t.SpanData.Annotations[0].Message)
	hook.Reset()
}

func (suite *GolaLoggerTestSuite) TestShouldPrintWarnfLogsInLogrusAndSpanSuccessfully() {
	logger, hook := logrusTest.NewNullLogger()
	t := testExporter{}
	trace.RegisterExporter(&t)
	ctx, s := trace.StartSpan(suite.context, "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))

	golaLoggerEntry := newLoggerEntry(ctx, logger.WithContext(context.Background()))
	golaLoggerEntry.Warnf("%s", "hi")
	s.End()

	suite.Equal("hi", hook.LastEntry().Message)
	suite.Equal(logrus.WarnLevel, hook.LastEntry().Level)
	suite.Equal("warning", t.SpanData.Annotations[0].Attributes["level"])
	suite.Equal("hi", t.SpanData.Annotations[0].Message)
	hook.Reset()
}

func (suite *GolaLoggerTestSuite) TestShouldPrintWarningfLogsInLogrusAndSpanSuccessfully() {
	logger, hook := logrusTest.NewNullLogger()
	t := testExporter{}
	trace.RegisterExporter(&t)
	ctx, s := trace.StartSpan(suite.context, "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))

	golaLoggerEntry := newLoggerEntry(ctx, logger.WithContext(context.Background()))
	golaLoggerEntry.Warningf("%s", "hi")
	s.End()

	suite.Equal("hi", hook.LastEntry().Message)
	suite.Equal(logrus.WarnLevel, hook.LastEntry().Level)
	suite.Equal("warning", t.SpanData.Annotations[0].Attributes["level"])
	suite.Equal("hi", t.SpanData.Annotations[0].Message)
	hook.Reset()
}

func (suite *GolaLoggerTestSuite) TestShouldPrintErrorfLogsInLogrusAndSpanSuccessfully() {
	logger, hook := logrusTest.NewNullLogger()
	t := testExporter{}
	trace.RegisterExporter(&t)
	ctx, s := trace.StartSpan(suite.context, "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))

	golaLoggerEntry := newLoggerEntry(ctx, logger.WithContext(context.Background()))
	golaLoggerEntry.Errorf("%s", "hi")
	s.End()

	suite.Equal("hi", hook.LastEntry().Message)
	suite.Equal(logrus.ErrorLevel, hook.LastEntry().Level)
	suite.Equal("error", t.SpanData.Annotations[0].Attributes["level"])
	suite.Equal("hi", t.SpanData.Annotations[0].Message)
	hook.Reset()
}

//func (suite *GolaLoggerTestSuite) TestShouldPrintFatalfLogsInLogrusAndSpanSuccessfully() {
//	logger, hook := logrusTest.NewNullLogger()
//	handler := func() {
//
//		recover()
//
//	}
//	logrus.RegisterExitHandler(handler)
//	t := testExporter{}
//	trace.RegisterExporter(&t)
//	ctx, s := trace.StartSpan(suite.context, "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))
//
//	golaLoggerEntry := newLoggerEntry(ctx, logger.WithContext(context.Background()))
//
//	suite.Panics(func() { golaLoggerEntry.Fatalf("%s", "hi") })
//	suite.Equal("hi", hook.LastEntry().Message)
//	s.End()
//	suite.Equal("panic", t.SpanData.Annotations[0].Attributes["level"])
//	suite.Equal("hi", t.SpanData.Annotations[0].Message)
//}

func (suite *GolaLoggerTestSuite) TestShouldPrintPanicfLogsInLogrusAndSpanSuccessfully() {
	logger, hook := logrusTest.NewNullLogger()
	t := testExporter{}
	trace.RegisterExporter(&t)
	ctx, s := trace.StartSpan(suite.context, "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))

	golaLoggerEntry := newLoggerEntry(ctx, logger.WithContext(context.TODO()))

	suite.Panics(func() { golaLoggerEntry.Panicf("%s", "hi") })
	suite.Equal("hi", hook.LastEntry().Message)
	s.End()
	suite.Equal("panic", t.SpanData.Annotations[0].Attributes["level"])
	suite.Equal("hi", t.SpanData.Annotations[0].Message)
}

func (suite *GolaLoggerTestSuite) TestShouldPrintTraceLogsInLogrusAndSpanSuccessfully() {
	logger, hook := logrusTest.NewNullLogger()
	t := testExporter{}
	trace.RegisterExporter(&t)
	ctx, s := trace.StartSpan(suite.context, "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))

	golaLoggerEntry := newLoggerEntry(ctx, logger.WithContext(context.Background()))
	_ = golaLoggerEntry.SetLevel("trace")
	golaLoggerEntry.Trace("hi")
	s.End()

	suite.Equal("hi", hook.LastEntry().Message)
	suite.Equal(logrus.TraceLevel, hook.LastEntry().Level)
	suite.Equal("trace", t.SpanData.Annotations[0].Attributes["level"])
	suite.Equal("hi", t.SpanData.Annotations[0].Message)
	hook.Reset()
}

func (suite *GolaLoggerTestSuite) TestShouldPrintDebugLogsInLogrusAndSpanSuccessfully() {
	logger, hook := logrusTest.NewNullLogger()
	t := testExporter{}
	trace.RegisterExporter(&t)
	ctx, s := trace.StartSpan(suite.context, "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))

	golaLoggerEntry := newLoggerEntry(ctx, logger.WithContext(context.Background()))
	golaLoggerEntry.SetLevel("debug")
	golaLoggerEntry.Debug("hi")
	s.End()

	suite.Equal("hi", hook.LastEntry().Message)
	suite.Equal(logrus.DebugLevel, hook.LastEntry().Level)
	suite.Equal("debug", t.SpanData.Annotations[0].Attributes["level"])
	suite.Equal("hi", t.SpanData.Annotations[0].Message)
	hook.Reset()
}

func (suite *GolaLoggerTestSuite) TestShouldPrintInfoLogsInLogrusAndSpanSuccessfully() {
	logger, hook := logrusTest.NewNullLogger()
	t := testExporter{}
	trace.RegisterExporter(&t)
	ctx, s := trace.StartSpan(suite.context, "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))

	golaLoggerEntry := newLoggerEntry(ctx, logger.WithContext(context.Background()))
	golaLoggerEntry.Info("hi")
	s.End()

	suite.Equal("hi", hook.LastEntry().Message)
	suite.Equal(logrus.InfoLevel, hook.LastEntry().Level)
	suite.Equal("info", t.SpanData.Annotations[0].Attributes["level"])
	suite.Equal("hi", t.SpanData.Annotations[0].Message)
	hook.Reset()
}

func (suite *GolaLoggerTestSuite) TestShouldPrintPrintLogsInLogrusAndSpanSuccessfully() {
	logger, hook := logrusTest.NewNullLogger()
	t := testExporter{}
	trace.RegisterExporter(&t)
	ctx, s := trace.StartSpan(suite.context, "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))

	golaLoggerEntry := newLoggerEntry(ctx, logger.WithContext(context.Background()))
	golaLoggerEntry.Print("hi")
	s.End()

	suite.Equal("hi", hook.LastEntry().Message)
	suite.Equal(logrus.InfoLevel, hook.LastEntry().Level)
	suite.Equal("info", t.SpanData.Annotations[0].Attributes["level"])
	suite.Equal("hi", t.SpanData.Annotations[0].Message)
	hook.Reset()
}

func (suite *GolaLoggerTestSuite) TestShouldPrintWarnLogsInLogrusAndSpanSuccessfully() {
	logger, hook := logrusTest.NewNullLogger()
	t := testExporter{}
	trace.RegisterExporter(&t)
	ctx, s := trace.StartSpan(suite.context, "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))

	golaLoggerEntry := newLoggerEntry(ctx, logger.WithContext(context.Background()))
	golaLoggerEntry.Warn("hi")
	s.End()

	suite.Equal("hi", hook.LastEntry().Message)
	suite.Equal(logrus.WarnLevel, hook.LastEntry().Level)
	suite.Equal("warning", t.SpanData.Annotations[0].Attributes["level"])
	suite.Equal("hi", t.SpanData.Annotations[0].Message)
	hook.Reset()
}

func (suite *GolaLoggerTestSuite) TestShouldPrintWarningLogsInLogrusAndSpanSuccessfully() {
	logger, hook := logrusTest.NewNullLogger()
	t := testExporter{}
	trace.RegisterExporter(&t)
	ctx, s := trace.StartSpan(suite.context, "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))

	golaLoggerEntry := newLoggerEntry(ctx, logger.WithContext(context.Background()))
	golaLoggerEntry.Warning("hi")
	s.End()

	suite.Equal("hi", hook.LastEntry().Message)
	suite.Equal(logrus.WarnLevel, hook.LastEntry().Level)
	suite.Equal("warning", t.SpanData.Annotations[0].Attributes["level"])
	suite.Equal("hi", t.SpanData.Annotations[0].Message)
	hook.Reset()
}

func (suite *GolaLoggerTestSuite) TestShouldPrintErrorLogsInLogrusAndSpanSuccessfully() {
	logger, hook := logrusTest.NewNullLogger()
	t := testExporter{}
	trace.RegisterExporter(&t)
	ctx, s := trace.StartSpan(suite.context, "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))

	golaLoggerEntry := newLoggerEntry(ctx, logger.WithContext(context.Background()))
	golaLoggerEntry.Error("hi")
	s.End()

	suite.Equal("hi", hook.LastEntry().Message)
	suite.Equal(logrus.ErrorLevel, hook.LastEntry().Level)
	suite.Equal("error", t.SpanData.Annotations[0].Attributes["level"])
	suite.Equal("hi", t.SpanData.Annotations[0].Message)
	hook.Reset()
}

//func (suite *GolaLoggerTestSuite) TestShouldPrintFatalLogsInLogrusAndSpanSuccessfully() {
//	logger, hook := logrusTest.NewNullLogger()
//	handler := func() {
//
//		recover()
//
//	}
//	logrus.RegisterExitHandler(handler)
//	t := testExporter{}
//	trace.RegisterExporter(&t)
//	ctx, s := trace.StartSpan(suite.context, "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))
//
//	golaLoggerEntry := newLoggerEntry(ctx, logger.WithContext(context.Background()))
//
//	suite.Panics(func() { golaLoggerEntry.Fatal("hi") })
//	suite.Equal("hi", hook.LastEntry().Message)
//	s.End()
//	suite.Equal("panic", t.SpanData.Annotations[0].Attributes["level"])
//	suite.Equal("hi", t.SpanData.Annotations[0].Message)
//}

func (suite *GolaLoggerTestSuite) TestShouldPrintPanicLogsInLogrusAndSpanSuccessfully() {
	logger, hook := logrusTest.NewNullLogger()
	t := testExporter{}
	trace.RegisterExporter(&t)
	ctx, s := trace.StartSpan(suite.context, "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))

	golaLoggerEntry := newLoggerEntry(ctx, logger.WithContext(context.TODO()))

	suite.Panics(func() { golaLoggerEntry.Panic("hi") })
	suite.Equal("hi", hook.LastEntry().Message)
	s.End()
	suite.Equal("panic", t.SpanData.Annotations[0].Attributes["level"])
	suite.Equal("hi", t.SpanData.Annotations[0].Message)
}

func (suite *GolaLoggerTestSuite) TestShouldPrintTracelnLogsInLogrusAndSpanSuccessfully() {
	logger, hook := logrusTest.NewNullLogger()
	t := testExporter{}
	trace.RegisterExporter(&t)
	ctx, s := trace.StartSpan(suite.context, "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))

	golaLoggerEntry := newLoggerEntry(ctx, logger.WithContext(context.Background()))
	_ = golaLoggerEntry.SetLevel("trace")
	golaLoggerEntry.Traceln("hi")
	s.End()

	suite.Equal("hi", hook.LastEntry().Message)
	suite.Equal(logrus.TraceLevel, hook.LastEntry().Level)
	suite.Equal("trace", t.SpanData.Annotations[0].Attributes["level"])
	suite.Equal("hi", t.SpanData.Annotations[0].Message)
	hook.Reset()
}

func (suite *GolaLoggerTestSuite) TestShouldPrintDebuglnLogsInLogrusAndSpanSuccessfully() {
	logger, hook := logrusTest.NewNullLogger()
	t := testExporter{}
	trace.RegisterExporter(&t)
	ctx, s := trace.StartSpan(suite.context, "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))

	golaLoggerEntry := newLoggerEntry(ctx, logger.WithContext(context.Background()))
	golaLoggerEntry.SetLevel("debug")
	golaLoggerEntry.Debugln("hi")
	s.End()

	suite.Equal("hi", hook.LastEntry().Message)
	suite.Equal(logrus.DebugLevel, hook.LastEntry().Level)
	suite.Equal("debug", t.SpanData.Annotations[0].Attributes["level"])
	suite.Equal("hi", t.SpanData.Annotations[0].Message)
	hook.Reset()
}

func (suite *GolaLoggerTestSuite) TestShouldPrintInfolnLogsInLogrusAndSpanSuccessfully() {
	logger, hook := logrusTest.NewNullLogger()
	t := testExporter{}
	trace.RegisterExporter(&t)
	ctx, s := trace.StartSpan(suite.context, "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))

	golaLoggerEntry := newLoggerEntry(ctx, logger.WithContext(context.Background()))
	golaLoggerEntry.Infoln("hi")
	s.End()

	suite.Equal("hi", hook.LastEntry().Message)
	suite.Equal(logrus.InfoLevel, hook.LastEntry().Level)
	suite.Equal("info", t.SpanData.Annotations[0].Attributes["level"])
	suite.Equal("hi", t.SpanData.Annotations[0].Message)
	hook.Reset()
}

func (suite *GolaLoggerTestSuite) TestShouldPrintPrintlnLogsInLogrusAndSpanSuccessfully() {
	logger, hook := logrusTest.NewNullLogger()
	t := testExporter{}
	trace.RegisterExporter(&t)
	ctx, s := trace.StartSpan(suite.context, "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))

	golaLoggerEntry := newLoggerEntry(ctx, logger.WithContext(context.Background()))
	golaLoggerEntry.Println("hi")
	s.End()

	suite.Equal("hi", hook.LastEntry().Message)
	suite.Equal(logrus.InfoLevel, hook.LastEntry().Level)
	suite.Equal("info", t.SpanData.Annotations[0].Attributes["level"])
	suite.Equal("hi", t.SpanData.Annotations[0].Message)
	hook.Reset()
}

func (suite *GolaLoggerTestSuite) TestShouldPrintWarnlnLogsInLogrusAndSpanSuccessfully() {
	logger, hook := logrusTest.NewNullLogger()
	t := testExporter{}
	trace.RegisterExporter(&t)
	ctx, s := trace.StartSpan(suite.context, "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))

	golaLoggerEntry := newLoggerEntry(ctx, logger.WithContext(context.Background()))
	golaLoggerEntry.Warnln("hi")
	s.End()

	suite.Equal("hi", hook.LastEntry().Message)
	suite.Equal(logrus.WarnLevel, hook.LastEntry().Level)
	suite.Equal("warning", t.SpanData.Annotations[0].Attributes["level"])
	suite.Equal("hi", t.SpanData.Annotations[0].Message)
	hook.Reset()
}

func (suite *GolaLoggerTestSuite) TestShouldPrintWarninglnLogsInLogrusAndSpanSuccessfully() {
	logger, hook := logrusTest.NewNullLogger()
	t := testExporter{}
	trace.RegisterExporter(&t)
	ctx, s := trace.StartSpan(suite.context, "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))

	golaLoggerEntry := newLoggerEntry(ctx, logger.WithContext(context.Background()))
	golaLoggerEntry.Warningln("hi")
	s.End()

	suite.Equal("hi", hook.LastEntry().Message)
	suite.Equal(logrus.WarnLevel, hook.LastEntry().Level)
	suite.Equal("warning", t.SpanData.Annotations[0].Attributes["level"])
	suite.Equal("hi", t.SpanData.Annotations[0].Message)
	hook.Reset()
}

func (suite *GolaLoggerTestSuite) TestShouldPrintErrorlnLogsInLogrusAndSpanSuccessfully() {
	logger, hook := logrusTest.NewNullLogger()
	t := testExporter{}
	trace.RegisterExporter(&t)
	ctx, s := trace.StartSpan(suite.context, "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))

	golaLoggerEntry := newLoggerEntry(ctx, logger.WithContext(context.Background()))
	golaLoggerEntry.Errorln("hi")
	s.End()

	suite.Equal("hi", hook.LastEntry().Message)
	suite.Equal(logrus.ErrorLevel, hook.LastEntry().Level)
	suite.Equal("error", t.SpanData.Annotations[0].Attributes["level"])
	suite.Equal("hi", t.SpanData.Annotations[0].Message)
	hook.Reset()
}

//func (suite *GolaLoggerTestSuite) TestShouldPrintFatallnLogsInLogrusAndSpanSuccessfully() {
//	logger, hook := logrusTest.NewNullLogger()
//	handler := func() {
//
//		recover()
//
//	}
//	logrus.RegisterExitHandler(handler)
//	t := testExporter{}
//	trace.RegisterExporter(&t)
//	ctx, s := trace.StartSpan(suite.context, "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))
//
//	golaLoggerEntry := newLoggerEntry(ctx, logger.WithContext(context.Background()))
//
//	suite.Panics(func() { golaLoggerEntry.Fatalln("hi") })
//	suite.Equal("hi", hook.LastEntry().Message)
//	s.End()
//	suite.Equal("panic", t.SpanData.Annotations[0].Attributes["level"])
//	suite.Equal("hi", t.SpanData.Annotations[0].Message)
//}

func (suite *GolaLoggerTestSuite) TestShouldPrintPaniclnLogsInLogrusAndSpanSuccessfully() {
	logger, hook := logrusTest.NewNullLogger()
	t := testExporter{}
	trace.RegisterExporter(&t)
	ctx, s := trace.StartSpan(suite.context, "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))

	golaLoggerEntry := newLoggerEntry(ctx, logger.WithContext(context.TODO()))

	suite.Panics(func() { golaLoggerEntry.Panicln("hi") })
	suite.Equal("hi", hook.LastEntry().Message)
	s.End()
	suite.Equal("panic", t.SpanData.Annotations[0].Attributes["level"])
	suite.Equal("hi", t.SpanData.Annotations[0].Message)
}

func (suite *GolaLoggerTestSuite) TestShouldPrintSpanLogsWithDataSuccessfully() {
	logger, hook := logrusTest.NewNullLogger()
	t := testExporter{}
	trace.RegisterExporter(&t)
	ctx, s := trace.StartSpan(suite.context, "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))

	golaLoggerEntry := newLoggerEntry(ctx, logger.WithContext(context.Background())).WithField("dummyKey", "dummyValue")
	golaLoggerEntry.Warningln("hi")
	s.End()

	suite.Equal("hi", hook.LastEntry().Message)
	suite.Equal("dummyValue", hook.LastEntry().Data["dummyKey"])
	suite.Equal(logrus.WarnLevel, hook.LastEntry().Level)
	suite.Equal("dummyValue", t.SpanData.Annotations[0].Attributes["dummyKey"])
	suite.Equal("warning", t.SpanData.Annotations[0].Attributes["level"])
	suite.Equal("hi", t.SpanData.Annotations[0].Message)
	hook.Reset()
}

func (suite *GolaLoggerTestSuite) TestShouldPrintLogsFromGinContextSuccessfully() {
	logger, _ := logrusTest.NewNullLogger()
	t := testExporter{}
	trace.RegisterExporter(&t)
	r, _ := http.NewRequest("GEt", "/dummy", nil)
	r.Header.Set(constants.TRACING_SESSION_HEADER_KEY, "12345")
	ctx, s := trace.StartSpan(suite.context, "test-span", trace.WithSampler(trace.AlwaysSample()), trace.WithSpanKind(trace.SpanKindServer))
	suite.context.Request = r.WithContext(ctx)
	newLoggerEntry(suite.context, logger.WithContext(context.TODO())).Error("hi")

	s.End()
	suite.Equal("error", t.SpanData.Annotations[0].Attributes["level"])
	suite.Equal("hi", t.SpanData.Annotations[0].Message)
}

type TestHook struct{}

func (TestHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.ErrorLevel}
}
func (TestHook) Fire(e *logrus.Entry) error {
	return nil
}
func (suite *GolaLoggerTestSuite) TestShouldEnsureHooksAreAdded() {

	logger := logrus.StandardLogger()
	golaLoggerEntry := newLoggerEntry(context.Background(), logger.WithContext(context.Background()))
	hk := TestHook{}
	golaLoggerEntry.AddHook(hk)

	hooks := golaLoggerEntry.stdEntry.Logger.Hooks[logrus.ErrorLevel]
	suite.Equal(hk, hooks[0])
	fmt.Println(hooks)
}

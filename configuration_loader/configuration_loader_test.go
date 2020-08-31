package configuration_loader

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type testStruct struct {
	SomeConfig      string `json:"some_config" validate:"required"`
	SomeOtherConfig string `json:"some_other_config" validate:"required"`
}

type ConfigLoaderTestSuite struct {
	suite.Suite
	config       testStruct
	configLoader ConfigLoader
}

func TestConfigLoaderTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigLoaderTestSuite))
}

func (suite *ConfigLoaderTestSuite) SetupTest() {
	suite.configLoader = NewConfigLoader()
}

func (suite *ConfigLoaderTestSuite) TearDownTest() {
}

func (suite *ConfigLoaderTestSuite) TestConfigLoader_LoadConfig_WhenValidConfig() {
	err := suite.configLoader.Load("validTestConfig.json", &suite.config)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "some value", suite.config.SomeConfig)
	assert.Equal(suite.T(), "some other value", suite.config.SomeOtherConfig)
}

func (suite *ConfigLoaderTestSuite) TestConfigLoader_LoadConfig_WhenInvalidFilePath() {
	err := suite.configLoader.Load("InvalidPath/InvalidName.json", &suite.config)

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "open InvalidPath/InvalidName.json: no such file or directory", err.Error())
}

func (suite *ConfigLoaderTestSuite) TestConfigLoader_LoadConfig_WhenInvalidFileContent() {
	err := suite.configLoader.Load("invalidContentTestConfig.json", &suite.config)

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "invalid character 's' after object key", err.Error())
}

func (suite *ConfigLoaderTestSuite) TestConfigLoader_LoadConfig_WhenConfigValidationFails() {
	err := suite.configLoader.Load("invalidValidationTestConfig.json", &suite.config)

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "Key: 'testStruct.SomeConfig' Error:Field validation for 'SomeConfig' failed on the 'required' tag\n"+
		"Key: 'testStruct.SomeOtherConfig' Error:Field validation for 'SomeOtherConfig' failed on the 'required' tag", err.Error())
}

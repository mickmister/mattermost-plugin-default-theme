package main

import (
	"errors"
	"testing"

	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestOnConfigurationChangeValidTheme(t *testing.T) {
	theme := `{"awayIndicator":"#eae679","buttonBg":"#00aa55","buttonColor":"#041418","centerChannelBg":"#011010","centerChannelColor":"#1fbfec","codeTheme":"solarized-dark","dndIndicator":"#b79fa7","errorTextColor":"#dd2c45","linkColor":"#00aa55","mentionBg":"#00c463","mentionBj":"#00c463","mentionColor":"#001e27","mentionHighlightBg":"#032727","mentionHighlightLink":"#146c86","newMessageSeparator":"#146c86","onlineIndicator":"#0be3d6","sidebarBg":"#001e27","sidebarHeaderBg":"#001e27","sidebarHeaderTextColor":"#1fbfec","sidebarTeamBarBg":"#00181f","sidebarText":"#8ae9ff","sidebarTextActiveBorder":"#00c463","sidebarTextActiveColor":"#00c463","sidebarTextHoverBg":"#1a5c70","sidebarUnreadText":"#1fbfec"}`

	testAPI := &plugintest.API{}
	plugin := Plugin{}
	plugin.SetAPI(testAPI)

	testAPI.On("LoadPluginConfiguration", mock.AnythingOfType("*main.configuration")).Run(func(args mock.Arguments) {
		config := args.Get(0).(*configuration)
		config.CustomTheme = theme
	}).Return(nil)

	err := plugin.OnConfigurationChange()
	require.Nil(t, err)

	config := plugin.getConfiguration()
	require.Equal(t, config.CustomTheme, theme)
}

func TestOnConfigurationChangeBlankTheme(t *testing.T) {
	theme := ""

	testAPI := &plugintest.API{}
	plugin := Plugin{}
	plugin.SetAPI(testAPI)

	testAPI.On("LoadPluginConfiguration", mock.AnythingOfType("*main.configuration")).Run(func(args mock.Arguments) {
		config := args.Get(0).(*configuration)
		config.CustomTheme = theme
	}).Return(nil)

	err := plugin.OnConfigurationChange()
	require.Nil(t, err)

	config := plugin.getConfiguration()
	require.Equal(t, config.CustomTheme, theme)
}

func TestOnConfigurationChangeInvalidTheme(t *testing.T) {
	theme := `{"invalid": "json"`

	testAPI := &plugintest.API{}
	plugin := Plugin{}
	plugin.SetAPI(testAPI)

	testAPI.On("LoadPluginConfiguration", mock.AnythingOfType("*main.configuration")).Run(func(args mock.Arguments) {
		config := args.Get(0).(*configuration)
		config.CustomTheme = theme
	}).Return(nil)

	err := plugin.OnConfigurationChange()
	require.NotNil(t, err)
	require.Equal(t, err.Error(), "error parsing default theme while saving config: unexpected end of JSON input")

	config := plugin.getConfiguration()
	require.Equal(t, config.CustomTheme, "")
}

func TestOnConfigurationChangeErrorLoadingConfig(t *testing.T) {
	testAPI := &plugintest.API{}
	plugin := Plugin{}
	plugin.SetAPI(testAPI)

	testAPI.On("LoadPluginConfiguration", mock.AnythingOfType("*main.configuration")).Return(errors.New("some error"))

	err := plugin.OnConfigurationChange()
	require.NotNil(t, err)
	require.Equal(t, err.Error(), "failed to load plugin configuration: some error")

	config := plugin.getConfiguration()
	require.Equal(t, config.CustomTheme, "")
}

package main

import (
	"testing"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserHasBeenCreatedValidTheme(t *testing.T) {
	theme := `{"awayIndicator":"#eae679","buttonBg":"#00aa55","buttonColor":"#041418","centerChannelBg":"#011010","centerChannelColor":"#1fbfec","codeTheme":"solarized-dark","dndIndicator":"#b79fa7","errorTextColor":"#dd2c45","linkColor":"#00aa55","mentionBg":"#00c463","mentionBj":"#00c463","mentionColor":"#001e27","mentionHighlightBg":"#032727","mentionHighlightLink":"#146c86","newMessageSeparator":"#146c86","onlineIndicator":"#0be3d6","sidebarBg":"#001e27","sidebarHeaderBg":"#001e27","sidebarHeaderTextColor":"#1fbfec","sidebarTeamBarBg":"#00181f","sidebarText":"#8ae9ff","sidebarTextActiveBorder":"#00c463","sidebarTextActiveColor":"#00c463","sidebarTextHoverBg":"#1a5c70","sidebarUnreadText":"#1fbfec"}`

	testAPI := &plugintest.API{}
	plugin := Plugin{}
	plugin.SetAPI(testAPI)
	plugin.setConfiguration(&configuration{
		CustomTheme: theme,
	})

	user := &model.User{
		Id: "userid",
	}

	testAPI.On("UpdatePreferencesForUser", mock.AnythingOfType("string"), mock.AnythingOfType("[]model.Preference")).Run(func(args mock.Arguments) {
		userID := args.Get(0).(string)
		require.Equal(t, "userid", userID)

		prefs := args.Get(1).([]model.Preference)
		require.Len(t, prefs, 1)

		require.Equal(t, model.Preference{
			UserId:   "userid",
			Category: "theme",
			Name:     "",
			Value:    theme,
		}, prefs[0])
	}).Once().Return(nil)
	plugin.UserHasBeenCreated(nil, user)
}

func TestUserHasBeenCreatedBlankTheme(t *testing.T) {
	theme := ""

	testAPI := &plugintest.API{}
	plugin := Plugin{}
	plugin.SetAPI(testAPI)
	plugin.setConfiguration(&configuration{
		CustomTheme: theme,
	})

	user := &model.User{
		Id: "userid",
	}
	plugin.UserHasBeenCreated(nil, user)
}

func TestUserHasBeenCreatedInvalidTheme(t *testing.T) {
	theme := `{"invalid": "json"`

	testAPI := &plugintest.API{}
	plugin := Plugin{}
	plugin.SetAPI(testAPI)
	plugin.setConfiguration(&configuration{
		CustomTheme: theme,
	})

	user := &model.User{
		Id: "userid",
	}

	testAPI.On("LogError", "error parsing theme during user creation: unexpected end of JSON input").Once()
	plugin.UserHasBeenCreated(nil, user)
}

func TestUserHasBeenCreatedErrorSettingPreferences(t *testing.T) {
	theme := `{"awayIndicator":"#eae679","buttonBg":"#00aa55","buttonColor":"#041418","centerChannelBg":"#011010","centerChannelColor":"#1fbfec","codeTheme":"solarized-dark","dndIndicator":"#b79fa7","errorTextColor":"#dd2c45","linkColor":"#00aa55","mentionBg":"#00c463","mentionBj":"#00c463","mentionColor":"#001e27","mentionHighlightBg":"#032727","mentionHighlightLink":"#146c86","newMessageSeparator":"#146c86","onlineIndicator":"#0be3d6","sidebarBg":"#001e27","sidebarHeaderBg":"#001e27","sidebarHeaderTextColor":"#1fbfec","sidebarTeamBarBg":"#00181f","sidebarText":"#8ae9ff","sidebarTextActiveBorder":"#00c463","sidebarTextActiveColor":"#00c463","sidebarTextHoverBg":"#1a5c70","sidebarUnreadText":"#1fbfec"}`

	testAPI := &plugintest.API{}
	plugin := Plugin{}
	plugin.SetAPI(testAPI)
	plugin.setConfiguration(&configuration{
		CustomTheme: theme,
	})

	user := &model.User{
		Id: "userid",
	}

	testAPI.On("UpdatePreferencesForUser", mock.AnythingOfType("string"), mock.AnythingOfType("[]model.Preference")).Once().Return(&model.AppError{Message: "preference error"})
	testAPI.On("LogError", "error setting preferences for user userid. err=: preference error, ").Once()
	plugin.UserHasBeenCreated(nil, user)
}

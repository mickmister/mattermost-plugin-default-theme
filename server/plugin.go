package main

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration
}

// UserHasBeenCreated is invoked after a user was created.
func (p *Plugin) UserHasBeenCreated(c *plugin.Context, user *model.User) {
	theme := p.getConfiguration().CustomTheme
	if theme == "" {
		return
	}

	m := map[string]interface{}{}
	err := json.Unmarshal([]byte(theme), &m)
	if err != nil {
		p.API.LogError("error parsing theme during user creation: " + err.Error())
		return
	}

	pref := model.Preference{
		UserId:   user.Id,
		Category: model.PREFERENCE_CATEGORY_THEME,
		Name:     "",
		Value:    theme,
	}
	prefs := []model.Preference{pref}

	appErr := p.API.UpdatePreferencesForUser(user.Id, prefs)
	if appErr != nil {
		errString := fmt.Sprintf("error setting preferences for user %s. err=%s", user.Id, appErr.Error())
		p.API.LogError(errString)
	}
}

// See https://developers.mattermost.com/extend/plugins/server/reference/

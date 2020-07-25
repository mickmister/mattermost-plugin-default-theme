package main

import (
	"fmt"
	"net/http"
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

// ServeHTTP demonstrates a plugin that handles HTTP requests by greeting the world.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}

// UserHasBeenCreated is invoked after a user was created.
func (p *Plugin) UserHasBeenCreated(c *plugin.Context, user *model.User) {
	theme := p.getConfiguration().CustomTheme
	if theme == "" {
		return
	}

	pref := model.Preference{
		UserId:   user.Id,
		Category: model.PREFERENCE_CATEGORY_THEME,
		Name:     "",
		Value:    theme,
	}
	prefs := []model.Preference{pref}
	p.API.UpdatePreferencesForUser(user.Id, prefs)
}

// See https://developers.mattermost.com/extend/plugins/server/reference/

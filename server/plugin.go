package main

import (
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/pkg/errors"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration

	BotUserID string
	store     Store
	router    *mux.Router
}

// ServeHTTP demonstrates a plugin that handles HTTP requests by greeting the world.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	r.Header.Add("Mattermost-Plugin-ID", c.SourcePluginId)
	w.Header().Set("Content-Type", "application/json")

	p.router.ServeHTTP(w, r)
}

// See https://developers.mattermost.com/extend/plugins/server/reference/
func (p *Plugin) OnActivate() error {
	botID, err := p.Helpers.EnsureBot(&model.Bot{
		Username:    "badges",
		DisplayName: "Badges Bot",
		Description: "Created by the Badges plugin.",
	})
	if err != nil {
		return errors.Wrap(err, "failed to ensure badges bot")
	}
	p.BotUserID = botID
	p.store = NewStore(p.API)
	p.initializeAPI()

	return p.API.RegisterCommand(p.getCommand())
}

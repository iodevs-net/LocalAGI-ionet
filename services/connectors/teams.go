package connectors

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/mudler/LocalAGI/core/agent"
	"github.com/mudler/LocalAGI/core/types"
	"github.com/mudler/LocalAGI/pkg/config"
	"github.com/mudler/xlog"
)

// TeamsConnector - KISS integration with Microsoft Teams via Incoming Webhook
// Only sends messages TO Teams (notifications, alerts, summaries)
// For receiving messages FROM Teams, use email connector (already implemented)
type TeamsConnector struct {
	webhookURL string
}

// TeamsMessage - Simple message format for Teams Incoming Webhook
type TeamsMessage struct {
	Type    string `json:"@type"`
	Text    string `json:"text"`
}

// NewTeamsConnector - Creates a new Teams connector from config map
// config["webhookUrl"]: Incoming Webhook URL from Teams Channel
// Get webhook URL: Teams Channel > Connectors > Incoming Webhook > Configure
func NewTeamsConnector(config map[string]string) *TeamsConnector {
	webhookURL := config["webhookUrl"]
	if webhookURL == "" {
		xlog.Error("Teams connector: webhookUrl is required")
		return nil
	}
	return &TeamsConnector{
		webhookURL: webhookURL,
	}
}

// SendMessage - Sends a text message to the configured Teams channel
func (t *TeamsConnector) SendMessage(text string) error {
	msg := TeamsMessage{
		Type: "Message",
		Text: text,
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		xlog.Error("Teams connector: failed to marshal message", "error", err)
		return err
	}

	resp, err := http.Post(t.webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		xlog.Error("Teams connector: failed to send message", "error", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		xlog.Error("Teams connector: webhook returned non-200 status", "status", resp.StatusCode)
		return nil // Don't fail the agent if webhook fails
	}

	xlog.Info("Teams connector: message sent successfully")
	return nil
}

func (t *TeamsConnector) AgentResultCallback() func(state types.ActionState) {
	return nil
}

func (t *TeamsConnector) AgentReasoningCallback() func(state types.ActionCurrentState) bool {
	return nil
}

func (t *TeamsConnector) Start(a *agent.Agent) {
	// Teams connector currently only supports outgoing messages via webhook
}

// TeamsConfigMeta - Returns configuration fields for Teams connector
func TeamsConfigMeta() []config.Field {
	return []config.Field{
		{
			Name:        "webhookUrl",
			Label:       "Webhook URL",
			Type:        config.FieldTypeText,
			HelpText:    "Incoming Webhook URL from Teams Channel (Connectors > Incoming Webhook)",
			Required:    true,
		},
	}
}

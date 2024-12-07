package templates

import (
	"bytes"
	"html/template"

	"example.com/m/store"
)

const (
	messagesHTML = `
	<div id="messages">
		{{ range .Messages }}
		<div>{{ .ClientID }}: {{ .Text }}</div>
		{{ end }}
	</div>`
	notificationsHTML = `
	<div id="notifications">
		{{ range .Notifications }}
		<div>{{ . }}</div>
		{{ end }}
	</div>`
)

var messagesTemplate = template.Must(template.New("messages").Parse(messagesHTML))
var notificationsTemplate = template.Must(template.New("notifications").Parse(notificationsHTML))

func MessagesHtml(messages ...*store.Message) (string, error) {
	var buf bytes.Buffer
	data := struct {
		Messages []*store.Message
	}{
		Messages: messages,
	}
	if err := messagesTemplate.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func NotificationsHtml(notifications ...string) (string, error) {
	var buf bytes.Buffer
	data := struct {
		Notifications []string
	}{
		Notifications: notifications,
	}
	if err := notificationsTemplate.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

package templates

import (
	"fmt"

	"example.com/m/store"
)

func MessagesHtml(messages ...*store.Message) string {
	res := fmt.Sprintf("<div id=\"messages\">")
	for _, msg := range messages {
		res += wrapInDiv(fmt.Sprintf("[%s]: %s", msg.ClientID.String(), msg.Text))
	}
	return res + "</div"
}

func NotificationsHtml(notifications ...string) string {
	res := fmt.Sprintf("<div id=\"notifications\">")
	for _, notif := range notifications {
		res += wrapInDiv(notif)
	}
	return res + "</div"

}

func wrapInDiv(content string) string {
	return "<div>" + content + "</div>"
}

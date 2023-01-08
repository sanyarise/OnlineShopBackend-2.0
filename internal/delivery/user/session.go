package user

import (
	"go.uber.org/zap"
	"net/http"

	"github.com/dghubble/sessions"
)

const (
	sessionName    = "example-google-app"
	sessionSecret  = "example cookie signing secret"
	sessionUserKey = "key"
	sessionUserID  = "userId"
)

var sessionStore = sessions.NewCookieStore([]byte(sessionSecret), nil)

func IssueSession(l *zap.Logger,userId string) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		session := sessions.NewSession(sessionStore, sessionName)
		if err := session.Save(w); err != nil {
			l.Error("unable to save session")
		}
		session.Values[sessionUserID] = userId
	}
	return http.HandlerFunc(fn)
}


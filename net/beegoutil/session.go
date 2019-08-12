package beegoutil

import (
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/session"
)

const (
	BeegoSessionCookieNameCfgVar  string = "sessioncookiename"
	BeegoSessionCookieNameDefault string = "gosessionid"
	BeegoSessionProviderCfgVar    string = "sessionprovidername"
	BeegoSessionProviderDefault   string = "memory"
)

// InitSession creates a starts session management https://beego.me/docs/module/session.md
func InitSession(sessionProvider string, sessionConfig *session.ManagerConfig) {
	sessionProvider = strings.TrimSpace(sessionProvider)
	if len(sessionProvider) == 0 {
		sessionProvider = strings.TrimSpace(
			beego.AppConfig.String(BeegoSessionProviderCfgVar))
		if len(sessionProvider) == 0 {
			sessionProvider = BeegoSessionProviderDefault
		}
	}

	if sessionConfig != nil {
		globalSessions, _ := session.NewManager(sessionProvider, sessionConfig)
		go globalSessions.GC()
		return
	}

	sessionConfig = &session.ManagerConfig{Gclifetime: 3600}

	cfgCookieName := strings.TrimSpace(
		beego.AppConfig.String(BeegoSessionCookieNameCfgVar))
	if len(cfgCookieName) > 0 {
		sessionConfig.CookieName = cfgCookieName
	} else {
		sessionConfig.CookieName = BeegoSessionCookieNameDefault
	}

	globalSessions, _ := session.NewManager(sessionProvider, sessionConfig)
	go globalSessions.GC()
}

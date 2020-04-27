package sessions

import (
	"github.com/golangext/datastructures/session"
	"github.com/labstack/echo"
)

const session_key = "golangext_datastructures_middleware_sessions_session"

func Get(c echo.Context) session.Session {
	sess := c.Get(session_key)
	return sess.(session.Session)
}

func Set(c echo.Context, sess session.Session) {
	c.Set(session_key, sess)
}
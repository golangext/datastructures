package sessions

import (
	"net/http"
	"time"

	"github.com/golangext/datastructures/session"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

type cacheRestorer struct {
	cookieName string
	cache      session.Cache
}

func (restorer *cacheRestorer) do(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		currentSession := Get(c)
		if currentSession == nil {
			cook, err := c.Cookie(c.cookieName)
			if err != nil {
				log.Debugf("No cookie '%v' available: %v", restorer.cookieName, err)
			} else if cook == nil {
				log.Debugf("No cookie '%v' available", restorer.cookieName)
			} else {
				sess := restorer.cache.Get(cook.Value)
				if sess == nil {
					log.Debugf("Cookie '%v' references session '%v' not in cache", restorer.cookieName, sess.ID())
				} else {
					Set(c, sess)
					log.Debugf("Session '%v' restored from cookie '%v'", sess.ID(), restorer.cookieName)
				}
			}
		} else {
			log.Debugf("No need to restore session from cache. Current session is '%v'", currentSession.ID())
		}
		return next(c)
	}
}

type func(sessionID string) Session SessionSourceFunc

func RestoreFromCache(cache session.Cache, cookieName string) echo.MiddlewareFunc {
	v := cacheRestorer{cookieName: cookieName, cache: cache}
	return v.do
}

func RestoreFromSource(cache session.Cache, cookieName string, source SessionSourceFunc) echo.MiddlewareFunc {
	v := cacheRestorer{cookieName: cookieName, cache: cache}
	return v.do
}

func SetTempCookie(c echo.Context, cookieName string, sessionID string) {
	cookie := new(http.Cookie)
	cookie.Name = cookieName
	cookie.Value = sessionID
	c.SetCookie(cookie)
}

func SetPersistentCookie(c echo.Context, cookieName string, sessionID string, duration time.Duration) {
	cookie := new(http.Cookie)
	cookie.Name = cookieName
	cookie.Value = sessionID
	cookie.Expires = time.Now().Add(duration)
	c.SetCookie(cookie)
}

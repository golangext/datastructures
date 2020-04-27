package sessions

type funcRestorer struct {
	cookieName string
	cache      session.Cache
	fun SessionSourceFunc
}

func (restorer *funcRestorer) do(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		currentSession := Get(c)
		if currentSession == nil {
			cook, err := c.Cookie(restorer.cookieName)
			if err != nil {
				log.Debugf("No cookie '%v' available: %v", restorer.cookieName, err)
			} else if cook == nil {
				log.Debugf("No cookie '%v' available", restorer.cookieName)
			} else {
				sess := restorer.fun(cook.Value)
				if sess == nil {
					log.Debugf("Cookie '%v' references session '%v' not in cache", restorer.cookieName, sess.ID())
				} else {
					Set(c, sess)
					restorer.cache.Resurect(sess)
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

func RestoreFromFunc(cache session.Cache, cookieName string, fun SessionSourceFunc) echo.MiddlewareFunc {
	v := funcRestorer{cookieName: cookieName, cache: cache, fun: fun}
	return v.do
}


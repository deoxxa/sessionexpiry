package sessionexpiry

import (
	"net/http"
	"time"

	"fknsrs.biz/p/cookiesession"
)

type SessionExpiry struct {
	Store *cookiesession.Store
	TTL   time.Duration
}

func (s SessionExpiry) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if sess := s.Store.Get(r); sess.Valid {
		if time.Now().After(sess.Time.Add(s.TTL / 2)) {
			if err := s.Store.Save(rw, &sess); err != nil {
				panic(err)
			}
		}

		rw.Header().Set("x-session-lifetime", s.TTL.String())
		rw.Header().Set("x-session-expires", sess.Time.Add(s.TTL).Format(time.RFC3339))
		rw.Header().Set("x-session-remaining", sess.Time.Add(s.TTL).Sub(time.Now()).String())
	}

	next(rw, r)
}

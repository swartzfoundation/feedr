package model

import (
	"encoding/base32"
	"net/http"
	"strings"
	"time"

	"github.com/swartzfoundation/feedr/pkg/config"

	"log/slog"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

const defaultMaxAge = int(30 * 24 * time.Hour / time.Second) // 2592000 seconds = 30 days
const sessionIDLen = 32

var sessionStore *DatabaseStore

type Session struct {
	ID        string `sql:"unique_index"`
	Data      string `sql:"type:text"`
	CreatedAt int64
	UpdatedAt int64
	ExpiresAt int64 `sql:"index"`
}

type DatabaseStore struct {
	Codecs      []securecookie.Codec
	SessionOpts *sessions.Options
}

func NewSessionStore(options *sessions.Options, keyPairs ...[]byte) *DatabaseStore {
	ds := &DatabaseStore{
		Codecs: securecookie.CodecsFromPairs(keyPairs...),
		SessionOpts: &sessions.Options{
			Domain:   options.Domain,
			Path:     options.Path,
			MaxAge:   defaultMaxAge,
			SameSite: options.SameSite,
		},
	}

	// set cookie same site to none if debug mode
	if config.Config.DEBUG {
		ds.SessionOpts.Secure = true
		ds.SessionOpts.SameSite = http.SameSiteNoneMode
	}
	if err := db.AutoMigrate(&Session{}); err != nil {
		slog.Error("db: migrating session table", "error", err)
	}
	sessionStore = ds

	return ds
}

func GetSessionsStore() *DatabaseStore {
	return sessionStore
}

// MaxLength restricts the maximum length of new sessions to l.
// If l is 0 there is no limit to the size of a session, use with caution.
// The default is 4096 (default for securecookie)
func (ds *DatabaseStore) MaxLength(l int) {
	for _, c := range ds.Codecs {
		if codec, ok := c.(*securecookie.SecureCookie); ok {
			codec.MaxLength(l)
		}
	}
}

// Get returns a session for the given name after adding it to the registry.
func (ds *DatabaseStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(sessionStore, name)
}

// MaxAge sets the maximum age for the store and the underlying cookie
// implementation. Individual sessions can be deleted by setting
// Options.MaxAge = -1 for that session.
func (ds *DatabaseStore) MaxAge(age int) {
	ds.SessionOpts.MaxAge = age
	for _, codec := range ds.Codecs {
		if sc, ok := codec.(*securecookie.SecureCookie); ok {
			sc.MaxAge(age)
		}
	}
}

// New creates a session with name without adding it to the registry.
func (ds *DatabaseStore) New(r *http.Request, name string) (*sessions.Session, error) {
	session := sessions.NewSession(ds, name)
	opts := *ds.SessionOpts
	session.Options = &opts
	session.IsNew = true

	ds.MaxAge(ds.SessionOpts.MaxAge)

	// try fetch from db if there is a cookie
	s := ds.getSessionFromCookie(r, session.Name())
	if s != nil {
		if err := securecookie.DecodeMulti(session.Name(), s.Data, &session.Values, ds.Codecs...); err != nil {
			return session, nil
		}
		session.ID = s.ID
		session.IsNew = false
	}

	return session, nil
}

// Save session and set cookie header
func (ds *DatabaseStore) Save(
	r *http.Request,
	w http.ResponseWriter,
	session *sessions.Session,
) error {
	s := ds.getSessionFromCookie(r, session.Name())

	// delete if max age is < 0
	if session.Options.MaxAge < 0 {
		if s != nil {
			if err := db.Delete(s).Error; err != nil {
				return err
			}
		}
		http.SetCookie(w, sessions.NewCookie(session.Name(), "", session.Options))
		return nil
	}

	data, err := securecookie.EncodeMulti(session.Name(), session.Values, ds.Codecs...)
	if err != nil {
		return err
	}
	now := time.Now()

	expire := now.Add(time.Duration(session.Options.MaxAge) * time.Second)

	if s == nil {
		// generate random session ID key suitable for storage in the db
		session.ID = strings.TrimRight(
			base32.StdEncoding.EncodeToString(
				securecookie.GenerateRandomKey(sessionIDLen)), "=")
		s = &Session{
			ID:        session.ID,
			Data:      data,
			CreatedAt: now.Unix(),
			UpdatedAt: now.Unix(),
			ExpiresAt: expire.Unix(),
		}
		if err := db.Create(s).Error; err != nil {
			return err
		}
	} else {
		s.Data = data
		s.UpdatedAt = now.Unix()
		s.ExpiresAt = expire.Unix()
		if err := db.Save(s).Error; err != nil {
			return err
		}
	}

	// set session id cookie
	id, err := securecookie.EncodeMulti(session.Name(), s.ID, ds.Codecs...)
	if err != nil {
		return err
	}
	http.SetCookie(w, sessions.NewCookie(session.Name(), id, session.Options))
	return nil
}

// getSessionFromCookie looks for an existing Session from a session ID stored inside a cookie
func (ds *DatabaseStore) getSessionFromCookie(r *http.Request, name string) *Session {
	if cookie, err := r.Cookie(name); err == nil {
		sessionID := ""
		if err := securecookie.DecodeMulti(name, cookie.Value, &sessionID, ds.Codecs...); err != nil {
			return nil
		}
		s := &Session{}
		result := db.Where("id = ? AND expires_at > ?", sessionID, time.Now().Unix()).Limit(1).Find(s)
		if result.Error != nil || result.RowsAffected == 0 {
			return nil
		}
		return s
	}
	return nil
}

func (ds *DatabaseStore) Remove(r *http.Request, w http.ResponseWriter, session *sessions.Session) {
	if err := db.Delete(&Session{}, "id <= ?", session.ID).Error; err != nil {
		slog.Error("removing session:", "error", err)
	}
	options := &sessions.Options{
		Domain:   config.Config.Session.SessionCookieDomain,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	}
	if config.Config.DEBUG {
		options.Secure = true
		options.SameSite = http.SameSiteNoneMode
	}
	http.SetCookie(w, sessions.NewCookie(session.Name(), "", options))
}

// Cleanup deletes expired sessions
func (ds *DatabaseStore) Cleanup() {
	slog.Warn("Cleaning up expired sessions")
	if err := db.Delete(&Session{}, "expires_at <= ?", time.Now().Unix()).Error; err != nil {
		slog.Error("cleaning up expired sessions", "error", err)
	}
}

// PeriodicCleanup runs Cleanup every interval. Close quit channel to stop.
func (ds *DatabaseStore) PeriodicCleanup(interval time.Duration, quit <-chan struct{}) {
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			ds.Cleanup()
		case <-quit:
			return
		}
	}
}

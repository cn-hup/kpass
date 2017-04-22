package model

import (
	"fmt"
	"strconv"
	"time"

	"github.com/seccom/kpass/src/auth"
	"github.com/seccom/kpass/src/schema"
	"github.com/seccom/kpass/src/service"
	"github.com/seccom/kpass/src/util"
	"github.com/teambition/gear"
	"github.com/tidwall/buntdb"
)

// User is database access oject for users
type User struct {
	db *service.DB
}

// Init ...
func (m *User) Init(db *service.DB) *User {
	m.db = db
	return m
}

// CheckID ...
func (m *User) CheckID(id string) error {
	if len(id) < 3 {
		return gear.ErrBadRequest.WithMsg(fmt.Sprintf(`invalid user id "%s"`, id))
	}
	err := m.db.DB.View(func(tx *buntdb.Tx) error {
		if _, e := tx.Get(schema.UserKey(id)); e == nil {
			return gear.ErrConflict.WithMsg(fmt.Sprintf(`user "%s" exists`, id))
		}
		return nil
	})
	return dbError(err)
}

// CheckLogin ...
func (m *User) CheckLogin(id, pass string) (user *schema.User, err error) {
	verified := false
	err = m.db.DB.Update(func(tx *buntdb.Tx) error {
		userKey := schema.UserKey(id)
		value, e := tx.Get(userKey)
		if e != nil {
			return e
		}
		user, e = schema.UserFrom(value)
		if e != nil {
			return e
		}

		attempts, _ := tx.Get(schema.LoginAttemptKey(id))
		i, _ := strconv.Atoi(attempts)
		if user.IsBlocked {
			return gear.ErrForbidden.WithMsg("user blocked")
		}
		if i > 5 {
			return gear.ErrForbidden.WithMsg("too many login attempts, please retry after 2 hours")
		}

		if !auth.VerifyPass(id, pass, user.Pass) {
			i++
			tx.Set(schema.LoginAttemptKey(id), strconv.Itoa(i), &buntdb.SetOptions{
				Expires: true,
				TTL:     2 * time.Hour,
			})
			return nil
		} else if i > 0 {
			tx.Delete(schema.LoginAttemptKey(id))
		}
		verified = true
		return nil
	})

	if !verified && err == nil {
		err = gear.ErrBadRequest.WithMsg("user id or password error")
	}
	if err != nil {
		return nil, dbError(err)
	}
	return
}

// Create ...
func (m *User) Create(userID, pass string) (user *schema.User, err error) {
	err = m.db.DB.Update(func(tx *buntdb.Tx) error {
		userKey := schema.UserKey(userID)
		_, e := tx.Get(userKey)
		if e == nil {
			return gear.ErrConflict.WithMsg(fmt.Sprintf(`user "%s" exists`, userID))
		}

		user = &schema.User{
			ID:      userID,
			Pass:    auth.SignPass(userID, pass),
			Created: util.Time(time.Now()),
		}
		user.Updated = user.Created
		_, _, e = tx.Set(userKey, user.String(), nil)
		return e
	})

	if err != nil {
		return nil, dbError(err)
	}
	return
}

// Find ...
func (m *User) Find(id string) (user *schema.User, err error) {
	err = m.db.DB.View(func(tx *buntdb.Tx) error {
		res, e := tx.Get(schema.UserKey(id))
		if e == nil {
			user, e = schema.UserFrom(res)
		}
		return e
	})
	if err != nil {
		return nil, dbError(err)
	}
	return
}

// Update ...
func (m *User) Update(user *schema.User) error {
	err := m.db.DB.Update(func(tx *buntdb.Tx) error {
		user.Updated = util.Time(time.Now())
		_, _, e := tx.Set(schema.UserKey(user.ID), user.String(), nil)
		return e
	})
	return dbError(err)
}

// FindUsers ...
func (m *User) FindUsers(ids ...string) (users []*schema.UserResult, err error) {
	err = m.db.DB.View(func(tx *buntdb.Tx) error {
		users = IdsToUsers(tx, ids)
		return nil
	})
	if err != nil {
		return nil, dbError(err)
	}
	return
}

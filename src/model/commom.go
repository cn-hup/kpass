package model

import (
	"github.com/seccom/kpass/src/schema"
	"github.com/seccom/kpass/src/service"
	"github.com/teambition/gear"
	"github.com/tidwall/buntdb"
)

func dbError(err error) error {
	if err == nil {
		return nil
	}
	if err == buntdb.ErrNotFound {
		return gear.ErrNotFound.WithMsg(err.Error())
	}
	return gear.ErrInternalServerError.From(err)
}

// IdsToUsers ...
func IdsToUsers(tx *buntdb.Tx, ids []string) (users []*schema.UserResult) {
	for _, id := range ids {
		if res, e := tx.Get(schema.UserKey(id)); e == nil {
			if user, e := schema.UserFrom(res); e == nil {
				users = append(users, user.Result())
			}
		}
	}
	return
}

// Models ....
type Models struct {
	Entry  *Entry
	File   *File
	Secret *Secret
	Share  *Share
	Team   *Team
	User   *User
}

// Init ...
func (ms *Models) Init(db *service.DB) *Models {
	ms.Entry = new(Entry).Init(db)
	ms.File = new(File).Init(db)
	ms.Secret = new(Secret).Init(db)
	ms.Share = new(Share).Init(db)
	ms.Team = new(Team).Init(db)
	ms.User = new(User).Init(db)
	return ms
}

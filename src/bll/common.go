package bll

import (
	"github.com/seccom/kpass/src/model"
	"github.com/seccom/kpass/src/service"
)

// Bll is Business Logic Layer with all models
type Bll struct {
	Models *model.Models
}

// Blls ...
type Blls struct {
	Models *model.Models
	Team   *Team
	Entry  *Entry
	Secret *Secret
}

// Init ...
func (bs *Blls) Init(db *service.DB) *Blls {
	bs.Models = new(model.Models).Init(db)
	b := &Bll{bs.Models}

	bs.Team = &Team{b}
	bs.Entry = &Entry{b}
	bs.Secret = &Secret{b}
	return bs
}

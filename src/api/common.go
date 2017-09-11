package api

import (
	"github.com/seccom/kpass/src/bll"
	"github.com/seccom/kpass/src/model"
)

// APIs - all APIs
type APIs struct {
	Entry  *Entry
	Secret *Secret
	Share  *Share
	Team   *Team
	User   *User
}

// CommonAPI -
type CommonAPI struct {
	blls   *bll.Blls
	models *model.Models
}

// Init - 初始化
func (a *APIs) Init(blls *bll.Blls) *APIs {
	api := CommonAPI{blls, blls.Models}
	*a = APIs{
		Entry:  new(Entry).Init(api),
		Secret: new(Secret).Init(api),
		Share:  new(Share).Init(api),
		Team:   new(Team).Init(api),
		User:   new(User).Init(api),
	}

	return a
}

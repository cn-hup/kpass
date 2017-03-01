package api

import (
	"github.com/seccom/kpass/src/auth"
	"github.com/seccom/kpass/src/dao"
	"github.com/seccom/kpass/src/schema"
	"github.com/seccom/kpass/src/service"
	"github.com/seccom/kpass/src/util"
	"github.com/teambition/gear"
)

// Team is API oject for teams
//
// @Name Team
// @Description Team API
// @Accepts json
// @Produces json
type Team struct {
	team *dao.Team
	file *dao.File
}

// NewTeam returns a Team API instance
func NewTeam(db *service.DB) *Team {
	return &Team{dao.NewTeam(db), dao.NewFile(db)}
}

type tplTeamCreate struct {
	Name string `json:"name" swaggo:"true,team name,Teambition"`
}

func (t *tplTeamCreate) Validate() error {
	if t.Name == "" {
		return &gear.Error{Code: 400, Msg: "invalid team name"}
	}
	return nil
}

// Create ...
//
// @Title Create
// @Summary Create a team
// @Description Create a team
// @Param Authorization header string true "access_token"
// @Param body body tplTeamCreate true "team body"
// @Success 200 schema.TeamResult
// @Failure 400 string
// @Failure 401 string
// @Router POST /api/teams
func (a *Team) Create(ctx *gear.Context) error {
	body := new(tplTeamCreate)
	if err := ctx.ParseBody(body); err != nil {
		return ctx.Error(err)
	}

	key, err := auth.KeyFromCtx(ctx)
	if err != nil {
		return ctx.Error(err)
	}
	userID, _ := auth.UserIDFromCtx(ctx)
	teamPass := util.RandPass(20, 3, 5)
	res, err := a.team.Create(userID, teamPass, &schema.Team{
		Name:       body.Name,
		UserID:     userID,
		Visibility: "member",
		Members:    []string{userID},
	})

	if err != nil {
		return ctx.Error(err)
	}

	if err = a.file.SaveTeamPass(res.ID, userID, key, teamPass); err != nil {
		return ctx.Error(err)
	}
	return ctx.JSON(200, res)
}

type tplTeamUpdate map[string]interface{}

// Validate ...
func (t *tplTeamUpdate) Validate() error {
	empty := true
	for key, val := range *t {
		empty = false

		switch key {
		case "name":
			v, ok := val.(string)
			if !ok || v == "" {
				return &gear.Error{Code: 400, Msg: "invalid team name"}
			}
		case "isFrozen":
			_, ok := val.(bool)
			if !ok {
				return &gear.Error{Code: 400, Msg: "invalid team isFrozen"}
			}
		default:
			return &gear.Error{Code: 400, Msg: "invalid team property"}
		}
	}

	if empty {
		return &gear.Error{Code: 400, Msg: "no content"}
	}
	return nil
}

// Update ...
//
// @Title Update
// @Summary Update the team
// @Description only the team owner can update the team
// @Param Authorization header string true "access_token"
// @Param teamID path string true "team ID"
// @Param body body tplTeamUpdate true "team body"
// @Success 200 schema.TeamResult
// @Failure 400 string
// @Failure 401 string
// @Router PUT /api/teams/{teamID}
func (a *Team) Update(ctx *gear.Context) (err error) {
	TeamID, err := util.ParseOID(ctx.Param("teamID"))
	if err != nil {
		return ctx.ErrorStatus(400)
	}

	userID, _ := auth.UserIDFromCtx(ctx)
	body := new(tplTeamUpdate)
	if err = ctx.ParseBody(body); err != nil {
		return ctx.Error(err)
	}

	team, err := a.team.Find(TeamID, false)
	if err != nil {
		return ctx.Error(err)
	}
	if team.UserID != userID {
		return ctx.ErrorStatus(403)
	}

	changed := false
	for key, val := range *body {
		switch key {
		case "name":
			if name := val.(string); name != team.Name {
				changed = true
				team.Name = name
			}
		case "isFrozen":
			if isFrozen := val.(bool); isFrozen != team.IsFrozen {
				changed = true
				team.IsFrozen = isFrozen
			}
		}
	}

	if !changed {
		return ctx.End(204)
	}

	res, err := a.team.Update(TeamID, team)
	if err != nil {
		return ctx.Error(err)
	}
	return ctx.JSON(200, res)
}

type tplTeamMembers struct {
	Push []string `json:"$push" swaggo:"false,add some team members,[\"joe\"]"`
	Pull []string `json:"$pull" swaggo:"false,remove some team members,[]"`
}

// Validate ...
func (t *tplTeamMembers) Validate() error {
	if len(t.Push) == 0 && len(t.Pull) == 0 {
		return &gear.Error{Code: 400, Msg: "no content"}
	}
	if len(t.Push) > 100 || len(t.Pull) > 100 {
		return &gear.Error{Code: 400, Msg: "too many members"}
	}
	return nil
}

// Members ...
//
// @Title Members
// @Summary Add or remove team members
// @Description only the team owner can update the team members
// @Param Authorization header string true "access_token"
// @Param teamID path string true "team ID"
// @Param body body tplTeamMembers true "team members"
// @Success 200 schema.TeamResult
// @Failure 400 string
// @Failure 401 string
// @Router PUT /api/teams/{teamID}/members
func (a *Team) Members(ctx *gear.Context) (err error) {
	TeamID, err := util.ParseOID(ctx.Param("teamID"))
	if err != nil {
		return ctx.ErrorStatus(400)
	}

	userID, _ := auth.UserIDFromCtx(ctx)
	body := new(tplTeamMembers)
	if err = ctx.ParseBody(body); err != nil {
		return ctx.Error(err)
	}

	res, err := a.team.UpdateMembers(userID, TeamID, body.Pull, body.Push)
	if err != nil {
		return ctx.Error(err)
	}
	return ctx.JSON(200, res)
}

// Delete ...
//
// @Title Delete
// @Summary Delete the team
// @Description only the team owner can delete the team
// @Param Authorization header string true "access_token"
// @Param teamID path string true "team ID"
// @Success 204
// @Failure 400 string
// @Failure 401 string
// @Router DELETE /api/entries/{teamID}
func (a *Team) Delete(ctx *gear.Context) (err error) {
	TeamID, err := util.ParseOID(ctx.Param("teamID"))
	if err != nil {
		return ctx.ErrorStatus(400)
	}

	userID, _ := auth.UserIDFromCtx(ctx)
	team, err := a.team.Find(TeamID, false)
	if err != nil {
		return ctx.Error(err)
	}
	if team.UserID != userID {
		return ctx.ErrorStatus(403)
	}
	if team.Visibility == "private" {
		return ctx.Error(&gear.Error{Code: 403, Msg: "private team can't be deleted"})
	}

	team.IsDeleted = true
	if _, err = a.team.Update(TeamID, team); err != nil {
		return ctx.Error(err)
	}
	return ctx.End(204)
}

// Undelete ...
//
// @Title Undelete
// @Summary Undelete the team
// @Description only the team owner can undelete the team
// @Param Authorization header string true "access_token"
// @Param teamID path string true "entry ID"
// @Success 204
// @Failure 400 string
// @Failure 401 string
// @Router POST /api/teams/{teamID}:undelete
func (a *Team) Undelete(ctx *gear.Context) (err error) {
	TeamID, err := util.ParseOID(ctx.Param("teamID"))
	if err != nil {
		return ctx.ErrorStatus(400)
	}

	userID, _ := auth.UserIDFromCtx(ctx)
	team, err := a.team.Find(TeamID, true)
	if err != nil {
		return ctx.Error(err)
	}
	if team.UserID != userID {
		return ctx.ErrorStatus(403)
	}

	team.IsDeleted = false
	res, err := a.team.Update(TeamID, team)
	if err != nil {
		return ctx.Error(err)
	}
	return ctx.JSON(200, res)
}

// FindByMember return teams for current user
//
// @Title FindByMember
// @Summary Get teams for current user
// @Description Get teams for current user.
// @Param Authorization header string true "access_token"
// @Success 200 []schema.TeamResult
// @Failure 400 string
// @Failure 401 string
// @Router GET /api/teams
func (a *Team) FindByMember(ctx *gear.Context) (err error) {
	userID, err := auth.UserIDFromCtx(ctx)
	if err != nil {
		return ctx.Error(err)
	}
	teams, err := a.team.FindByMemberID(userID)
	if err != nil {
		return ctx.Error(err)
	}
	return ctx.JSON(200, teams)
}

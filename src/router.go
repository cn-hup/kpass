package src

import (
	"github.com/seccom/kpass/src/api"
	"github.com/seccom/kpass/src/auth"
	"github.com/seccom/kpass/src/bll"
	"github.com/seccom/kpass/src/ctl"
	"github.com/seccom/kpass/src/service"
	"github.com/teambition/gear"
)

func noOp(ctx *gear.Context) error {
	return gear.ErrNotFound.WithMsg("noOp")
}

func newRouter(db *service.DB) (Router *gear.Router) {
	Router = gear.NewRouter()

	blls := new(bll.Blls).Init(db)
	apis := new(api.APIs).Init(blls)
	fileCtl := new(ctl.File).Init(blls)

	// GET /download/fileID?refType=user&refID=userID
	// GET /download/fileID?refType=team&refID=teamID
	// GET /download/fileID?refType=entry&refID=entryID&signed=xxxx
	Router.Get("/download/:fileID", fileCtl.Download)

	Router.Post("/upload/avatar", auth.Middleware, fileCtl.UploadAvatar)
	Router.Post("/upload/team/:teamID/logo", auth.Middleware, fileCtl.UploadLogo)
	Router.Post("/upload/entry/:entryID/file", auth.Middleware, fileCtl.UploadFile)

	// generate a random password
	Router.Get("/api/password", apis.User.Password)

	// Create a new user
	Router.Post("/api/join", apis.User.Join)
	Router.Post("/api/login", apis.User.Login)
	// Return the user publicly info
	Router.Get("/api/user/:userID", auth.Middleware, apis.User.Find)

	// Create a team
	Router.Post("/api/teams", auth.Middleware, apis.Team.Create)
	// // Return current user's teams joined.
	Router.Get("/api/teams", auth.Middleware, apis.Team.FindByMember)
	// Join a team by invite code
	Router.Post("/api/teams/join", auth.Middleware, apis.Team.Join)
	// Undelete the team
	Router.Post(`/api/teams/:teamID+:undelete`, auth.Middleware, apis.Team.Undelete)
	// Return the team's entries list
	Router.Get("/api/teams/:teamID/entries", auth.Middleware, apis.Entry.FindByTeam)
	// Create a new entry for team
	Router.Post("/api/teams/:teamID/entries", auth.Middleware, apis.Entry.Create)
	// Invite a user to the team
	Router.Post("/api/teams/:teamID/invite", auth.Middleware, apis.Team.Invite)
	// Update the team
	Router.Put("/api/teams/:teamID", auth.Middleware, apis.Team.Update)
	// remove the team's member
	Router.Delete("/api/teams/:teamID/members/:userID", auth.Middleware, apis.Team.RemoveMember)
	// Delete the team
	Router.Delete("/api/teams/:teamID", auth.Middleware, apis.Team.Delete)
	// Return the team's shares list
	Router.Get("/api/teams/:teamID/shares", auth.Middleware, noOp)

	// Undelete the entry
	Router.Post("/api/entries/:entryID+:undelete", auth.Middleware, apis.Entry.Undelete)
	// Get the full entry, with all secrets
	Router.Get("/api/entries/:entryID", auth.Middleware, apis.Entry.Find)
	// Update the entry
	Router.Put("/api/entries/:entryID", auth.Middleware, apis.Entry.Update)
	// Delete the entry
	Router.Delete("/api/entries/:entryID", auth.Middleware, apis.Entry.Delete)
	// Add a secret to the entry
	Router.Post("/api/entries/:entryID/secrets", auth.Middleware, apis.Secret.Create)
	// Update the secret
	Router.Put("/api/entries/:entryID/secrets/:secretID", auth.Middleware, apis.Secret.Update)
	// Delete the secret
	Router.Delete("/api/entries/:entryID/secrets/:secretID", auth.Middleware, apis.Secret.Delete)
	// Add a share to the entry
	Router.Post("/api/entries/:entryID/shares", auth.Middleware, apis.Share.Create)
	// Get shares list of the entry
	// Router.Get("/api/entries/:entryID/shares", auth.Middleware, shareAPI.FindByEntry)
	// Delete the share
	// Router.Delete("/api/entries/:entryID/shares/:shareID", auth.Middleware, apis.Entry.DeleteShare)

	// Delete the file
	Router.Delete("/api/entries/:entryID/files/:fileID", auth.Middleware, apis.Entry.DeleteFile)

	// Get shares list of the team
	// Router.Get("/api/teams/:teamID/shares", auth.Middleware, shareAPI.FindByTeam)
	// Get shares list to the current user
	// Router.Get("/api/shares/me", auth.Middleware, shareAPI.FindByUser)
	// Get the share
	// Router.Get("/api/shares/:shareID", auth.Middleware, shareAPI.ReadShare)

	// Router.Get("/api/debug", func(ctx *gear.Context) error {
	// 	user, err := user.Current()
	// 	if err != nil {
	// 		panic(fmt.Errorf("Unable to determine user's home directory: %s", err))
	// 	}
	// 	pick := func(x, y interface{}) interface{} {
	// 		return x
	// 	}

	// 	return ctx.JSON(200, map[string]interface{}{
	// 		"Uid":        user.Uid,
	// 		"Gid":        user.Gid,
	// 		"UserName":   user.Username,
	// 		"Name":       user.Name,
	// 		"HomeDir":    user.HomeDir,
	// 		"Env":        os.Environ(),
	// 		"Executable": pick(os.Executable()),
	// 		"Hostname":   pick(os.Hostname()),
	// 		"WD":         pick(os.Getwd()),
	// 	})
	// })
	return
}

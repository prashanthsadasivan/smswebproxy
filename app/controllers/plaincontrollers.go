package controllers

import (
	"github.com/revel/revel"
	"smswebproxy/app/appconnections"
	"smswebproxy/app/models"
)

type App struct {
	GorpController
}

// ========== BEGIN PLAIN CONTROLLERS ==========

func (c App) Room(num string) revel.Result {
	return c.Render(num)
}

func (c App) Home() revel.Result {
	return c.Render()
}

// ==========  END PLAIN CONTROLLERS   ==========

func (c App) GorcAppConnection(num string) *appconnections.AppConnection {
	appconnection := appconnections.GetAppConnection(num)
	if appconnection == nil {
		soleUser := models.GetSoleUser(c.Txn)
		appconnection = appconnections.New(soleUser.GcmId, soleUser.Number)
	}

	return appconnection
}

func (c App) Send(sender, num_to, msgToTextOut string) revel.Result {
	return nil
}

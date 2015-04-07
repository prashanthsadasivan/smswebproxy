package controllers

import (
	"github.com/revel/revel"
)

func (c App) Ping() revel.Result {
	return c.RenderText("pong")
}

package controllers

import (
	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Ping() revel.Result {
	return c.RenderText("PONG")
}


func (c App) Index() revel.Result {
	return c.Render()
}
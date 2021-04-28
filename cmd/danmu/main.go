package main

import (
	"bytes"
	"image/color"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"golang.org/x/image/font/gofont/gomono"
)

type Text struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

type Scene struct{}

func (*Scene) Type() string { return "Scene" }

func (*Scene) Preload() {
	// TTF, err := ioutil.ReadFile("assets/fonts/OpenSans-Semibold.ttf")
	// if err != nil {
	// 	log.Fatal("Unable to read TTF font: ", err)
	// }

	engo.Files.LoadReaderData("go.ttf", bytes.NewReader(gomono.TTF))
}

func (*Scene) Setup(u engo.Updater) {
	world, _ := u.(*ecs.World)
	world.AddSystem(&common.RenderSystem{})

	fnt := &common.Font{
		URL:  "go.ttf",
		FG:   color.White,
		Size: 24,
	}
	fnt.CreatePreloaded()

	text := Text{BasicEntity: ecs.NewBasic()}
	text.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: 10, Y: 10},
		Width:    200,
		Height:   200,
	}

	text.RenderComponent = common.RenderComponent{
		Drawable: common.Text{
			Font: fnt,
			Text: "Hello World",
		},
	}
	text.RenderComponent.SetZIndex(1001)

	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&text.BasicEntity, &text.RenderComponent, &text.SpaceComponent)
		}
	}
}

func main() {
	opts := engo.RunOptions{
		Title:  "Hello World",
		Width:  400,
		Height: 400,
	}

	engo.Run(opts, &Scene{})
}

package main

import (
	"image/color"
	"log"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type Text struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

type Scene struct{}

func (*Scene) Type() string { return "Scene" }

func (*Scene) Preload() {
	err := engo.Files.Load("fonts/NotoSansSC-Regular.ttf")
	if err != nil {
		log.Fatal("Unable to load font: ", err)
	}
}

func (*Scene) Setup(u engo.Updater) {
	world, _ := u.(*ecs.World)
	world.AddSystem(&common.RenderSystem{})
	world.AddSystem(&common.AnimationSystem{})

	common.UnicodeCap = 40000

	fnt := &common.Font{
		URL:  "fonts/NotoSansSC-Regular.ttf",
		FG:   color.White,
		Size: 24,
	}
	err := fnt.CreatePreloaded()
	if err != nil {
		log.Fatal("Unable to preload: ", err)
	}

	text := Text{BasicEntity: ecs.NewBasic()}
	text.RenderComponent.Drawable = common.Text{
		Font: fnt,
		Text: "我是凯阳",
	}
	text.SetShader(common.HUDShader)

	text.RenderComponent.SetZIndex(1001)
	text.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: 10, Y: 10},
		Width:    200,
		Height:   200,
	}

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

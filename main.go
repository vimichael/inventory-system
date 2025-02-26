package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	Inventory *Inventory
	AtlasImg  *ebiten.Image
}

func NewGame() *Game {
	img, _, err := ebitenutil.NewImageFromFile("assets/atlas.png")
	if err != nil {
		log.Fatal(err)
	}

	inventory := NewInventory()
	inventory.Hand.ItemId = 1
	inventory.Hand.Amount = 10
	inventory.PlaceItems(1)
	inventory.Hand.ItemId = 2
	inventory.Hand.Amount = 10
	inventory.PlaceItems(12)

	return &Game{
		Inventory: inventory,
		AtlasImg:  img,
	}
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		return ebiten.Termination
	}
	g.Inventory.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{120, 180, 255, 255})
	g.Inventory.Draw(screen, g.AtlasImg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("window")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}

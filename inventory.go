package main

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const EmptyCell = 0

const (
	rowSize  = 9
	cellSize = 64
	// gap between
	gapSize = 2
	// rendering offset so inventory isnt at the topleft
	offsetX = 200
	offsetY = 200
)

type Cell struct {
	ItemId int
	Amount uint
	// used for both rendering and mouse hit detection
	Hitbox image.Rectangle
}

type Inventory struct {
	Cells []*Cell
	// the hand acts as a temporary cell
	Hand Cell
}

func GetHitbox(index int) image.Rectangle {
	gridX := index % rowSize
	gridY := index / rowSize
	// grid_pos * grid_pixel_size + render_offset
	pixelX := gridX*(cellSize+gapSize) + offsetX
	pixelY := gridY*(cellSize+gapSize) + offsetY

	// (x0, y0, x1, y1)
	return image.Rect(pixelX, pixelY, (pixelX + cellSize), pixelY+cellSize)
}

func NewInventory() *Inventory {
	cells := make([]*Cell, 27)
	for i := range cells {
		cells[i] = &Cell{
			Hitbox: GetHitbox(i),
		}
	}

	return &Inventory{
		Cells: cells,
		Hand:  Cell{},
	}
}

func (i *Inventory) Update() {
	mouseX, mouseY := ebiten.CursorPosition()
	// update the pos of the hand to the cursor
	i.Hand.Hitbox = image.Rect(mouseX, mouseY, (mouseX + cellSize), mouseY+cellSize)

	mouseClicked := inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0)
	rightMouseClicked := inpututil.IsMouseButtonJustPressed(ebiten.MouseButton2)

	if mouseClicked {
		// check if clicked cell
		mouseP := image.Point{mouseX, mouseY}
		// do something
		for index, cell := range i.Cells {
			if mouseP.In(cell.Hitbox) {
				i.PlaceItems(index)
			}
		}
	}

	if rightMouseClicked {
		// check if clicked cell
		mouseP := image.Point{mouseX, mouseY}
		// do something
		for index, cell := range i.Cells {
			if mouseP.In(cell.Hitbox) {
				i.PlaceOneItem(index)
			}
		}
	}
}

func (i *Inventory) DrawHand(screen *ebiten.Image, atlas *ebiten.Image) {
	if i.Hand.ItemId != EmptyCell {

		id := i.Hand.ItemId - 1
		// the number of items per row in the image. the magical 8.0 is the
		// size in pixels of an item. Change this to your tile size.
		numItemsPerRow := atlas.Bounds().Dx() / 8.0
		gX := id % numItemsPerRow
		gY := id / numItemsPerRow

		// convert the grid coord to pixel coord and add the tilesize
		// to get the bottom right coord
		srect := image.Rect(
			gX*8.0,
			gY*8.0,
			(gX+1)*8.0,
			(gY+1)*8.0,
		)

		ops := ebiten.DrawImageOptions{}
		ops.GeoM.Scale(8.0, 8.0)
		ops.GeoM.Translate(
			float64(i.Hand.Hitbox.Min.X-(cellSize/2)),
			float64(i.Hand.Hitbox.Min.Y)-(cellSize/2),
		)

		screen.DrawImage(
			atlas.SubImage(srect).(*ebiten.Image),
			&ops,
		)

		ebitenutil.DebugPrintAt(
			screen,
			fmt.Sprintf("%d", i.Hand.Amount),
			i.Hand.Hitbox.Min.X-(cellSize/2),
			i.Hand.Hitbox.Min.Y-(cellSize/2),
		)
	}
}

func (i *Inventory) Draw(screen *ebiten.Image, atlas *ebiten.Image) {
	for _, cell := range i.Cells {
		// always draw the cell
		vector.DrawFilledRect(
			screen,
			float32(cell.Hitbox.Min.X),
			float32(cell.Hitbox.Min.Y),
			float32(cell.Hitbox.Dx()),
			float32(cell.Hitbox.Dy()),
			color.RGBA{20, 20, 20, 255},
			true,
		)

		if cell.ItemId != EmptyCell {
			// since 1 = first item, we need to subtract 1 to get the coord.
			// this is because if index = 1, index % items_per_row = 1, not 0.
			id := cell.ItemId - 1
			numItemsPerRow := atlas.Bounds().Dx() / 8.0
			gX := id % numItemsPerRow
			gY := id / numItemsPerRow

			// the number of items per row in the image. the magical 8.0 is the
			// size in pixels of an item. Change this to your tile size.
			srect := image.Rect(
				gX*8.0,
				gY*8.0,
				(gX+1)*8.0,
				(gY+1)*8.0,
			)

			ops := ebiten.DrawImageOptions{}
			ops.GeoM.Scale(8.0, 8.0)
			ops.GeoM.Translate(float64(cell.Hitbox.Min.X), float64(cell.Hitbox.Min.Y))

			screen.DrawImage(
				atlas.SubImage(srect).(*ebiten.Image),
				&ops,
			)

			ebitenutil.DebugPrintAt(
				screen,
				fmt.Sprintf("%d", cell.Amount),
				cell.Hitbox.Min.X,
				cell.Hitbox.Min.Y,
			)
		}
	}

	i.DrawHand(screen, atlas)
}

func (i *Inventory) PrintInventory() {
	for _, cell := range i.Cells {
		fmt.Printf("cell: %d, %d\n", cell.ItemId, cell.Amount)
	}
	fmt.Printf("Hand: %d %d\n", i.Hand.ItemId, i.Hand.Amount)
}

// should probably call this `HandleRightClick`
func (i *Inventory) PlaceOneItem(index int) bool {
	// validate the index
	if index < 0 || index > len(i.Cells)-1 {
		return false
	}

	// grab a pointer to the target cell (don't copy!!!)
	cell := i.Cells[index]

	// check if we can merge
	if cell.ItemId == i.Hand.ItemId {
		cell.Amount += 1
		i.Hand.Amount--
		if i.Hand.Amount < 1 {
			i.Hand.ItemId = EmptyCell
		}
		return true
	}

	// check if cell is empty
	if cell.ItemId == EmptyCell {
		cell.ItemId = i.Hand.ItemId
		cell.Amount += 1
		i.Hand.Amount--
		if i.Hand.Amount < 1 {
			i.Hand.ItemId = EmptyCell
		}
		return true
	}

	// splitting stack in half
	if i.Hand.ItemId == EmptyCell && cell.Amount > 1 {
		half := cell.Amount / 2
		i.Hand.ItemId = cell.ItemId
		i.Hand.Amount = half
		cell.Amount -= half
		return true
	}

	// dont match and cell isnt empty, so swap
	temp := cell.ItemId
	cell.ItemId = i.Hand.ItemId
	i.Hand.ItemId = temp

	tempAmt := cell.Amount
	cell.Amount = i.Hand.Amount
	i.Hand.Amount = tempAmt

	return true
}

// should probably call this `HandleLeftClick`
func (i *Inventory) PlaceItems(index int) bool {
	// validate the index
	if index < 0 || index > len(i.Cells)-1 {
		return false
	}

	cell := i.Cells[index]

	// check if we can merge
	if cell.ItemId == i.Hand.ItemId {
		cell.Amount += i.Hand.Amount
		i.Hand.Amount = 0
		i.Hand.ItemId = EmptyCell
		return true
	}

	// check if cell is empty
	if cell.ItemId == EmptyCell {
		cell.ItemId = i.Hand.ItemId
		i.Hand.ItemId = EmptyCell
		cell.Amount = i.Hand.Amount
		i.Hand.Amount = 0
		return true
	}

	// dont match and cell isnt empty, so swap
	temp := cell.ItemId
	cell.ItemId = i.Hand.ItemId
	i.Hand.ItemId = temp

	tempAmt := cell.Amount
	cell.Amount = i.Hand.Amount
	i.Hand.Amount = tempAmt

	return true
}

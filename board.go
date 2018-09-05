package main

type Board struct {
	OffsetX, OffsetY int
	Rows             int
	Columns          int
	Positions        []Position
}

type Position struct {
	PTile *Tile
}

type Coord struct {
	Row    int
	Column int
}

type Tile struct {
	Type     string
	Row, Col int
	X, Y     int
}

func (b *Board) AddTile(x, y int) {
	tx, ty := fit_grid(x, y)
	iscale := int(scale)
	r, c := ty/iscale, tx/iscale
	ox, oy := b.OffsetX, b.OffsetY
	r, c = r-oy, c-ox
	pos := b.Position(Coord{r, c})
	if pos != nil {
		pos.PTile = &Tile{"grass", r, c, tx, ty}
	}
}

func fit_grid(x, y int) (tx, ty int) {
	nxt := offset * 2
	tx, ty = (x/nxt)*nxt, (y/nxt)*nxt
	return
}

func (b *Board) Position(c Coord) *Position {
	if c.Row < 0 || c.Row >= b.Rows || c.Column < 0 || c.Column >= b.Columns {
		return nil
	}
	return &b.Positions[c.Row*b.Columns+c.Column]
}

package main

import (
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/nsf/termbox-go"
)

type Cell struct {
	X    int
	Y    int
	Dead bool
}

type GameOfLife struct {
	cells      [][]Cell
	Generation int
}

type boolgen struct {
	src       rand.Source
	cache     int64
	remaining int
}

func (b *boolgen) Bool() bool {
	if b.remaining == 0 {
		b.cache, b.remaining = b.src.Int63(), 125
	}

	result := b.cache&0x01 == 1
	b.cache >>= 1
	b.remaining--

	return result
}

var (
	maxRow, maxCol int
)

func New() *boolgen {
	return &boolgen{src: rand.NewSource(time.Now().UnixNano())}
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	maxRow, maxCol = termbox.Size()
	maxCol -= 5
	rand.Seed(time.Now().UTC().UnixNano())

	b := initiate()
	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()
loop:
	for {
		select {
		case ev := <-eventQueue:
			if ev.Key == termbox.KeyCtrlC {
				break loop
			}
			if ev.Key == termbox.KeyCtrlN {
				b = initiate()
			}
		default:
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			b.NextGen()
			b.Print()
			termbox.Flush()
			time.Sleep(1 * time.Millisecond)
		}
	}
}

func initiate() *GameOfLife {
	b := &GameOfLife{
		cells: make([][]Cell, int(math.Max(float64(maxCol), float64(maxRow)))),
	}
	for i := 0; i <= maxCol; i++ {
		for j := 0; j <= maxRow; j++ {
			dead := true
			if randInt(0, 100) < 10 {
				dead = false
			}
			c := Cell{
				X:    i,
				Y:    j,
				Dead: dead,
			}
			b.cells[i] = append(b.cells[i], c)
		}
	}

	return b
}
func (b *GameOfLife) te() {
	b.cells[0][0].Dead = true
}
func (b *GameOfLife) NextGen() {
	duplicate := make([][]Cell, len(b.cells))
	for i := range b.cells {
		duplicate[i] = make([]Cell, len(b.cells[i]))
		copy(duplicate[i], b.cells[i])
	}

	for i := 0; i <= maxCol; i++ {
		for j := 0; j <= maxRow; j++ {
			ncnt := b.Nighbors(b.cells[i][j])
			if ncnt < 2 || ncnt > 3 {
				duplicate[i][j].Dead = true
			}
			if b.cells[i][j].Dead && ncnt == 3 {
				duplicate[i][j].Dead = false
			}
			if ncnt == 2 {
				duplicate[i][j] = b.cells[i][j]
			}
		}
	}
	b.Generation++
	b.cells = duplicate
}

func (b *GameOfLife) Nighbors(cell Cell) int {
	cnt := 0
	for x1 := cell.X - 1; x1 <= cell.X+1; x1++ {
		if x1 < 0 || x1 >= maxCol {
			continue
		}
		for y1 := cell.Y - 1; y1 <= cell.Y+1; y1++ {
			if x1 == cell.X && y1 == cell.Y {
				continue
			}
			if y1 < 0 || y1 >= maxRow {
				continue
			}
			if !b.cells[x1][y1].Dead {
				cnt++
			}
		}
	}
	return cnt
}

func (b *GameOfLife) Print() {
	for i := 0; i <= maxCol; i++ {
		for j := 0; j <= maxRow; j++ {
			ch := '*'
			if b.cells[i][j].Dead {
				ch = ' '
			}
			termbox.SetCell(j, i, ch, termbox.ColorDefault, termbox.ColorDefault)
		}
	}

	status := "Generation: " + strconv.Itoa(b.Generation)
	for x, ch := range status {
		termbox.SetCell(x+1, maxCol-2, ch, termbox.ColorDefault, termbox.ColorDefault)
	}
}

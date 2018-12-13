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
	Alives     int
}

var (
	maxRow, maxCol int
	pause          bool
	speed          time.Duration
	random         bool
)

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func main() {
	pause = true
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	speed = 100
	random = false
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputMouse)
	maxRow, maxCol = termbox.Size()
	maxCol -= 3
	rand.Seed(time.Now().UTC().UnixNano())
	b := initiate(false)
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
			if ev.Type == termbox.EventKey {
				if ev.Key == termbox.KeyCtrlC {
					break loop
				}

				if ev.Key == termbox.KeyCtrlS {
					pause = !pause
				}

				if ev.Key == termbox.KeyCtrlN {
					b = initiate(false)
					pause = true
				}

				if ev.Key == termbox.KeyCtrlI {
					speed++
				}

				if ev.Key == termbox.KeyCtrlD {
					speed--
				}

				if ev.Key == termbox.KeyCtrlR {
					b = initiate(true)
				}
			}
			if ev.Type == termbox.EventMouse {
				if ev.MouseX <= maxRow && ev.MouseY <= maxCol {
					b.cells[ev.MouseY][ev.MouseX].Dead = false
				}
			}
			b.Print()
			termbox.Flush()

		default:
			if !pause {
				b.NextGen()
				time.Sleep(speed * time.Millisecond)
			}
			b.Print()
			termbox.Flush()
		}
	}
}

func initiate(random bool) *GameOfLife {
	b := &GameOfLife{
		cells: make([][]Cell, int(math.Max(float64(maxCol), float64(maxRow)))),
	}
	for i := 0; i <= maxCol; i++ {
		for j := 0; j <= maxRow; j++ {
			dead := true
			if random {
				if randInt(0, 100) < 50 {
					dead = false
				}
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
	b.Alives = 0
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

			if !duplicate[i][j].Dead {
				b.Alives++
			}
		}
	}
	b.Generation++
	b.cells = duplicate
}

func (b *GameOfLife) Nighbors(cell Cell) int {
	cnt := 0
	for x1 := cell.X - 1; x1 <= cell.X+1; x1++ {
		for y1 := cell.Y - 1; y1 <= cell.Y+1; y1++ {
			if x1 == cell.X && y1 == cell.Y {
				continue
			}
			if x1 < 0 || x1 >= maxCol {
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
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	for i := 0; i <= maxCol; i++ {
		for j := 0; j <= maxRow; j++ {
			ch := '█'
			if b.cells[i][j].Dead {
				ch = '░'
			}
			termbox.SetCell(j, i, ch, termbox.ColorDefault, termbox.ColorDefault)
		}
	}
	b.PrintStatus()
}

func (b *GameOfLife) PrintStatus() {
	statuses := []string{
		"Generation: " + strconv.Itoa(b.Generation) + " | Alives: " + strconv.Itoa(b.Alives) + " | Speed: " + (speed * time.Microsecond).String(),
		"Shortcuts: New (Ctrl+N) | Start/Pause(Ctrl+S) | Random (Ctrl+R) | Set Cell Alive (Mouse Click) | Speed Up (Ctrl+D) | Speed Down (Ctrl+I)",
	}
	for c, status := range statuses {
		for x, ch := range status {
			termbox.SetCell(x+1, maxCol+c+1, ch, termbox.ColorCyan, termbox.ColorDefault)
		}
	}
}

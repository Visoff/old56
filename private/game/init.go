package game

import (
	"fmt"
	"strings"

	"github.com/gliderlabs/ssh"
)

type RenderEvent func(*Game, *[][]Pixel) error
type UpdateEvent func(*Game, []byte) error
type StartEvent func(*Game) error

type GameFactory struct {
	renderers []RenderEvent
	updates   []UpdateEvent
    starters  []StartEvent
}

func NewGameFactory() *GameFactory {
    return &GameFactory{
        renderers: make([]RenderEvent, 0),
        updates:   make([]UpdateEvent, 0),
    }
}

func (gf *GameFactory) AddUpdateListener(f UpdateEvent) {
	gf.updates = append(gf.updates, f)
}

func (gf *GameFactory) AddRenderObject(f RenderEvent) {
	gf.renderers = append(gf.renderers, f)
}

func (gf *GameFactory) AddStartEvent(f StartEvent) {
    gf.starters = append(gf.starters, f)
}

type Game struct {
	S       ssh.Session
	factory *GameFactory
    State   map[string]interface{}
    NeedUpdate bool
}

func (gf *GameFactory) New(s ssh.Session) *Game {
    g := &Game{
		S:       s,
		factory: gf,
        State:   make(map[string]interface{}),
        NeedUpdate: true,
	}
    for _, f := range gf.starters {
        err := f(g)
        if err != nil {
            panic(err)
        }
    }
    return g
}

func (g *Game) Close() {
	g.S.Close()
}

type Pixel struct {
	ch   rune
	tags []int
}

func NewPixel(ch rune, tags []int) Pixel {
    return Pixel{ch: ch, tags: tags}
}

func (g *Game) Update() error {
    g.NeedUpdate = false
	data := make([]byte, 8)
	g.S.Read(data)
	if len(data) == 0 {
		return nil
	}
	for _, u := range g.factory.updates {
		err := u(g, data)
		if err != nil {
			return err
		}
	}
    g.NeedUpdate = true
	return nil
}

func (g *Game) Render(w, h int) (string, error) {
	data := make([][]Pixel, h)
	for i := range data {
		data[i] = make([]Pixel, w)
		for j := range data[i] {
			data[i][j] = Pixel{ch: ' ', tags: make([]int, 0)}
		}
	}
	for _, r := range g.factory.renderers {
		err := r(g, &data)
		if err != nil {
			return "", err
		}
	}
	res := make([]string, h)
	for i := range res {
		res[i] = "\033[0m"
		for j := range data[i] {
			for _, t := range data[i][j].tags {
				res[i] += "\033[" + fmt.Sprint(t) + "m"
			}
			res[i] += string(data[i][j].ch)
		}
	}
	return strings.Join(res, "\n"), nil
}

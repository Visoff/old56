package main

import (
	"fmt"

	"github.com/Visoff/old56/private/game"
	"github.com/Visoff/old56/private/server"
)


func main() {
    gf := game.NewGameFactory()
    gf.AddUpdateListener(func(g *game.Game, b []byte) error {
        if b[0] == 'q' {
            g.S.Exit(0)
            return nil
        }
        if '1' <= b[0] && b[0] <= '3' {
            g.State["active_pane"] = int(b[0] - '0')
        }
        return nil
    })
    gf.AddRenderObject(func(g *game.Game, p *[][]game.Pixel) error {
        i := len(*p)-2
        for j := range (*p)[i] {
            (*p)[i][j] = game.NewPixel(
                '-',
                []int{3},
            )
        }
        i++
        money := fmt.Sprintf("%d$", g.State["money"].(int))
        j := 0
        for _, ch := range money {
            (*p)[i][j] = game.NewPixel(
                ch,
                []int{},
            )
            j++
        }
        return nil
    })
    tag_if_selected := func (sel, curr int) int {
        if sel == curr {
            return 92
        }
        return 0
    }
    gf.AddRenderObject(func(g *game.Game, p *[][]game.Pixel) error {
        w := (len((*p)[0]) - 4 - 4) / 3
        h := len(*p) - 2 - 2 - 2
        selected_pane := g.State["active_pane"].(int)
        game.Rect(2, 2, w, h, []int{tag_if_selected(1, selected_pane)}, p)
        game.Rect(4+w, 2, w, h, []int{tag_if_selected(2, selected_pane)}, p)
        game.Rect(6+2*w, 2, w, h, []int{tag_if_selected(3, selected_pane)}, p)
        return nil
    })
    gf.AddStartEvent(func(g *game.Game) error {
        g.State["money"] = 0
        g.State["workers"] = 0
        g.State["active_pane"] = 1
        return nil
    })
    err := server.Init(gf)
    if err != nil {
        panic(err)
    }
}

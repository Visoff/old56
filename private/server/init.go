package server

import (
	"io"
	"time"

	"github.com/Visoff/old56/private/game"
	"github.com/gliderlabs/ssh"
)

func Init(gf *game.GameFactory) error {
    ssh.Handle(func(s ssh.Session) {
        io.WriteString(s, "\033[H\033[2J")
        g := gf.New(s)
        p, window, _ := s.Pty()
        w, h := p.Window.Width, p.Window.Height
        defer g.Close()
        ticker := time.NewTicker(time.Millisecond * 100)
        for {
            select {
            case win := <-window:
                w, h = win.Width, win.Height
            case <-s.Context().Done():
                return
            case <-ticker.C:
                if g.NeedUpdate {
                    go func () {
                        g.Update()
                    }()
                }
                data, err := g.Render(w, h)
                if err != nil {
                    s.Exit(1)
                    return
                }
                io.WriteString(s, data)
            }
        }
    })

    return ssh.ListenAndServe(":8080", nil)
}

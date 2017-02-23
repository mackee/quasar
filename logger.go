package quasar

import (
	"bufio"
	"fmt"
	"io"
	"sync"
	"time"

	colorable "github.com/mattn/go-colorable"
)

type color int

const (
	red color = iota + 31
	green
	yellow
	blue
	magenta
	cyan

	colornum    = 5
	coloroffset = 31
)

var (
	logger      = colorable.NewColorableStderr()
	isColorable = true
	logcnt      = 0
	logmux      = sync.Mutex{}
)

func setColorable(b bool) {
	isColorable = b
}

func logging(r io.Reader, name string) {
	if c, ok := r.(io.Closer); ok {
		defer c.Close()
	}
	sc := bufio.NewScanner(r)

	logmux.Lock()
	clr := color(logcnt%colornum + coloroffset)
	logid := logcnt
	logcnt++
	logmux.Unlock()

	for sc.Scan() {
		hour, min, sec := time.Now().Clock()
		if isColorable {
			fmt.Fprintf(logger, "\x1b[%dm%02d:%02d:%02d %-10s %3d |\x1b[0m %s\n",
				clr,
				hour, min, sec,
				name, logid,
				sc.Text(),
			)
		} else {
			fmt.Fprintf(logger, "%02d:%02d:%02d %-10d %3d | %s\n",
				hour, min, sec,
				name, logid,
				sc.Text(),
			)
		}
	}
}

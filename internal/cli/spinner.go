package cli

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type Spinner struct {
	writer  io.Writer
	message string
	frames  []string
	stop    chan struct{}
	done    chan struct{}
	once    sync.Once
	enabled bool
}

func NewSpinner(
	writer io.Writer,
	message string,
	enabled bool,
) *Spinner {
	return &Spinner{
		writer:  writer,
		message: message,
		frames:  []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧"},
		stop:    make(chan struct{}),
		done:    make(chan struct{}),
		enabled: enabled,
	}
}

func (s *Spinner) Start() {
	if !s.enabled {
		return
	}

	go func() {
		defer close(s.done)

		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		frame := 0

		for {
			select {
			case <-ticker.C:
				fmt.Fprintf(
					s.writer,
					"\r%s %s",
					s.frames[frame],
					s.message,
				)

				frame = (frame + 1) % len(s.frames)

			case <-s.stop:
				return
			}
		}
	}()
}

func (s *Spinner) Stop() {
	if !s.enabled {
		return
	}

	s.once.Do(func() {
		close(s.stop)
		<-s.done

		// Limpia la linea del spinner.
		fmt.Fprint(s.writer, "\r\033[2K")
	})
}

func isTerminal(file *os.File) bool {
	info, err := file.Stat()
	if err != nil {
		return false
	}

	return info.Mode()&os.ModeCharDevice != 0
}

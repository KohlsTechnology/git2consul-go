package exit

import (
	"os"
	"os/signal"
	"syscall"
)

func Listen(fn func(os.Signal)) {
	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, os.Interrupt, os.Kill, syscall.SIGTERM)

		sig := <-sigchan
		if fn != nil {
			fn(sig)
		}
	}()
}

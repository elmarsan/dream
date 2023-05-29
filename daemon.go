package main

import (
	"flag"
	"log"
	"os"
	"syscall"

	"github.com/sevlyar/go-daemon"
)

var (
	stopSignal = flag.String("s", "", "Send signal to stop daemon")

	stop = make(chan struct{})
	done = make(chan struct{})
)

// run starts Dream as unix daemon.
func run(process func()) error {
	flag.Parse()
	daemon.AddCommand(daemon.StringFlag(stopSignal, "stop"), syscall.SIGTERM, termHandler)

	ctx := &daemon.Context{
		PidFileName: "dream.pid",
		PidFilePerm: 0644,
		LogFileName: "dream.log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
	}

	if len(daemon.ActiveFlags()) > 0 {
		d, err := ctx.Search()
		if err != nil {
			log.Fatalf("Unable send signal to the daemon: %s", err.Error())
		}

		return daemon.SendCommands(d)
	}

	_, err := ctx.Reborn()
	if err != nil {
		return err
	}

	defer func(daemonContext *daemon.Context) error {
		err := daemonContext.Release()
		if err != nil {
			return err
		}

		return nil
	}(ctx)

	go func() {
		for {
			select {
			case <-stop:
				break
			default:
				process()
			}
		}
	}()

	err = daemon.ServeSignals()
	if err != nil {
		return err
	}

	log.Println("Dream daemon terminated")
	close(stop)
	return nil
}

func termHandler(sig os.Signal) error {
	log.Println("Terminating dream daemon...")
	stop <- struct{}{}
	if sig == syscall.SIGQUIT {
		<-done
	}

	return daemon.ErrStop
}

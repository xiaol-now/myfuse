package cmd

import (
	"github.com/spf13/cobra"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func ObserveSignal(fn func()) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGUSR1, syscall.SIGUSR2)
	go func() {
		for s := range c {
			log.Println("接收信号关闭", s.String())
			fn()
			return
		}
	}()
}

func ShowHelp(err io.Writer) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		cmd.SetErr(err)
		cmd.HelpFunc()(cmd, args)
		return nil
	}
}

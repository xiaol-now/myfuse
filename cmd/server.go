package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"myfuse/logic"
	"myfuse/utils"
	"net"
	"strconv"
)

var ServerCmd = &cobra.Command{
	Use:   "server [PATH]",
	Short: "start myfuse server",
	Args:  cobra.MaximumNArgs(1),
	//RunE: func(cmd *cobra.Command, args []string) error {
	//	return ShowHelp(os.Stderr)(cmd, args)
	//},
	Run: func(cmd *cobra.Command, args []string) {
		utils.Mkdir(args[0])
		listen, _ := net.Listen("tcp", ":"+strconv.Itoa(port))
		log.Printf("Myfuse Server %s started", listen.Addr().String())
		s := logic.NewRpcServer(args[0])
		s.RegisterNotifyService()
		s.Serve(listen)
	},
}

var port int

func init() {
	ServerCmd.Flags().IntVarP(&port, "port", "p", 10000, "listening port")
}

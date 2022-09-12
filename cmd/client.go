package cmd

import (
	"context"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"myfuse/fusefs"
	"myfuse/logic"
	"myfuse/utils"
	"path/filepath"
	"time"
)

var ClientCmd = &cobra.Command{
	Use:   "client [MOUNT PATH]",
	Short: "start myfuse client",
	Args:  cobra.MaximumNArgs(1),
	//RunE: func(cmd *cobra.Command, args []string) error {
	//	return ShowHelp(os.Stderr)(cmd, args)
	//},
	Run: func(cmd *cobra.Command, args []string) {
		pathMount := args[0]
		id := utils.HashStr(pathMount)
		pathData := filepath.Join(dataPath, id)
		utils.Mkdir(pathMount)
		utils.Mkdir(pathData)

		rpcClient, err1 := logic.NewRpcClient(context.Background(), addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err1 != nil {
			log.Fatalln(err1)
		}
		serviceClient, err2 := rpcClient.NewNotifyServiceClient(id)
		if err2 != nil {
			log.Fatalln(err2)
		}
		StartFuseMount(serviceClient, pathMount, pathData, debug)
	},
}

func StartFuseMount(cli *logic.NotifyServiceClient, pathMount, pathData string, debug bool) {
	loopbackRoot, err3 := fusefs.NewLoopbackRoot(pathData, cli)
	if err3 != nil {
		log.Fatalln(err3)
	}
	sec := time.Second
	opts := &fs.Options{
		// These options are to be compatible with libfuse defaults,
		// making benchmarking easier.
		AttrTimeout:  &sec,
		EntryTimeout: &sec,
	}
	opts.Debug = debug
	opts.AllowOther = true
	opts.MountOptions.Options = append(opts.MountOptions.Options, "default_permissions")
	// First column in "df -T": original dir
	opts.MountOptions.Options = append(opts.MountOptions.Options, "fsname="+pathData)
	// Second column in "df -T" will be shown as "fuse." + Name
	opts.MountOptions.Name = "loopback"
	// Leave file permissions on "000" files as-is
	//opts.NullPermissions = true
	server, err := fs.Mount(pathMount, loopbackRoot, opts)
	if err != nil {
		log.Fatalf("Mount fail: %v\n", err)
	}
	ObserveSignal(func() {
		_ = server.Unmount()
	})
	server.Wait()
}

var debug bool
var dataPath string
var addr string

func init() {
	ClientCmd.Flags().BoolVar(&debug, "debug", false, "debug model")
	ClientCmd.Flags().StringVar(&dataPath, "data_path", "/tmp/myfuse", "Client local dataPath path")
	ClientCmd.Flags().StringVar(&addr, "addr", "", "Server addr")
}

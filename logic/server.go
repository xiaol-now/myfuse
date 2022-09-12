package logic

import (
	"google.golang.org/grpc"
	"log"
	"myfuse/logic/proto"
	"net"
)

type RpcServer struct {
	Path   string
	server *grpc.Server
}

func NewRpcServer(path string, opt ...grpc.ServerOption) *RpcServer {
	return &RpcServer{
		Path:   path,
		server: grpc.NewServer(opt...),
	}
}

func (s *RpcServer) RegisterNotifyService() {
	proto.RegisterNotifyServiceServer(s.server, NewNotifyService(s.Path))
}

func (s *RpcServer) Serve(listener net.Listener) {
	log.Fatalln(s.server.Serve(listener))
}

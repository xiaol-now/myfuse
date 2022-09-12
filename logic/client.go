package logic

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"myfuse/logic/proto"
	"path/filepath"
)

type RpcClient struct {
	conn *grpc.ClientConn
	ctx  context.Context
}

func NewRpcClient(ctx context.Context, addr string, opt ...grpc.DialOption) (*RpcClient, error) {
	conn, err := grpc.DialContext(ctx, addr, opt...)
	if err != nil {
		return nil, err
	}
	return &RpcClient{conn: conn, ctx: ctx}, nil
}

func (c *RpcClient) Close() error {
	return c.conn.Close()
}

type NotifyServiceClient struct {
	ctx context.Context
	cli proto.NotifyServiceClient
	Id  string
}

func (c *RpcClient) NewNotifyServiceClient(id string) (*NotifyServiceClient, error) {
	cli := proto.NewNotifyServiceClient(c.conn)
	_, err := cli.BindPath(context.Background(), &proto.BindPathRequest{Id: id})
	if err != nil {
		return nil, err
	}
	return &NotifyServiceClient{ctx: c.ctx, cli: cli, Id: id}, nil
}

func (n *NotifyServiceClient) Mknod(name string, mode uint32, dev uint32) {
	resp, err := n.cli.Mknod(n.ctx, &proto.MknodRequest{
		Name: filepath.Join(n.Id, name),
		Mode: mode,
		Dev:  dev,
	})
	log.Printf("Mknod event: %t %v", resp.Success, err)
}

func (n *NotifyServiceClient) Mkdir(name string, mode uint32) {
	resp, err := n.cli.Mkdir(n.ctx, &proto.MkdirRequest{
		Name: filepath.Join(n.Id, name),
		Mode: mode,
	})
	log.Printf("Mkdir event: %t %v", resp.Success, err)
}

func (n *NotifyServiceClient) Rmdir(name string) {
	resp, err := n.cli.Rmdir(context.TODO(), &proto.RmdirRequest{Name: filepath.Join(n.Id, name)})
	log.Printf("Rmdir event: %t %v", resp.Success, err)
}

func (n *NotifyServiceClient) Rename(oldName string, newName string) {
	resp, err := n.cli.Rename(context.TODO(), &proto.RenameRequest{
		OldName: filepath.Join(n.Id, oldName),
		NewName: filepath.Join(n.Id, newName),
	})
	log.Printf("Rename event: %t %v", resp.Success, err)
}

func (n *NotifyServiceClient) Create(name string, mode uint32, flags uint32) {
	resp, err := n.cli.Create(context.TODO(), &proto.CreateRequest{
		Name:  filepath.Join(n.Id, name),
		Mode:  mode,
		Flags: flags,
	})
	log.Printf("Create event: %t %v", resp.Success, err)
}

func (n *NotifyServiceClient) Write(name string, data []byte, offset int64) {
	resp, err := n.cli.Write(context.TODO(), &proto.WriteRequest{
		Name:   filepath.Join(n.Id, name),
		Data:   data,
		Offset: offset,
	})
	log.Printf("Write event: %t %v", resp.Success, err)
}

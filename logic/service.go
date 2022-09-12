package logic

import (
	"context"
	"log"
	"myfuse/logic/proto"
	"myfuse/utils"
	"os"
	"path/filepath"
	"syscall"
)

type NotifyService struct {
	Path string
	km   *utils.KeyMutex
}

func (n *NotifyService) BindPath(ctx context.Context, req *proto.BindPathRequest) (*proto.Response, error) {
	utils.Mkdir(filepath.Join(n.Path, req.Id))
	return &proto.Response{Success: true}, nil
}

func NewNotifyService(path string) *NotifyService {
	return &NotifyService{Path: path, km: new(utils.KeyMutex)}
}

func (n *NotifyService) Rmdir(ctx context.Context, req *proto.RmdirRequest) (*proto.Response, error) {
	p := filepath.Join(n.Path, req.Name)
	if err := syscall.Rmdir(p); err != nil {
		log.Println("Rmdir", err.Error())
		return nil, err
	}
	return &proto.Response{Success: true}, nil
}

func (n *NotifyService) Rename(ctx context.Context, req *proto.RenameRequest) (*proto.Response, error) {
	p1 := filepath.Join(n.Path, req.OldName)
	p2 := filepath.Join(n.Path, req.NewName)
	if err := syscall.Rename(p1, p2); err != nil {
		log.Println("Rename", err.Error())
		return nil, err
	}
	return &proto.Response{Success: true}, nil
}

func (n *NotifyService) Mknod(ctx context.Context, req *proto.MknodRequest) (*proto.Response, error) {
	p := filepath.Join(n.Path, req.Name)
	err := syscall.Mknod(p, req.Mode, int(req.Dev))
	if err != nil {
		log.Println("Mknod", err.Error())
		return nil, err
	}
	return &proto.Response{Success: true}, nil
}

func (n *NotifyService) Mkdir(ctx context.Context, req *proto.MkdirRequest) (*proto.Response, error) {
	p := filepath.Join(n.Path, req.Name)
	err := os.Mkdir(p, os.FileMode(req.Mode))
	if err != nil {
		log.Println("Mkdir", err.Error())
		return nil, err
	}
	return &proto.Response{Success: true}, nil
}

func (n *NotifyService) Create(ctx context.Context, req *proto.CreateRequest) (*proto.Response, error) {
	p := filepath.Join(n.Path, req.Name)
	fd, err := syscall.Open(p, int(req.Flags)|os.O_CREATE, req.Mode)
	if err != nil {
		log.Println("Create", err.Error())
		return nil, err
	}
	_ = syscall.Close(fd)
	return &proto.Response{Success: true}, nil
}

func (n *NotifyService) Write(ctx context.Context, req *proto.WriteRequest) (*proto.Response, error) {
	p := filepath.Join(n.Path, req.Name)
	fd, err1 := syscall.Open(p, os.O_WRONLY, 0755)
	if err1 != nil {
		log.Println("Open file", p, err1.Error())
		return nil, err1
	}
	defer func() { _ = syscall.Close(fd) }()

	unlock := n.km.Lock(p)
	defer unlock()

	_, err2 := syscall.Pwrite(fd, req.Data, req.Offset)
	if err2 != nil {
		log.Println("Write file", p, err2.Error())
		return nil, err2
	}
	return &proto.Response{Success: true}, nil
}

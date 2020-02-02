package main

import (
	"GitCode/Go-gRPC-Demo/grpcDemo/config"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"practicProject/myTest/redisDemo/web-grpc-redis/grpc/model"
	login "practicProject/myTest/redisDemo/web-grpc-redis/grpc/proto"
)

type loginService struct{}
var LoginService = loginService{}

// 实现SayHello服务
func (l loginService) Login(ctx context.Context, in *login.LoginRequest) (*login.LoginResponse, error) {
	username := in.Username
	password := in.Password

	client := model.NewRedisClient()
	status, userName := model.UserLogin(client, username, "", password, "yan")

	resp := &login.LoginResponse{}
	resp.Msg = status
	if status == "login sucessfully" {
		resp.Code = "200"
		resp.Data = fmt.Sprintf("{username: %s}", userName)
	} else {
		resp.Code = "400"
		resp.Data = fmt.Sprintf("{}")
	}

	return resp, nil
}

func (l loginService) Register(ctx context.Context, r *login.RegisterRequest) (*login.RegisterResponse, error) {
	userInfo := make(map[string]string)
	userInfo["username"] = r.Username
	userInfo["password"] = r.Password
	userInfo["email"] = r.Email
	userInfo["phone"] = r.Phone

	client := model.NewRedisClient()
	status := model.RedisRegister(client, userInfo)

	resp := &login.RegisterResponse{}
	resp.Message = status

	if status == "register successfully" {
		resp.Status = "200"
	} else {
		resp.Status = "400"
	}
	return resp, nil
}

func main() {
	listen, err := net.Listen("tcp", config.Address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 实例化grpc Server
	srv := grpc.NewServer()

	// 注册HelloService
	//hello.RegisterHelloServer(srv, &HelloService)
	login.RegisterLoginSrvServer(srv, &LoginService)

	fmt.Println("Listen on " + config.Address)
	// 等待网络连接
	srv.Serve(listen)
}
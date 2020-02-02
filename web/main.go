package main

import (
	"GitCode/Go-gRPC-Demo/grpcDemo/config"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"net/http"
	login "practicProject/myTest/redisDemo/web-grpc-redis/grpc/proto"
)

type Login struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Register struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

func HelloHandler(ctx *gin.Context) {
	name := ctx.DefaultQuery("name", "Mike")
	address := ctx.Query("address")

	fmt.Printf("name=%s, address=%s\n", name, address)
	ctx.String(http.StatusOK, "name=%s, address=%s", name, address)
}

func HiHandler(ctx *gin.Context) {
	message := ctx.PostForm("message")
	nick := ctx.DefaultPostForm("nick", "anonymous") // 此方法可以设置默认值

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "posted",
		"message": message,
		"nick":    nick,
	})
}

func LoginHandler(ctx *gin.Context) {
	var login Login
	if err := ctx.ShouldBind(&login); err != nil {
		log.Println("ctx should bind error:", err)
	}

	login.grpcLogin()
	ctx.JSON(http.StatusOK, gin.H{
		"code": 1000,
		"msg":  "success",
		"data": login,
	})
}

func RegisterHandler(ctx *gin.Context) {
	var register Register
	if err := ctx.ShouldBind(&register); err != nil {
		log.Println("ctx should bind error:", err)
	}

	register.grpcRegister()
	ctx.JSON(http.StatusOK, gin.H{
		"code": 1000,
		"msg":  "success",
		"data": register,
	})
}

func (r Register) grpcRegister() {
	// 客户带连接服务器
	conn, err := grpc.Dial(config.Address, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	// 初始化客户端（获得client句柄）
	client := login.NewLoginSrvClient(conn)

	// 实例对象
	reqBody := &login.RegisterRequest{}
	reqBody.Username = r.UserName
	reqBody.Password = r.Password
	reqBody.Email = r.Email
	reqBody.Phone = r.Phone

	// 通过句柄调用函数
	resp, err := client.Register(context.Background(), reqBody)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(resp)
}

func (l Login) grpcLogin() {
	// 客户带连接服务器
	conn, err := grpc.Dial(config.Address, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	// 初始化客户端（获得client句柄）
	//client := hello.NewHelloClient(conn)
	client := login.NewLoginSrvClient(conn)

	// 实例对象
	//reqBody := &hello.HelloRequest{}
	reqBody := &login.LoginRequest{}
	reqBody.Username = l.UserName
	reqBody.Password = l.Password
	// 通过句柄调用函数
	//resp, err := client.SayHello(context.Background(), reqBody)
	resp, err := client.Login(context.Background(), reqBody)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(resp)
}

func main() {
	router := gin.Default()

	router.GET("/hello", HelloHandler)
	router.POST("/hi", HiHandler)

	v1 := router.Group("/v1")
	{
		v1.POST("/login", LoginHandler)
		v1.POST("/register", RegisterHandler)
	}

	if err := router.Run(":8080"); err != nil {
		log.Fatal("gin router run error:", err)
	}

}

package api

import (
	"context"
	"distributed_streaming/cmd/adapter/app/datatype/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/fx"

	"google.golang.org/grpc"

	uuid "github.com/satori/go.uuid"

	pb "distributed_streaming/cmd/persist/app/datatype/pb"
)

var Modual = fx.Options(

	fx.Provide(
		NewGrpcClient,
	),
	fx.Invoke(
		InitAPI,
		InitGrpcClient,
	),
)

var client pb.PersistServiceClient

func GetUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"eventType": "GetUsers",
	})
}

func CreateUser(c *gin.Context) {

	var req json.CreateUserRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Binding error",
		})
	}

	uid := uuid.Must(uuid.NewV4(), err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.CreateUser(
		ctx,
		&pb.CreateUserRequest{
			Id:      uid.String(),
			Name:    req.Name,
			Balance: uint64(req.Balance),
		},
	)
	if err != nil {
		log.Println(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"eventType": "CreateUser",
		"status":    r.GetReply(),
		"payload": gin.H{
			"uid": uid,
		},
	})
}

func GetUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

func UpdateUser(c *gin.Context) {

	var req json.UpdateUserRequest
	err := c.ShouldBindJSON(&req)
	uid := c.Param("id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Binding error",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"eventType": "UpdateUser",
		"payload": gin.H{
			"uid": uid,
		},
	})
}

func DeleteUser(c *gin.Context) {

	uid := c.Param("id")

	c.JSON(http.StatusOK, gin.H{
		"eventType": " DeleteUser",
		"payload": gin.H{
			"uid": uid,
		},
	})
}

func GetTxs(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

func CreateTxs(c *gin.Context) {

	var req json.CreateTxRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Binding error",
		})
	}

	uid := uuid.Must(uuid.NewV4(), err)
	fmt.Println(uid)

	c.JSON(http.StatusOK, gin.H{
		"eventType": "CreateTx",
		"payload":   gin.H{},
	})
}

func NewGrpcClient() *grpc.ClientConn {

	conn, err := grpc.Dial(
		viper.GetString("grpc.port"),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	// defer conn.Close()

	return conn
}

func InitGrpcClient(conn *grpc.ClientConn) {
	client = pb.NewPersistServiceClient(conn)
}

func InitAPI(router *gin.Engine) {

	router.GET("/users", GetUsers)
	router.POST("/users", CreateUser)

	router.GET("/user/:id", GetUser)
	router.PUT("/user/:id", UpdateUser)
	router.DELETE("/user/:id", DeleteUser)

	router.GET("/txs", GetTxs)
	router.POST("/txs", CreateTxs)

}

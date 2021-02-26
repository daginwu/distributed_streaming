package api

import (
	"context"
	"distributed_streaming/cmd/adapter/app/datatype/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	cestan "github.com/cloudevents/sdk-go/protocol/stan/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/client"

	uuid "github.com/satori/go.uuid"
)

var clientSet map[string]client.Client

var Modual = fx.Options(

	fx.Provide(),
	fx.Invoke(
		InitAPI,
	),
)

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

	EmitEvent(
		clientSet["users"],
		"adapter",
		"create",
		map[string]interface{}{
			"uid":     uid.String(),
			"name":    req.Name,
			"balance": req.Balance,
		},
	)

	c.JSON(http.StatusOK, gin.H{
		"eventType": "CreateUser",
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

	EmitEvent(
		clientSet["users"],
		"adapter",
		"update",
		map[string]interface{}{
			"uid":     uid,
			"name":    req.Name,
			"balance": req.Balance,
		},
	)

	c.JSON(http.StatusOK, gin.H{
		"eventType": "UpdateUser",
		"payload": gin.H{
			"uid": uid,
		},
	})
}

func DeleteUser(c *gin.Context) {

	uid := c.Param("id")

	EmitEvent(
		clientSet["users"],
		"adapter",
		"delete",
		map[string]interface{}{
			"uid": uid,
		},
	)

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

	EmitEvent(
		clientSet["txs"],
		"adapter",
		"create",
		map[string]interface{}{
			"uid":   uid.String(),
			"from":  req.From,
			"to":    req.To,
			"money": req.Money,
		},
	)

	c.JSON(http.StatusOK, gin.H{
		"eventType": "CreateTx",
		"payload":   gin.H{},
	})
}

func NewStreamingClient(subject string) client.Client {

	var err error
	clientID := uuid.Must(uuid.NewV4(), err)
	s, err := cestan.NewSender(
		"test-cluster",
		clientID.String(),
		subject,
		cestan.StanOptions(),
	)
	if err != nil {
		log.Println("failed to create protocol: %v", err)
	}

	c, err := cloudevents.NewClient(
		s,
		cloudevents.WithTimeNow(),
		cloudevents.WithUUIDs(),
	)
	if err != nil {
		log.Println("failed to create client: %v", err)
	}

	return c
}

func EmitEvent(client client.Client, source string, eventType string, payload map[string]interface{}) {

	e := cloudevents.NewEvent()
	e.SetType(eventType)
	e.SetSource(source)
	_ = e.SetData(cloudevents.ApplicationJSON, payload)

	if result := client.Send(context.Background(), e); cloudevents.IsUndelivered(result) {
		log.Println("failed to send")
	} else {
		log.Println(cloudevents.IsACK(result))
	}
}

func InitAPI(router *gin.Engine) {

	clientSet = map[string]client.Client{
		"users": NewStreamingClient("users"),
		"txs":   NewStreamingClient("txs"),
	}

	log.Println(clientSet)

	router.GET("/users", GetUsers)
	router.POST("/users", CreateUser)

	router.GET("/user/:id", GetUser)
	router.PUT("/user/:id", UpdateUser)
	router.DELETE("/user/:id", DeleteUser)

	router.GET("/txs", GetTxs)
	router.POST("/txs", CreateTxs)

}

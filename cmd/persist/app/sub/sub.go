package sub

import (
	"fmt"

	badger "github.com/dgraph-io/badger/v3"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/fx"

	"context"
	"log"

	json "distributed_streaming/cmd/persist/app/datatype"

	cestan "github.com/cloudevents/sdk-go/protocol/stan/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/client"
)

var database *badger.DB
var clientSet *client.Client

var Modual = fx.Options(

	fx.Provide(
		NewSub,
	),
	fx.Invoke(
		InitSub,
	),
)

func NewSubClient(subject string) client.Client {

	var err error
	clientID := uuid.Must(uuid.NewV4(), err)
	r, err := cestan.NewConsumer(
		"test-cluster",
		clientID.String(),
		subject,
		cestan.StanOptions(),
	)
	if err != nil {
		log.Println("failed to create protocol: %v", err)
	}

	c, err := cloudevents.NewClient(
		r,
		cloudevents.WithTimeNow(),
		cloudevents.WithUUIDs(),
	)
	if err != nil {
		log.Println("failed to create client: %v", err)
	}

	return c
}

func PersistData(_ context.Context, event cloudevents.Event) {

	fmt.Printf("%s", event)
	meta := &json.EventMeta{}
	if err := event.DataAs(meta); err != nil {
		log.Println(err)
		log.Println(event)
	}
	switch meta.EventType {
	case "create":
		log.Println("create")
	case "update":
		log.Println("update")
	case "delete":
		log.Println("delete")
	default:

	}
}

func NewSub(lc fx.Lifecycle) map[string]client.Client {
	clientSet := make(map[string]client.Client)

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			for _, v := range clientSet {
				go func() {
					err := v.StartReceiver(context.TODO(), PersistData)
					if err != nil {
						log.Println("failed to start receiver: ", err)
					} else {
						log.Println("receiver stopped\n")
					}
				}()

			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
	return clientSet
}

func InitSub(db *badger.DB, clientSet map[string]client.Client) {

	database = db

	clientSet["users"] = NewSubClient("users")
	clientSet["txs"] = NewSubClient("txs")

	log.Println(clientSet)

}

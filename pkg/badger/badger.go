package keyvalue

import (
	"log"

	badger "github.com/dgraph-io/badger/v3"
	"go.uber.org/fx"
)

var Modual = fx.Options(

	fx.Provide(
		NewKeyValue,
	),
	fx.Invoke(
		InitNewKeyValue,
	),
)

func NewKeyValue(lc fx.Lifecycle) *badger.DB {

	log.Println("[Distributed_Streaming] Badger Key-Value database start")

	db, err := badger.Open(badger.DefaultOptions("datastore"))
	if err != nil {
		log.Println(err)
	}

	return db
}

func InitNewKeyValue(db *badger.DB) error {
	return nil
}

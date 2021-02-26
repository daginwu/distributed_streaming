package streaming

import (
	"go.uber.org/fx"
)

var Modual = fx.Options(

	fx.Provide(NewStreamingClient),
	fx.Invoke(InitStreamingClient),
)

func NewStreamingClient(lc fx.Lifecycle) {

}

func InitStreamingClient() {

}

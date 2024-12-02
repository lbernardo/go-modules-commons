package rabbitmq

import "go.uber.org/fx"

var Module = fx.Module("gmc.rabbitmq", fx.Provide(NewRabbitMQ))

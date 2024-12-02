package rabbitmq

type Message struct {
	Body        []byte
	ConsumerTag string
}

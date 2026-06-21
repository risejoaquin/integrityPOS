package events

type Broker struct {
	Broadcast      chan []byte
	NewClients     chan chan []byte
	ClosingClients chan chan []byte
	clients        map[chan []byte]bool
}

func NewBroker() *Broker {
	return &Broker{
		Broadcast:      make(chan []byte, 100),
		NewClients:     make(chan chan []byte),
		ClosingClients: make(chan chan []byte),
		clients:        make(map[chan []byte]bool),
	}
}

func (broker *Broker) Start() {
	for {
		select {
		case s := <-broker.NewClients:
			broker.clients[s] = true
		case s := <-broker.ClosingClients:
			delete(broker.clients, s)
			close(s)
		case event := <-broker.Broadcast:
			for clientMessageChan := range broker.clients {
				select {
				case clientMessageChan <- event:
				default:
					// Detach client safely if buffer is filled
					delete(broker.clients, clientMessageChan)
					close(clientMessageChan)
				}
			}
		}
	}
}

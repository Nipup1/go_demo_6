package event

const(
	EventLinkVisited = "link.visited"
)

//структура события
type Event struct{
	Type string
	Data any
}

//наш канал для событий
type EventBus struct{
	bus chan Event
}

func NewEventBus() *EventBus{
	return &EventBus{
		bus: make(chan Event),
	}
}

//метод добавления события
func (e *EventBus) Publish(event Event){
	e.bus <- event
}

//метод подписки на события
func (e *EventBus) Subscribe() <- chan Event{
	return e.bus
}
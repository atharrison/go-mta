package main

func handleEnvelope(env envelope) {
	Info.Println("Dispatcher dispatching message...")
	Debug.Println("Received Envelope:", env)

	go send(env)
}


package client

//import "testing"

func ExampleClient() {
	client := NewClient(":6666", "cli1")
	defer client.Close()

	client.Go()
}

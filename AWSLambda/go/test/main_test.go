package test

import (
	"context"
	"log"
	"main/handlers"
	"testing"
)

func TestHelloWorld(t *testing.T) {
	expected_result := "Hello from Lambda!"

	ctx := context.TODO()
	defer ctx.Done()
	result, err := handlers.SimpleHandler(ctx)

	if err != nil {
		t.Error("The hello world funtion breaks")
	}

	if result != expected_result {
		t.Error("Result is not the same as expected")
	}

	log.Print("Test works Correctly")
}

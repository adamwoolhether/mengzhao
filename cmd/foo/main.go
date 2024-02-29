package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/replicate/replicate-go"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// You can also provide a token directly with
	// `replicate.NewClient(replicate.WithToken("r8_..."))`
	r8, err := replicate.NewClient(replicate.WithTokenFromEnv())
	if err != nil {
		log.Fatal(err)
	}

	//https://replicate.com/stability-ai/stable-diffusion
	//version := "stability-ai/stable-diffusion:ac732df83cea7fff18b8472768c88ad041fa750ff7682a21affe81863cbe77e4"
	//
	input := replicate.PredictionInput{
		"prompt": "an astronaut riding a horse on mars, hd, dramatic lighting",
	}

	webhook := replicate.Webhook{
		URL:    os.Getenv("WEBHOOK_URL"),
		Events: []replicate.WebhookEventType{"start", "completed"},
	}

	// Run a model and wait for its output
	//output, err := r8.Run(ctx, version, input, &webhook)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println("output: ", output)

	//// Create a prediction
	version := "ac732df83cea7fff18b8472768c88ad041fa750ff7682a21affe81863cbe77e4"
	prediction, err := r8.CreatePrediction(ctx, version, input, &webhook, false)
	if err != nil {
		log.Fatal(err)
	}
	//
	// Wait for the prediction to finish
	err = r8.Wait(ctx, prediction)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("output: ", prediction)
}

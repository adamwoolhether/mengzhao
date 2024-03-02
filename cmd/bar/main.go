package main

/*
import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudflare/cloudflare-go"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	client, err := cloudflare.New(os.Getenv("CLOUDFLARE_API_TOKEN"), os.Getenv("CLOUDFLARE_EMAIL"))
	if err != nil {
		log.Fatal(err)
	}

	u, err := client.UserDetails(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", u)

	imgParams := cloudflare.UploadImageParams{
		//File:              nil,
		URL:  "https://cdn.pixabay.com/photo/2024/02/27/00/13/heliconia-8599119_1280.jpg",
		Name: "",
		//RequireSignedURLs: false,
		Metadata: map[string]interface{}{
			"key": "value",
		},
	}

	rc := cloudflare.ResourceContainer{
		Level:      cloudflare.AccountRouteLevel,
		Identifier: os.Getenv("CLOUDFLARE_ACCOUNT_ID"),
		Type:       cloudflare.AccountType,
	}

	img, err := client.UploadImage(ctx, &rc, imgParams)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", img)
}
*/

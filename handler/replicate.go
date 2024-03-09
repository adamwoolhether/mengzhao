package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/cloudflare/cloudflare-go"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"golang.org/x/sync/errgroup"

	"mengzhao/db"
	"mengzhao/types"
)

const (
	successWebhookResponse    = "succeeded"
	processingWebhookResponse = "processing"
)

type ReplicateResp struct {
	Input struct {
		Prompt string `json:"prompt"`
	} `json:"input"`
	Status string   `json:"status"`
	Output []string `json:"output"`
}

func ReplicateCallback(w http.ResponseWriter, r *http.Request) error {
	var resp ReplicateResp

	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return fmt.Errorf("decoding body: %w", err)
	}

	if resp.Status == processingWebhookResponse {
		return nil
	}

	if resp.Status != successWebhookResponse {
		return fmt.Errorf("replicate callback failed: %s", resp.Status)
	}

	batchID, err := uuid.Parse(chi.URLParam(r, "batch_id"))
	if err != nil {
		return fmt.Errorf("replicate callback batch_id invalid: %w", err)
	}

	images, err := db.GetImagesByBatchID(r.Context(), batchID)
	if err != nil {
		return fmt.Errorf("replicate callback getting batch[%s] images: %w", batchID, err)
	}

	if len(images) != len(resp.Output) {
		return fmt.Errorf("replicate callback batch[%s] images and output len: %d!= %d", batchID, len(images), len(resp.Output))
	}

	// CLOUDFLARE UPDLOAD AND GET URL HERE.
	cf, err := cloudflare.New(os.Getenv("CLOUDFLARE_API_TOKEN"), os.Getenv("CLOUDFLARE_EMAIL"))
	if err != nil {
		return fmt.Errorf("replicate callback cloudflare client: %w", err)
	}

	txFunc := func(ctx context.Context, tx bun.Tx) error {
		var eg errgroup.Group

		eg.Go(func() error {
			for i, imageURL := range resp.Output {
				newURL, err := uploadToCloudFlare(r.Context(), cf, imageURL, images[i])
				if err != nil {
					slog.Error("upload to cloudflare", "imageURL", imageURL, "batchID", batchID, "err", err)
					images[i].Status = types.ImageStatusFailed
				} else {
					images[i].Status = types.ImageStatusCompleted
					images[i].ImgLoc = newURL
				}

				if err := db.UpdateImage(r.Context(), tx, &images[i]); err != nil {
					return err
				}
			}

			return nil
		})

		if err := eg.Wait(); err != nil {
			return fmt.Errorf("batch update images: %w", err)
		}

		return nil
	}

	if err := db.Bun.RunInTx(r.Context(), &sql.TxOptions{}, txFunc); err != nil {
		return err
	}

	return nil
}

func uploadToCloudFlare(ctx context.Context, client *cloudflare.API, replicateURL string, image types.Image) (string, error) {
	imgParams := cloudflare.UploadImageParams{
		//File:              nil,
		URL:  replicateURL,
		Name: "",
		//RequireSignedURLs: false,
		Metadata: map[string]interface{}{
			"userID":  image.UserID,
			"batchID": image.BatchID,
			"srcURL":  image.ImgLoc,
		},
	}

	rc := cloudflare.ResourceContainer{
		Level:      cloudflare.AccountRouteLevel,
		Identifier: os.Getenv("CLOUDFLARE_ACCOUNT_ID"),
		Type:       cloudflare.AccountType,
	}

	img, err := client.UploadImage(ctx, &rc, imgParams)
	if err != nil {
		return "", fmt.Errorf("uploading image[%d] batch[%s] from srcURL[%s]: %w", image.ID, image.BatchID, image.ImgLoc, err)
	}

	return img.Variants[0], nil
}

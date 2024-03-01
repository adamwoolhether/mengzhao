package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/uptrace/bun"

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

	txFunc := func(ctx context.Context, tx bun.Tx) error {
		for i, imageURL := range resp.Output {
			images[i].Status = types.ImageStatusCompleted
			images[i].ImgLoc = imageURL
			images[i].Prompt = resp.Input.Prompt

			if err := db.UpdateImage(r.Context(), &images[i]); err != nil {
				return err
			}
		}

		return nil
	}

	if err := db.Bun.RunInTx(r.Context(), &sql.TxOptions{}, txFunc); err != nil {
		return err
	}

	return nil
}

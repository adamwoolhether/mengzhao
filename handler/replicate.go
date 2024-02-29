package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"mengzhao/db"
	"mengzhao/types"
)

const successWebhookResponse = "succeeded"

type ReplicateResp struct {
	Status string   `json:"status"`
	Output []string `json:"output"`
}

func ReplicateCallback(w http.ResponseWriter, r *http.Request) error {
	var resp ReplicateResp

	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return fmt.Errorf("decoding body: %w", err)
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

	for i, imageURL := range resp.Output {
		images[i].Status = types.ImageStatusCompleted
		images[i].ImgLoc = imageURL

	}

	return nil
}

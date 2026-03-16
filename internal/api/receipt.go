package api

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/christianselig/apollo-backend/internal/domain"
	"github.com/christianselig/apollo-backend/internal/itunes"
)

func (a *api) checkReceiptHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	vars := mux.Vars(r)
	apns := vars["apns"]

	// Discard the receipt body — no validation is performed on a self-hosted server
	// because sideloaded apps cannot use App Store in-app purchases.
	_, _ = ioutil.ReadAll(r.Body)

	// Always treat the device as having a valid lifetime subscription.
	if apns != "" {
		dev, err := a.deviceRepo.GetByAPNSToken(ctx, apns)
		if err == nil {
			dev.ExpiresAt = time.Now().Add(domain.DeviceActiveAfterReceitCheckDuration)
			dev.GracePeriodExpiresAt = dev.ExpiresAt.Add(domain.DeviceGracePeriodAfterReceiptExpiry)
			_ = a.deviceRepo.Update(ctx, &dev)
		}
	}

	info := itunes.ClientVerificationInfo{
		Products: []itunes.VerificationProduct{
			{Name: "ultra", Status: "LIFETIME", SubscriptionType: "LIFETIME"},
			{Name: "pro", Status: "LIFETIME"},
			{Name: "community_icons", Status: "LIFETIME"},
			{Name: "spca", Status: "LIFETIME"},
		},
	}

	w.WriteHeader(http.StatusOK)
	bb, _ := json.Marshal(info)
	_, _ = w.Write(bb)
}

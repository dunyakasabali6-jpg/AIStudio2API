package app

import (
	"context"
	"errors"
	"testing"
)

func TestNorthflankUpdateServiceDoesNotPollOrApplyBinaryUpdates(t *testing.T) {
	service := newSelfUpdater("v1.2.3", "", nil, nil, nil)

	status, err := service.UpdateStatus(context.Background(), true)
	if err != nil {
		t.Fatalf("UpdateStatus() error = %v", err)
	}
	if status.Mode != "container" || status.Current != "v1.2.3" || status.Available {
		t.Fatalf("UpdateStatus() = %#v, want current container mode without release availability", status)
	}

	_, err = service.ApplyUpdate(context.Background())
	if err == nil {
		t.Fatal("ApplyUpdate() error = nil, want container update conflict")
	}
	var statusError interface{ HTTPStatus() int }
	if !errors.As(err, &statusError) || statusError.HTTPStatus() != 409 {
		t.Fatalf("ApplyUpdate() error = %v, want HTTP 409", err)
	}
}

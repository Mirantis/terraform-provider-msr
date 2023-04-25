package connect_test

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	client "github.com/Mirantis/terraform-provider-msr/mirantis/msr/client"
	connect "github.com/Mirantis/terraform-provider-msr/mirantis/msr/connect"
)

func TestCreateValidResourceIDs(t *testing.T) {
	ctx := context.Background()

	orgID := "mke"
	teamID := "test"
	id := connect.CreateResourceID(ctx, orgID, teamID)

	expected := fmt.Sprintf("%s%s%s", orgID, connect.IdDelimiter, teamID)

	if !reflect.DeepEqual(id, expected) {
		t.Errorf("expected (%v), got (%v)", expected, id)
	}
}

func TestCreateInvalidResourceIDs(t *testing.T) {
	ctx := context.Background()

	orgID := "mke"
	teamID := "test"
	id := connect.CreateResourceID(ctx, orgID, teamID)

	expected := fmt.Sprintf("%s%s%swrong", orgID, connect.IdDelimiter, teamID)

	if reflect.DeepEqual(id, expected) {
		t.Errorf("expected id: (%v), got (%v)", expected, id)
	}
}

func TestParseValidResourceIDs(t *testing.T) {
	ctx := context.Background()

	eOrgID := "mke"
	eResID := "test"
	id := fmt.Sprintf("%s%s%s", eOrgID, connect.IdDelimiter, eResID)

	orgID, resourceID, err := connect.ExtractResourceIDs(ctx, id)
	if err != nil {
		t.Errorf("resource ID is invalid format '%s'", id)
	}
	if !reflect.DeepEqual(orgID, eOrgID) {
		t.Errorf("expected (%v), got (%v)", eOrgID, orgID)
	}

	if !reflect.DeepEqual(resourceID, eResID) {
		t.Errorf("expected (%v), got (%v)", eResID, resourceID)
	}
}

func TestParseInvalidResourceIDs(t *testing.T) {
	ctx := context.Background()

	orgID := "mke"
	teamID := "test"
	id := fmt.Sprintf("%s.%s", orgID, teamID)
	expectedErr := client.ErrInvalidResourceIDFormat
	orgID, resourceID, err := connect.ExtractResourceIDs(ctx, id)
	if !errors.Is(err, expectedErr) {
		t.Errorf("expected error: (%v),\n got (%v)", expectedErr, err)
	}

	fmt.Printf("orgID: %s\n", orgID)
	fmt.Printf("resourceID: %s\n", resourceID)
	if reflect.DeepEqual(orgID, nil) {
		t.Errorf("expected id: (%v), got (%v)", "", orgID)
	}
	if reflect.DeepEqual(resourceID, nil) {
		t.Errorf("expected id: (%v), got (%v)", "", resourceID)
	}
}

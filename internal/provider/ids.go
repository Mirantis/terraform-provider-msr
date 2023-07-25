package provider

// import (
// 	"context"
// 	"fmt"
// 	"strings"

// 	"github.com/Mirantis/terraform-provider-msr/internal/client"
// )

// const IdDelimiter = ":"

// func ExtractResourceIDs(ctx context.Context, id string) (orgID string, resourceID string, err error) {
// 	ids := strings.SplitN(id, IdDelimiter, 2)
// 	if len(ids) != 2 {
// 		return "", "", fmt.Errorf("%w '%s'", client.ErrInvalidResourceIDFormat, id)
// 	}
// 	return ids[0], ids[1], nil
// }

// func CreateResourceID(ctx context.Context, orgID string, teamID string) (id string) {
// 	return fmt.Sprintf("%s%s%s", orgID, IdDelimiter, teamID)
// }

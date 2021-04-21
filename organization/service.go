package organization

import (
	"context"
)

type Service interface {
	CreateOrganization(ctx context.Context, name string, phone string, website string) (string, error)
	GetOrganization(ctx context.Context, id string) (Organization, error)
}

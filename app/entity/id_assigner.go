package entity

import "context"

type IDAssigner interface {
	ID(ctx context.Context, packageName string, suuid string) (uint, error)
}

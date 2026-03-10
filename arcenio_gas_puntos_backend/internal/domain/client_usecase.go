package domain

import "context"

// ClientUsecase define las operaciones de negocio para clientes
type ClientUsecase interface {
	Create(ctx context.Context, client *Client) (*Client, error)
	GetByCedula(ctx context.Context, cedula string) (*Client, error)
	List(ctx context.Context) ([]*Client, error)
	Update(ctx context.Context, id string, client *Client) (*Client, error)
}

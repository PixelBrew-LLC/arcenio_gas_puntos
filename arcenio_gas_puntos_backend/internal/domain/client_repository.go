package domain

import "context"

// ClientRepository define las operaciones de persistencia de clientes
type ClientRepository interface {
	Create(ctx context.Context, client *Client) error
	GetByID(ctx context.Context, id string) (*Client, error)
	GetByCedula(ctx context.Context, cedula string) (*Client, error)
	List(ctx context.Context) ([]*Client, error)
	Update(ctx context.Context, client *Client) error
}

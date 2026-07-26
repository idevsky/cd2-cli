package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	pb "github.com/clouddrive/cd2-cli/pkg/proto"
)

type Client struct {
	conn          *grpc.ClientConn
	client        pb.CloudDriveFileSrvClient
	token         string
	address       string
	timeout       time.Duration
	skipVerifyTLS bool
	mu            sync.RWMutex
}

type Config struct {
	Address       string
	Token         string
	Timeout       time.Duration
	UseTLS        bool
	SkipVerifyTLS bool
}

func NewClient(cfg Config) (*Client, error) {
	if cfg.Address == "" {
		return nil, fmt.Errorf("address is required")
	}

	var creds credentials.TransportCredentials
	if cfg.UseTLS {
		creds = credentials.NewTLS(&tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: cfg.SkipVerifyTLS,
		})
	} else {
		creds = insecure.NewCredentials()
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
	}

	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	conn, err := grpc.NewClient(cfg.Address, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	return &Client{
		conn:          conn,
		client:        pb.NewCloudDriveFileSrvClient(conn),
		token:         cfg.Token,
		address:       cfg.Address,
		timeout:       timeout,
		skipVerifyTLS: cfg.SkipVerifyTLS,
	}, nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *Client) SetToken(token string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.token = token
}

func (c *Client) GetToken() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.token
}

func (c *Client) SetConn(conn *grpc.ClientConn) {
	c.conn = conn
	c.client = pb.NewCloudDriveFileSrvClient(conn)
}

func (c *Client) withAuth(ctx context.Context) context.Context {
	token := c.GetToken()
	if token == "" {
		return ctx
	}
	return metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token)
}

func (c *Client) withTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout == 0 {
		timeout = c.timeout
	}
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	return context.WithTimeout(ctx, timeout)
}

func (c *Client) Public() *PublicAPI {
	return &PublicAPI{c: c}
}

func (c *Client) Auth() *AuthAPI {
	return &AuthAPI{c: c}
}

func (c *Client) File() *FileAPI {
	return &FileAPI{c: c}
}

func (c *Client) Mount() *MountAPI {
	return &MountAPI{c: c}
}

func (c *Client) Transfer() *TransferAPI {
	return &TransferAPI{c: c}
}

func (c *Client) CloudAPI() *CloudAPIAPI {
	return &CloudAPIAPI{c: c}
}

func (c *Client) Backup() *BackupAPI {
	return &BackupAPI{c: c}
}

func (c *Client) WebDAV() *WebDAVAPI {
	return &WebDAVAPI{c: c}
}

func (c *Client) Token() *TokenAPI {
	return &TokenAPI{c: c}
}

func (c *Client) Session() *SessionAPI {
	return &SessionAPI{c: c}
}

func (c *Client) System() *SystemAPI {
	return &SystemAPI{c: c}
}

func (c *Client) Offline() *OfflineAPI {
	return &OfflineAPI{c: c}
}

func (c *Client) Webhook() *WebhookAPI {
	return &WebhookAPI{c: c}
}

func (c *Client) Local() *LocalAPI {
	return &LocalAPI{c: c}
}

func (c *Client) RemoteUpload() *RemoteUploadAPI {
	return &RemoteUploadAPI{c: c}
}

func (c *Client) Copy() *CopyAPI {
	return &CopyAPI{c: c}
}

func (c *Client) Stream() *StreamAPI {
	return &StreamAPI{c: c}
}

func (c *Client) Sync() *SyncAPI {
	return &SyncAPI{c: c}
}

func (c *Client) Promotion() *PromotionAPI {
	return &PromotionAPI{c: c}
}

func collectStream[T any](stream grpc.ClientStream) ([]T, error) {
	var items []T
	for {
		var item T
		err := stream.RecvMsg(&item)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

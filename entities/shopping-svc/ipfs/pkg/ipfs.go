// Package pkg implements a stand-in for an IPFS content client. It is
// deliberately tiny: this PoC only needs ipfs to be an independent Go
// module with a real dependency on shared-lib.
package pkg

import "github.com/abhi4u1947/go-multimodule-poc/entities/shared-lib/logger"

// Client is a minimal placeholder for an IPFS content-addressing client.
type Client struct {
	log *logger.Logger
}

// NewClient returns a Client that logs via shared-lib's logger.
func NewClient() *Client {
	return &Client{log: logger.New("ipfs")}
}

// Cat "fetches" the content behind a CID (stand-in implementation).
func (c *Client) Cat(cid string) string {
	c.log.Info("fetching CID %s", cid)
	return "content-for-" + cid
}

// Pin "pins" a CID so it is retained locally (stand-in implementation).
//
// Added in ipfs v0.18.2.
func (c *Client) Pin(cid string) {
	c.log.Info("pinning CID %s", cid)
}

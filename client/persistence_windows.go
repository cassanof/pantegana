//go:build !linux || !darwin
// +build !linux !darwin

package client

// For some reason Persistence is broken on windows.
// Meanwhile I find something else, we got this going...

func (c *Client) SetupPersistence() error {
	c.Persistence = nil
	return nil
}

func (c *Client) UnPersist() {
	return
}

func (c *Client) Persist() {
	return
}

package client

// TODO: Generate this client from stellaryard-core's /api/openapi.yaml
// using openapi-generator or oapi-codegen.
//
// DO NOT hand-write HTTP calls. Always go through the generated client.
// This is the mechanism that prevents CLI and dashboard from silently
// drifting apart in what they assume core's API looks like.
//
// If the generated client doesn't support what you need, that's a signal
// core's spec is missing something — flag it, don't route around the client.

// Client is a typed HTTP client for stellaryard-core's API.
type Client struct {
	BaseURL string
}

// New creates a new API client.
func New(baseURL string) *Client {
	return &Client{BaseURL: baseURL}
}

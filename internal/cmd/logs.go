package cmd

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

// validLogContainers contains allowed container names for logs.
var validLogContainers = map[string]bool{
	"horizon":     true,
	"soroban-rpc": true,
}

// toWebSocketURL translates an HTTP/HTTPS URL into a WS/WSS URL while preserving path.
func toWebSocketURL(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	switch u.Scheme {
	case "http":
		u.Scheme = "ws"
	case "https":
		u.Scheme = "wss"
	case "ws", "wss":
		// Already WebSocket scheme
	default:
		return "", fmt.Errorf("unsupported URL scheme %q", u.Scheme)
	}
	return u.String(), nil
}

var logsCmd = &cobra.Command{
	Use:   "logs <container-name>",
	Short: "Stream container logs",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		if !validLogContainers[name] {
			fatalf(ExitArgError, "invalid container name %q (supported: horizon, soroban-rpc)", name)
		}

		follow, _ := cmd.Flags().GetBool("follow")
		tail, _ := cmd.Flags().GetString("tail")

		wsBase, err := toWebSocketURL(coreURL)
		if err != nil {
			fatalf(ExitArgError, "invalid core-url: %v", err)
		}

		wsURL := fmt.Sprintf("%s/api/v1/containers/%s/logs?follow=%v", strings.TrimRight(wsBase, "/"), name, follow)
		if tail != "" {
			wsURL = fmt.Sprintf("%s&tail=%s", wsURL, url.QueryEscape(tail))
		}

		headers := make(http.Header)
		if apiKey != "" {
			headers.Set("Authorization", "Bearer "+apiKey)
		}

		dialer := websocket.DefaultDialer
		conn, resp, err := dialer.Dial(wsURL, headers)
		if err != nil {
			if resp != nil {
				if resp.StatusCode == http.StatusUnauthorized {
					fatalf(ExitAppError, "authentication required or invalid API key (HTTP 401)")
				}
				if resp.StatusCode == http.StatusBadRequest {
					fatalf(ExitArgError, "bad request (HTTP 400)")
				}
				fatalf(ExitAppError, "failed to connect: HTTP %d", resp.StatusCode)
			}
			fatalf(ExitCoreUnreachable, "cannot connect to core at %s: %v", coreURL, err)
		}
		defer conn.Close()

		// Handle interrupt / termination signals cleanly
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

		done := make(chan struct{})

		// Goroutine to read streaming log messages
		go func() {
			defer close(done)
			for {
				msgType, message, err := conn.ReadMessage()
				if err != nil {
					// Connection closed normally or error
					return
				}
				if msgType == websocket.TextMessage {
					fmt.Println(string(message))
				}
			}
		}()

		select {
		case <-done:
			// Normal EOF / finish
			return
		case <-sigCh:
			// Gracefully send close frame to server
			_ = conn.WriteControl(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, "client interrupted"),
				time.Now().Add(time.Second),
			)
			time.Sleep(50 * time.Millisecond)
			os.Exit(ExitSuccess)
		}
	},
}

func init() {
	logsCmd.Flags().Bool("follow", false, "follow log output (like tail -f)")
	logsCmd.Flags().String("tail", "100", "number of lines to show from end of logs")

	rootCmd.AddCommand(logsCmd)
}

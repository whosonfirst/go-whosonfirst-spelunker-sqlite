package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/whosonfirst/go-whosonfirst-spelunker-httpd/www"
)

func descendantsHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupWWWOnce.Do(setupWWW)

	if setupWWWError != nil {
		slog.Error("Failed to set up common configuration", "error", setupWWWError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupWWWError)
	}

	opts := &www.DescendantsHandlerOptions{
		Spelunker:     sp,
		Authenticator: authenticator,
		Templates:     html_templates,
		URIs:          uris_table,
	}

	return www.DescendantsHandler(opts)
}

func idHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupWWWOnce.Do(setupWWW)

	if setupWWWError != nil {
		slog.Error("Failed to set up common configuration", "error", setupWWWError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupWWWError)
	}

	opts := &www.IdHandlerOptions{
		Spelunker:     sp,
		Authenticator: authenticator,
		Templates:     html_templates,
		URIs:          uris_table,
	}

	return www.IdHandler(opts)
}

func searchHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupWWWOnce.Do(setupWWW)

	if setupWWWError != nil {
		slog.Error("Failed to set up common configuration", "error", setupWWWError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupWWWError)
	}

	opts := &www.SearchHandlerOptions{
		Spelunker:     sp,
		Authenticator: authenticator,
		Templates:     html_templates,
		URIs:          uris_table,
	}

	return www.SearchHandler(opts)
}

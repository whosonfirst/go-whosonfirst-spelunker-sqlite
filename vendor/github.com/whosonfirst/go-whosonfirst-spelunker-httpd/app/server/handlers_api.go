package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/whosonfirst/go-whosonfirst-spelunker-httpd/api"
)

func geoJSONHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupCommonOnce.Do(setupCommon)

	if setupCommonError != nil {
		slog.Error("Failed to set up common configuration", "error", setupCommonError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupCommonError)
	}

	opts := &api.GeoJSONHandlerOptions{
		Spelunker: sp,
	}

	return api.GeoJSONHandler(opts)
}

func svgHandlerFunc(ctx context.Context) (http.Handler, error) {

	setupCommonOnce.Do(setupCommon)

	if setupCommonError != nil {
		slog.Error("Failed to set up common configuration", "error", setupCommonError)
		return nil, fmt.Errorf("Failed to set up common configuration, %w", setupCommonError)
	}

	sz := api.DefaultSVGSizes()

	opts := &api.SVGHandlerOptions{
		Spelunker: sp,
		Sizes:     sz,
	}

	return api.SVGHandler(opts)
}

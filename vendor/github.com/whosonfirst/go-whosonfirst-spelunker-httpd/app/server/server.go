package server

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/aaronland/go-http-server"
	"github.com/aaronland/go-http-server/handler"
	"github.com/whosonfirst/go-whosonfirst-spelunker-httpd"
)

func Run(ctx context.Context, logger *slog.Logger) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs, logger)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet, logger *slog.Logger) error {

	opts, err := RunOptionsFromFlagSet(ctx, fs)

	if err != nil {
		return fmt.Errorf("Failed to derive run options from flagset, %w", err)
	}

	return RunWithOptions(ctx, opts, logger)
}

func RunWithOptions(ctx context.Context, opts *RunOptions, logger *slog.Logger) error {

	slog.SetDefault(logger)

	// First create a local copy of RunOptions that can't be
	// modified after the fact. 'run_options' is defined in vars.go

	v, err := opts.Clone()

	if err != nil {
		return fmt.Errorf("Failed to create local run options, %w", err)
	}

	run_options = v

	// To do: Move this in to RunOptionsFromFlagSet

	uris_table = &httpd.URIs{
		// WWW/human-readable
		Descendants: "/id/{id}/descendants",		
		Id: "/id/",
		// Descendants: "/descendants/", // FIX ME: Update to use improved syntax in Go 1.22
		Search:      "/search/",

		// API/machine-readable
		GeoJSON: "/geojson/",
		SVG:     "/svg/",
	}

	// To do: Add/consult "is enabled" flags

	handlers := map[string]handler.RouteHandlerFunc{

		// WWW/human-readable
		uris_table.Descendants: descendantsHandlerFunc,
		uris_table.Id:          idHandlerFunc,
		uris_table.Search:      searchHandlerFunc,

		// API/machine-readable
		uris_table.GeoJSON: geoJSONHandlerFunc,
		uris_table.SVG:     svgHandlerFunc,
	}

	go func() {
		for uri, h := range handlers {
			slog.Info("Enable handler", "uri", uri, "handler", fmt.Sprintf("%T", h))
		}
	}()

	route_handler, err := handler.RouteHandler(handlers)

	if err != nil {
		return fmt.Errorf("Failed to create route handlers, %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", route_handler)

	s, err := server.NewServer(ctx, run_options.ServerURI)

	if err != nil {
		return fmt.Errorf("Failed to create new server, %w", err)
	}

	slog.Info("Listening for requests", "address", s.Address())

	err = s.ListenAndServe(ctx, mux)

	if err != nil {
		return fmt.Errorf("Failed to start server, %w", err)
	}

	return nil
}

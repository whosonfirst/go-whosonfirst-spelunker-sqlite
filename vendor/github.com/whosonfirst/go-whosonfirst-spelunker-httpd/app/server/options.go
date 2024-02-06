package server

import (
	"context"
	"flag"
	"fmt"

	"github.com/mitchellh/copystructure"
	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	ServerURI        string `json:"server_uri"`
	SpelunkerURI     string `json:"spelunker_uri"`
	AuthenticatorURI string `json:"authenticator_uri"`
}

func (o *RunOptions) Clone() (*RunOptions, error) {

	v, err := copystructure.Copy(o)

	if err != nil {
		return nil, fmt.Errorf("Failed to create local run options, %w", err)
	}

	return v.(*RunOptions), nil
}

func RunOptionsFromFlagSet(ctx context.Context, fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	err := flagset.SetFlagsFromEnvVars(fs, "SPELUNKER")

	if err != nil {
		return nil, fmt.Errorf("Failed to assign flags from environment variables, %w", err)
	}

	opts := &RunOptions{
		ServerURI:        server_uri,
		AuthenticatorURI: authenticator_uri,
		SpelunkerURI:     spelunker_uri,
	}

	return opts, nil
}

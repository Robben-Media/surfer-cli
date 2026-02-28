package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/builtbyrobben/surfer-cli/internal/outfmt"
)

type AuditsCmd struct {
	List   AuditsListCmd   `cmd:"" help:"List all audits"`
	Create AuditsCreateCmd `cmd:"" help:"Create a new audit for a URL"`
}

type AuditsListCmd struct{}

func (cmd *AuditsListCmd) Run(ctx context.Context) error {
	return fmt.Errorf("listing all audits is not supported by the Surfer API; use 'surfer-cli audits create <url>' to create a new audit")
}

type AuditsCreateCmd struct {
	URL string `arg:"" required:"" help:"URL to audit"`
}

func (cmd *AuditsCreateCmd) Run(ctx context.Context) error {
	client, err := getSurferClient()
	if err != nil {
		return err
	}

	result, err := client.Audits().Create(ctx, cmd.URL)
	if err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, result)
	}

	if outfmt.IsPlain(ctx) {
		headers := []string{"ID", "URL", "STATE"}
		rows := [][]string{{fmt.Sprintf("%d", result.ID), result.URL, result.State}}

		return outfmt.WritePlain(os.Stdout, headers, rows)
	}

	fmt.Fprintf(os.Stderr, "Created audit\n\n")
	fmt.Printf("ID: %d\n", result.ID)
	if result.URL != "" {
		fmt.Printf("URL: %s\n", result.URL)
	}
	if result.State != "" {
		fmt.Printf("State: %s\n", result.State)
	}

	return nil
}

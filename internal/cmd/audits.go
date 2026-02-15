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

type AuditsListCmd struct {
	Page     int `help:"Page number" default:"1"`
	PageSize int `help:"Page size" default:"10"`
}

func (cmd *AuditsListCmd) Run(ctx context.Context) error {
	client, err := getSurferClient()
	if err != nil {
		return err
	}

	result, err := client.Audits().List(ctx, cmd.Page, cmd.PageSize)
	if err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, result)
	}

	if len(result.Data) == 0 {
		fmt.Fprintln(os.Stderr, "No audits found")
		return nil
	}

	fmt.Fprintf(os.Stderr, "Showing %d of %d audits (page %d)\n\n", len(result.Data), result.TotalItems, result.Page)
	for _, audit := range result.Data {
		fmt.Printf("ID: %s\n", audit.ID)
		if audit.URL != "" {
			fmt.Printf("  URL: %s\n", audit.URL)
		}
		if audit.State != "" {
			fmt.Printf("  State: %s\n", audit.State)
		}
		if audit.Score > 0 {
			fmt.Printf("  Score: %d\n", audit.Score)
		}
		fmt.Println()
	}

	return nil
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

	fmt.Fprintf(os.Stderr, "Created audit\n\n")
	fmt.Printf("ID: %s\n", result.ID)
	if result.URL != "" {
		fmt.Printf("URL: %s\n", result.URL)
	}
	if result.State != "" {
		fmt.Printf("State: %s\n", result.State)
	}

	return nil
}

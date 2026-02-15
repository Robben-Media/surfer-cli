package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/builtbyrobben/surfer-cli/internal/outfmt"
	"github.com/builtbyrobben/surfer-cli/internal/secrets"
	"github.com/builtbyrobben/surfer-cli/internal/surfer"
)

type EditorsCmd struct {
	List   EditorsListCmd   `cmd:"" help:"List all content editors"`
	Create EditorsCreateCmd `cmd:"" help:"Create a new content editor"`
	Get    EditorsGetCmd    `cmd:"" help:"Get a content editor by ID"`
	Score  EditorsScoreCmd  `cmd:"" help:"Get content score for an editor"`
}

type EditorsListCmd struct {
	Page     int `help:"Page number" default:"1"`
	PageSize int `help:"Page size" default:"10"`
}

func (cmd *EditorsListCmd) Run(ctx context.Context) error {
	client, err := getSurferClient()
	if err != nil {
		return err
	}

	result, err := client.ContentEditors().List(ctx, cmd.Page, cmd.PageSize)
	if err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, result)
	}

	if len(result.Data) == 0 {
		fmt.Fprintln(os.Stderr, "No content editors found")
		return nil
	}

	fmt.Fprintf(os.Stderr, "Showing %d of %d editors (page %d)\n\n", len(result.Data), result.TotalItems, result.Page)
	for _, editor := range result.Data {
		fmt.Printf("ID: %d\n", editor.ID)
		if len(editor.Keywords) > 0 {
			fmt.Printf("  Keywords: %s\n", strings.Join(editor.Keywords, ", "))
		}
		if editor.State != "" {
			fmt.Printf("  State: %s\n", editor.State)
		}
		if editor.Language != "" {
			fmt.Printf("  Language: %s\n", editor.Language)
		}
		fmt.Println()
	}

	return nil
}

type EditorsCreateCmd struct {
	Keywords string `required:"" help:"Comma-separated keywords"`
	Language string `help:"Language code (e.g., en, es, de)" default:"en"`
	Location string `help:"Location (e.g., United States)"`
	Device   string `help:"Device type (desktop or mobile)" default:"desktop"`
}

func (cmd *EditorsCreateCmd) Run(ctx context.Context) error {
	client, err := getSurferClient()
	if err != nil {
		return err
	}

	keywords := strings.Split(cmd.Keywords, ",")
	for i, kw := range keywords {
		keywords[i] = strings.TrimSpace(kw)
	}

	req := surfer.CreateContentEditorRequest{
		Keywords: keywords,
		Language: cmd.Language,
		Location: cmd.Location,
		Device:   cmd.Device,
	}

	result, err := client.ContentEditors().Create(ctx, req)
	if err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, result)
	}

	fmt.Fprintf(os.Stderr, "Created content editor\n\n")
	fmt.Printf("ID: %d\n", result.ID)
	if len(result.Keywords) > 0 {
		fmt.Printf("Keywords: %s\n", strings.Join(result.Keywords, ", "))
	}
	if result.URL != "" {
		fmt.Printf("URL: %s\n", result.URL)
	}

	return nil
}

type EditorsGetCmd struct {
	ID string `arg:"" required:"" help:"Content editor ID"`
}

func (cmd *EditorsGetCmd) Run(ctx context.Context) error {
	client, err := getSurferClient()
	if err != nil {
		return err
	}

	result, err := client.ContentEditors().Get(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, result)
	}

	fmt.Printf("ID: %d\n", result.ID)
	if len(result.Keywords) > 0 {
		fmt.Printf("Keywords: %s\n", strings.Join(result.Keywords, ", "))
	}
	if result.State != "" {
		fmt.Printf("State: %s\n", result.State)
	}
	if result.Language != "" {
		fmt.Printf("Language: %s\n", result.Language)
	}
	if result.Location != "" {
		fmt.Printf("Location: %s\n", result.Location)
	}
	if result.Device != "" {
		fmt.Printf("Device: %s\n", result.Device)
	}
	if result.URL != "" {
		fmt.Printf("URL: %s\n", result.URL)
	}

	return nil
}

type EditorsScoreCmd struct {
	ID string `arg:"" required:"" help:"Content editor ID"`
}

func (cmd *EditorsScoreCmd) Run(ctx context.Context) error {
	client, err := getSurferClient()
	if err != nil {
		return err
	}

	result, err := client.ContentEditors().Score(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if outfmt.IsJSON(ctx) {
		return outfmt.WriteJSON(os.Stdout, result)
	}

	fmt.Printf("Content Score: %d\n", result.ContentScore)

	return nil
}

func getSurferClient() (*surfer.Client, error) {
	// Check for environment variable override first
	apiKey := os.Getenv("SURFER_API_KEY")

	if apiKey == "" {
		// Try to get from keyring
		store, err := secrets.OpenDefault()
		if err != nil {
			return nil, fmt.Errorf("open credential store: %w", err)
		}

		apiKey, err = store.GetAPIKey()
		if err != nil {
			return nil, fmt.Errorf("get API key: %w (set SURFER_API_KEY or run 'surfer-cli auth set-key --stdin')", err)
		}
	}

	return surfer.NewClient(apiKey), nil
}

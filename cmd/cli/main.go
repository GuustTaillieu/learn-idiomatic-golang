package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	_ "modernc.org/sqlite"

	"github.com/spf13/cobra"

	"github.com/GuustTaillieu/idiomatic-go/internal/domain"
	database "github.com/GuustTaillieu/idiomatic-go/internal/sqlite"
)

func main() {
	rootCmd := NewRootCommand()

	if err := rootCmd.Execute(); err != nil {
		slog.Error("Command execution failed", "error", err)
		os.Exit(1)
	}
}

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "idiomatic-go",
		Short: "An idiomatic Go application",
	}

	rootCmd.AddCommand(SeedDatabaseCommand())

	return rootCmd
}

func SeedDatabaseCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "seed",
		Short: "Seed the database with initial data",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
			defer stop()

			db, err := sql.Open("sqlite", "file:orders.db?cache=shared&mode=rwc")
			if err != nil {
				slog.Error("Failed to open database", "error", err)
				os.Exit(1)
			}
			defer db.Close()

			itemStore, err := database.NewItemStore(db)
			if err != nil {
				slog.Error("Failed to create item store", "error", err)
				os.Exit(1)
			}
			inventoryStore, err := database.NewInventoryStore(db)
			if err != nil {
				slog.Error("Failed to create inventory store", "error", err)
				os.Exit(1)
			}

			// Seed the database with initial items
			items := []string{"item1", "item2", "item3"}
			for _, itemName := range items {
				item := domain.NewItem(itemName)
				if err := itemStore.Save(ctx, item); err != nil {
					slog.Error("Failed to save item", "error", err)
					os.Exit(1)
				}
				stock := domain.NewStock(item.ID, 10)
				if err := inventoryStore.AddStock(ctx, stock); err != nil {
					slog.Error("Failed to add stock", "error", err)
					os.Exit(1)
				}
			}

			fmt.Println("Database seeded successfully.")
		},
	}
}

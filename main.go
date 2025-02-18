package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgo/handler/middleware"
	"github.com/disgoorg/snowflake/v2"

	"github.com/lucastso10/bot_maid_nsfw_V2/bot/commands"
	"github.com/lucastso10/bot_maid_nsfw_V2/bot/components"
)

var (
	token   = os.Getenv("disgo_token")
	guildID = snowflake.GetEnv("disgo_guild_id")
)

func main() {
	slog.Info("starting bot_maid...")
	slog.Info("disgo version", slog.String("version", disgo.Version))

	application_handler := handler.New()
	application_handler.Use(middleware.Logger)

	// Uses the struct commands to declare all the commands on the list
	for _, command := range commands.Commands {

		slog.Info("Loading command:", slog.String("command", command.ApplicationCommand.CommandName()))

		application_handler.Command("/"+command.ApplicationCommand.CommandName(), command.Handler)

		// Only loads if there is an autocomplete function for the command
		if command.AutoComplete != nil {
			slog.Info("Loading autocomplete for command:", slog.String("command", command.ApplicationCommand.CommandName()))

			application_handler.Autocomplete("/"+command.ApplicationCommand.CommandName(), command.AutoComplete)
		}
	}

	// TODO: Make components declaration similar to commands
	application_handler.Component("/test-component", components.TestComponent)

	application_handler.NotFound(handleNotFound)

	client, err := disgo.New(
		token,
		bot.WithDefaultGateway(),
		bot.WithEventListeners(application_handler),
	)

	if err != nil {
		slog.Error("error while building bot", slog.Any("err", err))
		return
	}

	// Pulls all the commands
	var commands_to_sync []discord.ApplicationCommandCreate
	for _, command := range commands.Commands {
		commands_to_sync = append(commands_to_sync, command.ApplicationCommand)
	}

	if err = handler.SyncCommands(client, commands_to_sync, []snowflake.ID{guildID}); err != nil {
		slog.Error("Error while syncing commands", slog.Any("err", err))
		return
	}

	defer client.Close(context.TODO())

	if err = client.OpenGateway(context.TODO()); err != nil {
		slog.Error("Error while connecting to gateway", slog.Any("err", err))
	}

	slog.Info("bot_maid is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}

func handleNotFound(event *handler.InteractionEvent) error {
	return event.CreateMessage(discord.MessageCreate{Content: "⚠️ Erro! Commando não encontrado!"})
}

package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

type Command struct {
	ApplicationCommand discord.ApplicationCommandCreate
	Handler            func(e *handler.CommandEvent) error
	AutoComplete       func(e *handler.AutocompleteEvent) error
}

// Declare all the bots commands here
var Commands = []Command{
	{test, TestHandler, TestAutoComplete},
}

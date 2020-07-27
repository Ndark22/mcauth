package bot

import (
	"fmt"
	util "github.com/Floor-Gang/utilpkg/botutil"
	dg "github.com/bwmarrin/discordgo"
	"github.com/dhghf/mcauth/internal/common"
	"github.com/dhghf/mcauth/internal/common/db"
	"strconv"
)

// Ban a Discord user / Minecraft player
func (bot *Bot) cmdBan(msg *dg.MessageCreate, args []string) {
	// args is at least 3
	if len(args) < 3 {
		return
	}

	// if they mentioned a player
	// then -> args = [prefix, ban, @discord user]
	if len(msg.Mentions) > 0 {
		mentioned := msg.Mentions[0]
		playerID, _ := bot.store.Links.GetPlayerID(mentioned.ID)

		if len(playerID) == 0 {
			util.Reply(
				bot.client,
				msg.Message,
				fmt.Sprintf("%s isn't linked with anything", user.Mention()),
			)
			return
		}

		link := db.LinkedAcc{
			DiscordID: mentioned.ID,
			PlayerID:  playerID,
		}
		err := bot.store.Bans.Ban(link)

		if err != nil {
			util.Reply(
				bot.client,
				msg.Message,
				fmt.Sprintf(
					"%s (%s) is already banned",
					user.Mention(),
					playerID,
				),
			)
		} else {
			util.Reply(
				bot.client,
				msg.Message,
				fmt.Sprintf("%s (%s) is now banned", user.Mention(), playerID),
			)
		}
		return
	}

	// else -> args = [prefix, ban, mc player name]
	playerName := args[2]
	playerID := common.GetPlayerID(playerName)

	if len(playerID) == 0 {
		util.Reply(
			bot.client,
			msg.Message,
			fmt.Sprintf("%s isn't a valid player", playerName),
		)
		return
	}

	userID, _ := bot.store.Links.GetDiscordID(playerID)

	if len(userID) == 0 {
		util.Reply(
			bot.client,
			msg.Message,
			fmt.Sprintf("%s isn't linked with a user", playerName),
		)
		return
	}

	link := db.LinkedAcc{
		DiscordID: userID,
		PlayerID:  playerID,
	}
	err := bot.store.Bans.Ban(link)

	if err != nil {
		util.Reply(
			bot.client,
			msg.Message,
			fmt.Sprintf(
				"%s (%s) is already banned",
				user.Mention(),
				playerID,
			),
		)
	} else {
		util.Reply(
			bot.client,
			msg.Message,
			fmt.Sprintf("%s (%s) is now banned", user.Mention(), playerID),
		)
	}
}

// See the status of the bot
func (bot *Bot) cmdStatus(msg *dg.Message) {
	embed := &dg.MessageEmbed{
		Title: fmt.Sprintf("MCAuth Status [%s]", common.Version),
		URL:   "https://github.com/dhghf/mcauth",
		Color: 0xfc4646,
	}

	playerCount := bot.countPlayersOnline()
	playersOnline := &dg.MessageEmbedField{
		Name:   "Players Online",
		Value:  strconv.Itoa(playerCount),
		Inline: true,
	}

	linkedAccCount := bot.countLinkedAccounts()
	linkedAccounts := &dg.MessageEmbedField{
		Name:   "Linked Accounts",
		Value:  strconv.Itoa(linkedAccCount),
		Inline: true,
	}

	allPending := bot.countPendingAuthCodes()
	pendingAuthCodes := &dg.MessageEmbedField{
		Name:   "Pending Auth Codes",
		Value:  strconv.Itoa(allPending),
		Inline: true,
	}

	altAccsCount := bot.countAltAccounts()
	altAccsField := &dg.MessageEmbedField{
		Name:   "Alt Accounts",
		Value:  strconv.Itoa(altAccsCount),
		Inline: true,
	}

	whitelistedList := bot.getWhitelistedRoles()
	whitelisted := &dg.MessageEmbedField{
		Name:   "Whitelisted Roles",
		Value:  whitelistedList,
		Inline: true,
	}

	adminRolesList := bot.getAdminRoles()
	adminRoles := &dg.MessageEmbedField{
		Name:   "Admin Roles",
		Value:  adminRolesList,
		Inline: true,
	}

	embed.Fields = []*dg.MessageEmbedField{
		playersOnline, linkedAccounts, pendingAuthCodes,
		altAccsField, adminRoles, whitelisted,
	}

	_, err := bot.client.ChannelMessageSendEmbed(
		msg.ChannelID,
		embed,
	)

	if err != nil {
		log.Println("Failed to send status", err.Error())
	}
}
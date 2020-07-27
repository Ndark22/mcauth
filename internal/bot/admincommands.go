package bot

import (
	"fmt"
	dg "github.com/bwmarrin/discordgo"
	"github.com/dhghf/mcauth/internal/common"
	"strconv"
)

func cmdBan(msg *dg.Message, args []string) {

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

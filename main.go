package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Main funciton
func main() {
	//Gets the authfile key
	authfile, err := os.Open("authkey.txt")

	if err != nil {
		fmt.Println("authkey.txt is not able to be opened and/or does not exist", err)
	}
	defer authfile.Close()

	s := bufio.NewScanner(authfile)

	s.Scan()

	authkey := s.Text()

	start := "Bot " + authkey
	d, err := discordgo.New(start)

	if err != nil {
		fmt.Println("failed to create discord session", err)
	}

	bot, err := d.User("@me")

	if err != nil {
		fmt.Println("failed to access account", err)
	}
	username := bot.Username

	d.AddHandler(handleCmd)
	err = d.Open()

	if err != nil {
		fmt.Println("unable to establish connection", err)
	} else {
		fmt.Printf("%s is up and running! :D", username)
	}

	defer d.Close()

	<-make(chan struct{})
}

//handleCmd handles all the messages sent from discord.
func handleCmd(d *discordgo.Session, msg *discordgo.MessageCreate) {
	bot, _ := d.User("@me")

	user := msg.Author
	//Checks if the author of message is itself and if it is a bot.
	if user.ID == bot.ID || user.Bot {
		return
	}
	//Randomly generates a number between 1 - 10
	randomInt := rand.Intn(10)
	//Checks if user "Belle" sends a message, and if the random number generated was 5.
	//If so, react with a bell emote
	if user.ID == "263417510253035530" && randomInt == 5 {
		d.MessageReactionAdd(msg.ChannelID, msg.ID, "🔔")
	}
	content := msg.Content
	var command string

	prefix := "!"
	//Checks if the command starts with "!"
	//if it does, then remove the "!" and store it in command
	if strings.HasPrefix(content, prefix) {
		command = strings.TrimPrefix(content, prefix)
	} else {
		return
	}
	//Checks for emotes in the message
	var emotePattern, err = regexp.Compile("[0-9]+")
	if err != nil {
		fmt.Println(err)
		return
	}

	staticemote := "https://cdn.discordapp.com/emojis/%s.png"
	gifemote := "https://cdn.discordapp.com/emojis/%s.gif"
	var emote string

	//Ping checks if the command is online, if it is, then it sends pong back
	if command == "ping" {
		d.ChannelMessageSend(msg.ChannelID, "pong")
	} else if emotePattern.MatchString(command) {
		//If a user does emote an emote after "!" then it will enlarge it
		emoteID := emotePattern.FindString(command)
		if isValidLink(fmt.Sprintf(gifemote, emoteID)) {
			emote = gifemote

		} else {
			emote = staticemote

		}
		emoteURL := fmt.Sprintf(emote, emoteID)
		image := &discordgo.MessageEmbed{
			Image: &discordgo.MessageEmbedImage{
				URL: emoteURL,
			},
		}
		d.ChannelMessageDelete(msg.ChannelID, msg.ID)

		d.ChannelMessageSendEmbed(msg.ChannelID, image)
	}

}

//isValidLink checks if the link provided is valid through http GET request.
//returns true or false
func isValidLink(link string) bool {
	response, err := http.Get(link)
	if err != nil {
		fmt.Println("Link is not a http link", err)
	}
	defer response.Body.Close()
	if response.StatusCode == 200 {
		return true
	}
	return false
}

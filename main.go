package main

import (
    "fmt"
    "github.com/bwmarrin/discordgo"
    "os"
    "bufio"
)

// our main function
func main() {
    authfile,err := os.Open("authkey.txt");

    if err != nil {
      fmt.Print(err);
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
    }else{
      fmt.Printf("%s is up and running! :D",username)
    }

    defer d.Close()

    <-make(chan struct{})
}

// our command handler function
func handleCmd(d *discordgo.Session, msg *discordgo.MessageCreate) {
    bot,_ := d.User("@me")

    user := msg.Author
    if user.ID ==  bot.ID|| user.Bot {
        return
    }

    content := msg.Content

    if (content == "!test") {
        d.ChannelMessageSend(msg.ChannelID, "Testing..")
    }

    fmt.Printf("Message: %+v\n", msg.Message)
}

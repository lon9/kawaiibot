package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dshardmanager"
)

// SearchURL is url for searching images
const SearchURL string = "https://safebooru.org/index.php?page=dapi&s=post&q=index&json=1&limit=100&tags=%s"

// OriginalURL is original image url
const OriginalURL string = "https://safebooru.org/images/%s/%s"

var (
	emojis = []string{
		":hearts:",
		":stars:",
		":heart:",
		":yellow_heart:",
		":green_heart:",
		":purple_heart:",
		":blue_heart:",
		":two_hearts:",
		":revolving_hearts:",
		":heartbeat:",
		":heartpulse:",
		":sparkling_heart:",
		":cupid:",
		":black_heart:",
		":dark_sunglasses:",
		":dancer:",
	}
	colors = []int{
		0xe57373,
		0xf06292,
		0xba68c8,
		0x9575cd,
		0x7986cb,
		0x64b5f6,
		0x4fc3f7,
		0x4dd0e1,
		0x4db6ac,
		0x81c784,
		0xaed581,
		0xdce775,
		0xfff176,
		0xffd54f,
		0xffb74d,
		0xff8a65,
	}
)

// Image is object for image
type Image struct {
	ID        int    `json:"id"`
	Directory string `json:"directory"`
	Name      string `json:"image"`
}

// Bot is struct
type Bot struct {
	dshardmanager.Manager
	clientID string
}

// NewBot is constructor
func NewBot(token, clientID, env string) (*Bot, error) {
	b := &Bot{
		Manager:  *dshardmanager.New("Bot " + token),
		clientID: clientID,
	}
	b.Name = "kawaiibot"
	if env == "development" {
		b.SetNumShards(1)
	}
	b.AddHandler(b.ready)
	b.AddHandler(b.messageCreate)
	return b, b.Start()
}

// Close closes connections.
func (b *Bot) Close() error {
	return b.StopAll()
}

func (b *Bot) ready(s *discordgo.Session, event *discordgo.Ready) {
	s.UpdateStatus(0, "!moe (q={keyword})")
}

func (b *Bot) messageCreate(s *discordgo.Session, event *discordgo.MessageCreate) {
	m := strings.TrimSpace(event.Content)

	if m == "!moe invite" {
		s.ChannelMessageSend(event.ChannelID, fmt.Sprintf("https://discordapp.com/oauth2/authorize?client_id=%s&scope=bot&permissions=0", b.clientID))
		return
	}
	if strings.HasPrefix(m, "!moe q=") {
		b.sendEmbed(s, event.ChannelID, strings.TrimSpace(strings.Replace(m, "!moe q=", "", -1)))
		return
	}
	if strings.HasPrefix(m, "!moe") {
		b.sendEmbed(s, event.ChannelID, "")
		return
	}
}

func (b *Bot) sendEmbed(s *discordgo.Session, channelID, q string) {
	image, err := b.getImage(q)
	if err != nil {
		log.Println(err)
		s.ChannelMessageSend(channelID, fmt.Sprintf("There is no image for the query %s", q))
		return
	}
	emoji := emojis[rand.Intn(len(emojis))]
	s.ChannelMessageSendEmbed(
		channelID,
		NewEmbed().
			SetTitle(fmt.Sprintf("%s %s %s", emoji, emoji, emoji)).
			SetColor(colors[rand.Intn(len(colors))]).
			SetImage(fmt.Sprintf(OriginalURL, image.Directory, image.Name)).MessageEmbed,
	)
}

func (b *Bot) getImage(query string) (image *Image, err error) {
	var images []Image
	u := fmt.Sprintf(SearchURL, query)
	resp, err := http.Get(u)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&images)
	if len(images) == 0 {
		return nil, errors.New("No image")
	}
	image = &images[rand.Intn(len(images))]
	return
}

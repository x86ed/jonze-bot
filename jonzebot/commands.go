package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

//NowPlaying is the current movie playing
type NowPlaying struct {
	Start   time.Time
	Feature Feature
}

//Command struct for jonze commands
type Command struct {
	Trigger []string
	Action  string
	Values  []string
}

var playtime time.Time

//Parse maps a string to a command
func (c *Command) Parse(s string) (out bool) {
	for _, v := range c.Trigger {
		if strings.HasPrefix(s, v) {
			out = true
			fmt.Println(s, v)
			c.Action = actionMap[v]
			s = strings.Replace(s, v, "", -1)
			sa := strings.Split(s, " ")
			if len(sa) >= 1 {
				c.Values = sa
				if c.Action == sa[0] {
					c.Values = sa[1:]
				}
			}
		}
	}
	return
}

var help string = "```yaml\n" +
	`Commands:
* !jonze help/!jonzehelp - show this. 
* !jonze timestamp/!jonzets Video Name/01:01:11/tag name - will time tag a time in the selected video.
* !jonze add/!jonzeadd  url - will create an entry for a skate video at the url specfied.
* !jonze play/!jonzeplay video name/timestamp - will play the video requested at the timestamp specified.
* !jonze nominate/!jonzenom  feature film *&* skate video - nominate a feature film and skate video for the Sk8turday Matinee.
* !jonze vote/!jonzevote (number) - Without a number list the nominees for the week. With a number vote for that entry.
* !jonze list/!jonzelist  - list all movies in the vault.
* !jonze search/!jonzesearch string - search for a value in the vault.

* !jonze sk8/!sk8urday alert/play/schedule/delete int - movie party commands (video mod only)
` + "```"

var actionMap = map[string]string{
	"!jonze help":      "help",
	"!jonzehelp":       "help",
	"!jonze timestamp": "timestamp",
	"!jonzets":         "timestamp",
	"!jonze play":      "play",
	"!jonzeplay":       "play",
	"!jonze nominate":  "nominate",
	"!jonzenom":        "nominate",
	"!jonze vote":      "vote",
	"!jonzevote":       "vote",
	"!jonze add":       "add",
	"!jonzeadd":        "add",
	"!jonze list":      "list",
	"!jonzelist":       "list",
	"!jonze search":    "search",
	"!jonzesearch":     "search",
	"!sk8turday":       "sk8urday",
	"!jonze sk8":       "sk8urday",
}

var helpc = Command{
	Trigger: []string{"!jonze help", "!jonzehelp"},
	Action:  "help",
}

var timestampc = Command{
	Trigger: []string{"!jonze timestamp", "!jonzets"},
	Action:  "timestamp",
}

var playc = Command{
	Trigger: []string{"!jonze play", "!jonzeplay"},
	Action:  "play",
}

var nominatec = Command{
	Trigger: []string{"!jonze nominate", "!jonzenom"},
	Action:  "nominate",
}

var votec = Command{
	Trigger: []string{"!jonze vote", "!jonzevote"},
	Action:  "vote",
}

var listc = Command{
	Trigger: []string{"!jonze list", "!jonzelist"},
	Action:  "list",
}

var searchc = Command{
	Trigger: []string{"!jonze search", "!jonzesearch"},
	Action:  "search",
}

var addc = Command{
	Trigger: []string{"!jonze add", "!jonzeadd"},
	Action:  "add",
}

var sk8urdayc = Command{
	Trigger: []string{"!sk8urday", "!jonze sk8"},
	Action:  "sk8urday",
}

var commands = []Command{
	helpc,
	timestampc,
	playc,
	nominatec,
	votec,
	addc,
	sk8urdayc,
}

//Process processes a command object
func (c *Command) Process(s *discordgo.Session, m *discordgo.MessageCreate) {
	cont := m.Content
	short := c.Parse(cont)
	if !short {
		return
	}
	switch c.Action {
	case "add":
		c.add(s, m)
	case "help":
		c.help(s, m)
	case "nominate":
		c.nominate(s, m)
	case "play":
		c.play(s, m)
	case "timestamp":
		c.timestamp(s, m)
	case "list":
		c.list(s, m)
	case "search":
		c.search(s, m)
	case "vote":
		c.vote(s, m)
	case "sk8urday":
		c.sk8(s, m)
	}
}

func (c *Command) add(s *discordgo.Session, m *discordgo.MessageCreate) {
	if len(c.Values) < 2 {
		s.ChannelMessageSend(m.ChannelID, "nah. needs more info.")
		return
	}
	URL := c.Values[1]
	// n := strings.Join(c.Values[2:], " ")
	new := Video{
		URL:       URL,
		TimeCodes: map[string]Timecode{},
	}
	new.getName()
	if len(new.Name) > 0 {
		vault[new.URL] = new
		archiveJSON(os.Getenv("VIDEOVAULT"), &vault)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s was added to the film vault.", new.Name))
		return
	}

	s.ChannelMessageSend(m.ChannelID, "Whelp, that didn't work.")
}

func (c *Command) help(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, help)
}

func (c *Command) nominate(s *discordgo.Session, m *discordgo.MessageCreate) {
	ts := strings.Join(c.Values, " ")
	t := strings.Split(ts, "*&*")
	if len(t) < 2 {
		s.ChannelMessageSend(m.ChannelID, "nah. needs more films.")
		return
	}
	ms := t[0]
	ms = strings.TrimSpace(ms)
	sv := t[1]
	sv = strings.TrimSpace(sv)
	fmt.Println("urls ", ms, sv)
	f := Feature{}
	if isValidURL(ms) {
		var new Video
		if len(vault[ms].Name) > 0 {
			new = vault[ms]
		} else {
			new = Video{URL: ms, TimeCodes: map[string]Timecode{}}
			new.getName()
			vault[new.URL] = new
			archiveJSON(os.Getenv("VIDEOVAULT"), &vault)
		}
		f.Movie = new.Name
		f.MovieURL = new.URL
	} else {
		f.Movie = ms
	}

	if isValidURL(sv) {
		var new Video
		if len(vault[sv].Name) > 0 {
			new = vault[sv]
		} else {
			new = Video{URL: sv, TimeCodes: map[string]Timecode{}}
			new.getName()
			vault[new.URL] = new
			archiveJSON(os.Getenv("VIDEOVAULT"), &vault)
		}
		f.SkateVid = new.Name
		f.SkateURL = new.URL
	} else {
		f.SkateVid = sv
	}
	if !checkCombo(f.Movie, f.SkateVid) {
		playbill = append(playbill, f)
		archiveJSON(os.Getenv("PLAYBILL"), &playbill)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("**%s & %s** was added to the playbill.", f.Movie, f.SkateVid))
		return
	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("**%s & %s** already is on the playbill.", f.Movie, f.SkateVid))
}

func (c *Command) play(s *discordgo.Session, m *discordgo.MessageCreate) {
	if len(c.Values) < 2 {
		vals := listVault()
		for _, v := range vals {
			s.ChannelMessageSend(m.ChannelID, v)
		}
		return
	}
	rt := strings.Join(c.Values[1:], " ")
	p := strings.Split(rt, "/")
	if len(p) < 1 {
		s.ChannelMessageSend(m.ChannelID, "We don't have that one yet.")
		return
	}
	for _, v := range vault {
		if strings.Index(strings.ToLower(v.Name), strings.ToLower(p[0])) > -1 {
			fURL := "Check it out...\n"
			fURL += fmt.Sprintf("**%s**\n", v.Name)
			qSep := "?t="
			tc := ""
			if strings.Index(fURL, "?") > -1 {
				qSep = "&t="
			}
			if v.Service == "vimeo" {
				qSep = "#t="
			}
			if len(p) > 1 && len(v.TimeCodes[p[1]].Desc) > 1 {
				fURL += fmt.Sprintf("Part: %s\n", p[1])
				code := v.TimeCodes[p[1]]
				tc = qSep + code.toString()
			}
			fURL += v.URL + tc
			s.ChannelMessageSend(m.ChannelID, fURL)
			return
		}
	}
	s.ChannelMessageSend(m.ChannelID, "Couldn't find it.")
}

func (c *Command) timestamp(s *discordgo.Session, m *discordgo.MessageCreate) {
	if len(c.Values) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Not enough info. See ya!")
		return
	}
	rts := strings.Split(strings.Join(c.Values[1:], " "), "/")
	if len(rts) != 3 {
		s.ChannelMessageSend(m.ChannelID, "I don't know what to do with that. Later!")
		return
	}
	new := Timecode{
		Desc: rts[2],
	}
	new.fromString(rts[1])
	var key string
	for _, v := range vault {
		if strings.Index(strings.ToLower(v.Name), strings.ToLower(rts[0])) > -1 {
			key = v.URL
		}
	}
	if len(key) < 1 {
		s.ChannelMessageSend(m.ChannelID, "Couldn't find that film")
		return
	}
	match := vault[key]
	fmt.Println("vault", vault)
	fmt.Println("match", match, rts)
	match.TimeCodes[rts[2]] = new
	vault[key] = match
	archiveJSON(os.Getenv("VIDEOVAULT"), &vault)
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("**%s** was added to **%s**.", rts[2], rts[0]))
}

func (c *Command) vote(s *discordgo.Session, m *discordgo.MessageCreate) {
	if len(c.Values) < 2 {
		movies := listPlaybill()
		for _, v := range movies {
			s.ChannelMessageSend(m.ChannelID, v)
		}
		return
	}
	i, e := strconv.Atoi(c.Values[1])
	if e != nil {
		s.ChannelMessageSend(m.ChannelID, "That's not a valid number.")
		return
	}
	ov := len(playbill[i-1].Votes)
	playbill[i-1].addVote(m.Author.ID)
	if ov == len(playbill[i-1].Votes) {
		s.ChannelMessageSend(m.ChannelID, "Hey! You already voted for that.")
		return
	}
	archiveJSON(os.Getenv("PLAYBILL"), &playbill)
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s & %s has **%d** votes.", playbill[i-1].Movie, playbill[i-1].SkateVid, len(playbill[i-1].Votes)))
	s.ChannelMessageSend(m.ChannelID, currentTopVotes())
}

func checkRole(m *discordgo.MessageCreate) bool {
	role := "324575381581463553"
	var badRoles = []string{
		"359852475181694976",
		"636246344855453696",
		"716536970309927033",
		"416424719487860736",
		"697875957964603403",
	}
	if containsVal(m.Message.Member.Roles, role) > -1 {
		for _, v := range badRoles {
			if containsVal(m.Message.Member.Roles, v) > -1 {
				return false
			}
		}
		return true
	}
	return false
}

func (c *Command) sk8(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !checkRole(m) {
		s.ChannelMessageSend(m.ChannelID, "Oof. You don't have permission to do this.")
	}
	if currentMovie.Start.Local().Add(6*time.Hour).Unix() < time.Now().Unix() {
		currentMovie = NowPlaying{}
		archiveJSON(os.Getenv("NOWPLAYING"), &currentMovie)
	}
	if len(c.Values) < 2 {
		s.ChannelMessageSend(m.ChannelID, currentTopVotes())
		return
	}
	switch c.Values[1] {
	case "play":
		if len(c.Values) > 2 {
			c.sk8play(s, m)
			return
		}
	case "delete":
		if len(c.Values) > 2 {
			c.sk8del(s, m)
			return
		}
	case "schedule":
		if len(c.Values) > 3 {
			c.sk8sched(s, m)
			return
		}
	case "alert":
		c.sk8alert(s, m)
		return
	}
	s.ChannelMessageSend(m.ChannelID, "I don't know what you're trying to do.")
}

func (c *Command) sk8play(s *discordgo.Session, m *discordgo.MessageCreate) {
	l := "Jan 2, 2006 at 3:04pm (MST)"
	i, e := strconv.Atoi(c.Values[2])
	if e != nil {
		s.ChannelMessageSend(m.ChannelID, "That's not a valid number.")
		return
	}
	i--
	play := playbill[i]
	currentMovie = NowPlaying{
		Start:   playtime,
		Feature: play,
	}
	archiveJSON(os.Getenv("NOWPLAYING"), &currentMovie)
	copy(playbill[i:], playbill[i+1:])
	playbill[len(playbill)-1] = Feature{}
	playbill = playbill[:len(playbill)-1]
	s.ChannelMessageSend("715651111297613905", fmt.Sprintf("%s & %s is the next Sk8turday movie. Join us at %s in #Voice-chat to watch., ", play.Movie, play.SkateVid, playtime.Format(l)))
}

func (c *Command) sk8del(s *discordgo.Session, m *discordgo.MessageCreate) {
	i, e := strconv.Atoi(c.Values[2])
	if e != nil {
		s.ChannelMessageSend(m.ChannelID, "That's not a valid number.")
		return
	}
	i--
	play := playbill[i]
	copy(playbill[i:], playbill[i+1:])
	playbill[len(playbill)-1] = Feature{}
	playbill = playbill[:len(playbill)-1]
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s & %s was deleted.", play.Movie, play.SkateVid))
}

func (c *Command) sk8sched(s *discordgo.Session, m *discordgo.MessageCreate) {
	l := "Jan 2, 2006 at 3:04pm (MST)"
	ts := strings.Join(c.Values[1:], " ")
	t, e := time.Parse(l, ts)
	if e != nil {
		s.ChannelMessageSend(m.ChannelID, "Bad timestamp.")
		return
	}
	playtime = t
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Join us at %s in #Voice-chat to watch the next Sk8urday Movie., ", currentMovie.Start.Format(l)))
}

func (c *Command) sk8alert(s *discordgo.Session, m *discordgo.MessageCreate) {
	l := "Jan 2, 2006 at 3:04pm (MST)"
	if len(currentMovie.Feature.Movie) > 0 {
		s.ChannelMessageSend("715651111297613905", fmt.Sprintf("%s & %s is the next Sk8turday movie. Join us at %s in #Voice-chat to watch., ", currentMovie.Feature.Movie, currentMovie.Feature.SkateVid, currentMovie.Start.Format(l)))
	}
}

func (c *Command) search(s *discordgo.Session, m *discordgo.MessageCreate) {
}

func (c *Command) list(s *discordgo.Session, m *discordgo.MessageCreate) {
	vals := listVault()
	for _, v := range vals {
		s.ChannelMessageSend(m.ChannelID, v)
	}
	return
}

package main

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

//Feature movie votes object
type Feature struct {
	Movie    string
	MovieURL string
	SkateVid string
	SkateURL string
	Votes    []string
}

type features []Feature

func (a features) Len() int           { return len(a) }
func (a features) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a features) Less(i, j int) bool { return len(a[i].Votes) < len(a[j].Votes) }

func (f *Feature) addSkate(v Video) {
	f.SkateVid = v.Name
	f.SkateURL = v.URL
}

func (f *Feature) addMovie(v Video) {
	f.Movie = v.Name
	f.MovieURL = v.URL
}

func (f *Feature) addVote(s string) {
	var already bool
	for _, v := range f.Votes {
		if v == s {
			already = true
		}
	}
	if !already {
		f.Votes = append(f.Votes, s)
	}
}

func (f *Feature) getVotes() int {
	return len(f.Votes)
}

func (f *Feature) resetVote() {
	f.Votes = []string{}
}

//Video online video object
type Video struct {
	Name      string
	URL       string
	Service   string
	TimeCodes map[string]Timecode
}

func (v *Video) getName() {
	if len(v.URL) < 1 {
		return
	}
	if len(v.Service) < 1 {
		v.getService()
		if len(v.Service) < 1 {
			return
		}
	}
	t, err := getTitle(v.URL)
	if err != nil {
		return
	}
	v.Name = t
}

func (v *Video) getService() {
	if len(v.URL) < 1 {
		return
	}
	if strings.Index(v.URL, "youtu") > -1 {
		v.Service = "youtube"
	}
	if strings.Index(v.URL, "vimeo") > -1 {
		v.Service = "vimeo"
	}
}

//Timecode is a timestamp for a particular label
type Timecode struct {
	Offset time.Duration
	Desc   string
}

func (t *Timecode) toString() string {
	return t.Offset.String()
}

func (t *Timecode) fromString(s string) {
	sp := strings.Split(s, ":")
	var out string
	if len(sp) > 3 {
		return
	}
	for len(sp) < 3 {
		sp = append([]string{"0"}, sp...)
	}
	for i := len(sp) - 1; i > -1; i-- {
		hms := []string{"h", "m", "s"}
		out += sp[i] + hms[i]
	}
	val, err := time.ParseDuration(out)
	if err != nil {
		return
	}
	t.Offset = val
}

func listVault() (out []string) {
	out[0] += "*Movies:*\n"
	i := 1
	for _, v := range vault {
		chunk := i / 25
		out[chunk] += fmt.Sprintf("%d. **%s**\n", i, v.Name)
		i++
	}
	return
}

func listPlaybill() (out []string) {
	out[0] += "*Features:*\n"
	for i, v := range playbill {
		chunk := i / 25
		out[chunk] += fmt.Sprintf("%d. **%s** & **%s**\n", i+1, v.Movie, v.SkateVid)
	}
	return
}

func currentTopVotes() (out string) {
	f := features(playbill)
	sort.Sort(f)
	playbill = []Feature(f)
	for i := 0; i < len(playbill); i++ {
		out += fmt.Sprintf("%d. %s & %s - %d votes.\n", i+1, playbill[i].Movie, playbill[i].SkateVid, len(playbill[i].Votes))
		if i > 2 {
			return
		}
	}
	return
}

var playbill = []Feature{}
var vault = map[string]Video{}
var currentMovie = NowPlaying{}

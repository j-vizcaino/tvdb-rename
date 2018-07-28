package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/manifoldco/promptui"
	"github.com/nemith/tvdb"
	"path"
)

type EpisodeFile struct {
	Filename      string
	SeasonNumber  int
	EpisodeNumber int
}

// Look for S00E00 patterns
var reSeasonEpisode = regexp.MustCompile(`\b(?i:s(\d+)e(\d+))\b`)

func NewEpisodeFile(filename string) (EpisodeFile, error) {
	matches := reSeasonEpisode.FindStringSubmatch(filename)
	if len(matches) != (reSeasonEpisode.NumSubexp() + 1) {
		return EpisodeFile{}, fmt.Errorf("unable to find S..E.. pattern in filename")
	}

	season, _ := strconv.ParseInt(matches[1], 10, 32)
	episode, _ := strconv.ParseInt(matches[2], 10, 32)
	return EpisodeFile{
		Filename:      filename,
		SeasonNumber:  int(season),
		EpisodeNumber: int(episode),
	}, nil
}

func getSeriesName() string {
	prompt := promptui.Prompt{
		Label: "Search TV show",
	}
	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}
	return result
}

func getLanguage() string {
	prompt := promptui.Select{
		Label: "TV show language",
		Items: []string{"en", "fr"},
	}
	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}
	return result
}

func findTvShow(clt *tvdb.Client) (*tvdb.Series, error) {
	showName := getSeriesName()
	language := getLanguage()

	fmt.Println("Querying The TVDB...")
	series, err := clt.SearchSeries(showName, language)
	if err != nil {
		return nil, err
	}

	if len(series) == 0 {
		return nil, fmt.Errorf("no TV show found")
	}

	templates := promptui.SelectTemplates{
		Active:   `+ {{ .Name | cyan | bold }}`,
		Inactive: `  {{ .Name | cyan }}`,
		Selected: `{{ "✔" | green | bold }} {{ "Selected TV show" | bold }}: {{ .Name | cyan }}`,
	}

	prompt := promptui.Select{
		Label:     "Select TV show from list: ",
		Items:     series,
		Templates: &templates,
	}

	idx, _, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	return clt.SeriesByID(series[idx].ID, language)
}

func rename(epFile EpisodeFile, clt *tvdb.Client, show *tvdb.Series) {
	ep, err := clt.EpisodeBySeries(show.ID, epFile.SeasonNumber, epFile.EpisodeNumber, show.Language)
	if err != nil {
		fmt.Printf("Cannot find %s episode %d for season %d: %s", show.Name, epFile.EpisodeNumber, epFile.SeasonNumber, err)
		return
	}

	base := path.Base(epFile.Filename)
	if base == "." {
		base = ""
	}
	finalName := fmt.Sprintf("%s%s - S%02dE%02d - %s.%s", base, show.Name, ep.SeasonNumber, ep.EpisodeNumber, ep.EpisodeName, path.Ext(epFile.Filename))
	fmt.Printf("%s -> %s", epFile.Filename, finalName)
}

func main() {
	clt := tvdb.NewClient("")

	show, err := findTvShow(clt)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	toProcess := make([]EpisodeFile, 0, len(os.Args[1:]))
	for _, inFile := range os.Args[1:] {
		ep, err := NewEpisodeFile(inFile)
		if err != nil {
			fmt.Printf("%s: %s", inFile, err)
			continue
		}
		toProcess = append(toProcess, ep)
	}

	for _, f := range toProcess {
		rename(f, clt, show)
	}
}

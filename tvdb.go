package main

import (
	"fmt"
	"github.com/j-vizcaino/tvdb"
	"github.com/manifoldco/promptui"
	"os"
)

type TVDB struct {
	client      *tvdb.Client
	seriesCache map[string]*tvdb.Series
}

func NewTVDB(lang string) (*TVDB, error) {
	clt, err := tvdb.NewClient(tvdb.ClientOptions{
		APIKey: os.Getenv("TVDB_API_KEY"),
	})
	if err != nil {
		return nil, err
	}

	return &TVDB{
		client:      clt.WithLanguage(lang),
		seriesCache: make(map[string]*tvdb.Series),
	}, nil

}

func promptSeries(name string, series []tvdb.SeriesSearchResult) int {
	templates := promptui.SelectTemplates{
		Active:   `+ {{ .SeriesName | cyan | bold }}`,
		Inactive: `  {{ .SeriesName | cyan }}`,
		Selected: `{{ "✔" | green | bold }} {{ "Selected TV show" | bold }}: {{ .SeriesName | cyan }}`,
	}

	prompt := promptui.Select{
		Label:     fmt.Sprintf("The TVDB found the following series for %s, select the right one: ", name),
		Items:     series,
		Templates: &templates,
	}

	idx, _, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}
	return idx
}

func (t *TVDB) FindSeries(name string) (*tvdb.Series, error) {
	existing, ok := t.seriesCache[name]
	if ok {
		return existing, nil
	}

	series, err := t.client.SearchSeriesByName(name)
	if err != nil {
		return nil, err
	}

	idx := 0
	if len(series) == 0 {
		return nil, nil
	}

	if len(series) > 1 {
		idx = promptSeries(name, series)
	}

	show, err := t.client.SeriesByID(series[idx].ID)
	if err != nil {
		return nil, err
	}

	t.seriesCache[name] = show
	return show, nil
}

func (t *TVDB) FindEpisode(seriesID int, season int, episode int) (*tvdb.Episode, error) {
	episodes, err := t.client.EpisodesBySeriesID(seriesID,
		tvdb.WithAiredSeasonNumber(season),
		tvdb.WithAiredEpisodeNumber(episode),
	)
	if err != nil {
		return nil, err
	}

	if len(episodes) != 1 {
		return nil, fmt.Errorf("multiple episode candidates (%d) for season %d, episode %d", len(episodes), season, episode)
	}

	return &episodes[0], nil
}

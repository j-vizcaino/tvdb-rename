package main

import (
	"fmt"
	"os"
)

func prepareFile(t *TVDB, filepath string) *FileMetaData {
	fm := NewFileMetaData(filepath)

	series, err := t.FindSeries(fm.SeriesName)
	if err != nil {
		return nil
	}
	if series == nil {
		fmt.Printf("Unable to find a TVDB series matching '%s'.\n", fm.SeriesName)
		return nil
	}
	fm.SeriesName = series.SeriesName

	episode, err := t.FindEpisode(series.ID, fm.SeasonNumber, fm.EpisodeNumber)
	if err != nil {
		fmt.Printf("Unable to retrieve episode information: %s", err)
		return nil
	}
	fm.Title = episode.EpisodeName

	return &fm
}

func main() {
	t, err := NewTVDB("en")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, inFile := range os.Args[1:] {
		f := prepareFile(t, inFile)
		if f != nil {
			fmt.Println(*f)
		}
	}
}

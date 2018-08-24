package main

import (
	"fmt"
	"html/template"
	"os"
	"strings"
)

var (
	filenameFormat   = `{{ .SeriesName }} - S{{ .SeasonNumber | printf "%02d" }}E{{ .EpisodeNumber | printf "%02d" }} - {{ .Title }} {{.Source }}-{{ .Quality }}`
	filenameTemplate = template.Must(template.New("filename").Parse(filenameFormat))
)

func rename(f *FileMetaData) {
	var buf strings.Builder
	err := filenameTemplate.Execute(&buf, *f)
	if err != nil {
		fmt.Println(err)
		return
	}
	finalName := fmt.Sprintf("%s/%s%s",
		f.Directory,
		buf.String(),
		f.Extension)
	if finalName == f.Path {
		return
	}
	fmt.Printf("%40s -> %40s\n", f.Path, finalName)
	if err = os.Rename(f.Path, finalName); err != nil {
		fmt.Println(err)
	}
}

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
			rename(f)
		}
	}
}

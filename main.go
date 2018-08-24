package main

import (
	"fmt"
	"github.com/manifoldco/promptui"
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

func promptForQuality() string {
	templates := promptui.PromptTemplates{
		Prompt:  ` + {{ . }}`,
		Valid:   ` {{ "✔" | green | bold }} {{ "Quality" | bold }}: {{ . }}`,
		Success: ` {{ "✔" | green | bold }} {{ "Quality" | bold }}: {{ . }}`,
	}
	prompt := promptui.Prompt{
		Label:     "Unknown quality, please provide one: ",
		Templates: &templates,
	}

	res, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}
	return res
}

func promptForSource() string {
	templates := promptui.SelectTemplates{
		Active:   ` + {{ . | cyan | bold }}`,
		Inactive: `   {{ .| cyan }}`,
		Selected: ` {{ "✔" | green | bold }} {{ "Selected source" | bold }}: {{ . | cyan }}`,
	}

	prompt := promptui.Select{
		Label:     "Unknown source, select one from list: ",
		Items:     sourcesSet,
		Templates: &templates,
	}

	idx, _, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}
	return sourcesSet[idx]
}

func prepareFile(t *TVDB, filepath string) *FileMetaData {
	fmt.Printf(" * Peparing file '%s'...\n", filepath)
	fm := NewFileMetaData(filepath)

	if fm.SeriesName == "" {
		fmt.Println("[!] Could not detect series name from file.")
		return nil
	}

	series, err := t.FindSeries(fm.SeriesName)
	if err != nil {
		fmt.Printf("[!] Unable to query TVDB: %s\n", err)
		return nil
	}
	if series == nil {
		fmt.Printf("[!] Unable to find a TVDB series matching '%s'.\n", fm.SeriesName)
		return nil
	}
	fm.SeriesName = series.SeriesName

	episode, err := t.FindEpisode(series.ID, fm.SeasonNumber, fm.EpisodeNumber)
	if err != nil {
		fmt.Printf("[!] Unable to retrieve episode information: %s", err)
		return nil
	}
	fm.Title = episode.EpisodeName

	if fm.Source == "" {
		fm.Source = promptForSource()
	}

	if fm.Quality == "" {
		fm.Quality = promptForQuality()
	}
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

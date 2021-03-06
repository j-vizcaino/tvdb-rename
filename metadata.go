package main

import (
	"path"
	"regexp"
	"strconv"
	"strings"
)

var (
	sourcesMap map[string]string
	sourcesSet []string
)

func init() {
	sourcesMap = map[string]string{
		"DVD":    "DVD",
		"DVDRip": "DVD",
		"HDTV":   "HDTV",
		"SDTV":   "SDTV",
		"BDRip":  "BluRay",
		"BluRay": "BluRay",
		"WEB":    "WEB",
		"WEB-DL": "WEB",
		"WEBDL":  "WEB",
		"WEBRip": "WEB",
	}

	m := make(map[string]interface{})
	for _, v := range sourcesMap {
		m[v] = nil
	}
	sourcesSet = make([]string, 0, len(m))
	for v, _ := range m {
		sourcesSet = append(sourcesSet, v)
	}
}

type FileMetaData struct {
	Path      string
	Filename  string
	Directory string
	Extension string
	EpisodeMetaData
}

type EpisodeMetaData struct {
	SeriesName    string
	Title         string
	Source        string
	Quality       string
	SeasonNumber  int
	EpisodeNumber int
}

func NewFileMetaData(filepath string) FileMetaData {
	filename := path.Base(filepath)
	ext := path.Ext(filename)
	file := filename[0 : len(filename)-len(ext)]
	fm := FileMetaData{
		Path:            filepath,
		Filename:        filename,
		Extension:       ext,
		Directory:       path.Dir(filepath),
		EpisodeMetaData: NewEpisodeMetaData(file),
	}
	return fm
}

func NewEpisodeMetaData(s string) EpisodeMetaData {
	splitter := regexp.MustCompile(`[^A-Za-z0-9]+`)
	elts := splitter.Split(s, -1)
	season, episode, idx := findSeasonAndEpisode(elts)

	var seriesName string
	if idx > 0 {
		seriesName = strings.Join(elts[0:idx], " ")
	}
	md := EpisodeMetaData{
		SeriesName:    seriesName,
		Source:        findSource(elts),
		Quality:       findQuality(elts),
		SeasonNumber:  season,
		EpisodeNumber: episode,
	}

	return md
}

func findSeasonAndEpisode(elts []string) (int, int, int) {
	// Look for S00E00 patterns
	reSeasonEpisode := regexp.MustCompile(`^(?i:s(\d+)e(\d+))$`)

	for idx, elt := range elts {
		matches := reSeasonEpisode.FindStringSubmatch(elt)
		if len(matches) != (reSeasonEpisode.NumSubexp() + 1) {
			continue
		}

		season, _ := strconv.ParseInt(matches[1], 10, 32)
		episode, _ := strconv.ParseInt(matches[2], 10, 32)

		return int(season), int(episode), idx
	}
	return -1, -1, -1
}

func findSource(elts []string) string {
	for _, elt := range elts {
		for source, finalSource := range sourcesMap {
			if strings.EqualFold(elt, source) {
				return finalSource
			}
		}
	}
	return ""
}

func findQuality(elts []string) string {
	re := regexp.MustCompile(`[1-9][0-9]+p`)
	for _, elt := range elts {
		if re.MatchString(elt) {
			return elt
		}
	}
	return ""
}

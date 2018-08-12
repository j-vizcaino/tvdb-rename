package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewEpisodeMetaData(t *testing.T) {
	type testData struct {
		input    string
		expected EpisodeMetaData
	}

	data := []testData{
		{
			input: "Gotham.S04E19.BDRip.x264-DEMAND",
			expected: EpisodeMetaData{
				SeriesName:    "Gotham",
				SeasonNumber:  4,
				EpisodeNumber: 19,
				Source:        "BluRay",
			},
		}, {
			input: "Texas.Flip.N.Move.S09E13.Fine.Design.vs.Rough.Cabin.1080p.WEB.h264-CAFFEiNE",
			expected: EpisodeMetaData{
				SeriesName:    "Texas Flip N Move",
				SeasonNumber:  9,
				EpisodeNumber: 13,
				Quality:       "1080p",
				Source:        "WEB",
			},
		}, {
			input: "The Americans (2013) - S06E05 - The Great Patriotic War WEBDL-720p",
			expected: EpisodeMetaData{
				SeriesName:    "The Americans 2013",
				SeasonNumber:  6,
				EpisodeNumber: 5,
				Quality:       "720p",
				Source:        "WEB",
			},
		},
	}

	for _, d := range data {
		t.Run(d.input, func(tt *testing.T) {
			assert.Equal(tt, NewEpisodeMetaData(d.input), d.expected)
		})
	}
}

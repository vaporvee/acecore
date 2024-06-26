package custom

import (
	"strconv"
	"strings"
)

var color map[string]string = map[string]string{
	"red":     "#FF6B6B",
	"yellow":  "#FFD93D",
	"orange":  "#F6995C",
	"green":   "#6BCB77",
	"blue":    "#008DDA",
	"primary": "#211951",
}

const Gh_url string = "https://github.com/vaporvee/acecore/blob/main/"

func GetColor(s string) int {
	hexColor := strings.TrimPrefix(color[s], "#")
	decimal, err := strconv.ParseInt(hexColor, 16, 64)
	if err != nil {
		return 0
	}
	return int(decimal)
}

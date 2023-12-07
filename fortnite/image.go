package fortnite

import (
	"bytes"
	"image"
	"image/draw"
	"image/png"
	"io"
	"net/http"

	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/storage"
	"github.com/nfnt/resize"
)

var (
	PlaylistImages = map[string][]byte{
		"Playlist_DefaultSolo": {},
		// "Playlist_DefaultDuo": {},
		// "Playlist_DefaultSquad": {},
	}
)

func GenerateSoloImage() {
	background := *storage.Asset("background.png")

	itemFound := Cosmetics.GetRandomItemByType("AthenaCharacter")
	for itemFound.Images.Featured == "" {
		itemFound = Cosmetics.GetRandomItemByType("AthenaCharacter")
	}
	aid.Print(itemFound.Images.Featured)

	res, err := http.Get(itemFound.Images.Featured)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	bg, _, err := image.Decode(bytes.NewReader(background))
	if err != nil {
		panic(err)
	}

	soloPlayer, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}

	m := image.NewRGBA(bg.Bounds())
	draw.Draw(m, m.Bounds(), bg, image.Point{0, 0}, draw.Src)

	resized := resize.Resize(0, uint(float64(m.Bounds().Dy()) * 1.4), soloPlayer, resize.Lanczos3)
	centre := image.Point{
		m.Bounds().Dx()/2 - resized.Bounds().Dx()/2,
		(m.Bounds().Dy()/2 - resized.Bounds().Dy()/2) + 200,
	}
	draw.Draw(m, resized.Bounds().Add(centre), resized, image.Point{0, 0}, draw.Over)

	var bytes bytes.Buffer
	err = png.Encode(&bytes, m)
	if err != nil {
		panic(err)
	}

	PlaylistImages["Playlist_DefaultSolo"] = bytes.Bytes()
}
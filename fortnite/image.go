package fortnite

import (
	"bytes"
	"image"
	"image/draw"
	"image/png"
	"io"
	"math"
	"net/http"
	"strings"

	"github.com/ectrc/snow/aid"
	"github.com/ectrc/snow/storage"
	"github.com/nfnt/resize"
)

var (
	PlaylistImages = map[string][]byte{
		"Playlist_DefaultSolo": {},
		"Playlist_DefaultDuo": {},
		"Playlist_DefaultTrio": {},
		"Playlist_DefaultSquad": {},
	}
)

type colours struct {
	averageRed   uint8
	averageGreen uint8
	averageBlue  uint8
}

var SETS_NOT_ALLOWED = []string{
	"Soccer",
	"Football",
	"Waypoint",
}

func getAverageColour(img image.Image) colours {
	var red, green, blue uint64
	var count uint64

	bounds := img.Bounds()
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			r, g, b, _ := img.At(x, y).RGBA()
			red += uint64(r)
			green += uint64(g)
			blue += uint64(b)
			count++
		}
	}

	return colours{
		averageRed:   uint8(red / count),
		averageGreen: uint8(green / count),
		averageBlue:  uint8(blue / count),
	}
}

func GetAverageHexColour(img image.Image) string {
	colour := getAverageColour(img)
	return "#" + aid.ToHex(int(colour.averageRed)) + aid.ToHex(int(colour.averageGreen)) + aid.ToHex(int(colour.averageBlue))
}

func colorDifference(c1, c2 colours) float64 {
	diffRed := int(c1.averageRed) - int(c2.averageRed)
	diffGreen := int(c1.averageGreen) - int(c2.averageGreen)
	diffBlue := int(c1.averageBlue) - int(c2.averageBlue)

	return math.Sqrt(float64(diffRed*diffRed + diffGreen*diffGreen + diffBlue*diffBlue))
}

func GetCharacterImage(characterId string) image.Image {
	character, ok := Cosmetics.Items[characterId]
	if !ok {
		return getRandomCharacterImage()
	}
	
	response, err := http.Get(character.Images.Featured)
	if err != nil {
		return getRandomCharacterImage()
	}
	defer response.Body.Close()

	b, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	image, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		panic(err)
	}

	return image
}

func getRandomCharacterImage() image.Image {
	found := false
	var character FAPI_Cosmetic
	for !found {
		character = Cosmetics.GetRandomItemByType("AthenaCharacter")

		if character.Images.Featured == "" {
			continue
		}

		continueLoop := false
		for _, set := range SETS_NOT_ALLOWED {
			if strings.Contains(character.Set.Value, set) {
				continueLoop = true
				break
			}
		}

		if continueLoop {
			continue
		}

		if character.Introduction.BackendValue < 2 {
			continue
		}

		for _, tag := range character.GameplayTags {
			if strings.Contains(tag, "StarterPack") {
				continueLoop = true
				break
			}
		}
		if continueLoop {
			continue
		}

		found = true
	}

	response, err := http.Get(character.Images.Featured)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	b, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	image, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		panic(err)
	}

	return image
}

func getRandomCharacterImageWithSimilarColour(colour colours) image.Image {
	character := getRandomCharacterImage()
	characterColour := getAverageColour(character)

	if colorDifference(characterColour, colour) <= 140 {
		return character
	}

	return getRandomCharacterImageWithSimilarColour(colour)
}

func GenerateSoloImage() {
	background := *storage.Asset("background.png")

	itemFound := Cosmetics.GetRandomItemByType("AthenaCharacter")
	for itemFound.Images.Featured == "" {
		itemFound = Cosmetics.GetRandomItemByType("AthenaCharacter")
	}

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

func GenerateDuoImage() {
	background := *storage.Asset("background.png")

	player1 := getRandomCharacterImage()
	player2 := getRandomCharacterImageWithSimilarColour(getAverageColour(player1))
	
	bg, _, err := image.Decode(bytes.NewReader(background))
	if err != nil {
		panic(err)
	}


	m := image.NewRGBA(bg.Bounds())
	draw.Draw(m, m.Bounds(), bg, image.Point{0, 0}, draw.Src)

	resizedPlayerLeft := resize.Resize(0, uint(float64(m.Bounds().Dy()) * 1.25), player1, resize.Lanczos3)
	leftPosition := image.Point{
		m.Bounds().Dx()/2 - resizedPlayerLeft.Bounds().Dx()/2 - 200,
		(m.Bounds().Dy()/2 - resizedPlayerLeft.Bounds().Dy()/2) + 200,
	}
	draw.Draw(m, resizedPlayerLeft.Bounds().Add(leftPosition), resizedPlayerLeft, image.Point{0, 0}, draw.Over)

	resizedPlayerRight := resize.Resize(0, uint(float64(m.Bounds().Dy()) * 1.4), player2, resize.Lanczos3)
	rightPosition := image.Point{
		m.Bounds().Dx()/2 - resizedPlayerRight.Bounds().Dx()/2 + 240,
		(m.Bounds().Dy()/2 - resizedPlayerRight.Bounds().Dy()/2) + 200,
	}
	draw.Draw(m, resizedPlayerRight.Bounds().Add(rightPosition), resizedPlayerRight, image.Point{0, 0}, draw.Over)

	var bytes bytes.Buffer
	err = png.Encode(&bytes, m)
	if err != nil {
		panic(err)
	}

	PlaylistImages["Playlist_DefaultDuo"] = bytes.Bytes()
}

func GenerateTrioImage() {
	background := *storage.Asset("background.png")
	glow := *storage.Asset("glow.png")

	player1 := getRandomCharacterImage()
	player2 := getRandomCharacterImageWithSimilarColour(getAverageColour(player1))
	player3 := getRandomCharacterImageWithSimilarColour(getAverageColour(player1))

	bg, _, err := image.Decode(bytes.NewReader(background))
	if err != nil {
		panic(err)
	}

	glowImage, _, err := image.Decode(bytes.NewReader(glow))
	if err != nil {
		panic(err)
	}

	m := image.NewRGBA(bg.Bounds())
	draw.Draw(m, m.Bounds(), bg, image.Point{0, 0}, draw.Src)

	resizedPlayer1 := resize.Resize(0, uint(float64(m.Bounds().Dy()) * 1), player1, resize.Lanczos3)
	player1Position := image.Point{
		m.Bounds().Dx()/2 - resizedPlayer1.Bounds().Dx()/2 - 400,
		(m.Bounds().Dy()/2 - resizedPlayer1.Bounds().Dy()/2) + 150,
	}
	draw.Draw(m, resizedPlayer1.Bounds().Add(player1Position), resizedPlayer1, image.Point{0, 0}, draw.Over)

	resizedPlayer2 := resize.Resize(0, uint(float64(m.Bounds().Dy()) * 1.1), player2, resize.Lanczos3)
	player2Position := image.Point{
		m.Bounds().Dx()/2 - resizedPlayer2.Bounds().Dx()/2 + 350,
		(m.Bounds().Dy()/2 - resizedPlayer2.Bounds().Dy()/2) + 100,
	}
	draw.Draw(m, resizedPlayer2.Bounds().Add(player2Position), resizedPlayer2, image.Point{0, 0}, draw.Over)

	centre := image.Point{
		m.Bounds().Dx()/2 - glowImage.Bounds().Dx()/2,
		m.Bounds().Dy()/2 - glowImage.Bounds().Dy()/2,
	}
	draw.Draw(m, glowImage.Bounds().Add(centre), glowImage, image.Point{0, 0}, draw.Over)

	resizedPlayer3 := resize.Resize(0, uint(float64(m.Bounds().Dy()) * 1.4), player3, resize.Lanczos3)
	player3Position := image.Point{
		m.Bounds().Dx()/2 - resizedPlayer3.Bounds().Dx()/2,
		(m.Bounds().Dy()/2 - resizedPlayer3.Bounds().Dy()/2) + 200,
	}
	draw.Draw(m, resizedPlayer3.Bounds().Add(player3Position), resizedPlayer3, image.Point{0, 0}, draw.Over)

	var bytes bytes.Buffer
	err = png.Encode(&bytes, m)
	if err != nil {
		panic(err)
	}

	PlaylistImages["Playlist_DefaultTrio"] = bytes.Bytes()
}

func GenerateSquadImage() {
	background := *storage.Asset("background.png")
	glow := *storage.Asset("glow.png")

	player1 := getRandomCharacterImage()
	player2 := getRandomCharacterImageWithSimilarColour(getAverageColour(player1))
	player3 := getRandomCharacterImageWithSimilarColour(getAverageColour(player1))
	player4 := getRandomCharacterImageWithSimilarColour(getAverageColour(player1))

	bg, _, err := image.Decode(bytes.NewReader(background))
	if err != nil {
		panic(err)
	}

	glowImage, _, err := image.Decode(bytes.NewReader(glow))
	if err != nil {
		panic(err)
	}

	m := image.NewRGBA(bg.Bounds())
	draw.Draw(m, m.Bounds(), bg, image.Point{0, 0}, draw.Src)

	resizedPlayer4 := resize.Resize(0, uint(float64(m.Bounds().Dy()) * 0.8), player4, resize.Lanczos3)
	player4Position := image.Point{
		m.Bounds().Dx()/2 - resizedPlayer4.Bounds().Dx()/2 - 600,
		(m.Bounds().Dy()/2 - resizedPlayer4.Bounds().Dy()/2) + 110,
	}
	draw.Draw(m, resizedPlayer4.Bounds().Add(player4Position), resizedPlayer4, image.Point{0, 0}, draw.Over)

	resizedPlayer1 := resize.Resize(0, uint(float64(m.Bounds().Dy()) * 1), player1, resize.Lanczos3)
	player1Position := image.Point{
		m.Bounds().Dx()/2 - resizedPlayer1.Bounds().Dx()/2 - 350,
		(m.Bounds().Dy()/2 - resizedPlayer1.Bounds().Dy()/2) + 100,
	}
	draw.Draw(m, resizedPlayer1.Bounds().Add(player1Position), resizedPlayer1, image.Point{0, 0}, draw.Over)

	resizedPlayer2 := resize.Resize(0, uint(float64(m.Bounds().Dy()) * 1.1), player2, resize.Lanczos3)
	player2Position := image.Point{
		m.Bounds().Dx()/2 - resizedPlayer2.Bounds().Dx()/2 + 400,
		(m.Bounds().Dy()/2 - resizedPlayer2.Bounds().Dy()/2) + 100,
	}
	draw.Draw(m, resizedPlayer2.Bounds().Add(player2Position), resizedPlayer2, image.Point{0, 0}, draw.Over)

	centre := image.Point{
		m.Bounds().Dx()/2 - glowImage.Bounds().Dx()/2,
		m.Bounds().Dy()/2 - glowImage.Bounds().Dy()/2,
	}
	draw.Draw(m, glowImage.Bounds().Add(centre), glowImage, image.Point{0, 0}, draw.Over)

	resizedPlayer3 := resize.Resize(0, uint(float64(m.Bounds().Dy()) * 1.4), player3, resize.Lanczos3)
	player3Position := image.Point{
		m.Bounds().Dx()/2 - resizedPlayer3.Bounds().Dx()/2 + 150,
		(m.Bounds().Dy()/2 - resizedPlayer3.Bounds().Dy()/2) + 200,
	}
	draw.Draw(m, resizedPlayer3.Bounds().Add(player3Position), resizedPlayer3, image.Point{0, 0}, draw.Over)

	var bytes bytes.Buffer
	err = png.Encode(&bytes, m)
	if err != nil {
		panic(err)
	}

	PlaylistImages["Playlist_DefaultSquad"] = bytes.Bytes()
}

func GeneratePlaylistImages() {
	GenerateSoloImage()
	GenerateDuoImage()
	GenerateTrioImage()
	GenerateSquadImage()

	aid.Print("(snow) generated playlist images")
}
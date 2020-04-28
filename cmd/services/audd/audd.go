package audd

import (
	"fmt"
	"log"
	"net/http"
	"os"

	auddgo "github.com/AudDMusic/audd-go"

	"github.com/AnthonyNixon/raymond-api/cmd/utils/httperr"
)

var AUDD_TOKEN string

func Initialize() {
	AUDD_TOKEN = os.Getenv("RAYMOND_AUDD_TOKEN")
	if AUDD_TOKEN == "" {
		log.Fatal("No AUDD token present.")
	}
}

func IdentifyFile() (title string, artist string, error httperr.HttpErr) {
	fmt.Printf("token: %s\n", AUDD_TOKEN)
	parameters := map[string]string{
		"return": "timecode,apple_music,deezer,spotify",
	}
	file, err := os.Open("song.mp3")
	if err != nil {
		return "", "", httperr.New(http.StatusInternalServerError, "Failed to open song file", err.Error())
	}

	result, err := auddgo.RecognizeByFile(file, AUDD_TOKEN, parameters)
	if err != nil {
		return "", "", httperr.New(http.StatusInternalServerError, "Failed to Recognize File with AUDD", err.Error())
	}
	file.Close()

	song := result.Result
	fmt.Printf("%s - %s.\nTimecode: %s, album: %s. ℗ %s, %s\n\n Listen on: Apple Music: %s", song.Artist, song.Title, song.Timecode, song.Album, song.Label, song.ReleaseDate, song.AppleMusic.URL)
	return song.Title, song.Artist, nil
}

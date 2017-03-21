package main

import (
	"gobot.io/x/gobot/platforms/raspi"
	"os"
	"encoding/csv"
	"time"
	"github.com/stacktic/dropbox"
	_ "github.com/joho/godotenv/autoload"
	"fmt"
	"path/filepath"
)

func main() {
	pi := raspi.NewAdaptor()

	db := newDropbox(os.Getenv("DROPBOX_CLIENT_ID"), os.Getenv("DROPBOX_ACCESS_TOKEN"))

	folder := os.Getenv("TMP_FOLDER")
	dropboxFolder := os.Getenv("DROPBOX_FOLDER")

	current := currentFilename(folder)
	writer, file := newWriter(folder, current)

	makeSure(pi.Connect())

	fmt.Printf("Starting up - %s\n", current)

	go uploadPending(folder, current)

	for {
		filename := currentFilename(folder)
		if filename != current {
			writer.Flush()
			file.Close()
			go upload(current, dropboxFolder, db)
			current = filename
			writer, file = newWriter(folder, current)
		}

		writer.Write(dataFor(capture(pi)))

		time.Sleep(10 * time.Second)
	}
}
func upload(filename string, dropboxFolder string, db *dropbox.Dropbox) {
	_, name := filepath.Split(filename)
	dbFilename := filepath.Join(dropboxFolder, name)
	fmt.Printf("Uploading - %s -> %s\n", filename, dbFilename)

	_, err := db.UploadFile(filename, dbFilename, false, "")
	if err != nil {
		// TODO retry?
		fmt.Printf("Couldn't upload - %s\n", err.Error())
	}
	os.Remove(filename)
}

func uploadPending(folder string, ignoreFile string) {
	// TODO list pending
	// TODO upload
	// TODO delete

}

func capture(adaptor *raspi.Adaptor) Readings {
	humidity, _ := adaptor.DigitalRead("1")
	temp, _ := adaptor.DigitalRead("2")

	return Readings{
		Humidity:    humidity,
		Temperature: temp,
		Time:        time.Now(),
	}
}

func newDropbox(clientId string, accessToken string) *dropbox.Dropbox {
	db := dropbox.NewDropbox()
	db.SetAppInfo(clientId, "")
	db.SetAccessToken(accessToken)

	return db
}

func newWriter(folder string, filename string) (*csv.Writer, *os.File) {
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		makeSure(os.MkdirAll(folder, 0700))
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			panic(err)
		}
		writer := csv.NewWriter(file)
		writer.Write(headers())
		writer.Flush()

		return writer, file
	} else {
		file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND, 0660)
		if err != nil {
			panic(err)
		}
		writer := csv.NewWriter(file)
		return writer, file
	}
}

func makeSure(res error) {
	if res != nil {
		panic(res)
	}
}

func headers() []string {
	return []string{"humidity", "temperature", "ambient_light", "noise_level", "movement", "time" }
}

func dataFor(r Readings) []string {
	return []string{
		string(r.Humidity),
		string(r.Temperature),
		string(r.AmbientLight),
		string(r.NoiseLevel),
		string(r.Movement),
		r.Time.String(),
	}
}

// one file per minute
func currentFilename(folder string) string {
	return filepath.Join(
		folder,
		time.Now().Format("2006_01_02_15_04.csv"),
	)
}

type Readings struct {
	Humidity     int
	Temperature  int
	AmbientLight int
	NoiseLevel   int
	Movement     int
	Time         time.Time
}

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

	current := currentFilename(folder)
	writer, file := newWriter(folder, current)

	makeSure(pi.Connect())

	fmt.Printf("Starting up - %s\n", current)

	for {
		// TODO get readings
		//writer.Write()

		filename := currentFilename(folder)
		if filename != current {
			writer.Flush()
			file.Close()
			go func(filename string) {
				fmt.Printf("Uploading - %s\n", current)
				_, err := db.UploadFile(filename, filename, false, "")
				if err != nil {
					// TODO retry?
					fmt.Printf("Couldn't upload - %s\n", err.Error())
				}
				os.Remove(filename)
			}(filename)
			current = filename
			writer, file = newWriter(folder, current)
		}
		time.Sleep(10 * time.Second)
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
	return []string{"humidity", "time" }
}

// one file per minute
func currentFilename(folder string) string {
	return filepath.Join(
		folder,
		time.Now().Format("2006_01_02_15_04.csv"),
	)
}

//type Readings struct {
//	Humidity     int
//	Temperature  int
//	AmbientLight int
//	NoiseLevel   int
//	Movement     int
//	Time         int64
//}

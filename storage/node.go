package storage

import (
	"gosync/config"
	"gosync/prototypes"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func getNodePath(listener, local_path string) string {
	return strings.TrimPrefix(local_path, getBaseDir(listener))
}

func GetNodeCopy(item prototypes.DataTable, listener string) bool {
	cfg := config.GetConfig()
	lConf := cfg.Listeners[listener]
	log.Println("Downloading file...")

	rawURL := "http://" + item.HostUpdated + ":" + cfg.ServerConfig.ListenPort + "/" + listener + getNodePath(listener, item.Path)

	_, err := url.Parse(rawURL)

	if err != nil {
		panic(err)
	}

	check := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	resp, err := check.Get(rawURL) // add a filter to check redirect

	if err != nil {
		log.Println(err)
		return false

	}
	defer resp.Body.Close()
	log.Println(resp.Status)

	if resp.Status == "404" {
		log.Fatalf("File not found %s", rawURL)
	}

	if _, err := os.Stat(item.Path); err == nil {
		// We wipe the file as we need to replace with a new one
		err := os.Remove(item.Path)
		if err != nil {
			log.Fatalf("Could not remove file %s : %+v", item.Path, err.Error())
		}

	}

	file, err := os.Create(item.Path)

	if err != nil {
		log.Println(err)
		panic(err)
	}
	defer file.Close()

	size, err := io.Copy(file, resp.Body)

	if err != nil {
		panic(err)
	}

	file.Chown(lConf.Uid, lConf.Gid)
	//file.Chmod(os.FileMode(item.Perms))

	log.Printf("%s with %v bytes downloaded", item.Path, size)
	return true
}

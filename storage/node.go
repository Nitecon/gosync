package storage

import (
	"gosync/config"
	"gosync/prototypes"
	"log"
	"net/http"
	"net/url"
    "gosync/utils"
)

func GetNodeCopy(item prototypes.DataTable, listener string, uid, gid int, perms string) bool {
	cfg := config.GetConfig()
	log.Println("Downloading file...")

	rawURL := "http://" + item.HostUpdated + ":" + cfg.ServerConfig.ListenPort + "/" + listener + utils.GetRelativePath(listener, item.Path)

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

    size, fserr := utils.FileWrite(item.Path, resp.Body, uid, gid, perms)
    if fserr != nil {
        log.Fatalf("Error occurred writing file (%s): %+v", fserr.Error(), fserr)
    }


	log.Printf("%s with %v bytes downloaded", item.Path, size)
	return true
}

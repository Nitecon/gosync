package storage

import (
	"gosync/prototypes"
	"net/http"
	"net/url"
    "gosync/utils"
    "runtime"
)

func GetNodeCopy(item prototypes.DataTable, listener string, uid, gid int, perms string) bool {
	cfg := utils.GetConfig()
    utils.WriteLn("Downloading file...")

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

	utils.Check("Unable to download from: "+rawURL, 404, err)
	runtime.Goexit()
    defer resp.Body.Close()
    utils.WriteLn(resp.Status)

	if resp.Status == "404" {
        utils.WriteF("File not found %s", rawURL)
	}

    size, fserr := utils.FileWrite(item.Path, resp.Body, uid, gid, perms)
    utils.Check("Cannot write file...", 500, fserr)


    utils.WriteF("%s with %v bytes downloaded", item.Path, size)
	return true
}

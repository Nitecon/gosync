package storage

import (
	"net/http"
	"net/url"
    "gosync/utils"
    "gosync/nodeinfo"
)

func GetNodeCopy(item utils.DataTable, listener string, uid, gid int, perms string) bool {
	cfg := utils.GetConfig()
    aliveNodes, aliveErr := nodeinfo.GetNodes()
    utils.LogWriteF("Nodes: %+v \n\nAlive Errors: %s\n\n\n", aliveNodes, aliveErr.Error())
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
	if !utils.ErrorCheckF(err, 404, "Do not allow redirects: %s ", rawURL){
        defer resp.Body.Close()
        if resp.Status == "404" {
            utils.LogWriteF("File not found: %s", rawURL)
        }
        size, err := utils.FileWrite(item.Path, resp.Body, uid, gid, perms)
        utils.ErrorCheckF(err, 500, "Cannot write file: %s", item.Path)
        utils.LogWriteF("%s with %v bytes downloaded", item.Path, size)
        return true
    }
    return false
}

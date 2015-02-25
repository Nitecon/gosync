package storage

type GDrive struct {
	config string
}

func (g *GDrive) Upload(local_path, remote_path string) error {
	fmt.Println("Doing an upload to GDrive")
}

func (g *GDrive) Download(remote_path, local_path string) error {
	fmt.Println("Doing a download from GDrive")
}

func (g *GDrive) CheckMD5(local_path, remote_path string) bool {
	fmt.Println("Doing md5 check on GDrive")
	return true
}

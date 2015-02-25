package storage

type S3 struct {
	config string
}

func (g *S3) Upload(local_path, remote_path string) error {
	fmt.Println("Doing an upload to S3")
	return nil
}

func (g *S3) Download(remote_path, local_path string) error {
	fmt.Println("Doing a download from S3")
	return nil
}

func (g *S3) CheckMD5(local_path, remote_path string) bool {
	fmt.Println("Doing md5 check on S3")
	return true
}

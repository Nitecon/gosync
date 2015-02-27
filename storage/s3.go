package storage

import "log"

type S3 struct {
	config string
}

func (g *S3) Upload(local_path, remote_path string) error {
	log.Printf("S3 Uploading %s -> %s", local_path, remote_path)
	return nil
}

func (g *S3) Download(remote_path, local_path string) error {
	log.Printf("S3 Downloading %s -> %s", remote_path, local_path)
	return nil
}

func (g *S3) CheckMD5(local_path, remote_path string) bool {
	log.Printf("S3 MD5 Check %s -> %s", local_path, remote_path)
	return true
}

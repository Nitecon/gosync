package storage

import (
	"github.com/rlmcpherson/s3gof3r"
	"gosync/config"
	"log"
	"net/http"
	"os"
    "io"
)

type S3 struct {
	config   *config.Configuration
	bucket   *s3gof3r.Bucket
	listener string
}

type Keys struct {
    AccessKey     string
    SecretKey     string
    SecurityToken string
}

func (s *S3) Upload(local_path, remote_path string) error {
    log.Printf("S3 Uploading %s -> %s", local_path, remote_path)
	conf := new(s3gof3r.Config)
	*conf = *s3gof3r.DefaultConfig
    keys := new(s3gof3r.Keys)
    keys.AccessKey = s.config.S3Config.Key
    keys.SecretKey = s.config.S3Config.Secret

	s3 := s3gof3r.New("s3.amazonaws.com", *keys)
	b := s3.Bucket(s.config.Listeners[s.listener].Bucket)
	conf.Concurrency = 10
	// Setting debug to true (last var below)
	s3gof3r.SetLogger(os.Stderr, "", log.LstdFlags, true)
	r, err := os.Open(local_path)
	if err != nil {
        log.Fatalf("Error occurred opening local file (%s): %+v", err.Error(), err)
	}
	defer r.Close()
    header := make(http.Header)

	w, err := b.PutWriter(remote_path, header, conf)
	if err != nil {
        log.Fatalf("Error getting put writer (%s): %+v", err.Error(), err)
	}
	if _, err = io.Copy(w, r); err != nil {
        log.Fatalf("Error occurred with io.Copy (%s): %+v", err.Error(), err)
	}
	if err = w.Close(); err != nil {
        log.Fatalf("Error occurred closing s3 writer (%s): %+v", err.Error(), err)
	}
	log.Printf("S3 Uploading %s -> %s", local_path, remote_path)
	return nil
}

func (s *S3) Download(remote_path, local_path string) error {
	log.Printf("S3 Downloading %s -> %s", remote_path, local_path)
	return nil
}

func (s *S3) CheckMD5(local_path, remote_path string) bool {
	log.Printf("S3 MD5 Check %s -> %s", local_path, remote_path)
	return true
}
func getListener(dir string) string {
	cfg := config.GetConfig()
	var listener = ""
	for lname, ldata := range cfg.Listeners {
		if ldata.Directory == dir {
			listener = lname
		}
	}
	return listener
}

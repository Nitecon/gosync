package storage

import (
	"github.com/rlmcpherson/s3gof3r"
	"gosync/config"
    "gosync/utils"
	"log"
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
    conf,keys := s.GetS3Config()

    // Open bucket to put file into
    s3 := s3gof3r.New("", *keys)
    b := s3.Bucket(s.config.Listeners[s.listener].Bucket)

    // open file to upload
    file, err := os.Open(local_path)
    if err != nil {
        return err
    }

    // Open a PutWriter for upload
    w, err := b.PutWriter(remote_path, nil, conf)
    if err != nil {
        return err
    }
    if _, err = io.Copy(w, file); err != nil { // Copy into S3
        return err
    }
    if err = w.Close(); err != nil {
        return err
    }
    return nil

}

func (s *S3) Download(remote_path, local_path string) error {
    log.Printf("S3 Downloading %s -> %s", remote_path, utils.GetRelativePath(s.listener, local_path))
    conf,keys := s.GetS3Config()

    // Open bucket to put file into
    s3 := s3gof3r.New("", *keys)
    b := s3.Bucket(s.config.Listeners[s.listener].Bucket)

    r, h, err := b.GetReader(remote_path, conf)
    if err != nil {
        return err
    }
    // stream to standard output
    if _, err = io.Copy(os.Stdout, r); err != nil {
        return err
    }
    err = r.Close()
    if err != nil {
        return err
    }
    log.Println(h) // print key header data

    return nil
}

func (s *S3) CheckMD5(local_path, remote_path string) bool {
	log.Printf("S3 MD5 Check %s -> %s", local_path, remote_path)
	return true
}

func (s *S3) GetS3Config() (*s3gof3r.Config, *s3gof3r.Keys){
    conf := new(s3gof3r.Config)
    *conf = *s3gof3r.DefaultConfig
    keys := new(s3gof3r.Keys)
    keys.AccessKey = s.config.S3Config.Key
    keys.SecretKey = s.config.S3Config.Secret
    conf.Concurrency = 10
    // Setting debug to true (last var below)
    s3gof3r.SetLogger(os.Stderr, "", log.LstdFlags, true)
    return conf, keys
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

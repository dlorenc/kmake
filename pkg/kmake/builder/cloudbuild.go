package builder

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	cstorage "cloud.google.com/go/storage"
	"golang.org/x/oauth2/google"
	cloudbuild "google.golang.org/api/cloudbuild/v1"
	storage "google.golang.org/api/storage/v1"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func RemoteBuild(imageName, dockerfile, projectId string, t Tagger) (string, error) {
	tag := fmt.Sprintf("gcr.io/%s/%s", projectId, imageName)

	//TODO(@r2d4): make context configurable
	dir := filepath.Dir(dockerfile)

	cbBucket := projectId + "_cloudbuild"
	buildObject := fmt.Sprintf("source/%s-%s.tar.gz", imageName, randomID())

	ctx := context.Background()
	client, err := google.DefaultClient(ctx, cloudbuild.CloudPlatformScope)
	if err != nil {
		return "", errors.Wrap(err, "getting client")
	}

	hc, err := google.DefaultClient(ctx, storage.CloudPlatformScope)
	if err != nil {
		return "", errors.Wrap(err, "getting authenticated http client")
	}

	logrus.Infof("Pushing code to gs://%s/%s", cbBucket, buildObject)

	if err := uploadTar(ctx, dir, hc, cbBucket, buildObject); err != nil {
		return "", errors.Wrap(err, "uploading source tarball")
	}

	cb, err := cloudbuild.New(client)
	if err != nil {
		return "", errors.Wrap(err, "getting builder")
	}

	var steps []*cloudbuild.BuildStep
	steps = append(steps, &cloudbuild.BuildStep{
		Name: "gcr.io/cloud-builders/docker",
		Args: []string{"build", "--tag", tag, "."},
	})
	call := cb.Projects.Builds.Create(projectId, &cloudbuild.Build{
		LogsBucket: cbBucket,
		Source: &cloudbuild.Source{
			StorageSource: &cloudbuild.StorageSource{
				Bucket: cbBucket,
				Object: buildObject,
			},
		},
		Steps:  steps,
		Images: []string{tag},
	})
	op, err := call.Context(ctx).Do()
	if err != nil {
		return "", errors.Wrap(err, "could not create build")
	}

	remoteID, err := getBuildID(op)
	if err != nil {
		return "", errors.Wrapf(err, "getting build ID from op")
	}

	logrus.Infof("Logs at https://console.cloud.google.com/m/cloudstorage/b/%s/o/log-%s.txt", cbBucket, remoteID)

	fail := false
	var digest string
	for {
		b, err := cb.Projects.Builds.Get(projectId, remoteID).Do()
		if err != nil {
			return "", errors.Wrap(err, "getting build status")
		}

		if s := b.Status; s != "WORKING" && s != "QUEUED" {
			if b.Status == "FAILURE" {
				fail = true
			}
			logrus.Infof("Build status: %v", s)
			digest, err = getImageID(b)
			if err != nil {
				return "", errors.Wrap(err, "getting image id from finished build")
			}
			break
		}

		time.Sleep(time.Second)
	}

	c, err := cstorage.NewClient(ctx)
	if err != nil {
		return "", errors.Wrap(err, "getting cloud storage client")
	}
	defer c.Close()
	if err := c.Bucket(cbBucket).Object(buildObject).Delete(ctx); err != nil {
		return "", errors.Wrap(err, "cleaning up source tar after build")
	}
	logrus.Print("Deleted source targz")
	if fail {
		return "", errors.Wrap(err, "cloud build failed")
	}

	if err != nil {
		return "", errors.Wrap(err, "getting result image digest")
	}
	logrus.Infof("Image built at %s", fmt.Sprintf("%s@%s", imageName, digest))

	return fmt.Sprintf("%s@%s", tag, digest), nil
}

func getBuildID(op *cloudbuild.Operation) (string, error) {
	if op.Metadata == nil {
		return "", errors.New("missing Metadata in operation")
	}
	var buildMeta cloudbuild.BuildOperationMetadata
	if err := json.Unmarshal([]byte(op.Metadata), &buildMeta); err != nil {
		return "", err
	}
	if buildMeta.Build == nil {
		return "", errors.New("missing Build in operation metadata")
	}
	return buildMeta.Build.Id, nil
}

func getImageID(b *cloudbuild.Build) (string, error) {
	if b.Results == nil || len(b.Results.Images) == 0 {
		spew.Dump(b)
		return "", errors.New("missing build result image metadata")
	}
	return b.Results.Images[0].Digest, nil
}

// https://github.com/upspin/gcp/blob/a0be1026ff5e1367fe35314b9a10230c00d5d493/cmd/upspin-deploy-gcp/cdbuild.go#L143
func uploadTar(ctx context.Context, root string, hc *http.Client, bucket string, objectName string) error {
	c, err := cstorage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer c.Close()

	w := c.Bucket(bucket).Object(objectName).NewWriter(ctx)
	gzw := gzip.NewWriter(w)
	tw := tar.NewWriter(gzw)

	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if path == root {
			return nil
		}
		relpath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		info = renamingFileInfo{info, relpath}

		hdr, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		hdr.Name = rel
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(tw, f)
		return err
	}); err != nil {
		w.CloseWithError(err)
		return err
	}
	if err := tw.Close(); err != nil {
		w.CloseWithError(err)
		return err
	}
	if err := gzw.Close(); err != nil {
		w.CloseWithError(err)
		return err
	}
	return w.Close()
}

type renamingFileInfo struct {
	os.FileInfo
	name string
}

func (fi renamingFileInfo) Name() string {
	return fi.name
}

func randomID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", b)
}

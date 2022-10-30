package docker_image

import (
	"archive/tar"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/distribution/manifest/schema2"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/google/go-containerregistry/pkg/crane"
	log "github.com/sirupsen/logrus"
)

type DockerImageManager struct {
	image      string
	docker     *client.Client
	dockerAuth types.AuthConfig
}

func NewDockerImageManager(docker *client.Client, dockerAuth types.AuthConfig, image string) *DockerImageManager {
	return &DockerImageManager{
		docker:     docker,
		dockerAuth: dockerAuth,
		image:      image,
	}
}

func (manager *DockerImageManager) pullImage() {
	js, _ := json.Marshal(manager.dockerAuth)
	authstr := base64.URLEncoding.EncodeToString(js)

	log.WithField("image", manager.image).Info("Pulling image")

	rc, err := manager.docker.ImagePull(context.Background(), manager.image, types.ImagePullOptions{RegistryAuth: string(authstr)})
	if err != nil {
		panic(err)
	}

	io.Copy(ioutil.Discard, rc)
	rc.Close()
	log.WithField("image", manager.image).Info("Pulled image")
}

func (manager *DockerImageManager) GetImageID() (string, error) {
	bytes, err := crane.Manifest(manager.image)
	if err != nil {
		return "", err
	}

	var manifest schema2.DeserializedManifest
	err = manifest.UnmarshalJSON(bytes)
	if err != nil {
		return "", err
	}

	return strings.TrimPrefix(manifest.Manifest.Config.Digest.String(), "sha256:"), nil
}

func (manager *DockerImageManager) extractImageContents(dir string) {
	manager.pullImage()
	resp, err := manager.docker.ContainerCreate(context.Background(), &container.Config{Image: manager.image}, nil, nil, nil, "")
	if err != nil {
		panic(err)
	}
	defer manager.docker.ContainerRemove(context.Background(), resp.ID, types.ContainerRemoveOptions{})

	log.WithField("image", manager.image).WithField("container", resp.ID).Info("Created temp container")

	tarStream, err := manager.docker.ContainerExport(context.Background(), resp.ID)
	if err != nil {
		panic(err)
	}

	untar(dir, tarStream)
	tarStream.Close()
}

func (manager *DockerImageManager) GetImageContents() string {
	dir, err := ioutil.TempDir("", "docker-packages")
	if err != nil {
		panic(err)
	}
	log.WithField("image", manager.image).WithField("dir", dir).Info("Created temp dir")

	manager.extractImageContents(dir)
	return dir
}

// Untar takes a destination path and a reader; a tar reader loops over the tarfile
// creating the file structure at 'dst' along the way, and writing any files
func untar(dst string, r io.Reader) error {
	tr := tar.NewReader(r)

	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return nil

		// return any other error
		case err != nil:
			log.WithError(err).Error("Untar error")
			return err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(dst, header.Name)

		// the following switch could also be done using fi.Mode(), not sure if there
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}

			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			f.Close()

		case tar.TypeSymlink:
			err := os.Symlink(header.Linkname, target)
			if err != nil {
				return err
			}
		}
	}
}

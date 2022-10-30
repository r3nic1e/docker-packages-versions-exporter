# packages-versions-exporter

## Plan
* Get running images from Prometheus metrics
* Check if images have changed
* Load changed images FS and chroot into it
* Detect OS
* Get packages list depending on OS
* Load packages list to media

## Media filetree

```
docker-packages/
  { image name without version}/
    { image version }/
      { image id } - file with packages list
```

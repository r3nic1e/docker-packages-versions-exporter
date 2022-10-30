package analyze

import (
	"os"
	"path"
	"syscall"
)

func chroot(dir string) (func() error, error) {
	root, err := os.Open("/")
	if err != nil {
		return nil, err
	}

	err = mkchr(path.Join(dir, "dev/null"), mkdev(1, 3))
	if err != nil {
		return nil, err
	}
	err = mkchr(path.Join(dir, "dev/urandom"), mkdev(1, 9))
	if err != nil {
		return nil, err
	}

	if err := syscall.Chroot(dir); err != nil {
		root.Close()
		return nil, err
	}

	return func() error {
		defer root.Close()
		if err := root.Chdir(); err != nil {
			return err
		}

		if err := syscall.Chroot("."); err != nil {
			return err
		}

		return nil
	}, nil
}

func mkchr(dir string, device int) error {
	err := syscall.Mknod(dir, syscall.S_IFCHR|0666, device)
	if err != nil {
		return err
	}
	err = syscall.Chmod(dir, 0666)
	if err != nil {
		return err
	}
	return nil
}

// mkdev is used to build the value of linux devices (in /dev/) which specifies major
// and minor number of the newly created device special file.
// Linux device nodes are a bit weird due to backwards compat with 16 bit device nodes.
// They are, from low to high: the lower 8 bits of the minor, then 12 bits of the major,
// then the top 12 bits of the minor.
func mkdev(major int, minor int) int {
	dev := (uint64(major) & 0x00000fff) << 8
	dev |= (uint64(major) & 0xfffff000) << 32
	dev |= (uint64(minor) & 0x000000ff) << 0
	dev |= (uint64(minor) & 0xffffff00) << 12
	return int(dev)
}

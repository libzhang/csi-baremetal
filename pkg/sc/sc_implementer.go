package sc

// StorageClassImplementer is an interface with different methods for a volume creation depending on Storage Class
// It should be used across node level operations
type StorageClassImplementer interface {
	CreateFileSystem(fsType FileSystem, device string) error
	DeleteFileSystem(device string) error

	CreateTargetPath(path string) error
	DeleteTargetPath(path string) error

	IsMounted(device string) (bool, error)
	BindMount(device, dir string, mountDevice bool) error
	Unmount(path string) error

	IsMountPoint(path string) (bool, error)
	// atomic methods for using in NodePublish
	PrepareVolume(device, targetPath string) (bool, error)
}

// FileSystem defines Linux filesystem
type FileSystem string

// Filesystem which can be used for CSI
const (
	XFS FileSystem = "xfs"
)

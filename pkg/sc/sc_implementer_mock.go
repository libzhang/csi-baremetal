package sc

import (
	"github.com/stretchr/testify/mock"
)

type ImplementerMock struct {
	mock.Mock
}

func (m *ImplementerMock) CreateFileSystem(fsType FileSystem, device string) error {
	args := m.Mock.Called(fsType, device)
	return args.Error(0)
}

func (m *ImplementerMock) DeleteFileSystem(device string) error {
	args := m.Mock.Called(device)
	return args.Error(0)
}

func (m *ImplementerMock) CreateTargetPath(path string) error {
	args := m.Mock.Called(path)
	return args.Error(0)
}

func (m *ImplementerMock) DeleteTargetPath(path string) error {
	args := m.Mock.Called(path)
	return args.Error(0)
}

func (m *ImplementerMock) IsMounted(device string) (bool, error) {
	args := m.Mock.Called(device)
	return args.Bool(0), args.Error(1)
}

func (m *ImplementerMock) BindMount(device, dir string, mountDevice bool) error {
	args := m.Mock.Called(device, dir, mountDevice)
	return args.Error(0)
}

func (m *ImplementerMock) Unmount(path string) error {
	args := m.Mock.Called(path)
	return args.Error(0)
}

func (m *ImplementerMock) IsMountPoint(path string) (bool, error) {
	args := m.Mock.Called(path)
	return args.Bool(0), args.Error(1)
}

func (m *ImplementerMock) PrepareVolume(device, targetPath string) (bool, error) {
	args := m.Mock.Called(device, targetPath)
	return args.Bool(0), args.Error(1)
}

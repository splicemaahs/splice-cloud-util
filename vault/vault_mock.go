package vault

import (
	"encoding/json"

	"github.com/stretchr/testify/mock"
)

type vaultMock struct {
	mock.Mock
}

// NewVaultMock - Mocking the vault interactions
func NewVaultMock() Client {

	return &vaultMock{}

}

func (vm *vaultMock) GetVersion() (string, error) {

	return "1.2.3", nil
}

func (vm *vaultMock) GetData(keypath string) (DataRecord, error) {

	var dataRecord DataRecord
	var newData string

	switch keypath {
	case "/test/1":
		newData = `{
			"data": {
				"data": {
					"value": "vault-test1"
				},
				"metadata": {
					"created_time": "2019-08-30T02:43:30.607941986Z",
					"deletion_time": "",
					"destroyed": false,
					"version": 1
				}
			}
		}`
	case "/test/2":
		newData = `{
			"data": {
				"data": {
					"value": "vault-test2"
				},
				"metadata": {
					"created_time": "2019-08-30T02:43:30.607941986Z",
					"deletion_time": "",
					"destroyed": false,
					"version": 1
				}
			}
		}`
	default:
		newData = `{}`
	}

	marshErr := json.Unmarshal([]byte(newData), &dataRecord)
	if marshErr != nil {
		return DataRecord{}, marshErr
	}

	return dataRecord, nil

}
func (vm *vaultMock) GetPaths(keypath string) (map[string]Paths, error) {

	pathList := make(map[string]Paths, 0)

	switch keypath {
	case "/test":
		pathList["/test/1"] = Paths{
			Type:     vaultData,
			Path:     "1",
			Parent:   "/test",
			FullPath: "/test/1",
		}
		pathList["/test/2"] = Paths{
			Type:     vaultData,
			Path:     "2",
			Parent:   "/test",
			FullPath: "/test/2",
		}
		pathList["/test/folder1"] = Paths{
			Type:     vaultFolder,
			Path:     "folder1",
			Parent:   "/test",
			FullPath: "/test/folder1",
		}
		pathList["/test/folder2"] = Paths{
			Type:     vaultFolder,
			Path:     "folder2",
			Parent:   "/test",
			FullPath: "/test/folder2",
		}
	case "/test/folder1":
		pathList["/test/folder1/data"] = Paths{
			Type:     vaultData,
			Path:     "data",
			Parent:   "/test/folder1",
			FullPath: "/test/folder1/data",
		}
	case "/test/folder2":
		pathList["/test/folder2/data"] = Paths{
			Type:     vaultData,
			Path:     "data",
			Parent:   "/test/folder2",
			FullPath: "/test/folder2/data",
		}
	}
	return pathList, nil
}

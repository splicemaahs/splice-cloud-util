package vault

import (
	"encoding/json"
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestNewVault(t *testing.T) {

	mockVault := NewVaultMock()

	vaultVersion, err := mockVault.GetVersion()
	if err != nil {
		assert.Equal(t, err, nil)
	}
	assert.Equal(t, vaultVersion, "1.2.3")

}

func TestGetVersion(t *testing.T) {

	mockVault := NewVaultMock()

	vaultVersion, err := mockVault.GetVersion()
	if err != nil {
		assert.Equal(t, err, nil)
	}
	assert.Equal(t, vaultVersion, "1.2.3")

}

func TestGetData(t *testing.T) {
	var dataRecord DataRecord

	mockVault := NewVaultMock()

	newData := `{
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

	marshErr := json.Unmarshal([]byte(newData), &dataRecord)
	if marshErr != nil {
		// Do we just fail the test here?
		assert.Equal(t, marshErr, nil)
	}

	vaultData, err := mockVault.GetData("/test/1")
	if err != nil {
		assert.Equal(t, marshErr, nil)
	}

	assert.Equal(t, dataRecord, vaultData)

	emptyData, err := mockVault.GetData("/test/99")
	if err != nil {
		assert.Equal(t, marshErr, nil)
	}
	assert.Equal(t, DataRecord{}, emptyData)

}

func TestGetPaths(t *testing.T) {

	newPaths := make(map[string]Paths, 0)

	mockVault := NewVaultMock()

	newPaths["/test/1"] = Paths{
		Type:     "Folder",
		Path:     "1",
		Parent:   "/test",
		FullPath: "/test/1",
	}
	newPaths["/test/2"] = Paths{
		Type:     "Folder",
		Path:     "2",
		Parent:   "/test",
		FullPath: "/test/2",
	}

	vaultPaths, err := mockVault.GetPaths("/test")
	if err != nil {
		assert.Equal(t, err, nil)
	}

	assert.Equal(t, newPaths, vaultPaths)

}

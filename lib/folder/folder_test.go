package folder

import "testing"

func TestListFolders(t *testing.T) {
	folders, err := ListFolders("./folder_test")
	if err != nil {
		t.Error(err)
	}

	expected := map[string]bool{"folder1": true, "folder2": true}

	for _, folder := range folders {
		if _, ok := expected[folder]; !ok {
			t.Errorf("folder %s doesn't exists", folder)
		}
	}
}

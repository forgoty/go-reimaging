package validators

import (
	"testing"
)

func TestValidateDownloadDirWhenUnvalidPathIsProvided(t *testing.T) {
	_, error := ValidateDownloadDir("$^%*@#")
	if error == nil {
		t.Errorf("Error should be raised")
	}
}

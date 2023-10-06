package shred

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestShred(t *testing.T) {
	// Define test cases
	tests := []struct {
		name      string
		setupFunc func() (string, error)
		expectErr bool
	}{
		{
			name: "shred existing file",
			setupFunc: func() (string, error) {
				file, err := ioutil.TempFile("", "shredder_test")
				if err != nil {
					return "", err
				}
				defer file.Close()
				return file.Name(), nil
			},
			expectErr: false,
		},
		{
			name: "shred non-existing file",
			setupFunc: func() (string, error) {
				return "nonexistentfile", nil
			},
			expectErr: true,
		},
		{
			name: "shred read-only file",
			setupFunc: func() (string, error) {
				file, err := ioutil.TempFile("", "shredder_test")
				if err != nil {
					return "", err
				}
				file.Close()
				os.Chmod(file.Name(), 0444)
				return file.Name(), nil
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath, err := tt.setupFunc()
			if err != nil {
				t.Fatalf("Failed to set up %q test: %v", tt.name, err)
			}

			err = Shred(filePath)

			if !tt.expectErr && err != nil {
				t.Errorf("Shred() error = %v, expectErr %v", err, tt.expectErr)
				return
			}

			if tt.expectErr && err == nil {
				t.Errorf("Expected error, but got none")
				return
			}
		})
	}
}

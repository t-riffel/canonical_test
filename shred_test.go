package shred

import (
	"io/ioutil"
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
				// Create a temp file
				file, err := ioutil.TempFile("", "shredder_test")
				if err != nil {
					return "", err
				}
				defer file.Close()
				return file.Name(), nil
			},
			expectErr: false,
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

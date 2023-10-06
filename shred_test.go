package shred

import (
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
				file, err := os.CreateTemp("", "shredder_test")
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
				file, err := os.CreateTemp("", "shredder_test")
				if err != nil {
					return "", err
				}
				file.Close()
				os.Chmod(file.Name(), 0444)
				return file.Name(), nil
			},
			expectErr: true,
		},
		{
			name: "shred a directory",
			setupFunc: func() (string, error) {
				dir, err := os.MkdirTemp("", "shredder_test")
				if err != nil {
					return "", err
				}
				return dir, nil
			},
			expectErr: true,
		},
		{
			name: "shred a symlink",
			setupFunc: func() (string, error) {
				file, err := os.CreateTemp("", "shredder_test")
				if err != nil {
					return "", err
				}
				file.Close()
				symlink := file.Name() + "_symlink"
				// Create a symlink to the temporary file
				err = os.Symlink(file.Name(), symlink)
				if err != nil {
					return "", err
				}
				return symlink, nil
			},
			expectErr: false,
		},
		{
			name: "shred a device file",
			setupFunc: func() (string, error) {
				return "/dev/null", nil // use /dev/null as an example device file
			},
			expectErr: true,
		},
		{
			name: "shred an empty file",
			setupFunc: func() (string, error) {
				file, err := os.CreateTemp("", "shredder_test")
				if err != nil {
					return "", err
				}
				file.Close()
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

package vault

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("with valid map", func(t *testing.T) {
		input := map[string]string{
			"test1": "/path/1",
			"test2": "/path/2",
		}
		vaults := New(input)

		if !vaults.Has("test1") {
			t.Error("expected vault 'test1' to exist")
		}
		if !vaults.Has("test2") {
			t.Error("expected vault 'test2' to exist")
		}
	})

	t.Run("with nil map", func(t *testing.T) {
		vaults := New(nil)

		if vaults.Has("anything") {
			t.Error("expected empty vaults")
		}
	})

	t.Run("with empty map", func(t *testing.T) {
		vaults := New(map[string]string{})

		if vaults.Has("anything") {
			t.Error("expected empty vaults")
		}
	})
}

func TestDefaultVaults(t *testing.T) {
	vaults := DefaultVaults()

	expectedNames := []string{"develop", "private", "work"}
	actualNames := vaults.Names()

	if !reflect.DeepEqual(expectedNames, actualNames) {
		t.Errorf("expected names %v, got %v", expectedNames, actualNames)
	}

	if !vaults.Has("develop") {
		t.Error("expected 'develop' vault to exist")
	}
	if !vaults.Has("work") {
		t.Error("expected 'work' vault to exist")
	}
	if !vaults.Has("private") {
		t.Error("expected 'private' vault to exist")
	}
}

func TestHas(t *testing.T) {
	vaults := New(map[string]string{
		"test1": "/path/1",
		"test2": "/path/2",
	})

	tests := []struct {
		name      string
		vaultName string
		expected  bool
	}{
		{"existing vault test1", "test1", true},
		{"existing vault test2", "test2", true},
		{"non-existing vault", "test3", false},
		{"empty name", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := vaults.Has(tt.vaultName)
			if result != tt.expected {
				t.Errorf("expected Has(%q) = %v, got %v", tt.vaultName, tt.expected, result)
			}
		})
	}
}

func TestPath(t *testing.T) {
	vaults := New(map[string]string{
		"dev":  "/home/dev",
		"prod": "/home/prod",
	})

	tests := []struct {
		name         string
		vaultName    string
		expectedPath string
		expectedOk   bool
	}{
		{"existing vault dev", "dev", "/home/dev", true},
		{"existing vault prod", "prod", "/home/prod", true},
		{"non-existing vault", "staging", "", false},
		{"empty name", "", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, ok := vaults.Path(tt.vaultName)
			if ok != tt.expectedOk {
				t.Errorf("expected ok = %v, got %v", tt.expectedOk, ok)
			}
			if path != tt.expectedPath {
				t.Errorf("expected path %q, got %q", tt.expectedPath, path)
			}
		})
	}
}

func TestNames(t *testing.T) {
	tests := []struct {
		name          string
		vaults        map[string]string
		expectedNames []string
	}{
		{
			name: "multiple vaults sorted",
			vaults: map[string]string{
				"zebra": "/z",
				"alpha": "/a",
				"beta":  "/b",
			},
			expectedNames: []string{"alpha", "beta", "zebra"},
		},
		{
			name:          "empty vaults",
			vaults:        map[string]string{},
			expectedNames: []string{},
		},
		{
			name: "single vault",
			vaults: map[string]string{
				"only": "/path",
			},
			expectedNames: []string{"only"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := New(tt.vaults)
			names := v.Names()

			if !reflect.DeepEqual(names, tt.expectedNames) {
				t.Errorf("expected names %v, got %v", tt.expectedNames, names)
			}
		})
	}
}

func TestNewVaults(t *testing.T) {
	vaults := NewVaults(map[string]string{
		"test": "/path",
	})

	if !vaults.Has("test") {
		t.Error("NewVaults should create valid vaults")
	}

	path, ok := vaults.Path("test")
	if !ok || path != "/path" {
		t.Errorf("expected path '/path', got %q (ok=%v)", path, ok)
	}
}

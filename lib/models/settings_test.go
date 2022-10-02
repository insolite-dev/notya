//
// Copyright 2021-present Insolite. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.
//

package models_test

import (
	"testing"

	"github.com/insolite-dev/notya/lib/models"
)

func TestInitSettings(t *testing.T) {
	tests := []struct {
		testname string
		expected models.Settings
	}{
		{
			testname: "should return initial settings properly",
			expected: models.Settings{Editor: models.DefaultEditor, LocalPath: models.DefaultLocalPath},
		},
	}

	for _, td := range tests {
		t.Run(td.testname, func(t *testing.T) {
			got := models.InitSettings("notya")

			if got.Editor != td.expected.Editor || got.LocalPath != td.expected.LocalPath {
				t.Errorf("InitSettings's sum was different: Want: %v | Got: %v", got, td.expected)
			}
		})
	}
}

func TestToByte(t *testing.T) {
	tests := []struct {
		testname       string
		model          models.Settings
		expectedLength int
	}{
		{
			testname:       "should return initial settings properly",
			model:          models.Settings{Editor: "mvim"},
			expectedLength: 43,
		},
	}

	for _, td := range tests {
		t.Run(td.testname, func(t *testing.T) {
			got := td.model.ToByte()

			if len(got) != td.expectedLength {
				t.Errorf("ToByte's length sum was different: Want: %v | Got: %v", td.expectedLength, len(got))
			}
		})
	}
}

func TestToJSON(t *testing.T) {
	tests := []struct {
		model    models.Settings
		expected map[string]interface{}
	}{
		{
			model: models.Settings{
				Name:               models.DefaultAppName,
				Editor:             models.DefaultEditor,
				LocalPath:          "~notya",
				FirebaseProjectID:  "notya",
				FirebaseAccountKey: "~notya/key.json",
				FirebaseCollection: "notya-notes",
			},
			expected: map[string]interface{}{
				"name":             models.DefaultAppName,
				"editor":           models.DefaultEditor,
				"local_path":       "~notya",
				"fire_project_id":  "notya",
				"fire_account_key": "~notya/key.json",
				"fire_collection":  "notya-notes",
			},
		},
	}

	for _, td := range tests {
		got := td.model.ToJSON()

		for key, value := range td.expected {

			if got[key] != value {
				t.Errorf("SettingsToJSON's sum was different: Want: %v | Got: %v", got, td.expected)
			}
		}
	}
}

func TestDecodeSettings(t *testing.T) {
	tests := []struct {
		testname      string
		argumentValue string
		expected      models.Settings
	}{
		{
			testname:      "should generate settings model from json properly",
			argumentValue: `{"editor": "vi"}`,
			expected:      models.Settings{Editor: models.DefaultEditor},
		},
	}

	for _, td := range tests {
		t.Run(td.testname, func(t *testing.T) {
			got := models.DecodeSettings(td.argumentValue)

			if got.Editor != td.expected.Editor {
				t.Errorf("DecodeSettings's sum was different: Want: %v | Got: %v", got, td.expected)
			}
		})
	}
}

func TestFirePath(t *testing.T) {
	tests := []struct {
		model    models.Settings
		expected string
	}{
		{
			model:    models.Settings{},
			expected: "notya",
		},
		{
			model:    models.Settings{Name: "notya"},
			expected: "notya",
		},
		{
			model:    models.Settings{FirebaseCollection: "notya-notes", Name: "notya"},
			expected: "notya-notes",
		},
	}

	for _, td := range tests {
		got := td.model.FirePath()

		if got != td.expected {
			t.Errorf("FirePath's sum was different: Want: %v | Got: %v", got, td.expected)
		}
	}
}

func TestIsValid(t *testing.T) {
	tests := []struct {
		testname string
		settings models.Settings
		expected bool
	}{
		{
			testname: "should check settings validness correctly | [valid]",
			settings: models.InitSettings("/usr/mock/localpath"),
			expected: true,
		},
		{
			testname: "should check settings validness correctly | [invalid]",
			settings: models.Settings{},
			expected: false,
		},
	}

	for _, td := range tests {
		t.Run(td.testname, func(t *testing.T) {
			got := td.settings.IsValid()

			if got != td.expected {
				t.Errorf("IsValid sum was different: Want: %v | Got: %v", got, td.expected)
			}
		})
	}
}

func TestIsFirebaseEnabled(t *testing.T) {
	tests := []struct {
		settings models.Settings
		expected bool
	}{
		{
			settings: models.InitSettings("/usr/mock/localpath"),
			expected: false,
		},
		{
			settings: models.Settings{FirebaseProjectID: "mock-project-id"},
			expected: true,
		},
	}

	for _, td := range tests {
		got := td.settings.IsFirebaseEnabled()

		if got != td.expected {
			t.Errorf("IsFirebaseEnabled sum was different: Want: %v | Got: %v", got, td.expected)
		}
	}
}

func TestIsUpdated(t *testing.T) {
	tests := []struct {
		testname     string
		old, current models.Settings
		expected     bool
	}{
		{
			testname: "should check properly if fulls settings is updated",
			old:      models.Settings{Editor: models.DefaultEditor},
			current:  models.Settings{Editor: models.DefaultEditor},
			expected: false,
		},
		{
			testname: "should check properly if fulls settings is updated",
			old:      models.Settings{Editor: "code"},
			current:  models.Settings{Editor: models.DefaultEditor},
			expected: true,
		},
	}

	for _, td := range tests {
		t.Run(td.testname, func(t *testing.T) {
			got := models.IsUpdated(td.old, td.current)

			if got != td.expected {
				t.Errorf("IsUpdated sum was different: Want: %v | Got: %v", got, td.expected)
			}
		})
	}
}

func TestIsPathUpdated(t *testing.T) {
	tests := []struct {
		serviceType  string
		old, current models.Settings
		expected     bool
	}{
		{
			serviceType: "LOCAL",
			old:         models.Settings{LocalPath: "test/path"},
			current:     models.Settings{LocalPath: "test/path"},
			expected:    false,
		},
		{
			serviceType: "LOCAL",
			old:         models.Settings{LocalPath: "test/path"},
			current:     models.Settings{LocalPath: "new/test/path"},
			expected:    true,
		},
		{
			serviceType: "LOCAL",
			old:         models.Settings{Editor: "code"},
			current:     models.Settings{Editor: models.DefaultEditor},
			expected:    false,
		},
		{
			serviceType: "FIREBASE",
			old:         models.Settings{FirebaseCollection: "test/path"},
			current:     models.Settings{FirebaseCollection: "test/path"},
			expected:    false,
		},
		{
			serviceType: "FIREBASE",
			old:         models.Settings{FirebaseCollection: "test/path"},
			current:     models.Settings{FirebaseCollection: "new/test/path"},
			expected:    true,
		},
		{
			serviceType: "undefined",
			old:         models.Settings{FirebaseCollection: "test/path"},
			current:     models.Settings{FirebaseCollection: "new/test/path"},
			expected:    false,
		},
	}

	for i, td := range tests {
		got := models.IsPathUpdated(td.old, td.current, td.serviceType)

		if got != td.expected {
			t.Errorf("IsUpdated[%v] sum was different: Want: %v | Got: %v", i, got, td.expected)
		}
	}
}

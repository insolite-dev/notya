// Copyright 2021-2022 present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package assets_test

import (
	"testing"

	"github.com/AlecAivazis/survey/v2"
	"github.com/anonistas/notya/assets"
	"github.com/anonistas/notya/lib/models"
)

func TestChoseNotePrompt(t *testing.T) {
	type arguments struct {
		msg     string
		options []string
	}

	tests := []struct {
		testname string
		args     arguments
		expected survey.Select
	}{
		{
			testname: "should generate choosing-note prompt properly",
			args: arguments{
				msg:     "edit",
				options: []string{"1", "2", "3"},
			},
			expected: survey.Select{
				Message: "Choose a note to edit:",
				Options: []string{"1", "2", "3"},
			},
		},
	}

	for _, td := range tests {
		t.Run(td.testname, func(t *testing.T) {
			got := assets.ChooseNotePrompt(td.args.msg, td.args.options)

			// Closure function to check if options are different or not.
			var isDiffArr = func() bool {
				var a1, a2 = got.Options, td.expected.Options
				if len(a1) != len(a2) {
					return true
				}

				for i := 0; i < len(a1); i++ {
					if a1[i] != a2[i] {
						return true
					}
				}

				return false
			}()

			if got.Message != td.expected.Message || isDiffArr {
				t.Errorf("Sum of ChooseNotePrompt was different: Want: %v | Got: %v", td.expected, got)
			}
		})
	}
}

func TestNewNamePrompt(t *testing.T) {
	tests := []struct {
		testname     string
		defaultValue string
		expected     survey.Input
	}{
		{
			testname:     "should generate new-name-prompt properly",
			defaultValue: "default-name",
			expected: survey.Input{
				Message: "New name: ",
				Help:    "Enter new note name/title (don't forget putting type of it, like: `renamed_note.txt`)",
				Default: "default-name",
			},
		},
	}

	for _, td := range tests {
		t.Run(td.testname, func(t *testing.T) {
			got := assets.NewNamePrompt(td.defaultValue)

			if got.Message != td.expected.Message || got.Help != td.expected.Help || got.Default != td.expected.Default {
				t.Errorf("Sum of NewNamePrompt was different: Want: %v | Got: %v", td.expected, got)
			}
		})
	}
}

func TestSettingsEditPromptQuestions(t *testing.T) {
	tests := []struct {
		testname        string
		defaultSettings models.Settings
		expected        []*survey.Question
	}{
		{
			testname:        "should generate settings-edit prompt questions properly",
			defaultSettings: models.InitSettings("default_path"),
			expected: []*survey.Question{
				{
					Name: "editor",
					Prompt: &survey.Input{
						Default: models.InitSettings("default_path").Editor,
						Message: "Editor",
						Help:    "Editor for notya. --> vim/nvim/code/code-insiders ...",
					},
					Validate: survey.MinLength(1),
				},
				{
					Name: "local_path",
					Prompt: &survey.Input{
						Default: models.InitSettings("default_path").LocalPath,
						Message: "Local Path",
						Help:    "Local path of notya base working directory",
					},
				},
			},
		},
	}

	for _, td := range tests {
		t.Run(td.testname, func(t *testing.T) {
			got := assets.SettingsEditPromptQuestions(td.defaultSettings)

			//
			var isDiff = func() bool {
				if len(got) != len(td.expected) {
					return true
				}

				for i := 0; i < len(got); i++ {
					var a1, a2 = got[i], td.expected[i]
					if a1.Name != a2.Name {
						return true
					}
				}

				return false
			}()

			if isDiff {
				t.Errorf("Sum of SettingsEditPromptQuestions was different: Want: %v | Got: %v", got, td.expected)
			}
		})
	}
}
// Copyright 2021-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package pkg

import "github.com/AlecAivazis/survey/v2"

// Version is current version of application.
const Version = "v0.1.0"

var (
	// Custom configuration for survey icons and colors.
	// See [https://github.com/mgutz/ansi#style-format] for details.
	SurveyIconsConfig = func(icons *survey.IconSet) {
		icons.Question.Format = "cyan"
		icons.Question.Text = "[?]"
		icons.Help.Format = "blue"
		icons.Help.Text = "Help ->"
		icons.Error.Format = "yellow"
		icons.Error.Text = "Warning ->"
	}
)

// CreateAnswers is a model structure of the [CreateNoteQuestions].
type CreateAnswers struct {
	Title    string
	EditNote bool `survey:"edit-note"`
}

var (
	// CreateNoteQuestions is a list of questions for create command.
	CreateNoteQuestions = []*survey.Question{
		{
			Name: "title",
			Prompt: &survey.Input{
				Message: "Enter name of new note: ",
				Help:    "Append to your note any name you want  and then, complete file name with special file name type | e.g: new_note.md",
			},
			Validate: survey.MinLength(1),
		},
		{
			Name: "edit-note",
			Prompt: &survey.Confirm{
				Message: "Do you wanna open note with Vi/Vim, to edit file?",
				Default: true,
			},
		},
	}
)

// Copyright 2021-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package services

import (
	"errors"
	"fmt"
	"os"

	"github.com/anonistas/notya/lib/models"
	"github.com/anonistas/notya/pkg"
)

// LocalService is a class implementation of service repo.
type LocalService struct {
	notyaPath string
	stdargs   models.StdArgs
}

// Set [LocalService] as [ServiceRepo].
var _ ServiceRepo = &LocalService{}

// NewLocalService, creates new local service by given arguments.
func NewLocalService(notyapath string, stdargs models.StdArgs) *LocalService {
	return &LocalService{notyaPath: notyapath, stdargs: stdargs}
}

// Init creates notya working directory into running machine.
func (l *LocalService) Init() error {
	// Check if working directory already exists
	if pkg.FileExists(l.notyaPath) {
		return errors.New("Notya already initialized before")
	}

	// Create new notya working directory.
	creatingErr := pkg.NewFolder(l.notyaPath)
	if creatingErr != nil {
		return creatingErr
	}

	return nil
}

// CreateNote, creates new note at [notya notes path],
// and fills it's data by given note model.
func (l *LocalService) CreateNote(note models.Note) error {
	notePath := l.notyaPath + note.Title

	// Check if file already exists.
	if pkg.FileExists(notePath) {
		alreadyExists := "A file with the name " + fmt.Sprintf("`%v`", note.Title) + " already exists"
		return errors.New(alreadyExists)
	}

	// Create new file inside notes.
	creatingErr := pkg.WriteNote(notePath, []byte(note.Body))
	if creatingErr != nil {
		return creatingErr
	}

	return nil
}

// ViewNote, opens note-file from given [note.Name], then takes it body,
// and returns new fully-filled note.
func (l *LocalService) ViewNote(note models.Note) (*models.Note, error) {
	notePath := l.notyaPath + note.Title

	// Open and read body of note.
	res, err := pkg.ReadBody(notePath)
	if err != nil {
		return nil, err
	}

	// Re-generate note with path and body.
	modifiedNote := models.Note{Title: note.Title, Path: notePath, Body: *res}

	return &modifiedNote, nil
}

// EditNote, overwrites exiting file from [notya notes path].
func (l *LocalService) EditNote(note models.Note) error {
	notePath := l.notyaPath + note.Title

	// Check if file exists or not.
	if !pkg.FileExists(notePath) {
		notExists := fmt.Sprintf("File not exists at: notya/%v", note.Title)
		return errors.New(notExists)
	}

	// Edit note's body
	if writingErr := pkg.WriteNote(notePath, []byte(note.Body)); writingErr != nil {
		return writingErr
	}

	return nil
}

// Rename, changes given note's name.
func (l *LocalService) Rename(editnote models.EditNote) error {
	// Check if it's same titles.
	if editnote.Current.Title == editnote.New.Title {
		return errors.New("Current and new name are same")
	}

	// Check if file exists at new note path.
	if pkg.FileExists(editnote.New.Path) {
		alreadyExists := fmt.Sprintf("A file exists at: notya/%v, please provide a unique name", editnote.New.Title)
		return errors.New(alreadyExists)
	}

	// Rename given note.
	if err := os.Rename(editnote.Current.Path, editnote.New.Path); err != nil {
		return err
	}

	return nil
}

// Remove, deletes given note file.
func (l *LocalService) Remove(note models.Note) error {
	notePath := l.notyaPath + note.Title

	// Delete provided file.
	if err := pkg.Delete(notePath); err != nil {
		return err
	}

	return nil
}

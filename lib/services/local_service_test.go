//
// Copyright 2021-present Insolite. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.
//

package services_test

import (
	"errors"
	"testing"

	"github.com/insolite-dev/notya/assets"
	"github.com/insolite-dev/notya/lib/models"
	"github.com/insolite-dev/notya/lib/services"
	"github.com/insolite-dev/notya/pkg"
)

// Define a mock local service implementation.
//
// Note:
// Tests are based on current-machine's local storage.
// Mocking techniques not used.
var ls = services.LocalService{
	NotyaPath: "./",
	Config:    models.Settings{LocalPath: "./", Editor: "vi"},
	Stdargs:   models.StdArgs{},
}

func TestNewLocalService(t *testing.T) {
	tests := []struct {
		stdargs  models.StdArgs
		expected services.LocalService
	}{
		{
			stdargs:  models.StdArgs{},
			expected: services.LocalService{Stdargs: models.StdArgs{}},
		},
		{
			stdargs:  ls.Stdargs,
			expected: services.LocalService{Stdargs: models.StdArgs{}},
		},
	}

	for _, td := range tests {
		got := services.NewLocalService(td.stdargs)
		if got.Stdargs != td.expected.Stdargs {
			t.Errorf("Sum of [NewLocalService] is different: Got: %v | Want: %v", got, td.expected)
		}
	}
}

func TestGeneratePath(t *testing.T) {
	tests := []struct {
		ls       services.LocalService
		title    string
		expected string
	}{
		{
			ls:       ls,
			title:    "new-note.txt",
			expected: ls.Config.LocalPath + "new-note.txt",
		},
		{
			ls: services.LocalService{
				Config: models.Settings{LocalPath: ".", Editor: "vi"},
			},
			title:    "new-note.txt",
			expected: ls.Config.LocalPath + "new-note.txt",
		},
	}

	for _, td := range tests {
		got := td.ls.GeneratePath(models.Node{Title: td.title})

		if got != td.expected {
			t.Errorf("Sum of [GeneratePath] is different: Got: %v | Want: %v", got, td.expected)
		}
	}
}

func TestType(t *testing.T) {
	tests := []struct {
		expected string
	}{
		{expected: services.LOCAL.ToStr()},
	}

	for _, td := range tests {
		got := ls.Type()

		if got != td.expected {
			t.Errorf("Sum of [Type] is different: Got: %v | Want: %v", got, td.expected)
		}
	}
}

func TestPath(t *testing.T) {
	tests := []struct {
		expected string
	}{
		{expected: ls.Path()},
	}

	for _, td := range tests {
		got := ls.Path()

		if got != td.expected {
			t.Errorf("Sum of [Path] is different: Got: %v | Want: %v", got, td.expected)
		}
	}
}

func TestStateConfig(t *testing.T) {
	tests := []struct {
		expected models.Settings
	}{
		{expected: ls.Config},
	}

	for _, td := range tests {
		got := ls.StateConfig()

		if got != td.expected {
			t.Errorf("Sum of [StateConfig] is different: Got: %v | Want: %v", got, td.expected)
		}
	}
}

func TestInit(t *testing.T) {
	tests := []struct {
		localService services.LocalService
		beforeAct    func()
		afterAct     func()
		expected     error
	}{
		{
			localService: services.LocalService{
				NotyaPath: "mock/local-path/",
				Config:    models.Settings{LocalPath: "mock/local-path/"},
			},
			beforeAct: func() {},
			afterAct: func() {
				notyaPath, _ := pkg.NotyaPWD(models.Settings{})
				_ = pkg.Delete(*notyaPath + "/" + models.SettingsName)
			},
			expected: errors.New("mkdir mock/local-path/: no such file or directory"),
		},
		{
			localService: services.LocalService{
				Config: models.Settings{LocalPath: ".notya-mocktests"},
			},
			beforeAct: func() {
				notyaPath, _ := pkg.NotyaPWD(models.Settings{LocalPath: ".notya-mocktests"})
				_ = pkg.NewFolder(*notyaPath)

				s := models.InitSettings("")
				_ = pkg.WriteNote(*notyaPath+"/"+models.SettingsName, s.ToByte())
			},
			afterAct: func() {
				notyaPath, _ := pkg.NotyaPWD(models.Settings{LocalPath: ".notya-mocktests"})
				pkg.Delete(*notyaPath + "/" + models.SettingsName)
				pkg.Delete(*notyaPath + "/")
			},
			expected: nil,
		},
		{
			localService: services.LocalService{
				Config: models.Settings{LocalPath: ".notya-mocktests"},
			},
			beforeAct: func() {
				notyaPath, _ := pkg.NotyaPWD(models.Settings{LocalPath: ".notya-mocktests"})
				_ = pkg.NewFolder(*notyaPath)
			},
			afterAct: func() {
				notyaPath, _ := pkg.NotyaPWD(models.Settings{LocalPath: ".notya-mocktests"})
				pkg.Delete(*notyaPath + "/" + models.SettingsName)
				pkg.Delete(*notyaPath + "/")
			},
			expected: nil,
		},
	}

	for _, td := range tests {
		td.beforeAct()
		got := td.localService.Init()
		td.afterAct()

		if (got == nil || td.expected == nil) && got != td.expected ||
			(got != nil && td.expected != nil) && got.Error() != td.expected.Error() {
			t.Errorf("Sum of [Init] is different: Got: %v | Want: %v", got, td.expected)
		}
	}
}

func TestSettings(t *testing.T) {
	tests := []struct {
		localService  services.LocalService
		arg           *string
		beforeAct     func()
		afterAct      func()
		expectedError error
		expected      models.Settings
	}{
		{
			localService: services.LocalService{
				Config:    models.Settings{LocalPath: "./"},
				NotyaPath: "./",
			},
			beforeAct: func() {
				s := models.InitSettings("")
				_ = pkg.WriteNote(models.SettingsName, s.ToByte())
			},
			afterAct: func() {
				_ = pkg.Delete(models.SettingsName)
			},
			expectedError: nil,
			expected:      models.InitSettings(""),
		},
	}

	for _, td := range tests {
		td.beforeAct()
		got, err := td.localService.Settings(td.arg)
		td.afterAct()

		if got.Editor != td.expected.Editor || got.LocalPath != td.expected.LocalPath {
			t.Errorf("Sum of [Settigns] is different: Got: %v | Want: %v", got, td.expected)
		}

		if err != td.expectedError {
			t.Errorf("Error Sum of [Settigns] is different: Got: %v | Want: %v", err, td.expectedError)
		}
	}
}

func TestWriteSettings(t *testing.T) {
	tests := []struct {
		model        models.Settings
		localService services.LocalService
		beforeAct    func()
		afterAct     func()
		expected     error
	}{
		{
			model: models.Settings{},
			localService: services.LocalService{
				Config:    models.Settings{LocalPath: "./"},
				NotyaPath: "./",
			},
			beforeAct: func() {},
			afterAct:  func() {},
			expected:  assets.InvalidSettingsData,
		},
		{
			model: models.InitSettings("./"),
			localService: services.LocalService{
				Config:    models.Settings{LocalPath: "./"},
				NotyaPath: "./",
			},
			beforeAct: func() {},
			afterAct: func() {
				_ = pkg.Delete(models.SettingsName)
			},
			expected: nil,
		},
	}

	for _, td := range tests {
		td.beforeAct()
		got := td.localService.WriteSettings(td.model)
		td.afterAct()

		if got != td.expected {
			t.Errorf("Sum of [WriteSettings] is different: Got: %v | Want: %v", got, td.expected)
		}
	}
}

func TestOpenSettings(t *testing.T) {
	ls := services.LocalService{
		NotyaPath: "./",
		Config:    models.Settings{LocalPath: "./", Editor: "vi"},
		Stdargs:   models.StdArgs{},
	}

	tests := []struct {
		settings     models.Settings
		localService services.LocalService
		beforeAct    func(title string)
		afterAct     func(title string)
		expected     error
	}{
		{
			settings:     models.Settings{ID: "somerandomnotethatnotexists"},
			localService: ls,
			beforeAct:    func(title string) {},
			afterAct:     func(title string) {},
			expected:     assets.NotExists("somerandomnotethatnotexists", "File or Directory"),
		},
		{
			settings:     models.Settings{ID: ""},
			localService: ls,
			beforeAct:    func(title string) {},
			afterAct:     func(title string) {},
			expected:     assets.NotExists(models.SettingsName, "File or Directory"),
		},
		{

			settings:     models.Settings{ID: "somerandomdirthatexists"},
			localService: ls,
			beforeAct: func(title string) {
				path := ls.GeneratePath(models.Node{Title: title})
				_ = pkg.WriteNote(path, []byte{})
			},
			afterAct: func(title string) {
				path := ls.GeneratePath(models.Node{Title: title})
				_ = pkg.Delete(path)
			},
			expected: errors.New("exit status 2"),
		},
	}

	for _, td := range tests {
		td.beforeAct(td.settings.ID)
		got := td.localService.OpenSettings(td.settings)
		td.afterAct(td.settings.ID)

		if (got == nil || td.expected == nil) && got != td.expected ||
			(got != nil && td.expected != nil) && got.Error() != td.expected.Error() {
			t.Errorf("Sum of [OpenSettings] is different: Got: %v | Want: %v", got, td.expected)
		}
	}
}

func TestOpen(t *testing.T) {
	ls := services.LocalService{
		NotyaPath: "./",
		Config:    models.Settings{LocalPath: "./", Editor: "vi"},
		Stdargs:   models.StdArgs{},
	}

	tests := []struct {
		node         models.Node
		localService services.LocalService
		beforeAct    func(node models.Node)
		afterAct     func(node models.Node)
		expected     error
	}{
		{
			node:         models.Node{Title: "somerandomnotethatnotexists"},
			localService: ls,
			beforeAct:    func(node models.Node) {},
			afterAct:     func(node models.Node) {},
			expected:     assets.NotExists("somerandomnotethatnotexists", "File or Directory"),
		},
		{
			node:         models.Node{Title: ""},
			localService: ls,
			beforeAct:    func(node models.Node) {},
			afterAct:     func(node models.Node) {},
			expected:     errors.New("exit status 2"),
		},
		{
			node:         models.Node{Title: "somerandomnote.txt"},
			localService: ls,
			beforeAct: func(node models.Node) {
				path := ls.GeneratePath(node)
				_ = pkg.WriteNote(path, []byte{})
			},
			afterAct: func(node models.Node) {
				path := ls.GeneratePath(node)
				_ = pkg.Delete(path)
			},
			expected: errors.New("exit status 2"),
		},
	}

	for _, td := range tests {
		td.beforeAct(td.node)
		got := td.localService.Open(td.node)
		td.afterAct(td.node)

		if (got == nil || td.expected == nil) && got != td.expected ||
			(got != nil && td.expected != nil) && got.Error() != td.expected.Error() {
			t.Errorf("Sum of [Open] is different: Got: %v | Want: %v", got, td.expected)
		}
	}
}

func TestRemove(t *testing.T) {
	ls := services.LocalService{
		NotyaPath: "./test/",
		Config:    models.Settings{LocalPath: "./test/", Editor: "vi"},
		Stdargs:   models.StdArgs{},
	}

	tests := []struct {
		node         models.Node
		localService services.LocalService
		beforeAct    func(node models.Node)
		afterAct     func(node models.Node)
		expected     error
	}{
		{
			node:         models.Node{Title: "somerandomnotethatnotexists"},
			localService: ls,
			beforeAct:    func(node models.Node) {},
			afterAct:     func(node models.Node) {},
			expected:     assets.NotExists("somerandomnotethatnotexists", "File or Directory"),
		},
		{
			node:         models.Node{Title: "newfile"},
			localService: ls,
			beforeAct:    func(node models.Node) {},
			afterAct:     func(node models.Node) {},
			expected:     assets.NotExists("newfile", "File or Directory"),
		},
		{
			node:         models.Node{Title: ".mock-folder"},
			localService: ls,
			beforeAct: func(node models.Node) {
				path := ls.GeneratePath(node)
				_ = pkg.NewFolder(ls.Path())
				_ = pkg.NewFolder(path)
				_ = pkg.WriteNote(path+"/"+"mock_note.txt", []byte{})
			},
			afterAct: func(node models.Node) {
				path := ls.GeneratePath(node)
				_ = pkg.Delete(path + "/" + "mock_note.txt")
				_ = pkg.Delete(path)
				_ = pkg.Delete(ls.Path())
			},
			expected: nil,
		},
		{
			node:         models.Node{Title: "somerandomnote.txt"},
			localService: ls,
			beforeAct: func(node models.Node) {
				path := ls.GeneratePath(node)
				_ = pkg.NewFolder(ls.Path())
				_ = pkg.WriteNote(path, []byte{})
			},
			afterAct: func(node models.Node) {
				_ = pkg.Delete(ls.Path())
			},
			expected: nil,
		},
	}

	for _, td := range tests {
		td.beforeAct(td.node)
		got := td.localService.Remove(td.node)
		td.afterAct(td.node)

		if (got == nil || td.expected == nil) && got != td.expected ||
			(got != nil && td.expected != nil) && got.Error() != td.expected.Error() {
			t.Errorf("Sum of [Remove] is different: Got: %v | Want: %v", got, td.expected)
		}
	}
}

func TestCreate(t *testing.T) {
	tests := []struct {
		note         models.Note
		localService services.LocalService
		beforeAct    func(note models.Note)
		afterAct     func(note models.Note)
		expected     error
	}{
		{
			note:         models.Note{Title: "somerandomnotethatexists"},
			localService: ls,
			beforeAct: func(note models.Note) {
				path := ls.GeneratePath(note.ToNode())
				_ = pkg.WriteNote(path, []byte{})
			},
			afterAct: func(note models.Note) {
				path := ls.GeneratePath(note.ToNode())
				_ = pkg.Delete(path)
			},
			expected: assets.AlreadyExists("somerandomnotethatexists", "file"),
		},
		{
			note:         models.Note{Title: "mocknote.txt"},
			localService: ls,
			beforeAct: func(note models.Note) {
			},
			afterAct: func(note models.Note) {
				_ = pkg.Delete(ls.GeneratePath(note.ToNode()))
			},
			expected: nil,
		},
	}

	for _, td := range tests {
		td.beforeAct(td.note)
		_, got := td.localService.Create(td.note)
		td.afterAct(td.note)

		if (got == nil || td.expected == nil) && got != td.expected ||
			(got != nil && td.expected != nil) && got.Error() != td.expected.Error() {
			t.Errorf("Sum of [Create] is different: Got: %v | Want: %v", got, td.expected)
		}
	}
}

func TestView(t *testing.T) {
	tests := []struct {
		note         models.Note
		localService services.LocalService
		beforeAct    func(note models.Note)
		afterAct     func(note models.Note)
		expected     *models.Note
		expectedErr  error
	}{
		{
			note:         models.Note{Title: "somerandomnotethatnotexists"},
			localService: ls,
			beforeAct: func(note models.Note) {
				path := ls.GeneratePath(note.ToNode())
				_ = pkg.Delete(path)
			},
			afterAct:    func(note models.Note) {},
			expected:    nil,
			expectedErr: assets.NotExists("somerandomnotethatnotexists", "File"),
		},
		{
			note:         models.Note{Title: "mocknote.txt"},
			localService: ls,
			beforeAct: func(note models.Note) {
				path := ls.GeneratePath(note.ToNode())
				_ = pkg.WriteNote(path, []byte{})
			},
			afterAct: func(note models.Note) {
				path := ls.GeneratePath(note.ToNode())
				_ = pkg.Delete(path)
			},
			expected:    &models.Note{Title: "mocknote.txt", Body: string([]byte{})},
			expectedErr: nil,
		},
	}

	for _, td := range tests {
		td.beforeAct(td.note)
		gotRes, gotErr := td.localService.View(td.note)
		td.afterAct(td.note)

		if (gotRes == nil || td.expected == nil) && gotRes != td.expected ||
			(gotRes != nil && td.expected != nil) && (gotRes.Title != td.expected.Title || gotRes.Body != td.expected.Body) {
			t.Errorf("Sum of {res}[View] is different: Got: %v | Want: %v", gotRes, td.expected)
		}

		if (gotErr == nil || td.expectedErr == nil) && gotErr != td.expectedErr ||
			(gotErr != nil && td.expectedErr != nil) && gotErr.Error() != td.expectedErr.Error() {
			t.Errorf("Sum of {error}[View] is different: Got: %v | Want: %v", gotErr, td.expectedErr)
		}
	}
}

func TestEdit(t *testing.T) {
	tests := []struct {
		note         models.Note
		localService services.LocalService
		beforeAct    func(note models.Note)
		afterAct     func(note models.Note)
		expected     *models.Note
		expectedErr  error
	}{
		{
			note:         models.Note{Title: "somerandomnotethatnotexists"},
			localService: ls,
			beforeAct: func(note models.Note) {
				path := ls.GeneratePath(note.ToNode())
				_ = pkg.Delete(path)
			},
			afterAct:    func(note models.Note) {},
			expected:    nil,
			expectedErr: assets.NotExists("somerandomnotethatnotexists", "File"),
		},
		{
			note:         models.Note{Title: "mocknote.txt", Body: "empty-body"},
			localService: ls,
			beforeAct: func(note models.Note) {
				path := ls.GeneratePath(note.ToNode())
				_ = pkg.WriteNote(path, []byte{})
			},
			afterAct: func(note models.Note) {
				path := ls.GeneratePath(note.ToNode())
				_ = pkg.Delete(path)
			},
			expected:    &models.Note{Title: "mocknote.txt", Body: "empty-body"},
			expectedErr: nil,
		},
	}

	for _, td := range tests {
		td.beforeAct(td.note)
		gotRes, gotErr := td.localService.Edit(td.note)
		td.afterAct(td.note)

		if (gotRes == nil || td.expected == nil) && gotRes != td.expected ||
			(gotRes != nil && td.expected != nil) && (gotRes.Title != td.expected.Title || gotRes.Body != td.expected.Body) {
			t.Errorf("Sum of {res}[Edit] is different: Got: %v | Want: %v", gotRes, td.expected)
		}

		if (gotErr == nil || td.expectedErr == nil) && gotErr != td.expectedErr ||
			(gotErr != nil && td.expectedErr != nil) && gotErr.Error() != td.expectedErr.Error() {
			t.Errorf("Sum of {error}[Edit] is different: Got: %v | Want: %v", gotErr, td.expectedErr)
		}
	}
}

func TestCopy(t *testing.T) {
	tests := []struct {
		note         models.Note
		localService services.LocalService
		beforeAct    func(note models.Note)
		afterAct     func(note models.Note)
		expected     error
	}{
		{
			note:         models.Note{Title: "somerandomnotethatexists"},
			localService: ls,
			beforeAct: func(note models.Note) {
				path := ls.GeneratePath(note.ToNode())
				_ = pkg.WriteNote(path, []byte{})
			},
			afterAct: func(note models.Note) {
				path := ls.GeneratePath(note.ToNode())
				_ = pkg.Delete(path)
			},
			expected: nil,
		},
		{
			note:         models.Note{Title: "mocknote.txt"},
			localService: ls,
			beforeAct: func(note models.Note) {
			},
			afterAct: func(note models.Note) {
			},
			expected: assets.NotExists("mocknote.txt", "File"),
		},
	}

	for _, td := range tests {
		td.beforeAct(td.note)
		got := td.localService.Copy(td.note)
		td.afterAct(td.note)

		if (got == nil || td.expected == nil) && got != td.expected ||
			(got != nil && td.expected != nil) && got.Error() != td.expected.Error() {
			t.Errorf("Sum of [Copy] is different: Got: %v | Want: %v", got, td.expected)
		}
	}
}

func TestRename(t *testing.T) {
	tests := []struct {
		editnode     models.EditNode
		localService services.LocalService
		beforeAct    func(ed models.EditNode)
		afterAct     func(ed models.EditNode)
		expected     error
	}{
		{
			editnode: models.EditNode{
				Current: models.Node{Title: ".current-note"},
				New:     models.Node{Title: ".new-note"},
			},
			localService: ls,
			beforeAct: func(ed models.EditNode) {
				_ = pkg.Delete(ls.GeneratePath(ed.Current))
			},
			afterAct: func(ed models.EditNode) {},
			expected: assets.NotExists(".current-note", "File or Directory"),
		},
		{
			editnode: models.EditNode{
				Current: models.Node{Title: ".same-name-note"},
				New:     models.Node{Title: ".same-name-note"},
			},
			localService: ls,
			beforeAct: func(ed models.EditNode) {
				path := ls.GeneratePath(ed.Current)
				_ = pkg.WriteNote(path, []byte{})
			},
			afterAct: func(ed models.EditNode) {
				_ = pkg.Delete(ls.GeneratePath(ed.Current))
			},
			expected: assets.SameTitles,
		},
		{
			editnode: models.EditNode{
				Current: models.Node{Title: ".current-note"},
				New:     models.Node{Title: ".new-note"},
			},
			localService: ls,
			beforeAct: func(ed models.EditNode) {
				_ = pkg.WriteNote(ls.GeneratePath(ed.Current), []byte{})
				_ = pkg.WriteNote(ls.GeneratePath(ed.New), []byte{})
			},
			afterAct: func(ed models.EditNode) {
				_ = pkg.Delete(ls.GeneratePath(ed.Current))
				_ = pkg.Delete(ls.GeneratePath(ed.New))
			},
			expected: assets.AlreadyExists(".new-note", "File or Directory"),
		},
		{
			editnode: models.EditNode{
				Current: models.Node{Title: ".current-note"},
				New:     models.Node{Title: ".new-note"},
			},
			localService: ls,
			beforeAct: func(ed models.EditNode) {
				_ = pkg.WriteNote(ls.GeneratePath(ed.Current), []byte{})
			},
			afterAct: func(ed models.EditNode) {
				_ = pkg.Delete(ls.GeneratePath(ed.New))
			},
			expected: nil,
		},
	}

	for _, td := range tests {
		td.beforeAct(td.editnode)
		got := td.localService.Rename(td.editnode)
		td.afterAct(td.editnode)

		if (got == nil || td.expected == nil) && got != td.expected ||
			(got != nil && td.expected != nil) && got.Error() != td.expected.Error() {
			t.Errorf("Sum of [Rename] is different: Got: %v | Want: %v", got, td.expected)
		}
	}
}

func TestMkdir(t *testing.T) {
	tests := []struct {
		dir          models.Folder
		localService services.LocalService
		beforeAct    func(dir models.Folder)
		afterAct     func(dir models.Folder)
		expected     error
	}{
		{
			dir:          models.Folder{Title: "somerandomdirthatexists"},
			localService: ls,
			beforeAct: func(dir models.Folder) {
				path := ls.GeneratePath(dir.ToNode())
				_ = pkg.NewFolder(path)
			},
			afterAct: func(dir models.Folder) {
				path := ls.GeneratePath(dir.ToNode())
				_ = pkg.Delete(path)
			},
			expected: assets.AlreadyExists("./somerandomdirthatexists/", "directory"),
		},
		{
			dir:          models.Folder{Title: "mocknote"},
			localService: ls,
			beforeAct: func(dir models.Folder) {
				path := ls.GeneratePath(dir.ToNode())
				_ = pkg.Delete(path)
			},
			afterAct: func(dir models.Folder) {
				path := ls.GeneratePath(dir.ToNode())
				_ = pkg.Delete(path)
			},
			expected: nil,
		},
	}

	for _, td := range tests {
		td.beforeAct(td.dir)
		_, got := td.localService.Mkdir(td.dir)
		td.afterAct(td.dir)

		if (got == nil || td.expected == nil) && got != td.expected ||
			(got != nil && td.expected != nil) && got.Error() != td.expected.Error() {
			t.Errorf("Sum of [Mkdir] is different: Got: %v | Want: %v", got, td.expected)
		}
	}
}

func TestGetAll(t *testing.T) {
	gLS := services.LocalService{
		NotyaPath: "./.testmocks/",
		Config:    models.Settings{LocalPath: "./.testmocks/"},
	}

	tests := []struct {
		localService services.LocalService
		beforeAct    func(dir string)
		afterAct     func(dir string)
		expected     []models.Note
		expectedErr  error
	}{
		{
			localService: gLS,
			beforeAct: func(dir string) {
				_ = pkg.NewFolder(dir)
			},
			afterAct: func(dir string) {
				_ = pkg.Delete(dir)
			},
			expected:    nil,
			expectedErr: assets.EmptyWorkingDirectory,
		},
		{
			localService: gLS,
			beforeAct: func(dir string) {
				_ = pkg.NewFolder(dir)
				_ = pkg.WriteNote(dir+".new-note.txt", []byte{})
				_ = pkg.WriteNote(dir+".new-note-1.txt", []byte{})
			},
			afterAct: func(dir string) {
				pkg.Delete(dir + ".new-note.txt")
				pkg.Delete(dir + ".new-note-1.txt")
				pkg.Delete(dir)
			},
			expected: []models.Note{
				{Title: ".new-note-1.txt", Path: gLS.NotyaPath + ".new-note-1.txt"},
				{Title: ".new-note.txt", Path: gLS.NotyaPath + ".new-note.txt"},
			},
			expectedErr: nil,
		},
	}

	for _, td := range tests {
		td.beforeAct(td.localService.NotyaPath)
		gotRes, _, gotErr := td.localService.GetAll("", models.NotyaIgnoreFiles)
		td.afterAct(td.localService.NotyaPath)

		for i, got := range gotRes {
			if got.Title != td.expected[i].Title || got.Path != td.expected[i].Path {
				t.Errorf("Sum of {res -> index:%v}[GetAll] is different: Got: %v | Want: %v", i, got, td.expected[i])
			}
		}

		if (gotErr == nil || td.expectedErr == nil) && gotErr != td.expectedErr ||
			(gotErr != nil && td.expectedErr != nil) && gotErr.Error() != td.expectedErr.Error() {
			t.Errorf("Sum of {error}[GetAll] is different: Got: %v | Want: %v", gotErr, td.expectedErr)
		}
	}
}

func TestMoveNotes(t *testing.T) {
	tests := []struct {
		settings     models.Settings
		localService services.LocalService
		beforeAct    func(oldS, newS models.Settings)
		afterAct     func(oldS, newS models.Settings)
		expected     error
	}{
		{
			settings: models.Settings{LocalPath: ""},
			localService: services.LocalService{
				NotyaPath: "./.testmocks/",
				Config:    models.Settings{LocalPath: "./.testmocks/"},
			},
			beforeAct: func(oldS, newS models.Settings) {
				_ = pkg.NewFolder(oldS.LocalPath)
			},
			afterAct: func(oldS, newS models.Settings) {
				_ = pkg.Delete(oldS.LocalPath)
			},
			expected: assets.EmptyWorkingDirectory,
		},
		{
			localService: services.LocalService{
				NotyaPath: "./.testmocks/",
				Config:    models.Settings{LocalPath: "./.testmocks/"},
			},
			settings: models.Settings{LocalPath: "./.testmocks-1/"},
			beforeAct: func(oldS, newS models.Settings) {
				_ = pkg.NewFolder(oldS.LocalPath)
				_ = pkg.WriteNote(oldS.LocalPath+".note.txt", []byte{})
				_ = pkg.NewFolder(newS.LocalPath)
				_ = pkg.WriteNote(newS.LocalPath+".note.txt", []byte{})
			},
			afterAct: func(oldS, newS models.Settings) {
				_ = pkg.Delete(oldS.LocalPath + ".note.txt")
				_ = pkg.Delete(newS.LocalPath + ".note.txt")
				_ = pkg.Delete(newS.LocalPath)
				_ = pkg.Delete(oldS.LocalPath)
			},
			expected: nil,
		},
	}

	for _, td := range tests {
		td.beforeAct(td.localService.Config, td.settings)
		got := td.localService.MoveNotes(td.settings)
		td.afterAct(td.localService.Config, td.settings)

		if (got == nil || td.expected == nil) && got != td.expected ||
			(got != nil && td.expected != nil) && got.Error() != td.expected.Error() {
			t.Errorf("Sum of {error}[MoveNotes] is different: Got: %v | Want: %v", got, td.expected)
		}
	}
}

func TestFetch(t *testing.T) {
	// TODO: Implement tests by mocking.
}

func TestPush(t *testing.T) {
	// TODO: Implement tests by mocking.
}

func TestMigrate(t *testing.T) {
	// TODO: Implement tests by mocking.
}

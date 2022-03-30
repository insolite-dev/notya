package services_test

import (
	"errors"
	"testing"

	"github.com/anonistas/notya/assets"
	"github.com/anonistas/notya/lib/models"
	"github.com/anonistas/notya/lib/services"
	"github.com/anonistas/notya/pkg"
)

// Define a mock local service implementation.
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
		note     models.Note
		expected string
	}{
		{
			note:     models.Note{Path: "../random-path-note.txt"},
			expected: "../random-path-note.txt",
		},
		{
			note:     models.Note{Title: "new-note.txt"},
			expected: ls.Config.LocalPath + "new-note.txt",
		},
	}

	for _, td := range tests {
		got := ls.GeneratePath(td.note)

		if got != td.expected {
			t.Errorf("Sum of [GeneratePath] is different: Got: %v | Want: %v", got, td.expected)
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
		got, err := td.localService.Settings()
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

func TestOpen(t *testing.T) {
	ls := services.LocalService{
		NotyaPath: "./",
		Config:    models.Settings{LocalPath: "./", Editor: "vi"},
		Stdargs:   models.StdArgs{},
	}

	tests := []struct {
		note         models.Note
		localService services.LocalService
		beforeAct    func(note models.Note)
		afterAct     func(note models.Note)
		expected     error
	}{
		{
			note:         models.Note{Title: "somerandomnotethatnotexists"},
			localService: ls,
			beforeAct:    func(note models.Note) {},
			afterAct:     func(note models.Note) {},
			expected:     assets.NotExists("somerandomnotethatnotexists"),
		},
		{
			note:         models.Note{Title: ""},
			localService: ls,
			beforeAct:    func(note models.Note) {},
			afterAct:     func(note models.Note) {},
			expected:     assets.NotExists(""),
		},
		{
			note:         models.Note{Title: "somerandomnote.txt"},
			localService: ls,
			beforeAct: func(note models.Note) {
				path := ls.GeneratePath(note)
				_ = pkg.WriteNote(path, []byte{})
			},
			afterAct: func(note models.Note) {
				path := ls.GeneratePath(note)
				_ = pkg.Delete(path)
			},
			expected: errors.New("signal: abort trap"),
		},
	}

	for _, td := range tests {
		td.beforeAct(td.note)
		got := td.localService.Open(td.note)
		td.afterAct(td.note)

		if (got == nil || td.expected == nil) && got != td.expected ||
			(got != nil && td.expected != nil) && got.Error() != td.expected.Error() {
			t.Errorf("Sum of [Open] is different: Got: %v | Want: %v", got, td.expected)
		}
	}
}

func TestRemove(t *testing.T) {
	ls := services.LocalService{
		NotyaPath: "./",
		Config:    models.Settings{LocalPath: "./", Editor: "vi"},
		Stdargs:   models.StdArgs{},
	}

	tests := []struct {
		note         models.Note
		localService services.LocalService
		beforeAct    func(note models.Note)
		afterAct     func(note models.Note)
		expected     error
	}{
		{
			note:         models.Note{Title: "somerandomnotethatnotexists"},
			localService: ls,
			beforeAct:    func(note models.Note) {},
			afterAct:     func(note models.Note) {},
			expected:     assets.NotExists("somerandomnotethatnotexists"),
		},
		{
			note:         models.Note{Title: ""},
			localService: ls,
			beforeAct:    func(note models.Note) {},
			afterAct:     func(note models.Note) {},
			expected:     assets.NotExists(""),
		},
		{
			note:         models.Note{Title: ".mock-folder"},
			localService: ls,
			beforeAct: func(note models.Note) {
				path := ls.GeneratePath(note)
				_ = pkg.NewFolder(path)
				_ = pkg.WriteNote(path+"/"+"mock_note.txt", []byte{})
			},
			afterAct: func(note models.Note) {
				path := ls.GeneratePath(note)
				_ = pkg.Delete(path + "/" + "mock_note.txt")
				_ = pkg.Delete(path)
			},
			expected: errors.New("remove ./.mock-folder: directory not empty"),
		},
		{
			note:         models.Note{Title: "somerandomnote.txt"},
			localService: ls,
			beforeAct: func(note models.Note) {
				path := ls.GeneratePath(note)
				_ = pkg.WriteNote(path, []byte{})
			},
			afterAct: func(note models.Note) {},
			expected: nil,
		},
	}

	for _, td := range tests {
		td.beforeAct(td.note)
		got := td.localService.Remove(td.note)
		td.afterAct(td.note)

		if (got == nil || td.expected == nil) && got != td.expected ||
			(got != nil && td.expected != nil) && got.Error() != td.expected.Error() {
			t.Errorf("Sum of [Rename] is different: Got: %v | Want: %v", got, td.expected)
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
				path := ls.GeneratePath(note)
				_ = pkg.WriteNote(path, []byte{})
			},
			afterAct: func(note models.Note) {
				path := ls.GeneratePath(note)
				_ = pkg.Delete(path)
			},
			expected: assets.AlreadyExists("somerandomnotethatexists"),
		},
		{
			note:         models.Note{Title: "mocknote.txt"},
			localService: ls,
			beforeAct:    func(note models.Note) {},
			afterAct: func(note models.Note) {
				path := ls.GeneratePath(note)
				_ = pkg.Delete(path)
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
				path := ls.GeneratePath(note)
				_ = pkg.Delete(path)
			},
			afterAct:    func(note models.Note) {},
			expected:    nil,
			expectedErr: assets.NotExists("somerandomnotethatnotexists"),
		},
		{
			note:         models.Note{Title: "mocknote.txt"},
			localService: ls,
			beforeAct: func(note models.Note) {
				path := ls.GeneratePath(note)
				_ = pkg.WriteNote(path, []byte{})
			},
			afterAct: func(note models.Note) {
				path := ls.GeneratePath(note)
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
				path := ls.GeneratePath(note)
				_ = pkg.Delete(path)
			},
			afterAct:    func(note models.Note) {},
			expected:    nil,
			expectedErr: assets.NotExists("somerandomnotethatnotexists"),
		},
		{
			note:         models.Note{Title: "mocknote.txt", Body: "empty-body"},
			localService: ls,
			beforeAct: func(note models.Note) {
				path := ls.GeneratePath(note)
				_ = pkg.WriteNote(path, []byte{})
			},
			afterAct: func(note models.Note) {
				path := ls.GeneratePath(note)
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

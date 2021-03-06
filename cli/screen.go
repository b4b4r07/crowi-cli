package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/b4b4r07/crowi/api" // TODO (spinner)
	"github.com/crowi/go-crowi"
)

var Extention string = ".md"

type Screen struct {
	Text  string
	ID    func(string) string
	Pages *crowi.Pages
}

func NewScreen() (*Screen, error) {
	s := api.NewSpinner("Fetching...")
	s.Start()
	defer s.Stop()

	user := Conf.Crowi.User
	if user == "" {
		return &Screen{}, errors.New("user is not defined")
	}

	client, err := NewClient()
	if err != nil {
		return &Screen{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), api.Timeout)
	defer cancel()

	res, err := client.Pages.List(ctx, "", user, &crowi.PagesListOptions{
		crowi.ListOptions{Pagenation: Conf.Crowi.Paging},
	})
	if err != nil {
		return &Screen{}, err
	}

	if !res.OK {
		return &Screen{}, errors.New(res.Error)
	}

	text := ""
	for _, pi := range res.Pages {
		text += fmt.Sprintf("%s\n", pi.Path)
	}

	id := func(path string) (id string) {
		for _, pi := range res.Pages {
			if pi.Path == path {
				return pi.ID
			}
		}
		return ""
	}

	return &Screen{
		Text:  text,
		ID:    id,
		Pages: res,
	}, nil
}

type Line struct {
	Path      string
	URL       string
	ID        string
	LocalPath string
}

type Lines []Line

func (s *Screen) parseLine(line string) *Line {
	return &Line{
		Path:      line, // for now
		URL:       path.Join(Conf.Crowi.BaseURL, line),
		ID:        s.ID(line),
		LocalPath: filepath.Join(Conf.Crowi.LocalPath, s.ID(line)+Extention),
	}
}

func (s *Screen) Select() (lines Lines, err error) {
	if s.Text == "" {
		err = errors.New("no text to display")
		return
	}
	selectcmd := Conf.Core.SelectCmd
	if selectcmd == "" {
		err = errors.New("no selectcmd specified")
		return
	}

	var buf bytes.Buffer
	err = runFilter(selectcmd, strings.NewReader(s.Text), &buf)
	if err != nil {
		return
	}

	if buf.Len() == 0 {
		err = errors.New("no lines selected")
		return
	}

	selectedLines := strings.Split(buf.String(), "\n")
	for _, line := range selectedLines {
		if line == "" {
			continue
		}
		parsedLine := s.parseLine(line)
		lines = append(lines, *parsedLine)
	}

	if len(lines) == 0 {
		err = errors.New("no lines selected")
		return
	}

	return
}

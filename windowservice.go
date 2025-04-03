package main

import (
	"blive-vup-layer/util"
	"fmt"
)

func (s *Service) ShowWindow(windowID string) {
	defer util.Recover()

	window, ok := s.app.WindowMap[windowID]
	if !ok {
		s.writeResultError(ResultTypeWindow, CodeBadRequest, fmt.Sprintf("window id %s not found", windowID))
		return
	}
	window.Window.Show()
	window.Window.Restore()
	window.Window.Focus()
	s.writeResultOK(ResultTypeWindow, nil)
}

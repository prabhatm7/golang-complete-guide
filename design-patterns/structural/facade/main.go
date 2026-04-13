package main

import "fmt"

/*
SUBSYSTEMS
*/
type TV struct{}

func (t *TV) On() { fmt.Println("TV ON") }

type SoundSystem struct{}

func (s *SoundSystem) On() { fmt.Println("Sound ON") }

type Lights struct{}

func (l *Lights) Dim() { fmt.Println("Lights DIM") }

/*
FACADE
*/
type HomeTheaterFacade struct {
	tv     *TV
	sound  *SoundSystem
	lights *Lights
}

func NewHomeTheater() *HomeTheaterFacade {
	return &HomeTheaterFacade{
		tv:     &TV{},
		sound:  &SoundSystem{},
		lights: &Lights{},
	}
}

func (h *HomeTheaterFacade) WatchMovie() {
	h.lights.Dim()
	h.tv.On()
	h.sound.On()
}

func main() {
	home := NewHomeTheater()
	home.WatchMovie()
}

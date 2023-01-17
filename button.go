package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/nuttech/bell"
)

type Button struct {
	sprite   *ebiten.Image
	position Point
	isActive bool
	pressed  func()
}

// Clear all consts properly
// PressDetect is called on LMB clicked and detects whether
// a click landed on a parent button.
func (b *Button) PressDetect(message bell.Message) {
	pressedX := message.Value.(Point).x / 16
	pressedY := message.Value.(Point).y / 16
	bWidth := 1.0
	bHegith := 1.0
	bPosX := (*b).position.x + bWidth / 2
	bPosY := (*b).position.y + bHegith / 2

	if pressedX < bPosX + bWidth / 2 &&
		pressedX > bPosX - bWidth / 2 &&
		pressedY < bPosY + bHegith / 2 &&
		pressedY >  bPosY - bHegith / 2 {
		(*b).pressed()
	}
}

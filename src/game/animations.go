package game

import (
	"errors"
)

type AnimationType int32
type AnimationStatus int32

const (
	ANIM_CREATE_CELL AnimationType = iota
)

const (
	ANIM_CREATED AnimationStatus = iota
	ANIM_RUNNING
	ANIM_FINISHED
)

type Animation interface {
	GetType() AnimationType
	GetStatus() AnimationStatus
	GetLength() int
	Step() error
}

type CreateAnimation struct {
	animationType   AnimationType
	animationStatus AnimationStatus
	duration        int
	position        Cell
	startingSize    int
	targetSize      int
	currentSize     int
	step            int
}

// impl CreateAnimation

func (animation *CreateAnimation) GetType() AnimationType {
	return ANIM_CREATE_CELL
}

func (animation *CreateAnimation) GetStatus() AnimationStatus {
	return animation.animationStatus
}

func (animation *CreateAnimation) GetLength() int {
	return animation.duration
}

func (animation *CreateAnimation) Step() error {
	if animation.animationStatus == ANIM_FINISHED {
		return errors.New("step called on already finished animation")
	}

	animation.animationStatus = ANIM_RUNNING
	animation.currentSize += animation.step
	if animation.currentSize >= animation.targetSize {
		animation.animationStatus = ANIM_FINISHED
	}

	return nil
}

// end impl CreateAnimation

func CreateCellAnimation(cell Cell, durationInTicks int) *CreateAnimation {
	animation := CreateAnimation{
		animationType:   ANIM_CREATE_CELL,
		animationStatus: ANIM_CREATED,
		duration:        durationInTicks,
		position:        cell,
		startingSize:    CELL_SIZE / 2,
		targetSize:      CELL_SIZE,
	}

	animation.currentSize = animation.startingSize
	animation.step = max((animation.targetSize-animation.startingSize)/durationInTicks, 1)
	return &animation
}

func OnAnimationFinished(animation Animation, g *Game) {
	switch animation.GetType() {
	case ANIM_CREATE_CELL:
		a := animation.(*CreateAnimation)
		cellPos := a.position

		c := &g.board.cells[cellPos.pos_x][cellPos.pos_y]
		c.isRendered = true
		c.val = 2
	}
}

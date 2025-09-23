package game

import (
	"errors"
)

type AnimationType int32
type AnimationStatus int32

const (
	ANIM_CREATE_CELL AnimationType = iota
	ANIM_MOVE_CELL
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
	OnFinish(Animation, *Game)
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

type MoveAnimation struct {
	animationType   AnimationType
	animationStatus AnimationStatus
	duration        int
	start_pos       Vec2
	end_pos         Vec2
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

func (animation *CreateAnimation) OnFinish(anim Animation, g *Game) {
	a := anim.(*CreateAnimation)

	cellPos := a.position
	c := &g.board.cells[cellPos.pos_x][cellPos.pos_y]
	c.isRendered = true
	c.val = 2
}

// end

// impl MoveAnimation

func (animation *MoveAnimation) GetType() AnimationType {
	return ANIM_MOVE_CELL
}

func (animation *MoveAnimation) GetStatus() AnimationStatus {
	return animation.animationStatus
}

func (animation *MoveAnimation) GetLength() int {
	return animation.duration
}

func (animation *MoveAnimation) Step() error {
	// TODO
	return nil
}

func (animation *MoveAnimation) OnFinish() {

}

// end

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

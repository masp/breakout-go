package	main

import (
	"fmt"
	"math"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/go-gl/mathgl/mgl32"
)

var LEFT = -1.0;
var RIGHT = 1.0;
type Paddle struct {
	Rectangle
	paddleSpeed float64
}

type Collision struct {
	delta mgl64.Vec2
	normal mgl64.Vec2
}

func (p *Paddle) move(dir float64, dt float64) {
	p.pos = mgl64.Vec2{p.pos.X() + dir * dt * p.paddleSpeed, p.pos.Y()}
}

type Ball struct {
	Rectangle
	vel mgl64.Vec2
}

func deflect(b, n mgl64.Vec2) mgl64.Vec2 {
	// r=d−2(d⋅n)n
	return b.Sub(n.Mul(2*(b.Dot(n))))
}

func (b *Ball) update(block *Block, paddle Paddle, dt float64) bool {
	newPos := mgl64.Vec2{b.pos.X() + b.vel.X() * dt, b.pos.Y() + b.vel.Y() * dt}
	newBall := Rectangle{newPos, b.width, b.height, b.color}
	collision := Collision{}
	updateBlock := false
	for i, box := range block.boxes {
		if box.getCollision(newBall).normal.Len() > 0 {
			collision = box.getCollision(newBall)
			block.boxes[i] = block.boxes[len(block.boxes)-1]
			block.boxes = block.boxes[:len(block.boxes)-1]
			updateBlock = true
		}
	}
	if paddle.getCollision(newBall).normal.Len() > 0 {
		collision = paddle.getCollision(newBall);
	}

	// Check walls
	if newPos.X() < 0 {
		collision = Collision{mgl64.Vec2{-newPos.X(), 0}, mgl64.Vec2{1, 0}}
	} else if newPos.Y() < 0 {
		collision = Collision{mgl64.Vec2{0, -newPos.Y()}, mgl64.Vec2{0, 1}}
	} else if newPos.X() + b.width > width {
		collision = Collision{mgl64.Vec2{width - (newPos.X() + b.width), 0}, mgl64.Vec2{-1, 0}}
	} else if newPos.Y() + b.height > height {
		collision = Collision{mgl64.Vec2{0, height - (newPos.Y() + b.height)}, mgl64.Vec2{0, -1}}
	}

	if collision.normal.Len() > 0 {
		b.pos = newPos.Add(collision.delta);
		b.vel = deflect(b.vel, collision.normal);
	} else {
		b.pos = newPos
	}
	return updateBlock

}

type Rectangle struct {
	pos mgl64.Vec2
	width, height float64
	color mgl64.Vec4
}

func (r Rectangle) getCollision(other Rectangle) Collision {
	dx := r.center().X() - other.center().X()
	px := (r.width / 2 + other.width / 2) - math.Abs(dx)
	dy := r.center().Y() - other.center().Y()
	py := (r.height / 2 + other.height / 2) - math.Abs(dy)

	signX := math.Copysign(1, -dx);
	signY := math.Copysign(1, -dy);

	if px <= 0 || py <= 0 {
		return Collision{};
	}

	if px < py {
		return Collision{mgl64.Vec2{px * signX, 0}, mgl64.Vec2{signX, 0}}
	} else {
		return Collision{mgl64.Vec2{0, py * signY}, mgl64.Vec2{0, signY}}
	}
	return Collision{};
}

func (r Rectangle) center() mgl64.Vec2 {
	return mgl64.Vec2{r.pos.X() + r.width / 2, r.pos.Y() + r.height / 2}
}

func (r Rectangle) getVerts() ([]float32) {
	bottomLeftCorner := []float32 {
		0,       0,
		0,       float32(r.height),
		float32(r.width), 0,
	}

	upperRightCorner := []float32 {
		0,       float32(r.height),
		float32(r.width), float32(r.height),
		float32(r.width), 0,
	}
	return append(bottomLeftCorner, upperRightCorner...)
}

func (r Rectangle) getModelMatrix() mgl32.Mat4 {
	return mgl32.Translate3D(float32(r.pos.X()), float32(r.pos.Y()), 0)
}

func (r Rectangle) getColor() mgl32.Vec4 {
	return mgl32.Vec4{float32(r.color.X()), float32(r.color.Y()), float32(r.color.Z()), float32(r.color.W())}
}

func (r Rectangle) String() string {
	return fmt.Sprintf("%v, %v", r.pos.X(), r.pos.Y())
}

type Block struct {
	boxes []Rectangle
}

func colorFor(x, y int) mgl64.Vec4 {
	var r, g, b float64;
	r = 1.0
	if x % 2 == 0 {
		r = 0.5
	}
	if y % 2 == 0 {
		g = 1.0
	}
	if y % 2 == 0 && x % 2 == 0 {
		b = 1.0
	}
	return mgl64.Vec4{ r, g, b, 1.0 }
}

func buildMap() Block {
	rects := make([]Rectangle, 0)
	startX := 50.0
	startY := 300.0
	for x := 0; x < 5; x++ {
		for y := 0; y < 3; y++ {
			rects = append(rects, Rectangle{
				mgl64.Vec2{ startX + float64(x) * 105.0, startY + float64(y) * 50.0 },
				100, 45, colorFor(x, y)})
		}
	}

	return Block{rects}

}

func (b *Block) getVerts() ([]float32) {
	sumVerts := make([]float32, 0)
	for _, block := range b.boxes {
		sumVerts = append(sumVerts, block.getVerts()...)
	}
	return sumVerts
}


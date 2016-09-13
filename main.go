package main

import (
	"runtime"
	"fmt"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl64"
)

const width = 640
const height = 480

var paddle = Paddle{Rectangle{mgl64.Vec2{0, 100}, 200, 50, mgl64.Vec4{1.0, 0.0, 0.0, 1.0}}, 300.0}
var ball = Ball{Rectangle{mgl64.Vec2{width/2, height/2}, 10, 10, mgl64.Vec4{0.5, 0.5, 0.5, 1.0}}, mgl64.Vec2{-100.0, 200.0}}

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func onKey(w *glfw.Window, key glfw.Key, scancode int,
		action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Press {
		w.SetShouldClose(true)
	}
}

func AddObjs(r *Renderer, paddle *Paddle, ball *Ball, block *Block) {
	r.AddObj(paddle)
	r.AddObj(ball)
	for _, box := range block.boxes {
		r.AddObj(box)
	}
}

func main() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(width, height, "Breakout", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	window.SetKeyCallback(onKey)
	block := buildMap()
	renderer := Renderer{}
	AddObjs(&renderer, &paddle, &ball, &block)
	renderer.Init()

	glerr := gl.GetError()
	if glerr != gl.NO_ERROR {
		fmt.Printf("Error: %x", glerr)
	}

	previousTime := glfw.GetTime()
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		time := glfw.GetTime()
		elapsed := time - previousTime
		previousTime = time

		renderer.Render()

		if window.GetKey(glfw.KeyD) == glfw.Press {
			paddle.move(RIGHT, elapsed)
		} else if window.GetKey(glfw.KeyA) == glfw.Press {
			paddle.move(LEFT, elapsed)
		}

		shouldUpdateBoxes := ball.update(&block, paddle, elapsed)
		if (shouldUpdateBoxes) {
			renderer.objects = []Renderable{}
			AddObjs(&renderer, &paddle, &ball, &block);
			renderer.createObjects();
		}

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

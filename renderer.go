package main

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/gl/v3.3-core/gl"
	"strings"
	"fmt"
)

var vertexShader = `
#version 330
in vec2 position;
uniform mat4 projection;
uniform mat4 model;

void main()
{
	gl_Position = projection * model * vec4(position, 0.0, 1.0);
}
` + "\x00"

var fragmentShader = `
#version 330

uniform vec4 color = vec4(1.0, 0.0, 0.0, 1.0);
out vec4 fragColor;

void main()
{
	fragColor = color;
}
` + "\x00"

type Renderable interface {
	getVerts() []float32
	getModelMatrix() mgl32.Mat4
	getColor() mgl32.Vec4
}

type Renderer struct {
	objects []Renderable
	vao, vbo, prog, vshader, fshader, colorUniform, modelUniform uint32
	initd bool
}

func (r *Renderer) Init() {
	var err error
	r.vshader, err = createShader(gl.VERTEX_SHADER, vertexShader)
	if err != nil {
		panic(err)
	}
	r.fshader, err = createShader(gl.FRAGMENT_SHADER, fragmentShader)
	if err != nil {
		panic(err)
	}
	r.prog, err = createProgram(r.vshader, r.fshader)
	if err != nil {
		panic(err)
	}

	gl.GenVertexArrays(1, &r.vao)
	r.createObjects()
}

func (r *Renderer) AddObj(obj Renderable) {
	r.objects = append(r.objects, obj)
}

func (r *Renderer) Render() {
	gl.UseProgram(r.prog)
	gl.BindVertexArray(r.vao)
	projection := mgl32.Ortho2D(0, width, 0, height)
	projectionUniform := gl.GetUniformLocation(r.prog, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	modelUniform := gl.GetUniformLocation(r.prog, gl.Str("model\x00"))
	colorUniform := gl.GetUniformLocation(r.prog, gl.Str("color\x00"))
	var currIndx int32 = 0
	for _, obj := range r.objects {
		modelMat := obj.getModelMatrix()
		colorVec := obj.getColor()
		gl.UniformMatrix4fv(modelUniform, 1, false, &modelMat[0])
		gl.ProgramUniform4fv(r.prog, colorUniform, 1, &colorVec[0])
		verticesLength := int32(len(obj.getVerts())) / 2
		gl.DrawArrays(gl.TRIANGLES, currIndx, verticesLength)
		currIndx += verticesLength
	}
}

func (r *Renderer) createObjects() {
	vertices := r.allVertices()
	gl.BindVertexArray(r.vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices) * 4, gl.Ptr(vertices), gl.STATIC_DRAW)
	if r.vbo != 0 {
		gl.DeleteBuffers(1, &r.vbo)
	}
	r.vbo = vbo
	vertAttrib := uint32(gl.GetAttribLocation(r.prog, gl.Str("position\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))
}

func (r *Renderer) allVertices() []float32 {
	verts := make([]float32, 0)
	for _, obj := range r.objects {
		verts = append(verts, obj.getVerts()...)
	}
	return verts
}

func createShader(stype uint32, shaderText string) (uint32, error) {
	shader := gl.CreateShader(stype)
	csources, free := gl.Strs(shaderText)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var compiled, maxLength int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &compiled)
	if compiled == gl.FALSE {
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &maxLength)
		log := strings.Repeat("\x00", int(maxLength + 1))
		gl.GetShaderInfoLog(shader, maxLength, &maxLength, gl.Str(log))
		return 0, fmt.Errorf("failed to compile %v: %v", shaderText, log)
	}
	return shader, nil
}

func createProgram(vs uint32, fs uint32) (uint32, error) {
	prog := gl.CreateProgram()
	gl.AttachShader(prog, vs)
	gl.AttachShader(prog, fs)

	gl.LinkProgram(prog)
	var linked, maxLength int32
	gl.GetProgramiv(prog, gl.LINK_STATUS, &linked)
	if linked == gl.FALSE {
		gl.GetProgramiv(prog, gl.INFO_LOG_LENGTH, &maxLength)
		log := strings.Repeat("\x00", int(maxLength + 1))
		gl.GetProgramInfoLog(prog, maxLength, &maxLength, gl.Str(log))
		return 0, fmt.Errorf("failed to link: %v", log)
	}

	return prog, nil
}



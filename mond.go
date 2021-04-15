package main

import (
"fmt"
"golang.org/x/image/colornames"
"image"
_ "image/png"
"math"
"math/rand"
"os"
"time"

"github.com/faiface/pixel"
"github.com/faiface/pixel/pixelgl"
//"golang.org/x/image/colornames"
)

// Gravity Physics Stuff
const N  = 1500
const G = 0.1
const eps = 0.1
const metric = 800
const zoom = 0.4
const winHeight = 800
const winWidth = 800

// Vector Definition
type Vector struct {
	x, y float64
}
// Vector Methods
func (this Vector) Neg() Vector{
	return Vector{this.x*-1, this.y*-1}
}
func (this Vector) Norm() float64{
	return math.Sqrt(this.x*this.x + this.y*this.y)
}
func (this Vector) Add(that Vector) Vector{
	return Vector{this.x + that.x,this.y + that.y}
}
func (this Vector) Scale(C float64) Vector{
	return Vector{this.x*C,this.y*C}
}

//Body Struct
type Body struct{
	mass float64
	pos Vector
	vel Vector
	acc Vector
	force Vector
}

// Body Methods
func (this Body) getGravity(Bodies [N]Body) Vector{
	this.acc  = Vector{0,0}
	for i := 0; i<N; i++ {
		var r = Bodies[i].pos.Add(this.pos.Neg()).Norm()
		if r == 0 {
			this.acc.Add(Vector{0, 0})
		}
	}
	if this.pos.Norm() > 800{
		this.acc = this.acc.Neg().Scale(100)
	}
	return this.acc
}
func (this Body) getVel() Vector{
	this.vel.x += this.acc.x*eps
	this.vel.y += this.acc.y*eps
	return this.vel
}
func (this Body) getPos() Vector{
	this.pos.x += this.vel.x*eps
	this.pos.y += this.vel.y*eps
	return this.pos
}

// Loading the image given the path
func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

// Running and Setip of the
func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Giri's Gravity Simulations!",
		Bounds: pixel.R(0, 0, winHeight, winWidth),
		//VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	//win.SetSmooth(true)
	pic, err := loadPicture("hiking.png")
	if err != nil {
		panic(err)
	}

	var bodyFrames [N]*pixel.Sprite
	for i:=0;i<N;i++{
		bodyFrames[i] = pixel.NewSprite(pic, pic.Bounds())
	}

	// Camera Setup
	var(
		camPos = Vector{0,0}
		camSpeed float64 = 600
		camZoomSpeed = 1.2
		camZoom = 1.0
	)

	// Body Initialisation
	var Bodies [N]Body
	for i := 0; i<N; i++ {
		c := 0.0
		Bodies[i] = Body{rand.Float64()*30 + 10, Vector{rand.Float64()*metric-metric/2, rand.Float64()*metric-metric/2}, Vector{c,c}, Vector{0,0}, Vector{0,0}}
	}

	// Fps Setup
	var (
		frames = 0
		second = time.Tick(time.Second)
	)

	// Computing of Trajectories
	last := time.Now()

	for !win.Closed() {
		win.Clear(colornames.Midnightblue)

		dt := time.Since(last).Seconds()
		last = time.Now()

		// Keypress Activities
		if win.Pressed(pixelgl.KeyLeft) {
			camPos.x += camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyRight) {
			camPos.x -= camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyUp) {
			camPos.y -= camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyDown) {
			camPos.y += camSpeed * dt
		}
		camZoom *= math.Pow(camZoomSpeed, win.MouseScroll().Y)

		// Computation of Trajectories
		for i:=0;i<N;i++ {
			Bodies[i].acc = Bodies[i].getGravity(Bodies)
			Bodies[i].vel = Bodies[i].getVel()
			Bodies[i].pos = Bodies[i].getPos()
			var x, y = Bodies[i].pos.x+winWidth/2, Bodies[i].pos.y+winHeight/2
			mat := pixel.IM
			mat = mat.Moved(pixel.V(x,y)).Scaled(pixel.V(x,y),zoom).Scaled(pixel.V(x,y),Bodies[i].mass/25).Moved(pixel.V(camPos.x,camPos.y)).Scaled(pixel.V(winWidth/2,winHeight/2),camZoom)
			bodyFrames[i].Draw(win, mat)
		}
		win.Update()

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}
	}
}

// Finally running the program
func main() {
	pixelgl.Run(run)
}


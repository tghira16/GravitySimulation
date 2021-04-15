package main

import "math"
var (
	N  = 2
	G = 100
	eps = 0.01
	)


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
		if r == 0{
			this.acc.Add(Vector{0,0})
		}else if r<this.mass{
			this.acc = ((Bodies[i].pos.Add(this.pos.Neg()).Scale((G*Bodies[i].mass)/(r*r*r))).Add(this.acc))
		}else{
			this.acc = ((Bodies[i].pos.Add(this.pos.Neg()).Scale((G*Bodies[i].mass)/(r*r*r))).Add(this.acc))
		}
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
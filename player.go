// Copyright 2022 Anıl Konaç

package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/jakecoffman/cp"
	"golang.org/x/image/colornames"
)

const (
	playerWidth  = 40.0
	playerHeight = 100.0
	gunWidth     = playerHeight / 2.0
	gunHeight    = playerWidth / 3.0
)

const (
	gravity = 500.0
	// Taken from cp-examples/player and modified
	playerVelocity        = 400.0
	playerGroundAccelTime = 0.1
	playerGroundAccel     = playerVelocity / playerGroundAccelTime
	playerAirAccelTime    = 0.5
	playerAirAccel        = playerVelocity / playerAirAccelTime
	jumpHeight            = 100.0
	//
)

var (
	imagePlayer = ebiten.NewImage(1, 1)
	imageGun    = ebiten.NewImage(1, 1)
)

func init() {
	imagePlayer.Fill(colornames.Slategray)
	imageGun.Fill(colornames.Orange)
}

type player struct {
	pos               cp.Vector
	posGun            cp.Vector
	angleGun          float64
	shape             *cp.Shape
	body              *cp.Body
	drawOptionsPlayer ebiten.DrawImageOptions
	drawOptionsGun    ebiten.DrawImageOptions
	onGround          bool
}

func newPlayer(pos cp.Vector) *player {
	player := &player{
		pos: pos,
	}

	player.body = cp.NewBody(1, cp.INFINITY)
	player.body.SetPosition(pos)
	player.body.SetVelocityUpdateFunc(playerUpdateVelocity)
	player.shape = cp.NewBox(player.body, playerWidth, playerHeight, playerHeight/2.0)
	// player.shape.SetElasticity(0.1)

	return player
}

func (p *player) update(input *input) {
	// Update position
	p.pos = p.body.Position()

	// Update gun position
	p.posGun = p.pos.Add(cp.Vector{X: playerWidth / 2.0, Y: playerHeight / 2.0})

	// Update gun angle
	distX := input.cursorPos.X - p.posGun.X
	distY := input.cursorPos.Y - p.posGun.Y
	p.angleGun = math.Atan2(distY, distX)

	// Handle inputs
	var surfaceV cp.Vector
	if input.right {
		surfaceV.X = -playerVelocity
	} else if input.left {
		surfaceV.X = playerVelocity
	}
	p.shape.SetSurfaceV(surfaceV)
	if p.onGround { // TODO
		p.shape.SetFriction(playerGroundAccel / gravity) // Taken from cp-examples/player
	} else {
		p.shape.SetFriction(0)
	}

	// Grab the grounding normal from last frame - Taken from cp-examples/player
	groundNormal := cp.Vector{}
	p.body.EachArbiter(func(arb *cp.Arbiter) {
		n := arb.Normal() //.Neg()

		if n.Y > groundNormal.Y {
			groundNormal = n
		}
	})
	p.onGround = groundNormal.Y > 0

	if input.up && p.onGround {
		jumpV := math.Sqrt(2.0 * jumpHeight * gravity) // Taken from cp-examples/player
		p.body.SetVelocityVector(p.body.Velocity().Add(cp.Vector{X: 0, Y: -jumpV}))
	}
	// Apply air control if not on ground
	if !p.onGround {
		// p.body.SetVelocity(cp.LerpConst(v.X, -surfaceV.X, playerAirAccel*deltaTime), v.Y)
		p.body.SetVelocityVector(p.body.Velocity().Add(cp.Vector{X: -surfaceV.X * deltaTime, Y: 0}))
	}

	// v := p.body.Velocity()
	// fmt.Printf("Friction: %.2f\tVel X: %.2f\tVel Y: %.2f\n", p.shape.Friction(), v.X, v.Y)
	p.updateGeometryMatrices()
}

func (p *player) updateGeometryMatrices() {
	// Player
	p.drawOptionsPlayer.GeoM.Reset()
	p.drawOptionsPlayer.GeoM.Scale(playerWidth, playerHeight)
	p.drawOptionsPlayer.GeoM.Translate(p.pos.X, p.pos.Y)

	// Gun
	p.drawOptionsGun.GeoM.Reset()
	p.drawOptionsGun.GeoM.Scale(gunWidth, gunHeight)
	p.drawOptionsGun.GeoM.Translate(0, -gunHeight/2.0)
	p.drawOptionsGun.GeoM.Rotate(p.angleGun)
	p.drawOptionsGun.GeoM.Translate(p.posGun.X, p.posGun.Y)
}

func (p *player) draw(dst *ebiten.Image) {
	// Draw prototype player
	dst.DrawImage(imagePlayer, &p.drawOptionsPlayer)

	// Draw prototype gun
	dst.DrawImage(imageGun, &p.drawOptionsGun)
}

func playerUpdateVelocity(body *cp.Body, gravity cp.Vector, damping, dt float64) {
	body.UpdateVelocity(gravity, damping, dt)
}

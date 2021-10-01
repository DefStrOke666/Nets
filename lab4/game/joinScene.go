package game

import (
	"github.com/borodun/nsu-nets/lab4/game/proto"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"image/color"
	"strconv"
)

var (
	serverHeight = 50
)

type JoinScene struct {
	backgroundPics []*Picture

	buttonPics    []*Picture
	exitButtonPic *Picture
	background    *ebiten.Image

	servers []proto.GameMessage_AnnouncementMsg

	serverImg []*Picture
	infoImg   []*Picture

	selectedServer int
	canJoin        bool
}

func NewJoinScene() *JoinScene {
	scene := &JoinScene{}

	scene.backgroundPics = make([]*Picture, 3)
	scene.buttonPics = make([]*Picture, 2)
	scene.servers = generateServers(3)
	scene.canJoin = false

	scene.updateImages()

	return scene
}

func (j *JoinScene) createServerImage(w, h int, msg *proto.GameMessage_AnnouncementMsg, backgroundClr color.Color, textClr color.Color) *ebiten.Image {
	img := ebiten.NewImage(w, h)
	img.Fill(backgroundClr)

	println("Player count: ", strconv.Itoa(len(msg.Players.Players)))
	playerCountImg := createStringImage(strconv.Itoa(len(msg.Players.Players)), getMenuFonts(3), textClr)
	plX, plY := playerCountImg.Size()

	str := "Can't join"
	if *msg.CanJoin {
		str = "Can join"
	}
	canJoinImg := createStringImage(str, getMenuFonts(3), textClr)

	op := &ebiten.DrawImageOptions{}
	margin := int(margin)
	op.GeoM.Translate(float64(margin), float64((h-plY)/2))
	img.DrawImage(playerCountImg, op)
	op.GeoM.Translate(float64(margin+plX), 0)
	img.DrawImage(canJoinImg, op)
	return img
}

func (j *JoinScene) updateServersPictures(w, x, y int) {
	j.serverImg = nil
	for i := range j.servers {
		selected := i
		p := NewPicture(
			j.createServerImage(w, serverHeight, &j.servers[i], serverBackgroundIdleColor, serverTextIdleColor),
			j.createServerImage(w, serverHeight, &j.servers[i], serverBackgroundActiveColor, serverTextActiveColor),
		)
		p.SetRect(p.GetIdleImage().Bounds().Add(image.Point{X: x, Y: y + i*serverHeight}))
		p.SetHandler(func(state *GameState) {
			j.selectedServer = selected
			j.canJoin = *j.servers[selected].CanJoin
		})
		j.serverImg = append(j.serverImg, p)
	}
}

func (j *JoinScene) createInfoImage(w, h, idx int, msg *proto.GameMessage_AnnouncementMsg, backgroundClr color.Color, textClr color.Color) *ebiten.Image {
	img := ebiten.NewImage(w, h)
	img.Fill(backgroundClr)

	playerCountImg := createStringImage(strconv.Itoa(idx), getMenuFonts(3), textClr)
	plX, plY := playerCountImg.Size()

	str := "Can't join"
	if *msg.CanJoin {
		str = "Can join"
	}
	canJoinImg := createStringImage(str, getMenuFonts(3), textClr)

	op := &ebiten.DrawImageOptions{}
	margin := int(margin)
	op.GeoM.Translate(float64(margin), float64((h-plY)/2))
	img.DrawImage(playerCountImg, op)
	op.GeoM.Translate(float64(margin+plX), 0)
	img.DrawImage(canJoinImg, op)
	return img
}

func (j *JoinScene) updateInfoPictures(w, h, x, y int) {
	j.infoImg = nil
	for i := range j.servers {
		p := NewPicture(
			j.createInfoImage(w, h, i, &j.servers[i], serverBackgroundIdleColor, serverTextIdleColor),
			j.createInfoImage(w, h, i, &j.servers[i], serverBackgroundActiveColor, serverTextActiveColor),
		)
		p.SetRect(p.GetIdleImage().Bounds().Add(image.Point{X: x, Y: y}))
		j.infoImg = append(j.infoImg, p)
	}
}

func (j *JoinScene) updateImages() {
	margin := int(margin)
	spacingsV := margin * 3
	spacingsH := margin * 3

	titleH := textHeight("Servers", getMenuFonts(8)) + margin

	widthUnit := (screenWidth - spacingsH) / 10
	heightUnit := (screenHeight - titleH - spacingsV) / 6

	servListW := widthUnit * 7
	servListH := heightUnit * 5

	infoW := widthUnit * 3
	infoH := servListH

	j.background = getRoundRect(screenWidth, screenHeight, backgroundColor)
	j.backgroundPics[0] = NewPicture(
		createStringImage("Servers", getMenuFonts(8), titleIdleColor),
		createStringImage("Servers", getMenuFonts(8), titleActiveColor))
	j.backgroundPics[1] = NewPicture(
		getRoundRectWithBorder(servListW, servListH, centreIdleColor, lineIdleColor),
		getRoundRectWithBorder(servListW, servListH, centreActiveColor, lineActiveColor))
	j.backgroundPics[2] = NewPicture(
		getRoundRectWithBorder(infoW, infoH, centreIdleColor, lineIdleColor),
		getRoundRectWithBorder(infoW, infoH, centreActiveColor, lineActiveColor))

	j.backgroundPics[0].SetRect(j.backgroundPics[0].GetIdleImage().Bounds().Add(image.Point{X: margin, Y: margin}))
	j.backgroundPics[1].SetRect(j.backgroundPics[1].GetIdleImage().Bounds().Add(image.Point{X: margin, Y: titleH + margin}))
	j.backgroundPics[2].SetRect(j.backgroundPics[2].GetIdleImage().Bounds().Add(image.Point{X: margin*2 + servListW, Y: titleH + margin}))

	buttonW := widthUnit * 3
	buttonH := heightUnit

	j.buttonPics[0] = NewPicture(
		getRoundRectWithBorder(buttonW, buttonH, centreIdleColor, lineIdleColor),
		getRoundRectWithBorder(buttonW, buttonH, centreActiveColor, lineActiveColor))
	j.buttonPics[1] = NewPicture(
		getRoundRectWithBorder(buttonW, buttonH, centreIdleColor, lineIdleColor),
		getRoundRectWithBorder(buttonW, buttonH, centreActiveColor, lineActiveColor))
	j.exitButtonPic = NewPicture(
		getRoundRectWithBorder(buttonW, buttonH, centreIdleColor, lineIdleColor),
		getRoundRectWithBorder(buttonW, buttonH, centreActiveColor, lineActiveColor),
	).SetHandler(func(state *GameState) {
		state.SceneManager.GoTo(NewTitleScene())
	})

	j.buttonPics[0].SetRect(j.buttonPics[0].GetIdleImage().Bounds().Add(image.Point{X: margin, Y: titleH + margin*2 + servListH}))
	j.buttonPics[1].SetRect(j.buttonPics[1].GetIdleImage().Bounds().Add(image.Point{X: margin*2 + buttonW, Y: titleH + margin*2 + servListH}))
	j.exitButtonPic.SetRect(j.exitButtonPic.GetIdleImage().Bounds().Add(image.Point{X: screenWidth - margin - buttonW, Y: titleH + margin*2 + servListH}))

	j.updateServersPictures(servListW-int(lineThickness*2), margin+int(lineThickness), margin+int(radius)+titleH)
	j.updateInfoPictures(infoW-int(lineThickness*2), infoH-int(radius*2), margin*2+servListW+int(lineThickness), margin+int(radius)+titleH)
}

func (j *JoinScene) Update(state *GameState) error {
	if sizeChanged {
		j.updateImages()
	}

	for i := range j.buttonPics {
		if j.canJoin {
			j.buttonPics[i].Update(state)
		}
	}
	j.exitButtonPic.Update(state)

	for i := range j.backgroundPics {
		j.backgroundPics[i].Update(state)
	}

	for i := range j.serverImg {
		j.serverImg[i].Update(state)
	}

	for i := range j.infoImg {
		j.infoImg[i].Update(state)
	}

	return nil
}

func (j *JoinScene) Draw(screen *ebiten.Image) {
	screen.Fill(fillColor)
	screen.DrawImage(j.background, &ebiten.DrawImageOptions{})

	for i := range j.backgroundPics {
		j.backgroundPics[i].Draw(screen)
	}

	for i := range j.serverImg {
		j.serverImg[i].Draw(screen)
	}

	for i := range j.buttonPics {
		j.buttonPics[i].Draw(screen)
	}
	j.exitButtonPic.Draw(screen)

	j.infoImg[j.selectedServer].Draw(screen)

	println("Selected server: ", j.selectedServer)
}

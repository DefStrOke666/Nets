package snakes

import (
	"github.com/borodun/nsu-nets/lab4/snakes/proto"
	"github.com/borodun/nsu-nets/lab4/snakes/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"image/color"
	"strconv"
	"strings"
)

var (
	serverHeight = 50
)

const (
	maxDatagramSize = 8192
)

type JoinScene struct {
	backgroundPics []*utils.Picture

	buttonPics    []*utils.Picture
	exitButtonPic *utils.Picture
	background    *ebiten.Image

	servers []*proto.GameMessage_AnnouncementMsg

	serverImg []*utils.Picture
	infoImg   []*utils.Picture

	selectedServer int
	canJoin        bool

	serversUpdated bool
	exit           bool
}

func NewJoinScene() *JoinScene {
	scene := &JoinScene{}

	scene.backgroundPics = make([]*utils.Picture, 3)
	scene.buttonPics = make([]*utils.Picture, 2)
	scene.canJoin = false

	go scene.receiveAnnouncements()

	scene.updateImages()
	return scene
}

func (j *JoinScene) createServerImage(w, h int, msg *proto.GameMessage_AnnouncementMsg, backgroundClr color.Color, textClr color.Color) *ebiten.Image {
	img := ebiten.NewImage(w, h)
	img.Fill(backgroundClr)

	playerCountImg := utils.CreateStringImage(strconv.Itoa(len(msg.Players.Players)), utils.GetMenuFonts(3), textClr)
	plX, plY := playerCountImg.Size()

	str := "Can't join"
	if *msg.CanJoin {
		str = "Can join"
	}
	str2 := str + "  " + strconv.Itoa(len(msg.Players.Players)) + " players"
	canJoinImg := utils.CreateStringImage(str2, utils.GetMenuFonts(3), textClr)

	op := &ebiten.DrawImageOptions{}
	Margin := int(utils.Margin)
	op.GeoM.Translate(float64(Margin), float64((h-plY)/2))
	img.DrawImage(playerCountImg, op)
	op.GeoM.Translate(float64(Margin+plX), 0)
	img.DrawImage(canJoinImg, op)
	return img
}

func (j *JoinScene) updateServersPictures(w, h, x, y int) {
	j.serverImg = nil
	for i := range j.servers {
		if i*serverHeight > h {
			break
		}
		selected := i
		p := utils.NewPicture(
			j.createServerImage(w, serverHeight, j.servers[i], utils.ServerBackgroundIdleColor, utils.ServerTextIdleColor),
			j.createServerImage(w, serverHeight, j.servers[i], utils.ServerBackgroundActiveColor, utils.ServerTextActiveColor),
		)
		p.SetRect(p.GetIdleImage().Bounds().Add(image.Point{X: x, Y: y + i*serverHeight}))
		p.SetHandler(func() {
			j.selectedServer = selected
			j.canJoin = *j.servers[selected].CanJoin
		})
		j.serverImg = append(j.serverImg, p)
	}
}

func (j *JoinScene) createInfoImage(w, h int, msg *proto.GameMessage_AnnouncementMsg, backgroundClr color.Color, textClr color.Color) *ebiten.Image {
	img := ebiten.NewImage(w, h)
	img.Fill(backgroundClr)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(utils.Margin, utils.Margin)

	configStr := strings.Split(msg.Config.String(), ",")
	for _, s := range configStr {
		if s != "" {
			textH := utils.TextHeight(s, utils.GetMenuFonts(3))
			img.DrawImage(utils.CreateStringImage(s, utils.GetMenuFonts(3), textClr), op)
			op.GeoM.Translate(0, float64(textH)+utils.Margin)
		}
	}

	return img
}

func (j *JoinScene) updateInfoPictures(w, h, x, y int) {
	j.infoImg = nil
	for i := range j.servers {
		p := utils.NewPicture(
			j.createInfoImage(w, h, j.servers[i], utils.ServerBackgroundIdleColor, utils.ServerTextIdleColor),
			j.createInfoImage(w, h, j.servers[i], utils.ServerBackgroundActiveColor, utils.ServerTextActiveColor),
		)
		p.SetRect(p.GetIdleImage().Bounds().Add(image.Point{X: x, Y: y}))
		j.infoImg = append(j.infoImg, p)
	}
}

func (j *JoinScene) updateImages() {
	Margin := int(utils.Margin)
	spacingsV := Margin * 3
	spacingsH := Margin * 3

	titleH := utils.TextHeight("Servers", utils.GetMenuFonts(8)) + Margin

	widthUnit := (screenWidth - spacingsH) / 10
	heightUnit := (screenHeight - titleH - spacingsV) / 6

	servListW := widthUnit * 7
	servListH := heightUnit * 5

	infoW := widthUnit * 3
	infoH := servListH

	j.background = utils.GetRoundRect(screenWidth, screenHeight, utils.BackgroundColor)
	j.backgroundPics[0] = utils.NewPicture(
		utils.CreateStringImage("Servers", utils.GetMenuFonts(8), utils.TitleIdleColor),
		utils.CreateStringImage("Servers", utils.GetMenuFonts(8), utils.TitleActiveColor))
	j.backgroundPics[1] = utils.NewPicture(
		utils.GetRoundRectWithBorder(servListW, servListH, utils.CentreIdleColor, utils.LineIdleColor),
		utils.GetRoundRectWithBorder(servListW, servListH, utils.CentreActiveColor, utils.LineActiveColor))
	j.backgroundPics[2] = utils.NewPicture(
		utils.GetRoundRectWithBorder(infoW, infoH, utils.CentreIdleColor, utils.LineIdleColor),
		utils.GetRoundRectWithBorder(infoW, infoH, utils.CentreActiveColor, utils.LineActiveColor))

	j.backgroundPics[0].SetRect(j.backgroundPics[0].GetIdleImage().Bounds().Add(image.Point{X: Margin, Y: Margin}))
	j.backgroundPics[1].SetRect(j.backgroundPics[1].GetIdleImage().Bounds().Add(image.Point{X: Margin, Y: titleH + Margin}))
	j.backgroundPics[2].SetRect(j.backgroundPics[2].GetIdleImage().Bounds().Add(image.Point{X: Margin*2 + servListW, Y: titleH + Margin}))

	buttonW := widthUnit * 3
	buttonH := heightUnit

	j.buttonPics[0] = utils.NewPicture(
		utils.BorderedRoundRectWithText(buttonW, buttonH, utils.CentreIdleColor, utils.LineIdleColor, "Join", utils.GetMenuFonts(4)),
		utils.BorderedRoundRectWithText(buttonW, buttonH, utils.CentreActiveColor, utils.LineActiveColor, "Join", utils.GetMenuFonts(4)),
	).SetHandler(func() {
		admin := j.servers[j.selectedServer].Players.Players[0]
		addr := admin.GetIpAddress() + ":" + strconv.Itoa(int(admin.GetPort()))
		conf := j.servers[j.selectedServer].Config
		j.exit = true
		sceneManager.GoTo(NewGameScene(conf, addr, false))
	})
	j.buttonPics[1] = utils.NewPicture(
		utils.BorderedRoundRectWithText(buttonW, buttonH, utils.CentreIdleColor, utils.LineIdleColor, "View", utils.GetMenuFonts(4)),
		utils.BorderedRoundRectWithText(buttonW, buttonH, utils.CentreActiveColor, utils.LineActiveColor, "View", utils.GetMenuFonts(4)),
	).SetHandler(func() {
		admin := j.servers[j.selectedServer].Players.Players[0]
		addr := admin.GetIpAddress() + ":" + strconv.Itoa(int(admin.GetPort()))
		conf := j.servers[j.selectedServer].Config
		j.exit = true
		sceneManager.GoTo(NewGameScene(conf, addr, true))
	})
	j.exitButtonPic = utils.NewPicture(
		utils.BorderedRoundRectWithText(buttonW, buttonH, utils.CentreIdleColor, utils.LineIdleColor, "Return", utils.GetMenuFonts(4)),
		utils.BorderedRoundRectWithText(buttonW, buttonH, utils.CentreActiveColor, utils.LineActiveColor, "Return", utils.GetMenuFonts(4)),
	).SetHandler(func() {
		j.exit = true
		sceneManager.GoTo(NewTitleScene())
	})

	j.buttonPics[0].SetRect(j.buttonPics[0].GetIdleImage().Bounds().Add(image.Point{X: Margin, Y: titleH + Margin*2 + servListH}))
	j.buttonPics[1].SetRect(j.buttonPics[1].GetIdleImage().Bounds().Add(image.Point{X: Margin*2 + buttonW, Y: titleH + Margin*2 + servListH}))
	j.exitButtonPic.SetRect(j.exitButtonPic.GetIdleImage().Bounds().Add(image.Point{X: screenWidth - Margin - buttonW, Y: titleH + Margin*2 + servListH}))

	j.updateServersPictures(servListW-int(utils.LineThickness*2), servListH-int(utils.Radius*2)-int(utils.LineThickness*2), Margin+int(utils.LineThickness), Margin+int(utils.Radius)+int(utils.LineThickness)+titleH)
	j.updateInfoPictures(infoW-int(utils.LineThickness*2), infoH-int(utils.Radius*2)-int(utils.LineThickness*2), Margin*2+servListW+int(utils.LineThickness), Margin+int(utils.Radius)+int(utils.LineThickness)+titleH)
}

func (j *JoinScene) Update(state *GameState) error {
	if sizeChanged || j.serversUpdated {
		j.updateImages()
		j.serversUpdated = false
	}

	for i := range j.buttonPics {
		if j.canJoin {
			j.buttonPics[i].Update()
		}
	}
	j.exitButtonPic.Update()

	for i := range j.backgroundPics {
		j.backgroundPics[i].Update()
	}

	for i := range j.serverImg {
		j.serverImg[i].Update()
	}

	for i := range j.infoImg {
		j.infoImg[i].Update()
	}

	return nil
}

func (j *JoinScene) Draw(screen *ebiten.Image) {
	screen.Fill(utils.FillColor)
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

	if len(j.infoImg) != 0 {
		j.infoImg[j.selectedServer].Draw(screen)
	}
}

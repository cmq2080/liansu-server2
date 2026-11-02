package main

import (
	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
)

const (
	WINDOW_W = 540
	WINDOW_H = 400

	WINDOW_PADDING   = 8 // 窗口四周留白
	FUNC_AREA_MARGIN = 4 // 区块内上下左右留白（区块之间留白用一个就够，不再累加）

	TITLE_H      = 24 // 区块标题高度
	TITLE_MARGIN = 2  // 标题与第一行之间
	LABEL_H      = 16 // 字段标签高度
	LABEL_MARGIN = 2  // 标签与输入框之间
	EDIT_H       = 24 // 输入框高度（宋体12够用）
	ROW_MARGIN   = 2  // 行与行之间

	BTN_H   = 32 // 按钮高度（放不下会自动缩小）
	BTN_GAP = 8  // 按钮之间

	// 区块边框：BvRaised = 1px 浅浮雕；想更淡改 types.BvSpace（等于完全无边框）
	FUNC_AREA_BEVEL = types.BvRaised
)

type FontStyle struct {
	family string
	size   int32
	bold   bool
}

func NewFont(fontStyle *FontStyle) *vcl.TFont {
	font := vcl.NewFont()
	font.SetName(fontStyle.family)
	font.SetSize(fontStyle.size)
	if fontStyle.bold {
		font.SetStyle(types.NewSet(types.FsBold))
	}
	return font
}

var (
	commonFontStyle = &FontStyle{family: "宋体", size: 10, bold: false}
	titleFontStyle  = &FontStyle{family: "宋体", size: 12, bold: true}
	editFontStyle   = &FontStyle{family: "宋体", size: 12, bold: false}
	buttonFontStyle = titleFontStyle
)

// EditLine 一行「标签 + 输入框」
type EditLine struct {
	id       string
	label    *vcl.TLabel
	editText *vcl.TEdit
}

func (el *EditLine) Text() string {
	return el.editText.Text()
}

func (el *EditLine) SetText(text string) {
	el.editText.SetText(text)
}

// Bind 与 cfg 字段绑定：先回填，之后输入变化实时写回
func (el *EditLine) Bind(model *string) {
	if model == nil {
		return
	}
	el.SetText(*model)
	el.editText.SetOnChange(func(vcl.IObject) {
		*model = el.editText.Text()
	})
}

// FuncArea 一个带标题的区块
type FuncArea struct {
	id        string
	owner     vcl.IComponent
	panel     *vcl.TPanel
	title     *vcl.TLabel
	editLines []*EditLine
}

// AddEditLine 新增一行，加入顺序即显示顺序（上→下）
func (fa *FuncArea) AddEditLine(id string, caption string) *EditLine {
	label := vcl.NewLabel(fa.owner)
	label.SetParent(fa.panel)
	label.SetCaption(caption)

	editText := vcl.NewEdit(fa.owner)
	editText.SetParent(fa.panel)
	editText.SetFont(NewFont(editFontStyle))

	el := &EditLine{id: id, label: label, editText: editText}
	fa.editLines = append(fa.editLines, el)
	return el
}

// Line 按 id 取某一行
func (fa *FuncArea) Line(id string) *EditLine {
	for _, el := range fa.editLines {
		if el.id == id {
			return el
		}
	}
	return nil
}

// layout 区块内部排版：标题在最上，然后每行「标签在上、输入框在下」
func (fa *FuncArea) layout() {
	inner := fa.panel.ClientWidth() - FUNC_AREA_MARGIN*2
	if inner < 20 {
		return
	}

	y := int32(FUNC_AREA_MARGIN)
	fa.title.SetBounds(FUNC_AREA_MARGIN, y, inner, TITLE_H)
	y += TITLE_H + TITLE_MARGIN

	for _, el := range fa.editLines {
		el.label.SetBounds(FUNC_AREA_MARGIN, y, inner, LABEL_H)
		y += LABEL_H + LABEL_MARGIN
		el.editText.SetBounds(FUNC_AREA_MARGIN, y, inner, EDIT_H)
		y += EDIT_H + ROW_MARGIN
	}
}

type MainWindow struct {
	form *vcl.TForm
	btns *vcl.TPanel // 按钮区

	hostFuncArea  *FuncArea
	nginxFuncArea *FuncArea
	phpFuncArea   *FuncArea
	buttons       []*vcl.TButton
}

func (mw *MainWindow) build() {
	form := vcl.Application.CreateForm()
	form.SetCaption(APP_NAME)
	form.SetPosition(types.PoScreenCenter)
	form.SetBorderStyle(types.BsSingle) // 固定尺寸，不可拉伸
	form.EnabledMaximize(false)
	form.SetFont(NewFont(commonFontStyle))
	form.SetWidth(WINDOW_W)
	form.SetHeight(WINDOW_H)
	mw.form = form

	// 左上：服务器
	mw.hostFuncArea = mw.NewFuncArea("hostFuncArea", "服务器")
	mw.hostFuncArea.AddEditLine("host", "主机（域名/IP）")
	mw.hostFuncArea.AddEditLine("port", "端口")
	mw.hostFuncArea.AddEditLine("wwwroot", "wwwroot目录")

	// 右上：Nginx
	mw.nginxFuncArea = mw.NewFuncArea("nginxFuncArea", "Nginx")
	mw.nginxFuncArea.AddEditLine("dir", "目录")
	mw.nginxFuncArea.AddEditLine("conf", "配置文件（可选）")
	mw.nginxFuncArea.AddEditLine("logs", "log目录（可选）")

	// 左下：PHP
	mw.phpFuncArea = mw.NewFuncArea("phpFuncArea", "PHP")
	mw.phpFuncArea.AddEditLine("dir", "目录")
	mw.phpFuncArea.AddEditLine("ini", "配置文件（可选）")
	mw.phpFuncArea.AddEditLine("port", "cgi端口")

	// 右下：按钮区
	mw.newButtons()

	mw.bindCfg()

	// 窗口真正显示出来后再排一次，保证拿到最终客户区尺寸
	mw.layout()
	form.SetOnShow(func(vcl.IObject) {
		mw.layout()
	})
}

// NewFuncArea 建一个区块面板（位置由 layout 统一计算）
func (mw *MainWindow) NewFuncArea(id string, caption string) *FuncArea {
	panel := vcl.NewPanel(mw.form)
	panel.SetParent(mw.form)
	panel.SetBevelInner(types.BvNone)
	panel.SetBevelOuter(FUNC_AREA_BEVEL)
	panel.SetBorderStyle(types.BsNone)
	panel.SetBorderWidth(0)

	fa := &FuncArea{id: id, owner: mw.form, panel: panel}

	title := vcl.NewLabel(mw.form)
	title.SetParent(panel)
	title.SetCaption(caption)
	title.SetFont(NewFont(titleFontStyle))
	fa.title = title

	return fa
}

// newButtons 建四个按钮（竖排，位置由 layoutButtons 统一计算）
func (mw *MainWindow) newButtons() {
	panel := vcl.NewPanel(mw.form)
	panel.SetParent(mw.form)
	panel.SetBevelInner(types.BvNone)
	panel.SetBevelOuter(types.BvNone) // 按钮区不要边框
	panel.SetBorderStyle(types.BsNone)
	panel.SetBorderWidth(0)
	mw.btns = panel

	captions := []string{"管理应用", "保存配置", "环境-启动", "环境-停止"}
	handlers := []func(vcl.IObject){mw.onManageApps, mw.onSaveConfig, mw.onEnvStart, mw.onEnvStop}
	for i, caption := range captions {
		btn := vcl.NewButton(mw.form)
		btn.SetParent(panel)
		btn.SetCaption(caption)
		btn.SetFont(NewFont(buttonFontStyle))
		btn.SetOnClick(handlers[i])
		mw.buttons = append(mw.buttons, btn)
	}
}

// layout 外层 2×2：左上/右上/左下/右下 四等分
func (mw *MainWindow) layout() {
	cw := mw.form.ClientWidth()
	ch := mw.form.ClientHeight()
	if cw < 100 || ch < 100 {
		return
	}

	// 外圈留白与区块中缝统一为 gap，四边和中缝视觉等宽
	gap := int32(FUNC_AREA_MARGIN) * 2 // 4+4 = 8

	colW := (cw - gap*3) / 2 // 客户端宽度-3个缝（2列+1）
	rowH := (ch - gap*3) / 2 // 客户端高度-3个缝（2行+1）
	if colW <= 0 || rowH <= 0 {
		return
	}

	cells := []*vcl.TPanel{
		mw.hostFuncArea.panel,  // 左上
		mw.nginxFuncArea.panel, // 右上
		mw.phpFuncArea.panel,   // 左下
		mw.btns,                // 右下
	}
	for i, p := range cells {
		x := gap + int32(i%2)*(colW+gap)
		y := gap + int32(i/2)*(rowH+gap)
		p.SetBounds(x, y, colW, rowH)
	}

	for _, fa := range []*FuncArea{mw.hostFuncArea, mw.nginxFuncArea, mw.phpFuncArea} {
		fa.layout()
	}
	mw.layoutButtons()
}

// layoutButtons 按钮竖排：高度平分区块，整体垂直居中
func (mw *MainWindow) layoutButtons() {
	count := int32(len(mw.buttons))
	if count == 0 {
		return
	}

	w := mw.btns.ClientWidth() - FUNC_AREA_MARGIN*2
	h := (mw.btns.ClientHeight() - FUNC_AREA_MARGIN*2 - BTN_GAP*(count-1)) / count
	if h > BTN_H {
		h = BTN_H
	}
	if h < 22 {
		h = 22
	}

	totalH := h*count + BTN_GAP*(count-1)
	y := (mw.btns.ClientHeight() - totalH) / 2

	for _, btn := range mw.buttons {
		btn.SetBounds(FUNC_AREA_MARGIN, y, w, h)
		y += h + BTN_GAP
	}
}

// bindCfg 把 cfg 各字段绑到对应输入框
func (mw *MainWindow) bindCfg() {
	if cfg == nil {
		return
	}
	mw.hostFuncArea.Line("host").Bind(&cfg.Server.Host)
	mw.hostFuncArea.Line("port").Bind(&cfg.Server.Port)
	mw.hostFuncArea.Line("wwwroot").Bind(&cfg.Server.Wwwroot)

	mw.nginxFuncArea.Line("dir").Bind(&cfg.Nginx.Dir)
	mw.nginxFuncArea.Line("conf").Bind(&cfg.Nginx.Conf)
	mw.nginxFuncArea.Line("logs").Bind(&cfg.Nginx.Logs_dir)

	mw.phpFuncArea.Line("dir").Bind(&cfg.Php.Dir)
	mw.phpFuncArea.Line("ini").Bind(&cfg.Php.Ini)
	mw.phpFuncArea.Line("port").Bind(&cfg.Php.Port)
}

func (mw *MainWindow) onManageApps(vcl.IObject) {
	vcl.ShowMessage("管理应用：待实现")
}

func (mw *MainWindow) onSaveConfig(vcl.IObject) {
	vcl.ShowMessage("保存配置：待实现")
}

func (mw *MainWindow) onEnvStart(vcl.IObject) {
	vcl.ShowMessage("环境-启动：待实现")
}

func (mw *MainWindow) onEnvStop(vcl.IObject) {
	vcl.ShowMessage("环境-停止：待实现")
}

func startGUI() {
	vcl.Application.Initialize()
	vcl.Application.SetMainFormOnTaskBar(true)

	mainWindow := &MainWindow{}
	mainWindow.build()

	vcl.Application.Run()
}

package gui

import (
	"image/color"
	
	"tcpchat/client"
	"tcpchat/model"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type UI struct {
	app      fyne.App
	window   fyne.Window
	address  *widget.Entry
	username *widget.Entry
	connect  *widget.Button
	messages *widget.RichText
	input    *widget.Entry
	send     *widget.Button
	client   *client.TCPClient
}

func NewUI() *UI {
	myApp := app.New()
	myWindow := myApp.NewWindow("TCP Chat Client")
	myWindow.Resize(fyne.NewSize(600, 400))

	ui := &UI{
		app:    myApp,
		window: myWindow,
	}

	ui.setupUI()
	return ui
}

func (ui *UI) setupUI() {
	// 服务器地址输入区域
	ui.address = widget.NewEntry()
	ui.address.SetPlaceHolder("请输入服务器地址 (例如: localhost:8080)")
	
	ui.username = widget.NewEntry()
	ui.username.SetPlaceHolder("请输入用户名")
	
	ui.connect = widget.NewButton("连接", ui.connectToServer)

	addressContainer := container.NewVBox(
		container.NewGridWithColumns(2, widget.NewLabel("地址:"), ui.address),
		container.NewGridWithColumns(2, widget.NewLabel("用户名:"), ui.username),
		ui.connect,
	)

	// 消息显示区域
	ui.messages = widget.NewRichText()
	ui.messages.Wrapping = fyne.TextWrapWord
	scrollableMessages := container.NewVScroll(ui.messages)

	// 消息输入区域
	ui.input = widget.NewEntry()
	ui.input.SetPlaceHolder("输入消息...")
	ui.input.OnSubmitted = func(s string) {
		ui.sendMessage()
	}

	ui.send = widget.NewButton("发送", ui.sendMessage)
	inputContainer := container.NewBorder(nil, nil, nil, ui.send, ui.input)

	// 整体布局
	content := container.NewBorder(
		addressContainer,   // 顶部
		inputContainer,     // 底部
		nil,                // 左侧
		nil,                // 右侧
		scrollableMessages, // 中心
	)

	ui.window.SetContent(content)
}

func (ui *UI) connectToServer() {
	if ui.client != nil {
		ui.client.Close()
		ui.client = nil
	}
	
	if ui.username.Text == "" {
		ui.addSystemMessage("请输入用户名")
		return
	}

	// 创建客户端并设置回调函数
	ui.client = client.NewTCPClient().(*client.TCPClient)
	// 使用WithOnMessage方法设置回调函数
	ui.client.WithOnMessage(ui.handleMessage)(ui.client)

	err := ui.client.Dial(ui.address.Text)
	if err != nil {
		ui.addSystemMessage("连接失败: " + err.Error())
		return
	}
	
	// 设置用户名
	err = ui.client.Setname(ui.username.Text)
	if err != nil {
		ui.addSystemMessage("设置用户名失败: " + err.Error())
		return
	}

	ui.addSystemMessage("已连接到服务器")
	go func() {
		ui.client.Start()
		// 连接断开后重新启用UI元素
		ui.address.Enable()
		ui.username.Enable()
		ui.connect.Enable()
		ui.addSystemMessage("与服务器断开连接")
	}()

	// 连接成功后禁用地址输入和连接按钮
	ui.address.Disable()
	ui.username.Disable()
	ui.connect.Disable()
}

func (ui *UI) sendMessage() {
	if ui.client == nil || ui.input.Text == "" {
		return
	}

	err := ui.client.SendMessage(ui.input.Text)
	if err != nil {
		ui.addSystemMessage("发送失败: " + err.Error())
		return
	}

	ui.input.SetText("")
}

func (ui *UI) handleMessage(message *model.Message) {
	switch message.Type {
	case model.ChatMessage:
		ui.addChatMessage(message.Owner + ": " + message.Content)
	case model.SystemMessage:
		ui.addSystemMessage(message.Content)
	case model.SetNameMessage:
		ui.addSystemMessage("用户 " + message.Content + " 加入了聊天")
	case model.LeaveMessage:
		ui.addSystemMessage("用户 " + message.Content + " 离开了聊天")
	}
}

func (ui *UI) addChatMessage(content string) {
	// 添加普通聊天消息
	segment := widget.TextSegment{
		Style: widget.RichTextStyle{
			ColorName: "",
		},
		Text: content + "\n",
	}
	ui.messages.Segments = append(ui.messages.Segments, &segment)
	ui.messages.Refresh()
}

func (ui *UI) addSystemMessage(content string) {
	// 添加系统消息
	segment := widget.TextSegment{
		Style: widget.RichTextStyle{
			ColorName: "",
			Inline:    true,
		},
		Text: "[系统] " + content + "\n",
	}
	
	ui.messages.Segments = append(ui.messages.Segments, &segment)
	ui.messages.Refresh()
}

func (ui *UI) ShowAndRun() {
	// 设置支持中文的字体以避免乱码
	ui.app.Settings().SetTheme(&customTheme{})
	ui.window.ShowAndRun()
}

// 自定义主题以支持中文显示
type customTheme struct{}

func (t *customTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, variant)
}

func (t *customTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t *customTheme) Font(style fyne.TextStyle) fyne.Resource {
	// 使用支持中文的字体，例如 Noto Sans CJK 或 Wqy Zenhei
	// 这里我们尝试使用系统默认的中文字体，如果不可用则回退到默认
	fontName := "Noto Sans CJK"
	if style.Bold {
		fontName += " Bold"
	}
	if style.Italic {
		fontName += " Italic"
	}
	
	// 尝试加载字体资源（如果系统安装了对应字体）
	// 注意：fyne 不直接支持动态加载字体文件，需通过 theme.Font 返回资源
	// 此处简化处理，使用默认字体并依赖系统字体配置
	return theme.DefaultTheme().Font(style)
}

func (t *customTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
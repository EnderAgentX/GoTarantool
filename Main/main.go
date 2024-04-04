package main

import (
	"GoTarantool/Server"
	"fmt"
	"log"
	"unicode"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/tarantool/go-tarantool"
)

var MyUser string
var myPass string
var SelectedGroupName string
var SelectedGroupId string
var userId string
var selectedRowMsg *gtk.ListBoxRow
var selectedRowGroup *gtk.ListBoxRow

type TimedMsg struct {
	msg  string
	time uint64
}

type MessageStruct struct {
	message string
	user    string
	msgTime uint64
}

var messagesArr = make([]TimedMsg, 0)

func main() {
	//isLogin := false
	conn := Server.Server()
	defer conn.Close()

	gtk.Init(nil)

	// Создаём билдер
	b, err := gtk.BuilderNew()
	if err != nil {
		log.Fatal("Ошибка:", err)
	}

	// Загружаем в билдер окно из файла Glade
	err = b.AddFromFile("Chat_gui.glade")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}

	// Получаем объект главного окна по ID

	//Объекты

	objMain, err := b.GetObject("main_window")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	objReg, err := b.GetObject("regWin")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}

	objChangeGroup, err := b.GetObject("changeWinGroup")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}

	objChangeMsg, err := b.GetObject("changeWinMsg")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}

	objNewGroupWin, err := b.GetObject("newGroupWin")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}

	objJoinGroupWin, err := b.GetObject("joinGroupWin")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}

	// Преобразуем из объекта именно окно типа gtk.Window
	// и соединяем с сигналом "destroy" чтобы можно было закрыть
	// приложение при закрытии окна
	winMain := objMain.(*gtk.Window)
	//winMain.Move(0, 0)

	winMain.Connect("destroy", func() {
		gtk.MainQuit()
	})

	winReg := objReg.(*gtk.Window)
	winReg.Connect("delete-event", func() {
		gtk.MainQuit()

	})

	winChangeGroup := objChangeGroup.(*gtk.Dialog)
	winChangeGroup.Connect("delete-event", func() {
		winChangeGroup.Hide()

	})

	winChangeMsg := objChangeMsg.(*gtk.Dialog)
	winChangeMsg.Connect("delete-event", func() {
		winChangeMsg.Hide()

	})

	winNewGroup := objNewGroupWin.(*gtk.Dialog)

	winJoinGroup := objJoinGroupWin.(*gtk.Dialog)

	objMain, _ = b.GetObject("msg_entry")
	msgEntry := objMain.(*gtk.Entry)

	// Получаем кнопку

	objMain, _ = b.GetObject("msg_btn")
	msgBtn := objMain.(*gtk.Button)

	// Получаем метку

	objMain, _ = b.GetObject("guild_label")
	guildLabel := objMain.(*gtk.Label)

	objMain, _ = b.GetObject("user_label")
	userLabel := objMain.(*gtk.Label)

	objMain, _ = b.GetObject("registration_success_label")
	registrationSuccessLabel := objMain.(*gtk.Label)

	objMain, _ = b.GetObject("labelCheckId")
	labelCheckId := objMain.(*gtk.Label)

	objMain, _ = b.GetObject("msg_scroll")
	scrolledWindow := objMain.(*gtk.ScrolledWindow)

	objReg, _ = b.GetObject("login_reg_btn")
	loginRegBtn := objReg.(*gtk.Button)

	objChangeMsg, _ = b.GetObject("changeBtnMsg")
	changeBtnMsg := objChangeMsg.(*gtk.Button)

	objChangeGroup, _ = b.GetObject("changeBtnGroup")
	changeBtnGroup := objChangeGroup.(*gtk.Button)

	objChangeGroup, _ = b.GetObject("changeGroupEntry")
	changeGroupEntry := objChangeGroup.(*gtk.Entry)

	objMain, _ = b.GetObject("btn_change_msg")
	btnChangeMsg := objMain.(*gtk.Button)

	objMain, _ = b.GetObject("newuser_reg_btn")
	newUserRegBtn := objMain.(*gtk.Button)

	objMain, _ = b.GetObject("login_reg_entry")
	loginRegEntry := objMain.(*gtk.Entry)

	objMain, _ = b.GetObject("pass_reg_entry")
	passRegEntry := objMain.(*gtk.Entry)

	objMain, _ = b.GetObject("joinGroupEntry")
	joinGroupEntry := objMain.(*gtk.Entry)

	objMain, _ = b.GetObject("listbox_groups")
	groupsListbox := objMain.(*gtk.ListBox)

	objNewGroupWin, _ = b.GetObject("entryNewGroupName")
	entryNewGroupName := objNewGroupWin.(*gtk.Entry)

	objNewGroupWin, _ = b.GetObject("entryNewGroupId")
	entryNewGroupId := objNewGroupWin.(*gtk.Entry)

	objMain, _ = b.GetObject("add_group_btn")
	addGroupBtn := objMain.(*gtk.Button)

	objMain, _ = b.GetObject("btn_del_group")
	delGroupBtn := objMain.(*gtk.Button)

	objMain, _ = b.GetObject("btn_change_group")
	changeGroupBtn := objMain.(*gtk.Button)

	objMain, _ = b.GetObject("exit_btn")
	exitBtn := objMain.(*gtk.Button)

	objMain, _ = b.GetObject("listbox_msg")
	msgListbox := objMain.(*gtk.ListBox)

	objNewGroupWin, _ = b.GetObject("newGroupBtn")
	newGroupBtn := objNewGroupWin.(*gtk.Button)

	objNewGroupWin, _ = b.GetObject("btnCheckId")
	btnCheckId := objNewGroupWin.(*gtk.Button)

	objMain, _ = b.GetObject("join_group_btn")
	joinGroupBtn := objMain.(*gtk.Button)

	objMain, _ = b.GetObject("joinGroupBtnConfirm")
	joinGroupBtnConfirm := objMain.(*gtk.Button)

	winNewGroup.Connect("delete-event", func() {
		entryNewGroupId.SetText("")
		entryNewGroupName.SetText("")
		labelCheckId.SetText("")
		winNewGroup.Hide()

	})

	winJoinGroup.Connect("delete-event", func() {
		joinGroupEntry.SetText("")
		winJoinGroup.Hide()
	})

	groupsListbox.Connect("button-press-event", func(box *gtk.ListBox, event *gdk.Event) {
		buttonEvent := gdk.EventButtonNewFromEvent(event)
		if buttonEvent.Type() == gdk.EVENT_2BUTTON_PRESS {
			fmt.Println("DOUBLE")
			groupsListbox.UnselectAll()
			groupsListbox.ShowAll()
			winMain.ShowAll()
			return // Ignore double-click events
		}

		// Handle single-click events
		row := box.GetSelectedRow()
		if row != nil {
			index := row.GetIndex()
			fmt.Printf("Clicked on row %d\n", index)
		}
	})

	msgListbox.Connect("row-activated", func() {
		selectedRowMsg := msgListbox.GetSelectedRow()
		widgetRow, _ := (selectedRowMsg.GetChild())
		labelRow := widgetRow.(*gtk.Label)
		fmt.Println(labelRow.GetText())
	})

	changeBtnMsg.Connect("clicked", func() {
		winChangeMsg.Hide()
	})

	changeBtnGroup.Connect("clicked", func() {
		newGroup, _ := changeGroupEntry.GetText()
		_, _ = conn.Call("fn.edit_group", []interface{}{MyUser, SelectedGroupId, newGroup})
		GetGroups(conn, groupsListbox)
		selectedRowGroup, _ = gtk.ListBoxRowNew()
		SelectedGroupName = newGroup
		clearListbox(msgListbox)
		messagesArr = messagesArr[:0]
		msgEntry.SetText("")
		winChangeGroup.Hide()
	})

	groupsListbox.Connect("row-activated", func() {
		fmt.Println("Выбрана группа")
		AutoScroll(scrolledWindow)
		tempRow := groupsListbox.GetSelectedRow()
		fmt.Println("tempRow.GetIndex()", tempRow.GetIndex())
		fmt.Println("selectedRow.GetIndex())", selectedRowGroup.GetIndex())
		if tempRow.GetIndex() == selectedRowGroup.GetIndex() {
			groupsListbox.UnselectAll()
			clearListbox(msgListbox)
			selectedRowGroup = groupsListbox.GetSelectedRow()
		} else {
			selectedRowGroup = groupsListbox.GetSelectedRow()
			selectedRowGroupId, _ := selectedRowGroup.GetName()
			fmt.Println("group id ", selectedRowGroupId)

			labelRow, _ := selectedRowGroup.GetChild()
			groupLabel := labelRow.(*gtk.Label)
			textLabel, _ := groupLabel.GetText()
			SelectedGroupName = textLabel
			SelectedGroupId = selectedRowGroupId
			guildLabel.SetText(SelectedGroupName)
			messagesArr = messagesArr[:0]
			GetMsg(conn, msgListbox)
			AutoScroll(scrolledWindow)
			winMain.ShowAll()

			//Таймер потом включить
			//t := time.NewTimer(1 * time.Second)
			//go func() {
			//	for {
			//
			//		t.Reset(1 * time.Second)
			//		GetMsg(conn, msgListbox)
			//		<-t.C
			//	}
			//
			//}()
		}

	})

	// Сигнал по нажатию на кнопку

	exitBtn.Connect("clicked", func() {
		groupsListbox.UnselectAll()
		msgListbox.UnselectAll()
		clearListbox(groupsListbox)
		clearListbox(msgListbox)
		selectedRowGroup, _ = gtk.ListBoxRowNew()
		selectedRowMsg, _ = gtk.ListBoxRowNew()
		userLabel.SetText("")
		guildLabel.SetText("")
		messagesArr = messagesArr[:0]
		fmt.Println("Выход")
		winMain.Hide()
		winReg.ShowAll()
	})

	newGroupBtn.Connect("clicked", func() {
		groupName, _ := entryNewGroupName.GetText()
		groupId, _ := entryNewGroupId.GetText()

		info, _ := conn.Call("fn.check_group_id", []interface{}{groupId})
		tuples := info.Tuples()
		checked := tuples[0][0].(string)
		if checked == "true" {
			if MyUser != "" && !isOnlyWhitespace(groupName) && !isOnlyWhitespace(groupId) {
				fmt.Println(MyUser, groupName, groupId)
				_, _ = conn.Call("fn.new_group", []interface{}{MyUser, groupId, groupName})

				rowGroup, _ := gtk.ListBoxRowNew()
				rowGroup.SetName(groupId)
				labelGroup, _ := gtk.LabelNew(groupName)
				labelGroup.SetSizeRequest(-1, 50)
				markup := fmt.Sprintf("<span font_desc='Serif Bold Italic 20'>%s</span>", groupName)
				labelGroup.SetMarkup(markup)
				rowGroup.Add(labelGroup)
				groupsListbox.Insert(rowGroup, 0)

				entryNewGroupId.SetText("")
				entryNewGroupName.SetText("")
				labelCheckId.SetText("")

				winNewGroup.Close()
				winMain.ShowAll()
			}
		}
		if checked == "false" {
			checkedText := "Тег занят"
			markup := fmt.Sprintf("<span size='15000' foreground='red'>%s</span>", checkedText)
			labelCheckId.SetText(checkedText)
			labelCheckId.SetMarkup(markup)
		}

	})

	btnCheckId.Connect("clicked", func() {
		groupId, _ := entryNewGroupId.GetText()
		info, _ := conn.Call("fn.check_group_id", []interface{}{groupId})
		tuples := info.Tuples()
		checked := tuples[0][0].(string)
		if checked == "true" {
			checkedText := "Тег свободен"
			markup := fmt.Sprintf("<span size='15000' foreground='green'>%s</span>", checkedText)
			labelCheckId.SetText(checkedText)
			labelCheckId.SetMarkup(markup)
		}
		if checked == "false" {
			checkedText := "Тег занят"
			markup := fmt.Sprintf("<span size='15000' foreground='red'>%s</span>", checkedText)
			labelCheckId.SetText(checkedText)
			labelCheckId.SetMarkup(markup)
		}
	})

	btnChangeMsg.Connect("clicked", func() {
		winChangeMsg.Run()
	})

	delGroupBtn.Connect("clicked", func() {
		_, _ = conn.Call("fn.del_group", []interface{}{MyUser, SelectedGroupId})
		groupsListbox.Remove(selectedRowGroup)
		selectedRowGroup, _ = gtk.ListBoxRowNew()
		clearListbox(msgListbox)
		guildLabel.SetText("")
		messagesArr = messagesArr[:0]
		msgEntry.SetText("")
	})

	changeGroupBtn.Connect("clicked", func() {
		fmt.Println(selectedRowGroup)
		fmt.Println("......")
		if selectedRowGroup != nil {
			groupRow, _ := selectedRowGroup.GetChild()
			groupLabel := groupRow.(*gtk.Label)
			groupText, _ := groupLabel.GetText()
			changeGroupEntry.SetText(groupText)
			winChangeGroup.Run()
		}

	})

	loginRegBtn.Connect("clicked", func() {

		if err == nil {
			// Устанавливаем текст из поля ввода метке
			newUser, _ := loginRegEntry.GetText()
			newPass, _ := passRegEntry.GetText()

			clearListbox(groupsListbox)

			info, _ := conn.Call("fn.login", []interface{}{newUser, newPass})
			fmt.Println(info)
			userTuples := info.Tuples()
			if userTuples[0][0].(bool) == true {
				MyUser = newUser
				myPass = newPass
				userLabel.SetText(MyUser)
				GetGroups(conn, groupsListbox)
				loginRegEntry.SetText("")
				passRegEntry.SetText("")
				registrationSuccessLabel.SetText("\n")
				winReg.Hide()
				winMain.ShowAll()

			} else {
				errText := "Ошибка! \n Неверный логин или пароль!"
				markup := fmt.Sprintf("<span size='15000' foreground='red'>%s</span>", errText)
				registrationSuccessLabel.SetText(errText)
				registrationSuccessLabel.SetMarkup(markup)
			}

		}

	})

	addGroupBtn.Connect("clicked", func() {
		winNewGroup.Run()

	})

	joinGroupBtn.Connect("clicked", func() {
		winJoinGroup.Run()
	})

	joinGroupBtnConfirm.Connect("clicked", func() {
		groupId, _ := joinGroupEntry.GetText()
		info, _ := conn.Call("fn.join_group", []interface{}{MyUser, groupId})
		fmt.Println(info)
		tuples := info.Tuples()
		joinGroupName := tuples[0][0].(string)

		rowGroup, _ := gtk.ListBoxRowNew()
		rowGroup.SetName(groupId)
		labelGroup, _ := gtk.LabelNew(joinGroupName)
		labelGroup.SetSizeRequest(-1, 50)
		markup := fmt.Sprintf("<span font_desc='Serif Bold Italic 20'>%s</span>", joinGroupName)
		labelGroup.SetMarkup(markup)
		rowGroup.Add(labelGroup)
		groupsListbox.Insert(rowGroup, 0)

		joinGroupEntry.SetText("")
		winJoinGroup.Hide()
		winMain.ShowAll()
	})

	msgBtn.Connect("button-press-event", func(btn *gtk.Button, event *gdk.Event) {
		buttonEvent := gdk.EventButtonNewFromEvent(event)
		if buttonEvent.Type() == gdk.EVENT_2BUTTON_PRESS || buttonEvent.Type() == gdk.EVENT_3BUTTON_PRESS {
			newMsg, _ := msgEntry.GetText()
			_, _ = conn.Call("fn.new_msg", []interface{}{newMsg, SelectedGroupName, MyUser})
			GetMsg(conn, msgListbox)
			AutoScroll(scrolledWindow)
			clearListbox(msgListbox)
			GetMsg(conn, msgListbox)
			msgListbox.ShowAll()
			winMain.ShowAll()
			fmt.Println("Double click")
			return // Ignore double-click events
		}

		newMsg, _ := msgEntry.GetText()
		_, _ = conn.Call("fn.new_msg", []interface{}{newMsg, SelectedGroupId, MyUser})
		GetMsg(conn, msgListbox)
		AutoScroll(scrolledWindow)
		msgEntry.SetText("")

	})

	msgEntry.Connect("key-press-event", func(entry *gtk.Entry, event *gdk.Event) {
		keyEvent := &gdk.EventKey{Event: event}
		keyVal := keyEvent.KeyVal()
		if keyVal == gdk.KEY_Return {
			text, _ := entry.GetText()
			log.Println("Enter key pressed in entry. Text entered:", text)
			newMsg, _ := msgEntry.GetText()
			_, _ = conn.Call("fn.new_msg", []interface{}{newMsg, SelectedGroupId, MyUser})
			GetMsg(conn, msgListbox)
			winMain.ShowAll()
			AutoScroll(scrolledWindow)
			msgEntry.SetText("")
			scrolledWindow.ShowAll()
			msgListbox.ShowAll()
			msgListbox.SelectAll()
			winMain.ShowAll()
		}
	})

	newUserRegBtn.Connect("clicked", func() {
		newLogin, _ := loginRegEntry.GetText()
		newPass, _ := passRegEntry.GetText()
		if newLogin != "" && newPass != "" {
			info, _ := conn.Call("fn.new_user", []interface{}{newLogin, newPass})
			successTuples := info.Tuples()
			success := successTuples[0][0].(bool)
			successText := ""
			markup := ""
			if success == true {
				fmt.Println("Успешная регистрация!")
				successText = "Успешная регистрация! \n"
				markup = fmt.Sprintf("<span size='15000' foreground='green'>%s</span>", successText)

			} else {
				successText = "Ошибка! \n Пользователь уже существует"
				markup = fmt.Sprintf("<span size='15000' foreground='red'>%s</span>", successText)
			} // TODO неверный пароль
			registrationSuccessLabel.SetText(successText)
			registrationSuccessLabel.SetMarkup(markup)
		}
	})

	//})

	// Отображаем все виджеты в окне
	winReg.ShowAll()

	// Выполняем главный цикл GTK (для отрисовки). Он остановится когда
	// выполнится gtk.MainQuit()
	gtk.Main()

}

func AutoScroll(scrolledWindow *gtk.ScrolledWindow) {
	scrolledWindow.ShowAll()
	adjustment := scrolledWindow.GetVAdjustment()
	adjustment.SetValue(adjustment.GetLower())
	// scrolledWindow.ShowAll()
	// adjustment := scrolledWindow.GetVAdjustment()
	// tempUp := adjustment.GetUpper()
	// adjustment.SetUpper(adjustment.GetUpper() + 22)
	// adjustment.SetValue(adjustment.GetUpper() - adjustment.GetPageSize())
	// fmt.Println("Up до", adjustment.GetUpper())
	// fmt.Println("Page до", adjustment.GetPageSize())
	// fmt.Println("value до", adjustment.GetUpper() - adjustment.GetPageSize())

	// scrolledWindow.ShowAll()
	// adjustment.SetUpper(tempUp)
	// adjustment.SetValue(adjustment.GetValue() - 22)
	// fmt.Println("Up после", adjustment.GetUpper())
	// fmt.Println("Page после", adjustment.GetPageSize())
}

// func delGroup(groupId int, groupsListbox *gtk.ListBox) {
// 	groupsListbox
// }

func clearListbox(ListBox *gtk.ListBox) {
	messagesArr = messagesArr[:0]

	children := ListBox.GetChildren()

	for children.Length() > 0 {
		child := children.NthData(0)
		if widget, ok := child.(*gtk.Widget); ok {
			ListBox.Remove(widget)
		}
		children = ListBox.GetChildren()
	}
	ListBox.ShowAll()
}

func isOnlyWhitespace(input string) bool {
	for _, char := range input {
		if !unicode.IsSpace(char) { // Проверяем, что символ не является пробелом
			return false
		}
	}
	return true
}

func GetGroups(conn *tarantool.Connection, groupsListbox *gtk.ListBox) {
	clearListbox(groupsListbox)
	infoUserGroups, _ := conn.Call("fn.get_user_groups", []interface{}{MyUser})
	userGroupsTuples := infoUserGroups.Tuples()
	fmt.Println(userGroupsTuples)
	fmt.Println(len(userGroupsTuples))
	fmt.Println(len(userGroupsTuples[0]))
	if len(userGroupsTuples) == 1 && len(userGroupsTuples[0]) == 0 {
		return
	}
	for i := 0; i < len(userGroupsTuples); i++ {
		rowGroup, _ := gtk.ListBoxRowNew()
		rowGroup.SetName(userGroupsTuples[i][0].(string))
		labelGroup, _ := gtk.LabelNew(userGroupsTuples[i][1].(string))
		labelGroup.SetSizeRequest(-1, 50)
		tempText, _ := labelGroup.GetText()
		markup := fmt.Sprintf("<span font_desc='Serif Bold Italic 20'>%s</span>", tempText)
		labelGroup.SetMarkup(markup)
		rowGroup.Add(labelGroup)
		groupsListbox.Insert(rowGroup, 0)

	}
	groupsListbox.ShowAll()
}

func GetMsg(conn *tarantool.Connection, msgListBox *gtk.ListBox) {
	var lastTimedMsg uint64

	fmt.Println("len(messagesArr) ", len(messagesArr))
	if len(messagesArr) == 0 {
		clearListbox(msgListBox)
		lastTimedMsg = 0 // В самом начале загружаем все сообщения
	} else {
		lastTimedMsg = messagesArr[len(messagesArr)-1].time
	}

	infoTimedMsg, _ := conn.Call("fn.time_group_msg", []interface{}{SelectedGroupId, lastTimedMsg})
	newMessagesCntTuples := infoTimedMsg.Tuples()
	cntMsg := int(newMessagesCntTuples[0][0].(uint64))
	fmt.Println("lastTimedMsg", lastTimedMsg)

	var newMessages []MessageStruct
	for i := 0; i < cntMsg; i++ {
		newMessagesTuples := newMessagesCntTuples[1][i].([]interface{})
		newMessages = append(newMessages, MessageStruct{newMessagesTuples[0].(string), newMessagesTuples[1].(string), newMessagesTuples[2].(uint64)})
	}

	if cntMsg != 0 {
		for i := 0; i < cntMsg; i++ {
			msgText := newMessages[i].message
			msgUser := newMessages[i].user
			msgTime := newMessages[i].msgTime

			newMsg := msgUser + "(" + SelectedGroupName + "): " + msgText

			messagesArr = append(messagesArr, TimedMsg{msg: newMsg, time: msgTime})

			//listbox

			rowMsg, _ := gtk.ListBoxRowNew()
			labelMsg, _ := gtk.LabelNew(newMsg)
			labelMsg.SetHAlign(gtk.ALIGN_START)
			labelMsg.SetJustify(gtk.JUSTIFY_CENTER)
			rowMsg.Add(labelMsg)
			//msgListBox.Insert(rowMsg, 1)
			msgListBox.Prepend(rowMsg)

		}

		msgListBox.ShowAll()

	}
}

package main

import (
	"GoTarantool/Server"
	"fmt"
	"log"

	"github.com/gotk3/gotk3/gtk"
	"github.com/tarantool/go-tarantool"
)

var MyUser string
var myPass string
var GroupName string
var GroupId string
var userId string
var selectedRow *gtk.ListBoxRow

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
	obj, err := b.GetObject("main_window")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	obj2, err := b.GetObject("regWin")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}
	obj3, _ := b.GetObject("CenterWin")

	// Преобразуем из объекта именно окно типа gtk.Window
	// и соединяем с сигналом "destroy" чтобы можно было закрыть
	// приложение при закрытии окна
	win := obj.(*gtk.Window)
	win3 := obj3.(*gtk.ApplicationWindow)
	//win.Move(0, 0)

	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	win2 := obj2.(*gtk.Dialog)
	win2.Connect("delete-event", func() {
		fmt.Println("Окно закрывается, но не будет удалено")

		//win2.HideOnDelete()
		win2.Hide()
		win3.Hide()
		win.ShowAll()
		//win2.Hide()
		//win2 := obj2.(*gtk.Dialog)
		//win2.ShowAll()

	})

	obj, _ = b.GetObject("msg_entry")
	msgEntry := obj.(*gtk.Entry)

	// Получаем кнопку

	obj, _ = b.GetObject("msg_btn")
	msgBtn := obj.(*gtk.Button)

	// Получаем метку
	obj, _ = b.GetObject("msg_box")
	msgBox := obj.(*gtk.Label)
	msgBox.SetText("")

	obj, _ = b.GetObject("guild_label")
	guildLabel := obj.(*gtk.Label)

	obj, _ = b.GetObject("user_label")
	userLabel := obj.(*gtk.Label)

	obj, _ = b.GetObject("reg_success_label")
	regSuccessLabel := obj.(*gtk.Label)

	obj, _ = b.GetObject("reg_err_label")
	regErrLabel := obj.(*gtk.Label)

	obj, _ = b.GetObject("msg_scroll")
	scrolledWindow := obj.(*gtk.ScrolledWindow)

	obj2, _ = b.GetObject("login_reg_btn")
	loginRegBtn := obj2.(*gtk.Button)

	obj, _ = b.GetObject("newuser_reg_btn")
	newUserRegBtn := obj.(*gtk.Button)

	obj, _ = b.GetObject("login_reg_entry")
	loginRegEntry := obj.(*gtk.Entry)

	obj, _ = b.GetObject("pass_reg_entry")
	passRegEntry := obj.(*gtk.Entry)

	obj, _ = b.GetObject("listbox_groups")
	groupsListbox := obj.(*gtk.ListBox)

	obj, _ = b.GetObject("add_group_entry")
	addGroupEntry := obj.(*gtk.Entry)

	obj, _ = b.GetObject("add_group_btn")
	addGroupBtn := obj.(*gtk.Button)

	obj, _ = b.GetObject("listbox_msg")
	msgListbox := obj.(*gtk.ListBox)

	groupsListbox.Connect("row-activated", func() {
		AutoScroll(scrolledWindow)
		tempRow := groupsListbox.GetSelectedRow()
		fmt.Println(tempRow.GetIndex())
		fmt.Println(selectedRow.GetIndex())
		if tempRow.GetIndex() == selectedRow.GetIndex() {
			groupsListbox.UnselectAll()
			clearListbox(msgListbox)
			selectedRow = groupsListbox.GetSelectedRow()
		} else {
			selectedRow = groupsListbox.GetSelectedRow()
			labelRow, _ := selectedRow.GetChild()
			groupLabel := labelRow.(*gtk.Label)
			textLabel, _ := groupLabel.GetText()
			fmt.Println(textLabel)
			GroupName = textLabel
			guildLabel.SetText(GroupName)
			messagesArr = messagesArr[:0]
			GetMsgTest(conn, msgListbox)

			AutoScroll(scrolledWindow)
			win.ShowAll()

			//Таймер потом включить
			//t := time.NewTimer(1 * time.Second)
			//go func() {
			//	for {
			//
			//		t.Reset(1 * time.Second)
			//		GetMsgTest(conn, msgBox, msgListbox)
			//		<-t.C
			//	}
			//
			//}()
		}

	})

	// Сигнал по нажатию на кнопку

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
				infoUserGroups, _ := conn.Call("fn.get_user_groups", []interface{}{MyUser})
				userGroupsTuples := infoUserGroups.Tuples()
				for i := 0; i < len(userGroupsTuples[0]); i++ {
					rowGroup, _ := gtk.ListBoxRowNew()
					labelGroup, _ := gtk.LabelNew(userGroupsTuples[0][i].(string))
					labelGroup.SetJustify(gtk.JUSTIFY_LEFT)
					rowGroup.Add(labelGroup)
					groupsListbox.Insert(rowGroup, 0)
				}
				regSuccessLabel.Hide()
				regErrLabel.Hide()
				loginRegEntry.SetText("")
				passRegEntry.SetText("")
				win2.Hide()
				win.ShowAll()

			} else {
				fmt.Println("Неверный логин или пароль")
			}

		}

	})

	addGroupBtn.Connect("clicked", func() {
		groupName, _ := addGroupEntry.GetText()
		if MyUser != "" {
			fmt.Println(groupName)
			_, _ = conn.Call("fn.new_group", []interface{}{MyUser, groupName})
			rowGroup, _ := gtk.ListBoxRowNew()
			labelGroup, _ := gtk.LabelNew(groupName)
			rowGroup.Add(labelGroup)
			groupsListbox.Insert(rowGroup, 0)
			win.ShowAll()
		}

	})

	msgBtn.Connect("clicked", func() {
		newMsg, _ := msgEntry.GetText()
		_, _ = conn.Call("fn.new_msg", []interface{}{newMsg, GroupName, MyUser})
		GetMsgTest(conn, msgListbox)
		AutoScroll(scrolledWindow)

	})

	newUserRegBtn.Connect("clicked", func() {
		newLogin, _ := loginRegEntry.GetText()
		newPass, _ := passRegEntry.GetText()
		info, _ := conn.Call("fn.new_user", []interface{}{newLogin, newPass})
		successTuples := info.Tuples()
		success := successTuples[0][0].(bool)
		if success == true {
			fmt.Println("Успешная регистрация!")
			regErrLabel.Hide()
			regSuccessLabel.Show()
		} else {
			regSuccessLabel.Hide()
			regErrLabel.Show()
		}

	})

	// closeRegBtn.Connect("clicked", func() {
	// 	regSuccessLabel.Hide()
	// 	regErrLabel.Hide()
	// 	loginRegEntry.SetText("")
	// 	passRegEntry.SetText("")
	// 	win2.Hide()

	//})

	// Отображаем все виджеты в окне
	win2.Run()

	// Выполняем главный цикл GTK (для отрисовки). Он остановится когда
	// выполнится gtk.MainQuit()
	gtk.Main()

}

func AutoScroll(scrolledWindow *gtk.ScrolledWindow) {
	adjustment := scrolledWindow.GetVAdjustment()
	adjustment.SetValue(adjustment.GetUpper() - adjustment.GetPageSize())
	adjustment.SetUpper(adjustment.GetPageSize() * 100)
}

func clearListbox(ListBox *gtk.ListBox) {
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


func GetMsgTest(conn *tarantool.Connection, msgListBox *gtk.ListBox) {
	var lastTimedMsg uint64

	fmt.Println("len(messagesArr) ",len(messagesArr))
	if len(messagesArr) == 0 {
		clearListbox(msgListBox)
		lastTimedMsg = 0 // В самом начале загружаем все сообщения
	} else {
		lastTimedMsg = messagesArr[len(messagesArr)-1].time
	}

	infoTimedMsg, _ := conn.Call("fn.time_group_msg", []interface{}{GroupName, lastTimedMsg})
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

			newMsg := msgUser + "(" + GroupName + "): " + msgText

			messagesArr = append(messagesArr, TimedMsg{msg: newMsg, time: msgTime})

			//listbox

			rowMsg, _ := gtk.ListBoxRowNew()
			labelMsg, _ := gtk.LabelNew(newMsg)
			labelMsg.SetHAlign(gtk.ALIGN_START)
			labelMsg.SetJustify(gtk.JUSTIFY_CENTER)
			rowMsg.Add(labelMsg)
			msgListBox.Insert(rowMsg, -1)

		}
		
		msgListBox.ShowAll()

	}
}

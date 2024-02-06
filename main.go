package main

import (
	"GoTarantool/Server"
	"fmt"
	"github.com/gotk3/gotk3/gtk"
	"github.com/tarantool/go-tarantool"
	"log"
)

var MyUser string
var myPass string
var GroupName string
var GroupId string
var userId string

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
	obj2, err := b.GetObject("regWin")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}

	// Преобразуем из объекта именно окно типа gtk.Window
	// и соединяем с сигналом "destroy" чтобы можно было закрыть
	// приложение при закрытии окна
	win := obj.(*gtk.Window)
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	win2 := obj2.(*gtk.Window)
	win2.Connect("delete-event", func() {
		win2.Hide()
	})

	// Получаем поле ввода
	obj, _ = b.GetObject("login_entry")
	loginEntry := obj.(*gtk.Entry)

	obj, _ = b.GetObject("password_entry")
	passEntry := obj.(*gtk.Entry)

	obj, _ = b.GetObject("msg_entry")
	msgEntry := obj.(*gtk.Entry)

	// Получаем кнопку
	obj, _ = b.GetObject("login_btn")
	loginBtn := obj.(*gtk.Button)

	obj, _ = b.GetObject("reg_btn")
	regBtn := obj.(*gtk.Button)

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

	obj2, _ = b.GetObject("close_reg_btn")
	closeRegBtn := obj2.(*gtk.Button)

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

	groupsListbox.Connect("row-activated", func() {
		AutoScroll(scrolledWindow)
		selectedRow := groupsListbox.GetSelectedRow()
		labelRow, _ := selectedRow.GetChild()
		groupLabel := labelRow.(*gtk.Label)
		textLabel, _ := groupLabel.GetText()
		fmt.Println(textLabel)
		GroupName = textLabel
		guildLabel.SetText(GroupName)
		GetMsgTest(conn, msgBox)
		AutoScroll(scrolledWindow)
		msgBox.SetText("")
		messagesArr = messagesArr[:0]
		GetMsgTest(conn, msgBox)
		AutoScroll(scrolledWindow)
		win.ShowAll()

		//t := time.NewTimer(1 * time.Second)
		//go func() {
		//	for {
		//
		//		t.Reset(1 * time.Second)
		//		GetMsgTest(conn, msgBox)
		//		<-t.C
		//	}
		//
		//}()

	})

	// Сигнал по нажатию на кнопку

	loginBtn.Connect("clicked", func() {

		if err == nil {
			// Устанавливаем текст из поля ввода метке
			newUser, _ := loginEntry.GetText()
			newPass, _ := passEntry.GetText()

			children := groupsListbox.GetChildren()

			for children.Length() > 0 {
				child := children.NthData(0)
				if widget, ok := child.(*gtk.Widget); ok {
					groupsListbox.Remove(widget)
				}
				children = groupsListbox.GetChildren()
			}
			win.ShowAll()

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
					rowGroup.Add(labelGroup)
					groupsListbox.Insert(rowGroup, 0)
				}
				win.ShowAll()

			} else {
				fmt.Println("Неверный логин или пароль")
			}

			//tuples := info.Tuples()
			//userId = tuples[0][0].(string)
			//GroupId = tuples[1][0].(string)
			//info, _ = conn.Call("fn.user_guild", []interface{}{MyUser})
			//tuples = info.Tuples()
			//GroupName = tuples[0][0].(string)
			//guildLabel.SetText(GroupName)
			//msgBox.SetText("")
			//messagesArr = messagesArr[:0]
			//GetMsgTest(conn, msgBox)
			////GetMsg(conn, msgBox)
			//AutoScroll(scrolledWindow)
			//
			//t := time.NewTimer(1 * time.Second)
			//if isLogin == false {
			//	isLogin = true
			//	go func() {
			//		for {
			//
			//			t.Reset(1 * time.Second)
			//			//GetMsg(conn, msgBox)
			//			GetMsgTest(conn, msgBox)
			//			AutoScroll(scrolledWindow)
			//			<-t.C
			//		}
			//
			//	}()
			//
			//	// Увеличиваем счетчик wait group на 1
			//}
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
		GetMsgTest(conn, msgBox)
		AutoScroll(scrolledWindow)

	})

	regBtn.Connect("clicked", func() {
		win2.ShowAll()

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

	closeRegBtn.Connect("clicked", func() {
		regSuccessLabel.Hide()
		regErrLabel.Hide()
		loginRegEntry.SetText("")
		passRegEntry.SetText("")
		win2.Hide()

	})

	// Отображаем все виджеты в окне
	win.ShowAll()

	// Выполняем главный цикл GTK (для отрисовки). Он остановится когда
	// выполнится gtk.MainQuit()
	gtk.Main()

}

func AutoScroll(scrolledWindow *gtk.ScrolledWindow) {
	adjustment := scrolledWindow.GetVAdjustment()
	adjustment.SetValue(adjustment.GetUpper() - adjustment.GetPageSize())
	adjustment.SetUpper(adjustment.GetPageSize() * 100)
	fmt.Println(adjustment.GetUpper(), adjustment.GetPageSize(), adjustment.GetLower())
}

func GetMsg(conn *tarantool.Connection, msgBox *gtk.Label) {
	info, _ := conn.Call("fn.guild_msg", []interface{}{GroupId})
	messagesArr = messagesArr[:0]
	messages := info.Tuples()
	if len(messages[0]) != 0 {
		allMsg := ""
		for i := range messages {
			msgText := messages[i][0].(string)
			msgUserId := messages[i][1]
			msgTime := messages[i][2].(uint64)
			msgUserNameTuples, _ := conn.Call("fn.get_name", []interface{}{msgUserId})
			msgUserName := msgUserNameTuples.Tuples()[0][0].(string)

			newMsg := msgUserName + "(" + GroupName + "): " + msgText

			messagesArr = append(messagesArr, TimedMsg{msg: newMsg, time: msgTime})

			allMsg = allMsg + newMsg + "\n"
		}
		msgBox.SetText(allMsg)
		fmt.Println(messagesArr)

		lastTimedMsg := messagesArr[len(messagesArr)-1].time
		fmt.Println(lastTimedMsg)
		infoTimedMsg, _ := conn.Call("fn.time_guild_msg", []interface{}{lastTimedMsg})
		newMessages := infoTimedMsg.Tuples()
		if len(newMessages[0]) != 0 {
			for i := range newMessages {
				fmt.Println(newMessages[i][1].(string))
			}
		}

	} else {
		msgBox.SetText("")
		messagesArr = messagesArr[:0]
	}
}

func GetMsgTest(conn *tarantool.Connection, msgBox *gtk.Label) {
	var lastTimedMsg uint64
	if len(messagesArr) == 0 {
		lastTimedMsg = 0 // В самом начале загружаем все сообщения
	} else {
		lastTimedMsg = messagesArr[len(messagesArr)-1].time
	}
	fmt.Println(lastTimedMsg)
	fmt.Println(GroupName)
	infoTimedMsg, _ := conn.Call("fn.time_group_msg", []interface{}{GroupName, lastTimedMsg})
	fmt.Println("Функция")
	newMessagesCntTuples := infoTimedMsg.Tuples()
	fmt.Println(".")
	fmt.Println(newMessagesCntTuples)
	cntMsg := int(newMessagesCntTuples[0][0].(uint64))
	fmt.Println(".")
	fmt.Println(lastTimedMsg)

	var newMessages []MessageStruct
	for i := 0; i < cntMsg; i++ {
		newMessagesTuples := newMessagesCntTuples[1][i].([]interface{})
		newMessages = append(newMessages, MessageStruct{newMessagesTuples[0].(string), newMessagesTuples[1].(string), newMessagesTuples[2].(uint64)})
	}
	fmt.Println(newMessages)
	fmt.Println("Кол-во", cntMsg)

	if cntMsg != 0 {
		allMsg := ""
		for i := 0; i < cntMsg; i++ {
			msgText := newMessages[i].message
			msgUser := newMessages[i].user
			msgTime := newMessages[i].msgTime
			fmt.Println("1test")

			//msgUserNameTuples, _ := conn.Call("fn.get_name", []interface{}{msgUser}) //ОШИБКА
			//msgUserName := msgUserNameTuples.Tuples()[0][0].(string)

			newMsg := msgUser + "(" + GroupName + "): " + msgText

			messagesArr = append(messagesArr, TimedMsg{msg: newMsg, time: msgTime})

			allMsg = allMsg + newMsg + "\n"
		}
		tText, _ := msgBox.GetText()
		msgBox.SetText(tText + allMsg)
		fmt.Println("Массив", messagesArr)

	}
}

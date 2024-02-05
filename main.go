package main

import (
	"GoTarantool/Server"
	"fmt"
	"github.com/gotk3/gotk3/gtk"
	"github.com/tarantool/go-tarantool"
	"log"
)

var myUser string
var myPass string
var guildName string
var guildId string
var userId string

type TimedMsg struct {
	msg  string
	time uint64
}

type MessageStruct struct {
	message string
	userId  string
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

	//obj, _ = b.GetObject("guild_label")
	//guildLabel := obj.(*gtk.Label)

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

	// Сигнал по нажатию на кнопку

	loginBtn.Connect("clicked", func() {

		if err == nil {
			// Устанавливаем текст из поля ввода метке
			myUser, _ = loginEntry.GetText()
			myPass, _ = passEntry.GetText()
			info, _ := conn.Call("fn.login", []interface{}{myUser, myPass})
			fmt.Println(info)
			userTuples := info.Tuples()
			if userTuples[0][0].(bool) == true {
				userLabel.SetText(myUser)
			} else {
				fmt.Println("Неверный логин или пароль")
			}

			//tuples := info.Tuples()
			//userId = tuples[0][0].(string)
			//guildId = tuples[1][0].(string)
			//info, _ = conn.Call("fn.user_guild", []interface{}{myUser})
			//tuples = info.Tuples()
			//guildName = tuples[0][0].(string)
			//guildLabel.SetText(guildName)
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
		fmt.Println(groupName)
		rowGroup, _ := gtk.ListBoxRowNew()
		labelGroup, _ := gtk.LabelNew(groupName)
		rowGroup.Add(labelGroup)
		groupsListbox.Insert(rowGroup, 0)
		win.ShowAll()

	})

	msgBtn.Connect("clicked", func() {
		newMsg, _ := msgEntry.GetText()
		_, _ = conn.Call("fn.new_msg", []interface{}{newMsg, guildId, userId})
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
	scrolledWindow.Connect("size-allocate", func() {
		adjustment := scrolledWindow.GetVAdjustment()
		adjustment.SetValue(adjustment.GetUpper() - adjustment.GetPageSize())
	})
}

func GetMsg(conn *tarantool.Connection, msgBox *gtk.Label) {
	info, _ := conn.Call("fn.guild_msg", []interface{}{guildId})
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

			newMsg := msgUserName + "(" + guildName + "): " + msgText

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
	infoTimedMsg, _ := conn.Call("fn.time_guild_msg", []interface{}{guildId, lastTimedMsg})
	newMessagesCntTuples := infoTimedMsg.Tuples()
	cntMsg := int(newMessagesCntTuples[0][0].(uint64))
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
			msgUserId := newMessages[i].userId
			msgTime := newMessages[i].msgTime

			msgUserNameTuples, _ := conn.Call("fn.get_name", []interface{}{msgUserId}) //ОШИБКА
			msgUserName := msgUserNameTuples.Tuples()[0][0].(string)

			newMsg := msgUserName + "(" + guildName + "): " + msgText

			messagesArr = append(messagesArr, TimedMsg{msg: newMsg, time: msgTime})

			allMsg = allMsg + newMsg + "\n"
		}
		tText, _ := msgBox.GetText()
		msgBox.SetText(tText + allMsg)
		fmt.Println("Массив", messagesArr)

	}
}

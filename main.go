package main

import (
	"GoTarantool/Server"
	"fmt"
	"github.com/gotk3/gotk3/gtk"
	"github.com/tarantool/go-tarantool"
	"log"
	"time"
)

var myUser string
var guildName string
var guildId uint64
var userId uint64

type TimedMsg struct {
	msg  string
	time uint64
}

var messagesArr = make([]TimedMsg, 0)

func main() {
	isLogin := false

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

	// Преобразуем из объекта именно окно типа gtk.Window
	// и соединяем с сигналом "destroy" чтобы можно было закрыть
	// приложение при закрытии окна
	win := obj.(*gtk.Window)
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	// Получаем поле ввода
	obj, _ = b.GetObject("login_entry")
	loginEntry := obj.(*gtk.Entry)

	obj, _ = b.GetObject("msg_entry")
	msgEntry := obj.(*gtk.Entry)

	// Получаем кнопку
	obj, _ = b.GetObject("login_btn")
	loginBtn := obj.(*gtk.Button)

	obj, _ = b.GetObject("msg_btn")
	msgBtn := obj.(*gtk.Button)

	// Получаем метку
	obj, _ = b.GetObject("msg_box")
	msgBox := obj.(*gtk.Label)
	msgBox.SetText("")

	obj, _ = b.GetObject("guild_label")
	guildLabel := obj.(*gtk.Label)

	obj, _ = b.GetObject("msg_scroll")
	scrolledWindow := obj.(*gtk.ScrolledWindow)

	// Сигнал по нажатию на кнопку

	loginBtn.Connect("clicked", func() {

		if err == nil {
			// Устанавливаем текст из поля ввода метке
			myUser, _ = loginEntry.GetText()
			info, _ := conn.Call("mm.login", []interface{}{myUser})
			tuples := info.Tuples()
			userId = tuples[0][0].(uint64)
			guildId = tuples[1][0].(uint64)
			info, _ = conn.Call("mm.user_guild", []interface{}{myUser})
			tuples = info.Tuples()
			guildName = tuples[0][0].(string)
			guildLabel.SetText(guildName)
			GetMsgTest(conn, msgBox)
			//GetMsg(conn, msgBox)
			AutoScroll(scrolledWindow)

			t := time.NewTimer(1 * time.Second)
			if isLogin == false {
				isLogin = true
				go func() {
					for {

						t.Reset(1 * time.Second)
						//GetMsg(conn, msgBox)
						GetMsgTest(conn, msgBox)
						AutoScroll(scrolledWindow)
						<-t.C
					}

				}()

				// Увеличиваем счетчик wait group на 1
			}
		}

	})

	msgBtn.Connect("clicked", func() {
		newMsg, _ := msgEntry.GetText()
		_, _ = conn.Call("mm.new_msg", []interface{}{newMsg, guildId, userId})
		GetMsg(conn, msgBox)
		AutoScroll(scrolledWindow)

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
	info, _ := conn.Call("mm.guild_msg", []interface{}{guildId})
	messagesArr = messagesArr[:0]
	messages := info.Tuples()
	if len(messages[0]) != 0 {
		allMsg := ""
		for i := range messages {
			msgText := messages[i][0].(string)
			msgUserId := messages[i][1]
			msgTime := messages[i][2].(uint64)
			msgUserNameTuples, _ := conn.Call("mm.get_name", []interface{}{msgUserId})
			msgUserName := msgUserNameTuples.Tuples()[0][0].(string)

			newMsg := msgUserName + "(" + guildName + "): " + msgText

			messagesArr = append(messagesArr, TimedMsg{msg: newMsg, time: msgTime})

			allMsg = allMsg + newMsg + "\n"
		}
		msgBox.SetText(allMsg)
		fmt.Println(messagesArr)

		lastTimedMsg := messagesArr[len(messagesArr)-1].time
		fmt.Println(lastTimedMsg)
		infoTimedMsg, _ := conn.Call("mm.time_guild_msg", []interface{}{lastTimedMsg})
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
		lastTimedMsg = 0
	} else {
		lastTimedMsg = messagesArr[len(messagesArr)-1].time
	}
	infoTimedMsg, _ := conn.Call("mm.time_guild_msg", []interface{}{lastTimedMsg})
	newMessages := infoTimedMsg.Tuples()

	if len(newMessages[0]) != 0 {
		allMsg := ""
		fmt.Println(newMessages[0])
		for i := range newMessages {
			msgText := newMessages[i][1].(string)
			msgUserId := newMessages[i][0]
			msgTime := newMessages[i][4].(uint64)

			fmt.Println(msgTime)

			msgUserNameTuples, _ := conn.Call("mm.get_name", []interface{}{msgUserId}) //ОШИБКА
			msgUserName := msgUserNameTuples.Tuples()[0][0].(string)
			fmt.Println("Проверка")

			newMsg := msgUserName + "(" + guildName + "): " + msgText

			messagesArr = append(messagesArr, TimedMsg{msg: newMsg, time: msgTime})

			allMsg = allMsg + newMsg + "\n"
		}
	}
}

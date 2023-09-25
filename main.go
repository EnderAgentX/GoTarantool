package main

import (
	"GoTarantool/Server"
	"fmt"
	"github.com/gotk3/gotk3/gtk"
	"log"
	"time"
)

func main() {
	var myUser string
	var guildName string
	var guildId uint64
	var userId uint64
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
			print(myUser)
			info, _ := conn.Call("mm.login", []interface{}{myUser})
			tuples := info.Tuples()
			userId = tuples[0][0].(uint64)
			guildId = tuples[1][0].(uint64)
			fmt.Println(tuples)
			info, _ = conn.Call("mm.user_guild", []interface{}{myUser})
			tuples = info.Tuples()
			guildName = tuples[0][0].(string)
			guildLabel.SetText(guildName)
			info, _ = conn.Call("mm.guild_msg", []interface{}{guildId})
			fmt.Println(info)
			messages := info.Tuples()
			allMsg := ""

			for i := range messages {
				msgText := messages[i][0].(string)
				msgUserId := messages[i][1]
				msgUserNameTuples, _ := conn.Call("mm.get_name", []interface{}{msgUserId})
				msgUserName := msgUserNameTuples.Tuples()[0][0].(string)

				newMsg := msgUserName + "(" + guildName + "): " + msgText
				allMsg = allMsg + newMsg + "\n"

			}
			//allMsg = myUser + "(" + guildName + "): " + allMsg
			msgBox.SetText(allMsg)
			scrolledWindow.Connect("size-allocate", func() {
				adjustment := scrolledWindow.GetVAdjustment()
				if err != nil {
					log.Fatal("Ошибка при получении VAdjustment:", err)
				}

				adjustment.SetValue(adjustment.GetUpper() - adjustment.GetPageSize())
			})
			fmt.Println(allMsg)

			t := time.NewTimer(1 * time.Second)
			if isLogin == false {
				isLogin = true
				go func() {
					for {

						t.Reset(1 * time.Second)
						info, _ = conn.Call("mm.guild_msg", []interface{}{guildId})
						fmt.Println(info)
						messages = info.Tuples()
						//msgBox.SetText("")
						allMsg = ""
						for i := range messages {
							msgText := messages[i][0].(string)
							msgUserId := messages[i][1]
							msgUserNameTuples, _ := conn.Call("mm.get_name", []interface{}{msgUserId})
							msgUserName := msgUserNameTuples.Tuples()[0][0].(string)

							newMsg := msgUserName + "(" + guildName + "): " + msgText
							allMsg = allMsg + newMsg + "\n"
						}
						//allMsg = myUser + "(" + guildName + "): " + allMsg
						msgBox.SetText(allMsg)
						scrolledWindow.Connect("size-allocate", func() {
							adjustment := scrolledWindow.GetVAdjustment()
							if err != nil {
								log.Fatal("Ошибка при получении VAdjustment:", err)
							}

							adjustment.SetValue(adjustment.GetUpper() - adjustment.GetPageSize())
						})
						fmt.Println(allMsg)

						<-t.C
					}

				}()

				// Увеличиваем счетчик wait group на 1
			}
		}

	})

	msgBtn.Connect("clicked", func() {
		newMsg, _ := msgEntry.GetText()
		fmt.Println(newMsg)
		fmt.Println(myUser)
		fmt.Println(guildName)
		_, _ = conn.Call("mm.new_msg", []interface{}{newMsg, guildId, userId})
		info, _ := conn.Call("mm.guild_msg", []interface{}{guildId})
		fmt.Println(info)
		messages := info.Tuples()
		//msgBox.SetText("")
		allMsg := ""
		for i := range messages {
			msgText := messages[i][0].(string)
			msgUserId := messages[i][1]
			msgUserNameTuples, _ := conn.Call("mm.get_name", []interface{}{msgUserId})
			msgUserName := msgUserNameTuples.Tuples()[0][0].(string)

			newMsg := msgUserName + "(" + guildName + "): " + msgText
			allMsg = allMsg + newMsg + "\n"
		}
		//allMsg = myUser + "(" + guildName + "): " + allMsg
		msgBox.SetText(allMsg)
		scrolledWindow.Connect("size-allocate", func() {
			adjustment := scrolledWindow.GetVAdjustment()
			if err != nil {
				log.Fatal("Ошибка при получении VAdjustment:", err)
			}

			adjustment.SetValue(adjustment.GetUpper() - adjustment.GetPageSize())
		})

	})

	// Отображаем все виджеты в окне
	win.ShowAll()

	// Выполняем главный цикл GTK (для отрисовки). Он остановится когда
	// выполнится gtk.MainQuit()
	gtk.Main()

}

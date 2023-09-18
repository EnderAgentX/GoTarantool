package main

import (
	"GoTarantool/Server"
	"bufio"
	"fmt"
	"github.com/gotk3/gotk3/gtk"
	"log"
	"os"
	"time"
)

func main() {
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

	// Получаем кнопку
	obj, _ = b.GetObject("login_btn")
	loginBtn := obj.(*gtk.Button)

	// Получаем метку
	obj, _ = b.GetObject("msg_box")
	msgBox := obj.(*gtk.Label)
	msgBox.SetText("")

	obj, _ = b.GetObject("guild_label")
	guildLabel := obj.(*gtk.Label)

	textLabel, _ := msgBox.GetText()

	t := time.NewTimer(3 * time.Second)
	go func() {
		for {
			t.Reset(1 * time.Second)
			textLabel, _ = msgBox.GetText()
			msgBox.SetText(textLabel + "*" + "\n")
			<-t.C
		}
	}()

	// Сигнал по нажатию на кнопку
	loginBtn.Connect("clicked", func() {

		if err == nil {
			// Устанавливаем текст из поля ввода метке
			myUser, _ := loginEntry.GetText()
			textLabel, _ = msgBox.GetText()
			print(myUser)
			info, _ := conn.Call("mm.login", []interface{}{myUser})
			tuples := info.Tuples()
			//userId := tuples[0][0]
			//guildId := tuples[1][0]
			//fmt.Println(userId)
			info, _ = conn.Call("mm.user_guild", []interface{}{myUser})
			tuples = info.Tuples()
			guildName := tuples[0][0].(string)
			guildLabel.SetText(guildName)

		}

	})

	// Отображаем все виджеты в окне
	win.ShowAll()

	// Выполняем главный цикл GTK (для отрисовки). Он остановится когда
	// выполнится gtk.MainQuit()
	gtk.Main()

	/////////////////////////////////////////////////////////////////////////////

	_, _ = conn.Call("mm.insertAll", []interface{}{})

	myscanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Введите имя пользователя: ")
	myscanner.Scan()
	myUser := myscanner.Text()
	fmt.Println(myUser)

	info, _ := conn.Call("mm.login", []interface{}{myUser})
	tuples := info.Tuples()
	userId := tuples[0][0]
	guildId := tuples[1][0]
	//fmt.Println(userId)
	guildName, _ := conn.Call("mm.user_guild", []interface{}{myUser})

	for {
		fmt.Print("Введите сообщение: ")
		myscanner.Scan()
		msg := myscanner.Text()
		if msg == "login" {
			fmt.Print("Введите имя пользователя: ")
			myscanner.Scan()
			myUser = myscanner.Text()
			fmt.Println(myUser)
			info, _ = conn.Call("mm.login", []interface{}{myUser})
			tuples = info.Tuples()
			userId = tuples[0][0]
			guildId = tuples[1][0]
			//fmt.Println(userId)
			guildName, _ = conn.Call("mm.user_guild", []interface{}{myUser})
		} else {
			fmt.Printf("%s(%s): %s", myUser, guildName.Tuples()[0][0], msg)
			fmt.Println("")
			_, _ = conn.Call("mm.new_msg", []interface{}{msg, guildId, userId})
			//fmt.Println(funcres)
		}
	}

}

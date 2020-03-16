package main

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/websocket"
	//"github.com/go-xorm/xorm"
	"xorm.io/xorm"
	_ "github.com/mattn/go-sqlite3"
    "encoding/json"
)


type User struct {
	ID		int64	`xorm:"not null pk autoincr INT(10)"`
    JH		string	`json:"警号" xorm:"varchar(50) not null unique"`
	XM		string	`json:"姓名" xorm:"varchar(50) not null"`
	MM		string	`json:"密码" xorm:"varchar(50) not null"`
    //Version string  `xorm:"varchar(200)"`
// 	ID		int64	`json:"id" xorm:"not null pk autoincr INT(10)"`
//	Version   string `xorm:"varchar(200)"`
//	Salt      string
//	Username  string
//	Password  string    `xorm:"varchar(200)"`
//	Languages string    `xorm:"varchar(200)"`
//	CreatedAt time.Time `xorm:"created"`
//	UpdatedAt time.Time `xorm:"updated"`
}
//
//type User1 struct {
//	//id		int64	`xorm:"not null pk autoincr INT(10)"`
//    JH		string	`json:"警号" xorm:"varchar(50) not null unique"`
//	XM		string	`json:"姓名" xorm:"varchar(50) not null"`
//	MM		string	`json:"密码" xorm:"varchar(50) not null"`
//}

func main() {
	app := iris.Default()

	app.Get("/pong", func(ctx iris.Context) {
		ctx.JSON(iris.Map{
			"message": "pong吴志伟",
		})
	})
	
	/////////////////////////////////////////////////////////////////
	///////////这里是真正的代码, 后面是参考代码 //////////////////////////
	/////////////////////////////////////////////////////////////////
	orm, err := xorm.NewEngine("sqlite3", "./test.db")
	if err != nil {
		app.Logger().Fatalf("orm failed to initialized: %v", err)
	}

	iris.RegisterOnInterrupt(func() {
		orm.Close()
	})
	
	//在控制台打印sql语句，默认为false
    orm.ShowSQL(true)

	err = orm.Sync2(new(User))
	//* 自动检测和创建表，这个检测是根据表的名字

	if err != nil {
		app.Logger().Fatalf("orm failed to initialized User table: %v", err)
	}	
	
	//网页上显示用户信息
	app.Get("/getusers", func(ctx iris.Context) {
        //users:= make([]User, 0)
        orm.Iterate(new(User), func(i int, bean interface{})error{
            user := bean.(*User)
            ctx.Writef("<H1></H1><H3>")
            ctx.JSON(user)
            ctx.Writef("</H3>")
            return nil
        })
	})
	
	//发送JSON格式的用户信息,以便EXCEL文件导出
	app.Post("/getusers", func(ctx iris.Context) {
		//ctx.Gzip(true)               // enable gzip for big files
		users:= make([]User, 0)
		err := orm.Find(&users)
		if err != nil{
			ctx.JSON(err)
		}else{
			//ctx.Writef("%#v", users)
			ctx.JSON(users)
		}
		
//		ctx.JSON(iris.Map{
//			"status":  "posted",
//			"message": "message",
//			"nick":    "nick",
//		})
	})
	
	//重定向到exportusers.html网页, 以便提供一致的用户体验
	app.Get("/exportusers", func(ctx iris.Context) {
		ctx.ServeFile("./web/exportusers.html", false) // true for gzip.
        //        ctx.ServeFile("./html/index.html", true) // true for gzip.

	})
	
	//重定向到importusers.html网页, 以便提供一致的用户体验
	app.Get("/importusers", func(ctx iris.Context) {
		ctx.ServeFile("./web/importusers.html", false) // true for gzip.
        //        ctx.ServeFile("./html/index.html", true) // true for gzip.

	})
	
	//接收JSON格式的用户信息,导入数据库
	app.Post("/importusers", func(ctx iris.Context) {
		//ctx.Gzip(true)               // enable gzip for big files
        
		msg := ctx.FormValue("msg")
		fmt.Printf("接收到的用户EXCEL表: %v\n", msg)
        
        //str := `{"page": 1, "fruits": ["apple", "peach"]}`
        users := &[]User{}
        if err := json.Unmarshal([]byte(msg), &users); err != nil {
			
			msg := fmt.Sprintf("%v", err) 
			msg = "解析错误:" + msg
			fmt.Printf("%v\n", msg)
			ctx.JSON(iris.Map{
				"msg": msg,	//这也是在$ajax.success中的
			})
			//ctx.JSON(msg)	//这样也可以, 只要接收端对应即可
			
			return
        }
        //fmt.Println(users)
		fmt.Printf("导入用户EXCEL表解析结果: %v\n", *users  )
		
		for _, v := range *users {
			//fmt.Println(v)
			user := &User{
				ID: v.ID,
				JH: v.JH,
				XM: v.XM,
				MM: v.MM}
			_,err := orm.Insert(user) //<--------插入
			if err != nil {
				msg := fmt.Sprintf("导入失败: %v", err) 
				fmt.Println( msg)			
				ctx.JSON(iris.Map{		
					"msg": msg,
				})					
				return
			}
		}
		
		
		
		
		
		ctx.JSON(iris.Map{		
			"msg": "导入成功",
		})		
        

//        for i, v := range users {
//            _ = i
//            fmt.Printf("%v\n", v)
//        }
		
//		var users []User
//		err := ctx.ReadJSON(&users)
//        fmt.Printf("%v\n", users)
//		if err != nil {
//            fmt.Printf("err\n")
//			ctx.StatusCode(iris.StatusBadRequest)
//			ctx.WriteString(err.Error())
//			return
//		}
        fmt.Printf("ok\n")
		//ctx.Writef("导入成功!%#+v\n",users)

	})
	
	//添加用户
	app.Get("/adduser/{jh:string}/{xm:string}", func(ctx iris.Context) {
		jh := ctx.Params().Get("jh")
		xm := ctx.Params().Get("xm")
		
		user := &User{
			JH: jh,
			XM: xm,
			MM: jh}
        _,err := orm.Insert(user) //<--------插入
        if err != nil {
            ctx.Writef("<H1></H1>")
            ctx.Writef("<H3>%#v</H3>", err)
        }else{		
            ctx.Writef("<H1></H1>")			//没有这一行,下面的<H3>格式显示错误
            ctx.Writef("<h3>添加用户成功, 请使用警号登录,初始密码即警号!</h3>")
        }
	})
	
	//删除某个用户
	app.Get("/deluser/{jh:string}", func(ctx iris.Context) {
		jh := ctx.Params().Get("jh")
		//xm := ctx.Params().Get("xm")
		
		user := &User{
			JH: jh,
			//XM: xm,
			//MM: jh
		}
		
		n, err := orm.Delete(user)	//<--------删除
        if err != nil {
            ctx.Writef("<H1></H1>")
            ctx.Writef("<H3>%#v</H3>", err)
		}else if n==0 {
            ctx.Writef("<H1></H1>")			//没有这一行,下面的<H3>格式显示错误
            ctx.Writef("<h3>没有找到!</h3>")			
		}else{		
            ctx.Writef("<H1></H1>")			//没有这一行,下面的<H3>格式显示错误
            ctx.Writef("<h3>删除成功!</h3>")
        }
	})
	
	//删除全部用户
	app.Get("/delallusers", func(ctx iris.Context) {
            ctx.ServeFile("./web/delallusers.html", false)//不压缩
        
	})
	
	app.Post("/delallusers", func(ctx iris.Context) {
		msg := ctx.FormValue("msg")
		fmt.Printf("接收到msg: %v\n", msg)   
		
		err := orm.DropTables(new(User))
		fmt.Println(err)
		err = orm.Sync2(new(User))
		if err != nil {
			app.Logger().Fatalf("orm failed to initialized User table: %v", err)
		}	
		
//		orm.Iterate(new(User), func(i int, bean interface{})error{
//            user := bean.(*User)
//			fmt.Println(user.JH)
//			deluser := &User{
//				JH: "000000",
//			}
//			affected, err := orm.Delete(deluser)
//			msg := fmt.Sprintf("删除了: %v 个, err: %v", affected, err) 
//			fmt.Println(msg)
//            return nil
//        })
		
	
		ctx.JSON("全部删除成功")
        
	})
	/////////////////////////////////////////////////////////////////
	///////////上面是真正的代码, 后面是参考代码 //////////////////////////
	/////////////////////////////////////////////////////////////////
	// http://localhost:8080/pong

	//Parameters in path
	// This handler will match /user/kataras but will not match neither /user/ or /user.
	app.Get("/user/{name}", func(ctx iris.Context) {
		name := ctx.Params().Get("name")
		ctx.Writef("<H1>Hello %s<H1>", name)
	})
	// http://localhost:8080/user/吴志伟
	// curl "http://localhost:8080/user/吴志伟"

	// This handles the /user/kataras/42
	// and fires 400 bad request if /user/kataras/string.
	// The "else 400" is optionally:
	// by-default it will fire 404 not found if alphanumeric instead
	// of number passed on the "age" parameter.
	app.Get("/user/{name:string}/{age:int else 400}", func(ctx iris.Context) {
		name := ctx.Params().Get("name")
		age, _ := ctx.Params().GetInt("age")
		ctx.Writef("%s is %d years old", name, age)
	})
	// http://localhost:8080/user/吴志伟/26
	// curl "http://localhost:8080/user/吴志伟/26"

	// However, this one will match /action/{user}/star and also /action/{user}/stars
	// or even /action/{user}/likes/page/2.
	// It should match anything after the /action/{user}/
	// except the /action/{user}/static which is handled by the below route.
	app.Get("/action/{user:string}/{action:path}", func(ctx iris.Context) {
		user := ctx.Params().Get("user")
		action := ctx.Params().Get("action")
		ctx.Writef("user: %s | action: %s", user, action)
	})
	// http://localhost:8080/action/吴志伟/26/上海/江苏

	// Unlike other frameworks and routers,
	// Iris is smart enough to understand that this is not the previous,
	// wildcard of type path route, it should only match the /action/{user}/static.
	app.Get("/action/{user:string}/static", func(ctx iris.Context) {
		user := ctx.Params().Get("user")
		ctx.Writef("static path for user: %s", user)
	})
	// http://localhost:8080/user/kataras
	// http://localhost:8080/user/kataras/25
	// http://localhost:8080/action/kataras/upgrade
	// http://localhost:8080/action/kataras/static
	// curl "http://localhost:8080/action/kataras/static"

	//Querystring parameters
	// Query string parameters are parsed using the existing underlying request object.
	// The request responds to a url matching:  /welcome?firstname=Jane&lastname=Doe.
	app.Get("/welcome", func(ctx iris.Context) {
		firstname := ctx.URLParamDefault("firstname", "Guest")
		// shortcut for ctx.Request().URL.Query().Get("lastname").
		lastname := ctx.URLParam("lastname")

		ctx.Writef("Hello %s %s", firstname, lastname)
	})
	//http://localhost:8080/welcome?firstname=Jane&lastname=Doe
	//或者 curl "http://localhost:8080/welcome?firstname=Jane&lastname=Doe"

	//////////////////////////////////////////////////////////////////////////////////////////////////
	//Multipart/Urlencoded Form
	app.Get("/testpost", func(ctx iris.Context) {
		ctx.HTML(`
		<html>
			<form action='/form_post' method='post'>
			    message: <input type="text" name="message">
   				nick: <input type="text" name="nick">
	 
	   			 <input type='submit' id='submit'>
			</form>
		</html>`)
	})
	app.Post("/form_post", func(ctx iris.Context) {
		message := ctx.FormValue("message")
		nick := ctx.FormValueDefault("nick", "anonymous")

		ctx.JSON(iris.Map{
			"status":  "posted",
			"message": message,
			"nick":    nick,
		})
	})
	//http://localhost:8080/testpost
	//或者 curl -X POST "http://localhost:8080/form_post" -F message=吴志伟 -F nick=wzw

	//Another example: query + post form
	app.Get("/testpostquery", func(ctx iris.Context) {
		ctx.HTML(`
		<html>
			<form action='/form_postquery?id=123&page=321' method='post'>
			    message: <input type="text" name="message">
   				name: <input type="text" name="name">
	 
	   			 <input type='submit' id='submit'>
			</form>
		</html>`)
	})
	app.Post("/form_postquery", func(ctx iris.Context) {
		id := ctx.URLParam("id")
		page := ctx.URLParamDefault("page", "0")
		name := ctx.FormValue("name")
		message := ctx.FormValue("message")
		// or `ctx.PostValue` for POST, PUT & PATCH-only HTTP Methods.

		app.Logger().Infof("id: %s; page: %s; name: %s; message: %s", id, page, name, message)
	})
	//http://localhost:8080/testpostquery
	//或者 curl -X POST -F message=hello -F name=吴志伟   "http://localhost:8080/form_postquery?page=abc&id=111"

	// Listen and serve on http://localhost:8080.

	////// *.html模板的使用 *.js /////////////////
	app.RegisterView(iris.HTML("./views", ".html"))

	app.Get("/mb", func(ctx iris.Context) {

		ctx.ViewData("message", "Hello world!")

		ctx.View("hello.html")
	})

	app.RegisterView(iris.HTML("./js", ".js"))

	app.Get("/app.js", func(ctx iris.Context) {

		ctx.View("app.js")
	})
	//http://localhost:8080/mb

/////////////一个静态网站,包含子目录都可以自动路由寻址//////////////////////
/////////////一个静态网站,包含子目录都可以自动路由寻址//////////////////////
/////////////一个静态网站,包含子目录都可以自动路由寻址//////////////////////
	app.StaticWeb("/", "./web") //<-----------------------设定网站根目录
        //app.StaticWeb("/", "./html") //<-----------------------设定网站根目录


	app.Get("/", func(ctx iris.Context) {
		ctx.ServeFile("./web/index.html", true) // true for gzip.
        //        ctx.ServeFile("./html/index.html", true) // true for gzip.

	})
	//http://localhost:8080                 (含有js子目录)
	//http://localhost:8080/index.html      (效果同上)
	//http://localhost:8080/hello.html


/////////////websockets 服务器/////////////////////////////////////////////
/////////////websockets 服务器/////////////////////////////////////////////
/////////////websockets 服务器/////////////////////////////////////////////
	app.Get("/ws", func(ctx iris.Context) {
		ctx.ServeFile("./html/websockets.html", false) // second parameter: enable gzip?
	})

	setupWebsocket(app)

	// x2
	// http://localhost:8080
	// http://localhost:8080
	// write something, press submit, see the result.


	fmt.Println("http://localhost:8101");
	fmt.Println("http://localhost:8101/ws");
	app.Run(iris.Addr(":8101"))//<--------------------------------------



}

func setupWebsocket(app *iris.Application) {
	// create our echo websocket server
	ws := websocket.New(websocket.Config{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	})
	ws.OnConnection(handleConnection) //这里只是设置,还没有开始响应

	// register the server on an endpoint.
	// see the inline javascript code in the websockets.html,
	// this endpoint is used to connect to the server.
	app.Get("/echo", ws.Handler()) //这个websockets是架设在ws://IP:8080/echo
	//如果要架设多个websockets服务,可以通过不同的路由设置来完成,即IP:PROT是相同的,只是后面的路由不同. 可以生成多个ws1 := websocket.New(),然后app.Get("/path", ws1.Handler())
	// serve the javascript built'n client-side library,
	// see websockets.html script tags, this path is used.
	app.Any("/iris-ws.js", websocket.ClientHandler()) //返回iris-ws.js文件给浏览器
}

func handleConnection(c websocket.Connection) {
	// Read events from browser
	c.On("chat", func(msg string) {
		// Print the message to the console, c.Context() is the iris's http context.
		fmt.Printf("%s sent: %s\n", c.Context().RemoteAddr(), msg)
		// Write message back to the client message owner with:
		c.Emit("chat", msg) //回给当前客户端
		// Write message to all except this client with:
		c.To(websocket.Broadcast).Emit("chat", msg) //发给所有客户端除了当前客户端
	})
}









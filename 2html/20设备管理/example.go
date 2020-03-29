package main

import (
	"fmt"
	"time"
	"math/rand"
	"github.com/kataras/iris"
	"github.com/kataras/iris/websocket"
	//"github.com/go-xorm/xorm"
	"encoding/json"

	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
)

type User struct {
	ID int64  `xorm:"not null pk autoincr unique INT(10)"`
	JH string `json:"警号" xorm:"varchar(50) not null unique"`
	XM string `json:"姓名" xorm:"varchar(50) not null"`
	MM string `json:"密码" xorm:"varchar(50) not null"`
	//json标签实现的是外部数据与struct之间的映射关系
	//xorm标签实现的是数据库与struct之间的映射关系
}

type Asset struct {
	ID      int64   `xorm:"not null pk autoincr unique INT(10)"`
	XH      string  `json:"序号" xorm:"varchar(50)"`
	BM      string  `json:"部门"` //string 对应 Varchar(255)
	ZCBH    string  `json:"资产编号"`
	ZCMC    string  `json:"资产名称"`
	PP      string  `json:"品牌"`
	GGXH    string  `json:"规格型号"`
	QDRQ    int64   `json:"取得日期"`
	SL      float64 `json:"数量" xorm:"Numeric"` //float64精度太高,不便于定位和查找
	DJ      float64 `json:"单价" xorm:"Numeric"`
	JE      float64 `json:"金额" xorm:"Numeric"`
	SYR     string  `json:"使用人"`
	CFDD    string  `json:"存放地点" xorm:"not null"`
	BZ      string  `json:"备注"`
	QRCODE  string  `json:"二维码" xorm:"not null unique"`
	//Version int     `xorm:"version"` // 乐观锁
}

type Recorder struct{
    ID          int64   `xorm:"not null pk autoincr unique INT(10)"`
    QRCODE      string  `json:"二维码" xorm:"not null"`
    ZCMC        string  `json:"资产名称"`
    CFDD        string  `json:"存放地点" xorm:"not null"`
    JH          string  `json:"警号" xorm:"varchar(50) not null"`
    XM          string  `json:"姓名" xorm:"varchar(50) not null"`
    BZ          string  `json:"备注"`
    CreatedAt   time.Time `json:"记录时间" xorm:"created"`
    
}

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

    //打开 User 数据表
	if err = orm.Sync2(new(User)); err != nil {
		app.Logger().Fatalf("orm failed to initialized User table: %v", err)
	}
    //打开 Asset 数据表 
	if err = orm.Sync2(new(Asset)); err != nil {
		app.Logger().Fatalf("orm failed to initialized Asset table: %v", err)
	}
    //打开 Recorder 数据表 
	if err = orm.Sync2(new(Recorder)); err != nil {
		app.Logger().Fatalf("orm failed to initialized User table: %v", err)
	}
/////////////////////////
//		asset1 := &Asset{CFDD: "311仓库"}
//		isaddr, _ := orm.Exist(asset1)		
//		assets := make([]Asset, 0)
//		err1 := orm.Find(&assets, asset1)
//		_ = err1
//		//fmt.Println("isaddr",isaddr,"assets",assets,"err",err1)	
//////////////////////////	
	
	//定时生成随机数,用于设置 admin 的 cookie
	//两小时换一个密码
	rand.Seed(time.Now().UnixNano())
    cookienum := rand.Intn(9999999999)
	fmt.Println("随机数:  ",cookienum)
    ticker := time.NewTicker(time.Second * 60 * 60 * 2)
    go func() {
        for _ = range ticker.C {
			rand.Seed(time.Now().UnixNano())
			cookienum = rand.Intn(9999999999)
			fmt.Println("随机数:  ",cookienum)
        }
    }()	
	
	// Post: login
	app.Post("/login", func(ctx iris.Context) {
		jh := ctx.FormValue("jh")
		pwd := ctx.FormValue("pwd")
		
		//cookie操作
		//ctx.SetCookieKV(name, value)
		//value := ctx.GetCookie(name)
		//ctx.RemoveCookie(name)
		//定制路径:ctx.SetCookieKV(name, value, iris.CookiePath("/custom/path/cookie/will/be/stored"))

		user := &User{JH: jh}
		has, _ := orm.Get(user)
		//fmt.Printf("找到:%v",has)
		if !has {
			fmt.Printf("用户不存在!\n")
			ctx.JSON(iris.Map{
				"status": "1",
				"msg":    "用户不存在",
			})
			return
		}
		if user.MM != pwd {
			fmt.Printf("密码错误!\n")
			ctx.JSON(iris.Map{
				"status": "2",
				"msg":    "密码错误!",
			})
			return
		}

		//如果是管理员则修改COOKIE
		cookiestr :="0000000000"
		fmt.Println("ctx.GetCookie(name)= ",  ctx.GetCookie("cookienum"))	
		if jh == "000000" {
			cookiestr = fmt.Sprintf("%v",cookienum)
			fmt.Println("set cookienum: ",cookiestr)
		}
		
		ctx.JSON(iris.Map{
			"status": "0",
			"data":    user,
			"cookienum": cookiestr,
		})

	})
	
	// Post: query 扫描设备查询
	app.Post("/query", func(ctx iris.Context) {
		query := ctx.FormValue("dat")

		asset := &Asset{QRCODE: query}
		isasset, _ := orm.Get(asset)		
		if isasset {
            //如果是资产设备还需要查询流转记录		
            recorders := make([]Recorder, 0)
            err := orm.Find(&recorders, &Recorder{QRCODE: query})
            //_ = err  
            fmt.Println("err: ",err,"  recorders: ", recorders)
            //数据回传
			ctx.JSON(iris.Map{
				"status": 	"0",
				"msg":		"扫描的是资产二维码",
				"data":    	asset,
                "recorders":recorders,
			})			
			return
		}		
		
		asset1 := &Asset{CFDD: query}
		isaddr, _ := orm.Exist(asset1)		
		assets := make([]Asset, 0)
		err := orm.Find(&assets, asset1)
		_ = err
		fmt.Println("isaddr",isaddr,"assets",assets,"err",err)
		
		if isaddr {
			ctx.JSON(iris.Map{
				"status": 	"1",
				"msg":		"扫描的是地点二维码",
				"data":    assets,				
			})
			return
		}
		
		ctx.JSON(iris.Map{
			"status": 	"2",
			"msg":		"二维码错误, 请重新扫描!",  				
		})		

	})
	
	// Post: /verifyaddr
	app.Post("/verifyaddr", func(ctx iris.Context) {
		query := ctx.FormValue("dat")
		
		asset1 := &Asset{CFDD: query}
		isaddr, _ := orm.Exist(asset1)		
		fmt.Println("isaddr",isaddr)
		
		if isaddr {
			ctx.JSON(iris.Map{
				"status": 	"1",
				"msg":		"完成设备存放地点改变!",				
			})
			return
		}
		
		ctx.JSON(iris.Map{
			"status": 	"2",
			"msg":		"二维码错误, 请重新扫描!",  				
		})		

	})
	
	// Post: /changeaddr
	app.Post("/changeaddr", func(ctx iris.Context) {
		asset_str:= ctx.FormValue("asset")
		addr_new := ctx.FormValue("addr")
		user_str := ctx.FormValue("user")
		comment  := ctx.FormValue("comment")
		fmt.Println(asset_str,addr_new,user_str,comment)
		
		user := &User{}
		asset:= &Asset{}
		json.Unmarshal([]byte(user_str), user)
		json.Unmarshal([]byte(asset_str), asset)
		fmt.Println(user,asset)
		
		if len(comment)==0{ //备注为空时强制更新
			comment = " "
		}
		
		//1,此处加入事务处理
		session := orm.NewSession()
		defer session.Close()
		err := session.Begin()
		msg := fmt.Sprintf("变更地址失败,请重试!\n%v", err)
		if err != nil{
			ctx.JSON(iris.Map{
				"status": 	"2",
				"msg":		msg, 			
		      })
			return
		}		
		
        //2,写Asset表
		affected, err := session.Id(asset.ID).Update(
			&Asset{CFDD:addr_new,
				  BZ:comment})
        fmt.Println("Update Asset affected: ",affected,"  err:",err)
		if affected != 1 || err != nil {	//失败
			msg := fmt.Sprintf("变更地址失败,请重试!\n%v", err)
			ctx.JSON(iris.Map{
				"status": 	"2",
				"msg":		msg, 			
		      })
			session.Rollback()				//失败则回滚
            return
        }	
                 
        //3,写Recorder表
        affected, err = session.Insert(&Recorder{
            QRCODE: asset.QRCODE,
            ZCMC:   asset.ZCMC,
            CFDD:   addr_new,
            JH:     user.JH,
            XM:     user.XM,
            BZ:     comment,            
        }) 
        fmt.Println("Update Recorder affected: ",affected,"  err:",err)
		if affected != 1 || err != nil {	//失败
			msg := fmt.Sprintf("变更地址失败,请重试!\n%v", err)
			ctx.JSON(iris.Map{
				"status": 	"2",
				"msg":		msg, 			
		      })
			session.Rollback()				//失败则回滚
            return
        }
		
		//4,事务提交
		err = session.Commit()				//事务提交. 一个事务到此结束
		if err != nil{
			msg := fmt.Sprintf("变更地址失败,请重试!\n%v", err)
			ctx.JSON(iris.Map{
				"status": 	"2",
				"msg":		msg, 			
		      })
			session.Rollback()				//失败则回滚
			return
		}
        
        //成功!!!
        ctx.JSON(iris.Map{
            "status": 	"1",
            "msg":		"成功!",				
        })
    })
	
	// Post: changepwd
	app.Post("/changepwd", func(ctx iris.Context) {
		jh := ctx.FormValue("jh")
		pwd := ctx.FormValue("pwd")
		
		user := &User{JH: jh}
		has, _ := orm.Get(user)
		if !has {
			fmt.Printf("用户不存在!\n")
			ctx.JSON(iris.Map{
				"status": "1",
				"msg":    "用户不存在",
			})
			return
		}
		user.MM = pwd
		orm.ID(user.ID).Update(user)
		ctx.JSON(iris.Map{
			"status": "0",
			"msg":    "密码修改成功!",
		})
		
		
	})		

	//网页上显示用户信息
	app.Get("/getusers", func(ctx iris.Context) {
		//操作确权, 如果不是admin,则引导到主页
		getcookienum := ctx.GetCookie("cookienum")
		cookiestr := fmt.Sprintf("%v",cookienum)
		if getcookienum != cookiestr {
			ctx.Redirect("/")
			return
		}
		
		//users:= make([]User, 0)
		orm.Iterate(new(User), func(i int, bean interface{}) error {
			user := bean.(*User)
			ctx.Writef("<H1></H1><H3>")
			ctx.JSON(user)
			ctx.Writef("</H3>")
			return nil
		})
	})

	app.Get("/admin", func(ctx iris.Context) {
		//操作确权, 如果不是admin,则引导到主页
		getcookienum := ctx.GetCookie("cookienum")
		cookiestr := fmt.Sprintf("%v",cookienum)
		if getcookienum != cookiestr {
			ctx.Redirect("/")
			return
		}
				
		ctx.ServeFile("./web/admin.html", false) // true for gzip.
		//        ctx.ServeFile("./html/index.html", true) // true for gzip.

	})
	//发送JSON格式的用户信息,以便EXCEL文件导出
	app.Post("/getusers", func(ctx iris.Context) {
		//操作确权, 如果不是admin,则引导到主页
		getcookienum := ctx.GetCookie("cookienum")
		cookiestr := fmt.Sprintf("%v",cookienum)
		if getcookienum != cookiestr {
			ctx.Redirect("/")
			return
		}
		
		users := make([]User, 0)
		err := orm.Find(&users)
		if err != nil {
			ctx.JSON(err)
		} else {
			//ctx.Writef("%#v", users)
			ctx.JSON(users)
		}

	})

	//重定向到exportusers.html网页, 以便提供一致的用户体验
	app.Get("/exportusers", func(ctx iris.Context) {
		//操作确权, 如果不是admin,则引导到主页
		getcookienum := ctx.GetCookie("cookienum")
		cookiestr := fmt.Sprintf("%v",cookienum)
		if getcookienum != cookiestr {
			ctx.Redirect("/")
			return
		}
		
		ctx.ServeFile("./web/exportusers.html", false) // true for gzip.
		//        ctx.ServeFile("./html/index.html", true) // true for gzip.

	})

	//重定向到importusers.html网页, 以便提供一致的用户体验
	app.Get("/importusers", func(ctx iris.Context) {
		//操作确权, 如果不是admin,则引导到主页
		getcookienum := ctx.GetCookie("cookienum")
		cookiestr := fmt.Sprintf("%v",cookienum)
		if getcookienum != cookiestr {
			ctx.Redirect("/")
			return
		}
		
		ctx.ServeFile("./web/importusers.html", false) // true for gzip.
		//        ctx.ServeFile("./html/index.html", true) // true for gzip.

	})

	//接收JSON格式的用户信息,导入数据库
	app.Post("/importusers", func(ctx iris.Context) {
		//操作确权, 如果不是admin,则引导到主页
		getcookienum := ctx.GetCookie("cookienum")
		cookiestr := fmt.Sprintf("%v",cookienum)
		if getcookienum != cookiestr {
			ctx.Redirect("/")
			return
		}
		

		msg := ctx.FormValue("msg")
		fmt.Printf("接收到的用户EXCEL表: %v\n", msg)

		//str := `{"page": 1, "fruits": ["apple", "peach"]}`
		users := &[]User{}
		if err := json.Unmarshal([]byte(msg), &users); err != nil {

			msg := fmt.Sprintf("%v", err)
			msg = "解析错误:" + msg
			fmt.Printf("%v\n", msg)
			ctx.JSON(iris.Map{
				"msg": msg, //这也是在$ajax.success中的
			})
			//ctx.JSON(msg)	//这样也可以, 只要接收端对应即可

			return
		}
		//fmt.Println(users)
		fmt.Printf("导入用户EXCEL表解析结果: %v\n", *users)

		for i, v := range *users {
			//fmt.Println(v)
			user := &User{
				//ID: v.ID,			//不用管
				JH: v.JH,
				XM: v.XM,
				MM: v.MM}
			_, err := orm.Insert(user) //<--------插入
			if err != nil {
				msg := fmt.Sprintf("导入失败: %v,  第 %v 行", err, i)
				fmt.Println(msg)
				ctx.JSON(iris.Map{
					"msg": msg,
				})
				return
			}
		}

		ctx.JSON(iris.Map{
			"msg": "导入成功",
		})

		fmt.Printf("ok\n")

	})
	
	//按警号查询流转记录
	app.Get("/queryrecorders/{jh:string}", func(ctx iris.Context){		
		jh := ctx.Params().Get("jh")
		
		recorders := make([]Recorder, 0)
		err := orm.Find(&recorders, &Recorder{JH: jh})
        fmt.Println("err: ",err,"  recorders: ", recorders)
		
		if err != nil {
			ctx.JSON(iris.Map{
				"status": 	"2",
				"msg":		err,  				
			})	
			return
		} else {
			ctx.JSON(iris.Map{
				"status": 	"0",
				//"msg":		"扫描的是资产二维码",
				//"data":    	asset,
                "recorders":recorders,
			})			
			return
		}	
		
	})//app.Get("/queryrecorders/{jh:string}", func(ctx iris.Context){

	//添加用户
	app.Get("/adduser/{jh:string}/{xm:string}", func(ctx iris.Context) {
		//操作确权, 如果不是admin,则引导到主页
		getcookienum := ctx.GetCookie("cookienum")
		cookiestr := fmt.Sprintf("%v",cookienum)
		if getcookienum != cookiestr {
			ctx.Redirect("/")
			return
		}
		
		jh := ctx.Params().Get("jh")
		xm := ctx.Params().Get("xm")

		user := &User{
			JH: jh,
			XM: xm,
			MM: jh}
		_, err := orm.Insert(user) //<--------插入
		if err != nil {
			ctx.Writef("<H1></H1>")
			ctx.Writef("<H3>%#v</H3>", err)
		} else {
			ctx.Writef("<H1></H1>") //没有这一行,下面的<H3>格式显示错误
			ctx.Writef("<h3>添加用户成功, 请使用警号登录,初始密码即警号!</h3>")
		}
	})

	//删除某个用户
	app.Get("/deluser/{jh:string}", func(ctx iris.Context) {
		//操作确权, 如果不是admin,则引导到主页
		getcookienum := ctx.GetCookie("cookienum")
		cookiestr := fmt.Sprintf("%v",cookienum)
		if getcookienum != cookiestr {
			ctx.Redirect("/")
			return
		}
		
		jh := ctx.Params().Get("jh")
		//xm := ctx.Params().Get("xm")

		user := &User{
			JH: jh,
			//XM: xm,
			//MM: jh
		}

		n, err := orm.Delete(user) //<--------删除
		if err != nil {
			ctx.Writef("<H1></H1>")
			ctx.Writef("<H3>%#v</H3>", err)
		} else if n == 0 {
			ctx.Writef("<H1></H1>") //没有这一行,下面的<H3>格式显示错误
			ctx.Writef("<h3>没有找到!</h3>")
		} else {
			ctx.Writef("<H1></H1>") //没有这一行,下面的<H3>格式显示错误
			ctx.Writef("<h3>删除成功!</h3>")
		}
	})

	//删除全部用户
	app.Get("/delallusers", func(ctx iris.Context) {
		//操作确权, 如果不是admin,则引导到主页
		getcookienum := ctx.GetCookie("cookienum")
		cookiestr := fmt.Sprintf("%v",cookienum)
		if getcookienum != cookiestr {
			ctx.Redirect("/")
			return
		}
		
		ctx.ServeFile("./web/delallusers.html", false) //不压缩

	})

	app.Post("/delallusers", func(ctx iris.Context) {
		//操作确权, 如果不是admin,则引导到主页
		getcookienum := ctx.GetCookie("cookienum")
		cookiestr := fmt.Sprintf("%v",cookienum)
		if getcookienum != cookiestr {
			ctx.Redirect("/")
			return
		}
		
		msg := ctx.FormValue("msg")
		fmt.Printf("接收到msg: %v\n", msg)

		err := orm.DropTables(new(User)) //删除数据及表结构
		fmt.Println(err)
		err = orm.Sync2(new(User)) //重建表结构
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

	//重定向到importassets.html网页, 以便提供一致的用户体验
	app.Get("/importassets", func(ctx iris.Context) {
		//操作确权, 如果不是admin,则引导到主页
		getcookienum := ctx.GetCookie("cookienum")
		cookiestr := fmt.Sprintf("%v",cookienum)
		if getcookienum != cookiestr {
			ctx.Redirect("/")
			return
		}
		
		ctx.ServeFile("./web/importassets.html", false) // true for gzip.
		//        ctx.ServeFile("./html/index.html", true) // true for gzip.

	})

	//接收JSON格式的资产信息,导入数据库
	app.Post("/importassets", func(ctx iris.Context) {
		//操作确权, 如果不是admin,则引导到主页
		getcookienum := ctx.GetCookie("cookienum")
		cookiestr := fmt.Sprintf("%v",cookienum)
		if getcookienum != cookiestr {
			ctx.Redirect("/")
			return
		}
		
		msg := ctx.FormValue("msg")
		fmt.Printf("接收到的用户EXCEL表: %v\n", msg)

		//str := `{"page": 1, "fruits": ["apple", "peach"]}`
		assets := &[]Asset{}
		if err := json.Unmarshal([]byte(msg), &assets); err != nil {

			msg := fmt.Sprintf("%v", err)
			msg = "解析错误:" + msg
			fmt.Printf("%v\n", msg)
			ctx.JSON(iris.Map{
				"msg": msg, //这也是在$ajax.success中的
			})
			//ctx.JSON(msg)	//这样也可以, 只要接收端对应即可

			return
		}
		//fmt.Println(users)
		fmt.Printf("导入资产EXCEL表解析结果: %v\n", *assets)

		for i, v := range *assets {
			//fmt.Println(v)
			asset := &Asset{
				XH:     v.XH,
				BM:     v.BM,
				ZCBH:   v.ZCBH,
				ZCMC:   v.ZCMC,
				PP:     v.PP,
				GGXH:   v.GGXH,
				QDRQ:   v.QDRQ,
				SL:     v.SL,
				DJ:     v.DJ,
				JE:     v.JE,
				SYR:    v.SYR,
				CFDD:   v.CFDD,
				BZ:     v.BZ,
				QRCODE: v.QRCODE,
			}
			_, err := orm.Insert(asset) //<--------插入
			if err != nil {
				msg := fmt.Sprintf("导入失败: %v,  第 %v 行", err, i)
				fmt.Println(msg)
				ctx.JSON(iris.Map{
					"msg": msg,
				})
				return
			}
		}

		ctx.JSON(iris.Map{
			"msg": "导入成功",
		})

		fmt.Printf("ok\n")

	})

	//删除全部资产表
	app.Get("/delallassets", func(ctx iris.Context) {
		//操作确权, 如果不是admin,则引导到主页
		getcookienum := ctx.GetCookie("cookienum")
		cookiestr := fmt.Sprintf("%v",cookienum)
		if getcookienum != cookiestr {
			ctx.Redirect("/")
			return
		}
		
		ctx.ServeFile("./web/delallassets.html", false) //不压缩

	})

	app.Post("/delallassets", func(ctx iris.Context) {
		//操作确权, 如果不是admin,则引导到主页
		getcookienum := ctx.GetCookie("cookienum")
		cookiestr := fmt.Sprintf("%v",cookienum)
		if getcookienum != cookiestr {
			ctx.Redirect("/")
			return
		}
		
		msg := ctx.FormValue("msg")
		fmt.Printf("接收到msg: %v\n", msg)

		err := orm.DropTables(new(Asset)) //删除数据及表结构
		fmt.Println(err)
		err = orm.Sync2(new(Asset)) //重建表结构
		if err != nil {
			app.Logger().Fatalf("orm failed to initialized User table: %v", err)
		}
		ctx.JSON("全部删除成功")

	})

	//重定向到exportassets.html网页, 以便提供一致的用户体验
	app.Get("/exportassets", func(ctx iris.Context) {
		//操作确权, 如果不是admin,则引导到主页
		getcookienum := ctx.GetCookie("cookienum")
		cookiestr := fmt.Sprintf("%v",cookienum)
		if getcookienum != cookiestr {
			ctx.Redirect("/")
			return
		}
		
		ctx.ServeFile("./web/exportassets.html", true) // true for gzip.
		//        ctx.ServeFile("./html/index.html", true) // true for gzip.

	})
	
	//重定向到exportrecorders.html网页, 以便提供一致的用户体验
	app.Get("/exportrecorders", func(ctx iris.Context) {
		//操作确权, 如果不是admin,则引导到主页
		getcookienum := ctx.GetCookie("cookienum")
		cookiestr := fmt.Sprintf("%v",cookienum)
		if getcookienum != cookiestr {
			ctx.Redirect("/")
			return
		}
		
		ctx.ServeFile("./web/exportrecorders.html", true) // true for gzip.
		//        ctx.ServeFile("./html/index.html", true) // true for gzip.

	})
	//发送JSON格式的流转记录表,以便EXCEL文件导出
	app.Post("/exportrecorders", func(ctx iris.Context) {
		//操作确权, 如果不是admin,则引导到主页
		getcookienum := ctx.GetCookie("cookienum")
		cookiestr := fmt.Sprintf("%v",cookienum)
		if getcookienum != cookiestr {
			ctx.Redirect("/")
			return
		}
		
		//ctx.Gzip(true)               // enable gzip for big files
		recorders := make([]Recorder, 0)
		err := orm.Find(&recorders)
		if err != nil {
			ctx.JSON(err)
		} else {
			//ctx.Writef("%#v", users)
			ctx.JSON(recorders)
		}

		//		ctx.JSON(iris.Map{
		//			"status":  "posted",
		//			"message": "message",
		//			"nick":    "nick",
		//		})
	})
	
	//重定向到admin.html网页, 以便提供一致的用户体验
	app.Get("/admin", func(ctx iris.Context) {
		//操作确权, 如果不是admin,则引导到主页
		getcookienum := ctx.GetCookie("cookienum")
		cookiestr := fmt.Sprintf("%v",cookienum)
		if getcookienum != cookiestr {
			ctx.Redirect("/")
			return
		}
		
		ctx.ServeFile("./web/admin.html", true) // true for gzip.
		//        ctx.ServeFile("./html/index.html", true) // true for gzip.

	})
	//发送JSON格式的用户信息,以便EXCEL文件导出
	app.Post("/getassets", func(ctx iris.Context) {
		//操作确权, 如果不是admin,则引导到主页
		getcookienum := ctx.GetCookie("cookienum")
		cookiestr := fmt.Sprintf("%v",cookienum)
		if getcookienum != cookiestr {
			ctx.Redirect("/")
			return
		}
		
		assets := make([]Asset, 0)
		err := orm.Find(&assets)
		if err != nil {
			ctx.JSON(err)
		} else {
			//ctx.Writef("%#v", users)
			ctx.JSON(assets)
		}

		//		ctx.JSON(iris.Map{
		//			"status":  "posted",
		//			"message": "message",
		//			"nick":    "nick",
		//		})
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

	fmt.Println("http://localhost:8101")
	fmt.Println("http://localhost:8101/ws")
	app.Run(iris.Addr(":8101")) //<--------------------------------------

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


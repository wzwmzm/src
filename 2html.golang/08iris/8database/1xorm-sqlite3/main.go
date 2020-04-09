// Package main shows how an orm can be used within your web app
// it just inserts a column and select the first.
package main

import (
	"time"

	"github.com/kataras/iris"

	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
	//_ "github.com/go-sql-driver/mysql"
)

/*
	go get -u github.com/mattn/go-sqlite3
	go get -u github.com/go-xorm/xorm

	If you're on win64 and you can't install go-sqlite3:
		1. Download: https://sourceforge.net/projects/mingw-w64/files/latest/download
		2. Select "x86_x64" and "posix"
		3. Add C:\Program Files\mingw-w64\x86_64-7.1.0-posix-seh-rt_v5-rev1\mingw64\bin
		to your PATH env variable.

	Docs: http://xorm.io/docs/
*/

// User is our user table structure.
type User struct {
	ID        int64  // auto-increment by-default by xorm
	Version   string `xorm:"varchar(200)"`
	Salt      string
	Username  string
	Password  string    `xorm:"varchar(200)"`
	Languages string    `xorm:"varchar(200)"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
}

func main() {
	app := iris.New()
/*
	x,err:=xorm.NewEngine("mysql", "root:111111@/sys?charset=utf8")
	engine.ShowWarn = true //显示警告
	1,增加：  插入一条新的记录，该记录必须是未存在的，否则会返回错误
			_, err := x.Insert(&Account{Name: name, Balance: balance})
	2,删除:   如果多个字段同时赋值，则是多个条件同时满足的记录都会被删除。
			_, err := x.Delete(&Account{Id: id})
	3,查询 单条数据 使用Get方法:
		A:
			a:=&Account{}
			has, err := x.Id(id).Get(a)		//has=true/false;
			//ａ中的非空域自动成为查询的条件
		B:
			a := new(Account)
			has, err := x.Where("name=?", "adn").Get(a)
			//ａ中的非空域自动成为查询的条件
		Ｃ:
			a := &Account{Id:1}
			has, err := x.Get(a)
	4,修改:
			a.Count = 101			//对上面的查询结果进行修改
	 		_, err = x.Update(a)	//注意，Update 接受的参数是指针
	5,批量获取信息:
			err = x.Desc("balance").Find(&as)
			//在这里，我们还调用了 Desc 方法对记录按照存款数额将账户从大到小排序。
			//Find 方法的第一个参数为 slice 的指针或 Map 指针，即为查询后返回的结果，
			//第二个参数可选，为查询的条件 struct 的指针。
	6,锁及事务:
			sess := x.NewSession()
			defer sess.Close()
	 		//开启事务
			if err = sess.Begin(); err != nil {
				return err
			}
			if _, err = sess.Update(a1); err != nil {
			// 发生错误时进行回滚
				sess.Rollback()
				return err
			}//事务完成
			err = session.Commit()
			if err != nil {
				return
			}
	7,统计:
			total,err = x.Count(a)
	8,Iterate方法:
			err := x.Where("id > ?=)", 30).Iterate(new(Account), func(i int, bean interface{})error{
				user := bean.(*Account)	//断言,完成类型转换
				do somthing use i and user
			})
			提供逐条执行查询到的记录的方法，他所能使用的条件和Find方法完全相同
	9,查询特定字段:
			x.Cols("name").Iterate(new(Account), printFn)
			var printFn = func(idx int, bean interface{}) error {
				dosomething
				return nil
			}
			查询出来的结构只有 Name 字段有值，其它字段均为零值。要注意的是，Cols 方法所接受的参数是数据表中对应的名称，而不是字段名称。
	10,排除特定字段:
			x.Omit("name").Iterate(new(Account), printFn)
			所查询出来的结构只有 Name 字段为零值。要注意的是，Omit 方法所接受的参数是数据表中对应的名称，而不是字段名称。
	11,查询结果偏移:
			x.Limit(3, 2).Iterate(new(Account), printFn)
			此处的查询结果为偏移 2 个后，最多取出 3 个记录。
	12,日志记录:
			x.ShowSQL = true	//开启日志功能
			//下面是日志写入文件
			f, err := os.Create("sql.log")
			if err != nil {
				log.Fatalf("Fail to create log file: %v\n", err)
				return
			}
			x.Logger = xorm.NewSimpleLogger(f)
	13,LRU 缓存: (将数据库缓存,以提高速度)
			//下面两句紧接 xorm.NewEngine() 后
			cacher := xorm.NewLRUCacher(xorm.NewMemoryStore(), 1000)
			x.SetDefaultCacher(cacher)
	14,事件钩子:
		A:
			func (a *Account) BeforeInsert() {
				log.Printf("before insert: %s", a.Name)
			}
		B:
			func (a *Account) AfterInsert() {
				log.Printf("after insert: %s", a.Name)
			}
*/
	orm, err := xorm.NewEngine("sqlite3", "./test.db")
	if err != nil {
		app.Logger().Fatalf("orm failed to initialized: %v", err)
	}

	iris.RegisterOnInterrupt(func() {
		orm.Close()
	})

	err = orm.Sync2(new(User))
	//* 自动检测和创建表，这个检测是根据表的名字
	//* 自动检测和新增表中的字段，这个检测是根据字段名，同时对表中多余的字段给出警告信息
	//* 自动检测，创建和删除索引和唯一索引，这个检测是根据索引的一个或多个字段名，而不根据索引名称。因此这里需要注意，如果在一个有大量数据的表中引入新的索引，数据库可能需要一定的时间来建立索引。
	//* 自动转换varchar字段类型到text字段类型，自动警告其它字段类型在模型和数据库之间不一致的情况。
	//* 自动警告字段的默认值，是否为空信息在模型和数据库之间不匹配的情况
	//
	//以上这些警告信息需要将`engine.ShowWarn` 设置为 `true` 才会显示。

	if err != nil {
		app.Logger().Fatalf("orm failed to initialized User table: %v", err)
	}

	//app.Get("/insert", func(ctx iris.Context) {
	app.Get("/", func(ctx iris.Context) {
		user := &User{Username: "kataras", Salt: "hash---", Password: "hashed", CreatedAt: time.Now(), UpdatedAt: time.Now()}
		orm.Insert(user) //<--------插入

		ctx.Writef("user inserted: %#v", user)
	})

	app.Get("/get", func(ctx iris.Context) {
		user := User{ID: 1}
		if ok, _ := orm.Get(&user); ok { //<--------查询
			ctx.Writef("user found: %#v", user)
		}
	})

	// http://localhost:8080/insert
	// http://localhost:8080/get
	app.Run(iris.Addr(":8101"), iris.WithoutServerError(iris.ErrServerClosed))
}

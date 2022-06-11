/***
原始文章：  golang获取基金数据： https://www.jianshu.com/p/674205726be7

其中使用了“golang调用javascript”技术，文章地址： https://www.jianshu.com/p/2d33545f00ab
使用此技术目的也就省了JSON数据解析，似乎毫无必要，多此一举。

调用的数据接口： 
	1，本程序用的是东方财富网的接口，
		返回数据用了类似var gs = [];的形式，更好地配合JS程序，
		地址： https://fund.eastmoney.com/js/fundcode_search.js
	2，也可以用《小熊同学》的接口： https://www.doctorxiong.club/
***/

package main

import (
    "fmt"
    "github.com/fatih/color"
    "github.com/robertkrimen/otto"
    "github.com/rodaine/table"
    "io/ioutil"
    "net/http"
)

//封装了JS脚本调用,功能是提取vm中变量name的值
//其实只省略了错误检测
func getJsVariable(vm *otto.Otto, name string) interface{} {
    v, err := vm.Get(name)
    if err != nil {
        fmt.Println(err)
        return nil
    }

    obj, _ := v.Export()
    return obj
}

/**
基金列表
*/
func fundList() {
    url := "https://fund.eastmoney.com/js/fundcode_search.js"
    method := "GET"

    client := &http.Client{}
    req, err := http.NewRequest(method, url, nil)	//1，封装网络请求

    if err != nil {
        fmt.Println(err)
        return
    }
    res, err := client.Do(req)				//2，执行网络请求，得到回复数据
    if err != nil {
        fmt.Println(err)
        return
    }
    defer res.Body.Close()

    body, err := ioutil.ReadAll(res.Body)		//3，取得回复数据
    if err != nil {
        fmt.Println(err)
        return
    }



    vm := otto.New()
    //js格式，只有一个变量r
    //body中数据是可执行的JS语句 var r = [["000001","HXCZHH","华夏成长混合","混合型-偏股","HUAXIACHENGZHANGHUNHE"]]
    vm.Run(string(body))

    dataList := getJsVariable(vm, "r")
    fundList := dataList.([][]string)

    headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
    columnFmt := color.New(color.FgYellow, color.CrossedOut).SprintfFunc()

    tbl := table.New("序号", "基金代码", "拼音缩写", "基金名称", "基金类型", "拼音全称")
    tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

    for i, fundItem := range fundList {
        if i < 100 {
            tbl.AddRow(i+1, fundItem[0], fundItem[1], fundItem[2], fundItem[3], fundItem[4])
        }
    }

    tbl.Print()
}

/**
基金公司列表
*/
func fundCompanyList() {
    url := "https://fund.eastmoney.com/js/jjjz_gs.js"
    //返回值如：var gs={op:[["80163340","安信基金"],["81052915","安信证券资产"]]}

    method := "GET"

    client := &http.Client{}
    req, err := http.NewRequest(method, url, nil)	//1，封装网络请求

    if err != nil {
        fmt.Println(err)
        return
    }
    res, err := client.Do(req)				//2，执行网络请求，取得返回值
    if err != nil {
        fmt.Println(err)
        return
    }
    defer res.Body.Close()

    body, err := ioutil.ReadAll(res.Body)		//3，取得数据主体
    if err != nil {
        fmt.Println(err)
        return
    }

    vm := otto.New()
    vm.Run(string(body))				//4,因为body中内容是var gs={...},所以是可执行语句
    gs := getJsVariable(vm, "gs")			//5,从vm中取得变量gs的值

    //以下是公司的格式化显示
    gsmap := gs.(map[string]interface{})
    op := gsmap["op"].([][]string)

    headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
    columnFmt := color.New(color.FgYellow, color.CrossedOut).SprintfFunc()

    tbl := table.New("序号", "公司编码", "公司名称")
    tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

    for i, companyItem := range op {
        tbl.AddRow(i+1, companyItem[0], companyItem[1])
    }

    tbl.Print()
}

/**
基金净值
*/
func fundWorth(code string) {
    url := fmt.Sprintf("https://fundgz.1234567.com.cn/js/%s.js", code)
    method := "GET"

    client := &http.Client{}
    req, err := http.NewRequest(method, url, nil)

    if err != nil {
        fmt.Println(err)
        return
    }
    res, err := client.Do(req)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer res.Body.Close()

    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        fmt.Println(err)
        return
    }

    vm := otto.New()
    //jsonpgz({"dwjz":"1.2350","fundcode":"000001","gsz":"1.2333","gszzl":"-0.14","gztime":"2021-11-05 15:00","jzrq":"2021-11-04","name":"华夏成长混合"});
    //这里设置了一个JS函数jsonpgz, 功能就是返回调用函数的参数
    vm.Set("jsonpgz", func(call otto.FunctionCall) otto.Value {
        return call.Argument(0)
    })

    result, _ := vm.Run(string(body))	//这里实际上执行了函数调用：jsonpgz({"dwjz":"1.2350"...});返回的是调用的参数
    resultExp, _ := result.Export()
    fundSummary := resultExp.(map[string]interface{})
    //map[dwjz:1.2350 fundcode:000001 gsz:1.2333 gszzl:-0.14 gztime:2021-11-05 15:00 jzrq:2021-11-04 name:华夏成长混合]
    fmt.Printf("名称：%v\n", fundSummary["name"])
    fmt.Printf("单位净值：%v\n", fundSummary["dwjz"])
    fmt.Printf("估算值：%v\n", fundSummary["gsz"])
}

/**
基金详情
返回值样例在文件 00基金详情.js 中，非常的详细
*/
func fundDetail(code string) {
    url := fmt.Sprintf("https://fund.eastmoney.com/pingzhongdata/%s.js", code)
    //返回值样例在文件 00基金详情.js 中，非常的详细
    method := "GET"

    client := &http.Client{}
    req, err := http.NewRequest(method, url, nil)

    if err != nil {
        fmt.Println(err)
        return
    }
    res, err := client.Do(req)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer res.Body.Close()

    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        fmt.Println(err)
        return
    }

    vm := otto.New()

    vm.Run(string(body))
    //收益率
    syl_1n := getJsVariable(vm, "syl_1n")
    syl_6y := getJsVariable(vm, "syl_6y")
    syl_3y := getJsVariable(vm, "syl_3y")
    syl_1y := getJsVariable(vm, "syl_1y")
    fmt.Printf("近一年收益率: %v\n", syl_1n)
    fmt.Printf("近6月收益率: %v\n", syl_6y)
    fmt.Printf("近三月收益率: %v\n", syl_3y)
    fmt.Printf("近一月收益率: %v\n", syl_1y)
}

func main() {
    fundCompanyList()
    fundList()
    fundWorth("000001")
    fundDetail("000001")
}


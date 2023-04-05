
package main

import (
    "fmt"
    "github.com/fatih/color"
    "github.com/robertkrimen/otto"
    "github.com/rodaine/table"
    "io/ioutil"
    "net/http"
)

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
    //js格式，只有一个变量r
    //var r = [["000001","HXCZHH","华夏成长混合","混合型-偏股","HUAXIACHENGZHANGHUNHE"]]
    vm.Run(string(body))

    dataList := getJsVariable(vm, "r")
    fundList := dataList.([][]string)

    headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
    columnFmt := color.New(color.FgYellow, color.CrossedOut).SprintfFunc()

    tbl := table.New("序号", "基金代码", "拼音缩写", "基金名称", "基金类型", "拼音全称")
    tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

    for i, fundItem := range fundList {
//        if i < 100 {		//取消100的限制,显示所有的基金
            tbl.AddRow(i+1, fundItem[0], fundItem[1], fundItem[2], fundItem[3], fundItem[4])
//        }
    }

    tbl.Print()
}

/**
基金公司列表
*/
func fundCompanyList() {
    url := "https://fund.eastmoney.com/js/jjjz_gs.js"
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
    gs := getJsVariable(vm, "gs")

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
    vm.Set("jsonpgz", func(call otto.FunctionCall) otto.Value {
        return call.Argument(0)
    })

    result, _ := vm.Run(string(body))
    resultExp, _ := result.Export()
    fundSummary := resultExp.(map[string]interface{})
    //map[dwjz:1.2350 fundcode:000001 gsz:1.2333 gszzl:-0.14 gztime:2021-11-05 15:00 jzrq:2021-11-04 name:华夏成长混合]
    fmt.Printf("名称：%v\n", fundSummary["name"])
    fmt.Printf("单位净值：%v\n", fundSummary["dwjz"])
    fmt.Printf("估算值：%v\n", fundSummary["gsz"])
}

/**
基金详情
*/
func fundDetail(code string) {
    url := fmt.Sprintf("https://fund.eastmoney.com/pingzhongdata/%s.js", code)
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
    fmt.Printf("基金公司列表\n\n")
    fundCompanyList()
    fmt.Printf("\n\n\n基金列表\n\n")
    fundList()
    fmt.Printf("\n\n\n基金净值\n\n")
    fundWorth("000001")
    fmt.Printf("\n\n\n基金详细\n\n")
    fundDetail("000001")
}


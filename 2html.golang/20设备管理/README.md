通用查询并展示的设想:
1, xorm的sql查询: 
    sql := "select * from userinfo"
    results, err := engine.Query(sql)
    返回结果是: 第一个返回值results为[]map[string][]byte的形式。
2, 如果直接 c.JSON(results),  那么 []byte 转 json 的时候将用base64转成字符串
3, 所以要先将 []byte 转成string, 互转示例:
    str2 := "hello"
    data2 := []byte(str2)
    fmt.Println(data2)
    str2 = string(data2[:])
    fmt.Println(str2)
4, 然后再 c.JSON(results)
5, H5 表格展示 JSON数组:
    <!DOCTYPE html>
    <html>
        <head>
            <meta charset="UTF-8">
            <title></title>
            <script src="JQuery/jquery-1.11.3.js"></script>
            <script type="text/javascript">
                var stu = [
                    { "name": "张三", "sex": "男", "age": "20" },
                    { "name": "李四", "sex": "男", "age": "18" },
                    { "name": "王五", "sex": "男", "age": "19" }
                ];
                $(document).ready(function() {
                    var s = "";
                    for(var i = 0; i < stu.length; i++) {
                        s="<tr>";//重新装载每一行
                        for(var key in stu[i]){	 //此时key=属性名               	
                            s=s+"<td>"+stu[i][key]+"</td>";	                	
                        }
                      s+="</tr>";  
                      $("#tab").append(s);//把遍历到的每一行 加入
                   }
                });
            </script>
        </head>
        <body>
            <table  border="1px" id="tab">//没有border就没有框线
                <tr>
                    <td>姓名</td>
                    <td>性别</td>
                    <td>年龄</td>
                </tr>
            </table>
        </body>
    </html>


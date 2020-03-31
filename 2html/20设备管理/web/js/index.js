let n_cameras = 0; //一共有几个摄像头
let cameraslist = null; //取得的摄像头组, 对于scanner是通用的
let usr; //存放当前用户对象JSON.parse(localStorage.usr)
let asset; //存放当前二维码扫描的资产
let assets; //存放当前地址二维码的资产列表
let myscanner; //通用扫描仪, 多页共用这一个扫描仪
let recorders; //存放设备的流转记录
if (localStorage.status == undefined) localStorage.status = "0"; //localStorage.status设置的是允许进入的最大页号
//创建COOKIE. cookie如果由服务器创建,缺省为HTTPONLY,客户端JS无法操作
if (getCookie("cookienum") == null) {
	document.cookie = "cookienum=0000000000"
}


//1,admin 权限管理(完成)
//2,流转记录EXCEL文件导出(完成)
//3,流转记录按警号查询输出(完成)
//4,配色整理
//5,PWA调整
//6,事务管理调整(完成)

//请登录提示
if (localStorage.usr == undefined) {
	$("#name_xs").html('<b style="color:MediumVioletRed ">请先登录</b>');
} else {
	usr = JSON.parse(localStorage.usr)
	$("#name_xs").text(usr["姓名"]);
}

//本人流转记录查询
$("#queryrecorders").click(function () {
	if (localStorage.usr == undefined) {
		alert("请先登录!");
	};

	$.get("queryrecorders/" + usr.警号,
		function (data, status) {
			switch (data.status) {
				case "0":
					let recorders1 = data.recorders;
					console.log("流转记录: ", recorders1);
					$(".recorderslist1").remove();
					let t1 = `<p class="recorderslist1"><br/><b>登记记录如下:</b>`
					//for (let key in recorders1) {
					for (let key = recorders1.length;
						(key - 1) >= 0; key--) {
						let val = recorders1[key - 1];
						console.log(val.记录时间);
						t1 += `<br/><b class="text-info">时间: ` + val.记录时间 + `</b>`;
						t1 += `<br/>资产名称: ` + val.资产名称;
						t1 += `<br/>资产二维码: ` + val.二维码;
						t1 += `<br/>变更地址: ` + val.存放地点;
						t1 += `<br/>备注:` + val.备注; // + `<br/>`;
					} //for (let key in recorders) {
					$("#queryrecorders").after(t1);
					break;
				case "2": //读取流转记录时发生错误
					alert("错误: ", data.msg);
					break;
			} //switch

		}) //$.get("queryrecorders/XXX",

}) //$("#queryrecorders").click(function(){

//修改密码
$("#changepwdbtm").click(function () {
	if (localStorage.usr == undefined) {
		alert("请先登录!");
		return;
	} else {
		console.log(JSON.parse(localStorage.usr).ID);
	}
	var word1 = document.getElementById("pw1").value;
	var word2 = document.getElementById("pw2").value;
	if (word1 == "" || word2 == "") {
		return;
	}
	//校验两次输入一致
	if (word1 != word2) {
		window.alert("两次输入的新密码不一致！");
		return;
	}

	$.post("/changepwd", {
			jh: JSON.parse(localStorage.usr)["警号"],
			pwd: word2
		},
		function (data, status) {
			alert(data.msg);
		});
})

//login
$("#loginbtm").click(function () {
	let jh = $("#jh").val();
	let pwd = $("#pwd").val();
	if (jh == "" || pwd == "") {
		return;
	}

	//alert("login begin")
	$.post("/login", {
			jh: jh,
			pwd: pwd
		},
		function (data, status) {
			if (data.status != "0") {
				alert(data.msg);
			} else {
				//alert(JSON.stringify(data.msg));
				//alert(data.msg["姓名"]);
				$("#name_xs").text(data.data["姓名"]);
				usr = data.data;
				localStorage.usr = JSON.stringify(data.data);
				document.cookie = "cookienum=" + data.cookienum;
				allow_1();
			}
		});

});

//changeaddr
$("#submitbtn").click(function () {
	//	let id = asset.ID;
	//	let addr_old = asset.存放地点;
	//	let addr_new = localStorage.txt;
	//	let comment = $("#comment").val();
	//	console.log(id," : ",addr_old," : ",addr_new," : ",comment);

	$.post("/changeaddr", {
			asset: localStorage.asset,
			addr: localStorage.newaddr,
			user: localStorage.usr,
			comment: $("#comment").val(),
		},
		function (data, status) {
			if (data.status != "1") { //失败
				alert(data.msg);
			} else { //成功
				alert(data.msg);
				localStorage.status = 1;
				swiper.slideTo(1);
				$("#comment").val("");
			}
		});

});

var swiper = new Swiper('.swiper-container', {
	on: {
		slideNextTransitionStart: function () {
			if (swiper.activeIndex == 2 && localStorage.status == "5") swiper.slideTo(5);
		},

		slidePrevTransitionStart: function () {
			if (swiper.activeIndex == 4 && localStorage.status == "5") swiper.slideTo(1);
		},

		slideChange: function () {
			//console.log("event: slideChange");
			//console.log('改变了，activeIndex为' + this.activeIndex);

			switch (this.activeIndex) {
				case 0:
					x = "0 页 : 设置";
					myscanner.stop().then(function () {
						console.log("myscanner0 stop.");
						try {
							scan1();
						} catch (err) {
							console.log("scan1()---错误:" + err.name + ":" + err.message);
						};
					});
					switchstatus();
					break;
				case 1:
					x = "1 页 : 第一次扫码";
					myscanner.stop().then(function () {
						console.log("myscanner1 stop.");
						try {
							scan2();
						} catch (err) {
							console.log("scan2()---错误:" + err.name + ":" + err.message);
						}
					});
					switchstatus();
					break;
				case 2:
					x = "2 页 : 显示扫描设备信息 ";
					myscanner.stop();
					switchstatus();

					asset = JSON.parse(localStorage.asset);
					$("#p0").text(asset.资产名称);
					$("#p1").text("存放地点: " + asset.存放地点);
					$("#p2").text("资产编号: " + asset.资产编号);
					$("#p3").text("规格型号: " + asset.规格型号);
					$("#p4").text("取得日期: " + formatDate(asset.取得日期, '-'));
					$("#p5").text("使 用 人: " + asset.使用人);
					$("#p6").text("资产二维码: " + asset.二维码);
					$("#p7").text("数   量: " + asset.数量);
					$("#p8").text("单   价: " + asset.单价);
					$("#p9").text("金   额: " + asset.金额);
					$("#p10").text("序   号: " + asset.序号);
					$("#p11").text("部   门: " + asset.部门);
					$("#p12").text("品   牌: " + asset.品牌);
					$("#p13").text("备   注: " + asset.备注);

					recorders = JSON.parse(localStorage.recorders);
					$(".recorderslist").remove();
					//$("#s0").text(" " + recorders[0].存放地点);
					for (let key in recorders) {
						let val = recorders[key];
						let t1 = `
                            <div class="card  recorderslist">
                                <div class="card-header">
                                    <h6 class="text-info">时间: ` + val.记录时间 + `</h6>
                                </div>
                                <div class="card-body">
                                    <p>人员: ` + val.姓名 + `</p>
                                    <p>地址: ` + val.存放地点 + `</p>
                                    <p>备注: ` + val.备注 + `</p>
                                </div>

                            </div>`;
						$("#liuzhuan").after(t1);
					};
					scrollTo(0, 0);


					break;
				case 3:
					x = "3 页 : 扫描设备更换地址码";
					myscanner.stop().then(function () {
						console.log("myscanner2 stop.");
						try {
							scan3();
						} catch (err) {
							console.log("scan3()---错误:" + err.name + ":" + err.message);
						}
					});
					switchstatus();
					break;
				case 4:
					x = "4 页 : 设备地址变更提交确认";
					myscanner.stop();
					switchstatus();

					asset = JSON.parse(localStorage.asset);
					$("#c0").text(asset.资产名称);
					$("#c1").text(asset.存放地点);
					$("#c2").text(localStorage.newaddr);
					break;
				case 5:
					x = "5 页 : 显示地址上设备信息";
					myscanner.stop();
					switchstatus();
					assets = JSON.parse(localStorage.assets);

					$(".assetslist").remove();
					$("#d0").text(" " + assets[0].存放地点);
					for (let key in assets) {
						let val = assets[key];
						let t1 = '<div class="card    assetslist"> <div class = "card-header text-info" ><h6>资产名称:' + val.资产名称 + '</h6> </div> <div class = "card-body" > ';
						let t2 = '<p>资产编号:' + val.资产编号 + '</p>';
						let t3 = '<p>规格型号:' + val.规格型号 + '</p>';
						let t4 = '<p>取得日期:' + formatDate(val.取得日期, '-') + '</p>';
						let t5 = '<p>使 用 人:' + val.使用人 + '</p>';
						let t6 = '<p>资产二维码:' + val.二维码 + '</p>';
						let t7 = '<p>数   量:' + val.数量 + '</p>';
						let t8 = '<p>单   价:' + val.单价 + '</p>';
						let t9 = '<p>金   额:' + val.金额 + '</p>';
						let t10 = '<p>序   号:' + val.序号 + '</p>';
						let t11 = '<p>部   门:' + val.部门 + '</p>';
						let t12 = '<p>品   牌:' + val.品牌 + '</p>';
						let t13 = '<p>备   注:' + val.备注 + '</p>';
						let ttt = t1 + t2 + t3 + t4 + t5 + t6 + t7 + t8 + t9 + t10 + t11 + t12 + t13 + '</div></div>';
						$("#d0").after(ttt);
					};



					break;

			}
			console.log(x);
		},


	},
});

//根据当前所在页和状态,调整页面的翻页开关
function switchstatus() {
	if ((swiper.activeIndex).toString() == localStorage.status) {
		swiper.allowSlideNext = false;
	} else {
		swiper.allowSlideNext = true;
	}
	//console.log("swiper.allowSlideNext = ", swiper.allowSlideNext);
}

//判断是否进入第1页.  用户名及摄像头设定前禁止滑动
function allow_1() {
	if (localStorage.status === "0" && localStorage.camera_id && localStorage.usr) {
		localStorage.status = "1";
		switchstatus();
	};
}

//////////  1,初始化第一个摄像头(设置页)
function scan1() {
	//	while (cameralock) { console.log(" cameralock, please wait")};
	//	cameralock = true;
	try {
		myscanner = new Instascan.Scanner({
			video: document.getElementById('webcamera1'),
			continuous: true, //持续扫描
			mirror: false, //镜像,使用前置摄像头时使用
			captureImage: false, //将图像数据作为结果的一部分
			backgroundScan: false, //选项卡未激活时是否扫描
			refractoryPeriod: 5000, //连续识别相同QR码之前的时间
			scanPeriod: 1 //两次扫描之间的周期
		});
	} catch (err) {
		console.log("scan3()中, new Instascan.Scanner 出错");
		alert("大哥! 你太快了, 我受不了了... 我们重新来!!!");
		location.reload();
	};

	myscanner.addListener('scan', function (data) {
		console.log("scanner1 is scanning.");
	});
	myscanner.addListener('active', function () {
		console.log("scanner1 is active.");
	});

	myscanner.addListener('inactive', function () {
		console.log("scanner1 is inactive.");
	});


}
//第一次扫描的摄像头
function scan2() {
	//	while (cameralock) { console.log(" cameralock, please wait")};
	//	cameralock = true;
	try {
		myscanner = new Instascan.Scanner({
			video: document.getElementById('webcamera2'),
			continuous: true, //持续扫描
			mirror: false, //镜像,使用前置摄像头时使用
			captureImage: false, //将图像数据作为结果的一部分
			backgroundScan: false, //选项卡未激活时是否扫描
			refractoryPeriod: 5000, //连续识别相同QR码之前的时间
			scanPeriod: 1 //两次扫描之间的周期
		});
	} catch (err) {
		console.log("scan2()中, new Instascan.Scanner 出错");
		alert("大哥! 你太快了, 我受不了了... 我们重新来!!!");
		location.reload();
	};

	myscanner.addListener('scan', function (data) {
		let content = utf82str(data); //解决中文乱码
		//alert("scanner3: " + content);
		console.log("scan2: " + content);
		$.post("/query", {
				dat: content
			},
			function (data, status) {
				switch (data.status) {
					case "0": //扫描的是资产二维码	
						//alert(JSON.stringify(data.msg));
						//alert(data.msg["姓名"]);
						//$("#name_xs").text(data.data["姓名"]);
						asset = data.data;
						localStorage.asset = JSON.stringify(data.data);
						recorders = data.recorders;
						localStorage.recorders = JSON.stringify(data.recorders);

						//$("#asset").text(localStorage.asset);

						localStorage.status = 3; //可以进入第3页扫描设备更换地址码"
						switchstatus();
						swiper.slideNext();
						break;
					case "1": //扫描的是地点二维码
						assets = data.data;
						localStorage.assets = JSON.stringify(data.data);

						//$("#assets").text(localStorage.assets);

						localStorage.status = 5; //可以进入第5页显示地址上设备信息
						switchstatus();
						swiper.slideTo(5);
						break;
					case "2": //错误的二维码
						alert(data.msg + "\n扫描内容: " + content);
						localStorage.status = 1; //退回到第一次扫描
						switchstatus();
						break;
				}
			}); //$.post("/query",
	}); //myscanner.addListener

	myscanner.addListener('active', function () {
		console.log("scanner2 is active.");
	});

	myscanner.addListener('inactive', function () {
		console.log("scanner2 is inactive.");
	});

	myscanner.start(cameraslist[localStorage.camera_id])
		.then(function () {
			console.log("现在是scan2();");
		})
		.catch(function (err) {
			console.log("scan2()错误:" + err.name + ":" + err.message);
			location.reload();
		});
}

//第二次扫描的摄像头
function scan3() {
	//	while (cameralock) { console.log(" cameralock, please wait")};
	//	cameralock = true;
	try {
		myscanner = new Instascan.Scanner({
			video: document.getElementById('webcamera3'),
			continuous: true, //持续扫描
			mirror: false, //镜像,使用前置摄像头时使用
			captureImage: false, //将图像数据作为结果的一部分
			backgroundScan: false, //选项卡未激活时是否扫描
			refractoryPeriod: 5000, //连续识别相同QR码之前的时间
			scanPeriod: 1 //两次扫描之间的周期
		});
	} catch (err) {
		console.log("scan3()中, new Instascan.Scanner 出错");
		alert("大哥! 你太快了, 我受不了了... 我们重新来!!!");
		location.reload();
	};

	//console.log ("myscanner3:"+myscanner);

	myscanner.addListener('scan', function (data) {
		let content = utf82str(data); //解决中文乱码
		//alert("scanner3: " + content);
		console.log("scan3: " + content);
		//	scanner3.stop();	//结束扫描

		$.post("/verifyaddr", {
				dat: content
			},
			function (data, status) {
				switch (data.status) {
					case "1": //完成地址校验,进入确认地址变更提交页,第4页
						localStorage.newaddr = content
						localStorage.status = 4;
						switchstatus();
						swiper.slideNext(); //进入确认提交页,第4页		
						break;
					case "2": //错误的二维码
						alert(data.msg + "\n扫描内容: " + content);
						localStorage.status = 3; //退回到第二次扫描
						switchstatus();
						break;
				}
			}); //$.post("/query",		



	});

	myscanner.addListener('active', function () {
		console.log("scanner3 is active.");
	});

	myscanner.addListener('inactive', function () {
		console.log("scanner3 is inactive.");
	});

	myscanner.start(cameraslist[localStorage.camera_id])
		.then(function () {
			console.log("现在是scan3();");
		})
		.catch(function (err) {
			console.log("scan3()错误:" + err.name + ":" + err.message);
			location.reload();
		});
}


//设置页, 摄像头列表
Instascan.Camera.getCameras().then(function (cameras) {
	switchstatus(); //摄像头初始化完成,可以开始了
	if (cameras.length > 0) {
		cameraslist = cameras; //取得摄像头列表
		// 显示 "摄像头: BisonCam, NB Pro"
		if (localStorage.camera_id) {
			$('#cameraid_c').text(cameraslist[localStorage.camera_id].name);
		}
		//列表开始
		n_cameras = cameras.length;
		for (var i = n_cameras; --i >= 0;) {
			let checked = ""
			if (localStorage.camera_id == i) {
				checked = "checked"
			}
			//onclick实现:
			//	1,设定localStorage.camera_id
			//	2,关闭并重新打开指定摄像头
			//	3,判断是否能进入第1页
			$("#cameras-div").prepend(`
			<div class="radio text-left">
				<label><input type="radio" name="camera_op_c"  ` + checked + `  onclick="localStorage.camera_id=` + i + `;$('#cameraid_c').text(cameraslist[` + i + `].name);myscanner.stop();myscanner.start(cameraslist[` + i + `]);  allow_1(); "> ` + cameraslist[i].name + `</label>
			</div>`);
		};
	} else {
		alert('No cameras found.');
	}
}).catch(function (e) {
	alert(e);
});

// UCS-2和UTF8都是unicode的一种编码方式
// js代码中使用的是UCS-2编码
//这里将扫描得到的UTF8码转成UCS-2码
function utf82str(utf) {
	var str = "";
	var tmp;

	for (var i = 0; i < utf.length; i++) {
		// 英文字符集合
		if (utf.charCodeAt(i) >> 7 == 0x00) {
			str += utf.charAt(i);
			continue;
		}
		// 其他字符集
		else if (utf.charCodeAt(i) >> 5 == 0x06) {
			tmp = ((utf.charCodeAt(i + 0) & 0x1f) << 6) |
				((utf.charCodeAt(i + 1) & 0x3f) << 0);
			str += String.fromCharCode(tmp);
			i++;
			continue;
		}
		// 中文字符集
		else if (utf.charCodeAt(i) >> 4 == 0x0e) {
			tmp = ((utf.charCodeAt(i + 0) & 0x0f) << 12) |
				((utf.charCodeAt(i + 1) & 0x3f) << 6) |
				((utf.charCodeAt(i + 2) & 0x3f) << 0);
			str += String.fromCharCode(tmp);
			i += 2;
			continue;
		}
		// 其他字符集
		else if (utf.charCodeAt(i) >> 3 == 0x1f) {
			tmp = ((utf.charCodeAt(i + 0) & 0x07) << 18) |
				((utf.charCodeAt(i + 1) & 0x3f) << 12) |
				((utf.charCodeAt(i + 2) & 0x3f) << 6);
			((utf.charCodeAt(i + 3) & 0x3f) << 0);
			str += String.fromCharCode(tmp);
			i += 3;
			continue;
		}
		// 非法字符集
		else {
			throw "不是UTF-8字符集"
		}
	}

	return str;
}

//整数转日期, 如formatDate(42333,'-')
function formatDate(numb, format) {
	const time = new Date((numb - 1) * 24 * 3600000 + 1)
	time.setYear(time.getFullYear() - 70)
	const year = time.getFullYear() + ''
	const month = time.getMonth() + 1 + ''
	const date = time.getDate() - 1 + ''
	if (format && format.length === 1) {
		return year + format + month + format + date
	}
	return year + (month < 10 ? '0' + month : month) + (date < 10 ? '0' + date : date)
}

//读取cookies
function getCookie(name) {
	var arr, reg = new RegExp("(^| )" + name + "=([^;]*)(;|$)");

	if (arr = document.cookie.match(reg))

		return unescape(arr[2]);
	else
		return null;
}

//switchstatus(); //进入指定状态
swiper.allowSlideNext = false; //等待摄像头初始化完成
scan1();

///////备用代码//////////////////////////////////
/*
window.onpagehide = function () {
	console.log("pagehide");
}
window.onpageshow = function () {
	console.log("pageshow");
}
window.addEventListener('blur', function () {
	console.log('window blur-2');
});
window.addEventListener('focus', function () {
	console.log('window focus-2');
});
$("#usr").focusin(function () {
	console.log("手风琴1 focusin");
});
$("#usr").focusout(function () {
	console.log("手风琴1 focusout");
});





$("#slide1").hide();


var mySwiper = new Swiper('.swiper-container', {

	initialSlide :2,		//1,设置为2后，画面初始停在第三个Slide
	direction : 'vertical',	//2,滑动方向,(horizontal)或垂直(vertical)
	observer:true,			//3,swiper元素更新后的显示自动更新, 
							//  默认false，不开启，可以使用update()手动更新。
	on: {					//4,事件注册, 有哪些具体事件可以查看手册
			slideChange: function () {	//A,(activeIndex发生改变时)
				console.log(this.activeIndex);
			},

			slideChangeTransitionEnd: function(){//B,切换结束时
				alert(this.activeIndex);
			},
	},
})

//Swiper的一些主要属性------------------------------
1,mySwiper.activeIndex		//返回当前活动块(激活块)的索引。
2,mySwiper.el.style.opacity=0.5;//el, 相当于 ${} 选择器
3,mySwiper.allowSlideNext	//设置/查看是否允许切换至下一个slide。
4,mySwiper.allowTouchMove	//设置/查看是否禁止触摸滑动。

//Swiper的一些主要方法------------------------------
1,mySwiper.slideNext()		//滑动到下一滑块
2,mySwiper.slideTo(index)	//滑动到指定滑块
3,mySwiper.updateSlides()	//重新计算Slides的数量，当你使用js来删除slide时可能会用到。使用mySwiper.removeSlide来删除slide则会自动计算，无需使用此方法。
4,mySwiper.update()			//更新Swiper，就像重新初始化一样
5,mySwiper.appendSlide(slides)//添加slide到slides的结尾。
6,mySwiper.prependSlide(slides)//添加slide到slides的第一个位置。
7,mySwiper.addSlide(index, slides)//在指定位置增加slide
8,mySwiper.removeSlide(index)//移除索引为index的slide。
9,mySwiper.removeAllSlides()//移除所有slides。
10,mySwiper.on(event,handler)//添加回调函数或者事件句柄。
11,


//1,添加删除页(没有更新显示)-------------------------
//注意: 修改swiper后如果要更新显示效果,需要调Update()方法,或设置observer:true,
swiper.appendSlide([
	'<div class="swiper-slide">Slide ' + (++appendNumber) + '</div>',
	'<div class="swiper-slide">Slide ' + (++appendNumber) + '</div>'
]);

mySwiper.removeAllSlides();      //移除全部
mySwiper.addSlide(index, slides);//在指定位置增加
mySwiper.appendSlide(slides);    //在尾部添加

//2,删除页后更新显示:--------------------------------
var mySwiper = new Swiper('.swiper-container',{
	pagination: {
		el: '.swiper-pagination',//分页器,就是几个点点
	},
	observer:true,	//<-----相当于自动初始化swiper
					//<-----默认false，不开启，可以使用update()方法更新。
})
$('#btn1').click(function(){
	$(".swiper-wrapper .swiper-slide1").remove();
})

//3,删除页后更新显示----------------------------------
var mySwiper = new Swiper('.swiper-container',{
	pagination : '.swiper-pagination',//分页器,就是几个点点
})
$('#btn1').click(function(){
	$(".swiper-wrapper .swiper-slide1").remove(); 
  mySwiper.update()	//<------更新显示(不然会有一个空白页)
})

//4,可以使用 div 的 显示/隐藏 属性实现页面的添加或删除(bootstrap 的辅助类show/hide)

*/

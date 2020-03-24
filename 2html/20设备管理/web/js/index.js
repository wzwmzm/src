//if (!(localStorage.camera_id)) localStorage.camera_id=0;
//选用哪个摄像头,  注意:  不初始化!!!
let n_cameras = 0; //一共有几个摄像头
let cameraslist = null; //取得的摄像头组, 对于scanner是通用的
let usr; //存放当前用户对象JSON.parse(localStorage.usr)
let asset; //存放当前二维码扫描的资产
let assets; //存放当前地址二维码的资产列表
if (localStorage.status == undefined) localStorage.status = "0"; //localStorage.status设置的是允许进入的最大页号

//请登录提示
if (localStorage.usr == undefined) {
	$("#name_xs").html('<b style="color:MediumVioletRed ">请先登录</b>');
} else {
	usr = JSON.parse(localStorage.usr)
	$("#name_xs").text(usr["姓名"]);
}
////请选择摄像头提示
//if (localStorage.camera_id == undefined) {
//	$("#cameraid_c").html('<b style="color:MediumVioletRed ">请选择摄像头</b>');
//}else{
//	usr = JSON.parse(localStorage.usr)
//	$("#name_xs").text(usr["姓名"]);
//}
//1:第一次扫码<--------状态变更
//2:显示扫描设备信息
//3:扫描变更地址码<-----状态变更
//4:确认地址变更提交<---状态变更
//5:显示扫描地址信息



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
				allow_1();
			}
		});

});



//状态机
//function switchto0() {
//	$("#slide1").hide();
//	$("#slide2").hide();
//	$("#slide3").hide();
//	$("#slide4").hide();
//	$("#slide5").hide();
//	swiper.allowTouchMove = false;
//}
//
//function switchto1() {
//	$("#slide1").show();
//	$("#slide2").hide();
//	$("#slide3").hide();
//	$("#slide4").hide();
//	$("#slide5").hide();
//	swiper.allowTouchMove = true;
//}
//切换到指定状态
function switchstatus() {
	if ((swiper.activeIndex).toString() == localStorage.status) {
		swiper.allowSlideNext = false;
	} else {
		swiper.allowSlideNext = true;
	}
	console.log("swiper.allowSlideNext = ", swiper.allowSlideNext);
}
//switchstatus();

var swiper = new Swiper('.swiper-container', {
	on: {
		//		touchMove: function (event) { //在滑动的过程中连续发出
		//			console.log("event: touchMove");
		//			//swiper.allowTouchMove= false;
		//
		//		},
		//
		//		sliderMove: function (event) { //在滑动的过程中连续发出
		//			console.log("event: sliderMove");
		//			//swiper.allowTouchMove= false;
		//
		//		},
		slideNextTransitionStart: function () {
			if (swiper.activeIndex==2 && localStorage.status=="5") swiper.slideTo(5);
		},

		slidePrevTransitionStart: function () {
			if (swiper.activeIndex==4 && localStorage.status=="5") swiper.slideTo(1);
		},

		slideChange: function () {
			console.log("event: slideChange");
			console.log('改变了，activeIndex为' + this.activeIndex);

			switch (this.activeIndex) {
				case 0:
					x = "0 页 : 设置";
					scanner1.stop();
					scanner2.stop();
					scanner3.stop();
					switchstatus();
					break;
				case 1:
					x = "1 页 : 第一次扫码";
					//$("#slide2").hide();
					//swiper.allowTouchMove= false;//设置
					//console.log(swiper.allowTouchMove); //提示
					scanner1.stop();
/////////////////////////////////////////////////
scanner2.addListener('scan', function (data) {
	let content = utf82str(data);	//解决中文乱码
	alert("scanner2: "+content);
	
	//<----可以调用scanner.stop()结束扫描

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

					$("#asset").text(localStorage.asset);

					localStorage.status = 3; //可以进入第3页扫描设备更换地址码"
					switchstatus();
					swiper.slideNext();
					break;
				case "1": //扫描的是地点二维码
					assets = data.data;
					localStorage.assets = JSON.stringify(data.data);

					$("#assets").text(localStorage.assets);

					localStorage.status = 5; //可以进入第5页显示地址上设备信息
					switchstatus();
					swiper.slideTo(5);
					break;
				case "2": //错误的二维码
					alert(data.msg);
					break;
			}
		});


});					
/////////////////////////////////////////////////					
					scanner2.start(cameraslist[localStorage.camera_id]);
					scanner3.stop();
					switchstatus();
					break;
				case 2:
					x = "2 页 : 显示扫描设备信息 ";
					scanner1.stop();
					scanner2.stop();
					scanner3.stop();
					switchstatus();
					
					asset = JSON.parse(localStorage.asset);
					$("#p0").text(asset.资产名称);
					$("#p1").text("存放地点: " + asset.存放地点);
					break;
				case 3:
					x = "3 页 : 扫描设备更换地址码";
					scanner1.stop();
					scanner2.stop();
					scanner3.start(cameraslist[localStorage.camera_id]);
					switchstatus();
					break;
				case 4:
					x = "4 页 : 设备地址变更提交确认";
					scanner1.stop();
					scanner2.stop();
					scanner3.stop();
					switchstatus();
					break;
				case 5:
					x = "5 页 : 显示地址上设备信息";
					//$("#slide2").show();
					scanner1.stop();
					scanner2.stop();
					scanner3.stop();
					switchstatus();
					break;

			}
			console.log(x);
		},


	},
});

//判断是否进入第1页.  用户名及摄像头设定前禁止滑动
function allow_1() {
	if (!(localStorage.camera_id) || !(localStorage.usr)) {
		localStorage.status = "0";
	} else {
		localStorage.status = "1";
	};
	switchstatus();
}

//allow_1();


//////////  1,初始化第一个摄像头(设置页)
//			与<vedio>图像输出硬关联(初始化后不能更改), 
//			与摄像头软关联
let scanner1 = new Instascan.Scanner({
	video: document.getElementById('webcamera1'),
	continuous: true, //持续扫描
	mirror: false, //镜像,使用前置摄像头时使用
	captureImage: false, //将图像数据作为结果的一部分
	backgroundScan: false, //选项卡未激活时是否扫描
	refractoryPeriod: 5000, //连续识别相同QR码之前的时间
	scanPeriod: 1 //两次扫描之间的周期
});

let scanner2 = new Instascan.Scanner({
	video: document.getElementById('webcamera2'),
	continuous: true, //持续扫描
	mirror: false, //镜像,使用前置摄像头时使用
	captureImage: false, //将图像数据作为结果的一部分
	backgroundScan: false, //选项卡未激活时是否扫描
	refractoryPeriod: 5000, //连续识别相同QR码之前的时间
	scanPeriod: 1 //两次扫描之间的周期
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



scanner2.addListener('scan', function (data) {
	let content = utf82str(data);	//解决中文乱码
	alert("scanner2: "+content);
	
	//<----可以调用scanner.stop()结束扫描

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

					$("#asset").text(localStorage.asset);

					localStorage.status = 3; //可以进入第3页扫描设备更换地址码"
					switchstatus();
					swiper.slideNext();
					break;
				case "1": //扫描的是地点二维码
					assets = data.data;
					localStorage.assets = JSON.stringify(data.data);

					$("#assets").text(localStorage.assets);

					localStorage.status = 5; //可以进入第5页显示地址上设备信息
					switchstatus();
					swiper.slideTo(5);
					break;
				case "2": //错误的二维码
					alert(data.msg);
					break;
			}
		});


});

//初始化第二次扫描摄像头
let scanner3 = new Instascan.Scanner({
	video: document.getElementById('webcamera3'),
	continuous: true, //持续扫描
	mirror: false, //镜像,使用前置摄像头时使用
	captureImage: false, //将图像数据作为结果的一部分
	backgroundScan: false, //选项卡未激活时是否扫描
	refractoryPeriod: 5000, //连续识别相同QR码之前的时间
	scanPeriod: 1 //两次扫描之间的周期
});
scanner3.addListener('scan', function (data) {
	let content = utf82str(data);	//解决中文乱码
	alert("scanner3: "+content);
	//<----可以调用scanner.stop()结束扫描
});

//设置页, 摄像头列表
Instascan.Camera.getCameras().then(function (cameras) {

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
					<label><input type="radio" name="camera_op_c"  ` + checked + `  onclick="localStorage.camera_id=` + i + `;$('#cameraid_c').text(cameraslist[` + i + `].name);scanner1.stop();scanner1.start(cameraslist[` + i + `]);  allow_1(); "> ` + cameraslist[i].name + `</label>
				</div>`);
		};
	} else {
		alert('No cameras found.');
	}
}).catch(function (e) {
	alert(e);
});

switchstatus(); //进入指定状态


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

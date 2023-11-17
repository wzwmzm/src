// What it does:
//
// 	This program outputs the current OpenCV library version to the console.
//
// How to run:
//
// 		go run ./cmd/version/main.go
//

package main

import (
	"fmt"

	"gocv.io/x/gocv"
	// "gocv.io/x/gocv/openvino/ie"
	"image"
)

// 模块一：version()展示
func b_a() {
	fmt.Println("----模块一: version()演示--------------")
	fmt.Printf("opencv lib version: %s \n", gocv.OpenCVVersion())
	fmt.Printf("gocv version: %s\n", gocv.Version())
	// fmt.Printf("OpenVINO Inference Engine version: %s\n", ie.Version())

}

// 模块二：这个代码读取了一张图片，遍历了所有像素，并对RG两个通道的像素值进行了互换操作。最后，将修改后的图片显示。
func b_b() {
	fmt.Println("----模块二展示的功能:--------------")
	fmt.Println("----1, 读取本地图片 img := gocv.IMRead('AI.jpeg', gocv.IMReadColor) --------------")
	fmt.Println("----2, 读取MAT对象的类型 img.Type()--------------")
	fmt.Println("----3, 读取某像素的值(灰度)-img.GetUCharAt()-------------")
	fmt.Println("----4, 读取某像素的三通道值-img.GetVecbAt()-------------")
	fmt.Println("----5, 彩色照片转成灰度照片 gocv.CvtColor(img, &grayImage, gocv.ColorBGRToGray)	--------------")
	fmt.Println("----6, 设置某像素的值(灰度照) 	grayImage.SetUCharAt(row, col, 255)--------------")
	fmt.Println("----7, 灰度照片转成彩色照片 gocv.CvtColor(grayImage, &colorImage, gocv.ColorGrayToBGR)--------------")
	fmt.Println("----8, 设置某像素的值(彩色照) colorImage.SetUCharAt(100, col*3, 255) -----------")
	fmt.Println("----                       colorImage.SetUCharAt(100, col*3+1, 45) ----------")
	fmt.Println("----                       colorImage.SetUCharAt(100, col*3+2, 182)--------------")
	fmt.Println("----8, 保存图片 gocv.IMWrite('output.jpg', img)--------------")
	fmt.Println("----9, 扣图 croppedMat := img.Region(image.Rect(100, 100, 200, 200))--------------")
	fmt.Println("----10, 取得视频的某一帧 vc.Set(gocv.VideoCapturePosFrames, frames)--------------")
	//fmt.Println("----9, 扣图 croppedMat := img.Region(image.Rect(100, 100, 200, 200))--------------")
	//fmt.Println("----9, 扣图 croppedMat := img.Region(image.Rect(100, 100, 200, 200))--------------")
	//fmt.Println("----9, 扣图 croppedMat := img.Region(image.Rect(100, 100, 200, 200))--------------")
	//fmt.Println("----9, 扣图 croppedMat := img.Region(image.Rect(100, 100, 200, 200))--------------")
	//fmt.Println("----9, 扣图 croppedMat := img.Region(image.Rect(100, 100, 200, 200))--------------")

	// 读取图像
	img := gocv.IMRead("AI.jpeg", gocv.IMReadColor)
	defer img.Close()
	// 检查图像是否正确读取
	if img.Empty() {
		fmt.Println("无法读取图像")
		return
	}
	//显示原始图
	window := gocv.NewWindow("原始图像")
	defer window.Close()
	window.IMShow(img) //注意此句后面必须跟有WaitKey(),否则不会显示.

	// 获取图像的宽度和高度
	rows, cols := img.Rows(), img.Cols()

	//扣图
	//对croppedMat的修改会直接影响到原图, 所以才需要后面的克隆
	croppedMat := img.Region(image.Rect(int(cols/4), int(rows/4), int(cols*3/4), int(rows*3/4)))
	resultMat := (croppedMat).Clone()
	defer resultMat.Close()
	//显示扣图
	window3 := gocv.NewWindow("扣图")
	defer window3.Close()
	window3.IMShow(resultMat)

	fmt.Println("读取MAT对象的类型 img.Type() = " + img.Type().String() + "//type = CV8UC3 = MatTypeCV8U + MatChannels3")
	//对原图操作
	// fmt.Printf("ddd-color: GetUCharAt3 = %v, %v, %v\n", img.GetUCharAt3(100, 100, 0), img.GetUCharAt3(100, 100, 1), img.GetUCharAt3(100, 100, 2))
	bgr := img.GetUCharAt(100, 100) //
	bbb := img.GetVecbAt(100, 100)
	fmt.Printf("color: GetUCharAt = %v\n", bgr) //Before: BGR = 32
	fmt.Printf("color: GetVecbAt = %v\n", bbb)  //Before: bbb = [145 131 129]

	//由原始彩色图像转换成灰度图像并显示
	grayImage := gocv.NewMat()
	defer grayImage.Close()
	gocv.CvtColor(img, &grayImage, gocv.ColorBGRToGray)
	/////////将某一行设置成特定颜色
	for col := 0; col < cols; col += 1 {
		grayImage.SetUCharAt(100, col, 32)
		// grayImage.SetUCharAt(100, col*3+1, 45)
		// grayImage.SetUCharAt(100, col*3+2, 182)
	}
	//显示灰图
	window1 := gocv.NewWindow("gray Image")
	defer window1.Close()
	window1.IMShow(grayImage)
	//对灰图操作
	bgr = grayImage.GetUCharAt(100, 100)
	bbb = grayImage.GetVecbAt(100, 100)
	fmt.Printf("color-gray: GetUCharAt = %v\n", bgr) //Before: BGR = 32
	fmt.Printf("color-gray: GetVecbAt = %v\n", bbb)  //Before: bbb = [145 131 129]

	//由灰度图像转换成彩成图像并显示
	colorImage := gocv.NewMat()
	defer colorImage.Close()
	gocv.CvtColor(grayImage, &colorImage, gocv.ColorGrayToBGR)
	/////////将某一行设置成特定颜色
	for col := 0; col < cols; col++ { //三通道颜色顺序为B,G,R
		colorImage.SetUCharAt(100, col*3, 255)
		colorImage.SetUCharAt(100, col*3+1, 45)
		colorImage.SetUCharAt(100, col*3+2, 182)

		//mat.SetUCharAt3(),此函数会引起内存出错，不建议使用
		// colorImage.SetUCharAt3(200, col, 0, 255)
		// colorImage.SetUCharAt3(200, col, 1, 45)
		// colorImage.SetUCharAt3(200, col, 2, 182)

	}
	//显示彩色图
	window2 := gocv.NewWindow("color Image")
	defer window2.Close()
	window2.IMShow(colorImage)
	//对彩色图操作
	bgr = colorImage.GetUCharAt(100, 100)
	bbb = colorImage.GetVecbAt(100, 100)
	fmt.Printf("gray-color: GetUCharAt = %v\n", bgr) //Before: BGR = 32
	fmt.Printf("gray-color: GetVecbAt = %v\n", bbb)  //Before: bbb = [145 131 129]

	gocv.WaitKey(0) //注意!调用主体不是window!
	//0:表示无限等待键盘的输入,直到有输入才进入下一步.

	//保存图片
	gocv.IMWrite("output.jpg", colorImage)

}

// 模块四：摄像头操作的演示代码
func b_d() {
	fmt.Println("----模块四:摄像头操作的演示代码--------------")
	// 打开摄像头（默认0号摄像头）
	webcam, err := gocv.VideoCaptureDevice(0)
	if err != nil {
		fmt.Printf("Error opening capture device: %v\n", err)
		return
	}
	defer webcam.Close()

	// 申请一个窗口
	window := gocv.NewWindow("Webcam")
	defer window.Close()
	// 申请一个矩阵，用于存放图像数据
	img := gocv.NewMat()
	defer img.Close()
	for {
		// 读取摄像头数据
		if ok := webcam.Read(&img); !ok {
			fmt.Printf("Error: cannot read device %d\n", 0)
			return
		}
		if img.Empty() {
			continue
		}

		// 获取第100行第100列的像素值
		//bgr := img.AtVec3b(100, 100)

		// bgr := img.GetUCharAt(100, 100)
		// bbb := img.GetVecbAt(100, 100)
		// fmt.Printf("Before: BGR = %v\n", bgr) //Before: BGR = 32
		// fmt.Printf("Before: bbb = %v\n", bbb) //Before: bbb = [145 131 129]

		// 显示图像
		window.IMShow(img)          //注意!此句后面必须跟有WaitKey(),否则不会显示
		if window.WaitKey(1) >= 0 { //等待时间必须大于0,这里为1ms. 注意!与前面不同的是调用主体是window!
			break
		}
	}

}

// 获取视频的一帧
func GetVideoMoment(filePath string, time float64) { // time单位:秒
	//load video
	vc, err := gocv.VideoCaptureFile(filePath)
	if err != nil {
		return
	}
	defer vc.Close()

	frames := vc.Get(gocv.VideoCaptureFrameCount) //取得视频总帧数
	fps := vc.Get(gocv.VideoCaptureFPS)           //取得视频每秒刷新率
	totaltime := frames / fps                     //视频总时长,单位:秒

	frames = (time / totaltime) * frames //计算得到要取的是第几帧

	// Set Video frames
	vc.Set(gocv.VideoCapturePosFrames, frames) //设置以帧序号的方式获取图像

	//vc.Set(gocv.VideoCapturePosMsec, time/1000) //设置以时间方式获取图像!!!此方法可以取代上面复杂的计算

	img := gocv.NewMat()
	defer img.Close()
	vc.Read(&img)

	// imageObject, err := img.ToImage()
	// if err != nil {
	// 	return i, err
	// }

	w := gocv.NewWindow("视频截图")
	defer w.Close()
	w.IMShow(img)
	gocv.WaitKey(0)
	//return &img
}

func main() {
	// 模块一：version()展示
	b_a()
	// 模块二：常用功能实例展示
	b_b()
	//摄像头操作的演示代码
	b_d()
	//视频中截取第5秒的一帧
	GetVideoMoment("./vedio.webm", 5)
	//defer img
	//w := gocv.NewWindow("视频截图")
	//w.IMShow(*img)
	//gocv.WaitKey(0)

}

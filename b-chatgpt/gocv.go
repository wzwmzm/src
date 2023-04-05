package main

import (
    "fmt"
    "image"
    "image/color"

    "gocv.io/x/gocv"
)

func main() {
    // 读取图像
    img := gocv.IMRead("image.jpg", gocv.IMReadColor)

    // 检查图像是否正确读取
    if img.Empty() {
        fmt.Println("无法读取图像")
        return
    }

    // 获取图像的宽度和高度
    rows, cols := img.Rows(), img.Cols()

    // 遍历每个像素并进行操作
    for i := 0; i < rows; i++ {
        for j := 0; j < cols; j++ {
            // 获取像素的颜色
            pixel := img.At(i, j).(color.RGBA)

            // 将像素的红色和绿色分量互换
            pixel.R, pixel.G = pixel.G, pixel.R

            // 将像素的颜色设置回去
            img.Set(i, j, pixel)
        }
    }

    // 显示修改后的图像
    window := gocv.NewWindow("Modified Image")
    defer window.Close()

    window.IMShow(img)
    gocv.WaitKey(0)
}

// index.js
// 获取应用实例
const app = getApp()

Page({
  data: {
    wording: "world"
    },
  // 事件处理函数
  onClick: function(){
    this.setData({
      wording:"hello"
    })
  }
})
 
// miniprogram/components/cloudTipModal/index.js
const { isMac } = require('../../envList.js');

Component({

  /**
   * 页面的初始数据
   */
  data: {           //这里是组件私有数据，对外不可见
    showUploadTip: false,
    tipText: isMac ? 'sh ./uploadCloudFunction.sh' : './uploadCloudFunction.bat'
  },
  properties: {     //这里的属性才是对外的接口
    showUploadTipProps: Boolean
  },
  observers: {      //数据监听，这里当监听到有showUploadTipProps改变时执行操作。
    showUploadTipProps: function(showUploadTipProps) {
      this.setData({
        showUploadTip: showUploadTipProps
      });
    }
  },
  methods: {        //组件中方法
    onChangeShowUploadTip() {
      this.setData({
        showUploadTip: !this.data.showUploadTip
      });
    },

    copyShell() {
      wx.setClipboardData({ //拷贝数据到剪贴板
        data: this.data.tipText,
      });
    },
  }

});

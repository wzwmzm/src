Page({

  /**
   * 页面的初始数据
   */
  data: {
    showUploadTip: false,
    haveGetCodeSrc: false,
    envId: '',
    codeSrc: ''
  },

  onLoad(options) {
    this.setData({
      envId: options.envId
    });
  },

  getCodeSrc() {
    wx.showLoading({
      title: '',
    });
    wx.cloud.callFunction({
      name: 'quickstartFunctions',
      config: {
        env: this.data.envId
      },
      data: {
        type: 'getMiniProgramCode'
      }
    }).then((resp) => {
      this.setData({
        haveGetCodeSrc: true,
        codeSrc: resp.result
      });
      console.log("codeSrc: ",this.data.codeSrc," (云函数取到二维码图片数据后将其存储在云端，返回的是图片在云存储中的地址)");
      wx.hideLoading();
    }).catch((e) => {
      console.log(e);
      this.setData({
        showUploadTip: true
      });
      wx.hideLoading();
    });
  },

  clearCodeSrc() {
    this.setData({
      haveGetCodeSrc: false,
      codeSrc: ''
    });
  }

});

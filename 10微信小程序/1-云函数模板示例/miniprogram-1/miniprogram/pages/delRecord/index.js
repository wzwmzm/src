Page({

    /**
     * 页面的初始数据
     */
    record : new Array(),

    data: {
      showUploadTip: false,
      haveGetRecord: false,
      envId: '',
      record: ''
    },
  
    onLoad(options) {
      this.setData({
        envId: options.envId
      });
      wx.showLoading({
        title: '',
      });
      wx.cloud.callFunction({
        name: 'quickstartFunctions',
        config: {
          env: this.data.envId
        },
        data: {
          type: 'selectRecord'
        }
      }).then((resp) => {
        this.setData({
          record: resp.result.data
        });
        wx.hideLoading();
      }).catch((e) => {
        console.log(e);
        this.setData({
          showUploadTip: true
        });
        wx.hideLoading();
      });
    },
  
    delRecord() {
      wx.showLoading({
        title: '',
      });
      console.log("要删除的记录： ",this.record);
      wx.cloud.callFunction({
        name: 'quickstartFunctions',
        config: {
          env: this.data.envId
        },
        data: {
          type: 'delRecord',
          data: this.record
        }
      }).then((resp) => {
        wx.navigateTo({
          url: `/pages/delRecord/index?envId=${this.data.envId}`,
        });
        wx.hideLoading();
      }).catch((e) => {
        console.log("删除失败： ", e);
        this.setData({
          showUploadTip: true
        });
        wx.hideLoading();
      });
    },
  
    // bindInput (e) {
    //   const index = e.currentTarget.dataset.index;
    //   const record = this.data.record;
    //   record[index].sales = Number(e.detail.value);
    //   console.log( e.detail.value );
    //   this.setData({
    //     record
    //   });
    // },

    checkboxChange(e) {
        console.log('checkbox发生change事件，携带value值为：', e.detail.value);
        
        this.record = [];
        e.detail.value.forEach(element => {
            //console.log(element);
            this.record.push(this.data.record[element]);
            //console.log(record);
            
        });
        // console.log(this.record);
        // this.setData({
        //     record
        // })



    }
    
  
  });
  
// ==FeHelperMonkey==
// @reminder        请不要删除这部分代码注释，这是FeHelper油猴脚本能正常工作的基本条件！当然，你可以按需修改这里的内容！
// @id              mf_1598692144021
// @name            专为腾讯视频播放时关闭静音
// @url-pattern     https://v.qq.com/*
// @enable          true
// @require-js      
// @auto-refresh    0
// @updated         2020-08-30 11:09:30
// ==/FeHelperMonkey==


(function(){
  
  $(document).ready(function(){
    
    console.log("油猴开始了1111");
    console.log("油猴开始了----");
    console.log("油猴开始了----");    

});
  
  $(()=>{
    console.log("油猴开始了22222");
    console.log("油猴开始了=====");
    console.log("油猴开始了=====");
    console.log("油猴开始了=====");
    
  })
  
  
    var sleep = ms => new Promise((resolve,reject) => setTimeout(resolve,ms));
  
    sleep(1000).then(() => $('txpdiv.txp_btn.txp_btn_volume').click() ); //在页面载入完成后再过3秒关闭静音
  
  
    console.log("油猴开始了33333");
    console.log("油猴开始了!!!!!");
    console.log("油猴开始了!!!!!!");
  
 	//$('txpdiv.txp_btn.txp_btn_volume').click();	//在页面载入完成后关闭静音

})();
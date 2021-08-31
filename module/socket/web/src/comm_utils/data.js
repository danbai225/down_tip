// 字符串转16进制
function str2hex(str) {
    var result = [];
    var list = str.split("");
    for(var i=0;i<list.length;i++){
      if(i != 0){
        //加空格，分割16进制
        result.push(" ");
      }
      var item = list[i];
      //将字符串转化为16进制数据
      var Str = item.charCodeAt().toString(16);
      result.push(Str);
    }  
    return result.join("");
}
// 16进制转字符串
function hex2str(str) {
    var result = [];
    //
    //通过空格来分开16进制的字符
    var list = str.split(" ");
    for(var i=0;i<list.length;i++){
       var item = list[i];
       //转为asciicode 码
       var asciiCode = parseInt(item,16);
       //转为文字
       var charValue = String.fromCharCode(asciiCode);
       //添加到集合中
       result.push(charValue);
    }
    //返回结果
    return result.join("");
}
//将字符串转化为二进制的数据
function str2Binary(str){
    var result = [];
    var list = str.split("");
    for(var i=0;i<list.length;i++){
      if(i != 0){
        //加空格，分割二进制
        result.push(" ");
      }
      var item = list[i];
      //将字符串转化为二进制数据
      var binaryStr = item.charCodeAt().toString(2);
      result.push(binaryStr);
    }  
    return result.join("");
  }
   
  //二进制转为字符串
  function binary2Str(str){
    var result = [];
    //
    //通过空格来分开二进制的字符
    var list = str.split(" ");
    for(var i=0;i<list.length;i++){
       var item = list[i];
       //转为asciicode 码
       var asciiCode = parseInt(item,2);
       //转为文字
       var charValue = String.fromCharCode(asciiCode);
       //添加到集合中
       result.push(charValue);
    }
    //返回结果
    return result.join("");
  }
export {str2hex,hex2str,str2Binary,binary2Str};

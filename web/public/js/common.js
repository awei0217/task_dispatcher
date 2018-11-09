


/**
 * 判断入参是否为空
 * @param str
 * @returns {boolean} 空: 返回true; 非空: 返回false
 */
function isNull(str){
    if(str == null || str == ''|| str=="" || str == undefined ||typeof(str)=="undefined" ){
        return true;
    }
    return false;
}
/**
 * 判断入参是否非空
 * @param str
 * @returns {boolean} 非空: 返回true; 空: 返回false
 */
function isNotNull(str){
    if(isNull(str)){
        return false;
    }
    return true;
}



//时间戳的处理

layui.laytpl.toDateString = function (d, format) {
    if(d==null){
        return "";
    }
    var date = new Date(d || new Date())
        , ymd = [
        this.digit(date.getFullYear(), 4)
        , this.digit(date.getMonth() + 1)
        , this.digit(date.getDate())
    ]
        , hms = [
        this.digit(date.getHours())
        , this.digit(date.getMinutes())
        , this.digit(date.getSeconds())
    ];

    format = format || 'yyyy-MM-dd HH:mm:ss';
    if (format ==1){
        format ='yyyy-MM-dd';
    }
    return format.replace(/yyyy/g, ymd[0])
        .replace(/MM/g, ymd[1])
        .replace(/dd/g, ymd[2])
        .replace(/HH/g, hms[0])
        .replace(/mm/g, hms[1])
        .replace(/ss/g, hms[2]);
};

//数字前置补零
layui.laytpl.digit = function (num, length, end) {
    var str = '';
    num = String(num);
    length = length || 2;
    for (var i = num.length; i < length; i++) {
        str += '0';
    }
    return num < Math.pow(10, length) ? str + (num | 0) : num;
};





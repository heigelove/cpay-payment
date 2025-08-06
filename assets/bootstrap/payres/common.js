// var REQ_URL = "http://localhost:8092";
function getTopLevelDomain() {
    const hostname = window.location.hostname; // 获取当前主机名
    const parts = hostname.split('.'); // 按“.”分割主机名
    const length = parts.length;

    // 检查是否为顶级域名与国家级顶级域名进行一起处理
    if (length > 2) {
        // 处理可能的次级域名（如 www、mail 等）
        const secondLevelDomain = parts[length - 2]; // 取倒数第二部分
        const topLevelDomain = parts[length - 1]; // 取最后一部分
        const host = `${secondLevelDomain}.${topLevelDomain}`; // 返回 "example.co.uk" 形式
        if (host == 'yondermedia.shop' || host == 'raarinfotechce.shop' || host == 'keptofintechpri.shop' || host == 'pages.dev' || host == 'newboxshop.shop') {
            return 'fastpay.life'
        }else {
            return host;
        }
    } else {
        // 如果没有次级域名，直接返回主机名
        return hostname;
    }
}

var REQ_URL = "https://api." + getTopLevelDomain();

function toast(msg, duration) {
    duration = isNaN(duration) ? 3000 : duration;
    var m = document.createElement('div');
    m.innerHTML = msg;
    m.style.cssText = "padding:0.06rem; background:#000000; opacity:0.72; color:#FFFFFF; text-align:center; border-radius:0.2rem; position:fixed; top:50%; left:50%; transform:translate(-50%,-50%); -webkit-transform:translate(-50%,-50%); z-index:999; font-size: 0.3rem; min-width:50%; max-width: 70%;";
    document.body.appendChild(m);
    setTimeout(function () {
        var d = 0.5;
        m.style.webkitTransition = '-webkit-transform ' + d + 's ease-in, opacity ' + d + 's ease-in';
        m.style.opacity = '0';
        setTimeout(function () {
            document.body.removeChild(m)
        }, d * 1000);
    }, duration);
}

function copy(text) {
    var copyInput = document.createElement('input');
    document.body.appendChild(copyInput);
    copyInput.setAttribute('value', text);
    copyInput.select();
    document.execCommand("Copy");
    copyInput.remove();
    toast("Copy success", 1500);
}

function decryptDES(message) {
    var keyHex = CryptoJS.enc.Utf8.parse("65e1ccc7f4c7b717");
    var decrypted = CryptoJS.DES.decrypt(message, keyHex, {
        mode: CryptoJS.mode.ECB,
        padding: CryptoJS.pad.Pkcs7,
    });
    return decrypted.toString(CryptoJS.enc.Utf8);
}

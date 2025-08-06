function getQueryString(name) {
    let reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)", "i");
    let r = window.location.search.substr(1).match(reg);
    if (r != null) {
        return decodeURIComponent(r[2]);
    }
    return null;
}

function getParams() {
    let params = window.location.href.split("?");
    if (params.length > 1) {
        return params[1]
    }
}
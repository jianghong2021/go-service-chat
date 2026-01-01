function asEncrypt(cont, k) {
    return CryptoJS.AES.encrypt(cont, k + 'c8').toString();
}

function asDecrypt(cont, k) {
    return CryptoJS.AES.decrypt(cont, k + 'c8').toString(CryptoJS.enc.Utf8);
}

function getOssUrl(file) {
    return new Promise((resolve, reject) => {
        $.get({
            url: '/getOssUrl?file=' + file,
            success(res) {

                if (res.code == 200) {
                    const url = res.result;
                    resolve(url)
                } else {
                    reject(res.msg || '未知错误')
                }
            },
            error(err) {
                reject(err.message || '网络错误')
            }
        })
    })
}
import Http from "./request";

let http = new Http({
    baseURL: "/",
    timeout: 30000
});

let base = process.env.REACT_APP_BASE_URL

export async function getLink(){
    console.log(http)
    let data = await http.get(base+"/sign/");
    console.log(data.data.data.sign)
    let resp = await http.get(base+"/login/user/qrcode/generate")
    console.log(resp.data.result)
    let codeURL = "https://login.xuexi.cn/login/qrcommit?showmenu=false&code="+
        resp.data.result+"&appId=dingoankubyrfkttorhpou"
    return {"url":codeURL, "sign":data.data.data.sign,"code":resp.data.result}
}

export async function checkToken() {
    let token = window.localStorage.getItem("xxqg_token")
    if (token === null) {
        return false
    }
    let responseData = await http.post(base + "/auth/check/"+token);
    return responseData.data;

}




export async function login(data) {
    let responseData = await http.post(base+"/auth/login",data);
    return responseData.data;
}

export async function checkQrCode(code) {
    let data = new FormData();
    data.append("qrCode",code)
    data.append("goto","https://oa.xuexi.cn")
    data.append("pdmToken","")
    let resp = await http.post(base+"/login/login/login_with_qr",data,{
        headers: {
            "content-type":"application/x-www-form-urlencoded;charset=UTF-8"
        }
    })
    return resp.data
}

export async function getAbout(){
    let resp = await http.get(base+"/about");
    return resp.data;
}

export async function getToken(code,sign){
    let token = window.localStorage.getItem("xxqg_token")
    let resp = await http.post(base+"/user?register_id="+token,{
        "code":code,
        "state":sign
    });
    return resp.data;
}

export async function deleteUser(uid){
    let resp = await http.delete(base+"/user?uid="+uid);
    return resp.data;
}

export async function getUsers(){
    let resp = await http.get(base+"/user");
    return resp.data
}

export async function getConfig() {
    let resp = await http.get(base+"/config/file");
    return resp.data;
}

export async function restart() {
    let resp = await http.post(base+"/restart");
    return resp.data;
}

export async function update() {
    let resp = await http.post(base+"/update");
    return resp.data;
}

export async function saveConfig(data) {
    let resp = await http.post(base+"/config/file",{
        "data":data
    });
    return resp.data;
}

export async function getExpiredUsers(){
    let resp = await http.get(base+"/user/expired");
    return resp.data
}

export async function getScore(token) {
    let resp = await http.get(base+"/score?token="+token);
    return resp.data;
}


export async function study(uid) {
    let resp = await http.post(base+"/study?uid="+uid);
    return resp.data;
}

export async function stopStudy(uid) {
    let resp = await http.post(base+"/stop_study?uid="+uid);
    return resp.data;
}

export async function getLog(uid) {
    let resp = await http.get(base+"/log");
    return resp.data;
}

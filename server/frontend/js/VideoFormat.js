class Videoclass {
    constructor(channel, title, show, releaseDate, duration, link, pageLink, fileName) {
        this.channel = channel;
        this.title = title;
        this.show = show;
        this.releaseDate = releaseDate;
        this.duration = duration;
        this.link = link;
        this.pageLink = pageLink;
        this.fileName = fileName;
    }
}

function createAjaxRequest(){
    let request;
    if(window.XMLHttpRequest){
        request = new XMLHttpRequest();
    }else{
        request = new ActiveXObject("Microsoft.XMLHTTP");
    }
    return request;
}

function Logout() {
    window.location.href = "/index.html";
    sendPostLogoutRequest();
}

function openProfil() {
    window.location.href="/Profil.html";
}
function openHome() {
    window.location.href="/index.html";

}

function sendPostLogoutRequest(){
    const request = createAjaxRequest();
    request.onreadystatechange = function () {
        if(4 === this.readyState){
            if(200 === this.status){
                alert(this.responseText)

                hideAvatar();
            }else{
                alert(this.status + ":" + this.responseText);
            }
        }
    }
    request.open("POST",/logout/,true);
    request.setRequestHeader("Content-Type","application/x-www-form-urlencoded");
    request.send("dummy=dummy");
}





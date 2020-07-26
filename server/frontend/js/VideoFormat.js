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

function Logout() {
    sessionStorage.clear();
    window.location.href = "/index.html";
    sendGetLogoutRequest(function (status) {
        if(200 === status.status){
            alert(status.responseText)
            hideAvatar();
        }else{
            alert(status.status + ":" + status.responseText);
        }
    });
}

function openProfil() {
    window.location.href="/profile.html";
}
function openHome() {
    sessionStorage.clear();
    window.location.href="/index.html";
}

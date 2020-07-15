let channelJson;
let channelName;
let start =0;
let end = 30;
let siteNumber;


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

function loadSenderPage(wert) {
    window.location.href = "/senderpage.html";
    channelName = wert;
    sessionStorage.setItem('channel', wert);
    sendGetVideos();

}

function sendGetVideos() {
    const request = createAjaxRequest();
    request.onreadystatechange = function () {
        if (4 === this.readyState) {
            if (200 === this.status) {
                channelJson = JSON.parse(this.responseText);
                channelName = sessionStorage.getItem("channel");
                console.log(channelJson);
                setPage();
                siteNumber = channelJson.length/30;
                console.log(siteNumber);
            } else {
                alert(this.status + ":" + this.responseText);
            }
        }
    }
    request.open("GET", "/getVideoByChannel" + "?channel=" + sessionStorage.getItem('channel'), true);
    request.send();
}


function createAjaxRequest() {
    let request;
    if (window.XMLHttpRequest) {
        request = new XMLHttpRequest();
    } else {
        request = new ActiveXObject("Microsoft.XMLHTTP");
    }
    return request;
}

function setPage() {
    console.log("start");
    const videodiv = document.getElementById("videos");
    let currentVideo = new Videoclass("", "", "", "", "", "", "", "");
    let lastVideo;
    lastVideo = channelJson[0];

    let show =document.createElement("div");
    show.id = lastVideo.show;
    console.log(show);
     for (let i = start; i < end; i++) {
         currentVideo = channelJson[i];
        if (currentVideo.show !== lastVideo.show) {
           if(i>0) {
              videodiv.appendChild(show);
           }
           show =  document.createElement("div");
           show.id= currentVideo.show;
           console.log(show);
        }

        let element = document.createElement("input");
        element.id = currentVideo.title;
        element.type = "image";
        element.src = "media/Sender-Logos/zdf.png";
        element.className = "sender"
        element.value = JSON.stringify(currentVideo);
        show.appendChild(element);
        lastVideo = currentVideo;
    }
    videodiv.appendChild(show);
    console.log("Ende");


}


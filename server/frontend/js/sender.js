let channelJson;
let channelName;
let start =0;
let end = 30;

let currentPage =1;

let lastPage;

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

}

function sendGetVideos() {
    const request = createAjaxRequest();
    request.onreadystatechange = function () {
        if (4 === this.readyState) {
            if (200 === this.status) {
                channelJson = JSON.parse(this.responseText);
                channelName = sessionStorage.getItem("channel");
                console.log(channelJson);
                lastPage = Math.round(channelJson.length/end)+1;
                console.log(lastPage);
                setPage();

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
    /*Todo
    - Schleife braucht noch eine IF-Bedingung für die letzte Seite
     */

    let videosDiv = document.getElementById("videos");
    videosDiv.remove();
    videosDiv = document.createElement("div");
    videosDiv.id = "videos";
    const vContainer = document.getElementById("videoContainer");
    vContainer.appendChild(videosDiv);

    let currentVideo = new Videoclass("", "", "", "", "", "", "", "");
    let lastVideo;

    let show =  document.createElement("div");
    lastVideo = channelJson[start+((currentPage-1)*30)];
    show.id = lastVideo.show;
    show.className= "showClass";
    let t = document.createTextNode(lastVideo.show);
    show.appendChild(t);
    show.appendChild(document.createElement('br'));
    show.appendChild(document.createElement("hr"));
    appendShow(lastVideo,show);
    for(let i =(start+1)+((currentPage-1)*end);i<end*currentPage;i++){
        currentVideo = channelJson[i];
        if(lastVideo.show !== currentVideo.show){
            videosDiv.appendChild(show);
            show =  document.createElement("div");
            show.id =  currentVideo.show;
            show.className= "showClass";
            t = document.createTextNode(currentVideo.show);
            show.appendChild(t);
            show.appendChild(document.createElement('br'));
            show.appendChild(document.createElement("hr"));
        }
        appendShow(currentVideo,show);
        lastVideo = currentVideo;
    }
    videosDiv.appendChild(show);

}

function appendShow(video,showdiv){
    const videoDiv = document.createElement("div");
    const header5 = document.createElement("h5");
    const header7 = document.createElement("h6");
    const img = document.createElement("img");
    const a = document.createElement("a");
    a.href=JSON.stringify(video);
    videoDiv.setAttribute("class","videoLink");
    img.setAttribute("src","/media/Sender-Logos/"+video.channel+".png");
    img.setAttribute("class","thumbnail");
    videoDiv.appendChild(a);
    header5.innerHTML = video.title;
    header7.innerHTML = video.duration;
    videoDiv.appendChild(img);
    videoDiv.appendChild(header5);
    videoDiv.appendChild(header7);
    showdiv.appendChild(videoDiv);
}

function previousPage(){
    if((currentPage-1)<1);
    else {
        currentPage = currentPage - 1;
        setPage();
        document.getElementById("inputButton").value=JSON.stringify(currentPage);

    }
}
function nextPage() {
    currentPage=currentPage+1;
    setPage();
    document.getElementById("inputButton").value=JSON.stringify(currentPage);
}
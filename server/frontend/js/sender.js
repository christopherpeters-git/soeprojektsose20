let channelJson;
let channelName;
const start =0;
const end = 30;

let currentPage =1;

let lastPage;

function sendGetSearchRequestChannel(){
    const request = createAjaxRequest();
    const searchString = document.getElementById("searchInput").value;
    let channel = channelName;
    console.log(channel + "  "+ searchString);
    request.onreadystatechange = function () {
        if(4 === this.readyState){
            if(200 === this.status){
                alert(this.responseText);
                console.log("Suchfunktion Channelpage");
            }else{
                alert(this.status + ":" + this.responseText);
            }
        }
    }

    request.open("GET","/search" +"?search="+searchString+"&"+"channel="+channel,true);
    request.send();
}


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

function setSenderPagePicture(channel) {
    let img = document.createElement("img");
    img.setAttribute("src","/media/Sender-Logos/"+channel.channel+".png");
    img.setAttribute("id","senderPicture");
    const senderPagePic =document.getElementById("senderPic");
    senderPagePic.appendChild(img);
}

function sendGetVideos() {
    const request = createAjaxRequest();
    request.onreadystatechange = function () {
        if (4 === this.readyState) {
            if (200 === this.status) {
                channelJson = JSON.parse(this.responseText);
                if(channelJson===null){
                    window.location.href="/index.html";
                }
                channelName = sessionStorage.getItem("channel");
                console.log(channelJson.length);
                lastPage = (Math.ceil(channelJson.length/end));
                console.log(lastPage);
                setPage();
                setSenderPagePicture(channelJson[1]);

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
    let tempStart=start;
    let tempEnd =end;
    let videosDiv = document.getElementById("videos");
    videosDiv.remove();
    videosDiv = document.createElement("div");
    videosDiv.id = "videos";
    const vContainer = document.getElementById("videoContainer");
    vContainer.appendChild(videosDiv);

    let currentVideo = new Videoclass("", "", "", "", "", "", "", "");
    let lastVideo;
    if(currentPage === lastPage){
        if(channelJson.length<30){
            tempEnd=channelJson.length;
        }
        else {
            tempEnd = (lastPage * 30) - channelJson.length;
        }
    }
    let show =  document.createElement("div");
    lastVideo = channelJson[start+((currentPage-1)*30)];
    show.id = lastVideo.show;
    show.className= "showClass";
    let t = document.createTextNode(lastVideo.show);
    show.appendChild(t);
    show.appendChild(document.createElement('br'));
    show.appendChild(document.createElement("hr"));
    appendShow(lastVideo,show)

    for(let i =(tempStart+1)+((currentPage-1)*tempEnd);i<tempEnd*currentPage;i++){
        console.log(tempEnd);
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
    header5.className="videoTitle";
    const header7 = document.createElement("h6");
    header7.className="videoDuration"
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
    videoDiv.addEventListener("click",openVideoPlayer,false);
    videoDiv.value = video;
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
    if((currentPage+1)>lastPage);
    else {
        currentPage = currentPage + 1;
        setPage();
        document.getElementById("inputButton").value = JSON.stringify(currentPage);
    }
}

function openVideoPlayer() {
    console.log(this.value);

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

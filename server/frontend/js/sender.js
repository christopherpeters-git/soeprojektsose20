let channelJson;
let channelName;
const start =0;
const end = 30;

let currentPage =1;

let lastPage;


function loadSenderPage(wert) {
    window.location.href = "/channel.html";
    channelName = wert;
    sessionStorage.setItem('channel', wert);

}
function openHomePage() {
    window.location.href = "/index.html";

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
    appendShow(lastVideo,show,(start+((currentPage-1)*30)));

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
        appendShow(currentVideo,show,i);
        lastVideo = currentVideo;
    }
    videosDiv.appendChild(show);

}

function appendShow(video,showdiv,i){
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
    videoDiv.value = [video,i];
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
    sessionStorage.setItem('video', JSON.stringify(this.value));
    console.log(this.value);
    window.location.href = "/videoPlayer.html";
}


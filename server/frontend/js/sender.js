let channelJson;
let channelName;
let start =0;
let end = 30;

let currentPage =1;
let lastPage =10;


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
                addButtons();
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
    const videosdiv = document.getElementById("videos");
    let currentVideo = new Videoclass("", "", "", "", "", "", "", "");
    let lastVideo;
    let show =  document.createElement("div");
    lastVideo = channelJson[start];
    show.id = lastVideo.show;
    appendShow(lastVideo,show);
    for(let i =(start+1)+((currentPage-1)*end);i<end*currentPage;i++){
        currentVideo = channelJson[i];
        if(lastVideo.show !== currentVideo.show){
            videosdiv.appendChild(show);
            show =  document.createElement("div");
            show.id =  currentVideo.show;
        }
        appendShow(currentVideo,show);
        lastVideo = currentVideo;
    }
    videosdiv.appendChild(show);

}

function appendShow(video,showdiv){
    const videoDiv = document.createElement("div");
    const header5 = document.createElement("h5");
    const header7 = document.createElement("h6");
    const img = document.createElement("img");
    const a = document.createElement("a");
    a.href=JSON.stringify(video);
    videoDiv.setAttribute("class","videoLink");
    img.setAttribute("src","/media/Sender-Logos/ard.png");
    img.setAttribute("class","thumbnail");
    videoDiv.appendChild(a);
    header5.innerHTML = video.title;
    header7.innerHTML = video.duration;
    videoDiv.appendChild(img);
    videoDiv.appendChild(header5);
    videoDiv.appendChild(header7);
    showdiv.appendChild(videoDiv);
}

function addButtons() {
    const buttonDiv = document.getElementById("buttons");
    for(let i =currentPage; i<=lastPage;i++){
        let button = document.createElement("button");
        button.className= "senderPageButtons";
        button.value = JSON.stringify(i);
        button.textContent= JSON.stringify(i);
        button.addEventListener('click',setPageWithNumbers,false);
        buttonDiv.appendChild(button);
    }
}

function setPageWithNumbers() {
    setPageNumbers(this.value);
    setPage();
    let buttonSet;
    buttonSet= document.getElementsByClassName("senderPageButtons");
    console.log(buttonSet);
    for(let i =0;i<=buttonSet.length;i++){
        buttonSet[i].value = currentPage+i;
        buttonSet[i].textContent = JSON.stringify(currentPage+i);
    }
}

 async function setPageNumbers(value) {
    currentPage = parseInt(value);
    lastPage = currentPage + 10;
    console.log(lastPage);
    let videosDiv = document.getElementById("videos");
    videosDiv.remove();
    videosDiv = document.createElement("div");
    videosDiv.id = "videos";
    const vContainer = document.getElementById("videoContainer");
    vContainer.appendChild(videosDiv);
    console.log(videosDiv);
     await new Promise((res, rej) => {
         setTimeout(() => res("Now it's done!"), 300)
     });

    return 1;
}
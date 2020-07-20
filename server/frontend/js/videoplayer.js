let channel;
sendGetVideos();


function initVideoPlayer() {
    const videoPlayer =document.getElementById("my-video");
    let video = JSON.parse(sessionStorage.getItem("video"));
    videoPlayer.children[0].setAttribute("src",video[0].link);
    document.title =video[0].title;
    const videoTitle = document.getElementById("videoTitle");
    videoTitle.textContent=video[0].title;
    addVideoinformation(video);
    document.getElementById("nextVideos").innerHTML="";
    fillNextVideos(video);
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

function clearVideoPlayer() {
}

function addVideoinformation(video) {
    const videoClick = document.getElementById("videoClick");
    let clickNumber = sendGetClickedVideos(video);
    videoClick.textContent = clickNumber +" Aufrufe• " + video[0].releaseDate;
    let shareButton = document.createElement("button");
    shareButton.id= "shareButton";
    shareButton.value= video.pageLink;
    shareButton.addEventListener("click",shareThisVideo,false);
    shareButton.textContent = "➦ Teilen";
    videoClick.appendChild(shareButton);
    let addToFavoritBtn = document.createElement("button");

}

function sendGetClickedVideos(video){
    let clickNumber=1;
    const request = createAjaxRequest();
    request.onreadystatechange = function () {
        if(4 === this.readyState){
            if(200 === this.status){
                clickNumber = this.responseText;
            }else{
                alert(this.status + ":" + this.responseText);
            }
        }
    }
    request.open("GET","/clickVideo" +"?videoTitle="+video[0].title,false);
    request.send();
    return clickNumber;
}

function shareThisVideo(){
    console.log(this.value);
}

function sendGetVideos() {
    const request = createAjaxRequest();
    request.onreadystatechange = function () {
        if (4 === this.readyState) {
            if (200 === this.status) {
                channel = JSON.parse(this.responseText);
                if(channel===null) {
                    window.location.href = "/index.html";
                }
            } else {
                alert(this.status + ":" + this.responseText);
            }
        }
    }
    request.open("GET", "/getVideoByChannel" + "?channel=" + sessionStorage.getItem('channel'), false);
    request.send();
}

function fillNextVideos(video) {
    const nxtVideos =document.getElementById("nextVideos");
    const start = video[1];
    let end=10+start;
    if((channel.length-start)<10){
        console.log((channel.length))
        end = channel.length-start;
    }
    if(channel.length<10){
        end =channel.length;
    }
    for(let i = start+1;i<end;i++) {
        const videoDiv = document.createElement("div");
        const header5 = document.createElement("h5");
        header5.className = "videoTitle";
        const header7 = document.createElement("h6");
        header7.className = "videoDuration"
        const img = document.createElement("img");
        const a = document.createElement("a");
        a.href = JSON.stringify(channel[i]);
        videoDiv.setAttribute("class", "videoLink");
        img.setAttribute("src", "/media/Sender-Logos/" + channel[i].channel + ".png");
        img.setAttribute("class", "thumbnail");
        videoDiv.appendChild(a);
        header5.innerHTML = channel[i].title;
        header7.innerHTML = channel[i].duration;
        videoDiv.appendChild(img);
        videoDiv.appendChild(header5);
        videoDiv.appendChild(header7);
        videoDiv.addEventListener("click", openVideoPlayer, false);
        videoDiv.value = [channel[i], i];
        nxtVideos.appendChild(videoDiv);
    }
}


function openVideoPlayer() {
    sessionStorage.setItem('video', JSON.stringify(this.value));
    initVideoPlayer();
}
let channel;
sendGetVideos();


function initVideoPlayer() {
    const videoPlayer =document.getElementById("my-video");
    let video = JSON.parse(sessionStorage.getItem("video"));
    videoPlayer.setAttribute("poster","media/Sender-Logos/"+video[0].channel+".png");
    videoPlayer.children[0].setAttribute("src",video[0].link);
    document.title =video[0].title;
    const videoTitle = document.getElementById("videoTitle");
    videoTitle.textContent=video[0].title;
    addVideoinformation(video);
    document.getElementById("nextVideos").innerHTML="";
    fillNextVideos(video);
    document.getElementById("moreInformation").innerHTML="";
    setMoreInformation(video);
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
    shareButton.value=JSON.stringify(video[0].pageLink);
    videoClick.appendChild(shareButton);
    let addToFavoritBtn = document.createElement("button");
    addToFavoritBtn.id = "Favbtn";
    addToFavoritBtn.textContent = "❤";
    addToFavoritBtn.value=JSON.stringify(video[0]);
    addToFavoritBtn.addEventListener("click",addVideoToFav,false);
    videoClick.appendChild(addToFavoritBtn);
    videoClick.appendChild(document.createElement("br")); videoClick.appendChild(document.createElement("br"));


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

function setMoreInformation(video) {
    const infoDiv = document.getElementById("moreInformation");
    let informationSet = document.createElement("div");
    const img = document.createElement("img");
    img.setAttribute("id","infoPic");
    img.setAttribute("src", "/media/Sender-Logos/" + video[0].channel + ".png");
    const header5 = document.createElement("h5");
    header5.id = "InfoTitle";
    header5.innerHTML = "Channel: "+video[0].channel+" Show: "+video[0].show;
    let tempdiv = document.createElement("div");
    tempdiv.id = "Info";
    let header1 = document.createElement("h1");
    let header2 = document.createElement("h2");
    let header3 = document.createElement("h3");
    const a = document.createElement("a");
    header1.innerHTML = "Titel: "+video[0].title;
    header2.innerHTML = "Dauer: "+video[0].duration;
    header3.innerHTML = "Seitenlink: "+video[0].pageLink;
    tempdiv.appendChild(header1);
    tempdiv.appendChild(header2);
    tempdiv.appendChild(header3);
    informationSet.appendChild(img);
    informationSet.appendChild(header5);
    informationSet.appendChild(tempdiv);
    infoDiv.appendChild(informationSet);
}

function openVideoPlayer() {
    sessionStorage.setItem('video', JSON.stringify(this.value));
    initVideoPlayer();
}
function addVideoToFav() {
    sendPostFavoriteRequest(this.value)

}

function sendPostFavoriteRequest(video){
    const request = createAjaxRequest();
    request.onreadystatechange = function () {
        if(4 === this.readyState){
            if(200 === this.status){
                alertSetterFunction("rgba(39,255,0,0.75)",this.responseText);
            }else{
                alertSetterFunction("rgba(255,0,30,0.75)",this.responseText);
            }
            console.log(this);
        }
    }
    request.open("POST",/addToFavorites/,true);
    request.setRequestHeader("Content-Type","application/x-www-form-urlencoded");
    request.send("video="+JSON.stringify(video));
}

function alertSetterFunction(color,message) {
    const alert = document.getElementById("alert");
    alert.textContent= message;
    alert.style.background=color;
    alert.style.display="block"
    setTimeout(function(){alert.style.display="none"},1500);
}
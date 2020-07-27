let channel;
getFlagGetVideos();


function setDefaultAutplay() {
    const slider = document.getElementsByClassName("switch");
    slider[0].children[0].checked=false;
}

function getFlagGetVideos() {
    console.log("flag: "+ sessionStorage.getItem(("favFlag")));
    if(sessionStorage.getItem(("favFlag"))==="1"){
        sendGetFetchFavoritesRequest((response)=>{
            channel = JSON.parse(response.responseText);
        },false)
    }
    else{
        const flag = JSON.parse(sessionStorage.getItem("pageFlag"));
        if(flag===0){
            sendGetVideos(callBackFunctionGetVideosForVideoPlayer,false);
        }
        else if(flag===1){
            console.log("pageFlag: "+flag);
            sendGetSearchRequest(function (status) {
                if(200 === status.status){
                    channel = JSON.parse(status.responseText);

                }else{
                    console.log(status.status + ":" + status.responseText);
                }
            },false);
        }

    }
}

function initVideoPlayer() {
    const videoPlayer =document.getElementById("my-video")
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


function clearVideoPlayer() {
    let myPlayer = document.getElementById("my-video");
    myPlayer = myPlayer.children[0];
    myPlayer.removeEventListener("ended",autoPlayFunction,false);

}

function addVideoinformation(video) {
    const videoClick = document.getElementById("videoClick");
    let clickNumber = sendGetClickedVideos(video,false);
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



function shareThisVideo(event){
    alertSetterFunction("#cccccc",this.value+" kopiert in Zwischenablage",3000);
    let temInput = document.getElementById("tempInput");
    temInput.value=this.value;
    temInput.style.display="block";
    let copyText = document.getElementById("tempInput");
    copyText.select();
    copyText.setSelectionRange(0, 99999)
    document.execCommand("copy");
    console.log(copyText.value)
    temInput.style.display="none";
}

function callBackFunctionGetVideosForVideoPlayer(status){
    if (200 === status.status) {
        channel = JSON.parse(status.responseText);
        if(channel===null) {
            window.location.href = "/index.html";
        }
    } else {
        alert(status.status + ":" + status.responseText);
    }
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
        videoDiv.appendChild(document.createElement("br"));
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
    let myPlayer = document.getElementById("my-video");
    myPlayer = myPlayer.children[0];
    myPlayer.play();

}
function addVideoToFav() {
    console.log(this.value);
    sendPostFavoriteRequest(encodeURIComponent(this.value));
}



function alertSetterFunction(color,message,timeout) {
    const alert = document.getElementById("alert");
    alert.textContent= message;
    alert.style.background=color;
    alert.style.display="block"
    if(timeout>1500) {
        let spanBtn = document.createElement("button");
        spanBtn.className = "closebtn";
        spanBtn.textContent = "✖"
        spanBtn.addEventListener("click", closeAlert, false);
        alert.appendChild(spanBtn);
    }
    setTimeout(function(){alert.style.display="none"},timeout);
}
function closeAlert() {
    this.parentElement.style.display='none';
}


function toggleAutoplayVideoplayer() {
    const slider = document.getElementsByClassName("switch");
    let myPlayer = document.getElementById("my-video");
    myPlayer = myPlayer.children[0];
    if(slider[0].children[0].checked){
        myPlayer.addEventListener("ended",autoPlayFunction,false);
    }
    else {
        myPlayer.removeEventListener("ended",autoPlayFunction,false);
    }
}

function autoPlayFunction() {
    console.log("started new Video");
    let listVideos = document.getElementById("nextVideos");
    if(listVideos!==null) {
        listVideos=listVideos.children[0];
        sessionStorage.setItem('video', JSON.stringify(listVideos.value));
        initVideoPlayer();
        let myPlayer = document.getElementById("my-video");
        myPlayer = myPlayer.children[0];
        myPlayer.play();
    }
}




let channelJson=null;
let channelName=null;
const start =0;
const end = 30;
let currentPage =1;
let lastPage;



function loadSenderPage() {
    let wert = this.value;
    //set pageFlag with 0
    sessionStorage.setItem("pageFlag","0");
    //Open channel.html
    window.location.href = "/channel.html";
    //set channelName with value of caller
    channelName = wert;
    //set channel with value of caller
    sessionStorage.setItem('channel', wert);
}

//sets the Picture of the channel thats calling the channel.html
function setSenderPagePicture(channel) {
    let img = document.createElement("img");
    img.setAttribute("src","/media/Sender-Logos/"+channel.channel+".png");
    img.setAttribute("id","senderPicture");
    const senderPagePic =document.getElementById("senderPic");
    senderPagePic.appendChild(img);
}

//callbackfunction of a Request
function callBackFunctionGetVideos(status){
    if (200 === status.status) {
        channelJson = JSON.parse(status.responseText);
        if(channelJson===null){
            //if channel Json == NULL
                //go back to Mainpage
            window.location.href="/index.html";
        }
        channelName = sessionStorage.getItem("channel");
        //getting the lastPage of channel
        lastPage = (Math.ceil(channelJson.length/end));
        setPage();
        setSenderPagePicture(channelJson[1]);

    } else {
        console.log(status.status + ":" + status.responseText);
    }
}

//callback Function of a Searchrequest
function callbackFunctionSetSearchRequest(status) {
    if(200 === status.status){
        channelJson = JSON.parse(status.responseText);
        lastPage = (Math.ceil(channelJson.length/end));
        setPage();

    }else{
        console.log(status.status + ":" + status.responseText);
    }
}



//sets the channel.html
function setPage() {
    if(channelJson.length<1){
       return 0;
    }
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
    //checks how many sites are left
    if(currentPage === lastPage){
        if(channelJson.length<30){
            tempEnd=channelJson.length;
        }
        else {
            tempEnd = (lastPage * 30) - channelJson.length;
        }
    }
    //adding a div for every show
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

//appending Videolinks from the same Show
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
    if(sessionStorage.getItem("pageFlag")==="0") {
        videoDiv.addEventListener("click", openVideoPlayerWithPageResults, false);
    }else if(sessionStorage.getItem("pageFlag")==="1"){
        videoDiv.addEventListener("click", openVideoPlayerWithSearchResults, false);
    }
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

// open Videoplayer
function openVideoPlayerWithSearchResults() {
    //set favFlag = 0
    sessionStorage.setItem("favFlag","0");
    //set video = value of caller
    sessionStorage.setItem('video', JSON.stringify(this.value));
    console.log(this.value);
    //open Videoplayer.html
    window.location.href = "/videoPlayer.html";
}

function openVideoPlayerWithPageResults() {
    //set favFlag = 0
    sessionStorage.setItem("favFlag","0");
    //set video = value of caller
    sessionStorage.setItem('video', JSON.stringify(this.value));
    console.log(this.value);
    //open Videoplayer.html
    window.location.href = "/videoPlayer.html";
}


//call Function depending on Flags
function checkFlag() {
    const flag = JSON.parse(sessionStorage.getItem("pageFlag"));
    console.log(flag);
    if(flag===0){
        sendGetVideos(callBackFunctionGetVideos);
    }else if(flag===1){
        sendGetSearchRequest(callbackFunctionSetSearchRequest);
    }
}

//sets the Searchtext for a search on a Channel
function setSearchTextWithChannel() {
    console.log(channelName);
    const searchValue = document.getElementById("searchInput").value;
    if(searchValue==="") return;
    let searchString;
    searchString= JSON.stringify([channelName, searchValue]);
    sessionStorage.setItem("searchString",searchString);
    window.location.href = "/searchResults.html";
}

function addEventListenerSenderPage() {
    const  searchInput = document.getElementById("searchInput");
    searchInput.addEventListener("keydown",searchOnEnter,false);

}

function searchOnEnter(event) {
    if(event.key ==="Enter") {
        setSearchTextWithChannel();
    }
}
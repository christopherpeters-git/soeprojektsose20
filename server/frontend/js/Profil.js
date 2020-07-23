let userInformation=sendPostCookieAuthRequest();
function Logout() {
    window.location.href = "/index.html";
    sendPostLogoutRequest();

}
function openHome() {
    window.location.href = "/index.html";
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

function sendPostCookieAuthRequest(){
    const request = createAjaxRequest();
    request.onreadystatechange = function () {
        if(4 === this.readyState){
            if(200 === this.status){
                userInformation=JSON.parse(this.responseText);
                console.log(userInformation);
                setFavList();
            }else{
                console.log(this.status + ":" + this.responseText);
                document.getElementById("Login_Screen").style.visibility="visible";

            }
        }
    }
    request.open("POST","/cookieAuth/",true);
    request.setRequestHeader("Content-Type","application/x-www-form-urlencoded");
    request.send("dummy=dummy");
}

function setFavList(){
    const favDiv= document.getElementById(" favorites");
    for(let i =0;i<userInformation.favoriteVideos.length;i++){
        appendFav(userInformation.favoriteVideos[i],favDiv,i,"openFavVideoPlayer");
    }
}


function appendFav(video,showdiv,i){
    const videoDiv = document.createElement("div");
    const header5 = document.createElement("h5");
    header5.className="videoTitle";
    const header7 = document.createElement("h6");
    header7.className="videoDuration"
    const img = document.createElement("img");
    const a = document.createElement("a");
    const checkBox = document.createElement("input");
    checkBox.type = "checkbox";
    checkBox.className = "checkBoxFav";
    a.href=JSON.stringify(video);
    videoDiv.setAttribute("class","videoLink");
    img.setAttribute("src","/media/Sender-Logos/"+video.channel+".png");
    img.setAttribute("class","thumbnail");
    videoDiv.appendChild(checkBox);
    videoDiv.appendChild(a);
    header5.innerHTML = video.title;
    header7.innerHTML = video.duration;
    videoDiv.appendChild(img);
    videoDiv.appendChild(header5);
    videoDiv.appendChild(header7);
    videoDiv.addEventListener("click",openFavVideoPlayer,false);
    videoDiv.value = [video,i];
    showdiv.appendChild(videoDiv);
}

function openFavVideoPlayer() {
    alert("test");

}

function tester(event) {

}
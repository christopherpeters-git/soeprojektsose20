let userInformation=sendGetCookieAuthRequest();
function Logout() {
    window.location.href = "/index.html";
    sendGetLogoutRequest();

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

function sendPostSaveProfilePicture(){
    const request = createAjaxRequest();
    const profilePicture = document.getElementById("ppUpload").files[0];
    if (profilePicture == null){
        alert("No picture set!");
        return;
    }
    const formData = new FormData();
    formData.append("profilepicture",profilePicture)
    request.onreadystatechange = function () {
        if(4 === this.readyState){
            if(200 === this.status){
                alert(this.responseText)
                loadProfilePicture()
            }else{
                alert(this.status + ":" + this.responseText);
            }
        }
    }
    request.open("POST","/setProfilePicture/",true);
    request.send(formData);
}

function loadProfilePicture(){
    document.getElementById("profilePicture").src = /getProfilePicture/
    document.getElementById("profilePictureMini").src = /getProfilePicture/
}

function sendGetFetchFavoritesRequest(){
    const request = createAjaxRequest();
    request.onreadystatechange = function () {
        if(4 === this.readyState){
            if(200 === this.status){
                userInformation.favoriteVideos = JSON.parse(this.responseText);
                console.log(userInformation.favoriteVideos);
            }else{
                alert(this.status + ":" + this.responseText);
            }
        }
    }
    request.open("GET","/getFavorites/",true);
    request.send();
}


function sendGetCookieAuthRequest(){
    const request = createAjaxRequest();
    request.onreadystatechange = function () {
        if(4 === this.readyState){
            if(200 === this.status){
                userInformation=JSON.parse(this.responseText);
                sendGetFetchFavoritesRequest();
                setFavList();
            }else{
                console.log(this.status + ":" + this.responseText);
                document.getElementById("Login_Screen").style.visibility="visible";

            }
        }
    }
    request.open("GET","/cookieAuth/",true);
    request.send();
}

function setFavList(){
    const favDiv= document.getElementById("favorites");
    favDiv.innerHTML="";
    favDiv.textContent="Favoriten";
    const editBtn =document.createElement("button");
    editBtn.textContent= "✎";
    editBtn.id= "editBtn";
    editBtn.className="favBtn";
    editBtn.addEventListener("click",setFav,false);
    const safeBtn =document.createElement("button");
    safeBtn.textContent= "✔";
    safeBtn.className="favBtn";
    safeBtn.id="safeBtn";
    safeBtn.addEventListener("click",startDeletingFav,false);
    const abortBtn =document.createElement("button");
    abortBtn.textContent= "✖";
    abortBtn.className="favBtn";
    abortBtn.id ="abortBtn";
    abortBtn.addEventListener("click",abortFav,false);
    const selectAllBtn =document.createElement("button");
    selectAllBtn.textContent= "All";
    selectAllBtn.className="favBtn";
    selectAllBtn.id ="selectAllBtn";
    selectAllBtn.addEventListener("click",selectAllFavorites,false);
    favDiv.appendChild(editBtn);
    favDiv.appendChild(document.createElement("br"));
    favDiv.appendChild(selectAllBtn);
    favDiv.appendChild(safeBtn);
    favDiv.appendChild(abortBtn);
    favDiv.appendChild(document.createElement("hr"));
    if (userInformation.favoriteVideos != null){
        for(let i =0;i<userInformation.favoriteVideos.length;i++){
            appendFav(userInformation.favoriteVideos[i],favDiv,i,"openFavVideoPlayer");
        }
        favDiv.appendChild(document.createElement("hr"));
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
    sessionStorage.setItem("favFlag","1");
    sessionStorage.setItem('video', JSON.stringify(this.value));
    console.log(this.value);
    window.location.href = "/videoPlayer.html";
}

function setFav(event) {
    const favList = document.getElementById("favorites");
    const checkedBtn = document.getElementsByClassName("checkBoxFav");
    const editBtn = document.getElementById("editBtn");
    const safeBtn = document.getElementById("safeBtn");
    const abortBtn = document.getElementById("abortBtn");
    const selectAll = document.getElementById("selectAllBtn");

    for(let i =0;i<checkedBtn.length;i++){
       checkedBtn[i].style.visibility="visible";
    }
    for(let i =0;i<favList.children.length;i++){
        favList.children[i].removeEventListener("click",openFavVideoPlayer,false);
    }
    editBtn.style.visibility="hidden";
    selectAll.style.visibility="visible";
    safeBtn.style.display="block";
    abortBtn.style.display="block";

}

function abortFav() {
    const editBtn = document.getElementById("editBtn");
    const safeBtn = document.getElementById("safeBtn");
    const abortBtn = document.getElementById("abortBtn");
    const checkedBtn = document.getElementsByClassName("checkBoxFav");
    editBtn.style.visibility="visible";
    safeBtn.style.display="none";
    abortBtn.style.display="none";
    for(let i =0;i<checkedBtn.length;i++){
        checkedBtn[i].style.visibility="hidden";
    }

}
function startDeletingFav() {
    const listFav = document.getElementById("favorites");
    console.log(listFav.children[5].value);
    for(let i =6;i<listFav.children.length-1;i++){
        if(listFav.children[i].children[0].checked){
           sendPostRemoveFavoriteRequest(listFav.children[i].value[0]);
        }
    }
    /*todo
    * delay after request
    * */
    setTimeout(() => { location.reload();  }, 200);
}


function sendPostRemoveFavoriteRequest(video){
    const request = createAjaxRequest();
    request.onreadystatechange = function () {
        if(4 === this.readyState){
            if(200 === this.status){
            }else{
                alert(this.status + ":" + this.responseText);
            }
            console.log(this);
        }
    }
    request.open("POST",/removeFromFavorites/,true);
    request.setRequestHeader("Content-Type","application/x-www-form-urlencoded");
    request.send("video="+encodeURIComponent(JSON.stringify(video)));
}

function selectAllFavorites() {
    const listFav = document.getElementById("favorites");
    console.log(listFav);
    for(let i =6;i<listFav.children.length-1;i++){
        listFav.children[i].children[0].checked = true;
    }
}
let userInformation;
sendGetCookieAuthRequest(callBackFunctionSetUserArray);
sendGetFetchFavoritesRequest((response) => {
    console.log(response.responseText)
    userInformation.favoriteVideos = JSON.parse(response.responseText);
    setFavList();
});


// Function to show the profile informations (Name, Username).
function displayProfileInformation(){
    const name = document.getElementById("displayName");
    const username = document.getElementById("displayUsername");
    const nameText = document.createElement("h3");
    const usernameText = document.createElement("h3");
    nameText.innerHTML = "Name: " + userInformation.name;
    usernameText.innerHTML = "Username: " + userInformation.username;
    name.appendChild(nameText);
    username.appendChild(usernameText);
}
//Button to upload a profile picture.
function toggleUploadButtons(){
    const uploadButtons = document.getElementById("uploadButtons");
    if (uploadButtons.style.visibility === "visible"){
        uploadButtons.style.visibility = "hidden";
    }else{
        uploadButtons.style.visibility = "visible";
    }
}
//Save the selected profile picture.
function saveProfilePicture(){
    sendPostSaveProfilePicture(()=>{
        loadProfilePicture();
        console.log("profile picture set successfully")
    });
}
// Load the selected profile picture.
function loadProfilePicture(){
    const pp = document.getElementById("profilePicture");
    pp.setAttribute("src","/getProfilePicture/");
    window.location.reload();
}

function Logout() {
    window.location.href = "/index.html";
    sendGetLogoutRequest();
}
function openHome() {
    window.location.href = "/index.html";
}

function callBackFunctionSetUserArray(status) {
    if(200 === status.status){
        userInformation=JSON.parse(status.responseText);
        displayProfileInformation();
    }else{
        console.log(status.status + ":" + status.responseText);
    }
}

// Function to create a favorite list.
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
    for(let i =0;i<userInformation.favoriteVideos.length;i++){
        appendFav(userInformation.favoriteVideos[i],favDiv,i,"openFavVideoPlayer");
    }
    favDiv.appendChild(document.createElement("hr"));
}

// Function to select a video, in the favorite list, with a checkbox.
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
// Function to open the Videoplayer page.
function openFavVideoPlayer() {
    sessionStorage.setItem("favFlag","1");
    sessionStorage.setItem('video', JSON.stringify(this.value));
    console.log(this.value);
    window.location.href = "/videoPlayer.html";
}
// Function to put videos in the favorite list.
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
// function for the three favorite list Button (Edit, Safe and Abort).
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
// Function to delete videos in the favorite list.
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
// Function to select all videos in the favorite list.
function selectAllFavorites() {
    const listFav = document.getElementById("favorites");
    console.log(listFav);
    for(let i =6;i<listFav.children.length-1;i++){
        listFav.children[i].children[0].checked = true;
    }
}
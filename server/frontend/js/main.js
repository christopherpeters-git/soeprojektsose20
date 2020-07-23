"use strict"
let searchResultsJSON;
let pageFlag;

class User{
    constructor(id,name,username,favoriteVideos) {
        this.id = id;
        this.username = username;
        this.name = name;
        this.favoriteVideos= favoriteVideos;
    }

}
class Video {
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

const video =  {
    "channel": "ARD",
    "title": "\"Plan B\" für Bayern",
    "show": "\"Plan B\" für Bayern",
    "releaseDate": "24.06.2020",
    "duration": "00:42:55",
    "link": "http://cdn-storage.br.de/b7/2020-06/24/159a1622b65711eabca2984be109059a_C.mp4",
    "pageLink": "https://www.ardmediathek.de/ard/player/Y3JpZDovL2JyLmRlL3ZpZGVvLzY3ZmY5YjUyLTI5NDYtNDEwMC04MDk1LTg2OTU1NjgxOTMyZA",
    "fileName": "72|X.mp4"
};


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
                hideVBlockerAndLogin();
                unhideAvatar();
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

function sendGetSearchRequest(){
    const request = createAjaxRequest();
    console.log("SearchRequest: " + sessionStorage.getItem("search"));
    const searchString = sessionStorage.getItem("search");
    let channel = "none";
    console.log(channel + "  "+ searchString);
    request.onreadystatechange = function () {
        if(4 === this.readyState){
            if(200 === this.status){
                searchResultsJSON = JSON.parse(this.responseText);
                setPage(searchResultsJSON);
            }else{
                alert(this.status + ":" + this.responseText);
            }
        }
    }

    request.open("GET","/search" +"?search="+searchString+"&"+"channel="+channel,true);
    request.send();
}

function sendGetSearchRequestSearchResults(){
    const request = createAjaxRequest();
    const searchString = document.getElementById("searchInput").value;
    let channel = "none";
    request.onreadystatechange = function () {
        if(4 === this.readyState){
            if(200 === this.status){
                searchResultsJSON = JSON.parse(this.responseText);
                setPage(searchResultsJSON);
            }else{
                alert(this.status + ":" + this.responseText);
            }
        }
    }

    request.open("GET","/search" +"?search="+searchString+"&"+"channel="+channel,true);
    request.send();
}

function sendPostLoginRequest(){
    const usernameInput = document.getElementById("usernameLogin").value;
    const passwordInput = document.getElementById("passwordLogin").value;
    const request = createAjaxRequest();
    request.onreadystatechange = function () {
        if(4 === this.readyState){
            if(200 === this.status){
                alert(this.responseText)
                let user = new User("","","",null);
                user = JSON.parse(this.responseText);
                console.log(user);
                hideVBlockerAndLogin();
                unhideAvatar();
            }else{
                alert(this.status + ":" + this.responseText);
            }
        }
    }
    request.open("POST",/login/,true);
    request.setRequestHeader("Content-Type","application/x-www-form-urlencoded");
    console.log("" + usernameInput + " " + passwordInput);
    request.send("usernameInput="+usernameInput+"&"+"passwordInput="+passwordInput);
}

function sendPostLogoutRequest(){
    const request = createAjaxRequest();
    request.onreadystatechange = function () {
        if(4 === this.readyState){
            if(200 === this.status){
                alert(this.responseText)
                unhideVBlockerAndLogin();
                hideAvatar();
            }else{
                alert(this.status + ":" + this.responseText);
            }
        }
    }
    request.open("POST",/logout/,true);
    request.setRequestHeader("Content-Type","application/x-www-form-urlencoded");
    request.send("dummy=dummy");
}

function sendPostRegisterRequest(){
    const name = document.getElementById("fullname").value;
    const usernameInput = document.getElementById("usernameReg").value;
    const passwordInput = document.getElementById("passwordReg").value;
    const request = createAjaxRequest();
    request.onreadystatechange = function () {
        if(4 === this.readyState){
            if(200 === this.status){
                alert(this.responseText)
                hideVBlockerAndLogin();
                unhideAvatar();
                loginAfterRegister();
            }else{
                alert(this.status + ":" + this.responseText);
            }
        }
    }
    request.open("POST",/register/,true);
    request.setRequestHeader("Content-Type","application/x-www-form-urlencoded");
    console.log("" + usernameInput + " " + passwordInput);
    request.send("usernameInput="+usernameInput+"&"+"passwordInput="+passwordInput+"&"+"nameInput="+name);
}

function sendPostFavoriteRequest(){
    const request = createAjaxRequest();
    request.onreadystatechange = function () {
        if(4 === this.readyState){
            if(200 === this.status){
                alert(this.responseText);
            }else{
                alert(this.status + ":" + this.responseText);
            }
            console.log(this);
        }
    }
    request.open("POST",/addToFavorites/,true);
    request.setRequestHeader("Content-Type","application/x-www-form-urlencoded");
    request.send("video="+JSON.stringify(video));
}



function sendGetClickedVideos(){
    const request = createAjaxRequest();
    request.onreadystatechange = function () {
        if(4 === this.readyState){
            if(200 === this.status){
                alert(this.responseText);
            }else{
                alert(this.status + ":" + this.responseText);
            }
        }
    }

    request.open("GET","/clickVideo" +"?videoTitle="+"\"Plan B\" für Bayern",true);
    request.send();
}

function hideVBlockerAndLogin() {
    var vblocker = document.getElementById("v_blocker");
    var loginscreen = document.getElementById("Login_Screen");
    loginscreen.style.visibility = "hidden";
    vblocker.style.visibility = "hidden";
}

function unhideVBlockerAndLogin() {
    var vblocker = document.getElementById("v_blocker");
    var loginscreen = document.getElementById("Login_Screen");
    vblocker.style.visibility = "visible";
    loginscreen.style.visibility = "visible";
}
function hideAvatar() {
    var avatar = document.getElementById("Dropdown");
    avatar.style.visibility = "hidden";

}
function unhideAvatar() {
    var avatar = document.getElementById("Dropdown");
    avatar.style.visibility = "visible";
}

function openProfil(){
    window.location.href = "/Profil.html";
}

function loginAfterRegister() {
   let userLogin = document.getElementById("usernameLogin");
   let userPass = document.getElementById("passwordLogin");
   userLogin.value = document.getElementById("usernameReg").value;
   userPass.value = document.getElementById("passwordReg").value;
   sendPostLoginRequest();
}

function openSearchPage() {
    let searchInp = document.getElementById("searchInput").value;
    sessionStorage.setItem('search',searchInp);
    console.log("OpenSearchPage: " + sessionStorage.getItem("search"));
    window.location.href = "/searchresults.html";
}

function fillSearchPage() {
    let videoEx = new Video("","","","","","","","");
    let videoResultsC = document.getElementById("videoResultContainer");
    let vidDiv = document.getElementById("videoResults");
    vidDiv.innerHTML = "";
    for(videoEx of searchResultsJSON) {
        let h2 = document.createElement("div");
        let h5 = document.createElement("h5");
        h5.innerText = videoEx.title;
        h2.appendChild(h5);
        vidDiv.appendChild(h2)
    }
}

function initChannelPage() {
    pageFlag = "channel";
    sessionStorage.setItem("pageflag",pageFlag);
    sendGetVideos();
    searchOnEnter();
    console.log(pageFlag);
}

function initMainPage() {
    pageFlag = "main";
    sessionStorage.setItem("pageflag",pageFlag);
    sendPostCookieAuthRequest();
    searchOnEnter();
    console.log(pageFlag);
}

function initSearchPage() {
    searchOnEnter();
    let whichPage = sessionStorage.getItem("pageflag");
    console.log("Suche wird aufgerufen: " + whichPage);
    if(!(whichPage.localeCompare("main"))) {
        sendGetSearchRequest();
    }else {
        console.log("Suche Channel");
        sendGetSearchRequestChannel();
    }
}



function openTab(evt, tabName)  {
    // Declare all variables
    var i, tabcontent, tablinks;

    // Get all elements with class="tabcontent" and hide them
    tabcontent = document.getElementsByClassName("tabcontent");
    for (i = 0; i < tabcontent.length; i++) {
        tabcontent[i].style.display = "none";
    }

    // Get all elements with class="tablinks" and remove the class "active"
    tablinks = document.getElementsByClassName("tablinks");
    for (i = 0; i < tablinks.length; i++) {
        tablinks[i].className = tablinks[i].className.replace(" active", "");
    }

    // Show the current tab, and add an "active" class to the button that opened the tab
    document.getElementById(tabName).style.display = "block";
    evt.currentTarget.className += " active";
}

function searchOnEnter() {
    const inputSearch = document.getElementById("searchInput");
    inputSearch.addEventListener("keyup", function (event) {
        if (event.key === "Enter") {
            event.preventDefault();
            document.getElementById("searchIcon").click();
        }
    })
}

function decideSearch() {

}
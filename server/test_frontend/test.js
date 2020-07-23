"use strict"

class User{
    constructor(id,name,username,favoriteVideos) {
        this.id = id;
        this.username = username;
        this.name = name;
        this.favoriteVideos= favoriteVideos;
    }

}

const video =  {
    "channel": "ARD",
    "title": "2 Mann f\u00fcr alle G\u00e4nge - Roastbeef mit Senfsaatsauce",
    "show":  "2 Mann f\u00fcr alle G\u00e4nge",
    "releaseDate": "28.01.2016",
    "duration": "00:30:00",
    "link": "http://mediastorage01.sr-online.de/Video/FS/ZMANN/2mann_20160123_180701_L.mp4",
    "pageLink": "http://www.ardmediathek.de/tv/2-Mann-f%C3%BCr-alle-G%C3%A4nge/2-Mann-f%C3%BCr-alle-G%C3%A4nge-Roastbeef-mit-Se/SR-Fernsehen/Video?bcastId=8638714&documentId=32938818",
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

function sendPostSaveProfilePicture(){
    const request = createAjaxRequest();
    const profilePicture = document.getElementById("ppUpload").value;
    const formData = new FormData();
    formData.append("profilepicture",profilePicture)
    request.onreadystatechange = function () {
        if(4 === this.readyState){
            if(200 === this.status){
                alert(this.responseText)
            }else{
                alert(this.status + ":" + this.responseText);
            }
        }
    }
    request.open("POST","/setProfilePicture/",true);
    request.send(formData);
}

function sendGetLoadProfilePicture(){
    const request = createAjaxRequest();
    request.onreadystatechange = function () {
        if(4 === this.readyState){
            if(200 === this.status){
                document.getElementById("pp").value = this.responseText
                console.log(this.responseText)
            }else{
                alert(this.status + ":" + this.responseText);
            }
        }
    }
    request.open("GET","/getProfilePicture/",true);
    request.send();
}

function sendPostCookieAuthRequest(){
    const request = createAjaxRequest();
    request.onreadystatechange = function () {
        if(4 === this.readyState){
            if(200 === this.status){
                alert(this.responseText)
            }else{
                alert(this.status + ":" + this.responseText);
            }
        }
    }
    request.open("POST","/cookieAuth/",true);
    request.setRequestHeader("Content-Type","application/x-www-form-urlencoded");
    request.send("dummy=dummy");
}

function sendPostFetchFavoritesRequest(){
    const request = createAjaxRequest();
    request.onreadystatechange = function () {
        if(4 === this.readyState){
            if(200 === this.status){
                alert(this.responseText)
            }else{
                alert(this.status + ":" + this.responseText);
            }
        }
    }
    request.open("POST","/getFavorites/",true);
    request.setRequestHeader("Content-Type","application/x-www-form-urlencoded");
    request.send("dummy=dummy");
}

function sendGetSearchRequest(){
    const request = createAjaxRequest();
    const searchString = document.getElementById("searchInput").value;
    const channelString = document.getElementById("channelInput").value;
    request.onreadystatechange = function () {
        if(4 === this.readyState){
            if(200 === this.status){
                alert(this.responseText);
            }else{
                alert(this.status + ":" + this.responseText);
            }
        }
    }

    request.open("GET","/search" +"?search="+searchString + "&channel="+channelString,true);
    request.send();
}

function sendPostLoginRequest(){
    const usernameInput = document.getElementById("usernameInput").value;
    const passwordInput = document.getElementById("passwordInput").value;
    const request = createAjaxRequest();
    request.onreadystatechange = function () {
        if(4 === this.readyState){
            if(200 === this.status){
                alert(this.responseText)
                let user = new User("","","",null);
                user = JSON.parse(this.responseText);
                console.log(user);
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
    const name = document.getElementById("nameInput2").value;
    const usernameInput = document.getElementById("usernameInput2").value;
    const passwordInput = document.getElementById("passwordInput2").value;
    const request = createAjaxRequest();
    request.onreadystatechange = function () {
        if(4 === this.readyState){
            if(200 === this.status){
                alert(this.responseText)

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

function sendPostRemoveFavoriteRequest(){
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
    request.open("POST",/removeFromFavorites/,true);
    request.setRequestHeader("Content-Type","application/x-www-form-urlencoded");
    request.send("video="+encodeURIComponent(JSON.stringify(video)));
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
    console.log()
    request.send("video="+encodeURIComponent(JSON.stringify(video)));
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

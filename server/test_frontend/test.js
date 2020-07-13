"use strict"

class User{
    constructor(id,name,username,favoriteVideos) {
        this.id = id;
        this.username = username;
        this.name = name;
        this.favoriteVideos= favoriteVideos;
    }

}

let username;
let password;

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
                username = usernameInput;
                password = passwordInput;
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

function sendPostFavoriteRequest(){
    const usernameInput = username;
    const passwordInput = password;
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
    request.open("POST",/addToFavorites/,true);
    request.setRequestHeader("Content-Type","application/x-www-form-urlencoded");
    console.log("" + usernameInput + " " + passwordInput);
    request.send("usernameInput="+usernameInput+"&"+"passwordInput="+passwordInput+"&"+"video="+JSON.stringify(video));
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

    request.open("GET","/clickVideo",true);
    request.send("?videoTitle="+"\"Plan B\" für Bayern");
}

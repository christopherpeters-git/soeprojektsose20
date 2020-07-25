function createAjaxRequest(){
    let request;
    if(window.XMLHttpRequest){
        request = new XMLHttpRequest();
    }else{
        request = new ActiveXObject("Microsoft.XMLHTTP");
    }
    return request;
}

function sendPostCookieAuthRequest(callbackFunction,async=true){
    const request = createAjaxRequest();
    request.onreadystatechange = function () {
        if(4 === this.readyState){
           callbackFunction(this);
        }
    }
    request.open("POST","/cookieAuth/",async);
    request.setRequestHeader("Content-Type","application/x-www-form-urlencoded");
    request.send("dummy=dummy");
}

function sendPostLoginRequest(callbackFunction){
    const usernameInput = document.getElementById("usernameLogin").value;
    const passwordInput = document.getElementById("passwordLogin").value;
    const request = createAjaxRequest();
    request.onreadystatechange = function () {
        if(4 === this.readyState){
            callbackFunction(this);
        }
    }
    request.open("POST",/login/,true);
    request.setRequestHeader("Content-Type","application/x-www-form-urlencoded");
    request.send("usernameInput="+usernameInput+"&"+"passwordInput="+passwordInput);
}

function sendPostLogoutRequest(callbackFunction){
    const request = createAjaxRequest();
    request.onreadystatechange = function () {
        if(4 === this.readyState){
            callbackFunction(this);
        }
    }
    request.open("POST",/logout/,true);
    request.setRequestHeader("Content-Type","application/x-www-form-urlencoded");
    request.send("dummy=dummy");
}

function sendPostRegisterRequest(callbackFunction){
    const name = document.getElementById("fullname").value;
    const usernameInput = document.getElementById("usernameReg").value;
    const passwordInput = document.getElementById("passwordReg").value;
    const request = createAjaxRequest();
    request.onreadystatechange = function () {
        if(4 === this.readyState){
            callbackFunction(this);
        }
    }
    request.open("POST",/register/,true);
    request.setRequestHeader("Content-Type","application/x-www-form-urlencoded");
    request.send("usernameInput="+usernameInput+"&"+"passwordInput="+passwordInput+"&"+"nameInput="+name);
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

function sendGetVideos(callbackFunction,async=true) {
    const request = createAjaxRequest();
    request.onreadystatechange = function () {
        if (4 === this.readyState) {
            callbackFunction(this);
        }
    }
    request.open("GET", "/getVideoByChannel" + "?channel=" + sessionStorage.getItem('channel'), async);
    request.send();
}

function sendGetSearchRequest(callBackFunction,async=true){
    const request = createAjaxRequest();
    const incomingString = JSON.parse(sessionStorage.getItem("searchString"));
    console.log(incomingString);
    let channelString = incomingString [0];
    let searchString =incomingString[1];
    request.onreadystatechange = function () {
        if(4 === this.readyState){
            callBackFunction(this);
        }
    }
    request.open("GET","/search" +"?search="+searchString + "&channel="+channelString,async);
    request.send();
}

function sendGetClickedVideos(video,async=true){
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
    request.open("GET","/clickVideo" +"?videoTitle="+video[0].title,async);
    request.send();
    return clickNumber;
}

function sendPostFavoriteRequest(video){
    const request = createAjaxRequest();
    request.onreadystatechange = function () {
        if(4 === this.readyState){
            if(200 === this.status){
                alertSetterFunction("rgba(39,255,0,0.75)",this.responseText,1500);
            }else{
                alertSetterFunction("rgba(255,0,30,0.75)",this.responseText,1500);
            }
        }
    }
    request.open("POST",/addToFavorites/,true);
    request.setRequestHeader("Content-Type","application/x-www-form-urlencoded");
    request.send("video="+video);
}
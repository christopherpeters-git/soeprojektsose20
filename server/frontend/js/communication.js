const illegalStrings = ["/",".","<",">"];

function createAjaxRequest(){
    let request;
    if(window.XMLHttpRequest){
        request = new XMLHttpRequest();
    }else{
        request = new ActiveXObject("Microsoft.XMLHTTP");
    }
    return request;
}

//allows user to choose picture from files and sends picture to server, callback reloads site on success with new profile picture
function sendPostSaveProfilePicture(callbackFunction){
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
                callbackFunction(this);
            }else{
                console.log(this.status + ":" + this.responseText);
            }
        }
    }
    request.open("POST","/setProfilePicture/",true);
    request.send(formData);
}
//checks with server for existing cookie, callback on success hides login and makes site useable
function sendGetCookieAuthRequest(callbackFunction, async=true){
    const request = createAjaxRequest();
    request.onreadystatechange = function () {
        if(4 === this.readyState){
           callbackFunction(this);
        }
    }
    request.open("GET","/cookieAuth/",async);
    request.send();
}
//opens request for user_favorites from server, callback creates favorites in html
function sendGetFetchFavoritesRequest(callbackFunction, async = true){
    const request = createAjaxRequest();
    request.onreadystatechange = function () {
        if(4 === this.readyState){
            if(200 == this.status){
                callbackFunction(this);
            }else{
                console.log(this.status + ":" + this.responseText);
            }
        }
    }
    request.open("GET","/getFavorites/",async);
    request.send();
}
//opens loginrequest for server, callback hides login on success
function sendPostLoginRequest(callbackFunction){
    const usernameInput = document.getElementById("usernameLogin").value;
    const passwordInput = document.getElementById("passwordLogin").value;
    const stringArray = [usernameInput,passwordInput];
    alert("Wir benutzen einen Cookie, um die automatische Authentifizierung zu ermöglichen.");
    if(areStringsIllegal(stringArray)){
        console.alert("Illegale strings");
        return;
    }
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
//user is logged out, callback opens login on success
function sendGetLogoutRequest(callbackFunction){
    const request = createAjaxRequest();
    request.onreadystatechange = function () {
        if(4 === this.readyState){
            callbackFunction(this);
        }
    }
    request.open("GET",/logout/,true);
    request.send();
}
//opens registerrequest, callback logs user in and hides login
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
    if(areStringsIllegal(incomingString)) {
        alert("Illegaler String!");
        return;
    }
    console.log(incomingString);
    let channelString = incomingString [0];
    if(channelString==null){
        channelString="none";
    }
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
//visible feedback for user for adding/deleting favorites
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

function areStringsIllegal(EnteredStringArray) {
    for (let i = 0; i<illegalStrings.length;i++) {
        if (EnteredStringArray[1].includes(illegalStrings[i])) {
            return true;
        }
    }
    return false;
}
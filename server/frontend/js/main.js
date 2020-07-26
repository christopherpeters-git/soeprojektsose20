"use strict"

class User{
    constructor(id,name,username,favoriteVideos) {
        this.id = id;
        this.username = username;
        this.name = name;
        this.favoriteVideos= favoriteVideos;
    }

}


function callBackFunctionCookieAuthRequest(status){
    if(status.status===200){
        hideVBlockerAndLogin();
        unhideAvatar();
        const currentSite = window.location;
        console.log(currentSite);
    }else{
        console.log(status.status + ":" + status.responseText);
        document.getElementById("Login_Screen").style.visibility="visible";
    }
}

function callBackFunctionLogin(status){
    if(200 === status.status){
        hideVBlockerAndLogin();
        unhideAvatar();
    }else{
        alert(status.status + ":" + status.responseText);
    }
}

function callBackFunctionLogout(status) {
    if(200 === status.status){
        unhideVBlockerAndLogin();
        hideAvatar();
    }else{
        alert(status.status + ":" + status.responseText);
    }
}

function callBackFunctionRegister(status) {
    if(200 === status.status){
        hideVBlockerAndLogin();
        unhideAvatar();
        loginAfterRegister();
    }else{
        alert(status.status + ":" + status.responseText);
    }
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

function openProfil() {
    window.location.href="/profile.html";
}



function loginAfterRegister() {
   let userLogin = document.getElementById("usernameLogin");
   let userPass = document.getElementById("passwordLogin");
   userLogin.value = document.getElementById("usernameReg").value;
   userPass.value = document.getElementById("passwordReg").value;
   sendPostLoginRequest(callBackFunctionLogin);
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

function setSearchtext() {
    const searchValue = document.getElementById("searchInput").value;
    if(searchValue==="") return;
    let searchString = JSON.stringify(["none",searchValue]);
    sessionStorage.setItem("searchString",searchString);
    console.log(sessionStorage.getItem("searchString"));
    window.location.href = "/searchResults.html";
}
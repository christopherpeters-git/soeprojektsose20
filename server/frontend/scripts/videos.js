const getJsonVideoTargetUrl = "localhost:80/getVideos/";
let videosJson;


class video {
    constructor(channel,title,show,releaseDate,duration,link,pageLink,fileName) {
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

function init() {

    //intialisiere Rest der Seite
    searchOnEnter();
}


function searchOnEnter() {
    const inputSearch = document.getElementById("searchInput");
    inputSearch.addEventListener("keyup", function (event) {
        if (event.key === "Enter") {
            event.preventDefault();
            document.getElementById("iconSearchbar").click();
        }
    })
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

function requestVideos() {
    var request = createAjaxRequest();
    request.onreadystatechange = function() {
        if (4 === this.readyState && 200 === this.status) {
            videosJson = JSON.parse(this.responseText);
            init();
        }

    }
    request.open("GET",getJsonVideoTargetUrl,true);
    request.send();
}

function searchThroughDatabase() {
    console.log("Starte Suche...");
}
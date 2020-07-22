class Video {
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
    checkAttributes(searchArray,results) {
        let titleLower = this.title.toLowerCase();
        let showLower = this.show.toLowerCase();

        for(let i = 0;i<searchArray.length;i++) {
            let currentSubstring = searchArray[i];
            if(titleLower.includes(currentSubstring) || showLower.includes(currentSubstring)) {
                for(let j = 0; j<results.length;j++) {
                    if(this === results[j]) {
                        return;
                    }
                }
                results.push(this);
            }
        }
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
            document.getElementById("searchIcon").click();
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
    let searchInput = document.getElementById("searchInput").value;
    if(searchInput === "") {
        console.log("Leerer Eintrag, keine Suche durchgeführt");
        return;
    }
    searchInput.toLowerCase();
    if(checkIfChannel(searchInput)) {
        return;
    }
    let searchArray = searchInput.split(" ");
    let currentVideo = new Video("","","","","","","","");
    let results = [];
    for(currentVideo of videosJson) {
        currentVideo.checkAttributes(searchArray,results);
    }
}

//******************************HelperFunctionsSearchbar**********************************************

function checkIfChannel(searchInput) {
    //TODO
    return false;
}


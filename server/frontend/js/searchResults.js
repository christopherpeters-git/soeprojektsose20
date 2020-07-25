let searchJson;


function initSearchResults() {
    sessionStorage.setItem("pageFlag","1");
    checkFlag();
}

function setSearchResult(){

}

function setSearchTextWithResults(){
    const searchValue = document.getElementById("searchInput").value;
    if(searchValue==="") return;
    let searchString;
    searchString= JSON.stringify([sessionStorage.getItem("channel"), searchValue]);
    sessionStorage.setItem("searchString",searchString);
    initSearchResults();
}
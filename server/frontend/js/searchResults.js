
function initSearchResults() {
    sessionStorage.setItem("pageFlag","1");
    checkFlag();
}


function setSearchTextWithResults(){
    currentPage=1;
    const inputValue = document.getElementById("inputButton");
    inputValue.value = currentPage;
    const searchValue = document.getElementById("searchInput").value;
    if(searchValue==="") return;
    let searchString;
    console.log(sessionStorage.getItem("channel"));
    if(sessionStorage.getItem("channel")===null){
        searchString= JSON.stringify(["none", searchValue]);
    }else{
        searchString= JSON.stringify([sessionStorage.getItem("channel"), searchValue]);
    }
    sessionStorage.setItem("searchString",searchString);
    initSearchResults();
}
//Function to initialize the result of the search.
function initSearchResults() {
    sessionStorage.setItem("pageFlag","1");
    checkFlag();
}
// Function to search on the SearchResult page.
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

function addEventListenerOnSearchPage() {
    const searchInput = document.getElementById("searchInput");
    searchInput.addEventListener("keydown", searchAgainOnEnter, false);

}
//Function to start the search with enter.
function searchAgainOnEnter(event) {
    if(event.key ==="Enter") {
        setSearchTextWithResults();
    }
}
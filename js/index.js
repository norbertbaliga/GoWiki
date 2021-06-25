function createPage() {
    var title = document.getElementById('title');
    if(title == null || title.value == '') {
        alert('Please give a page title!');
        return false;
    }
    location.href = '/edit/' + title.value;
}

function fillSearchInput() {
    var queryString = window.location.search;
    var urlParams = new URLSearchParams(queryString);
    var searchParam = urlParams.get('q');

    if(searchParam == null || searchParam == '') {
        return 
    }

    var searchInput = document.getElementById('search');
    if(searchInput != null) {
        searchInput.value = searchParam;
    }
}
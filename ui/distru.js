function searchThis() {
	if (event.keyCode == 13) {
		window.location += 'search/'+document.getElementById('search').value;
	}
}

function isEnter(e) {
    e = e || window.event || {};
    var charCode = e.charCode || e.keyCode || e.which;
        if (charCode == 13) {
			window.location += 'search/'+document.getElementById('search').value;
        }

 }
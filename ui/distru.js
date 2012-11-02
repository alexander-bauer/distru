function searchThis() {
	if (event.keyCode == 13) {
		window.location += 'search/'+document.getElementById('search').value;
	}
}
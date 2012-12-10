//  The function searchThis() works with Safari and Chrome
//  but not with FireFox and Opera. Instead of using a 
//  complete HTML5 form, we use an onkeydown function that
//  calls the javascript function to redirect us to search/x
//  where x is the search term(s)
function searchThis() {
	if (event.keyCode == 13) {
		window.location = '/search/'+document.getElementById('search').value.replace(/%/g, "%25");
	}
}

//  The function isEnter(e) works with FireFox and Opera
//  while the other function works with Safari and Chrome.
//  Instead of using a complete HTML5 form, we use an onkeypress
//  function that calls the javascript function to redirect
//  us to search/x where x is the search term(s)
function isEnter(e) {
    e = e || window.event || {};
    var charCode = e.charCode || e.keyCode || e.which;
        if (charCode == 13) {
			window.location = '/search/'+document.getElementById('search').value.replace(/%/g, "%25");
        }
}


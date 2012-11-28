//
//	parse.js
//


// Parse that shit on load.
window.onload=parse;
var url = window.location.pathname
var term = url.substring(url.indexOf(":9048/search/")+9);

// The function parse() decides what to do with the term.
// There are two things that can be done. First, if it
// is a mathematical equation, then do the math and add a
// div on the search page with the answer for the user.
// Second, if there is only one word, and it is a real word,
// then define it. 
function parse() {
	doMath(term);
	
}

// The function doMath(term) checks to see if the term
// is a math equation or not. First, it should parse
// words and change them into the corresponding symbols. 
// Ex: plus, minus, times ... +, -, *
// Then it should figure out if it is a real equation or not.
// If it is a real equation, do the math, and if it is
// not a real equation, then don't do the math.
function doMath(term) {
	try {
		
		// put it in lowercase so we don't have to deal
		// with this shit being case sensitive and whatnot
		term = term.toLowerCase()
		term = decodeURIComponent(term);
		
		var operators = {"plus" : "+", "and" : "+", "minus" : "-", "times" : "*", "x" : "*",
						 "over" : "/", "divide" : "/", "divided by" : "/", "mod" : "%", "modulus" : "%"};
		
		for (var val in operators) {
			term = term.replace(new RegExp(val, "g"), operators[val]);
		}
				
		var value = eval(term);
		if (!isNaN(value)) {
			// display some html
			 document.getElementById("blank").innerHTML = "<center><div class='calculate' onmouseover='unhideBubble();' onmouseout='hideBubble();'>" + term + " = <strong>" + value + "</strong></div><div class='bubble hidden'><strong>What's this?</strong><br/>What you serached seemed to us like it was math, so we did the math for you!</div></center>";
		}
	}
	catch (e) {}
}

function unhideBubble() {
	document.getElementsByClassName("bubble").item(0).style.opacity = "1";
}

function hideBubble() {
	document.getElementsByClassName("bubble").item(0).style.opacity = "0";
}
//
//	parse.js
//


// Parse that shit on load.
window.onload=parse;
var url = window.location.pathname
var term = url.substring(url.indexOf(":9048/search/")+9);

// The function parse() decides what to do with the term.
// There are a few things that can be done. First, if it
// is a mathematical equation, then do the math and add a
// div on the search page with the answer for the user.
function parse() {
	parseMath(term);
}

// The function parseMath(term) checks to see if the term
// is a math equation or not. First, it should parse
// words and change them into the corresponding symbols. 
// Ex: plus, minus, times, ... TO ... +, -, *, ...
// Then it should figure out if it is a real equation or not.
// If it is a real equation, do the math, and if it is
// not a real equation, then don't do the math.
function parseMath(term) {
	try {
		
		// put it in lowercase so we don't have to deal
		// with this shit being case sensitive and whatnot
		term = term.toLowerCase()
		term = decodeURIComponent(term);
		
		var operators = {"plus" : "+", "and" : "+", "minus" : "-", "times" : "*",
						 "over" : "/", "divide" : "/", "mod" : "%", "modulus" : "%"};
		
		for (var val in operators) {
			term = term.replace(new RegExp(val, "g"), operators[val]);
		}
		
		
		if (isHex(term)){
			doMath(16, "0x", term);
		} //close hex
		
		else if (isOctal(term)) {
			doMath(8, "0o", term);
		} //close octal
		
		else if (isBinary(term)) {
			doMath(2, "0b", term);
		} //close binary
		
		else
		{
		var value = eval(term);
		
			if (!isNaN(value)) {
				// display some html
				 document.getElementById("blank").innerHTML = "<center><div class='calculate' onmouseover='unhideBubble();' onmouseout='hideBubble();'>" + term + " = <strong>" + value + "</strong></div><div class='bubble'><strong>What's this?</strong><br/>What you serached seemed to us like it was math, so we did the math for you!</div></center>";
			}
		} //close else
	}
	catch (e) {}
}

// The function doMath(base, type, newterm) does math!
// More specifically, it does hex, octal, and binary 
// mathematical equations. 
function doMath(base, type, newterm) {
	var original = newterm;
	term = newterm;
	term = term.replace("+", " + ");
	term = term.replace("-", " - ");
	term = term.replace("*", " * ");
	term = term.replace("/", " / ");
	term = term.replace("%", " % ");
	
	term = term.replace(/^\s+|\s+$/g,'').replace(/\s+/g,' ');
	
	var array = term.split(" ");
	
	
	for (var i = 0; i < array.length; i++) {
		if (!(array[i] == "+" || array[i] == "-" || array[i] == "*" || array[i] == "/" || array[i] == "%")) {
			if (array[i].length > 2 && array[i].substring(0,2) == "0o" || array[i].substring(0,2) == "0b") {
				var temp = array[i].substring(2);
				array[i] = parseInt(temp, base);
			}
			else array[i] = parseInt(array[i], base);
		}
	}
	term = array.join(' ');			
	var value = eval(term);
	value = value.toString(base);
	
	document.getElementById("blank").innerHTML = "<center><div class='calculate' onmouseover='unhideBubble();' onmouseout='hideBubble();'>" + original + " = <strong>" + type + value + "</strong></div><div class='bubble'><strong>What's this?</strong><br/>What you serached seemed to us like it was math, so we did the math for you!</div></center>";
}

// The function isHex() determines if 
// it is a hex number or not.
function isHex() {
	if (term.indexOf("0x") !== -1)
		return true;
	else return false;
}

// The function isOctal() determines if
// it is an octal number or not
function isOctal() {
	if (term.indexOf("0o") !== -1)
		return true;
	return false;
}

// The function isBinary() determines if
// it is a binary number or not
function isBinary() {
	if (term.indexOf("0b") !== -1)
		return true;
	return false;
}

// The function unhideBubble() unhides the bubble!
function unhideBubble() {
	document.getElementsByClassName("bubble").item(0).style.opacity = "1";
}

// The function hideBubble() hides the bubble!
function hideBubble() {
	document.getElementsByClassName("bubble").item(0).style.opacity = "0";
}
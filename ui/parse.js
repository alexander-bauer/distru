// Parse that shit on load.
window.onload=parse;
var url = window.location.pathname;
var term = url.substring(url.indexOf(":9048/search/")+9);
var base;
var original;

// The function parse() decides what to do with the term.
// There are a few things that can be done. First, if it
// is a mathematical equation, then do the math and add a
// div on the search page with the answer for the user.
function parse() {
	parseMath(term);
} // Close parse

// The function parseMath(term) checks to see if the term
// is a math equation or not. First, it should parse
// words and change them into the corresponding symbols. 
// Ex: plus, minus, times, ... TO ... +, -, *, ...
// Then it should figure out if it is a real equation or not.
// If it is a real equation, do the math, and if it is
// not a real equation, then don't do the math.
function parseMath(term) {
	try {
		
		// TODO
		// Check and see if it's really math.
		// If it's really math, then fine, do the rest
		// of it, but if not, then it shouldn't try
		// and do math because that's just more calls
		// that we don't need.
				
		term = fixOperators(term);
		base = checkBase(term);
		term = removeTermBase(term);
		original = term;
		term = fixSpaces(term);
		term = convertToDecimal(term);
		term = evaluate(term);
		if (!isNaN(term)) {
			term = convertToBase(term, base);	
			putOnPage(term);
		}
	}
	catch (e) {
		// Fuck errors
	}
} // Close parseMath

// The function fixOperators(term) changes all word
// operators with real operators so when doing math
// everything is much, much easier. Plus (no pun intended),
// it gives people more ways to do math
function fixOperators(term) {
	// I'm not dealing with case-sensitive shit
	term = term.toLowerCase();
	// This gets rid of the %20s and shit
	term = decodeURIComponent(term);
	
	var operators = {"plus" : "+", "and" : "+", "minus" : "-", "times" : "*",
					 "over" : "/", "divide" : "/", "mod" : "%", "modulus" : "%"};
	
	for (var val in operators) term = term.replace(new RegExp(val, "g"), operators[val]);
	
	return term;
} // Close fixOperators

// The function checkBase(term) checks to see what
// base it should put the final result in at the end.
function checkBase(term) {
	if (term.indexOf("in binary") !== -1) return 2;
	else if (term.indexOf("in octal") !== -1) return 8;
	else if (term.indexOf("in decimal") !== -1) return 10;
	else if (term.indexOf("in hex") !== -1 || term.indexOf("in hexadecimal") !== -1) return 16;
	
	var binary = term.indexOf("0b");
	var octal = term.indexOf("0o");
	var hex = term.indexOf("0x");

	if (binary !== -1) {
		if (octal !== -1) {
			if (hex !== -1) {
				if (binary < octal && binary < hex) return 2;
				else if (octal < binary && octal < hex) return 8;
				else return 16;
			}
			if (binary < octal) return 2;
			else return 8;
		}
		return 2;
	}
	else if (octal !== -1) {
		if (hex !== -1) {
			if (octal < hex) return 8;
			else return 16;
		}
		return 8;
	}
	else if (hex !== -1) return 16;
	else return 10;
} // Close checkBase

// The function fixSpaces(term) fixes all
// of the damn spacing errors that people
// tend to do because they're assholes.
function fixSpaces(term) {
	term = term.replace("+", " + ");
	term = term.replace("-", " - ");
	term = term.replace("*", " * ");
	term = term.replace("/", " / ");
	term = term.replace(/\^/g, " ^ ");
	term = term.replace(/\!/g, " ! ");
	term = term.replace(/%/g, " % ");
	term = term.replace(/\)/g, " ");
	term = term.replace(/\(/g, " ");
	
	term = term.replace(/^\s+|\s+$/g,'').replace(/\s+/g,' ');
	
	return term;	 
} // Close fixSpaces

// The function convertToDecimal(term) converts
// each part of the equation from whatever base
// it was at to base 10
function convertToDecimal(term) {
	var array = term.split(" ");
	
	for (var i = 0; i < array.length; i++) {
		if (!(array[i] == "+" || array[i] == "*" || array[i] == "/" || array[i] == "%")) {
			if (array[i].length > 2) {
				var temp = array[i].substring(2);
				if (array[i].substring(0,2) == "0x") array[i] = parseInt(temp, 16);
				else if (array[i].substring(0,2) == "0o") array[i] = parseInt(temp, 8);
				else if (array[i].substring(0,2) == "0b") array[i] = parseInt(temp, 2);
			}
		}
	}
		
	term = array.join(' ');
	return term;
} // Close convertToDecimal

// The function evaluate(term) evaluates
// the function by doing different types
// of maths
function evaluate(term) {
	
	for (var i = 0; i < term.length; i++) {
		if (term.indexOf("!") !== -1) term = factorial(term);
		else if (term.indexOf("^") !== -1) term = power(term);
		else if (term.indexOf("abs") !== -1) term = absv(term);
		else if (term.indexOf("sqrt") !== -1) term = square(term);
	}
	
	term = eval(term);
	
	return term;
} // Close evaluate

// The function putOnPage(term) puts the
// term, evaluated and all, on the page
// for all to see!
function putOnPage(term) {
	document.getElementById("blank").innerHTML = "<center><div class='calculate' onmouseover='unhideBubble();' onmouseout='hideBubble();'>" + original + " = <strong>" + term + "</strong></div><div class='bubble'><strong>What's this?</strong><br/>What you serached seemed to us like it was math, so we did the math for you!</div></center>";
} // Close putOnPage

// The function power(term) checks the entire
// equation for exponents, and does all of the
// exponents possible, then returns the term
function power(term) {
	var array = term.split(" ");
	for (var i = 0; i < array.length; i++) {
		if (array[i].indexOf("^") !== -1) {
			var result = Math.pow(array[i-1], array[i+1]);
			array[i-1] = result;
			array[i] = "";
			array[i+1] = "";
		}
	}
	term = array.join(" ");
	term = fixSpaces(term);
	if (term.indexOf("^") !== -1)
		return power(term);
	return term;
} // Close power

// The function square(term) checks the entire
// equation for square roots, and does all of the
// square roots possible, then returns the term
function square(term) {
	var array = term.split(" ");
	for (var i = array.length-1; i >= 0; i--) {
		if (array[i].indexOf("sqrt") !== -1) {
			array[i] = Math.sqrt(array[i+1]);
			array[i+1] = "";
		}
	}
	
	term = array.join(" ");
	term = fixSpaces(term);
	return term;
} // Close square

// The function absv(term) checks the entire
// equation for absolute values, and does all of the
// absolute values possible, then returns the term
function absv(term) {
	var array = term.split(" ");
	for (var i = array.length-1; i >= 0; i--) {
		if (array[i].indexOf("abs") !== -1) {
			array[i] = Math.abs(array[i+2]);
			array[i+1] = "";
			array[i+2] = "";
		}
	}
	
	term = array.join(" ");
	term = fixSpaces(term);
	return term;
} // Close absv

// The function power(term) checks the entire
// equation for factorials, and does all of the
// factorials possible, then returns the term
function factorial(term) {
	var array = term.split(" ");
	for (var i = 0; i < array.length; i++) {
		if (array[i].indexOf("!") !== -1) {
			var result = 1;
			for (var j = array[i-1]; j > 0; j--) {
				result *= j;
			}
			array[i] = "";
			array[i-1] = result;
		}
	}
	term = array.join(" ");
	term = fixSpaces(term);
	return term;
} // Close factorial

// The function convertToBase(term, base) converts
// each part of the equation from base 10 to
// whatever base it needs to be in
function convertToBase(term, base) {
	term = term.toString(base);
	if (term.indexOf("-") == 0) {
		term = term.replace("-", "");
		if (base == 2) term = "0b" + term;
		else if (base == 8) term = "0o" + term;
		else if (base == 16) term = "0x" + term;
		term = "-" + term;
	}
	else {
		if (base == 2) term = "0b" + term;
		else if (base == 8) term = "0o" + term;
		else if (base == 16) term = "0x" + term;
	}
	return term;
} // Close convertToBase

// The function unhideBubble() unhides the bubble!
function unhideBubble() {
	document.getElementsByClassName("bubble").item(0).style.opacity = "1";
} // Close unhideBubble

// The function hideBubble() hides the bubble!
function hideBubble() {
	document.getElementsByClassName("bubble").item(0).style.opacity = "0";
} // Close hideBubble

// The function removeTermBase(term) removes
// any instance of "to binary", "in binary", ...
// all the way up until hex.
function removeTermBase(term) {
	if (term.indexOf("in binary") !== -1) {
		return term.replace("in binary","");
	}
	else if (term.indexOf("to binary") !== -1) {
		return term.replace("to binary","");
	}
	else if (term.indexOf("in octal") !== -1) {
		return term.replace("in octal","");
	}
	else if (term.indexOf("to octal") !== -1) {
		return term.replace("to octal","");
	}
	else if (term.indexOf("in decimal") !== -1) {
		return term.replace("in decimal","");
	}
	else if (term.indexOf("to decimal") !== -1) {
		return term.replace("to decimal","");
	}
	else if (term.indexOf("in hex") !== -1) {
		return term.replace("in hex","");
	}
	else if (term.indexOf("in hexadecimal") !== -1) {
		return term.replace("in hexadecimal","");
	}
	else if (term.indexOf("to hex") !== -1) {
		return term.replace("to hex","");
	}
	else if (term.indexOf("to hexadecimal") !== -1) {
		return term.replace("to hexadecimal","");
	}
	else return term;
} // Close removeTermBase
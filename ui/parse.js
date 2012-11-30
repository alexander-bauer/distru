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
		
		//Fix spacing errors and put each thing in an array
		term = term.replace(/^\s+|\s+$/g,'').replace(/\s+/g,' ');
		var original = term;
		var array = term.split(" ");
		
		//Check to see if there's some more math we can do.
		for (var i = 0; i < array.length; i++) {
			if (array[i].indexOf("^") !== -1) {
				array[i] = power(i, term, array);
			}
			if (array[i].indexOf("sqrt") !== -1) {
			array[i] = square(i, term, array);
			}
			if (array[i].indexOf("!") !== -1) {
				array[i] = factorial(array[i]);
			}
			//TODO: MORE.
		}
				
		//put everything back into a string from the array
		term = array.join(' ');		
				
		//Check to see if there's Hex involved.
		if (term.indexOf("0x") !== -1){
			doMath(16, "0x", term, original);
		} //close hex
		
		//Check to see if there's Octal involved.
		else if (term.indexOf("0o") !== -1) {
			doMath(8, "0o", term, original);
		} //close octal
		
		//Check to see if there's binary involved.
		else if (term.indexOf("0b") !== -1) {
			doMath(2, "0b", term, original);
		} //close binary
		
		else
		{
		var value = eval(term);
			if (!isNaN(value)) {
				// display some html
				 document.getElementById("blank").innerHTML = "<center><div class='calculate' onmouseover='unhideBubble();' onmouseout='hideBubble();'>" + original + " = <strong>" + value + "</strong></div><div class='bubble'><strong>What's this?</strong><br/>What you serached seemed to us like it was math, so we did the math for you!</div></center>";
			}
		} //close else
	}
	catch (e) {}
}

// The function doMath(base, type, newterm) does math!
// More specifically, it does hex, octal, and binary 
// mathematical equations. 
function doMath(base, type, newterm, original) {
	term = newterm;
	term = term.replace("+", " + ");
	term = term.replace("-", " - ");
	term = term.replace("*", " * ");
	term = term.replace("/", " / ");
	term = term.replace("%", " % ");
	
	//THIS fIXES ALL DAMN SPACING ERRORS.
	term = term.replace(/^\s+|\s+$/g,'').replace(/\s+/g,' ');
	
	var array = term.split(" ");
		
	// Thanks to Prestaul of stack overflow
	// http://stackoverflow.com/questions/57803/how-to-convert-decimal-to-hex-in-javascript
	
	for (var i = 0; i < array.length; i++) {
		if (!(array[i] == "+" || array[i] == "-" || array[i] == "*" || array[i] == "/" || array[i] == "%")) {
			if (array[i].length > 2 && array[i].substring(0,2) == "0o" || array[i].substring(0,2) == "0b") {
				var temp = array[i].substring(2);
				array[i] = parseInt(temp, base);
			}
		}
	}
	term = array.join(' ');			
	var value = eval(term);
	value = value.toString(base);
	
	document.getElementById("blank").innerHTML = "<center><div class='calculate' onmouseover='unhideBubble();' onmouseout='hideBubble();'>" + original + " = <strong>" + type + value + "</strong></div><div class='bubble'><strong>What's this?</strong><br/>What you serached seemed to us like it was math, so we did the math for you!</div></center>";
}

// The function unhideBubble() unhides the bubble!
function unhideBubble() {
	document.getElementsByClassName("bubble").item(0).style.opacity = "1";
}

// The function hideBubble() hides the bubble!
function hideBubble() {
	document.getElementsByClassName("bubble").item(0).style.opacity = "0";
}

function power(i, term, array) {
	var before = array[i].substring(0,array[i].indexOf("^"));
	var after = array[i].substring(array[i].indexOf("^")+1);
	var beforeo = false;
	var beforeb = false;
	var beforeh = false;
	var befored = false;
		
	//Check to see if there's Hex involved.
	if (array[i].indexOf("0x") !== -1 && before.indexOf("0x") !== -1) {
		before = before.substring(2);
		before = parseInt(before, 16);
		beforeh = true;
	} //close before hex	
		
	//Check to see if there's Octal involved.
	else if (array[i].indexOf("0o") !== -1 && before.indexOf("0o") !== -1) {
		before = before.substring(2);
		before = parseInt(before, 8);
		beforeo = true;
	} //close before octal
									
	//Check to see if there's Binary involved.
	else if (array[i].indexOf("0b") !== -1 && before.indexOf("0b") !== -1) {
		before = before.substring(2);
		before = parseInt(before, 2);
		beforeb = true;
	} //close before octal
	
	else {
		before = parseInt(before, 10);
		befored = true;
	}
	
	if (array[i].indexOf("0x") !== -1 && after.indexOf("0x") !== -1) {
		after = after.substring(2);
		after = parseInt(after, 16);
	} //close after hex
	
	else if (array[i].indexOf("0b") !== -1 && after.indexOf("0b") !== -1) {
		after = after.substring(2);
		after = parseInt(after, 2);
	} //close after octal
	
	else if (array[i].indexOf("0o") !== -1 && after.indexOf("0o") !== -1) {
		after = after.substring(2);
		after = parseInt(after, 8);
	} //close after octal
	
	else {
		after = parseInt(after, 10);
	}
					
	array[i] = Math.pow(before, after);
					
	if (beforeh) array[i] = "0x" + array[i].toString(16);
	if (beforeo) array[i] = "0o" + array[i].toString(8);
	if (beforeb) array[i] = "0b" + array[i].toString(2);
	
	return array[i];
}

function square(i, term, array) {
	var value = array[i].substring(array[i].indexOf("(")+1, array[i].indexOf(")"));
	var h = false;
	var o = false;
	var b = false;
	
	//Check to see if there's Hex involved.
	if (array[i].indexOf("0x") !== -1) {
		value = value.substring(2);
		value = parseInt(value, 16);
		h = true;
	} //close before hex
	
	//Check to see if there's Octal involved.
	else if (array[i].indexOf("0o") !== -1) {
		value = value.substring(2);
		value = parseInt(value, 8);
		o = true;
	} //close before octal
									
	//Check to see if there's Binary involved.
	else if (array[i].indexOf("0b") !== -1) {
		value = value.substring(2);
		value = parseInt(value, 2);
		b = true;
	} //close before octal
	
	else {
		value = parseInt(value, 10);
	}
	
	value = Math.sqrt(value);
	array[i] = value;
					
	if (h) array[i] = "0x" + array[i].toString(16);
	if (o) array[i] = "0o" + array[i].toString(8);
	if (b) array[i] = "0b" + array[i].toString(2);
	
	return array[i];
}

function factorial(n) { 
	var h = false;
	var o = false;
	var b = false;
	if (n.indexOf("!") !== -1) {
		n = n.substring(0,n.indexOf("!"));
	}
	
	//Check to see if there's Hex involved.
	if (n.indexOf("0x") !== -1) {
		n = n.substring(2);
		n = parseInt(n, 16);
		h = true;
	} //close before hex
	
	//Check to see if there's Octal involved.
	else if (n.indexOf("0o") !== -1) {
		n = n.substring(2);
		n = parseInt(n, 8);
		o = true;
	} //close before octal
									
	//Check to see if there's Binary involved.
	else if (n.indexOf("0b") !== -1) {
		n = n.substring(2);
		n = parseInt(n, 2);
		b = true;
	} //close before octal
	
	var result = 1;
	for (var i = n; i > 0; i--) {
		result *= i;
	}
	
	if (h) result = "0x" + result.toString(16);
	if (o) result = "0o" + result.toString(8);
	if (b) result = "0b" + result.toString(2);
	
	return result;
}
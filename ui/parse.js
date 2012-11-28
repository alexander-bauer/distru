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
		
		var operators = {"plus" : "+", "and" : "+", "minus" : "-", "times" : "*",
						 "over" : "/", "divide" : "/", "mod" : "%", "modulus" : "%"};
		
		for (var val in operators) {
			term = term.replace(new RegExp(val, "g"), operators[val]);
		}
		
		
		if (isHex(term)){
			var original = term;
			term = term.replace("+", " + ");
			term = term.replace("-", " - ");
			term = term.replace("*", " * ");
			term = term.replace("/", " / ");
			term = term.replace("%", " % ");
			
			term = term.replace(/^\s+|\s+$/g,'').replace(/\s+/g,' ');
			
			var array = term.split(" ");
			
			
			for (var i = 0; i < array.length; i++) {
				if (!(array[i] == "+" || array[i] == "-" || array[i] == "*" || array[i] == "/" || array[i] == "%")) {
					if (array[i].length > 2 && array[i].substring(0,2) == "0o") {
						var temp = array[i].substring(2);
						array[i] = parseInt(temp, 16);
					}
					else if (array[i].length > 2 && array[i].substring(0,2) == "0b") {
						var temp = array[i].substring(2);
						array[i] = parseInt(temp, 16);
					}
					else array[i] = parseInt(array[i], 16);
				}
			}
			term = array.join(' ');			
			var value = eval(term);
			value = value.toString(16);
			
			document.getElementById("blank").innerHTML = "<center><div class='calculate' onmouseover='unhideBubble();' onmouseout='hideBubble();'>" + original + " = <strong>0x" + value + "</strong></div><div class='bubble'><strong>What's this?</strong><br/>What you serached seemed to us like it was math, so we did the math for you!</div></center>";
		} //close hex
		
		else if (isOctal(term)) {
			var original = term;
			term = term.replace("+", " + ");
			term = term.replace("-", " - ");
			term = term.replace("*", " * ");
			term = term.replace("/", " / ");
			term = term.replace("%", " % ");
			
			term = term.replace(/^\s+|\s+$/g,'').replace(/\s+/g,' ');
			
			var array = term.split(" ");
			
			
			for (var i = 0; i < array.length; i++) {
				if (!(array[i] == "+" || array[i] == "-" || array[i] == "*" || array[i] == "/" || array[i] == "%")) {
					if (array[i].length > 2 && array[i].substring(0,2) == "0o") {
						var temp = array[i].substring(2);
						array[i] = parseInt(temp, 8);
					}
					else if (array[i].length > 2 && array[i].substring(0,2) == "0b") {
						var temp = array[i].substring(2);
						array[i] = parseInt(temp, 2);
					}
					else array[i] = parseInt(array[i], 8);
				}
			}
			term = array.join(' ');			
			var value = eval(term);
			value = value.toString(8);
			
			document.getElementById("blank").innerHTML = "<center><div class='calculate' onmouseover='unhideBubble();' onmouseout='hideBubble();'>" + original + " = <strong>0o" + value + "</strong></div><div class='bubble'><strong>What's this?</strong><br/>What you serached seemed to us like it was math, so we did the math for you!</div></center>";
		} //close octal
		
		else if (isBinary(term)) {
			var original = term;
			term = term.replace("+", " + ");
			term = term.replace("-", " - ");
			term = term.replace("*", " * ");
			term = term.replace("/", " / ");
			term = term.replace("%", " % ");
			
			term = term.replace(/^\s+|\s+$/g,'').replace(/\s+/g,' ');
			
			var array = term.split(" ");
			
			
			for (var i = 0; i < array.length; i++) {
				if (!(array[i] == "+" || array[i] == "-" || array[i] == "*" || array[i] == "/" || array[i] == "%")) {
					if (array[i].length > 2 && array[i].substring(0,2) == "0b") {
						var temp = array[i].substring(2);
						array[i] = parseInt(temp, 2);
					}
					else if (array[i].length > 2 && array[i].substring(0,2) == "0o") {
						var temp = array[i].substring(2);
						array[i] = parseInt(array[i], 2);
					}
					else array[i] = parseInt(array[i], 2);
				}
			}
			term = array.join(' ');			
			var value = eval(term);
			value = value.toString(2);
			
			document.getElementById("blank").innerHTML = "<center><div class='calculate' onmouseover='unhideBubble();' onmouseout='hideBubble();'>" + original + " = <strong>0b" + value + "</strong></div><div class='bubble'><strong>What's this?</strong><br/>What you serached seemed to us like it was math, so we did the math for you!</div></center>";
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

function isHex() {
	if (term.indexOf("0x") !== -1)
		return true;
	else return false;
}

function isOctal() {
	if (term.indexOf("0o") !== -1)
		return true;
	return false;
}

function isBinary() {
	if (term.indexOf("0b") !== -1)
		return true;
	return false;
}

function unhideBubble() {
	document.getElementsByClassName("bubble").item(0).style.opacity = "1";
}

function hideBubble() {
	document.getElementsByClassName("bubble").item(0).style.opacity = "0";
}
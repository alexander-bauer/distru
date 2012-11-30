#Math

Distru can parse math when it's inputed into the search bar and searched. This document is an explanation for the types of math that Distru can handle, and how to use it. Distru supports math for the base 2 (binary), 8 (octal), 10 (decimal), and 16 (hexadecimal). When doing any type of math, you can combine different bases and Distru will parse it for you; the first base will be the resulting base.

###Bases
Distru supports four different bases. In order to use base 2, 8, and 16, a prefix must be added to the number that you want to use.

```math
Hex: 0x
Octal: 0o
Binary: 0b
```

##Basic Math
Distru supports basic math! The math can be combinations of any type of bases, and the spaces between the number and operators don't matter. 

**Examples:**
```math
3 + 3
0xFF - 0b1
0o3 * 8
32 / 2
```

##Exponents
Distru supports exponents!

**Examples:**
```math
10^3
0xF^3
0o3^0xA
0b1101^3
```

##Factorial
Distru supports factorial!

**Examples:**
```math
5!
0xF!
0o4!
0b11001!
```

##Square Root
Distru supports square root!

**Examples:**
```math
sqrt(4)
sqrt(0x4)
sqrt(0o4)
sqrt(0b100)
```
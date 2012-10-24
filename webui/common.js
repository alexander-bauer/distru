$(document).ready(function(){
	$(window).resize(function(){
		
		$('.search').css({
			position:'absolute',
		    left: ($(window).width() 
			- $('.search').outerWidth())/2,
		    top: ($(window).height() 
			- $('.search').outerHeight())/2
			});
		
		$('.name').css({
			position:'absolute',
		    left: ($(window).width() 
			- $('.name').outerWidth())/2,
		    top: ($(window).height() 
			- $('.name').outerHeight())/2 - 50
			});
			
 		});
 	$(window).resize();
	$(window).resize();
});
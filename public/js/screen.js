var container, screen, panel, fullScreenButton;
var webSocket;
var panelHideTimer;
var fullscreenEnabled = false;

$(document).ready(function() {
	// get container element
	container = $('#container');
	container.mousemove(function(event) {
		clearInterval(panelHideTimer);
		panel.attr('display', 'block');
		panel.fadeIn(100);
		panelHideTimer = setInterval(hidePanel, 1500);
	});

	// get panel element
	panel = $('#panel');

	// get fullscreen butoon element
	fullScreenButton = $('#full-screen-btn');
	fullScreenButton.click(toggleFullscreen);
	
	// get screen element
	screen = $('#screen');
	screen.dblclick(toggleFullscreen);

	// connect to websocket
	webSocket = $.simpleWebSocket({
		url: socketAddress(),
		dataType: 'arraybuffer'
	});

	// start listening
	webSocket.listen(function(message) {
		updateScreen(message);
	});
});

function toggleFullscreen() {
	fullScreenButton.css('background', fullscreenEnabled
			? "url('/public/img/fullscreen-exit.png') no-repeat"
			: "url('/public/img/fullscreen.png') no-repeat");

	fullscreenEnabled ? disableFullscreen() : enableFullscreen(screen[0]);
	fullscreenEnabled = !fullscreenEnabled;
}

function hidePanel() {
	panel.fadeOut(500);
	panel.attr('display', 'none');
	clearInterval(panelHideTimer);
}

function enableFullscreen(elem) {
	if (elem.requestFullscreen) {
    	elem.requestFullscreen();
  	} else if (elem.mozRequestFullScreen) {
    	elem.mozRequestFullScreen();
  	} else if (elem.webkitRequestFullscreen) {
    	elem.webkitRequestFullscreen();
  	} else if (elem.msRequestFullscreen) {
    	elem.msRequestFullscreen();
  	}
}

function disableFullscreen() {
	if (document.exitFullscreen) {
    	document.exitFullscreen();
  	} else if (document.mozCancelFullScreen) {
    	document.mozCancelFullScreen();
  	} else if (document.webkitExitFullscreen) {
    	document.webkitExitFullscreen();
  	} else if (document.msExitFullscreen) {
    	document.msExitFullscreen();
  	}
}

function updateScreen(blob) {
	var reader = new FileReader();
	reader.readAsDataURL(blob);
	reader.onloadend = function() {
		screen.attr('src', reader.result);
	};
}
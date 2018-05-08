(function (){
	window.addEventListener('load', init, false)

	function init() {
		var xhr = new XMLHttpRequest()
		xhr.open("GET", "/html")
		xhr.onload = function () {
			window.app.parse(xhr.responseText)
		}
		xhr.send()
	}
})((window.app = window.app || {}));
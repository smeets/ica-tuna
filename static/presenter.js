(function (app) {

	function div(cc, text) {
		var s = document.createElement('div')
		s.innerHTML = text
		s.classList.add(cc)
		return s
	}

	function image(src) {
		var s = document.createElement('img')
		s.src = src
		return s
	}

	function build(offer) {
		var el = document.createElement('div')
		el.classList.add('offer')

		el.appendChild(div('name', offer.name))
		if (offer.amount[0] === '/') {
			el.appendChild(div('price', offer.price))
			el.appendChild(div('amount', offer.amount))
		} else {
			el.appendChild(div('amount', offer.amount))
			el.appendChild(div('price', offer.price))
		}
		el.appendChild(image(offer.img))
		el.appendChild(div('info', offer.info))

		return el
	}
	app.present = function present(offers) {
		offers.forEach(offer =>
			document.body.appendChild(build(offer))
		)
	}
})((window.app = window.app || {}));
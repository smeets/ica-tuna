(function (app) {

	function div(cc, text) {
		var s = document.createElement('div')
		var t = document.createTextNode(text)
		s.appendChild(t)
		s.classList.add(cc)
		return s
	}

	function image(src) {
		var s = document.createElement('img')

		if (src.includes("icanet")) {
			src = "/proxy?q=" + encodeURI(src)
		}

		s.src = src
		return s
	}

	function recordHandler(e) {
		const deal = e.target.parentNode
		if (!deal.classList.contains('deal')) {
			return
		}

		if (deal.classList.contains('record')) {
			return
		}

		deal.classList.add('record')
		fetch("/record", {
			method: 'POST', // *GET, POST, PUT, DELETE, etc.
			headers: {
			  'Content-Type': 'application/json'
			},
			body: JSON.stringify([{
				offer: deal.dataset.offer,
				label: deal.dataset.label
			}])
		}).then(function (res) {
			if (res.ok) return undefined
			return res.text()
		}).then(function (txt) {
			if (!txt) return
			window.alert(txt)
			deal.classList.remove('record')
		})
	}

	function build(offer) {
		var el = document.createElement('div')
		el.classList.add('offer')

		var deal = document.createElement('div')
		deal.dataset["label"] = offer.name
		deal.classList.add('deal')
		deal.addEventListener('click', recordHandler)
		el.appendChild(deal)

		if (offer.amount[0] === '/') {
			deal.appendChild(div('amount', "["))
			deal.appendChild(div('price', offer.price))
			deal.appendChild(div('amount', offer.amount+"]"))
			deal.dataset["offer"] = offer.price+offer.amount
		} else {
			var amount = offer.amount.substring(0, offer.amount.indexOf(' ')) + "st"
			deal.appendChild(div('amount', "["))
			deal.appendChild(div('price', offer.price))
			deal.appendChild(div('amount', "/"+amount+"]"))
			deal.dataset["offer"] = offer.price+"/"+amount
		}

		el.appendChild(div('name', offer.name))

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

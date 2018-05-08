(function (app) {

	function parse_offer(offer) {
		var name  = "??"
		var info  = "??"
		var price = "??"
		var amount= ""
		var img   = "??"

		{
			var name_start = offer.indexOf("<h2")
			var name_start_end = offer.indexOf(">", name_start)
			var name_end = offer.indexOf("</h2>", name_start_end)
			name = offer.substring(name_start_end + 1, name_end).trim()
		}

		{
			var info_start = offer.indexOf('<p class="offer-type__product-info">')
			info_start += '<p class="offer-type__product-info">'.length
			var info_end = offer.indexOf("</p>", info_start)
			info = offer.substring(info_start, info_end).trim()
		}

		{
			var price_start = offer.indexOf('<div class="product-price__price-value">')
			price_start += '<div class="product-price__price-value">'.length
			var price_end = offer.indexOf("</div>", price_start)
			price = offer.substring(price_start, price_end).trim()
		}

		{
			var amount_start = offer.indexOf('<div class="product-price__amount">')
			if (amount_start !== -1) {
				amount_start += '<div class="product-price__amount">'.length
			} else {
				amount_start = offer.indexOf('<div class="product-price__unit-item">')
				amount_start += '<div class="product-price__unit-item">'.length
			}
			var amount_end = offer.indexOf("</div>", amount_start)
			amount = offer.substring(amount_start, amount_end).trim()
		}

		{
			var img_start = offer.indexOf('<img class="lazy"')
			var src_start = offer.indexOf('data-original="', img_start) + 'data-original="'.length
			var src_end   = offer.indexOf('"', src_start)
			img = offer.substring(src_start, src_end)
				.trim()
				.replace(/&amp;/g, "&")
		}

		return {
			name:  name,
			price: price,
			info:  info,
			amount: amount,
			img:   img
		}
	}

	app.parse = function parse (news) {
		var end_length = "</div>".length
		var start_length = 'class=" offer-category__item">'.length
		var ranges = []
		var start  = 0
		var end    = 0

		while (true) {
			start = news.indexOf('class=" offer-category__item">', end)
			if (start === -1) break
			start += start_length

			var open = 1
			var loc  = start
			while (open > 0) {
				var div = news.substring(loc, loc + 6)
				if (div === "<div c") open++
				else if (div === "</div>") open--
				loc += 1
			}
			end = loc + end_length
			ranges.push({
				start: start,
				end: end
			})
		}

		console.log(ranges)

		var sections = ranges.map(range => news.substring(range.start, range.end))
		var offers = sections.map(parse_offer)
		window.app.present(offers)

		//document.body.innerHTML = news
	}
})((window.app = window.app || {}));
(function (app){
    app.analyse = function analyse(offers) {
        const deals = offers.map(offer => document.querySelector(`[data-label="${offer.name}"`))
            .filter(d => d)
            .map(d => d.dataset)
            .map(d => ({
                label: d.label,
                offer: d.offer,
            }))
        console.log(deals.length, deals)

        fetch("/compare", {
            method: 'POST', // *GET, POST, PUT, DELETE, etc.
            headers: {
              'Content-Type': 'application/json'
            },
            body: JSON.stringify(deals)
        }).then(function (res) {
            return res.json()
        }).then(function (data) {
            console.log(data)
        })


        // var xhr = new XMLHttpRequest()
        // xhr.open("GET", "/html")
        // xhr.onload = function () {
        //     window.app.parse(xhr.responseText)
        // }
        // xhr.send()
    }
})((window.app = window.app || {}));

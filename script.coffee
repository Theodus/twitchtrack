update = () -> 
	$.ajax
		url: "/data"
		dataType: "json"
		error: (jqKHR, text) ->
			console.log text
			update()
		success: (data) ->
			console.log data
			tbody = document.getElementById("tbody")
			$(tbody).empty()
			for ch in data.channels
				tr = tbody.appendChild(document.createElement("tr"))
				td0 = tr.appendChild(document.createElement("td"))
				td1 = tr.appendChild(document.createElement("td"))
				td2 = tr.appendChild(document.createElement("td"))
				a = td0.appendChild(document.createElement("a"))
				a.innerHTML = ch.channel
				a.href = ch.url
				viewers = ch.viewers
				if viewers > 0
					td1.innerHTML = ch.game
					td2.innerHTML = ch.stream
	setTimeout(update, 120000)

update()

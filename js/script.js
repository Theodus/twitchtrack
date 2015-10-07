$.getJSON('/refresh', function(data) {
	console.log(data);
    var tbody = document.getElementById("tbody");
    data.channels.forEach(function(element, index){
        var tr = tbody.appendChild(document.createElement("tr"));
        var td0 = tr.appendChild(document.createElement("td"));
        var td1 = tr.appendChild(document.createElement("td"));
        var td2 = tr.appendChild(document.createElement("td"));
        var viewers = data.viewers[index];
        if(viewers>0) {
            var a = td0.appendChild(document.createElement("a"));
            a.innerHTML = data.channels[index];
            a.href = data.links[index];
        } else {
            td0.innerHTML = data.channels[index];
        }
        td1.innerHTML = viewers;
        td2.innerHTML = data.streams[index];
    });
});
setTimeout((function(){
	document.location.reload(true);
}), 60000);

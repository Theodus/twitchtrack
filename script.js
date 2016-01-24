(function(){
    $.ajax({
        url: "/data",
        success: function(data){
            console.log(data);
            var tbody = document.getElementById("tbody");
            while(tbody.hasChildNodes()) {
                tbody.removeChild(tbody.firstChild);
            }
            data.channels.forEach(function(element, index) {
                var tr = tbody.appendChild(document.createElement("tr"));
                var td0 = tr.appendChild(document.createElement("td"));
                var td1 = tr.appendChild(document.createElement("td"));
                var td2 = tr.appendChild(document.createElement("td"));
                var a = td0.appendChild(document.createElement("a"));
                a.innerHTML = element.channel;
                a.href = element.url;
                var viewers = element.viewers;
                if(viewers > 0) {
                    td1.innerHTML = element.game;
                    td2.innerHTML = element.stream;
                }
            });
        },
        dataType: "json",
        timeout: 30000,
    });
setTimeout(arguments.callee, 120000);
})();

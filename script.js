$.getJSON('/refresh', function(data){
  var tbody = document.getElementById("tbody");
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
});

setTimeout(function(){
  location.reload()
}, 120000);

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

setInterval(function(){
  $.getJSON('/refresh', function(data){
    var tbody = document.getElementById("tbody");
    data.channels.forEach(function(element, index) {
      var tr = tbody.childNodes[index+1];
      var td0 = tr.childNodes[0];
      var td1 = tr.childNodes[1];
      var td2 = tr.childNodes[2];
      var a = td0.childNodes[0];
      a.innerHTML = element.channel;
      a.href = element.url;
      var viewers = element.viewers;
      if(viewers > 0) {
        td1.innerHTML = element.game;
        td2.innerHTML = element.stream;
      }
    });
  });
}, 60000);

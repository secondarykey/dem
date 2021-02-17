var dialog = document.querySelector('#ProjectDialog');
var showDialogButton = document.querySelector('#registerProject');
if (! dialog.showModal) {
  dialogPolyfill.registerDialog(dialog);
}

showDialogButton.addEventListener('click', function() {
  dialog.showModal();
});

dialog.querySelector('.entry').addEventListener('click', function() {

  var projectid = document.querySelector('#projectID').value;
  var endpoint = document.querySelector('#endpoint').value;

  var xhr = new XMLHttpRequest();
  xhr.open('POST',"/project/add.json");
  xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
  xhr.responseType = 'json';
  xhr.onload = function() {
    var resp = xhr.response;
      if (resp.Success) {
        location.href = resp.Redirect
      } else {
        alert(resp.Message);
        dialog.close();
      }
  };
  xhr.send("projectid=" + projectid + "&endpoint=" + endpoint);
});

dialog.querySelector('.close').addEventListener('click', function() {
  dialog.close();
});

document.addEventListener("DOMContentLoaded", function() {
});

var lists = document.querySelectorAll('.list-item')

var projectID = document.querySelector('#settingID').value;
var currentKind = "";
var currentCursor = "";

lists.forEach(function(value) {
  value.addEventListener('click', function(e) {

    currentKind = e.target.getAttribute("data-name");
    currentCursor = "";
    view(currentKind,currentCursor);

  });
});

function view(kind,cursor) {

  if ( cursor == "" ) {
      var th = document.getElementById('table-header');
      th.innerHTML = "";
      var td = document.getElementById('table-body');
      td.innerHTML = "";
  }

  var xhr = new XMLHttpRequest();
  var url = "/" + projectID + "/" + kind + "/" + cursor;

  xhr.open('POST',url);
  xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
  xhr.responseType = 'json';

  xhr.onload = function() {
    var resp = xhr.response;
      if (resp.Success) {

        var th = document.getElementById('table-header');
        for ( var i = 0; i < resp.Header.length; i++ ) {
          var elm = document.createElement("th");
          elm.classList.add("mdl-data-table__cell--non-numeric");
          var txt = document.createTextNode(resp.Header[i]);
          elm.appendChild(txt);
          th.appendChild(elm)
        }

        var tb = document.getElementById('table-body');
        for ( var i = 0; i < resp.Data.length; i++ ) {
          var elm = document.createElement("tr");
          var row = resp.Data[i];
          for ( var j = 0; i < row.Values.length; i++ ) {
            var td = document.createElement("td");
            td.classList.add("mdl-data-table__cell--non-numeric");
            var txt = document.createTextNode(row.Values[i]);
            td.appendChild(txt);
            elm.appendChild(td);
          }
          tb.appendChild(elm)
        }

         
      } else {
        alert(resp.Message);
      }
  };
  xhr.send();
}

var currentKind = "";
var projectID = document.querySelector('#settingID').value;

document.addEventListener("DOMContentLoaded", function() {

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

  //kind list
  var lists = document.querySelectorAll('.list-item');

  lists.forEach(function(value) {
    value.addEventListener('click', function(e) {
      currentKind = e.target.getAttribute("data-name");
      view(currentKind,true);
    });
  });

  function view(kind,first) {

    if ( first ) {
      var th = document.getElementById('table-header');
      th.innerHTML = "";
      var td = document.getElementById('table-body');
      td.innerHTML = "";
    }

    var xhr = new XMLHttpRequest();
    var url = "/entity/view/" + kind;

    xhr.open('POST',url);
    xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
    xhr.responseType = 'json';

    xhr.onload = function() {
      var resp = xhr.response;
      if (!resp.Success) {
        alert(resp.Message);
        return;
      }

      if ( first ) {
        createHeader(resp.Header);
        changeCheckAll();
      }

      addData(resp.Data);

      var table = document.getElementById("view-table");
      componentHandler.upgradeElement(table,'MaterialDataTable');

      clearCheck();
    };
    xhr.send("first=" + first);
  }

  function createCheckboxLabel(id) {

    var label = document.createElement("label");
    label.setAttribute("for",id);
    label.classList.add("mdl-checkbox");
    label.classList.add("mdl-js-checkbox");
    label.classList.add("mdl-js-ripple-effect");
    label.classList.add("mdl-data-table__select");
    
    var input = document.createElement("input");
    input.setAttribute("type","checkbox");
    input.setAttribute("id",id);
    input.classList.add("mdl-checkbox__input");

    var btns = document.querySelectorAll(".remove-btn");
    input.addEventListener('change', function(e) {
      if ( e.target.checked ) {
        for ( var i = 0; i < btns.length; i++ ) {
            btns[i].disabled = false;
        }
      }
    });
    label.appendChild(input);
    return label;
  }

  function clearCheck() {
    var all = document.querySelector('#select-all');
    if ( all.checked ) {
        all.checked = false;
        return;
    }

    var table = document.querySelector('table');
    var boxes = table.querySelectorAll('tbody .mdl-checkbox__input');
    for (var i = 0, length = boxes.length; i < length; i++) {
      boxes[i].checked = false;
    }
  }

  function changeCheckAll() {

    var table = document.querySelector('table');
    var headerCheckbox = table.querySelector('thead .mdl-data-table__select input');
    var headerCheckHandler = function(event) {
      var boxes = table.querySelectorAll('tbody .mdl-checkbox__input');
      if (event.target.checked) {
        for (var i = 0, length = boxes.length; i < length; i++) {
          boxes[i].checked = true;
        }
      } else {
        for (var i = 0, length = boxes.length; i < length; i++) {
          boxes[i].checked = false;
        }
      }
    };
    headerCheckbox.addEventListener('change', headerCheckHandler);
  }

  function createHeader(header) {
    var th = document.getElementById('table-header');
    var elm = document.createElement("th");
    var label = createCheckboxLabel("select-all")
    elm.appendChild(label);
    th.appendChild(elm)

    for ( var i = 0; i < header.length; i++ ) {
      var elm = document.createElement("th");
      elm.classList.add("mdl-data-table__cell--non-numeric");
      var txt = document.createTextNode(header[i]);
      elm.appendChild(txt);
      th.appendChild(elm)
    }
  }

  function addData(data) {

    var tb = document.getElementById('table-body');

    var colLen = 0;
    for ( var i = 0; i < data.length; i++ ) {

      var row = data[i];
      var elm = document.createElement("tr");

      var td = document.createElement("td");
      var label = createCheckboxLabel(row.Key);
      td.appendChild(label);
      elm.appendChild(td);

      colLen = row.Values.length
      for ( var j = 0; j < colLen; j++ ) {
        var td = document.createElement("td");
        td.classList.add("mdl-data-table__cell--non-numeric");
        var txt = document.createTextNode(row.Values[j]);
        td.appendChild(txt);
        elm.appendChild(td);
      }
      tb.appendChild(elm)
    }

    if ( colLen != 0 ) {
      var tr = document.createElement("tr");
      tr.id = "nextUpdate";
      tr.addEventListener("click",function(e) {
        tr.parentElement.removeChild(tr);
        view(currentKind,false);
      });

      var td = document.createElement("td");
      td.setAttribute("style","text-align:center;");
      td.setAttribute("colspan",colLen + 1);
      td.textContent = "Next";
      tr.appendChild(td);
      tb.appendChild(tr);
    }
  }

  var delBtns = document.querySelectorAll('.remove-btn')
  for ( var i = 0; i < delBtns.length; i++ )  {     
    delBtns[i].addEventListener('click', function(e) {
        deleteRows(delBtns);
    });
  }

  function deleteRows(delBtns) {

    for ( var i = 0; i < delBtns.length; i++ )  {     
        delBtns[i].disabled = true;
    }

    var table = document.querySelector('table');
    var boxes = table.querySelectorAll('tbody .mdl-checkbox__input');

    var ids = new Array();
    var rows = new Array();

    for (var i = 0, length = boxes.length; i < length; i++) {
      if ( boxes[i].checked ) {
        var row = boxes[i].parentNode.parentNode.parentNode;
        ids.push(boxes[i].id);
        rows.push(row);
      }
    } 

    var xhr = new XMLHttpRequest();
    xhr.open('POST',"/entity/remove/" + currentKind);
    xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
    xhr.responseType = 'json';
    xhr.onload = function() {
      var resp = xhr.response;
        if (resp.Success) {
          for (var i = 0, length = rows.length; i < length; i++) {
            var row = rows[i];
            row.parentNode.removeChild(row);
          }
        } else {
          alert(resp.Message);
        }
    };
    xhr.send("ids=" + JSON.stringify(ids));
  }

  var list = document.querySelectorAll('.limit-list');
  for ( var i = 0; i < list.length; ++i ) {
    list[i].addEventListener("click",function(e) {
      var limit = e.target.textContent;
      document.getElementById("limit-text").textContent = limit;

      var xhr = new XMLHttpRequest();
      xhr.open('POST',"/entity/limit/" + limit);
      xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
      xhr.responseType = 'json';

      xhr.onload = function() {
        var resp = xhr.response;
          if (resp.Success) {
            if ( currentKind != "" ) {
              view(currentKind,true);
            }
          } else {
            alert(resp.Message);
          }
      };
      xhr.send();
    });
  }

  var list = document.querySelectorAll('.ns-list');
  for ( var i = 0; i < list.length; ++i ) {
    list[i].addEventListener("click",function(e) {
      var ns = e.target.textContent;
      document.getElementById("ns-text").textContent = ns;

      var xhr = new XMLHttpRequest();
      xhr.open('POST',"/entity/limit/" + ns);
      xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
      xhr.responseType = 'json';

      xhr.onload = function() {
        var resp = xhr.response;
          if (resp.Success) {
            if ( currentKind != "" ) {
              view(currentKind,true);
            }
          } else {
            alert(resp.Message);
          }
      };
      xhr.send();
    });
  }

});

var darkMode = document.querySelector('#darkmode');
darkMode.addEventListener('change', function() {
  location.href = "/view/dark/" + darkMode.checked;
});


var layout = document.querySelector(".mdl-layout");
var clicked = false;
layout.addEventListener("mdl-componentupgraded",function(e) {
  if ( projectID == "empty" ) {
    var btn = document.querySelector(".mdl-layout__drawer-button");
    if (btn != null && !clicked) { 
      btn.click();
      clicked = true;
    }
  }
});



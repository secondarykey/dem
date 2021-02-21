var cursor = "";

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
    var params = new Object();
    var projectid = document.querySelector('#projectID').value;
    var endpoint = document.querySelector('#endpoint').value;
    params.projectid = projectid;
    params.endpoint = endpoint;
    request("/project/add.json",params,function(resp) {
      location.href = resp.Redirect
    });
  });

  dialog.querySelector('.close').addEventListener('click', function() {
    dialog.close();
  });

  //kind list
  var lists = document.querySelectorAll('.list-item');

  lists.forEach(function(value) {
    value.addEventListener('click', function(e) {
      setCurrent("kind",e.target.getAttribute("data-name"));
      view(true);
    });
  });

  function view(first) {

    if ( getCurrent("kind") == "" ) {
      return "";
    }

    if ( first ) {
      cursor = "";
      var th = document.getElementById('table-header');
      th.innerHTML = "";
      var td = document.getElementById('table-body');
      td.innerHTML = "";
    }

    var url = "/entity/view";
    var params = new Object();
    params.first = first;
    params.cursor = cursor;

    request(url,params,function(resp) {

      cursor = resp.Next;

      if ( first ) {
        createHeader(resp.Header);
        changeCheckAll();
      }

      addData(resp.Data);

      var table = document.getElementById("view-table");
      componentHandler.upgradeElement(table,'MaterialDataTable');

      clearCheck();
    });
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
        view(false);
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
      confirmDem("Delete?","Do you want to delete the selected data?",function() {
        deleteRows(delBtns);
      });
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

    var params = new Object();
    params.ids = JSON.stringify(ids);
    request("/entity/remove",params,function(resp) {
      for (var i = 0, length = rows.length; i < length; i++) {
        var row = rows[i];
        row.parentNode.removeChild(row);
      }
    });
  }

  var list = document.querySelectorAll('.limit-list');
  for ( var i = 0; i < list.length; ++i ) {
    list[i].addEventListener("click",function(e) {

      var limit = e.target.textContent;
      document.getElementById("limit-text").textContent = limit;
      setCurrent("limit",limit)
      request("/entity/limit/change",new Object(),function(resp) {
        view(true);
      });
    });
  }

  var list = document.querySelectorAll('.ns-list');
  for ( var i = 0; i < list.length; ++i ) {
    list[i].addEventListener("click",function(e) {
      var ns = e.target.textContent;
      document.getElementById("ns-text").textContent = ns;
      setCurrent("namespace",ns);
      request("/namespace/change",new Object(),function(resp) { 
        view(true);
      });
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

function request(url,params,successFunc) {

  var xhr = new XMLHttpRequest();

  xhr.open('POST',url);
  xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
  xhr.responseType = 'json';

  xhr.onload = function() {
    var resp = xhr.response;
    if ( resp.Success ) {
      successFunc(resp);
    } else {
      //Error
      Alert(resp.Message);
    }
  };

  //xhr.onerror = function() {
  //};
  params["ID"]        = document.getElementById("ID").value;
  params["kind"]      = document.getElementById("kind").value;
  params["limit"]     = document.getElementById("limit").value;
  params["namespace"] = document.getElementById("namespace").value;

  xhr.send(parsePostValue(params));
  return;
}

function setCurrent(id,value) {
  document.getElementById(id).value = value;
}

function getCurrent(id) {
  return document.getElementById(id).value;
}

function parsePostValue(params) {
  var value = "";
  Object.keys(params).forEach(function (k) {
    if ( value != "" ) {
      value += "&"
    }
    value += k + "=" + params[k];
  });
  return value;
}

var handler = (function(){
    var events = {},
    key = 0;
    return {
        add: function(target, type, listener, capture) {
            target.addEventListener(type, listener, capture);
            events[key] = {
                target: target,
                type: type,
                listener: listener,
                capture: capture
            };
            return key++;
        },
        remove: function(key) {
            if(key in events) {
                var e = events[key];
                e.target.removeEventListener(e.type, e.listener, e.capture);
            }
        }
    };
}());

function confirmDem(title,msg,yesFunc) {
  var dialog = document.querySelector('#ConfirmDialog');
  if (!dialog.showModal) {
    dialogPolyfill.registerDialog(dialog);
  }

  var titleElm = dialog.querySelector('#confirmTitle');
  var msgElm = dialog.querySelector('#confirmMessage');
  titleElm.textContent = title;
  msgElm.textContent = msg;

  var yes = dialog.querySelector('.yes');
  var no = dialog.querySelector('.no');
  dialog.hide = function() {
    dialog.close();
    handler.remove(dialog.yes);
    handler.remove(dialog.no);
  }

  var yesKey = handler.add(yes,"click",function() {
    yesFunc();
    dialog.hide();
  });

  var noKey = handler.add(no,"click",function() {
    dialog.hide();
  });

  dialog.yes = yesKey;
  dialog.no = noKey;

  dialog.showModal();
}



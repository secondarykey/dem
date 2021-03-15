var cursor = "";
var allMode = false;

document.addEventListener("DOMContentLoaded", function() {

  var dialog = document.querySelector('#projectDialog');
  if (! dialog.showModal) {
    dialogPolyfill.registerDialog(dialog);
  }

  var entry = dialog.querySelector('.entry');
  var close = dialog.querySelector('.close');

  var entryKey = handler.add(entry,'click', function() {
    var params = new Object();
    var projectid = document.querySelector('#projectID').value;
    var endpoint = document.querySelector('#endpoint').value;
    params.projectid = projectid;
    params.endpoint = endpoint;
    request("/project/add.json",params,function(resp) {
      location.href = resp.Redirect
    });
  });

  var closeKey = handler.add(close,'click', function() {
    dialog.close();
  });

  var registerProject = document.querySelector('#registerProject');
  registerProject.addEventListener('click', function() {
    dialog.showModal();
  });

  var lists = document.querySelectorAll('.remove-project');
  lists.forEach(function(value) {
    value.addEventListener('click', function(e) {
      var id = e.target.getAttribute("data-id");
      confirmDem("Delete?","Deleteing a project is irreversibele,is that okay?",function() {
        location.href = "/project/remove/" + id;
      }) ;
    });
  });

  //kind list
  var lists = document.querySelectorAll('.list-item');

  lists.forEach(function(value) {
    value.addEventListener('click', function(e) {
      setCurrent("kind",e.target.getAttribute("data-name"));
      view(true);
    });
  });

  var lock = false;
  function view(first) {

    if ( getCurrent("kind") == "" ) {
      return "";
    }

    if ( lock ) {
        return;
    }

    lock = true;

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
      lock = false;
    },function(err) {
      lock = false;
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

    input.addEventListener('change', function(e) {
      if ( allMode ) return;
      changeButton();
    });

    label.appendChild(input);
    return label;
  }

  function changeButton() {
    var table = document.querySelector('table');
    var boxes = table.querySelectorAll('tbody .mdl-checkbox__input');

    var remove = true;
    var view = true;

    for ( var i = 0; i < boxes.length; i++ ) {
      if ( boxes[i].checked ) {
        remove = false;
        if ( !view ) {
          view = true;
          break;
        }
        view = false;
      }
    }
    disabledButton(view,remove);
  }
      

  function disabledButton(view,remove) {
    var views = document.querySelectorAll('.view-btn')
    for ( var i = 0; i < views.length; i++ ) {
      views[i].disabled = view;
    }
    var dels = document.querySelectorAll('.remove-btn')
    for ( var i = 0; i < dels.length; i++ ) {
      dels[i].disabled = remove;
    }
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
      allMode = true;
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

      view = true;
      if ( boxes.length == 1 ) {
        view = false;
      }
      disabledButton(view,!event.target.checked);
      allMode = false;
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
        var txt = document.createTextNode(row.Values[j].View);
        td.appendChild(txt);
        elm.appendChild(td);
      }

      var cc = -1;
      elm.addEventListener("mousedown",function(e) {
         var t = e.target;
         if ( t.tagName != "TD" ) {
           return;
         } 
         var check = t.parentElement.querySelector("input");

         if ( cc == -1 ) {
           cc = 1;
           setTimeout( function() {
             if ( cc != 2 )  {
               check.checked = !check.checked;
               changeButton();
             }
             cc = -1;
           },200);
         } else {
           cc = 2;
           e.preventDefault();
           showEntityDialog(check.id);
         }
      });

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

      var icon = document.createElement("i");
      icon.classList.add("material-icons");
      icon.textContent = "next_plan";

      td.appendChild(icon)
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

  var viewBtns = document.querySelectorAll('.view-btn')
  for ( var i = 0; i < viewBtns.length; i++ )  {     
    viewBtns[i].addEventListener('click', function(e) {
      var table = document.querySelector('table');
      var boxes = table.querySelectorAll('tbody .mdl-checkbox__input');
      var id = "";
      for ( var i = 0; i < boxes.length; i++ ) {
        if (!boxes[i].checked) {
          continue;
        }
        id = boxes[i].id
        break;
      }
      showEntityDialog(id);
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
        view(true);
    });
  }

  var list = document.querySelectorAll('.limit-list');
  for ( var i = 0; i < list.length; ++i ) {
    list[i].addEventListener("click",function(e) {

      var limit = e.target.getAttribute("data-value");
      document.getElementById("limit-text").setAttribute("data-value",limit);
      document.getElementById("limit-text").textContent = e.target.textContent;
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

  function fitWindow() {
    var table = document.getElementById("view-table");
    table.style.height = (window.innerHeight - 140) + "px";
  }
  window.addEventListener('resize',function() {
      fitWindow();
  }, false );
  fitWindow();

});

var darkMode = document.querySelector('#darkmode');
darkMode.addEventListener('change', function() {
  location.href = "/view/dark/" + darkMode.checked;
});


var layout = document.querySelector(".mdl-layout");
var clicked = false;
layout.addEventListener("mdl-componentupgraded",function(e) {
  if ( getCurrent("ID") == "" ) {
    var btn = document.querySelector(".mdl-layout__drawer-button");
    if (btn != null && !clicked) { 
      btn.click();
      clicked = true;
    }
  }
});

function request(url,params,successFunc,errorFunc) {

  var xhr = new XMLHttpRequest();

  xhr.open('POST',url);
  xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
  xhr.responseType = 'json';

  xhr.onload = function() {
    var resp = xhr.response;
    if ( resp.Success ) {
      successFunc(resp);
    } else {
      alertDem(resp.Message,resp.Detail,function() {
        if ( errorFunc !== undefined ) {
            errorFunc(resp);
        }
      });
    }
  };

  xhr.onerror = function() {
    var resp = xhr.response;
    alertDem(resp.Message,resp.Detail);
    if ( errorFunc !== undefined ) {
      errorFunc(resp);
    }
  };

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


function alertDem(title,msg,okFunc) {
  var dialog = document.querySelector('#alertDialog');
  if (!dialog.showModal) {
    dialogPolyfill.registerDialog(dialog);
  }

  var titleElm = dialog.querySelector('#alertTitle');
  var msgElm = dialog.querySelector('#alertMessage');
  titleElm.textContent = title;
  msgElm.textContent = msg;

  var ok = dialog.querySelector('.ok');
  dialog.hide = function() {
    dialog.close();
    handler.remove(dialog.okKey);
  }

  var okKey = handler.add(ok,"click",function() {
    if ( okFunc !== undefined ) {
      okFunc();
    }
    dialog.hide();
  });

  dialog.okKey = okKey;
  dialog.showModal();
}

function confirmDem(title,msg,yesFunc,noFunc) {

  var dialog = document.querySelector('#confirmDialog');
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
    if ( yesFunc !== undefined ) {
      yesFunc();
    }
    dialog.hide();
  });

  var noKey = handler.add(no,"click",function() {
    if ( noFunc !== undefined ) {
      noFunc();
    }
    dialog.hide();
  });

  dialog.yes = yesKey;
  dialog.no = noKey;

  dialog.showModal();
}

function setViewData(tag,name,type,current) {

    var v = current.View;
    var r = current.Real;

    console.log(v);
    console.log(r);

    switch ( type ) {
        case 10: //Omitted
            //長過ぎるので本当の表示をクリックしたら表示
            tag.textContent = v;
            tag.addEventListener("click",function() {
                tag.textContent = r;
            });
            break;
        case 20: //Expand
            //構造体なので表示を別途行う
            break;
        case 30: //Slice
            //インデックスでの表示にする
            break;
        case 40: //Download
            //ダウンロードする
            tag.textContent = v;
            break;
        case -1: //Error
            //赤文字にする？
            break;
        default: //Normal
            //そのまま表示
            tag.textContent = v;
            break;
    }

    return;
}

function showEntityDialog(id) {
  var params = new Object();
  params.key = id;
  request("/entity/get",params,function(resp) {

    var dialog = document.querySelector('#entityDialog');
    if (!dialog.showModal) {
      dialogPolyfill.registerDialog(dialog);
    }

    var kind = getCurrent("kind");

    var keyElm = dialog.querySelector('#entityKey');
    keyElm.textContent = kind + "(Key=" + resp.Entity.Key + ")";

    var contentElm = dialog.querySelector('#entityContent');
    contentElm.innerHTML = "";

    for ( var i = 0; i < resp.Header.length; i++ ) {
      var names = document.createElement("div");
      names.classList.add("mdl-cell");
      names.classList.add("mdl-cell--4-col");

      var name = resp.Header[i];
      names.textContent = name;
      contentElm.appendChild(names);

      var values = document.createElement("div");
      values.classList.add("mdl-cell");
      values.classList.add("mdl-cell--8-col");

      values.style.whiteSpace = "pre-wrap";

      var current = resp.Entity.Values[i];
      var type = resp.Entity.Types[i];

      setViewData(values,name,type,current)

      contentElm.appendChild(values);
    }

    var edit = dialog.querySelector('.edit');
    var close = dialog.querySelector('.close');

    dialog.hide = function() {
      dialog.close();
      handler.remove(dialog.editKey);
      handler.remove(dialog.closeKey);
    }

    var editKey = handler.add(edit,"click",function() {
        //TODO
      dialog.hide();
    });

    var closeKey = handler.add(close,"click",function() {
      dialog.hide();
    });

    dialog.editKey = editKey;
    dialog.closeKey = closeKey;

    dialog.showModal();

  });
}

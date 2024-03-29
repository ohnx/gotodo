/* global localStorage, showdown */

let converter = new showdown.Converter({
  headerLevelStart: 2,
  strikethrough: true,
  tables: true,
  tasklists: true,
  emoji: true
});

// LocalStorage helper
const LOCALSTORAGE_KEYS = {
  TOKEN: 1,
  USERNAME: 2,
  TAGS: 3
};
const API_ROOT = "http://nuc.int.masonx.ca:8080/api";

// Global variables
var todos = [];
var tags = [];
var focus_id = -1;
var focus_values = {};

// Helper functions
function post(url, data, callback) {
  var xmlhttp = new XMLHttpRequest();
  xmlhttp.open("POST", API_ROOT + url, true);

  xmlhttp.onreadystatechange = function() {
    if (xmlhttp.readyState == 4) {
      // Nice reasons why errors occur
      if (xmlhttp.status >= 500) {
        callback("API server error");
      } else if (xmlhttp.status == 404) {
        callback("Incorrect server configuration");
      } else {
        callback(xmlhttp.responseText);
      }
    }
  };

  xmlhttp.setRequestHeader("Content-Type", "application/json");
  xmlhttp.send(JSON.stringify(data));
}

function get(url, callback) {
  var xmlhttp = new XMLHttpRequest();
  xmlhttp.open("GET", API_ROOT + url, true);

  xmlhttp.onreadystatechange = function() {
    if (xmlhttp.readyState == 4) {
      // Nice reasons why errors occur
      if (xmlhttp.status >= 500) {
        callback("API server error");
      } else if (xmlhttp.status == 404) {
        callback("Incorrect server configuration");
      } else {
        callback(xmlhttp.responseText);
      }
    }
  };

  xmlhttp.send();
}

let daysOfWeek = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"];

// date stuff - very cursed code that probably will break for other tz's tbh
function serverDateToBrowser(server_date) {
  var now = new Date(server_date);
  now.setMinutes(now.getMinutes() - now.getTimezoneOffset());
  return now.toISOString().slice(0,16);
}

function serverDateToPretty(server_date) {
  var now = new Date(server_date);
  return daysOfWeek[now.getDay()] + " " + now.toLocaleString();
}

function tomorrowAtNineAm() {
  var tomorrow = new Date();
  tomorrow.setDate(tomorrow.getDate() + 1);
  tomorrow.setHours(9);
  tomorrow.setMinutes(-tomorrow.getTimezoneOffset());
  return tomorrow.toISOString().slice(0,16);
}

function browserDateToServer(browser_date) {
  return new Date(browser_date).toISOString();
}

const MS_IN_A_DAY = 1000*60*60*24;
const MS_IN_A_WEEK = MS_IN_A_DAY*7;
function prettyPrintDue(server_date) {
  let now = new Date();
  let due_date = new Date(server_date);
  let time_delta = due_date - now;

  let timePrintFunc = function(dt) {
    return dt.toLocaleTimeString()
              .split("").reverse().join("")
              .replace(":00", "")
              .split("").reverse().join("");
  };

  if (time_delta < MS_IN_A_WEEK) {
    // due within a week
    if (time_delta < 0) {
      // due in past, say the date
      if ((now.getDate() == due_date.getDate()) && 
          (now.getMonth() == due_date.getMonth()) &&
          (now.getYear() == due_date.getYear())) {
        return timePrintFunc(due_date) + " today";
      } else {
        return (due_date.getMonth()+1) + "/" + due_date.getDate();
      }
    } else {
      // due in future, say day of week + date
      if (time_delta < MS_IN_A_DAY) {
        if (now.getDate() == due_date.getDate()) {
          // removes the :00 in the seconds
          return timePrintFunc(due_date) + " today";
        } else {
          return timePrintFunc(due_date) + " tomorrow";
        }
      }
      return timePrintFunc(due_date) + " " + daysOfWeek[due_date.getDay()];
    }
  } else {
    // too far in the future
    return false;
  }
}

// modal stuff
var modal_showing = null;
function showModal(name) {
  if (modal_showing) {
    hideModal(true);
  }

  modal_showing = "modal-" + name;
  document.getElementById("modal-backdrop").style.display = "block";
  document.getElementById(modal_showing).style.display = "block";
  document.body.style.overflowY = "hidden";
}

function hideModal(dont_hide_backdrop) {
  // hide the view
  document.getElementById(modal_showing).style.display = "none";

  // hide the backdrop
  if (!dont_hide_backdrop) {
    document.getElementById("modal-backdrop").style.display = "none";
    document.body.style.overflowY = "auto";
  }

  // set variable
  modal_showing = null;
}

function registerModalCloses() {
  var closes = document.getElementsByClassName("modal-closer");
  for (var i = 0; i < closes.length; i++) {
    closes[i].addEventListener('click', function(e) {
      hideModal();
      e.preventDefault();
    }, false);
  }
}

function notify(msg, isBad) {
  var theID = Math.floor(Math.random() * 10000);
  document.getElementById("notification-queue").innerHTML += '<div class="notification'+(isBad?' notification-bad':'')+'" id="notification-' + theID +'">'+msg+'</div>';
  // Remove the elenebt after 3 seconds
  setTimeout(function() {
    var element = document.getElementById("notification-"+theID);
    element.parentNode.removeChild(element);
  }, 3000);
}

// Own functions
function genTagColors() {
  var possibleColors = ["#07457E", "#FFC914", "#7C077E", "#17BEBB", "#F17300", "#76B041", "#8C0000"];
  var curColors = possibleColors.slice();
  var q;

  // load existing colours
  var oldColours = localStorage.getItem(LOCALSTORAGE_KEYS.TAGS);
  if (oldColours) {
    var oldColors = JSON.parse(oldColours);
    for (var i = 0; i < oldColors.length; i++) {
      for (var j = 0; j < tags.length; j++) {
        if (oldColors[i].id == tags[j].id) {
          // restore the existing color
          tags[j].color = oldColors[i].color;
          q = curColors.indexOf(tags[j].color);
          if (q > -1) curColors.splice(q, 1);
          break;
        }
      }
    }
  }

  // only generate for those tags that don't already have colours
  for (var i = 0; i < tags.length; i++) {
    if (!tags[i].hasOwnProperty("color")) {
      // Ran out of colors
      if (curColors.length == 0) curColors = possibleColors.slice();
      tags[i].color = curColors[Math.floor(Math.random()*curColors.length)];
      q = curColors.indexOf(tags[i].color);
      if (q > -1) curColors.splice(q, 1);
    }
  }

  // store the results
  localStorage.setItem(LOCALSTORAGE_KEYS.TAGS, JSON.stringify(tags));
}

function tagToColor(tag) {
  for (var i = 0; i < tags.length; i++) {
    if (tag == tags[i].id) return tags[i].color;
  }
  return "#000";
}

function todoLinkHook(e) {
  var elem = e.target;

  // fetch info about this todo
  focus_id = parseInt(elem.dataset.id);
  infoTodo();
  e.preventDefault();
}
function hookTodos() {
  var todoElems = document.getElementsByClassName("todo-item");

  for (var i = 0; i < todoElems.length; i++) {
    todoElems[i].addEventListener('click', todoLinkHook, false);
    todoElems[i].addEventListener('dragstart', function(e) {
      const dt = event.dataTransfer;
      let id = parseInt(e.target.dataset.id);
      dt.setData("application/x.todoapp", e.target.dataset.id);
      dt.setData("text/plain", e.target.dataset.id);
      // dt.setData("text/plain", "");
    });
  }
}
var selected = [];
function updateFilter() {
  var strs = ["", "", "", "", ""];
  for (var i = 0; i < todos.length; i++) {
    // first check if this todo is selected
    if (selected.indexOf(todos[i].tag_id) < 0) continue;

    // it is, append the data
    strs[todos[i].state] += "<li style=\"color: " + tagToColor(todos[i].tag_id) + "\" ";
    strs[todos[i].state] += "class=\"todo-item\" data-id=\"" + todos[i].id + "\" draggable=\"true\">" + todos[i].name;
    let dueStr = prettyPrintDue(todos[i].due_date);
    if (dueStr) {
      strs[todos[i].state] += "<div class=\"due-date\"\" data-id=\"" + todos[i].id + "\">(due " + dueStr + ")</div>";
    }
    strs[todos[i].state] += "</li>";
  }
  for (var i = 1; i < 5; i++) {
    document.getElementById("todos-"+i).innerHTML = strs[i];
  }
  setTimeout(hookTodos, 50);
}

function tagLinkHook(e) {
  var elem = e.target;

  // toggle color
  var temp = elem.style.backgroundColor;
  elem.style.backgroundColor = elem.style.color;
  elem.style.color = temp;

  // add to or remove from selected
  var val = parseInt(elem.dataset.value);
  var index = selected.indexOf(val);
  if (index > -1) {
    selected.splice(index, 1);
  } else {
    selected.push(val);
  }
  updateFilter();
  e.preventDefault();
}
function hookTags() {
  var tagElems = document.getElementsByClassName("tag-list-item");

  for (var i = 0; i < tagElems.length; i++) {
    tagElems[i].addEventListener('click', tagLinkHook, false);
  }
}
function syncTags() {
  var str = "";
  var str2 = "";

  genTagColors();

  for (var i = 0; i < tags.length; i++) {
    selected.push(tags[i].id);
    str += "<option value=\"" + tags[i].id + "\" style=\"color: " + tags[i].color + ";\">" + tags[i].name + "</option>";
    str2 += "<li class=\"tag-list-item\" data-value=\"" + tags[i].id + "\" style=\"background-color: " + tags[i].color + "; border: 1px solid " + tags[i].color + "; color: #fff;\">" + tags[i].name + "</li>";
  }

  document.getElementById("me-tagid").innerHTML = str;
  document.getElementById("mgmnt-tags").innerHTML = str2;
  setTimeout(hookTags, 50);
}

function loginOk() {
  document.getElementById("login-password").value = "";
  document.getElementById("mgmnt-panel-username").innerHTML = localStorage.getItem(LOCALSTORAGE_KEYS.USERNAME);
  document.getElementById("login-panel").style.display = "none";
  document.getElementById("mgmnt-panel").style.display = "block";
  document.getElementById("md-edit").style.display = "inline-block";

  updateTodos();
}
function logoutOk() {
  document.getElementById("login-panel").style.display = "block";
  document.getElementById("mgmnt-panel").style.display = "none";
  document.getElementById("md-edit").style.display = "none";
  updateTodos();
}

function fetchTags() {
  get("/tags/list", function (text) {
    try {
      var json = JSON.parse(text);
      if (json.error) {
        notify("Failed to fetch tags: " + json.error, true);
      } else {
        tags = json.tags ? json.tags : [];
        syncTags();
        updateFilter();
      }
    } catch (e) {
      notify("Failed to fetch tags: " + text, true);
    }
  });
}

function login() {
  post("/token/new", {
    type: 1,
    username: document.getElementById("login-username").value,
    password: document.getElementById("login-password").value
  }, function (text) {
    try {
      var json = JSON.parse(text);
      if (json.error) {
        // Error
        notify("Failed to authenticate: " + json.error, true);
      } else {
        // All good!
        localStorage.setItem(LOCALSTORAGE_KEYS.TOKEN, json.token);
        localStorage.setItem(LOCALSTORAGE_KEYS.USERNAME, document.getElementById("login-username").value);
        loginOk();
      }
    } catch (e) {
      notify("Failed to authenticate: " + text, true);
    }
  });
}

function logout() {
  post("/token/invalidate", {
    token: localStorage.getItem(LOCALSTORAGE_KEYS.TOKEN),
    authority: localStorage.getItem(LOCALSTORAGE_KEYS.TOKEN)
  }, function (text) {
    try {
      var json = JSON.parse(text);
      if (json.error) {
        notify("Failed to sign off: " + json.error, true);
      } else {
        // All good!
        localStorage.removeItem(LOCALSTORAGE_KEYS.TOKEN);
        localStorage.removeItem(LOCALSTORAGE_KEYS.USERNAME);
        logoutOk();
      }
    } catch (e) {
      notify("Failed to sign off: " + text, true);
    }
  });
}

function checkLogin() {
  if (!localStorage.getItem(LOCALSTORAGE_KEYS.TOKEN)) return;
  post("/token/type", {
    token: localStorage.getItem(LOCALSTORAGE_KEYS.TOKEN)
  }, function (text) {
    try {
      var json = JSON.parse(text);
      if (json.type == 9) {
        // Session invalid
      } else {
        // All good!
        loginOk();
      }
    } catch (e) {
      
    }
  });
}

function invalidateToken() {
  post("/token/invalidate", {
    token: document.getElementById("mt-tokenvalue").value,
    authority: localStorage.getItem(LOCALSTORAGE_KEYS.TOKEN)
  }, function (text) {
    try {
      var json = JSON.parse(text);
      if (json.error) {
        notify("Failed to invalidate token: " + json.error, true);
      } else {
        // All good!
        notify("Token invalidated");
      }
    } catch (e) {
      notify("Failed to invalidate token: " + text, true);
    }
  });
}

function createToken() {
  post("/token/new", {
    type: parseInt(document.getElementById("mt-tokentype").value),
    authority: localStorage.getItem(LOCALSTORAGE_KEYS.TOKEN)
  }, function (text) {
    try {
      var json = JSON.parse(text);
      if (json.error) {
        notify("Failed to create token: " + json.error, true);
      } else {
        // All good!
        notify("New token created");
        document.getElementById("mt-tokenvalue").value = json.token;
      }
    } catch (e) {
      notify("Failed to invalidate token: " + text, true);
    }
  });
}

function updateTodos() {
  var obj = {};
  var token = localStorage.getItem(LOCALSTORAGE_KEYS.TOKEN);
  if (token) {
    obj.authority = token;
  }
  post("/todos/list", obj, function (text) {
    try {
      var json = JSON.parse(text);
      if (json.todos) {
        todos = json.todos;
      } else {
        // Empty array
        todos = [];
      }
      updateFilter();
    } catch (e) {
      notify("Failed to fetch list of todos: " + text, true);
    }
  });
}

function updateTodo() {
  post("/todo/update", {
    todo: {
      id: focus_id,
      state: parseInt(document.getElementById("me-state").value),
      tag_id: parseInt(document.getElementById("me-tagid").value),
      public: document.getElementById("me-public").value == "yes",
      name: document.getElementById("me-name").value,
      due_date: browserDateToServer(document.getElementById("me-duedate").value),
      description: document.getElementById("me-description").value,
    },
    authority: localStorage.getItem(LOCALSTORAGE_KEYS.TOKEN),
  }, function(text) {
    try {
      var json = JSON.parse(text);
      if (json.error) {
        notify("Failed to " + (focus_id == -1 ? "create" : "update") + " todo: " + json.error, true);
      } else {
        notify("Successfully " + (focus_id == -1 ? "created" : "updated") + " todo");
        // Hide modal by default
        hideModal();
        updateTodos();
      }
    } catch (e) {
      notify("Failed to " + (focus_id == -1 ? "create" : "update") + " todo: " + text, true);
    }
  });
}

function deleteTodo() {
  post("/todo/remove", {
    todo: {
      id: focus_id,
    },
    authority: localStorage.getItem(LOCALSTORAGE_KEYS.TOKEN),
  }, function(text) {
    try {
      var json = JSON.parse(text);
      if (json.error) {
        notify("Failed to delete todo: " + json.error, true);
      } else {
        notify("Successfully deleted todo");
        hideModal();
        updateTodos();
      }
    } catch (e) {
      notify("Failed to delete todo: " + text, true);
    }
  });
}

function infoTodo() {
  var obj = {};
  var token = localStorage.getItem(LOCALSTORAGE_KEYS.TOKEN);
  if (token) {
    obj.authority = token;
  }
  obj.todo = {id: focus_id};
  post("/todo/info", obj, function(text) {
    try {
      var json = JSON.parse(text);
      if (json.error) {
        notify("Failed to fetch information for todo: " + json.error, true);
      } else {
        focus_values = json.todo;
        document.getElementById("md-name").innerHTML = focus_values.name;
        document.getElementById("md-duedate").innerHTML = serverDateToPretty(focus_values.due_date);
        document.getElementById("md-desc").innerHTML = converter.makeHtml(focus_values.description);
        showModal("detailedtodo");
      }
    } catch (e) {
      notify("Failed to fetch information for todo: " + text, true);
    }
  });
}

function startEditingTodo(is_new) {
  showModal("edittodo");
  document.getElementById("me-name").value = focus_values.name;
  document.getElementById("me-duedate").value = serverDateToBrowser(focus_values.due_date);
  document.getElementById("me-description").value = focus_values.description;
  document.getElementById("me-state").selectedIndex = focus_values.state - 1;
  document.getElementById("me-tagid").selectedIndex = focus_values.tag_id - 1;
  document.getElementById("me-public").checked = focus_values.public;

  setTimeout(function () {
    this.focus();
    this.setSelectionRange(this.value.length, this.value.length);
  }.bind(document.getElementById('me-description')), 10);
}

// Init functions
function registerUIButtons() {
  // Login button + enter press
  document.getElementById("login-btn").addEventListener('click', function(e) {
    login();
    e.preventDefault();
  }, false);
  document.getElementById("login-password").addEventListener('keydown', function (e) {
    if (e.which == 13) {
      login();
    }
  });

  // Token management button
  document.getElementById("mgmnt-token").addEventListener('click', function(e) {
    showModal("token");
    e.preventDefault();
  }, false);

  // New todo button
  document.getElementById("mgmnt-newtodo").addEventListener('click', function(e) {
    focus_id = -1;
    document.getElementById("me-name").value = "";
    document.getElementById("me-duedate").value = tomorrowAtNineAm();
    document.getElementById("me-description").value = "";
    document.getElementById("me-state").selectedIndex = "0";
    document.getElementById("me-tagid").selectedIndex = "0";
    document.getElementById("me-public").checked = false;
    showModal("edittodo");
    setTimeout(function () {
      this.focus();
      this.setSelectionRange(this.value.length, this.value.length);
    }.bind(document.getElementById('me-name')), 10);
    e.preventDefault();
  }, false);

  // Logout button
  document.getElementById("mgmnt-logout").addEventListener('click', function(e) {
    logout();
    e.preventDefault();
  }, false);

  // Modal - detailed todo - edit todo
  document.getElementById("md-edit").addEventListener('click', function (e) {
    startEditingTodo();
    e.preventDefault();
  }, false);

  // Modal - token management - invalidate token
  document.getElementById("mt-invalidate").addEventListener('click', function (e) {
    invalidateToken();
    e.preventDefault();
  }, false);
  // Modal - token management - create token
  document.getElementById("mt-create").addEventListener('click', function (e) {
    createToken();
    e.preventDefault();
  }, false);

  // Modal - edit todo - delete todo
  document.getElementById("me-delete").addEventListener('click', function (e) {
    deleteTodo();
    e.preventDefault();
  }, false);
  // Modal - edit todo - save changes
  document.getElementById("me-save").addEventListener('click', function (e) {
    updateTodo();
    e.preventDefault();
  }, false);
}



(function() {
  registerModalCloses();
  registerUIButtons();
  checkLogin();
  updateTodos();
  fetchTags();
  // nicer ux
  document.getElementById('modal-backdrop').addEventListener('click', function (e) {
    if (e.target == this) {
      e.preventDefault();
      hideModal();
    }
  });
  document.addEventListener('keypress', function (e) {
    if (e.target == document.body) {
      if (modal_showing) {
        if (modal_showing == 'modal-detailedtodo') {
          if (e.key == 'e') {
            startEditingTodo();
            e.preventDefault();
          }
        }
        if (e.key == 'Escape') {
          e.preventDefault();
          hideModal();
        }
      }
    } else if (modal_showing == 'modal-edittodo') {
      if (e.metaKey || e.ctrlKey) {
        if (e.key == 's') {
          updateTodo();
          e.preventDefault();
        } else if (e.key == 'x') {
          hideModal();
        }
      } else if (e.key == 'Escape') {
        e.preventDefault();
      }
    }
    if (e.key == 'Escape') {
      e.preventDefault();
    }
  });
  document.addEventListener("dragover", function(e) {
    e.preventDefault();
  });
  document.addEventListener('drop', function(e) {
    e.preventDefault();
    var data = e.dataTransfer.getData("application/x.todoapp");

    let yMin = document.getElementById("todos-container").getBoundingClientRect().top;
    if ((yMin > e.clientY) || (!data)) {
      // too high up
      return;
    }

    // figure out which box the user dragged into
    let destList = 1;
    for (var i = 1; i < 5; i++) {
      let rect = document.getElementById("todos-" + i).getBoundingClientRect();
      if (e.clientX > rect.x) {
        // further right than we are, so we keep going
        destList = i;
      } else {
        // to the left of us, so we know it's us now
        break;
      }
    }

    post("/todo/info", {
      authority: localStorage.getItem(LOCALSTORAGE_KEYS.TOKEN),
      todo: {id: parseInt(data)},
    }, function(text) {
      try {
        var json = JSON.parse(text);
        if (json.error) {
          notify("Failed to fetch information for todo: " + json.error, true);
        } else {
          json.todo.state = parseInt(destList);
          post("/todo/update", {
            todo: json.todo,
            authority: localStorage.getItem(LOCALSTORAGE_KEYS.TOKEN),
          }, function(text) {
            try {
              var json = JSON.parse(text);
              if (json.error) {
                notify("Failed to update todo: " + json.error, true);
              } else {
                notify("Successfully updated todo");
                updateTodos();
              }
            } catch (e) {
              notify("Failed to update todo: " + text, true);
            }
          });
        }
      } catch (e) {
        notify("Failed to fetch information for todo: " + text, true);
      }
    });
  });
})();

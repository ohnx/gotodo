/* global localStorage */

// LocalStorage helper
const LOCALSTORAGE_KEYS = {
  TOKEN: 1,
  USERNAME: 2
};
const API_ROOT = "/api"

// Global variables
var todos = [];
var tags = [];
var editing_todo_id;

// Helper functions
function post(url, data, callback) {
  var xmlhttp = new XMLHttpRequest();
  xmlhttp.open("POST", API_ROOT + url, true);

  xmlhttp.onreadystatechange = function() {
    if (xmlhttp.readyState == 4) {
      // Ready for stuff
      callback(xmlhttp.responseText);
    }
  };

  xmlhttp.setRequestHeader("Content-Type", "application/json");
  xmlhttp.send(JSON.stringify(data));
}

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
    closes[i].addEventListener('click', function() {
      hideModal();
    });
  }
}

function hookLink(callback) {
  
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
function loginOk() {
  document.getElementById("login-password").value = "";
  document.getElementById("mgmnt-panel-username").value = localStorage.getItem(LOCALSTORAGE_KEYS.USERNAME);
  document.getElementById("login-panel").style.display = "none";
  document.getElementById("mgmnt-panel").style.display = "block";
}
function logoutOk() {
  document.getElementById("login-panel").style.display = "block";
  document.getElementById("mgmnt-panel").style.display = "none";
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
      notify("Failed to authenticate: " + e, true);
    }
  });
}

function logout() {
  post("/type/invalidate", {
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
      notify("Failed to sign off: " + e, true);
    }
  });
}

function check_login() {
  if (!localStorage.getItem(LOCALSTORAGE_KEYS.TOKEN)) return;
  post("/token/type", {
    token: localStorage.getItem(LOCALSTORAGE_KEYS.TOKEN)
  }, function (text) {
    try {
      var json = JSON.parse(text);
      if (json.type != 9) {
        // Session invalid
      } else {
        // All good!
        loginOk();
      }
    } catch (e) {
      
    }
  });
}

// Init functions
function registerUIButtons() {
  // Login button + enter press
  document.getElementById("login-btn").addEventListener('click', function() {
    login();
  });
  document.getElementById("login-password").addEventListener('keydown', function (e) {
    if (e.which == 13) {
      login();
    }
  });

  // Token management button
  document.getElementById("mgmnt-token").addEventListener('click', function() {
    showModal("token");
  });

  // New todo button
  document.getElementById("mgmnt-newtodo").addEventListener('click', function() {
    editing_todo_id = 0;
    showModal("edittodo");
  });

  // Logout button
  document.getElementById("mgmnt-logout").addEventListener('click', function() {
    logout();
  });
}

(function() {
  registerModalCloses();
  registerUIButtons();
  check_login();
})();

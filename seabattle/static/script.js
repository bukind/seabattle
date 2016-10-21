"use strict";

function onBodyLoaded() {
  setMisteryReactions();
}

function setMisteryReactions() {
  var tds;
  var i;

  tds = document.getElementsByClassName("mist");

  for (i = 0; i < tds.length; i++) {
    setMisteryCell(tds[i])
  }
}

function setMisteryCell(mist) {
  var kids;
  var i;
  var attrs;
  var j;
  var msg;

  mist.onmouseover = misteryCellOn;
  mist.onmouseout = misteryCellOff;
  kids = mist.children;
  for (i = 0; i < kids.length; i++) {
    if (kids[i].tagName === "A") {
      if (mist.id === "pA10") {
        attrs = kids[i].attributes;
        msg = "Cell " + mist.id + " anchor attributes:<ul>\n";
        for (j = 0; j < attrs.length; j++) {
          msg += "<li>" + attrs[j].name + ": " + attrs[j].value + "</li>\n";
        }
        msg += "</ul>";
        showMessage(msg);
      }
      kids[i].setAttribute("x-href", kids[i].getAttribute("href"));
      kids[i].removeAttribute("href");
      mist.onclick = misteryCellOnClick;
      break;
    }
  }
}

function misteryCellOnClick() {
  showMessage("cell " + this.id + " pressed");
}

function misteryCellOn() {
  this.classList.add("undercursor");
  // showMessage("cell " + this.id + ":" + this.className);
}

function misteryCellOff() {
  this.classList.remove("undercursor");
}

function showMessage(msg) {
  console = document.getElementById("messages");
  console.innerHTML += "<li>" + msg + "</li>\n";
}

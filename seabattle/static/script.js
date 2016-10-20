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

  mist.onmouseover = misteryCellOn;
  mist.onmouseout = misteryCellOff;
  kids = mist.children;
  for (i = 0; i < kids.length; i++) {
    if (kids[i].tagName = "a") {
      mist.onclick = kids[i].onclick;
      kids[i].onclick = function() {};
      break;
    }
  }
}

function misteryCellOn() {
  this.classList.add("undercursor");
  // showMessage("cell " + this.id + ":" + this.className)
}

function misteryCellOff() {
  this.classList.remove("undercursor");
}

function showMessage(msg) {
  console = document.getElementById("messages");
  console.innerHTML += "<li>" + msg + "</li>\n";
}

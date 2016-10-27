"use strict";

// constants
var tick = 0;
var height = 30;
var width = 16;
var tickInterval = 100;
var tickTimer = null;    // the timer to stop at the end of the game
var falling = null;      // the current falling piece
var cells = [];
var colors = ["", "black", "red", "green", "blue"];

// helper function to convert position (i,j) into the cell id.
function mkId(i, j) {
  return "" + i.toString() + "-" + j.toString();
}

// when the body is loaded - initialization.
function onBodyLoaded() {
  var well;
  var mesh = "";
  var i, j;

  // create a well
  well = document.getElementById("well-holder");
  mesh += "<table id=\"well\" class=\"well\">\n";
  for (i = 0; i < height; i++) {
    mesh += "<tr>\n";
    for (j = 0; j < width; j++) {
      mesh += " <td id=\"" + mkId(i,j) + "\"></td>\n";
    }
    mesh += "</tr>\n";
  }
  mesh += "</table>\n";
  well.innerHTML = mesh;

  // fill the cells array
  for (i = 0; i < height; i++) {
    var row;
    row = [];
    for (j = 0; j < width; j++) {
      row.push(0);
    }
    cells.push(row)
  }

  tickTimer = window.setInterval(onTick, tickInterval);
}

function showPiece(show) {
  var td = document.getElementById(mkId(falling.posy, falling.posx));
  if (show) {
  	td.classList.add(colors[falling.color]);
  } else {
    td.classList.remove(colors[falling.color]);
  }
}

// executed on the tick
function onTick() {
  tick += 1;

  var p = document.getElementById("tick");
  p.innerHTML = tick.toString();

  if (falling === null) {
    // generate a new piece
    var pos, color;
    pos = Math.floor(Math.random() * width);
    if (cells[0][pos] != 0) {
      // this is already occupied, the well is full.
      window.clearInterval(tickTimer);
      return;
    }

    color = Math.floor(Math.random() * (colors.length - 1) + 1);

    falling = {posx: pos, posy: 0, color: color};
    showPiece(true);
    return;
  }

  // try to move a piece
  var nextposy = falling.posy + 1;
  if (nextposy < height) {
    // check if it hits the previous items ...
    if (cells[nextposy][falling.posx] === 0) {
      // ... no, so move it.
      showPiece(false);
      falling.posy = nextposy;
      showPiece(true);
      return;
    }
  }

  // ... yes, the falling piece hits the obstacle.
  // freeze it where it is now.
  cells[falling.posy][falling.posx] = falling.color;
  falling = null;
}

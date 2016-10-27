"use strict";

// constants
var tick = 0;
var linecount = 0;
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
  var p;    // element counting ticks
  var posy; // position on y of the falling block
  var i;
  var row;

  tick += 1;

  p = document.getElementById("tick");
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
  posy = falling.posy + 1;
  if (posy < height) {
    // check if it hits the previous items ...
    if (cells[posy][falling.posx] === 0) {
      // ... no, so move it.
      showPiece(false);
      falling.posy = posy;
      showPiece(true);
      return;
    }
  }

  // ... yes, the falling piece hits the obstacle.
  // freeze it where it is now.
  posy = falling.posy;
  cells[posy][falling.posx] = falling.color;
  falling = null;

  // check if the line is complete.
  for (i = 0; i < width; i++) {
    if (cells[posy][i] === 0) {
      // not fully filled
      return;
    }
  }

  // the whole line is filled, cut it out
  cells.splice(posy,1);
  row = [];
  for (i = 0; i < width; i++) {
    row.push(0);
  }
  cells.unshift(row);

  // redraw all elements above and including posy.
  for (i = posy; i >= 0; i--) {
    var j;
    for (j = 0; j < width; j++) {
      var td = document.getElementById(mkId(i,j));
      td.className = colors[cells[i][j]];
    }
  }

  // increment the linecount and show it.
  linecount += 1;
  document.getElementById("linecount").innerHTML = linecount.toString();
}

"use strict";

// when the body is loaded - initialization.
function onBodyLoaded() {
  var game = makeGame()
  game.run()
}


// create a game object
function makeGame() {
  var height = 30;
  var width = 16;
  var colors = ["", "black", "red", "green", "blue"];
  var i, j;

  var game = {
    // --- fields
    tick: 0,           // arbitrary time
    cells: [],         // the playing field
    falling: null,     // the falling piece or null
    tickInterval: 100; // interval b/w ticks, ms
    linecount: 0,      // number of full lines collected
    colors: colors,    // possible colors of pieces
    tickTimer: null,   // the ticking timer (to be init in run)

    // --- methods
    run: runGame,      // run a game
    onTick: onTickHandler,
    onKeyUp: onKeyUpHandler,
  };

  // initialize the game field
  for (i = 0; i < height; i++) {
    var row;
    row = [];
    for (j = 0; j < width; j++) {
      row.push(0);
    }
    game.cells.push(row)
  }

  return game;
}


// run the game.
// initialize the game and start the timer.
function runGame() {
  domShowField(this.cells.length, this.cells[0].length);
  this.tickTimer = window.setInterval(this.onTick, this.tickInterval);
  window.onkeyup = this.onKeyUp;
}


// show initial field.
function domShowField(height, width) {
  var well;
  var mesh = "";
  var i, j;

  // create a well
  well = document.getElementById("well-holder");
  mesh += "<table id=\"well\" class=\"well\">\n";
  for (i = 0; i < height; i++) {
    mesh += "<tr>\n";
    for (j = 0; j < width; j++) {
      mesh += " <td id=\"" + mkCellId(i,j) + "\"></td>\n";
    }
    mesh += "</tr>\n";
  }
  mesh += "</table>\n";
  well.innerHTML = mesh;
}


function onKeyUpHandler(e) {
  var key = e.keyCode ? e.keyCode : e.which;
  if (key == 37) { // leftkey
    // try to move left
  } else if (key == 39) { // rightkey
    // try to move right
  } else if (key == 65) { // A
  } else if (key == 68) { // D
  } else if (key == 32) { // space
    // drop it
    this.dropPiece(this.falling);
  }
}


// executed on the tick
function onTickHandler() {
  var p;    // element counting ticks
  var posy; // position on y of the falling block

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

  if (tick > 100) {
    dropPiece();
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
  stopPiece(falling.posy);
}

function stopPiece(posy) {
  var i;
  var row;

  // freeze the piece
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


// helper function to convert position (i,j) into the cell id.
function mkCellId(i, j) {
  return "" + i.toString() + "-" + j.toString();
}




function showPiece(show) {
  var td = document.getElementById(mkId(falling.posy, falling.posx));
  if (show) {
  	td.classList.add(colors[falling.color]);
  } else {
    td.classList.remove(colors[falling.color]);
  }
}

function dropPiece() {
  var posy;
  var i;
  if (falling === null) {
    return;
  }
  posy = falling.posy;
  for (i = falling.posy+1; i < height; i++) {
    if (cells[i][falling.posx] != 0) {
      break;
    }
    posy = i;
  }
  showPiece(false);
  falling.posy = posy;
  showPiece(true);

  stopPiece(posy);
}


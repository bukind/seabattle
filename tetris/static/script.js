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

  var game = {
    // --- fields
    tick: 0,           // arbitrary time
    cells: [],         // the playing field
    falling: null,     // the falling piece or null
    tickInterval: 100, // interval b/w ticks, ms
    linecount: 0,      // number of full lines collected
    tickTimer: null,   // the ticking timer (to be init in run)

    // --- methods
    width: function() { return this.cells[0].length; },
    height: function() { return this.cells.length; },
    run: null,      // run a game
  };

  // initialize the game field
  {
    var i, j;
    for (i = 0; i < height; i++) {
      var row;
      row = [];
      for (j = 0; j < width; j++) {
        row.push(0);
      }
      game.cells.push(row)
    }
  }

  var domLog = function(msg) {
    var td = document.getElementById("log")
    td.innerHTML += msg + "<br/>\n";
  }

  // helper function to convert cell position into the id.
  var domCellId = function(y, x) {
    return "" + y.toString() + "-" + x.toString();
  }

  // show the initial field
  var domMakeField = function(height, width) {
    var well;
    var mesh = "";
    var i, j;

    // create a well
    well = document.getElementById("well-holder");
    mesh += "<table id=\"well\" class=\"well\">\n";
    for (i = 0; i < height; i++) {
      mesh += "<tr>\n";
      for (j = 0; j < width; j++) {
        mesh += " <td id=\"" + domCellId(i,j) + "\"></td>\n";
      }
      mesh += "</tr>\n";
    }
    mesh += "</table>\n";
    well.innerHTML = mesh;
  }

  var domShowField = function(height, width, cells) {
    var i, j;
    for (i = 0; i < height; i++) {
      for (j = 0; j < width; j++) {
        var td = document.getElementById(domCellId(i,j));
        td.className = colors[cells[i][j]];
      }
    }
  }

  // show the tick value
  var domShowTick = function(tick) {
    var p = document.getElementById("tick");
    p.innerHTML = tick.toString();
  }

  var domShowLineCount = function(linecount) {
    document.getElementById("linecount").innerHTML = linecount.toString();
  }

  var domShowPiece = function(piece, show) {
    var td = document.getElementById(domCellId(piece.posy, piece.posx));
    if (show) {
      td.classList.add(colors[piece.color]);
    } else {
      td.classList.remove(colors[piece.color]);
    }
  }

  var domRedrawPiece = function(oldpiece, newpiece) {
    domShowPiece(oldpiece, false);
    domShowPiece(newpiece, true);
  }

  // the reaction on the key pressing
  var onKeyUp = function(e) {
    var key = e.keyCode ? e.keyCode : e.which;
    if (key == 37) { // leftkey
      // try to move left
    } else if (key == 39) { // rightkey
      // try to move right
    } else if (key == 65) { // A
    } else if (key == 68) { // D
    } else if (key == 32) { // space
      // drop it
      game.dropFallingPiece();
    }
  }

  var makeRandInt = function(min, max) {
    return Math.floor(Math.random() * (max-min)) + min;
  }

  // generate the falling piece
  game.genFallingPiece = function() {
    var posx = makeRandInt(0, this.width());
    var color = makeRandInt(1, colors.length);
    this.falling = {
      posy: 0,
      posx: posx,
      color: color,
      clone: function() {
        var np = {
          posx: this.posx,
          posy: this.posy,
          color: this.color,
          clone: this.clone,
        };
        return np;
      },
    };
  }

  // check one cell in the board
  game.isCellEmpty = function(y,x) {
    return this.cells[y][x] === 0;
  }

  // check if a piece can fit into the board
  game.canFitPiece = function() {
    var piece = this.falling;
    return (piece.posy < this.height() &&
      piece.posx >= 0 && piece.posx < this.width() &&
      this.isCellEmpty(piece.posy, piece.posx));
  }

  // end the game
  game.endGame = function() {
    window.clearInterval(this.tickTimer);
  }

  game.slideFallingPiece = function() {
    this.falling.posy += 1;
    if (!this.canFitPiece()) {
      this.falling.posy -= 1;
      return false;
    }
    return true;
  }

  game.dropFallingPiece = function() {
    if (this.falling === null) {
      return;
    }
    var oldpiece = this.falling.clone();
    for (i = 0; i < height; i++) {
      if (!this.slideFallingPiece()) {
        break;
      }
    }
    domRedrawPiece(oldpiece, this.falling);
    this.stopFallingPiece();
  }

  game.stopFallingPiece = function() {
    var i;
    var posy = this.falling.posy;

    // freeze the piece
    this.cells[posy][this.falling.posx] = this.falling.color;
    this.falling = null;

    // check if the line is complete.
    for (i = 0; i < this.width(); i++) {
      if (this.isCellEmpty(posy,i)) {
        return;
      }
    }

    // cut the complete line
    this.cells.splice(posy, 1);
    var row = [];
    for (i = 0; i < width; i++) {
      row.push(0);
    }
    this.cells.unshift(row);

    domShowField(posy+1, width, this.cells);

    this.linecount += 1;
    domShowLineCount(this.linecount);
  }

  // action performed periodically
  var onTick = function() {
    game.tick += 1;
    domShowTick(game.tick);

    if (game.falling === null) {
      game.genFallingPiece();
      if (!game.canFitPiece()) {
        // failed to fit
        game.falling = null;
        game.endGame();
        return;
      }
      domShowPiece(game.falling, true);
      return;
    }

    if (game.tick > 100) {
      game.dropFallingPiece();
      return;
    }

    var oldpiece = game.falling.clone();
    if (game.slideFallingPiece()) {
      domRedrawPiece(oldpiece, game.falling);
      return;
    }
  
    // freeze it where it is now.
    game.stopFallingPiece();
  }

  // running the game - field initial show + setting events.
  game.run = function() {
    domMakeField(this.height(), this.width());
    this.tickTimer = window.setInterval(onTick, this.tickInterval);
    window.onkeyup = onKeyUp;
  }

  return game;
}

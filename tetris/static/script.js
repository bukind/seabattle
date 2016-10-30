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
  var pieces = [
    [[0,-1], [0,0], [0,1], [0,2]], // line
    [[0,-1], [0,0], [1,0], [1,1]], // z
    [[1,-1], [1,0], [0,0], [0,1]], // s
    [[0,0], [0,1], [1,0], [1,1]],  // block
    [[0,-1], [0,0], [0,1], [1,0]], // taur
  ];

  var game = {
    // --- fields
    tick: 0,           // arbitrary time
    cells: [],         // the playing field
    falling: null,     // the falling piece or null
    tickInterval: 100, // interval b/w ticks, ms
    linecount: 0,      // number of full lines collected
    tickTimer: null,   // the ticking timer (to be init in run)

    // --- methods
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
  var domMakeField = function() {
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
    var xy = piece.getXY();
    var i;
    for (i = 0; i < xy.length; i++) {
      if (xy[i][0] < 0 || xy[i][0] >= width) {
        continue;
      }
      if (xy[i][1] < 0 || xy[i][1] >= height) {
        continue;
      }
      var td = document.getElementById(domCellId(xy[i][1], xy[i][0]));
      if (show) {
        td.classList.add(colors[piece.color]);
      } else {
        td.classList.remove(colors[piece.color]);
      }
    }
  }

  var domRedrawPiece = function(oldpiece, newpiece) {
    domShowPiece(oldpiece, false);
    domShowPiece(newpiece, true);
  }

  // the reaction on the key pressing
  var onKeyUp = function(e) {
    var key = e.keyCode ? e.keyCode : e.which;

    if (game.falling === null) {
      return;
    }

    if (key === 37 || key === 39) { // left/right
      var oldpiece = game.falling.clone();
      game.falling.posx += key - 38;
      if (!game.canFitPiece()) {
        game.falling = oldpiece;
      } else {
        domRedrawPiece(oldpiece, game.falling);
      }
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
    var posx = makeRandInt(0, width);
    var color = makeRandInt(1, colors.length);
    var pn = makeRandInt(0, pieces.length);
    this.falling = {
      posy: 0,
      posx: posx,
      color: color,
      piece: pieces[pn],
      getXY: function() {
        var xy = [];
        var i;
        for (i = 0; i < this.piece.length; i++) {
          xy.push([this.posx+this.piece[i][0],
                   this.posy+this.piece[i][1]])
        }
        return xy;
      },
      clone: function() {
        var np = {
          posx: this.posx,
          posy: this.posy,
          color: this.color,
          piece: this.piece,
          getXY: this.getXY,
          clone: this.clone,
        };
        return np;
      },
    };
  }

  // check one cell in the board
  game.isCellEmpty = function(y,x) {
    if (x < 0 || x >= width) { return false; }
    if (y < 0) { return true; }
    if (y >= height) { return false; }
    return this.cells[y][x] === 0;
  }

  // check if a piece can fit into the board
  game.canFitPiece = function() {
    var piece = this.falling.getXY();
    var i;
    var minx = piece[0][0];
    var maxx = maxx;
    // try to fix cases where x < 0 or x >= width
    for (i = 1; i < piece.length; i++) {
      var x = piece[i][0];
      if (x < minx) {
        minx = x;
      } else if (x > maxx) {
        maxx = x;
      }
    }
    if (minx < 0) {
      this.falling.posx -= minx;
      piece = this.falling.getXY();
    } else if (maxx >= width) {
      this.falling.posx -= (maxx - width) + 1;
      piece = this.falling.getXY();
    }
    for (i = 0; i < piece.length; i++) {
      if (!this.isCellEmpty(piece[i][1], piece[i][0])) {
        return false;
      }
    }
    return true;
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
    var j;
    var xy = this.falling.getXY();
    var miny = xy[0][1];
    var maxy = miny;

    // freeze the piece
    for (i = 0; i < xy.length; i++) {
      var y = xy[i][1];
      this.cells[y][xy[i][0]] = this.falling.color;
      if (y < miny) {
        miny = y;
      }
      if (y > maxy) {
        maxy = y;
      }
    }
    this.falling = null;

    // check if the line is complete.
    var ycomplete = [];
    for (i = miny; i <= maxy; i++) {
      var ok = true;
      for (j = 0; j < width; j++) {
        if (this.isCellEmpty(i,j)) {
          ok = false;
          break;
        }
      }
      if (ok) {
        ycomplete.push(i);
      }
    }

    if (ycomplete.length == 0) {
      return;
    }

    // cut the complete line
    for (i = ycomplete.length-1; i >= 0; i--) {
      this.cells.splice(ycomplete[i], 1);
    }
    for (j = 0; j < ycomplete.length; j++) {
      var row = [];
      for (i = 0; i < width; i++) {
        row.push(0);
      }
      this.cells.unshift(row);
    }

    domShowField(ycomplete[ycomplete.length-1]+1, width, this.cells);
    this.linecount += ycomplete.length;

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

    /*
    if (game.tick > 100) {
      game.dropFallingPiece();
      return;
    }
     */

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
    domMakeField(height, width);
    this.tickTimer = window.setInterval(onTick, this.tickInterval);
    window.onkeyup = onKeyUp;
  }

  return game;
}

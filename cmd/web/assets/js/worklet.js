registerPaint('worklet', class {
  static get inputProperties() {
    return ['--clr','--accent']
  }
  ctx;
  pw;
  ph;
  clr;
  accent;
  paint(ctx, {width, height}, props) {
    this.ctx = ctx;
    const gridCellSize = 50;
    ctx.lineWidth = 3;
    this.ph = n => (n*height)/100;
    this.pw = n => (n*width)/100;
    this.clr = props.get('--clr').toString();
    this.accent = props.get('--accent').toString();

    //dot matrix
    Array.from(Array(25)).forEach((_,i) => {
      Array.from(Array(9)).forEach((_,j) => {
        ctx.fillStyle = this.accent;
        ctx.arc(this.ph((i+1)*2) + 5 , this.pw((j+1)*2) + 5, 2, 0, 2 * Math.PI);
        ctx.closePath();
        ctx.fill();
      })
    })

    //rects
    ctx.strokeStyle = this.clr;
    const rX = 68,
          rY = 65;
    ctx.strokeRect(this.pw(rX), this.ph(rY), this.pw(20), this.pw(15));
    ctx.strokeStyle = this.accent;
    ctx.strokeRect(this.pw(rX-5), this.ph(rY+2), this.pw(15), this.pw(7));

    
  //lines
  // this.ctx.strokeStyle = this.clr;
  // ctx.beginPath();
  // const arcX = this.pw(4), arcY = this.ph(10);
  // ctx.moveTo(arcX, arcY);
  // ctx.quadraticCurveTo(arcX+50, arcY-50, arcX+100, arcY);
  // ctx.quadraticCurveTo(arcX+150, arcY+50, arcX+200, arcY);
  // ctx.quadraticCurveTo(arcX+250, arcY-50, arcX+300, arcY);
  // ctx.stroke();
  ctx.save();
  ctx.rotate(45 * Math.PI / 180);
  this.curves(4,10,2,40);
  this.curves(4,10,6,40);
  this.curves(4,14,10,40);
  ctx.restore();
  }
  curves(n,x,y,gap) {
    this.ctx.strokeStyle = this.clr;
    this.ctx.beginPath();
    let arcX = this.pw(x), 
        arcY = this.ph(y);
    this.ctx.moveTo(arcX, arcY);
    Array.from(Array(n)).forEach((_,i) => {
      const off = (i+1)*gap,
            to = ((i+1)*2)*gap,
            from = to-gap,
            dir = ((i%2)?gap:(-gap));
            console.log(`iter n${i+1}`,from,to,dir)
      this.ctx.quadraticCurveTo(
        arcX+from,
        arcY+dir, 
        arcX+to, 
        arcY
      );
      this.ctx.stroke();
    })
  }
})
// from x y
// to x +100 y - 100 arc x+100 arcy = y
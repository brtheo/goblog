import "./schemeSwitcher";
import "./TOC";

if (CSS.paintWorklet) 
  CSS.paintWorklet.addModule("/static/js/worklet.js");
        

window.addEventListener('load', () => {
  
  const customProps =  [
    {name: '--clr', syntax: '<color>', inherits: true, initialValue: '#cacaca'},
    {name: '--accent', syntax: '<color>', inherits: true, initialValue: '#779d46'},
  ];
  // customProps.forEach(CSS.registerProperty);
  // document.documentElement.addEventListener('mousemove', (evt) => {
  //   requestAnimationFrame(() => {
  //     el.style.setProperty('--mx', evt.x);
  //     el.style.setProperty('--my', evt.y);
  //   });
  // });
})
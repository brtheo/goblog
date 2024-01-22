import {LitElement, css, unsafeCSS} from 'lit'
import {customElement} from 'lit/decorators.js'
import {html, unsafeStatic} from 'lit/static-html.js'

@customElement('table-content')
export class TOC extends LitElement {
  static styles = [
    css`:host {
      height: max-content;
      position: sticky !important;
      top: 10px;
      width: max-content;
      display: flex;
      flex-direction: column;
      padding: .8rem !important;
      gap: .2rem !important;
    }
    [heading] {
      font-family: 'sofia-pro'; 
      margin-block: .3rem;
    }
    a {
      text-decoration: none;
      transition: color 100ms linear; 
    }
    a:hover {
      color: var(--accent) !important; 
    }`,
    Array(6).fill(0).map((_,i) => Object({
      h: i+1,
      fs: `${(8 - i)/10}rem`
    })).map(obj => unsafeCSS(`
      h${obj.h} {
        font-size: ${obj.fs};
        margin-inline-start: ${(obj.h - 1) * .5}rem;
      }
    `))
  ];

  headingsMapped;
  headings;
  order;
  makeObserver() {
    const getPrev = n => Array(n+1).fill(0);
    
    const observer = new IntersectionObserver(entries => {
      entries.forEach(e => {  
        const idx = this.order[e.target.id]  
        if(!e.isIntersecting) this.style.setProperty(`--heading-${idx}`, 'var(--clr)');
        else {
          getPrev(idx).forEach((_,i) => 
            this.style.setProperty(`--heading-${i}`, 'var(--accent)')
          )
        }
      })
    }, {
      rootMargin: "50px",
      threshold: 0,
    });

    this.headings.forEach(h => 
      observer.observe(h as HTMLElement)
    )
  }
  makeHeadings() {
    const makeHash = s => "#"+s.toLowerCase().replaceAll(/\s[!?]/g,'').replaceAll(' ','-')
    const article = document.querySelector('article')!;
    this.headings = article.querySelectorAll('h1, h2, h3, h4, h5, h6')
    this.order = Object.fromEntries([...this.headings].map((h,i) => [h.id, i]))
    this.headingsMapped = [...this.headings].map((e,i) => html`
      <${unsafeStatic(e.tagName)} heading>
        <a class="az" style="color: var(--heading-${i}, var(--clr))" href=${makeHash(e.innerHTML)}>${e.innerHTML}</a>
      </${unsafeStatic(e.tagName)}>
    `)
  }
  connectedCallback(): void {
  super.connectedCallback();
  this.makeHeadings()
  this.makeObserver()
  }
  
  
  render() {
    return html`${this.headingsMapped ?? ''}`
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'table-content': TOC
  }
}

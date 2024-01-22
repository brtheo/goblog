import {LitElement, html, css} from 'lit'
import {customElement, property, state} from 'lit/decorators.js'

@customElement('scheme-switcher')
export class SchemeSwitcher extends LitElement {
  static styles = css`
  :host {
    display: grid;
    place-content: center;
    cursor: pointer;
    width: 100%;
  }
    button {
      all: unset;
      font-size: 2rem;
      filter: grayscale(.7);
    }
  `
  connectedCallback() {
    super.connectedCallback();
    if(globalThis.sessionStorage.getItem(this.storageKeyName)) {
      !Boolean(+globalThis.sessionStorage.getItem(this.storageKeyName)!)
        ? document.documentElement.classList.remove(this.cssVarName)
        : document.documentElement.classList.add(this.cssVarName)
    }
    else if(globalThis.matchMedia("(prefers-color-scheme: dark)").matches) {
      document.documentElement.classList.add(this.cssVarName)
      globalThis.sessionStorage.setItem(this.storageKeyName,'1')
    }
    else 
      globalThis.sessionStorage.setItem(this.storageKeyName,'0')

    this.mode = Boolean(+globalThis.sessionStorage.getItem(this.storageKeyName)!);
  }
  switchDarkmode() {
    document.documentElement.classList.toggle(this.cssVarName)
    !Boolean(+globalThis.sessionStorage.getItem(this.storageKeyName)!) 
      ? globalThis.sessionStorage.setItem(this.storageKeyName,'1')
      : globalThis.sessionStorage.setItem(this.storageKeyName, '0')
    this.mode = Boolean(+globalThis.sessionStorage.getItem(this.storageKeyName)!);
  }
  @state() mode: boolean;
  @property({type: String, attribute: 'css-var-name'}) cssVarName: string = "dark";
  @property({type: String, attribute: 'storage-key-name'}) storageKeyName: string = this.cssVarName;
  render() {
    
    return html`
      <button @click=${this.switchDarkmode}>${this.mode ? html`🌑` : html`☀️`}</button>
    `
  }
}

declare global {
  interface HTMLElementTagNameMap {
    'scheme-switcher': SchemeSwitcher
  }
}
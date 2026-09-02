/* Zachowania brane dotąd z `design/zasoby/prototyp.js`: powiadomienie, widoki,
   zakładki, menu i okno modalne. Reszta tamtej warstwy nie wchodzi. */

import { godzinaZycia, KSZTALTY_ZNAKU, type RodzajOgloszenia } from '../model/powiadomienie.ts';

interface StykOkna {
  dnToast?: (tytul: string, tresc: string, rodzaj: string, ms: number) => void;
  dnPrzelaczWidok?: (nazwa: string, grupa: string | null) => void;
}

const PRZESTRZEN_ZNAKU = 'http://www.w3.org/2000/svg';

function pojemnikPowiadomien(): HTMLElement {
  const stojacy = document.querySelector<HTMLElement>('.dn-toasty');
  if (stojacy !== null) return stojacy;
  const pojemnik = document.createElement('div');
  pojemnik.className = 'dn-toasty';
  pojemnik.setAttribute('role', 'status');
  pojemnik.setAttribute('aria-live', 'polite');
  document.body.appendChild(pojemnik);
  return pojemnik;
}

/* Znak stoi bezpośrednio w powiadomieniu, bo `komponenty.css` barwi go regułą
   `.dn-toast > svg`; owinięty w element pośredni zostałby bez barwy wagi. */
function znakPowiadomienia(rodzaj: RodzajOgloszenia): SVGElement {
  const znak = document.createElementNS(PRZESTRZEN_ZNAKU, 'svg');
  znak.setAttribute('viewBox', '0 0 24 24');
  znak.setAttribute('fill', 'none');
  znak.setAttribute('stroke', 'currentColor');
  znak.setAttribute('stroke-width', '1.75');
  znak.setAttribute('stroke-linecap', 'round');
  znak.setAttribute('stroke-linejoin', 'round');
  znak.setAttribute('aria-hidden', 'true');
  for (const ksztalt of KSZTALTY_ZNAKU[rodzaj]) {
    const element = document.createElementNS(PRZESTRZEN_ZNAKU, ksztalt.element);
    for (const [cecha, wartosc] of Object.entries(ksztalt.cechy)) {
      element.setAttribute(cecha, wartosc);
    }
    znak.appendChild(element);
  }
  return znak;
}

function powiadom(tytul: string, tresc: string, rodzaj: string, ms: number): void {
  const waga: RodzajOgloszenia = rodzaj in KSZTALTY_ZNAKU
    ? (rodzaj as RodzajOgloszenia)
    : 'informacja';
  const wpis = document.createElement('div');
  wpis.className = `dn-toast dn-toast--${waga}`;
  wpis.appendChild(znakPowiadomienia(waga));
  const tresci = document.createElement('div');
  const naglowek = document.createElement('div');
  naglowek.className = 'dn-toast-tytul';
  naglowek.textContent = tytul;
  tresci.appendChild(naglowek);
  if (tresc !== '') {
    const zdanie = document.createElement('div');
    zdanie.className = 'dn-toast-tresc';
    zdanie.textContent = tresc;
    tresci.appendChild(zdanie);
  }
  wpis.appendChild(tresci);
  pojemnikPowiadomien().appendChild(wpis);
  globalThis.setTimeout(() => wpis.remove(), godzinaZycia(ms));
}

function przelaczWidok(nazwa: string, grupa: string | null): void {
  const wyborWidoku = grupa === null
    ? '[data-widok]:not([data-grupa-widoku])'
    : `[data-grupa-widoku="${grupa}"]`;
  for (const widok of document.querySelectorAll<HTMLElement>(wyborWidoku)) {
    const czynny = widok.dataset.widok === nazwa;
    widok.dataset.widokAktywny = czynny ? 'tak' : 'nie';
    if (czynny) {
      widok.classList.remove('pt-wejscie');
      void widok.offsetWidth;
      widok.classList.add('pt-wejscie');
    }
  }
  const wyborPrzelacznika = grupa === null
    ? '[data-przelacz]:not([data-grupa])'
    : `[data-przelacz][data-grupa="${grupa}"]`;
  for (const przelacznik of document.querySelectorAll<HTMLElement>(wyborPrzelacznika)) {
    const czynny = przelacznik.dataset.przelacz === nazwa;
    przelacznik.setAttribute('aria-selected', String(czynny));
    if (przelacznik.hasAttribute('aria-current') || przelacznik.classList.contains('dn-boczna-pozycja')) {
      if (czynny) przelacznik.setAttribute('aria-current', 'page');
      else przelacznik.removeAttribute('aria-current');
    }
    przelacznik.classList.toggle(
      'dn-zakladka--wybrana',
      czynny && przelacznik.classList.contains('dn-zakladka'),
    );
  }
}

function ustawWidokiPoczatkowe(): void {
  const grupy = new Map<string, HTMLElement[]>();
  for (const widok of document.querySelectorAll<HTMLElement>('[data-widok]')) {
    const nazwa = widok.dataset.grupaWidoku ?? '';
    const wykaz = grupy.get(nazwa) ?? [];
    wykaz.push(widok);
    grupy.set(nazwa, wykaz);
  }
  for (const wykaz of grupy.values()) {
    if (wykaz.some((widok) => widok.dataset.widokAktywny === 'tak')) continue;
    if (wykaz.length > 0) wykaz[0].dataset.widokAktywny = 'tak';
  }
}

/* Zakładki chodzą strzałkami: pas zakładek jest jednym przystankiem tabulacji,
   więc bez tej obsługi klawiatura sięgnęłaby wyłącznie zakładki wybranej. */
function przestawZakladke(zakladka: HTMLElement, klawisz: string): boolean {
  const pas = zakladka.closest('[role="tablist"]');
  if (pas === null) return false;
  const zakladki = [...pas.querySelectorAll<HTMLElement>('[role="tab"]')];
  const numer = zakladki.indexOf(zakladka);
  if (numer < 0) return false;
  let cel = -1;
  if (klawisz === 'ArrowRight') cel = (numer + 1) % zakladki.length;
  if (klawisz === 'ArrowLeft') cel = (numer - 1 + zakladki.length) % zakladki.length;
  if (klawisz === 'Home') cel = 0;
  if (klawisz === 'End') cel = zakladki.length - 1;
  if (cel < 0) return false;
  zakladki[cel].click();
  zakladki[cel].focus();
  return true;
}

function oznaczPorzadekZakladek(): void {
  for (const zakladka of document.querySelectorAll<HTMLElement>('[role="tablist"] [role="tab"]')) {
    zakladka.tabIndex = zakladka.getAttribute('aria-selected') === 'true' ? 0 : -1;
  }
}

function zamknijMenu(pozaTym: string): void {
  for (const panel of document.querySelectorAll<HTMLElement>('[data-menu-tresc][data-otwarte="tak"]')) {
    if (panel.id === pozaTym) continue;
    panel.dataset.otwarte = 'nie';
    document.querySelector(`[data-menu="${panel.id}"]`)?.setAttribute('aria-expanded', 'false');
  }
}

// Położenie rozwiniętego panelu liczy `zasoby/menu.js` po zmianie `data-otwarte`.
function przelaczMenu(przycisk: HTMLElement): void {
  const oznaczenie = przycisk.dataset.menu ?? '';
  const panel = document.getElementById(oznaczenie);
  const nazwa = przycisk.dataset.etykietka ?? przycisk.getAttribute('aria-label') ?? 'Lista';
  if (panel === null) {
    powiadom(nazwa, 'Lista wskazana tym przyciskiem nie stoi w oknie.', 'ostrzezenie', 0);
    return;
  }
  if (panel.querySelectorAll('.sta-menu-poz, [role="menuitem"]').length === 0) {
    powiadom(nazwa, panel.dataset.pustyKomunikat ?? 'Brak pozycji do pokazania.', 'informacja', 0);
    return;
  }
  const otwarte = panel.dataset.otwarte === 'tak';
  panel.dataset.otwarte = otwarte ? 'nie' : 'tak';
  przycisk.setAttribute('aria-expanded', String(!otwarte));
}

const styk = globalThis as StykOkna;
styk.dnToast = powiadom;
styk.dnPrzelaczWidok = przelaczWidok;

// Faza przechwytywania: wiązania okien zatrzymują zdarzenia na pozycjach menu.
document.addEventListener('click', (zdarzenie) => {
  const cel = zdarzenie.target;
  if (!(cel instanceof Element)) return;
  const przycisk = cel.closest<HTMLElement>('[data-menu]');
  zamknijMenu(przycisk?.dataset.menu ?? '');
  if (przycisk !== null) przelaczMenu(przycisk);
}, true);

document.addEventListener('click', (zdarzenie) => {
  const cel = zdarzenie.target;
  if (!(cel instanceof Element)) return;
  const modal = cel.closest<HTMLElement>('[data-otworz-modal]')?.dataset.otworzModal;
  if (modal !== undefined) document.querySelector<HTMLDialogElement>(`#${modal}`)?.showModal();
  const przelacznik = cel.closest<HTMLElement>('[data-przelacz]');
  if (przelacznik === null) return;
  przelaczWidok(przelacznik.dataset.przelacz ?? '', przelacznik.dataset.grupa ?? null);
});

document.addEventListener('keydown', (zdarzenie) => {
  const cel = zdarzenie.target;
  if (!(cel instanceof Element)) return;
  const zakladka = cel.closest<HTMLElement>('[role="tab"]');
  if (zakladka === null) return;
  if (przestawZakladke(zakladka, zdarzenie.key)) zdarzenie.preventDefault();
});

function postaw(): void {
  ustawWidokiPoczatkowe();
  oznaczPorzadekZakladek();
}

if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', postaw, { once: true });
} else {
  postaw();
}

// Wiązanie karty okna roboczego z oknem komunikacji rdzenia wraz z wpisami
// wspólnymi wiązaniom modułów. Znacznik niesie biblioteka Właściciela.
import { Command } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import { przypiszOknoKomunikacji } from './okna-robocze.ts';

export function zwiazOkno(
  kanal: Kanal,
  _kodModulu: string,
  nazwaSrodowiska: string,
  korzen: ParentNode,
): void {
  void opiszNaglowek(kanal, nazwaSrodowiska, korzen);
}

export function zwiazOknoStojace(
  kanal: Kanal,
  idKarty: string,
  nazwaSrodowiska: string,
  idOkna: string,
  korzen: ParentNode,
): void {
  if (!przypiszOknoKomunikacji(idKarty, idOkna)) {
    oglos('Okno modułu', 'Karta zeszła z pasma, zanim okno rdzenia do niej doszło — '
      + 'zamknięcie karty nie zamknie tego okna.', 'ostrzezenie');
  }
  void opiszNaglowek(kanal, nazwaSrodowiska, korzen);
}

/* Nagłówek prototypu niesie środowisko, model i wysiłek wpisane wprost. Pole
   o znanej wartości dostaje ją z rdzenia, pole bez pokrycia schodzi — inaczej
   okno w sesji TalkIn twierdzi, że stoi w WorkSpace. */
export async function opiszNaglowek(
  kanal: Kanal,
  nazwaSrodowiska: string,
  korzen: ParentNode,
): Promise<void> {
  const naglowek = korzen.querySelector('.sta-kom-naglowek');
  if (naglowek === null) return;
  const wynik = await wywolaj(kanal, Command.ChannelList, { enabledOnly: true });
  const kanalModelu = wynik.udany ? wynik.wynik?.channels[0] : undefined;
  const wartosci: Record<string, string> = {
    'Środowisko': nazwaSrodowiska,
    Model: kanalModelu?.model ?? kanalModelu?.name ?? '',
  };
  for (const pole of [...naglowek.querySelectorAll('.sta-kom-pole')]) {
    const podpis = Object.keys(wartosci).find(
      (kandydat) => pole.textContent?.startsWith(kandydat) === true,
    );
    wpiszPole([pole], podpis ?? '', podpis === undefined ? '' : wartosci[podpis]);
  }
}

export function wpiszPole(pola: Element[], podpis: string, wartosc: string): void {
  const pole = podpis === ''
    ? pola[0]
    : pola.find((kandydat) => kandydat.textContent?.startsWith(podpis) === true);
  if (pole === undefined) return;
  if (wartosc === '') {
    pole.remove();
    return;
  }
  // Pole środowiska w prototypie nie niesie klasy `.dane`, tylko wyróżnienie.
  const dane = pole.querySelector('.dane') ?? pole.querySelector('b');
  if (dane !== null) dane.textContent = wartosc;
}

export function wpiszTekst(wezel: Element, tekst: string): boolean {
  for (const dziecko of wezel.childNodes) {
    if (dziecko.nodeType !== Node.TEXT_NODE) continue;
    if ((dziecko.nodeValue ?? '').trim() === '') continue;
    dziecko.nodeValue = tekst;
    return true;
  }
  return false;
}

export function opiszWezel(wezel: Element | null, wartosc: string): void {
  if (wezel === null) return;
  const teksty = [...wezel.childNodes].filter((dziecko) => dziecko.nodeType === Node.TEXT_NODE);
  const cel = teksty.find((dziecko) => (dziecko.nodeValue ?? '').trim() !== '')
    ?? teksty[teksty.length - 1];
  if (cel === undefined) return;
  cel.nodeValue = wartosc;
}

export function sklonuj(wezel: Element | null): HTMLElement | null {
  return wezel instanceof HTMLElement ? (wezel.cloneNode(true) as HTMLElement) : null;
}

export function godzina(znacznik: number): string {
  return new Date(znacznik).toLocaleTimeString('pl-PL', { hour: '2-digit', minute: '2-digit' });
}

export function data(znacznik: number): string {
  return new Date(znacznik).toLocaleDateString('pl-PL');
}

export function zdejmijSterowanieWspolne(cialo: HTMLElement): void {
  // Rejestr uprawnień nie zna domyślności, którą opisuje podpis piątego trybu.
  cialo.querySelector('#pop-tryb-upr .sta-popover-wiersz small')?.remove();
}

// Panel bez pokrycia w rdzeniu zostaje: układ okna jest kompozycją Właściciela.
export function uzgodnijPrzelacznikiPaneli(cialo: HTMLElement, _bezPokrycia: string[]): void {
  for (const przelacznik of cialo.querySelectorAll<HTMLElement>('[data-panel-toggle]')) {
    const panel = cialo.querySelector(`#${przelacznik.dataset.panelToggle ?? ''}`);
    przelacznik.setAttribute('aria-checked', String(panel?.hasAttribute('hidden') === false));
  }
}

export function zdejmijTrescWspolna(cialo: HTMLElement): void {
  const historia = cialo.querySelector('.sta-kom-historia');
  historia?.replaceChildren();
  /* Wiadomość wysyła wyłącznie Studio, a Studio tej ścieżki nie wywołuje.
     Puste pole nadawania obok wybranego modelu obiecuje odpowiedź, której
     nikt nie da, więc miejsce rozmowy mówi wprost, jak jest. */
  if (historia !== null) {
    const zdanie = cialo.ownerDocument.createElement('p');
    zdanie.className = 'dn-tekst-ciagly dn-tekst-ciagly--drobny';
    zdanie.dataset.pustka = 'tak';
    zdanie.textContent = 'Rozmowa tego modułu nie jest jeszcze związana z rdzeniem.';
    historia.appendChild(zdanie);
  }
  cialo.querySelector('.sta-kom-monitor .sta-kom-monitor-tresc')?.replaceChildren();
  /* Żetony kontekstu i znacznik pracy prototyp wpisuje wprost — nazwy plików,
     adresy stron, „pracuje”. Bez zdjęcia stoją obok prawdziwych zer wykazu
     jako dwa stany wykluczające się, oba wzięte znikąd. */
  cialo.querySelector('.sta-kom-kontekst')?.replaceChildren();
  /* Belka panelu rozmowy niesie nazwę pracy wymyśloną na pokaz — raz żetonem,
     raz dopiskiem w cudzysłowie przy tytule panelu. */
  for (const zeton of cialo.querySelectorAll('.sta-kom > .sta-okno-belka .sta-chip')) {
    zeton.remove();
  }
  const tytulRozmowy = cialo.querySelector('.sta-kom > .sta-okno-belka .sta-okno-tytul b');
  if (tytulRozmowy !== null) {
    tytulRozmowy.textContent = (tytulRozmowy.textContent ?? '')
      .replace(/\s*[—–-]?\s*„[^”]*”/u, '').trim();
  }
  cialo.querySelector('.sta-kontekst-akcji')?.replaceChildren();
  cialo.querySelector('.sta-kom-stan')?.remove();
  /* Znaczniki paneli niosą w prototypie liczby wzięte znikąd. Wiązanie panelu,
     które ma czym je wypełnić, robi to po tym zdjęciu. */
  for (const znacznik of cialo.querySelectorAll('.sta-okno-znacznik')) {
    znacznik.textContent = '';
  }
  for (const znacznik of cialo.querySelectorAll('#menu-filtr .sta-menu-poz[aria-checked]')) {
    znacznik.removeAttribute('aria-checked');
  }
}

const MIEJSCA = new WeakMap<Element, { rodzic: Element; przed: Node | null }>();

export function zdejmijAlboPostaw(wezel: Element | null, obecny: boolean): void {
  if (wezel === null) return;
  if (!obecny) {
    if (wezel.parentElement === null) return;
    MIEJSCA.set(wezel, { rodzic: wezel.parentElement, przed: wezel.nextSibling });
    wezel.remove();
    return;
  }
  if (wezel.parentElement !== null) return;
  const miejsce = MIEJSCA.get(wezel);
  miejsce?.rodzic.insertBefore(wezel, miejsce.przed);
}

export function wpiszAlboZdejmij(wezel: Element | null, wartosc: string): void {
  if (wezel === null) return;
  if (wartosc !== '') wezel.textContent = wartosc;
  zdejmijAlboPostaw(wezel, wartosc !== '');
}

export function wpiszTekstAlboZdejmij(wezel: Element | null, wartosc: string): void {
  if (wezel === null) return;
  if (wartosc !== '' && !wpiszTekst(wezel, wartosc)) {
    wezel.remove();
    return;
  }
  zdejmijAlboPostaw(wezel, wartosc !== '');
}

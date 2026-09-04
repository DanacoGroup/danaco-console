// Centrum powiadomień otwierane dzwonkiem szyny narzędziowej: rejestr zdarzeń
// od najnowszego wraz z trzema czynnościami kontraktu — odczytaniem, zamknięciem
// i odłożeniem. Plakietka dzwonka liczy zdarzenia nowe w całym rejestrze, nie
// w widoku, bo tak stanowi opis komendy wykazu.

import { Command, type Notification, type NotificationState } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';

const NAGLOWEK = 'Powiadomienia';
const WYBOR_DZWONKA = '[data-etykietka="Powiadomienia"]';

/** Odłożenie liczone od chwili naciśnięcia; kontrakt żąda chwili powrotu, nie odstępu. */
const ODLOZENIE_MS = 60 * 60 * 1000;

const NAZWY_STANOW: Readonly<Record<NotificationState, string>> = {
  nowe: 'nowe',
  odczytane: 'odczytane',
  obsluzone: 'obsłużone',
  odlozone: 'odłożone',
};

let otwarte: HTMLElement | null = null;

export function zwiazPowiadomienia(kanal: Kanal): void {
  document.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element) || cel.closest(WYBOR_DZWONKA) === null) return;
    zdarzenie.stopPropagation();
    if (otwarte !== null) {
      zamknij();
      return;
    }
    void otworz(kanal);
  }, true);
  void odswiezPlakietke(kanal);
}

function zamknij(): void {
  otwarte?.remove();
  otwarte = null;
}

async function otworz(kanal: Kanal): Promise<void> {
  const wynik = await wywolaj(kanal, Command.NotificationList, {});
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił wykazu zdarzeń.', 'blad');
    return;
  }
  postawPanel(kanal, wynik.wynik.notifications);
  ustawPlakietke(wynik.wynik.unread);
}

async function odswiezPlakietke(kanal: Kanal): Promise<void> {
  const wynik = await wywolaj(kanal, Command.NotificationList, { limit: 1 });
  if (wynik.udany && wynik.wynik !== undefined) ustawPlakietke(wynik.wynik.unread);
}

/* Plakietka siada na dzwonku, więc znika razem z nim, gdy szyny w oknie nie ma. */
function ustawPlakietke(nowych: number): void {
  const dzwonek = document.querySelector<HTMLElement>(WYBOR_DZWONKA);
  if (dzwonek === null) return;
  const stara = dzwonek.querySelector('.dn-plakietka');
  if (nowych <= 0) {
    stara?.remove();
    return;
  }
  const plakietka = stara ?? document.createElement('span');
  plakietka.className = 'dn-plakietka dn-plakietka--informacja';
  plakietka.textContent = String(nowych);
  if (stara === null) dzwonek.append(plakietka);
}

function postawPanel(kanal: Kanal, zdarzenia: readonly Notification[]): void {
  const nakladka = document.createElement('div');
  nakladka.className = 'dn-modal-nakladka';
  nakladka.setAttribute('role', 'dialog');
  nakladka.setAttribute('aria-label', NAGLOWEK);

  const okno = document.createElement('div');
  okno.className = 'dn-modal';
  const belka = document.createElement('div');
  belka.className = 'dn-modal-belka';
  const tytul = document.createElement('b');
  tytul.textContent = NAGLOWEK;
  const zamknijPrzycisk = document.createElement('button');
  zamknijPrzycisk.className = 'dn-btn-ikona';
  zamknijPrzycisk.type = 'button';
  zamknijPrzycisk.setAttribute('aria-label', 'Zamknij');
  zamknijPrzycisk.textContent = '✕';
  belka.append(tytul, zamknijPrzycisk);
  okno.append(belka);

  const obszar = document.createElement('div');
  obszar.className = 'sta-obszar';
  if (zdarzenia.length === 0) obszar.append(zdaniePuste());
  else for (const zdarzenie of zdarzenia) obszar.append(zbudujWiersz(kanal, zdarzenie));
  okno.append(obszar);

  nakladka.append(okno);
  document.body.append(nakladka);
  otwarte = nakladka;

  zamknijPrzycisk.addEventListener('click', zamknij);
  nakladka.addEventListener('click', (zdarzenie) => {
    if (zdarzenie.target === nakladka) zamknij();
  });
}

function zdaniePuste(): HTMLElement {
  const wiersz = document.createElement('p');
  wiersz.textContent = 'Rejestr nie ma zdarzeń otwartych.';
  return wiersz;
}

function zbudujWiersz(kanal: Kanal, zdarzenie: Notification): HTMLElement {
  const wiersz = document.createElement('div');
  wiersz.className = 'dn-obszar-pozycja';
  wiersz.dataset.idZdarzenia = zdarzenie.id;

  const tresc = document.createElement('span');
  tresc.className = 'dn-obszar-pozycja-nazwa';
  tresc.textContent = zdarzenie.text;
  const opis = document.createElement('span');
  opis.className = 'dn-etyk-mono';
  opis.textContent = `${zdarzenie.class} · ${zdarzenie.weight} · ${NAZWY_STANOW[zdarzenie.state]}`;
  wiersz.append(tresc, opis);

  if (zdarzenie.state !== 'obsluzone') {
    wiersz.append(
      czynnosc('Odczytane', () => wywolaj(kanal, Command.NotificationAcknowledge, { ids: [zdarzenie.id] })),
      czynnosc('Zamknij', () => wywolaj(kanal, Command.NotificationResolve, { id: zdarzenie.id })),
      czynnosc('Odłóż o godzinę', () => wywolaj(kanal, Command.NotificationSnooze, {
        id: zdarzenie.id,
        until: Date.now() + ODLOZENIE_MS,
      })),
    );
  }
  return wiersz;
}

function czynnosc(podpis: string, wykonaj: () => Promise<{ udany: boolean; blad?: { message?: string } }>): HTMLElement {
  const przycisk = document.createElement('button');
  przycisk.className = 'dn-btn';
  przycisk.type = 'button';
  przycisk.textContent = podpis;
  przycisk.addEventListener('click', () => {
    void wykonaj().then((wynik) => {
      if (!wynik.udany) {
        oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił wykonania.', 'blad');
        return;
      }
      zamknij();
    });
  });
  return przycisk;
}

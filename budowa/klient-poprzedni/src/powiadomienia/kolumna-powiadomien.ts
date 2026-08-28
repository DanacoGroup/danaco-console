import './powiadomienia.css';

import { NotificationState, type Notification } from '../../../shared/contract';
import { elementIkony } from '../ikony/ikony';
import type { Kanal } from '../protokol/kanal';
import { utworzZrodloCentrum, type ZrodloCentrum } from './zrodlo-centrum';

// Centrum powiadomień to kolumna boczna z rejestrem zdarzeń, zmieniająca wyłącznie stan zdarzenia.

export interface KolumnaPowiadomien {
  element: HTMLElement;
  /** Otwiera kolumnę i odczytuje rejestr. */
  otworz(): void;
  zamknij(): void;
  /** Otwiera albo zamyka — pod wyzwalacz plakietki. */
  przelacz(): void;
  czyOtwarta(): boolean;
  /** Odłącza nasłuch kanału. */
  rozlacz(): void;
}

/** Ile odkłada Odłóż w rejestrze powiadomień — godzina, wprost z konwencji drogi potwierdzenia zdarzenia. */
const ODLOZENIE_MS = 60 * 60 * 1000;

export function utworzKolumnePowiadomien(
  kanal: Kanal,
  naLicznik?: (nowe: number) => void,
): KolumnaPowiadomien {
  const zrodlo: ZrodloCentrum = utworzZrodloCentrum(kanal);

  const element = document.createElement('aside');
  element.className = 'po-kolumna';
  element.hidden = true;
  element.setAttribute('aria-label', 'Centrum powiadomień');

  const tytul = document.createElement('h2');
  tytul.className = 'po-tytul';
  tytul.textContent = 'Centrum powiadomień';

  const oznaczWszystkie = document.createElement('button');
  oznaczWszystkie.type = 'button';
  oznaczWszystkie.className = 'dn-btn dn-btn--zarys po-zbiorcze';
  oznaczWszystkie.textContent = 'Oznacz wszystkie jako odczytane';
  oznaczWszystkie.addEventListener('click', () => {
    void zrodlo.odczytajZdarzenia().then(nanieOdmowe).then(() => odczytaj());
  });

  const zamknij = document.createElement('button');
  zamknij.type = 'button';
  zamknij.className = 'dn-btn-ikona po-zamknij';
  zamknij.setAttribute('aria-label', 'Zamknij centrum powiadomień');
  zamknij.append(elementIkony('zamknij', { rozmiar: 16 }));
  zamknij.addEventListener('click', () => zamknijKolumne());

  const naglowek = document.createElement('header');
  naglowek.className = 'po-naglowek';
  naglowek.append(tytul, zamknij);

  const lista = document.createElement('ul');
  lista.className = 'po-lista';

  const zdanie = document.createElement('p');
  zdanie.className = 'po-zdanie';
  zdanie.hidden = true;

  element.append(naglowek, oznaczWszystkie, zdanie, lista);

  let otwarta = false;

  /** Nanosi zdanie odmowy; puste chowa akapit. */
  function nanieOdmowe(odmowa: string): void {
    zdanie.hidden = odmowa === '';
    zdanie.dataset['powodzenie'] = 'false';
    zdanie.textContent = odmowa;
  }

  function odczytaj(): void {
    void zrodlo.odczytaj().then((stan) => {
      naLicznik?.(stan.nowe);
      if (stan.odmowa !== '') {
        nanieOdmowe(stan.odmowa);
        lista.replaceChildren();
        return;
      }
      nanieOdmowe('');
      lista.replaceChildren(...stan.zdarzenia.map(pozycjaZdarzenia));
      if (stan.zdarzenia.length === 0) lista.append(pustyStan());
    });
  }

  /** Pozycja jednego zdarzenia wraz z jego dwoma czynnościami stanu. */
  function pozycjaZdarzenia(zdarzenie: Notification): HTMLLIElement {
    const pozycja = document.createElement('li');
    pozycja.className = 'dn-karta po-pozycja';
    pozycja.dataset['klasa'] = zdarzenie.class;
    pozycja.dataset['waga'] = zdarzenie.weight;
    pozycja.dataset['stan'] = zdarzenie.state;

    const naglowekPozycji = document.createElement('div');
    naglowekPozycji.className = 'po-pozycja__naglowek';
    const klasa = document.createElement('span');
    klasa.className = 'dn-plakietka po-pozycja__klasa';
    klasa.textContent = zdarzenie.class;
    const czas = document.createElement('time');
    czas.className = 'po-pozycja__czas';
    czas.dateTime = new Date(zdarzenie.createdAt).toISOString();
    czas.textContent = zapisChwili(zdarzenie.createdAt);
    naglowekPozycji.append(klasa, czas);

    const tresc = document.createElement('p');
    tresc.className = 'po-pozycja__tresc';
    tresc.textContent = zdarzenie.text;

    const czynnosci = document.createElement('div');
    czynnosci.className = 'po-pozycja__czynnosci';

    // Zdarzenie już zamknięte nie dostaje czynności: nie ma czego zamykać ani odkładać.
    if (zdarzenie.state !== NotificationState.Obsluzone) {
      czynnosci.append(
        czynnosc('Obsłużone', () =>
          zrodlo.zamknij(zdarzenie.id).then(nanieOdmowe).then(() => odczytaj()),
        ),
        czynnosc('Odłóż o godzinę', () =>
          zrodlo
            .odloz(zdarzenie.id, Date.now() + ODLOZENIE_MS)
            .then(nanieOdmowe)
            .then(() => odczytaj()),
        ),
      );
    }
    if (zdarzenie.state === NotificationState.Nowe) {
      czynnosci.append(
        czynnosc('Odczytane', () =>
          zrodlo
            .odczytajZdarzenia([zdarzenie.id])
            .then(nanieOdmowe)
            .then(() => odczytaj()),
        ),
      );
    }

    pozycja.append(naglowekPozycji, tresc, czynnosci);

    // Działania zdarzenia pokazujemy jako nazwane drogi: wykonuje je rodzina właściwa, nie to okno.
    if (zdarzenie.actions !== undefined && zdarzenie.actions.length > 0) {
      const drogi = document.createElement('p');
      drogi.className = 'dn-pole-opis po-pozycja__drogi';
      drogi.textContent =
        'Działania zdarzenia: ' + zdarzenie.actions.map((czyn) => czyn.label).join(', ') + '.';
      pozycja.append(drogi);
    }
    return pozycja;
  }

  function czynnosc(napis: string, wykonaj: () => Promise<unknown>): HTMLButtonElement {
    const przycisk = document.createElement('button');
    przycisk.type = 'button';
    przycisk.className = 'dn-btn dn-btn--zarys po-czynnosc';
    przycisk.textContent = napis;
    przycisk.addEventListener('click', () => void wykonaj());
    return przycisk;
  }

  function pustyStan(): HTMLLIElement {
    const pusty = document.createElement('li');
    pusty.className = 'po-pusty';
    pusty.textContent =
      'Rejestr jest pusty — nic nie czeka na wiedzę ani na decyzję. O stan procesów ' +
      'zapytaj w oknie komunikacji.';
    return pusty;
  }

  function zamknijKolumne(): void {
    otwarta = false;
    element.hidden = true;
  }

  // Licznik nadąża zdarzeniami także przy zamkniętej kolumnie — po to jest plakietka.
  const odsubskrybuj = zrodlo.naZmiane((nowe) => {
    naLicznik?.(nowe);
    if (otwarta) odczytaj();
  });

  return {
    element,

    otworz() {
      otwarta = true;
      element.hidden = false;
      odczytaj();
    },

    zamknij: zamknijKolumne,

    przelacz() {
      if (otwarta) zamknijKolumne();
      else this.otworz();
    },

    czyOtwarta: () => otwarta,

    rozlacz: () => odsubskrybuj(),
  };
}

/** Godzina dnia bieżącego, data pełna dla dni wcześniejszych, pokazywana przy pozycji rejestru powiadomień. */
function zapisChwili(znacznik: number): string {
  const chwila = new Date(znacznik);
  const dzis = new Date();
  const tenSamDzien =
    chwila.getFullYear() === dzis.getFullYear() &&
    chwila.getMonth() === dzis.getMonth() &&
    chwila.getDate() === dzis.getDate();
  return tenSamDzien
    ? chwila.toLocaleTimeString('pl-PL', { hour: '2-digit', minute: '2-digit' })
    : chwila.toLocaleString('pl-PL', {
        day: 'numeric',
        month: 'short',
        hour: '2-digit',
        minute: '2-digit',
      });
}

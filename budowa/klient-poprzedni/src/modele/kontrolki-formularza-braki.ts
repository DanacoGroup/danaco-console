/**
 * Dokładka kontrolek formularza dostarcza kontrolki „gołe”, bez wiersza
 * z etykietą, przeznaczone do panelu akcji oraz paska narzędzi modułu, gdzie
 * nazwę elementu niesie atrybut opisujący.
 */

import { pokazKomunikat } from '../aplikacja/komunikaty';

/**
 * Przycisk akcji okna — odpowiednik hubowego `przycisk`, ale z klasą domyślną.
 * Osobna nazwa, bo `przycisk(tresc, klasa)` ma własnych odbiorców i zostaje.
 */
export function przyciskAkcji(etykieta: string, klasa = 'dn-btn'): HTMLButtonElement {
  const kontrolka = document.createElement('button');
  kontrolka.type = 'button';
  kontrolka.className = klasa;
  kontrolka.textContent = etykieta;
  return kontrolka;
}

/**
 * Przycisk bez pokrycia w kontrakcie pozostaje widoczny i w pełni klikalny,
 * z powodem podanym wprost: po naciśnięciu nazywa brakującą komendę zamiast
 * być wygaszony albo usunięty.
 */
export function przyciskBezKomendy(etykieta: string, powod: string): HTMLButtonElement {
  const kontrolka = przyciskAkcji(etykieta, 'dn-btn dn-btn--zarys');
  kontrolka.title = powod;
  kontrolka.dataset['brakKomendy'] = 'tak';
  kontrolka.setAttribute('aria-description', powod);
  kontrolka.addEventListener('click', () => {
    pokazKomunikat({ tytul: `${etykieta} — bez komendy w kontrakcie`, tresc: powod, waga: 'ostrz' });
  });
  return kontrolka;
}

/** Pole jednowierszowe formularza, którego nazwę niesie atrybut opisujący, ponieważ etykiety nad polem nie ma w tym miejscu układu. */
export function pole(etykieta: string, podpowiedz = ''): HTMLInputElement {
  const kontrolka = document.createElement('input');
  kontrolka.type = 'text';
  kontrolka.className = 'dn-pole-kontrolka';
  kontrolka.placeholder = podpowiedz;
  kontrolka.setAttribute('aria-label', etykieta);
  return kontrolka;
}

/** Pole liczbowe kontrolki służy do wprowadzenia wartości takich jak priorytet zlecenia, granica wykazu albo próg obiegów. */
export function poleLiczbowe(etykieta: string, podpowiedz = ''): HTMLInputElement {
  const kontrolka = pole(etykieta, podpowiedz);
  kontrolka.type = 'number';
  return kontrolka;
}

/**
 * Pole wielowierszowe — treść kroku, ładunek komendy, warunek, uzasadnienie.
 * Liczba wierszy jest obowiązkowa i stoi na drugim miejscu: moduły potrzebują
 * tu różnych wysokości, a wartość domyślna wspólna dla wszystkich zmieniałaby
 * układ okna po cichu.
 */
export function poleTresci(
  etykieta: string,
  wiersze: number,
  podpowiedz = '',
  klasaDodatkowa = '',
): HTMLTextAreaElement {
  const kontrolka = document.createElement('textarea');
  kontrolka.className =
    klasaDodatkowa === '' ? 'dn-pole-kontrolka' : `dn-pole-kontrolka ${klasaDodatkowa}`;
  kontrolka.rows = wiersze;
  kontrolka.placeholder = podpowiedz;
  kontrolka.setAttribute('aria-label', etykieta);
  return kontrolka;
}

/** Lista wyboru zbudowana z podanych par wartości i odpowiadającej im etykiety wyświetlanej w kontrolce. */
export function wybor(
  etykieta: string,
  pozycje: ReadonlyArray<readonly [string, string]>,
): HTMLSelectElement {
  const kontrolka = document.createElement('select');
  kontrolka.className = 'dn-wybor';
  kontrolka.setAttribute('aria-label', etykieta);
  for (const [wartosc, opis] of pozycje) {
    const pozycja = document.createElement('option');
    pozycja.value = wartosc;
    pozycja.textContent = opis;
    kontrolka.append(pozycja);
  }
  return kontrolka;
}

/** Przełącznik dwustanowy zapisywany do rdzenia, oznaczający stan „automatyka czynna” w bieżącej konfiguracji modułu. */
export function przelacznik(etykieta: string): HTMLInputElement {
  const kontrolka = document.createElement('input');
  kontrolka.type = 'checkbox';
  kontrolka.className = 'dn-przelacznik';
  kontrolka.setAttribute('aria-label', etykieta);
  return kontrolka;
}

/**
 * Przełącznik widoku — czynność wyłącznie kliencka (zawijanie, znaczniki czasu),
 * więc nośnikiem jest przycisk ze stanem w `aria-pressed`: nic nie idzie do rdzenia.
 */
export function przelacznikWidoku(etykieta: string, wlaczony: boolean): HTMLButtonElement {
  const kontrolka = przyciskAkcji(etykieta, 'dn-btn dn-btn--zarys');
  kontrolka.setAttribute('aria-pressed', String(wlaczony));
  kontrolka.dataset['wlaczony'] = String(wlaczony);
  return kontrolka;
}

/** Przestawia przełącznik widoku na stan przeciwny do stanu bieżącego i zwraca jego nowy stan po dokonanej zmianie. */
export function przestaw(kontrolka: HTMLButtonElement): boolean {
  const nowy = kontrolka.dataset['wlaczony'] !== 'true';
  kontrolka.dataset['wlaczony'] = String(nowy);
  kontrolka.setAttribute('aria-pressed', String(nowy));
  return nowy;
}

/** Wykaz pozycji, którego klasa rodziny modułu przychodzi z zewnątrz i określa wygląd właściwy dla wywołującego modułu. */
export function wykaz(etykieta: string, klasa: string): HTMLUListElement {
  const element = document.createElement('ul');
  element.className = klasa;
  element.setAttribute('aria-label', etykieta);
  return element;
}

/**
 * Pozycja wykazu: tytuł, opis i miejsce na przyciski czynności. `przedrostek`
 * to prefiks rodziny modułu (`da`, `dt`, `dw`, `dm`) — z niego powstaje komplet
 * klas BEM pozycji; arkusz rodziny należy do modułu.
 */
export function pozycjaWykazu(
  tytul: string,
  opis: string,
  przedrostek: string,
): { element: HTMLElement; akcje: HTMLElement } {
  const naglowek = document.createElement('strong');
  naglowek.className = `${przedrostek}-pozycja__tytul`;
  naglowek.textContent = tytul;

  const tresc = document.createElement('span');
  tresc.className = `${przedrostek}-pozycja__opis`;
  tresc.textContent = opis;

  const akcje = document.createElement('span');
  akcje.className = `${przedrostek}-pozycja__akcje`;

  const element = document.createElement('li');
  element.className = `${przedrostek}-pozycja`;
  element.append(naglowek, tresc, akcje);
  return { element, akcje };
}

/**
 * Wiersz pola: etykieta nad kontrolką, opcjonalne zdanie objaśniające. Klasa
 * i objaśnienie idą jednym zapisem, a nie dwoma napisami z rzędu, żeby nie dało
 * się zamienić ich miejscami w wywołaniu.
 */
export function wiersz(
  etykieta: string,
  kontrolka: HTMLElement,
  opcje: { klasa: string; objasnienie?: string },
): HTMLElement {
  const podpis = document.createElement('span');
  podpis.className = 'dn-pole-etykieta';
  podpis.textContent = etykieta;

  const element = document.createElement('label');
  element.className = opcje.klasa;
  element.append(podpis, kontrolka);
  if (opcje.objasnienie !== undefined && opcje.objasnienie !== '') {
    const opis = document.createElement('span');
    opis.className = 'dn-pole-opis';
    opis.textContent = opcje.objasnienie;
    element.append(opis);
  }
  return element;
}

/**
 * Pobranie treści jako pliku. Eksport definicji, transkryptu czy migawki nie ma
 * komendy kontraktu i mieć jej nie musi: treść jest już w oknie.
 */
export function pobierzPlik(nazwa: string, tresc: string, rodzaj = 'application/json'): void {
  const adres = URL.createObjectURL(new Blob([tresc], { type: `${rodzaj};charset=utf-8` }));
  const odnosnik = document.createElement('a');
  odnosnik.href = adres;
  odnosnik.download = nazwa;
  odnosnik.click();
  URL.revokeObjectURL(adres);
}

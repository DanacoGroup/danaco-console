/**
 * Kontrolki formularzy sekcji modeli — jeden kształt pola dla całej sekcji.
 *
 * Formularz konta i edytor tożsamości opisują byty o polach **stałych**,
 * wynikających wprost z kontraktu (`AccountAddRequest`, `IdentityDocumentSet`),
 * a nie z katalogu wierszy. Nie mogą więc korzystać z generatora pól okna
 * konfiguracji, który buduje kontrolkę z `SettingDefinition` — nie ma czego mu
 * podać. Zamiast dwóch równoległych sposobów budowania pola w dwóch plikach
 * sekcja ma jeden ten.
 *
 * Wygląd bierzemy w całości z biblioteki `komponenty/` (`dn-pole*`), więc plik
 * nie zna ani jednej barwy i ani jednego odstępu.
 */
export * from './kontrolki-formularza-braki'; // kontrolki okien operacyjnych

/** Pole formularza: etykieta, kontrolka, opcjonalny opis pod nią. */
export interface PoleFormularza<T extends HTMLElement> {
  /** Wiersz osadzany w formularzu. */
  element: HTMLElement;
  /** Kontrolka wewnątrz wiersza. */
  kontrolka: T;
}

/** Opis pola: etykieta, podpowiedź w kontrolce i zdanie wyjaśniające. */
export interface OpisPola {
  etykieta: string;
  podpowiedz?: string;
  opis?: string;
}

/** Identyfikator wiążący etykietę z kontrolką; unikalny w obrębie dokumentu. */
let licznik = 0;
function nowyIdentyfikator(przedrostek: string): string {
  licznik += 1;
  return `dm-${przedrostek}-${licznik}`;
}

/** Pole tekstowe jednowierszowe. */
export function poleTekstowe(opis: OpisPola): PoleFormularza<HTMLInputElement> {
  const kontrolka = document.createElement('input');
  kontrolka.type = 'text';
  kontrolka.className = 'dn-pole-kontrolka';
  kontrolka.id = nowyIdentyfikator('pole');
  if (opis.podpowiedz !== undefined) kontrolka.placeholder = opis.podpowiedz;
  return { element: obudowa(opis, kontrolka), kontrolka };
}

/**
 * Pole poświadczenia — jedyne pole sekcji, którego wartości nie da się
 * odczytać z rdzenia.
 *
 * Kontrakt przyjmuje `credential` w żądaniu i nie zwraca go żadną komendą.
 * Pole jest zatem wyłącznie wejściem: puste znaczy „nie zmieniaj", wypełnione
 * znaczy „zapisz nowe". Nigdy nie pokazuje wartości zapisanej, ponieważ klient
 * jej nie ma.
 */
export function poleTajne(opis: OpisPola): PoleFormularza<HTMLInputElement> {
  const kontrolka = document.createElement('input');
  kontrolka.type = 'password';
  kontrolka.className = 'dn-pole-kontrolka dm-pole--tajne';
  kontrolka.id = nowyIdentyfikator('tajne');
  kontrolka.autocomplete = 'new-password';
  kontrolka.spellcheck = false;
  if (opis.podpowiedz !== undefined) kontrolka.placeholder = opis.podpowiedz;
  return { element: obudowa(opis, kontrolka), kontrolka };
}

/** Pole tekstowe wielowierszowe. */
export function poleWielowierszowe(
  opis: OpisPola,
  wiersze: number,
): PoleFormularza<HTMLTextAreaElement> {
  const kontrolka = document.createElement('textarea');
  kontrolka.className = 'dn-pole-kontrolka dm-pole--tresc';
  kontrolka.id = nowyIdentyfikator('tresc');
  kontrolka.rows = wiersze;
  kontrolka.spellcheck = false;
  if (opis.podpowiedz !== undefined) kontrolka.placeholder = opis.podpowiedz;
  return { element: obudowa(opis, kontrolka), kontrolka };
}

/** Jedna pozycja listy wyboru. */
export interface PozycjaWyboru {
  wartosc: string;
  etykieta: string;
}

/** Lista wyboru zbudowana z podanych pozycji. */
export function poleWyboru(
  opis: OpisPola,
  pozycje: readonly PozycjaWyboru[],
): PoleFormularza<HTMLSelectElement> {
  const kontrolka = document.createElement('select');
  kontrolka.className = 'dn-pole-kontrolka';
  kontrolka.id = nowyIdentyfikator('wybor');
  ustawPozycje(kontrolka, pozycje);
  return { element: obudowa(opis, kontrolka), kontrolka };
}

/** Wymienia pozycje listy wyboru, utrzymując wybór, o ile nadal istnieje. */
export function ustawPozycje(
  kontrolka: HTMLSelectElement,
  pozycje: readonly PozycjaWyboru[],
): void {
  const poprzednia = kontrolka.value;
  kontrolka.replaceChildren(
    ...pozycje.map((pozycja) => {
      const element = document.createElement('option');
      element.value = pozycja.wartosc;
      element.textContent = pozycja.etykieta;
      return element;
    }),
  );
  if (pozycje.some((pozycja) => pozycja.wartosc === poprzednia)) {
    kontrolka.value = poprzednia;
  }
}

/** Pole logiczne: przełącznik wraz z etykietą po jego prawej. */
export function poleLogiczne(opis: OpisPola): PoleFormularza<HTMLInputElement> {
  const kontrolka = document.createElement('input');
  kontrolka.type = 'checkbox';
  kontrolka.className = 'dn-przelacznik';
  kontrolka.id = nowyIdentyfikator('logiczne');

  const etykieta = document.createElement('label');
  etykieta.className = 'dn-pole-etykieta dm-pole__etykieta-obok';
  etykieta.htmlFor = kontrolka.id;
  etykieta.textContent = opis.etykieta;

  const wiersz = document.createElement('div');
  wiersz.className = 'dn-pole dm-pole dm-pole--logiczne';
  wiersz.append(kontrolka, etykieta);
  if (opis.opis !== undefined) wiersz.append(zdanieOpisu(opis.opis));

  return { element: wiersz, kontrolka };
}

/** Przycisk formularza. */
export function przycisk(tresc: string, klasa: string): HTMLButtonElement {
  const element = document.createElement('button');
  element.type = 'button';
  element.className = klasa;
  element.textContent = tresc;
  return element;
}

/**
 * Wiersz odpowiedzi rdzenia stojący pod formularzem.
 *
 * Naciśnięcie zawsze daje odpowiedź. Powodzenie mówi, co się stało;
 * niepowodzenie mówi, co odpowiedział rdzeń, i zostaje na widoku aż do
 * następnego działania.
 */
export interface WierszOdpowiedzi {
  element: HTMLElement;
  pokaz(tresc: string, powodzenie: boolean): void;
  wyczysc(): void;
}

export function utworzWierszOdpowiedzi(): WierszOdpowiedzi {
  const element = document.createElement('p');
  element.className = 'dm-odpowiedz';
  element.hidden = true;

  return {
    element,
    pokaz(tresc, powodzenie) {
      element.textContent = tresc;
      element.hidden = tresc === '';
      element.dataset.powodzenie = String(powodzenie);
    },
    wyczysc() {
      element.textContent = '';
      element.hidden = true;
    },
  };
}

/** Etykieta nad kontrolką wraz ze zdaniem wyjaśniającym pod nią. */
function obudowa(opis: OpisPola, kontrolka: HTMLElement): HTMLElement {
  const element = document.createElement('div');
  element.className = 'dn-pole dm-pole';

  const etykieta = document.createElement('label');
  etykieta.className = 'dn-pole-etykieta';
  etykieta.htmlFor = kontrolka.id;
  etykieta.textContent = opis.etykieta;

  element.append(etykieta, kontrolka);
  if (opis.opis !== undefined) element.append(zdanieOpisu(opis.opis));
  return element;
}

function zdanieOpisu(tresc: string): HTMLElement {
  const element = document.createElement('p');
  element.className = 'dn-pole-opis';
  element.textContent = tresc;
  return element;
}

/**
 * Kontrolki formularzy sekcji modeli udostępniają jeden kształt pola dla całej
 * sekcji, przeznaczony dla bytów o polach stałych wynikających wprost
 * z kontraktu, a nie z katalogu wierszy konfiguracji.
 */
export * from './kontrolki-formularza-braki'; // kontrolki okien operacyjnych

/** Pole formularza obejmuje wiersz osadzany w układzie, kontrolkę wewnątrz wiersza oraz opcjonalny opis wyświetlany pod kontrolką. */
export interface PoleFormularza<T extends HTMLElement> {
  /** Wiersz osadzany w formularzu. */
  element: HTMLElement;
  /** Kontrolka wewnątrz wiersza. */
  kontrolka: T;
}

/** Opis pola formularza zawiera etykietę, podpowiedź wyświetlaną w kontrolce oraz opcjonalne zdanie wyjaśniające pod kontrolką. */
export interface OpisPola {
  etykieta: string;
  podpowiedz?: string;
  opis?: string;
}

/** Identyfikator wiążący etykietę z kontrolką, tworzony z przedrostka i licznika, pozostaje unikalny w obrębie całego dokumentu. */
let licznik = 0;
function nowyIdentyfikator(przedrostek: string): string {
  licznik += 1;
  return `dm-${przedrostek}-${licznik}`;
}

/** Pole tekstowe jednowierszowe formularza, budowane z opisu zawierającego etykietę, podpowiedź i opcjonalne wyjaśnienie. */
export function poleTekstowe(opis: OpisPola): PoleFormularza<HTMLInputElement> {
  const kontrolka = document.createElement('input');
  kontrolka.type = 'text';
  kontrolka.className = 'dn-pole-kontrolka';
  kontrolka.id = nowyIdentyfikator('pole');
  if (opis.podpowiedz !== undefined) kontrolka.placeholder = opis.podpowiedz;
  return { element: obudowa(opis, kontrolka), kontrolka };
}

/**
 * Pole poświadczenia jest jedynym polem sekcji, którego wartości nie da się
 * odczytać z rdzenia — służy wyłącznie jako wejście przy zapisie nowej
 * wartości.
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

/** Pole tekstowe wielowierszowe formularza, którego liczba wierszy jest parametrem wywołania, budowane z opisu pola. */
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

/** Jedna pozycja listy wyboru złożona z wartości zapisywanej oraz etykiety wyświetlanej w kontrolce wyboru. */
export interface PozycjaWyboru {
  wartosc: string;
  etykieta: string;
}

/** Lista wyboru zbudowana z podanych pozycji, z których każda niesie wartość zapisywaną oraz etykietę widoczną w kontrolce. */
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

/** Wymienia pozycje listy wyboru na podane od nowa, utrzymując dotychczasowy wybór, o ile odpowiadająca mu pozycja nadal istnieje. */
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

/** Pole logiczne formularza łączy przełącznik dwustanowy z etykietą umieszczoną po jego prawej stronie. */
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

/** Przycisk formularza budowany z podanej treści napisu oraz klasy stylu określającej jego wygląd w interfejsie. */
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

/** Obudowa pola stawia etykietę nad kontrolką wraz z opcjonalnym zdaniem wyjaśniającym umieszczonym pod nią. */
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

import { elementGodla } from '../ikony/ikony';
import type { MetodaWejscia } from './zrodlo-auth';

/**
 * Postać ekranu wejścia — teksty, pola i przesłona. Bez ani jednego wywołania
 * rdzenia i bez wiedzy o tym, co się po naciśnięciu stanie.
 *
 * Oddzielone od `ekran-logowania.ts`, który prowadzi rozmowę z rdzeniem:
 * rozpoznaje stan bramki, wysyła hasło, czyta odmowę, wpuszcza. Podział idzie
 * po szwie między tym, co widać, a tym, co się dzieje.
 *
 * Nazwy wewnętrzne — bramka, rdzeń, `auth.register` — zostają w kodzie
 * i w komentarzach; na ekran idzie zdanie do przeczytania raz.
 */

/** Tryb ekranu: wejście hasłem, pierwsze ustawienie hasła albo jego zmiana. */
export type Tryb =
  | 'wejscie'
  | 'zalozenie'
  | 'reset'
  // Krok drugi rejestracji: konto jest założone i czeka na drogę z listu.
  // Osobny tryb, a nie stan w `zalozenie`, bo formularz jest inny — pyta
  // wyłącznie o drogę i nie ma po co pokazywać pól, które już wypełniono.
  | 'potwierdzenie'
  // Odzyskanie konta: krok pierwszy pyta o adres, krok drugi o drogę z listu
  // i nowe hasło.
  | 'odzyskanie'
  | 'odzyskanie-haslo';

/**
 * Najkrótsze przyjmowane hasło — reguła formularza, nie rdzenia.
 *
 * Rdzeń odmawia w `auth.register` wyłącznie przy haśle pustym, a w
 * `auth.password.reset` przy pustym którymkolwiek z dwóch; długości nie
 * egzekwuje. Skoro ekran obiecuje minimum, ekran go też pilnuje. Sprawdzian
 * pada po naciśnięciu, nigdy jako blokada przycisku — tak samo jak zgodność
 * hasła z powtórzeniem.
 *
 * Gdyby rdzeń dostał własną regułę długości, ta stała ma zejść: dwie reguły
 * w dwóch miejscach rozjadą się przy pierwszej zmianie.
 */
export const MIN_ZNAKOW = 8;

/** Zdanie pomocnicze pod polami; pusty napis chowa akapit. */
export const OBJASNIENIA: Record<Tryb, string> = {
  wejscie: '',
  zalozenie:
    `Pierwsze uruchomienie. Podaj login, adres e-mail i hasło (min ${MIN_ZNAKOW} znaków). ` +
    'Na podany adres przyjdzie droga potwierdzenia — będzie on później jedyną drogą odzyskania konta.',
  reset:
    'Podaj hasło, którym logujesz się dziś, i nowe. Pozostałe zalogowania zostaną zamknięte.',
  potwierdzenie:
    'Konto czeka na potwierdzenie adresu. Przepisz drogę potwierdzenia z listu, ' +
    'który właśnie poszedł na podany adres. Droga wygasa po godzinie.',
  odzyskanie:
    'Podaj adres e-mail uwierzytelniający. Jeżeli pasuje do konta, pójdzie na niego ' +
    'droga potwierdzenia pozwalająca ustawić nowe hasło.',
  'odzyskanie-haslo':
    'Przepisz drogę potwierdzenia z listu i podaj nowe hasło. ' +
    'Wszystkie urządzenia zalogują się ponownie.',
};

/**
 * Napisy przycisku — jedna forma, czasownikowa, we wszystkich trzech trybach.
 * Przycisk nie jest nagłówkiem ani sterem, tylko czynnością, więc mówi, co się
 * stanie po naciśnięciu.
 *
 * „Załóż konto" nie obiecuje wejścia, bo rejestracja go nie daje: konto powstaje
 * niepotwierdzone, a token dostępu wydaje dopiero potwierdzenie adresu. Napis
 * obiecujący wejście kazałby czekać na coś, co nie nadejdzie.
 */
export const NAPISY_PRZYCISKU: Record<Tryb, string> = {
  wejscie: 'Zaloguj się',
  zalozenie: 'Załóż konto',
  reset: 'Zmień hasło',
  potwierdzenie: 'Potwierdź adres i wejdź',
  odzyskanie: 'Wyślij drogę odzyskania',
  'odzyskanie-haslo': 'Ustaw nowe hasło',
};

/**
 * Etykieta pierwszego pola zależy od trybu i od metody.
 *
 * Przy zmianie hasła to samo pole niesie hasło dotychczasowe, a nie nowe; jedna
 * etykieta na dwa znaczenia kazałaby zgadywać, które hasło się wpisuje. Przy
 * wejściu PIN-em niesie PIN, a nazwanie go hasłem kazałoby podać rzecz, której
 * rdzeń w tej metodzie nie porówna.
 */
export function etykietaPierwszego(tryb: Tryb, metoda: MetodaWejscia): string {
  if (tryb === 'reset') return 'Dotychczasowe hasło';
  if (tryb === 'wejscie' && metoda === 'pin') return 'PIN';
  return 'Hasło';
}

/**
 * Pole tekstowe biblioteki pól — login, adres e-mail i droga potwierdzenia.
 *
 * Osobne od `poleHasla`, bo te trzy wartości nie są sekretami i ukrywanie ich
 * kropkami utrudniałoby jedyną czynność, jaką się z nimi wykonuje: sprawdzenie,
 * czy przepisało się je bez pomyłki. Droga potwierdzenia jest długa i przepisuje
 * się ją z listu — zasłonięta byłaby nie do zweryfikowania okiem.
 */
export function poleTekstu(
  klasa: string,
  etykieta: string,
  autocomplete: AutoFill,
  rodzaj: 'text' | 'email' = 'text',
): PoleHasla {
  const opis = document.createElement('span');
  opis.className = 'dn-pole-etykieta';
  opis.textContent = etykieta;

  const kontrolka = document.createElement('input');
  kontrolka.type = rodzaj;
  kontrolka.className = `dn-pole-kontrolka ${klasa}__kontrolka`;
  kontrolka.autocomplete = autocomplete;

  const pole = document.createElement('label');
  pole.className = `dn-pole ${klasa}`;
  pole.append(opis, kontrolka);
  return {
    pole,
    kontrolka,
    ustawEtykiete: (tekst) => {
      opis.textContent = tekst;
    },
  };
}

/** Pole hasła w kształcie biblioteki pól (`komponenty/pole.css`). */
export interface PoleHasla {
  pole: HTMLElement;
  kontrolka: HTMLInputElement;
  /** Zmienia napis etykiety — ta sama kontrolka służy różnym trybom. */
  ustawEtykiete(tekst: string): void;
}

export function poleHasla(klasa: string, etykieta: string, autocomplete: AutoFill): PoleHasla {
  const opis = document.createElement('span');
  opis.className = 'dn-pole-etykieta';
  opis.textContent = etykieta;

  const kontrolka = document.createElement('input');
  kontrolka.type = 'password';
  kontrolka.className = `dn-pole-kontrolka ${klasa}__kontrolka`;
  kontrolka.autocomplete = autocomplete;

  const pole = document.createElement('label');
  pole.className = `dn-pole ${klasa}`;
  pole.append(opis, kontrolka);
  return {
    pole,
    kontrolka,
    ustawEtykiete: (tekst) => {
      opis.textContent = tekst;
    },
  };
}

/**
 * Pole „Nie wyloguj mnie" — przełącznik w kształcie biblioteki pól.
 *
 * Nastawa nazywa się `keepSignedIn` w kontrakcie i `trwanieWejscia` w rdzeniu,
 * a na ekranie ma jeden napis w całym produkcie: dwie nazwy jednej nastawy każą
 * zgadywać, czy to ta sama rzecz. Skutek jest dwustronny — magazyn w kliencie
 * i trwanie w rdzeniu.
 */
export function utworzNiewylogowuj(zaznaczone: boolean): {
  pole: HTMLElement;
  kontrolka: HTMLInputElement;
} {
  const kontrolka = document.createElement('input');
  kontrolka.type = 'checkbox';
  kontrolka.className = 'dn-przelacznik au-niewylogowuj__kontrolka';
  kontrolka.checked = zaznaczone;

  const opis = document.createElement('span');
  opis.className = 'au-niewylogowuj__opis';
  opis.textContent = 'Nie wyloguj mnie';

  const pole = document.createElement('label');
  pole.className = 'au-niewylogowuj';
  pole.append(kontrolka, opis);
  return { pole, kontrolka };
}

/** Segmenty metody wejścia — ster, nie wyświetlacz. */
export interface SegmentyMetody {
  element: HTMLElement;
  /**
   * Obsadza segmenty metodami, które otwierają bramkę na tej maszynie.
   *
   * @param metody metody czynne; jedna albo żadna chowa cały pas.
   * @param wybrana metoda bieżąca — segment wybrany niesie wartość nastawy.
   */
  pokaz(metody: MetodaWejscia[], wybrana: MetodaWejscia): void;
}

const NAZWY_METOD: Record<MetodaWejscia, string> = {
  haslo: 'Hasło',
  pin: 'PIN',
};

/**
 * Przełącznik metody wejścia w postaci segmentów pigułkowych.
 *
 * Segmenty, a nie zakładki „Zaloguj się / Zarejestruj": `auth.register` jest
 * wykonalna dokładnie raz i po pierwszym uruchomieniu odmawia trwale, więc
 * zakładka rejestracji byłaby przyciskiem pewnej odmowy. O tym, który formularz
 * pokazać, rozstrzyga rdzeń polem `gatewayConfigured` powitania. Segmenty niosą
 * wybór, który w rdzeniu naprawdę istnieje: hasło i PIN.
 *
 * Pas znika przy jednej metodzie, bo przełącznik z jedną pozycją jest napisem,
 * a nie sterem; segment wyszarzony byłby bramą.
 *
 * Windows Hello nie jest tu segmentem — ani czynnym, ani wyszarzonym. Metoda
 * odmawia z powodu pochodzenia dokumentu (WebAuthn wywodzi `rp_id` z adresu,
 * a interfejs stoi pod adresem IP), więc wróci wdrożeniem pod domeną po
 * `https`, nie dopisaniem kodu. Zdanie o przyczynie stoi w sekcji
 * „Uwierzytelnianie" Okna Ustawień, gdzie Hello się zakłada.
 */
export function utworzSegmentyMetody(naWybor: (metoda: MetodaWejscia) => void): SegmentyMetody {
  const element = document.createElement('div');
  element.className = 'au-metody';
  // Grupa radiowa, nie zakładki: zakładki obiecują czytnikowi ekranu, że pod
  // każdą leży inny panel treści, a tu leży ten sam formularz z inną nastawą.
  element.setAttribute('role', 'radiogroup');
  element.setAttribute('aria-label', 'Metoda wejścia');
  element.hidden = true;

  function pokaz(metody: MetodaWejscia[], wybrana: MetodaWejscia): void {
    element.hidden = metody.length < 2;
    element.replaceChildren(
      ...metody.map((metoda) => {
        const segment = document.createElement('button');
        segment.type = 'button';
        segment.className = 'au-metoda';
        segment.setAttribute('role', 'radio');
        segment.setAttribute('aria-checked', String(metoda === wybrana));
        segment.textContent = NAZWY_METOD[metoda];
        segment.addEventListener('click', () => naWybor(metoda));
        return segment;
      }),
    );
  }

  return { element, pokaz };
}

/**
 * Odnośnik pod formularzem — przestawia ekran na inny tryb i nic poza tym.
 *
 * Przycisk, nie `<a href>`: odnośnik nigdzie nie prowadzi, tylko przestawia ten
 * sam ekran na inny tryb, a `<a>` bez celu myliłby czytnik ekranu i środkowy
 * przycisk myszy.
 */
export function utworzOdnosnik(napis: string, naNacisniecie: () => void): HTMLButtonElement {
  const element = document.createElement('button');
  element.type = 'button';
  element.className = 'au-odnosnik';
  element.textContent = napis;
  element.addEventListener('click', naNacisniecie);
  return element;
}

/**
 * Odnośnik „Nie pamiętam hasła" — droga odzyskania konta listem.
 *
 * Naciska go ten, kto hasła nie pamięta, więc prowadzi do `auth.recover`,
 * a nie do zmiany hasła ze znanym hasłem dotychczasowym — tamta jest czynnością
 * Ustawień i wymaga hasła, którego tu z definicji nie ma.
 */
export function utworzOdnosnikResetu(naNacisniecie: () => void): HTMLButtonElement {
  return utworzOdnosnik('Nie pamiętam hasła', naNacisniecie);
}

/**
 * Odnośnik „Mam drogę potwierdzenia z listu" — powrót do kroku drugiego.
 *
 * Bez niego krok potwierdzenia jest osiągalny wyłącznie z udanej rejestracji
 * w tym samym oknie: odświeżenie strony między listem a przepisaniem drogi
 * zamykało konto niepotwierdzone przed platformą na głucho, bo rejestracja
 * drugi raz odmawia (`conflict`), a logowanie odmawia oczekiwaniem na
 * potwierdzenie adresu.
 */
export function utworzOdnosnikPotwierdzenia(naNacisniecie: () => void): HTMLButtonElement {
  return utworzOdnosnik('Mam drogę potwierdzenia z listu', naNacisniecie);
}

/**
 * Przesłona wejścia: pełny widok, karta pośrodku, znak marki nad nią.
 *
 * Znak bierze się z `ikony/marka.ts`, a nie z gołego nagłówka tekstowego.
 * Podłoże jasne, bo karta stoi na tle strony: znak ma barwy własne i odmianę
 * dobiera się do podłoża, nie do motywu.
 */
export function zlozPrzeslone(stany: HTMLElement, ponow: HTMLElement): HTMLElement {
  // Godło idzie osobno od nazwy, a nie odmianą złożoną: `zbudujElement` nadaje
  // znakowi bok kwadratowy, więc odmiana pionowa zostałaby ściśnięta.
  const godlo = elementGodla({ rozmiar: 56, podloze: 'jasne', etykieta: 'Danaco Console' });
  godlo.classList.add('au-godlo');

  const tytul = document.createElement('h1');
  tytul.className = 'au-tytul';
  tytul.textContent = 'Danaco Console';

  const naglowek = document.createElement('div');
  naglowek.className = 'au-naglowek';
  naglowek.append(godlo, tytul);

  const podtytul = document.createElement('p');
  podtytul.className = 'au-podtytul';
  podtytul.textContent = 'Twoja platforma AI Workspace OS';

  const karta = document.createElement('section');
  karta.className = 'dn-karta au-karta';
  karta.setAttribute('aria-label', 'Logowanie do Danaco Console');
  karta.append(naglowek, podtytul, stany, ponow);

  const element = document.createElement('div');
  element.className = 'au-brama';
  element.append(karta);
  return element;
}

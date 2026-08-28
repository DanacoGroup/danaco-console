import { elementGodla } from '../ikony/ikony';
import type { MetodaWejscia } from './zrodlo-auth';

/** Postać ekranu wejścia niesie teksty, pola i przesłonę, bez ani jednego wywołania rdzenia i bez wiedzy o tym, co się stanie po naciśnięciu; rozmowę z rdzeniem prowadzi ekran logowania osobno. */
export type Tryb =
  | 'wejscie'
  | 'zalozenie'
  | 'reset'
  // Krok drugi rejestracji: konto jest założone i czeka na drogę z listu.
  | 'potwierdzenie'
  // Odzyskanie konta: krok pierwszy pyta o adres, krok drugi o drogę z listu i nowe hasło.
  | 'odzyskanie'
  | 'odzyskanie-haslo';

/** Najkrótsze przyjmowane hasło jest regułą formularza, nie rdzenia: rdzeń długości nie egzekwuje, a ekran, skoro obiecuje minimum, sam go pilnuje. */
export const MIN_ZNAKOW = 8;

/** Zdanie pomocnicze pod polami formularza dla każdego trybu ekranu; pusty napis chowa cały akapit wyjaśnienia. */
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

/** Napisy przycisku mają jedną formę, czasownikową, we wszystkich trybach: przycisk mówi, co się stanie po naciśnięciu, nie jest nagłówkiem ani sterem. */
export const NAPISY_PRZYCISKU: Record<Tryb, string> = {
  wejscie: 'Zaloguj się',
  zalozenie: 'Załóż konto',
  reset: 'Zmień hasło',
  potwierdzenie: 'Potwierdź adres i wejdź',
  odzyskanie: 'Wyślij drogę odzyskania',
  'odzyskanie-haslo': 'Ustaw nowe hasło',
};

/** Etykieta pierwszego pola zależy od trybu i od metody, bo jedna etykieta na dwa znaczenia kazałaby zgadywać, które pole się wpisuje. */
export function etykietaPierwszego(tryb: Tryb, metoda: MetodaWejscia): string {
  if (tryb === 'reset') return 'Dotychczasowe hasło';
  if (tryb === 'wejscie' && metoda === 'pin') return 'PIN';
  return 'Hasło';
}

/** Pole tekstowe biblioteki pól służy loginowi, adresowi e-mail i drodze potwierdzenia, bo te wartości nie są sekretami i nie wymagają ukrycia treści. */
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

/** Pole hasła w kształcie biblioteki pól, z kontrolką typu hasła i metodą zmiany napisu etykiety dla różnych trybów ekranu. */
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

/** Pole nie wyloguj mnie jest przełącznikiem w kształcie biblioteki pól; nastawa ma jeden napis w całym produkcie, choć nosi dwie nazwy wewnętrzne. */
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

/** Segmenty metody wejścia są sterem, nie wyświetlaczem, z metodą obsadzenia segmentami dostępnymi na tej maszynie i wskazaniem metody bieżącej. */
export interface SegmentyMetody {
  element: HTMLElement;
  /** Obsadza segmenty metodami czynnymi; jedna metoda albo żadna chowa cały pas segmentów. */
  pokaz(metody: MetodaWejscia[], wybrana: MetodaWejscia): void;
}

const NAZWY_METOD: Record<MetodaWejscia, string> = {
  haslo: 'Hasło',
  pin: 'PIN',
};

/** Przełącznik metody wejścia w postaci segmentów pigułkowych zamiast zakładek, bo rejestracja jest wykonalna dokładnie raz i po pierwszym uruchomieniu odmawia trwale. */
export function utworzSegmentyMetody(naWybor: (metoda: MetodaWejscia) => void): SegmentyMetody {
  const element = document.createElement('div');
  element.className = 'au-metody';
  // Grupa radiowa, nie zakładki: tu leży ten sam formularz z inną nastawą, nie inny panel treści.
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

/** Odnośnik pod formularzem przestawia ekran na inny tryb i nic poza tym; przycisk, nie odnośnik adresowy, żeby nie mylić czytnika ekranu. */
export function utworzOdnosnik(napis: string, naNacisniecie: () => void): HTMLButtonElement {
  const element = document.createElement('button');
  element.type = 'button';
  element.className = 'au-odnosnik';
  element.textContent = napis;
  element.addEventListener('click', naNacisniecie);
  return element;
}

/** Odnośnik nie pamiętam hasła prowadzi do drogi odzyskania konta listem, nie do zmiany hasła ze znanym hasłem dotychczasowym. */
export function utworzOdnosnikResetu(naNacisniecie: () => void): HTMLButtonElement {
  return utworzOdnosnik('Nie pamiętam hasła', naNacisniecie);
}

/** Odnośnik mam drogę potwierdzenia z listu wraca do kroku drugiego rejestracji, bez którego krok ten byłby osiągalny wyłącznie z tego samego okna. */
export function utworzOdnosnikPotwierdzenia(naNacisniecie: () => void): HTMLButtonElement {
  return utworzOdnosnik('Mam drogę potwierdzenia z listu', naNacisniecie);
}

/** Przesłona wejścia składa pełny widok z kartą pośrodku i znakiem marki nad nią, na podłożu jasnym niezależnym od motywu interfejsu. */
export function zlozPrzeslone(stany: HTMLElement, ponow: HTMLElement): HTMLElement {
  // Godło idzie osobno od nazwy, a nie odmianą złożoną, żeby odmiana pionowa nie została ściśnięta.
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

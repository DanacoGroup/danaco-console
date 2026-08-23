/**
 * Katalog rodzajów sugestii Always On Display — rozdz. 4 opracowania
 * `docs/funkcje-globalne/always-on-display.md`.
 *
 * Rodzaje są cztery: `doradztwo`, `konfiguracja`, `problem`, `kolejny_krok`.
 * Wartości są wzięte wprost z opracowania i odpowiadają polu `rodzaj` encji
 * `sugestia_aod` modelu danych.
 *
 * WAŻNE — czego ten plik NIE robi. Struktura `AodSuggestion` kontraktu
 * (`shared/contract.ts`) niesie `id`, `text`, `commandType`, `windowId`
 * i `createdAt`. Pola rodzaju NIE MA. Ten katalog nie dokłada rodzaju do
 * kontraktu i nie zgaduje go z treści zdania: sugestia przychodząca z rdzenia
 * ma rodzaj nieznany i tak jest opisana Operatorowi. Rodzaj mają wyłącznie te
 * sugestie, których autorem jest nakładka — decyzje rozpoznane regułą
 * (`rozpoznanie-decyzji.ts`), bo tam nakładka wie, co rozpoznała.
 *
 * Waga sugestii (rozdz. 4.5) jest OSOBNĄ osią wobec wagi rozpoznania
 * (`WagaDecyzji`: pewna/sporna z `rozpoznanie-decyzji.ts`). Waga rozdz. 4.5
 * mówi, JAK GŁOŚNO sugestia ma się ujawnić; `WagaDecyzji` mówi, CZYJA to ocena.
 * Obie żyją obok siebie, żadna nie zastępuje drugiej.
 */

/** Rodzaj sugestii — pole `rodzaj` encji `sugestia_aod`, rozdz. 4. */
export const RodzajSugestii = {
  /** Rozdz. 4.1 — zaobserwowany stan i zalecane działanie. */
  Doradztwo: 'doradztwo',
  /** Rozdz. 4.2 — ustawienie, automatyka albo agent rozwiązujący powtarzalność. */
  Konfiguracja: 'konfiguracja',
  /** Rozdz. 4.3 — nazwanie problemu i jego umiejscowienia. */
  Problem: 'problem',
  /** Rozdz. 4.4 — nazwanie kroku następnego. */
  KolejnyKrok: 'kolejny_krok',
} as const;
export type RodzajSugestii = (typeof RodzajSugestii)[keyof typeof RodzajSugestii];

/** Waga sugestii i sposób jej ujawnienia — rozdz. 4.5. */
export const WagaUjawnienia = {
  /** Dymek otwierany samoczynnie, plakietka pulsująca. */
  Wysoka: 'wysoka',
  /** Plakietka liczbowa; dymek po kliknięciu awatara. */
  Srednia: 'srednia',
  /** Pozycja listy oczekujących; bez plakietki pulsującej. */
  Niska: 'niska',
} as const;
export type WagaUjawnienia = (typeof WagaUjawnienia)[keyof typeof WagaUjawnienia];

/** Nazwa rodzaju widziana przez Operatora. */
export const NAZWY_RODZAJOW: Readonly<Record<RodzajSugestii, string>> = {
  [RodzajSugestii.Doradztwo]: 'doradztwo',
  [RodzajSugestii.Konfiguracja]: 'sugestia konfiguracji',
  [RodzajSugestii.Problem]: 'wskazanie problemu',
  [RodzajSugestii.KolejnyKrok]: 'kolejny krok',
};

/** Kiedy sugestia danego rodzaju powstaje — wiersz „Kiedy powstaje" rozdz. 4. */
export const POWSTANIE_RODZAJU: Readonly<Record<RodzajSugestii, string>> = {
  [RodzajSugestii.Doradztwo]:
    'Zakończenie długiego zadania modułu, dostępny wynik do przeglądu, praca nad zasobem ' +
    'powiązanym z zadaniem oczekującym.',
  [RodzajSugestii.Konfiguracja]:
    'Powtarzalność czynności ręcznej powyżej progu, brak konfiguracji potrzebnej do wykonania ' +
    'polecenia, praca w module bez przypisanego wykonawcy.',
  [RodzajSugestii.Problem]:
    'Negatywny wynik kontroli jakości, zadanie w stanie błędu, powtórzone ponowienie tego ' +
    'samego zadania, kolejka zatrzymana, pętla przerwana.',
  [RodzajSugestii.KolejnyKrok]:
    'Zlecenie oczekujące na zatwierdzenie, nadejście terminu reguły harmonogramu, zakończony ' +
    'etap pracy z jednoznacznym następstwem.',
};

/** Nazwa wagi widziana przez Operatora. */
export const NAZWY_WAG: Readonly<Record<WagaUjawnienia, string>> = {
  [WagaUjawnienia.Wysoka]: 'waga wysoka',
  [WagaUjawnienia.Srednia]: 'waga średnia',
  [WagaUjawnienia.Niska]: 'waga niska',
};

/** Jedno działanie sugestii wraz z jego skutkiem — tabele „Działanie / Skutek" rozdz. 4. */
export interface DzialanieRodzaju {
  /** Napis działania widziany przez Operatora. */
  nazwa: string;
  /** Skutek działania opisany opracowaniem. */
  skutek: string;
}

/**
 * Komplet działań każdego rodzaju — rozdz. 4.1–4.4.
 *
 * Katalog jest opisem, nie wykonaniem: mówi, jakie działania rodzaj niesie
 * i co każde robi. Które z nich mają dziś za sobą komendę kontraktu,
 * rozstrzyga `cztery-stery.ts` przy konkretnej decyzji — i tam, gdzie komendy
 * nie ma, stoi zdanie nazywające granicę zamiast przycisku-atrapy.
 */
export const DZIALANIA_RODZAJU: Readonly<Record<RodzajSugestii, readonly DzialanieRodzaju[]>> = {
  [RodzajSugestii.Doradztwo]: [
    {
      nazwa: 'Przejdź do wyniku',
      skutek: 'Otwiera moduł i okno operacyjne, w którym wynik powstał; karta sesji pozostaje ta sama.',
    },
    {
      nazwa: 'Przenieś do Chat Window',
      skutek: 'Wstawia treść sugestii do pola polecenia Chat Window jako polecenie gotowe do wysłania.',
    },
    {
      nazwa: 'Odłóż',
      skutek: 'Pozycja wraca do listy oczekujących, dymek się zamyka, status pozostaje „nowa".',
    },
    {
      nazwa: 'Odrzuć',
      skutek: 'Status „odrzucona"; analogiczne zdarzenie nie tworzy sugestii do końca karty sesji.',
    },
  ],

  [RodzajSugestii.Konfiguracja]: [
    {
      nazwa: 'Otwórz ustawienie',
      skutek: 'Otwiera okno konfiguracji na wskazanej pozycji, z objaśnieniem kontekstowym.',
    },
    {
      nazwa: 'Utwórz automatykę',
      skutek: 'Otwiera moduł Automations z formularzem wypełnionym rozpoznaną sekwencją.',
    },
    {
      nazwa: 'Przypisz agenta',
      skutek: 'Otwiera wybór agenta z modułu Agents dla bieżącego kontekstu pracy.',
    },
    {
      nazwa: 'Odrzuć',
      skutek: 'Status „odrzucona"; funkcja nie powtarza tej sugestii dla tego samego kontekstu.',
    },
  ],

  [RodzajSugestii.Problem]: [
    {
      nazwa: 'Otwórz Execution Loop Window',
      skutek: 'Otwiera okno pętli wykonawczej na zadaniu, którego dotyczy problem.',
    },
    {
      nazwa: 'Ponów zadanie',
      skutek: 'Wysyła do Koordynatora decyzję o ponowieniu wskazanego zadania.',
    },
    {
      nazwa: 'Skoryguj zlecenie',
      skutek: 'Wstawia treść zlecenia do pola polecenia Chat Window w trybie edycji.',
    },
    {
      nazwa: 'Wstrzymaj proces',
      skutek: 'Wywołuje wstrzymanie kolejki albo procesu; w trybie obserwatora prowadzi do przełączenia trybu.',
    },
    {
      nazwa: 'Odrzuć',
      skutek: 'Status „odrzucona"; problem pozostaje widoczny w Execution Loop Window, sugestia nie wraca.',
    },
  ],

  [RodzajSugestii.KolejnyKrok]: [
    {
      nazwa: 'Zatwierdź krok',
      skutek: 'Przekazuje zatwierdzenie punktu decyzyjnego; pętla przechodzi do kroku następnego.',
    },
    {
      nazwa: 'Wznów proces',
      skutek: 'Wywołuje wznowienie wstrzymanej kolejki albo procesu.',
    },
    {
      nazwa: 'Przenieś do Chat Window',
      skutek: 'Wstawia krok jako polecenie do pola polecenia Chat Window.',
    },
    {
      nazwa: 'Odłóż',
      skutek: 'Pozycja wraca do listy oczekujących; punkt decyzyjny pozostaje otwarty w Execution Loop Window.',
    },
  ],
};

/** Status sugestii w jej cyklu życia — rozdz. 9.4 opracowania. */
export const StatusSugestii = {
  /** Sugestia utworzona, bez reakcji Operatora. */
  Nowa: 'nowa',
  /** Działanie wykonane. */
  Przyjeta: 'przyjeta',
  /** Odrzucona wprost albo po upływie czasu życia. */
  Odrzucona: 'odrzucona',
} as const;
export type StatusSugestii = (typeof StatusSugestii)[keyof typeof StatusSugestii];

/**
 * Czy sugestia o tej wadze ujawnia się samoczynnie dymkiem.
 *
 * Rozdz. 3.4: „Waga sugestii ujawnianej samoczynnie — wysoka". Wartość jest
 * ustawieniem konfiguracyjnym; do czasu, aż kontrakt poniesie ustawienia
 * zasięgu „Always On Display", obowiązuje wartość domyślna opracowania.
 */
export function czyUjawniaSieSamoczynnie(waga: WagaUjawnienia): boolean {
  return waga === WagaUjawnienia.Wysoka;
}

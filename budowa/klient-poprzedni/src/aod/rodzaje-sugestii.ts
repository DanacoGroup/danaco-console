/** Rodzaj sugestii, odpowiadający polu `rodzaj` encji `sugestia_aod` modelu danych: cztery rozróżnione wartości. */
export const RodzajSugestii = {
  /** Zaobserwowany stan i zalecane działanie. */
  Doradztwo: 'doradztwo',
  /** Ustawienie, automatyka albo agent rozwiązujący powtarzalność. */
  Konfiguracja: 'konfiguracja',
  /** Nazwanie problemu i jego umiejscowienia. */
  Problem: 'problem',
  /** Nazwanie kroku następnego. */
  KolejnyKrok: 'kolejny_krok',
} as const;
export type RodzajSugestii = (typeof RodzajSugestii)[keyof typeof RodzajSugestii];

/** Waga sugestii i sposób jej ujawnienia Operatorowi: dymkiem, plakietką liczbową albo samą pozycją listy. */
export const WagaUjawnienia = {
  /** Dymek otwierany samoczynnie, plakietka pulsująca. */
  Wysoka: 'wysoka',
  /** Plakietka liczbowa; dymek po kliknięciu awatara. */
  Srednia: 'srednia',
  /** Pozycja listy oczekujących; bez plakietki pulsującej. */
  Niska: 'niska',
} as const;
export type WagaUjawnienia = (typeof WagaUjawnienia)[keyof typeof WagaUjawnienia];

/** Nazwa rodzaju sugestii widziana przez Operatora zawsze w treści dymka i na liście sugestii oczekujących. */
export const NAZWY_RODZAJOW: Readonly<Record<RodzajSugestii, string>> = {
  [RodzajSugestii.Doradztwo]: 'doradztwo',
  [RodzajSugestii.Konfiguracja]: 'sugestia konfiguracji',
  [RodzajSugestii.Problem]: 'wskazanie problemu',
  [RodzajSugestii.KolejnyKrok]: 'kolejny krok',
};

/** Zdanie opisujące, kiedy sugestia danego rodzaju powstaje, osobno dla każdego z czterech rodzajów sugestii. */
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

/** Nazwa wagi ujawnienia sugestii widziana przez Operatora zawsze w treści dymka i na plakietce liczbowej. */
export const NAZWY_WAG: Readonly<Record<WagaUjawnienia, string>> = {
  [WagaUjawnienia.Wysoka]: 'waga wysoka',
  [WagaUjawnienia.Srednia]: 'waga średnia',
  [WagaUjawnienia.Niska]: 'waga niska',
};

/** Jedno działanie sugestii wraz z jego skutkiem, opisane napisem przycisku oraz zdaniem skutku działania. */
export interface DzialanieRodzaju {
  // Napis działania widziany przez Operatora.
  nazwa: string;
  // Zdanie opisujące skutek działania.
  skutek: string;
}

/** Komplet działań każdego rodzaju sugestii; katalog jest opisem, a nie wykonaniem tych działań w oknie. */
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

/** Status sugestii w jej cyklu życia: od utworzenia, przez przyjęcie działania, po odrzucenie sugestii. */
export const StatusSugestii = {
  /** Sugestia utworzona, bez reakcji Operatora. */
  Nowa: 'nowa',
  /** Działanie wykonane. */
  Przyjeta: 'przyjeta',
  /** Odrzucona wprost albo po upływie czasu życia. */
  Odrzucona: 'odrzucona',
} as const;
export type StatusSugestii = (typeof StatusSugestii)[keyof typeof StatusSugestii];

/** Rozstrzyga, czy sugestia o danej wadze ujawnia się samoczynnie dymkiem, czy tylko plakietką na awatarze. */
export function czyUjawniaSieSamoczynnie(waga: WagaUjawnienia): boolean {
  return waga === WagaUjawnienia.Wysoka;
}

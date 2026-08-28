import {
  StudioAuthor,
  StudioMarkupKind,
  StudioMarkupState,
  type StudioMarkup,
} from '../../../../shared/contract';

/** Interfejs BarwaZnacznika opisuje jedną barwę znacznika wyrażoną żetonem motywu: kod, nazwę widoczną i nazwę żetonu tła. */
export interface BarwaZnacznika {
  kod: string;
  nazwa: string;
  /** Nazwa żetonu tła; arkusz `przybornik-znakowania.css` czyta ją z `data-barwa`. */
  zeton: string;
}

/**
 * Paleta barw znaczników.
 *
 * Sześć barw, a nie dowolna: barwa znacznika ma coś znaczyć na pierwszy rzut
 * oka, a przy dwudziestu odcieniach nie znaczy nic. Wykaz jest jawny, więc
 * Operator widzi, z czego wybiera.
 */
export const BARWY_ZNACZNIKOW: readonly BarwaZnacznika[] = [
  { kod: 'sygnal', nazwa: 'Sygnałowa', zeton: '--dn-sygnal-tlo' },
  { kod: 'ostrzezenie', nazwa: 'Ostrzegawcza', zeton: '--dn-ostrzezenie-tlo' },
  { kod: 'blad', nazwa: 'Alarmowa', zeton: '--dn-blad-tlo' },
  { kod: 'sukces', nazwa: 'Potwierdzająca', zeton: '--dn-sukces-tlo' },
  { kod: 'informacja', nazwa: 'Informacyjna', zeton: '--dn-informacja-tlo' },
  { kod: 'obojetna', nazwa: 'Obojętna', zeton: '--dn-powierzchnia-3' },
];

/** Interfejs ZnacznikWlasny opisuje jeden znacznik założony na fragmencie dokumentu wraz z jego zakresem, barwą, autorem i stanem otwarcia. */
export interface ZnacznikWlasny {
  kod: string;
  /** Nazwa znacznika — pełna, bez numeracji wymyślonej. */
  nazwa: string;
  /** Kod barwy z palety. */
  barwa: string;
  /** Zakres w znakach treści; `null` znaczy „cały dokument". */
  zakres: { poczatek: number; koniec: number } | null;
  /** Kto znacznik założył — Operator albo model. */
  autor: StudioAuthor;
  /** Czy znacznik jest jeszcze otwarty (nieodhaczony). */
  otwarty: boolean;
  /** Czas założenia w milisekundach epoki. */
  czas: number;
}

/**
 * Stała NAZWY_ZNACZNIKOW_GOTOWE niesie nazwy znaczników podpowiadane Operatorowi; nie jest wykazem zamkniętym, bo pole nazwy przyjmuje dowolną treść wpisaną przez Operatora.
 */
export const NAZWY_ZNACZNIKOW_GOTOWE: readonly string[] = [
  'do sprawdzenia',
  'wymaga źródła',
  'gotowe',
];

/** Interfejs ZnacznikiWlasne udostępnia zbiór znaczników dokumentu wraz z czynnościami zakładania, przestawiania, zdejmowania i czyszczenia wykazu. */
export interface ZnacznikiWlasne {
  /** Znaczniki dokumentu w kolejności założenia. */
  wykaz(): readonly ZnacznikWlasny[];
  /** Zakłada znacznik i oddaje go; nazwa pusta jest odrzucana. */
  zaloz(
    nazwa: string,
    barwa: string,
    zakres: { poczatek: number; koniec: number } | null,
    autor: StudioAuthor,
  ): ZnacznikWlasny | null;
  /** Przestawia odhaczenie znacznika. */
  przestaw(kod: string, otwarty: boolean): void;
  /** Zdejmuje znacznik. */
  zdejmij(kod: string): void;
  /** Zdejmuje wszystkie znaczniki — po zamknięciu dokumentu. */
  wyczysc(): void;
}

/**
 * Funkcja utworzZnacznikiWlasne zakłada zbiór znaczników jednego dokumentu; kod znacznika bierze się z licznika i czasu założenia, nie jest numeracją widoczną dla Operatora.
 */
export function utworzZnacznikiWlasne(): ZnacznikiWlasne {
  const znaczniki: ZnacznikWlasny[] = [];
  let licznik = 0;

  return {
    wykaz: () => znaczniki,

    zaloz(nazwa, barwa, zakres, autor) {
      const oczyszczona = nazwa.trim();
      if (oczyszczona === '') return null;
      licznik += 1;
      const znacznik: ZnacznikWlasny = {
        kod: `znacznik-${licznik}-${Date.now()}`,
        nazwa: oczyszczona,
        barwa: BARWY_ZNACZNIKOW.some((pozycja) => pozycja.kod === barwa)
          ? barwa
          : (BARWY_ZNACZNIKOW[0]?.kod ?? 'sygnal'),
        zakres,
        autor,
        otwarty: true,
        czas: Date.now(),
      };
      znaczniki.push(znacznik);
      return znacznik;
    },

    przestaw(kod, otwarty) {
      const znacznik = znaczniki.find((pozycja) => pozycja.kod === kod);
      if (znacznik !== undefined) znacznik.otwarty = otwarty;
    },

    zdejmij(kod) {
      const gdzie = znaczniki.findIndex((pozycja) => pozycja.kod === kod);
      if (gdzie >= 0) znaczniki.splice(gdzie, 1);
    },

    wyczysc() {
      znaczniki.length = 0;
    },
  };
}

/** Funkcja przybornikZetonBarwy zwraca żeton barwy znacznika; barwa nieznana schodzi na barwę obojętną, a nie na pustkę. */
export function przybornikZetonBarwy(barwa: string): string {
  const pozycja = BARWY_ZNACZNIKOW.find((wpis) => wpis.kod === barwa);
  return pozycja?.zeton ?? '--dn-powierzchnia-3';
}

/** Stała NAZWY_RODZAJOW_ZNAKOWANIA niesie nazwę rodzaju znakowania trwałego widoczną dla Operatora w interfejsie. */
export const NAZWY_RODZAJOW_ZNAKOWANIA: Readonly<Record<StudioMarkupKind, string>> = {
  [StudioMarkupKind.Highlight]: 'Wyróżnienie barwą',
  [StudioMarkupKind.Mark]: 'Znacznik własny',
  [StudioMarkupKind.Suggestion]: 'Propozycja zmiany na marginesie',
};

/** Stała NAZWY_STANOW_ZNAKOWANIA niesie nazwę stanu znakowania trwałego widoczną dla Operatora w interfejsie. */
export const NAZWY_STANOW_ZNAKOWANIA: Readonly<Record<StudioMarkupState, string>> = {
  [StudioMarkupState.Open]: 'otwarte',
  [StudioMarkupState.Accepted]: 'przyjęte',
  [StudioMarkupState.Rejected]: 'odrzucone',
  [StudioMarkupState.Resolved]: 'rozwiązane',
};

/**
 * Zdanie o znakowaniu leżącym W RDZENIU.
 *
 * Odróżnia się od zdania o znaczniku sesyjnym tym, że nazywa wykonawcę:
 * znakowanie modelu jest podpisane jako `model`, a rdzeń niesie przy nim nazwę
 * eksperta, więc dwóch wykonawców nie zlewa się w jednego.
 */
export function przybornikOpiszZnakowanieRdzenia(znakowanie: StudioMarkup): string {
  const kto =
    znakowanie.author === StudioAuthor.Model
      ? `model — ${znakowanie.authorAgentName ?? 'wykonawca nienazwany'}`
      : 'Operator';
  const rodzajZnacznika =
    znakowanie.markType === undefined || znakowanie.markType === ''
      ? ''
      : ` · rodzaj „${znakowanie.markType}"`;
  const barwa =
    znakowanie.color === undefined || znakowanie.color === '' ? '' : ` · barwa ${znakowanie.color}`;
  return (
    `${NAZWY_RODZAJOW_ZNAKOWANIA[znakowanie.kind] ?? znakowanie.kind} · ${kto} · ` +
    `znaki ${znakowanie.rangeStart}–${znakowanie.rangeEnd} · ` +
    `${NAZWY_STANOW_ZNAKOWANIA[znakowanie.state] ?? znakowanie.state}${rodzajZnacznika}${barwa} · ` +
    `${new Date(znakowanie.createdAt).toLocaleString('pl-PL')}`
  );
}

/** Funkcja przybornikOpiszZnacznik zwraca zdanie o znaczniku wraz z jego zakresem, autorem, stanem otwarcia i czasem założenia. */
export function przybornikOpiszZnacznik(znacznik: ZnacznikWlasny): string {
  const zakres =
    znacznik.zakres === null
      ? 'dotyczy całego dokumentu'
      : `znaki ${znacznik.zakres.poczatek}–${znacznik.zakres.koniec}`;
  const autor = znacznik.autor === StudioAuthor.Model ? 'model' : 'Operator';
  const stan = znacznik.otwarty ? 'otwarty' : 'odhaczony';
  return `${autor} · ${zakres} · ${stan} · ${new Date(znacznik.czas).toLocaleString('pl-PL')}`;
}

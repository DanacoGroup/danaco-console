import { Command } from '../../../../shared/contract';

// Teksty dla Operatora w module Design: słownik nazw okien i wykaz funkcji bez drogi w kontrakcie.

/**
 * Kod modułu z kolumny `modul.kod` rdzenia.
 *
 * Jedyne miejsce w kliencie wiążące ten katalog z kodem modułu Design.
 * Nie jest kopią katalogu modułów: moduł mówi wyłącznie,
 * którym modułem sam jest — wykazu pozostałych czternastu tu nie ma.
 */
export const KOD_MODULU = 'design';

/** Nazwa i rola okna operacyjnego modułu Design, zgodnie z katalogiem okien prowadzonym przez rdzeń systemu. */
export interface OpisOkna {
  /** Kod okna w katalogu rdzenia (`okna_modulu.kod`). */
  kod: string;
  nazwa: string;
  rola: string;
}

export const OKNO_DESIGN_BOARD: OpisOkna = {
  kod: 'design-board',
  nazwa: 'Design Board',
  rola: 'wiodące',
};

export const OKNO_ASSETS_PANEL: OpisOkna = {
  kod: 'assets-panel',
  nazwa: 'Assets Panel',
  rola: 'zarządca',
};

export const OKNO_PROMPT_BUILDER: OpisOkna = {
  kod: 'prompt-builder',
  nazwa: 'Prompt Builder',
  rola: 'kreator',
};

/**
 * Czwarte okno modułu. W rejestrze okien operacyjnych definicja jest jedna,
 * pod bezmodułowym kodem `preview-window`, a przypięcia są dwa — do Studio
 * i do Design. Ten opis dotyczy przypięcia Designu.
 */
export const OKNO_PREVIEW_WINDOW: OpisOkna = {
  kod: 'preview-window',
  nazwa: 'Preview Window',
  rola: 'pomocnicze',
};

/**
 * Piąte okno operacyjne modułu, wymienione wraz z rolą i zawartością. Kod nie pochodzi z rejestru
 * okien rdzenia, lecz z tej samej reguły, którą stosują pozostałe moduły budujące okna spoza
 * rejestru: nazwa własna zapisana małymi literami z łącznikami.
 */
export const OKNO_TOKENS_SYSTEM_PANEL: OpisOkna = {
  kod: 'tokens-system-panel',
  nazwa: 'Tokens & System Panel',
  rola: 'zarządca',
};

/** Kody okien operacyjnych budowanych przez moduł Design, z pominięciem okna rozmowy prowadzonej z modelem. */
export const KODY_OKIEN: readonly string[] = [
  OKNO_DESIGN_BOARD.kod,
  OKNO_PREVIEW_WINDOW.kod,
  OKNO_ASSETS_PANEL.kod,
  OKNO_PROMPT_BUILDER.kod,
  OKNO_TOKENS_SYSTEM_PANEL.kod,
];

/** Funkcja panelu akcji modułu Design, dla której kontrakt aplikacji nie ma jeszcze przypisanej komendy rdzenia. */
export interface BrakDrogi {
  /** Kod używany w atrybucie `data-brak` — po nim pyta sprawdzian. */
  kod: string;
  /** Nazwa czynności widoczna na kontrolce. */
  nazwa: string;
  /** Zdanie mówiące, czego brakuje i dlaczego czynność nie idzie do rdzenia. */
  powod: string;
}

/**
 * Czynności panelu akcji modułu, których okno nie wykonuje, bo kontrakt nie niósł dla nich komendy;
 * powody opisują wyłącznie drogę do rdzenia, nie stan pokrycia, który mierzy pas uczciwości modułu.
 */
export const BRAKI: Readonly<Record<string, BrakDrogi>> = {
  // Przypisanie zasobu do kolekcji nie jest tu brakiem: to osobna czynność, nie odmiana etykietowania.
  eksportZbiorczy: {
    kod: 'eksport-zbiorczy',
    nazwa: 'Pobierz zbiorczo',
    powod:
      `Eksport zbiorczy MA JUŻ KOMENDĘ: ${Command.DesignAssetExportBatch} przyjmuje zestaw ` +
      'zasobów, format i komplet skal, a odmowa jednego zasobu nie wstrzymuje pozostałych. ' +
      'Brakuje jej uchwytu w rdzeniu, więc okno jeszcze jej nie woła. Stan mierzy pas ' +
      'uczciwości modułu — on pyta rdzeń, a to zdanie tylko opisuje drogę.',
  },
  wersjeKompozycji: {
    kod: 'wersje-kompozycji',
    nazwa: 'Wersje kompozycji',
    powod:
      `Wersje kompozycji MAJĄ JUŻ KOMENDY: ${Command.DesignBoardVersionSave} utrwala układ pod ` +
      `nazwą, ${Command.DesignBoardVersionList} oddaje ciąg jego postaci, ` +
      `a ${Command.DesignBoardVersionRestore} przywraca wybraną, zakładając przy tym wersję ` +
      'z układu sprzed cofnięcia. Kontrakt zna też zdarzenie zmiany kompozycji, więc drugie okno ' +
      'dowie się o zapisie. Brakuje uchwytów w rdzeniu; porównania dwóch wersji wprost nie ma ' +
      'i zestawia się je odczytem obu.',
  },
  eksportKompozycji: {
    kod: 'eksport-kompozycji',
    nazwa: 'Eksportuj kompozycję',
    powod:
      `Eksport kompozycji MA JUŻ KOMENDĘ: ${Command.DesignBoardExport} wyrysowuje całą tablicę ` +
      'albo wskazany obszar. Wyrys wymaga jednak treści zasobów leżących w warstwach, więc ' +
      'zależy od tej samej drogi po bajty, która gasi dziś podgląd. Brakuje uchwytów w rdzeniu.',
  },
  kursorWspolpracy: {
    kod: 'kursor-wspolpracy',
    nazwa: 'Kursor współpracy',
    powod:
      `Obecność MA JUŻ DROGĘ: ${Command.DesignPresenceReport} zgłasza położenie kursora ` +
      'i zaznaczenie, a rdzeń rozgłasza je pozostałym osobnym zdarzeniem obecności. Zgłoszenie ' +
      'jest ulotne i nie zapisuje się w bazie — położenie kursora sprzed godziny nie jest wiedzą ' +
      'o niczym. Brakuje uchwytu w rdzeniu.',
  },
  szablonPromptu: {
    kod: 'szablon-promptu',
    nazwa: 'Zapisz jako szablon',
    powod:
      `Szablony MAJĄ JUŻ KOMENDY: ${Command.DesignPromptTemplateSave} utrwala prompt pod nazwą, ` +
      `${Command.DesignPromptTemplateList} oddaje szablony okna, a ${Command.DesignPromptHistoryList} ` +
      'oddaje prompty wydane wraz z zasobami, które z nich powstały. Do czasu dobudowy uchwytów ' +
      'historia w tym oknie nadal ginie z zamknięciem karty przeglądarki.',
  },
  // Dwa braki Preview Window.
  eksportZasobu: {
    kod: 'eksport-zasobu',
    nazwa: 'Eksportuj zasób',
    powod:
      `Eksport pojedynczego zasobu MA JUŻ KOMENDĘ: ${Command.DesignAssetExport} oddaje bajty ` +
      'gotowe do zapisania poza produktem, wraz z formatem, skalą, jakością i rozstrzygnięciem ' +
      'o metadanych. Konwersja obrazu tego nie zastępowała i nie zastępuje: zakłada NOWY zasób ' +
      'w magazynie i tam się kończy. Brakuje uchwytu w rdzeniu. Liczby komend obszaru to zdanie ' +
      'nie podaje — podaje ją pas uczciwości modułu, mierząc ją odczytem.',
  },
  akceptacjaWyniku: {
    kod: 'akceptacja-wyniku',
    nazwa: 'Akceptuj wynik',
    powod:
      `Przyjęcie wyniku znaczy oznaczenie zasobu ulubionym, a komenda ${Command.DesignAssetFavoriteSet} ` +
      'weszła do kontraktu 14.08.2026 — brakiem jest więc już tylko DROGA Z TEGO OKNA. ' +
      'Kontrolka stoi w Assets Panel („Oznacz ulubionym", czynnosci-zasobu.ts) i dotyczy ' +
      'zasobu wskazanego w wykazie; podgląd wskazania nie prowadzi. Wskaż zasób w Assets ' +
      'Panel i oznacz go tam.',
  },
};

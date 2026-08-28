import { Command } from '../../../../shared/contract';
import { pokazKomunikat } from '../../aplikacja/komunikaty';
import type { Kanal } from '../../protokol/kanal';
import { utworzKatalogOkien, type KatalogOkien } from '../katalog-okien';
import { utworzPokrycieKomend, type PokrycieKomend } from '../pokrycie-komend';
import { KOD_MODULU, KODY_OKIEN } from './etykiety-designu';

/**
 * Pas uczciwości modułu Design — czego moduł nie umie, choć wygląda, jakby umiał, mierzone
 * odczytem z rdzenia, nie liczbą wpisaną w napis.
 */
export interface PasekUczciwosci {
  element: HTMLElement;
  /** Pyta rdzeń o katalog okien i wykaz komend. Woła się raz, po montażu. */
  odczytaj(): Promise<void>;
  /** Odpina kontrolki od bytów wspólnych. Wołane z `rozlacz()` modułu. */
  zamknij(): void;
}

/**
 * Komendy, którymi moduł żyje — osiem obszaru własnego i sześć obszaru obrazu, wzięte z generatu
 * kontraktu, nie wpisane wprost.
 */
const KOMENDY_MODULU: readonly string[] = [
  Command.DesignAssetGenerate,
  Command.DesignAssetList,
  Command.DesignAssetUpload,
  Command.DesignAssetTagSet,
  Command.DesignAssetFavoriteSet,
  Command.DesignAssetRemove,
  Command.DesignBoardUpdate,
  Command.DesignBoardList,
  Command.ImageInspect,
  Command.ImageTransform,
  Command.ImageAdjust,
  Command.ImageConvert,
  Command.ImageUpscale,
  Command.ImageBackgroundRemove,
];

/**
 * Czynności opracowania wraz z komendą, która je wykonuje, wzięte z generatu kontraktu, nie
 * z literałów wpisanych w oknie.
 */
const CZYNNOSCI_MODULU: readonly (readonly [string, string, string])[] = [
  [
    'Treść zasobu do przeglądarki',
    Command.DesignAssetContentGet,
    'oddanie bajtów zasobu, bez których podgląd nie ma czego pokazać',
  ],
  ['Eksport zasobu', Command.DesignAssetExport, 'wydanie pojedynczego zasobu plikiem'],
  ['Eksport zbiorczy', Command.DesignAssetExportBatch, 'wydanie wielu zasobów naraz'],
  ['Eksport kompozycji', Command.DesignBoardExport, 'wyrys całej kompozycji do pliku'],
  [
    'Wersje kompozycji',
    Command.DesignBoardVersionList,
    'historia stanów tablicy i powrót do wersji',
  ],
  ['Kolekcje zasobów', Command.DesignCollectionList, 'grupowanie zasobów w kolekcje tematyczne'],
  [
    'Szablony promptu',
    Command.DesignPromptTemplateSave,
    'trwały zapis promptu do wielokrotnego użycia',
  ],
  ['Kursor współpracy', Command.DesignPresenceReport, 'obecność drugiego Operatora na kompozycji'],
  ['Zestaw żetonów', Command.DesignTokensetSave, 'trwały zapis systemu projektowego'],
];

const TRESC_PREVIEW_WINDOW =
  'Katalog rdzenia niesie dla Preview Window jedną definicję o kodzie bezmodułowym ' +
  '(„preview-window", rola pomocnicze) i dwa przypięcia do niej: do Studia i do Designu. ' +
  'Rozstrzyga to połowę pytania — definicja jest jedna — i nic nie mówi o zawartości. ' +
  'Zawartości katalog zrównać nie może: Studio podgląda dokument w formacie wyjściowym, ' +
  'a Design zasób wizualny wraz z rodziną wariantów. Druga połowa — jedno okno ' +
  'konfigurowalne czy dwie osobne przestrzenie na wspólnej definicji — pozostaje otwarta.';

const TRESC_OBRAZU =
  'Zasób powstaje, obrazu wciąż nie widać — ale brak przesunął się o jeden szczebel. Rdzeń ' +
  'generuje: wysyła polecenie kanałem obrazowym (adapter „obrazy"), odbiera fragment obrazu, ' +
  'odkłada bajty w magazynie pod sumą kontrolną, mierzy format i wymiary z nagłówka utrwalonego ' +
  'pliku i dopiero wtedy zakłada wiersz zasobu — z odsyłaczem, formatem i wymiarami. Zasobu bez ' +
  'bajtów nie zakłada żadną drogą; każdy brak kończy się odmową nazywającą brak, nigdy obrazem ' +
  'zastępczym. Pole odsyłacza jest ścieżką w systemie plików RDZENIA — magazyn treści zapisuje ' +
  'bajty i oddaje ścieżkę — więc przeglądarka Operatora nie wczyta spod niego niczego, także ' +
  'gdy obraz powstał bez zarzutu. KOMENDA ODDAJĄCA BAJTY JUŻ W KONTRAKCIE JEST i jest jedna ' +
  'dla wszystkich modułów magazynu: zasób Design, plik Library i dokument Studia leżą w tym ' +
  'samym repozytorium, bo rodziny dokumentu, mediów i archiwum rozwiązują swoje zasoby przez ' +
  'nie. Brakuje jej UCHWYTU W RDZENIU — i to jest inne zdanie niż „nie ma czym", bo droga jest ' +
  'już opisana i czeka na dobudowę, a nie na projekt. Wykaz pod tym zdaniem mierzy ten stan ' +
  'sam, pytając rdzeń o jego komendy. Skutek dla Operatora nie zmienił się jeszcze ani trochę: ' +
  'zasób jest w Assets Panel, ma format i wymiary, a Preview Window pokazuje w jego miejscu ' +
  'zdanie o braku drogi po bajty. Wgranie jest drugą drogą, nie jedyną: wniesienie pliku ' +
  'wskazanego przez Operatora działa bez żadnego kanału modelu — te dwie drogi się uzupełniają.';

export function utworzPasekUczciwosci(kanal: Kanal): PasekUczciwosci {
  const katalog: KatalogOkien = utworzKatalogOkien(kanal, KOD_MODULU, KODY_OKIEN);
  const pokrycie: PokrycieKomend = utworzPokrycieKomend(kanal);

  const zdanie = document.createElement('p');
  zdanie.className = 'dn-pole-opis md-uczciwosc__opis';
  zdanie.textContent =
    'Zasób wizualny wchodzi do modułu dwiema drogami: wygenerowaniem kanałem obrazowym oraz ' +
    'wgraniem pliku przez Operatora. Pas mówi, czego przy zasobie mimo to nie ma. Rozróżnia ' +
    'przy tym dwa stany, które łatwo zlać w jeden: komendy nie ma w kontrakcie wcale, a komenda ' +
    'jest, lecz rdzeń nie ma dla niej uchwytu. Wszystkie liczby poniżej są mierzone odczytem ' +
    'z rdzenia, nie wpisane — dlatego zmieniają się same, gdy zmienia się produkt.';

  const rzadCzynnosci = document.createElement('div');
  rzadCzynnosci.className = 'md-braki__rzad';
  rzadCzynnosci.append(
    ...CZYNNOSCI_MODULU.map(([etykieta, komenda, czynnosc]) =>
      pokrycie.przycisk(etykieta, komenda, czynnosc),
    ),
  );

  const rzadPytan = document.createElement('div');
  rzadPytan.className = 'md-braki__rzad';
  rzadPytan.append(
    pozycja(
      'bajty-obrazu-bez-drogi',
      'Wygenerowany obraz jest w rdzeniu, a nie widać go w oknie',
      'Bajty zasobu bez drogi do przeglądarki',
      TRESC_OBRAZU,
    ),
    pozycja(
      'preview-window',
      'Preview Window — definicja wspólna ze Studiem',
      'Preview Window — definicja jedna, przypięcia dwa, pytanie otwarte',
      TRESC_PREVIEW_WINDOW,
    ),
  );

  const element = document.createElement('section');
  element.className = 'md-uczciwosc';
  element.setAttribute('aria-label', 'Czego moduł Design nie umie');
  element.append(
    zdanie,
    // Zdanie o katalogu okien liczy rozjazd: ile okien rdzeń przypisuje, ile moduł buduje i skąd.
    katalog.zdanieElement('dn-pole-opis md-uczciwosc__opis'),
    // Wykaz pokrycia mówi, czy rdzeń ma uchwyt każdej komendy, którą moduł woła.
    pokrycie.wykaz(KOMENDY_MODULU, 'md-uczciwosc__wykaz'),
    rzadCzynnosci,
    rzadPytan,
  );

  return {
    element,
    async odczytaj() {
      await Promise.all([katalog.odczytaj(), pokrycie.odczytaj()]);
    },
    zamknij() {
      katalog.zamknij();
      pokrycie.zamknij();
    },
  };
}

/** Pozycja pasa: klikalna zawsze, a jedyną jej reakcją jest powiedzenie prawdy dymkiem, nigdy odebranie klikalności. */
function pozycja(kod: string, napis: string, tytul: string, tresc: string): HTMLButtonElement {
  const element = document.createElement('button');
  element.type = 'button';
  element.className = 'dn-btn dn-btn--zarys dn-btn--sm';
  element.dataset['brak'] = kod;
  element.textContent = napis;
  element.addEventListener('click', () => pokazKomunikat({ tytul, tresc, waga: 'ostrz' }));
  return element;
}

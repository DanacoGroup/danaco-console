import './powloka.css';

import { pokazKomunikat } from '../aplikacja/komunikaty';
import { utworzKartySesji, type KartaSesji, type KartySesji, type SluchaczKarty } from './karty-sesji';
import { utworzNawigacjeModulow, type NawigacjaModulow, type SluchaczModulu } from './nawigacja-modulow';
import { utworzObszarRoboczy, type ObszarRoboczy } from './obszar-roboczy';
import { utworzPasekGorny, type PasekGorny } from './pasek-gorny';
import { SRODOWISKO_DOMYSLNE, type KluczSrodowiska } from './srodowiska';
import { utworzZrodloWyszukiwania } from './wyszukiwanie-globalne';
import type { ZaczepyPaska } from './zaczepy-paska';

/**
 * Powłoka środowiska — rama, w której mieszka wszystko inne.
 *
 * Jedna odpowiedzialność: złożenie czterech pasów w jeden układ i związanie
 * ich zdarzeniami. Żaden pas nie zna pozostałych; wiedzę o ich współpracy
 * trzyma wyłącznie ten plik.
 *
 * Układ czterech pasów:
 *   1. pasek górny 48 px na atramencie ramy — jedyny pas nieprzełączający się
 *      z motywem,
 *   2. poziome karty sesji o mechanice zakładek,
 *   3. pionowa, stała nawigacja modułów środowiska,
 *   4. obszar roboczy wypełniany przez inne widoki.
 *
 * Wybór modułu przeładowuje kartę sesji: tytuł karty czynnej równa się nazwie
 * otwartego modułu, kontekst paska pokazuje parę środowisko · moduł,
 * a obszar roboczy zapowiada moduł, którego okna wejdą w jego miejsce.
 */
export interface Powloka {
  /** Element powłoki; montuje go warstwa składająca aplikację. */
  element: HTMLElement;
  /** Pas 1 — pasek górny. */
  pasek: PasekGorny;
  /** Pas 2 — karty sesji. */
  karty: KartySesji;
  /** Pas 3 — boczna nawigacja modułów. */
  nawigacja: NawigacjaModulow;
  /** Pas 4 — obszar roboczy. */
  obszar: ObszarRoboczy;
  /** Przestawia powłokę na inne środowisko wraz z jego wykazem modułów. */
  ustawSrodowisko(klucz: KluczSrodowiska): void;
  /**
   * Podaje pasku czynności, których powłoka nie zna — otwieranie okien
   * platformy (Always On Display, Mobile, Modele, Konfiguracja…).
   *
   * Osobna metoda, a nie pole w `OpcjePowloki`, z tego samego powodu co
   * `nawigacja.podlaczZrodlo`: powłoka powstaje zanim istnieje połączenie
   * z rdzeniem, a czynności bez kanału podać się nie da. Wywołanie jest jedną
   * linijką dokładaną w widoku trasy i niczego w nim nie przestawia.
   *
   * Bez wywołania nic się nie psuje: przyciski obecności, pozycje menu profilu
   * i wyniki wyszukiwania zostają klikalne i mówią, że to pasek nie dostał
   * drogi — nie że okien nie ma (`zaczepy-paska.ts`).
   */
  podlaczZaczepy(zaczepy: ZaczepyPaska): void;
  /**
   * Zdejmuje nasłuchy dokumentu założone przez pasek (skrót Ctrl+K, zmiana
   * motywu, zwijanie menu po kliku poza nim).
   *
   * Powłoka aplikacji żyje tyle, co okno, więc w produkcie wywołanie nie
   * zachodzi ani razu — ale nasłuch na `document` bez drogi zdjęcia jest
   * wyciekiem, gdy powłoka powstaje wielokrotnie (podgląd, sprawdziany).
   */
  zamknij(): void;
  /** Zdarzenie: wybór karty sesji. */
  naWyborKarty(sluchacz: SluchaczKarty): void;
  /** Zdarzenie: zamknięcie karty sesji. */
  naZamknieciekarty(sluchacz: SluchaczKarty): void;
  /** Zdarzenie: naciśnięcie ＋, czyli założenie karty sesji. */
  naNowaKarte(sluchacz: SluchaczKarty): void;
  /** Zdarzenie: wybór modułu w bocznej nawigacji. */
  naWyborModulu(sluchacz: SluchaczModulu): void;
}

/** Ustawienia powłoki. Wszystkie mają wartość przyjmowaną domyślnie. */
export interface OpcjePowloki {
  /** Nazwa produktu na pasku górnym. */
  produkt?: string;
  /** Podpis Operatora pod awatarem. */
  operator?: string;
  /** Środowisko otwierane na starcie. */
  srodowisko?: KluczSrodowiska;
}

export function utworzPowloke(opcje: OpcjePowloki = {}): Powloka {
  const element = document.createElement('div');
  element.className = 'dn-srodowisko';

  const nawigacja = utworzNawigacjeModulow();
  const obszar = utworzObszarRoboczy();

  // Karta zakładana przyciskiem ＋ bierze tytuł z modułu wybranego w nawigacji —
  // tytuł karty jest nazwą otwartego modułu, nie napisem zastępczym.
  const karty = utworzKartySesji({
    tytulNowej: () => nawigacja.wybrana()?.nazwa ?? 'Nowa sesja',
  });

  /**
   * Czynności paska. Obiekt powstaje pusty i wypełnia go `podlaczZaczepy` —
   * kontrolki czytają z niego dopiero przy naciśnięciu (`zaczepy-paska.ts`).
   */
  const zaczepy: ZaczepyPaska = {};

  /**
   * Materiał wyszukiwania bierze się z tego, co powłoka już ma.
   *
   * To jedyne miejsce widzące naraz boczną nawigację (wykaz modułów pobrany
   * przez `module.list`) i pas kart sesji, więc źródło składa się tutaj — tak
   * samo jak tutaj wiąże się kontekst paska z wyborem modułu.
   *
   * Nowego odczytu z rdzenia nie ma: drugi `module.list` obok tego, który
   * zrobiła już nawigacja, byłby drugą prawdą o jednym wykazie, a dwie prawdy
   * rozjeżdżają się przy pierwszej zmianie po stronie rdzenia.
   */
  const zrodloWyszukiwania = utworzZrodloWyszukiwania({
    srodowisko: () => nawigacja.srodowisko(),
    karty: () => karty.wykaz(),
    wybierzModul: (klucz) => nawigacja.wybierz(klucz),
    wybierzKarte: (id) => karty.wybierz(id),
    zaczepy,
    naKomunikat: (tytul, tresc) => pokazKomunikat({ tytul, tresc, waga: 'info' }),
  });

  const pasek = utworzPasekGorny({
    ...(opcje.produkt === undefined ? {} : { produkt: opcje.produkt }),
    ...(opcje.operator === undefined ? {} : { operator: opcje.operator }),
    zaczepy,
    zrodloWyszukiwania,
  });

  const korpus = document.createElement('div');
  korpus.className = 'dn-srodowisko__korpus';
  korpus.append(nawigacja.element, obszar.element);

  element.append(pasek.element, karty.element, korpus);

  nawigacja.naWybor((pozycja, dane) => {
    element.dataset.srodowisko = dane.klucz;
    pasek.pokazKontekst(dane.nazwa, pozycja.nazwa);
    karty.przemianujCzynna(pozycja.nazwa);
    obszar.zapowiedz(pozycja, dane);
  });

  // Wybór karty sesji przywraca w obszarze zapowiedź modułu tej karty; póki
  // widoki modułów nie są osadzone, obszar pokazuje moduł wskazany nawigacją.
  karty.naWybor((karta: KartaSesji) => {
    element.dataset.karta = karta.id;
  });

  // Powłoka otwiera się z jedną kartą sesji. Tytułu nie nadaje się tutaj:
  // nada go za chwilę wybór modułu, bo tytuł karty jest nazwą modułu.
  karty.dodaj('Sesja');

  // Wykaz modułów i wszystkie związania ustawiają się dopiero teraz, gdy
  // słuchacze są już podpięci — wywołanie z konstruktora nawigacji ich nie miało.
  nawigacja.pokaz(opcje.srodowisko ?? SRODOWISKO_DOMYSLNE);
  pasek.ustawPowiadomienia(0);

  return {
    element,
    pasek,
    karty,
    nawigacja,
    obszar,

    ustawSrodowisko: (klucz) => nawigacja.pokaz(klucz),

    zamknij: () => pasek.zamknij(),

    podlaczZaczepy(nowe) {
      // Podmieniana jest zawartość, nie odniesienie: kontrolki paska trzymają ten
      // sam obiekt od montażu i czytają z niego przy każdym naciśnięciu.
      Object.assign(zaczepy, nowe);
    },

    naWyborKarty: (sluchacz) => karty.naWybor(sluchacz),
    naZamknieciekarty: (sluchacz) => karty.naZamkniecie(sluchacz),
    naNowaKarte: (sluchacz) => karty.naNowa(sluchacz),
    naWyborModulu: (sluchacz) => nawigacja.naWybor(sluchacz),
  };
}

export type { KartaSesji, KluczSrodowiska };

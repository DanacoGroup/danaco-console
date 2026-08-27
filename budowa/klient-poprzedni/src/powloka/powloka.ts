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
 * Powłoka środowiska — rama, w której mieszka wszystko inne: cztery pasy złożone w jeden wspólny
 * układ pracy.
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
  /** Podaje pasku czynności, których powłoka nie zna — otwieranie okien platformy poza kontraktem sesji. */
  podlaczZaczepy(zaczepy: ZaczepyPaska): void;
  /** Zdejmuje nasłuchy dokumentu założone przez pasek — skrót Ctrl+K, zmianę motywu, zwijanie menu. */
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

/** Ustawienia powłoki środowiska Danaco Console, wszystkie mają wartość przyjmowaną domyślnie, gdy ich nie podano. */
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

  // Karta zakładana przyciskiem ＋ bierze tytuł z modułu wybranego w nawigacji, nie napis zastępczy.
  const karty = utworzKartySesji({
    tytulNowej: () => nawigacja.wybrana()?.nazwa ?? 'Nowa sesja',
  });

  /** Czynności paska: obiekt powstaje pusty, a kontrolki czytają z niego dopiero przy naciśnięciu. */
  const zaczepy: ZaczepyPaska = {};

  /** Materiał wyszukiwania bierze się z tego, co powłoka już ma: nawigacji bocznej i pasa kart sesji. */
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

  // Wybór karty sesji przywraca w obszarze zapowiedź modułu tej karty, dopóki widoki nie są osadzone.
  karty.naWybor((karta: KartaSesji) => {
    element.dataset.karta = karta.id;
  });

  // Powłoka otwiera się z jedną kartą sesji; tytuł nada za chwilę wybór modułu, jako nazwę modułu.
  karty.dodaj('Sesja');

  // Wykaz modułów i wszystkie związania ustawiają się dopiero teraz, gdy słuchacze są już podpięci.
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
      // Podmieniana jest zawartość, nie odniesienie: kontrolki paska trzymają ten sam obiekt od montażu.
      Object.assign(zaczepy, nowe);
    },

    naWyborKarty: (sluchacz) => karty.naWybor(sluchacz),
    naZamknieciekarty: (sluchacz) => karty.naZamkniecie(sluchacz),
    naNowaKarte: (sluchacz) => karty.naNowa(sluchacz),
    naWyborModulu: (sluchacz) => nawigacja.naWybor(sluchacz),
  };
}

export type { KartaSesji, KluczSrodowiska };

import './pasek-gorny.css';

import { elementIkony, IKONA_PELNOKOLOROWA } from '../ikony/ikony';
import { utworzAkcjePaska, type AkcjePaska } from './akcje-paska';
import { opisNieczynnosci, wyjasnijNieczynnosc } from './nieczynne-w-pasku';
import { utworzWykazWynikow, type WykazWynikow } from './wykaz-wynikow';
import type { ZrodloWyszukiwania } from './wyszukiwanie-globalne';
import type { ZaczepyPaska } from './zaczepy-paska';

/**
 * Pas 1 powłoki — pasek górny wysokości 48 px na ramie kokpitu: tożsamość produktu, kontekst pracy,
 * pole poleceń i grupa akcji.
 */
export interface PasekGorny {
  /** Pasek montowany jako pierwszy wiersz powłoki. */
  element: HTMLElement;
  /** Pole poleceń — udostępnione, by inny moduł mógł podpiąć wykonanie. */
  polePolecen: HTMLInputElement;
  /** Wypisuje bieżące środowisko i moduł w kontekście paska. */
  pokazKontekst(srodowisko: string, modul: string): void;
  /** Ustawia licznik powiadomień. */
  ustawPowiadomienia(ile: number): void;
  /** Grupa akcji po prawej stronie. */
  akcje: AkcjePaska;
  /** Zdejmuje nasłuchy dokumentu. Obowiązkowe przy zejściu powłoki. */
  zamknij(): void;
}

/** Ustawienia paska: nazwa produktu, podpis Operatora, czynności obsługiwane przez pasek oraz materiał wyszukiwania globalnego. */
export interface OpcjePaska {
  produkt?: string;
  operator?: string;
  /** Czynności paska nieznane temu plikowi — dostarczane z zewnątrz jako zaczepy przy montażu powłoki. */
  zaczepy?: ZaczepyPaska;
  /** Skąd wyszukiwanie bierze materiał; pominięcie oznacza brak materiału do przeszukania. */
  zrodloWyszukiwania?: ZrodloWyszukiwania;
}

/** Wartości nazwy produktu i podpisu Operatora przyjmowane, gdy wywołanie budujące pasek ich nie podaje. */
const PRODUKT = 'Danaco Console';
const OPERATOR = 'Operator';

export function utworzPasekGorny(opcje: OpcjePaska = {}): PasekGorny {
  const produkt = opcje.produkt ?? PRODUKT;

  const element = document.createElement('header');
  element.className = 'dn-pasek dn-pasek-gorny';

  const tozsamosc = document.createElement('div');
  tozsamosc.className = 'dn-pasek-godlo dn-pasek-gorny__godlo';

  const godlo = elementIkony(IKONA_PELNOKOLOROWA, {
    rozmiar: 24,
    etykieta: 'Godło Danaco',
  });

  const nazwa = document.createElement('span');
  nazwa.className = 'dn-pasek-gorny__produkt';
  nazwa.textContent = produkt;

  tozsamosc.append(godlo, nazwa);

  const kontekst = document.createElement('p');
  kontekst.className = 'dn-pasek-gorny__kontekst';

  const { pole: polePolecen, odepnij: odepnijSkrotPola } = utworzPolePolecen();

  const skrot = document.createElement('kbd');
  skrot.className = 'dn-kod dn-kod--wiersz dn-pasek-gorny__skrot';
  skrot.textContent = 'Ctrl K';

  const pole = document.createElement('div');
  pole.className = 'dn-pasek-gorny__pole';
  pole.append(polePolecen, skrot);

  const wykaz = zwiazWyszukiwanie(polePolecen, pole, opcje.zrodloWyszukiwania);

  const akcje = utworzAkcjePaska(opcje.operator ?? OPERATOR, opcje.zaczepy ?? {});

  element.append(tozsamosc, kontekst, pole, akcje.element);

  return {
    element,
    polePolecen,
    akcje,

    pokazKontekst(srodowisko, modul) {
      kontekst.textContent = `${srodowisko} · ${modul}`;
      kontekst.title = `Środowisko ${srodowisko}, moduł ${modul}`;
    },

    ustawPowiadomienia: (ile) => akcje.ustawPowiadomienia(ile),

    zamknij() {
      odepnijSkrotPola();
      wykaz?.zamknij();
      akcje.zamknij();
    },
  };
}

/**
 * Pole poleceń wraz z jego czynnymi zachowaniami.
 *
 * Ctrl+K prowadzi do pola ognisko, więc podpowiedź skrótu obok pola mówi
 * prawdę. Zatwierdzenie treści daje odpowiedź zamiast ciszy — czym dokładnie,
 * rozstrzyga `zwiazWyszukiwanie`.
 */
function utworzPolePolecen(): { pole: HTMLInputElement; odepnij: () => void } {
  const pole = document.createElement('input');
  pole.type = 'search';
  pole.className = 'dn-pasek-szukaj dn-pasek-gorny__polecenie';
  pole.placeholder = 'Polecenie, moduł albo wyszukanie';
  pole.setAttribute('aria-label', 'Wyszukiwanie globalne i pole poleceń');
  pole.autocomplete = 'off';

  const skrotOgniska = (zdarzenie: KeyboardEvent): void => {
    if (!zdarzenie.ctrlKey || zdarzenie.key.toLowerCase() !== 'k') return;
    zdarzenie.preventDefault();
    pole.focus();
    pole.select();
  };
  document.addEventListener('keydown', skrotOgniska);

  return {
    pole,
    odepnij: () => document.removeEventListener('keydown', skrotOgniska),
  };
}

/**
 * Wiąże pole paska z wykazem wyników — pięć zachowań klawiatury, przy których ognisko stale
 * zostaje w polu.
 */
function zwiazWyszukiwanie(
  pole: HTMLInputElement,
  gospodarz: HTMLElement,
  zrodlo: ZrodloWyszukiwania | undefined,
): WykazWynikow | undefined {
  if (zrodlo === undefined) {
    pole.addEventListener('keydown', (zdarzenie: KeyboardEvent) => {
      if (zdarzenie.key !== 'Enter') return;
      zdarzenie.preventDefault();
      wyjasnijNieczynnosc('polecenie');
    });
    pole.setAttribute('aria-label', opisNieczynnosci('polecenie', 'Pole poleceń'));
    return undefined;
  }

  const wykaz = utworzWykazWynikow({
    zrodlo,
    // Wybór kończy szukanie: pole wraca puste, bo fraza opisywała drogę, którą Operator już przeszedł.
    naWybor: () => {
      pole.value = '';
    },
  });
  gospodarz.append(wykaz.element);

  pole.addEventListener('input', () => wykaz.ustawFraze(pole.value));
  pole.addEventListener('focus', () => wykaz.otworz());

  pole.addEventListener('keydown', (zdarzenie: KeyboardEvent) => {
    if (zdarzenie.key === 'ArrowDown' || zdarzenie.key === 'ArrowUp') {
      zdarzenie.preventDefault();
      wykaz.przesun(zdarzenie.key === 'ArrowDown' ? 1 : -1);
      return;
    }
    if (zdarzenie.key === 'Escape') {
      // Escape zwija wykaz i zostawia ognisko w polu — Operator pisze dalej zamiast szukać ogniska.
      if (!wykaz.otwarty()) return;
      zdarzenie.stopPropagation();
      wykaz.zamknij();
      return;
    }
    if (zdarzenie.key !== 'Enter') return;
    zdarzenie.preventDefault();
    if (wykaz.zatwierdz()) return;
    // Nie było czego zatwierdzić — jedyne miejsce, gdzie pole mówi o niezbudowanym centrum poleceń.
    wyjasnijNieczynnosc('polecenie');
  });

  return wykaz;
}

/**
 * Ekrany okna dostępu do konta — dziesięć odsłon: logowanie, niepowodzenie
 * logowania, zwłoka nałożona przez rdzeń, zakładanie konta, potwierdzenie
 * adresu, konto założone bez potwierdzenia oraz cztery odsłony odzyskiwania
 * dostępu.
 */

import { el, tekst, type Dziecko } from '../narzedzia.ts';
import type { OdslonaDostepu } from '../przebieg.ts';
import { baner } from '../skladniki/baner.ts';
import { frazaNawigacyjna } from '../skladniki/fraza-nawigacyjna.ts';
import { krokiOdzyskiwania } from '../skladniki/kroki-odzyskiwania.ts';
import { metodyLogowania } from '../skladniki/metody-logowania.ts';
import { miernikSily } from '../skladniki/miernik-sily.ts';
import { naglowekEkranu } from '../skladniki/naglowek-ekranu.ts';
import { pasDzialan, type Czynnosc } from '../skladniki/pas-dzialan.ts';
import { poleHasla } from '../skladniki/pole-hasla.ts';
import { poleKodu } from '../skladniki/pole-kodu.ts';
import { poleSesji } from '../skladniki/pole-sesji.ts';
import { poleTekstowe } from '../skladniki/pole-tekstowe.ts';
import { zakladkiPigulki, type Pigulka } from '../skladniki/zakladki-pigulki.ts';

/**
 * Ile znaków ma droga potwierdzenia i jak długo obowiązuje — obie wartości
 * zmierzone na rdzeniu, nie wzięte z wartości przykładowych prototypu.
 */
const ZNAKOW_DROGI = 6;
const WAZNOSC_DROGI_S = 3600;
const WAZNOSC_DROGI_MIN = WAZNOSC_DROGI_S / 60;

/**
 * Licznik zwłoki startuje od zera i wartość dostaje z pomiaru, nie ze stałej.
 * Rdzeń nie odmawia kolejnej próby — nakłada na nią zwłokę, a jej długość
 * okno mierzy czasem trwania próby poprzedniej, bo odmowa jej nie niesie.
 */
const ZWLOKA_Z_POMIARU = 0;

/** Klasa czynności głównej — wyjścia z odsłony, nadawana przyciskowi prowadzącemu dalej w torze dostępu. */
const KLASA_GLOWNA = 'dn-btn dn-btn--sygnal';

const ZAMKNIJ: Czynnosc = { klucz: 'dzialania.zamknijAplikacje', komunikat: 'zamkniecie' };

function panel(
  nazwa: OdslonaDostepu,
  aktywny: boolean,
  pigulka: Pigulka,
  nazwaDostepna: string,
  dzieci: Dziecko[],
): HTMLElement {
  return el(
    'div',
    {
      klasa: 'we-panel',
      id: `s-${nazwa}`,
      role: 'tabpanel',
      tabindex: '0',
      /* Obszar o roli `tabpanel` musi mieć nazwę — bez niej czytnik ekranu
         ogłasza „panel” i nic więcej. Nazwa bierze się z tytułu odsłony. */
      'aria-label': tekst(nazwaDostepna),
      dane: {
        widok: nazwa,
        'grupa-widoku': 'stan',
        'widok-aktywny': aktywny ? 'tak' : 'nie',
        pigulka,
      },
    },
    [zakladkiPigulki({ wybrana: pigulka }), ...dzieci],
  );
}

/** Miejsce na komunikaty jednej odsłony; przebieg wstawia tam usterki zgłoszone przez rdzeń albo przez klienta. */
function komunikaty(odslona: OdslonaDostepu): HTMLElement {
  return el('div', {
    klasa: 'we-komunikaty',
    'aria-live': 'polite',
    dane: { komunikaty: odslona },
  });
}

/**
 * Zdanie z adresem. Adres stoi krojem maszynowym — łatwiej porównać go ze
 * skrzynką znak po znaku — więc zdanie składa się z węzłów, nie z łańcucha.
 */
function zdanieZAdresem(klucz: string, odslona: OdslonaDostepu): Dziecko[] {
  const czesci = tekst(klucz).split('{adres}');
  return [
    czesci[0] ?? '',
    el('span', { klasa: 'au-kod-adres', dane: { adres: odslona } }),
    czesci[1] === undefined
      ? ''
      : czesci[1].replace('{minuty}', String(WAZNOSC_DROGI_MIN)),
  ];
}

function polaLogowania(przedrostek: string, blad: boolean): HTMLElement {
  return el('div', { klasa: 'au-pola' }, [
    poleTekstowe({
      etykieta: 'dostep.logowanie.login',
      id: `${przedrostek}-login`,
      uzupelnij: 'username',
      bledne: blad,
      opisuje: blad ? 'blad-logowania' : null,
    }),
    poleHasla({
      etykieta: 'dostep.logowanie.haslo',
      id: `${przedrostek}-haslo`,
      opis: blad ? 'dostep.logowanieBlad.capsLock' : null,
      bledne: blad,
      opisuje: blad ? `blad-logowania ${przedrostek}-haslo-opis` : null,
    }),
  ]);
}

function resetHasla(): HTMLElement {
  return frazaNawigacyjna({
    pytanie: 'dostep.logowanie.reset.pytanie',
    czynnosc: 'dostep.logowanie.reset.czynnosc',
    cel: 'odzyskiwanie-adres',
  });
}

export function ekranDostepu(): HTMLElement[] {
  return [
    /* ── logowanie ─────────────────────────────────────────────────────── */
    panel('logowanie', true, 'logowanie', 'dostep.logowanie.tytul', [
      naglowekEkranu({ tytul: 'dostep.logowanie.tytul', lid: 'dostep.logowanie.lid' }),
      komunikaty('logowanie'),
      polaLogowania('log', false),
      poleSesji({ id: 'log-sesja' }),
      resetHasla(),
      ...metodyLogowania({ dostepne: ['email'] }),
      frazaNawigacyjna({
        pytanie: 'dostep.logowanie.fraza.pytanie',
        czynnosc: 'dostep.logowanie.fraza.czynnosc',
        cel: 'rejestracja',
      }),
    ]),

    /* ── logowanie, dane nierozpoznane ─────────────────────────────────── */
    panel('logowanie-blad', false, 'logowanie', 'dostep.logowanie.tytul', [
      naglowekEkranu({ tytul: 'dostep.logowanie.tytul' }),
      el('div', { klasa: 'we-komunikaty', 'aria-live': 'polite' }, [
        baner({
          rodzaj: 'blad',
          ikona: 'ostrzezenie',
          id: 'blad-logowania',
          glowa: 'dostep.logowanieBlad.baner.glowa',
          tresc: 'dostep.logowanieBlad.baner.tresc',
        }),
      ]),
      komunikaty('logowanie-blad'),
      polaLogowania('blad', true),
      resetHasla(),
      ...metodyLogowania({ dostepne: ['email'] }),
      frazaNawigacyjna({ klucz: 'dostep.logowanieBlad.fraza' }),
    ]),

    /* ── zwłoka nałożona na kolejną próbę; NIE zapora ─────────────────── */
    panel('logowanie-wstrzymane', false, 'logowanie', 'dostep.logowanieWstrzymane.tytul', [
      naglowekEkranu({
        tytul: 'dostep.logowanieWstrzymane.tytul',
        lid: 'dostep.logowanieWstrzymane.lid',
      }),
      baner({
        rodzaj: 'ostrzezenie',
        ikona: 'zegar',
        glowa: 'dostep.logowanieWstrzymane.baner.glowa',
        tresc: 'dostep.logowanieWstrzymane.baner.tresc',
        daneGlowy: { czas: '{odliczanie}' },
        odliczanie: ZWLOKA_Z_POMIARU,
      }),
      /* Pola zostają na ekranie: zwłoka trwa kilka sekund, a odsłona bez
         formularza zostawiałaby Operatora na widoku, z którego nie ma dokąd
         przejść. Odliczanie stoi w banerze nad nimi. */
      polaLogowania('wstrzymane', false),
      frazaNawigacyjna({
        czynnosc: 'dostep.logowanieWstrzymane.odzyskaj',
        cel: 'odzyskiwanie-adres',
      }),
    ]),

    /* ── zakładanie konta ──────────────────────────────────────────────── */
    panel('rejestracja', false, 'rejestracja', 'dostep.rejestracja.tytul', [
      naglowekEkranu({ tytul: 'dostep.rejestracja.tytul', lid: 'dostep.rejestracja.lid' }),
      komunikaty('rejestracja'),
      el('div', { klasa: 'au-pola' }, [
        el('div', { klasa: 'au-para' }, [
          poleTekstowe({ etykieta: 'dostep.rejestracja.login', id: 'rej-login', uzupelnij: 'username' }),
          poleTekstowe({
            etykieta: 'dostep.rejestracja.email',
            id: 'rej-email',
            typ: 'email',
            uzupelnij: 'email',
          }),
        ]),
        el('div', { klasa: 'au-para' }, [
          el('div', {}, [
            poleHasla({
              etykieta: 'dostep.rejestracja.haslo',
              id: 'rej-haslo',
              uzupelnij: 'new-password',
            }),
            miernikSily({ dla: 'rej-haslo' }),
          ]),
          poleHasla({
            etykieta: 'dostep.rejestracja.hasloPowtorz',
            id: 'rej-haslo-2',
            uzupelnij: 'new-password',
          }),
        ]),
      ]),
      poleSesji({ id: 'rej-sesja' }),
      frazaNawigacyjna({
        pytanie: 'dostep.rejestracja.fraza.pytanie',
        czynnosc: 'dostep.rejestracja.fraza.czynnosc',
        cel: 'logowanie',
      }),
    ]),

    /* ── potwierdzenie adresu po założeniu konta ───────────────────────── */
    panel('kod', false, 'rejestracja', 'dostep.kod.tytul', [
      naglowekEkranu({ tytul: 'dostep.kod.tytul', lidWezly: zdanieZAdresem('dostep.kod.lid', 'kod') }),
      komunikaty('kod'),
      ...poleKodu({
        znakow: ZNAKOW_DROGI,
        odliczanie: WAZNOSC_DROGI_S,
        czynnosc: 'wklej',
        grupa: 'kod',
      }),
      el('div', { klasa: 'we-komunikaty', 'aria-live': 'polite' }, [
        baner({
          rodzaj: 'informacja',
          ikona: 'informacja',
          glowa: 'dostep.kod.pomoc.glowa',
          tresc: 'dostep.kod.pomoc.tresc',
          dane: { nadawca: tekst('dostep.kod.pomocDane.nadawca') },
        }),
      ]),
      frazaNawigacyjna({ czynnosc: 'dostep.kod.zmienAdres', cel: 'rejestracja' }),
    ]),

    /* ── konto założone bez konta nadawczego ──────────────────────────── */
    panel(
      'konto-bez-potwierdzenia',
      false,
      'rejestracja',
      'dostep.kontoBezPotwierdzenia.tytul',
      [
        naglowekEkranu({
          tytul: 'dostep.kontoBezPotwierdzenia.tytul',
          lid: 'dostep.kontoBezPotwierdzenia.lid',
        }),
        /* Ostrzeżenie jest WYMAGANE, nie zalecane: adres jest jedyną drogą
           odzyskania konta, a Operator, który nie przeczyta tego przy
           rejestracji, dowie się w dniu, w którym będzie tego potrzebował. */
        baner({
          rodzaj: 'ostrzezenie',
          ikona: 'ostrzezenie',
          glowa: 'dostep.kontoBezPotwierdzenia.baner.glowa',
          daneGlowy: { adres: '' },
          trescWezly: [
            tekst('dostep.kontoBezPotwierdzenia.baner.tresc'),
          ],
        }),
        frazaNawigacyjna({ klucz: 'dostep.kontoBezPotwierdzenia.nota' }),
      ],
    ),

    /* ── odzyskiwanie: adres ───────────────────────────────────────────── */
    panel('odzyskiwanie-adres', false, 'logowanie', 'dostep.odzyskiwanie.adres.tytul', [
      naglowekEkranu({
        tytul: 'dostep.odzyskiwanie.adres.tytul',
        poTytule: [krokiOdzyskiwania({ biezacy: 1 })],
        lid: 'dostep.odzyskiwanie.adres.lid',
      }),
      komunikaty('odzyskiwanie-adres'),
      el('div', { klasa: 'au-pola' }, [
        poleTekstowe({
          etykieta: 'dostep.odzyskiwanie.adres.pole',
          id: 'odz-email',
          typ: 'email',
          uzupelnij: 'email',
        }),
      ]),
      el('div', { klasa: 'we-komunikaty', 'aria-live': 'polite' }, [
        baner({
          rodzaj: 'informacja',
          ikona: 'tarcza',
          glowa: 'dostep.odzyskiwanie.adres.ostrzezenie.glowa',
          tresc: 'dostep.odzyskiwanie.adres.ostrzezenie.tresc',
        }),
      ]),
      frazaNawigacyjna({ czynnosc: 'dostep.odzyskiwanie.powrot', cel: 'logowanie' }),
    ]),

    /* ── odzyskiwanie: droga z listu ───────────────────────────────────── */
    panel('odzyskiwanie-kod', false, 'logowanie', 'dostep.kod.tytul', [
      naglowekEkranu({
        tytul: 'dostep.kod.tytul',
        poTytule: [krokiOdzyskiwania({ biezacy: 2 })],
        lidWezly: zdanieZAdresem('dostep.kod.lid', 'odzyskiwanie-kod'),
      }),
      komunikaty('odzyskiwanie-kod'),
      ...poleKodu({
        znakow: ZNAKOW_DROGI,
        odliczanie: WAZNOSC_DROGI_S,
        czynnosc: 'ponow',
        grupa: 'odzyskiwanie',
      }),
      baner({
        rodzaj: 'informacja',
        ikona: 'tarcza',
        glowa: 'dostep.kod.ostrzezenie.glowa',
        tresc: 'dostep.kod.ostrzezenie.tresc',
      }),
      frazaNawigacyjna({ czynnosc: 'dostep.kod.zmienAdres', cel: 'odzyskiwanie-adres' }),
    ]),

    /* ── odzyskiwanie: nowe hasło ──────────────────────────────────────── */
    panel('odzyskiwanie-haslo', false, 'logowanie', 'dostep.odzyskiwanie.haslo.tytul', [
      naglowekEkranu({
        tytul: 'dostep.odzyskiwanie.haslo.tytul',
        poTytule: [krokiOdzyskiwania({ biezacy: 3 })],
        lid: 'dostep.odzyskiwanie.haslo.lid',
      }),
      komunikaty('odzyskiwanie-haslo'),
      el('div', { klasa: 'au-pola' }, [
        el('div', { klasa: 'au-para' }, [
          poleHasla({
            etykieta: 'dostep.odzyskiwanie.haslo.nowe',
            id: 'odz-haslo',
            uzupelnij: 'new-password',
          }),
          poleHasla({
            etykieta: 'dostep.odzyskiwanie.haslo.powtorz',
            id: 'odz-haslo-2',
            uzupelnij: 'new-password',
          }),
          miernikSily({ dla: 'odz-haslo' }),
        ]),
      ]),
      poleSesji({ id: 'odz-sesja' }),
      frazaNawigacyjna({ czynnosc: 'dostep.odzyskiwanie.powrot', cel: 'logowanie' }),
    ]),

    /* ── zwłoka nałożona na kolejne wysłanie; NIE zapora ──────────────── */
    panel(
      'odzyskiwanie-wstrzymane',
      false,
      'logowanie',
      'dostep.odzyskiwanieWstrzymane.tytul',
      [
        naglowekEkranu({
          tytul: 'dostep.odzyskiwanieWstrzymane.tytul',
          poTytule: [krokiOdzyskiwania({ biezacy: 1 })],
          lid: 'dostep.odzyskiwanieWstrzymane.lid',
        }),
        baner({
          rodzaj: 'ostrzezenie',
          ikona: 'zegar',
          glowa: 'dostep.odzyskiwanieWstrzymane.baner.glowa',
          tresc: 'dostep.odzyskiwanieWstrzymane.baner.tresc',
          dane: { minuty: WAZNOSC_DROGI_MIN },
          daneGlowy: { czas: '{odliczanie}' },
          odliczanie: ZWLOKA_Z_POMIARU,
        }),
        frazaNawigacyjna({ czynnosc: 'dostep.odzyskiwanie.powrot', cel: 'logowanie' }),
      ],
    ),
  ];
}

/**
 * Pasy działań odsłon dostępu. Czynność główna albo prowadzi do kolejnej
 * odsłony, albo wywołuje komendę rdzenia; bez jednego i drugiego przycisk
 * jest martwy.
 */
export function pasyDostepu(): HTMLElement[] {
  function pas(widok: OdslonaDostepu, aktywny: boolean, klucz: string, czynnosc: string): HTMLElement {
    return pasDzialan({
      widok,
      aktywny,
      czynnosci: [ZAMKNIJ, { klucz, klasa: KLASA_GLOWNA, czynnosc }],
    });
  }
  return [
    pas('logowanie', true, 'dzialania.zaloguj', 'zaloguj'),
    pas('logowanie-blad', false, 'dzialania.zalogujPonownie', 'zaloguj-ponownie'),
    pas('rejestracja', false, 'dzialania.utworzKonto', 'zarejestruj'),
    pas('kod', false, 'dzialania.potwierdzKonto', 'potwierdz-adres'),
    pas('konto-bez-potwierdzenia', false, 'dostep.kontoBezPotwierdzenia.wejdz', 'wejdz-bez-potwierdzenia'),
    pas('odzyskiwanie-adres', false, 'dzialania.wyslijKod', 'popros-o-odzyskanie'),
    pas('odzyskiwanie-kod', false, 'dzialania.potwierdzDroge', 'przejdz-do-hasla'),
    pas('odzyskiwanie-haslo', false, 'dzialania.potwierdzHaslo', 'ustaw-nowe-haslo'),
    /* Odsłony wstrzymania nie mają czynności głównej — nie ma czego wykonać,
       póki godzina nie minie. Zostaje wyjście z programu. */
    pas('logowanie-wstrzymane', false, 'dzialania.zaloguj', 'zaloguj'),
    pasDzialan({ widok: 'odzyskiwanie-wstrzymane', czynnosci: [ZAMKNIJ] }),
  ];
}

/**
 * EKRANY OKNA URUCHOMIENIA.
 *
 * Trzy odsłony tego samego okna: łączenie, powrót z ważnym tokenem, błąd
 * połączenia. Leżą razem, bo dzielą całą oprawę — belkę, kolumnę tożsamości
 * i wykaz czterech etapów. Różni je stan etapów i to, co stoi pod wykazem.
 *
 * Łączenie i powrót z tokenem nie mają czynności głównej: przechodzą dalej
 * same, gdy rdzeń odpowie. Błąd ją ma, bo tam jest co rozstrzygnąć.
 */

import { el } from '../narzedzia.ts';
import type { OdslonaUruchomienia, StanEtapu } from '../przebieg.ts';
import { baner } from '../skladniki/baner.ts';
import { frazaNawigacyjna } from '../skladniki/fraza-nawigacyjna.ts';
import { listaEtapow } from '../skladniki/lista-etapow.ts';
import { naglowekEkranu } from '../skladniki/naglowek-ekranu.ts';
import { pasDzialan, type Czynnosc } from '../skladniki/pas-dzialan.ts';

/** Ile sekund okno czeka, zanim samo ponowi próbę połączenia. */
const PONOWIENIE_S = 15;

/** Klasa czynności głównej — jedynego wyjścia z odsłony. */
const KLASA_GLOWNA = 'dn-btn dn-btn--sygnal';

/** Czynność obecna w każdej odsłonie przed uwierzytelnieniem. */
const ZAMKNIJ: Czynnosc = { klucz: 'dzialania.zamknijAplikacje', komunikat: 'zamkniecie' };

function panel(nazwa: OdslonaUruchomienia, aktywny: boolean, dzieci: HTMLElement[]): HTMLElement {
  return el(
    'div',
    {
      klasa: 'we-panel',
      id: `s-${nazwa}`,
      role: 'tabpanel',
      tabindex: '0',
      dane: {
        widok: nazwa,
        'grupa-widoku': 'wariant',
        'widok-aktywny': aktywny ? 'tak' : 'nie',
      },
    },
    dzieci,
  );
}

/** Wykaz etapów łączenia; treść odświeża się z przebiegu, oprawa zostaje. */
function wykazEtapow(stany: StanEtapu[], miary: string[]): HTMLElement {
  return listaEtapow({
    nazwy: 'uruchomienie.etapy',
    miary: 'uruchomienie.stany',
    obszar: 'uruchomienie.laczenie.tytul',
    stany,
    wartosciMiar: miary.map((klucz) => ({ klucz })),
  });
}

export function ekranUruchomienia(): HTMLElement[] {
  return [
    panel('w-laczenie', true, [
      naglowekEkranu({
        nadtytul: 'uruchomienie.nadtytul',
        tytul: 'uruchomienie.laczenie.tytul',
        lid: 'uruchomienie.laczenie.lid',
      }),
      el('div', { dane: { 'etapy-laczenia': 'w-laczenie' } }, [
        wykazEtapow(
          ['pracuje', 'oczekuje', 'oczekuje', 'oczekuje'],
          ['wToku', 'oczekuje', 'oczekuje', 'oczekuje'],
        ),
      ]),
    ]),

    panel('w-token', false, [
      naglowekEkranu({
        nadtytul: 'uruchomienie.nadtytul',
        tytul: 'uruchomienie.token.tytul',
        lid: 'uruchomienie.token.lid',
      }),
      el('div', { dane: { 'etapy-laczenia': 'w-token' } }, [
        wykazEtapow(
          ['gotowy', 'gotowy', 'gotowy', 'pracuje'],
          ['nawiazane', 'zgodna', 'zaufane', 'wToku'],
        ),
      ]),
      baner({
        rodzaj: 'sukces',
        ikona: 'tarcza',
        glowa: 'uruchomienie.token.baner.glowa',
        tresc: 'uruchomienie.token.baner.tresc',
      }),
      frazaNawigacyjna({ klucz: 'uruchomienie.token.fraza' }),
    ]),

    panel('w-blad', false, [
      naglowekEkranu({
        nadtytul: 'uruchomienie.nadtytul',
        tytul: 'uruchomienie.blad.tytul',
        lid: 'uruchomienie.blad.lid',
      }),
      el('div', { dane: { 'etapy-laczenia': 'w-blad' } }, [
        wykazEtapow(
          ['blad', 'oczekuje', 'oczekuje', 'oczekuje'],
          ['nieudane', 'oczekuje', 'oczekuje', 'oczekuje'],
        ),
      ]),
      // Miejsce na odmowę rdzenia albo na rozjazd wersji protokołu; puste,
      // dopóki przebieg nie ma czego w nie wstawić.
      el('div', { klasa: 'we-komunikaty', 'aria-live': 'polite', dane: { komunikaty: 'w-blad' } }),
      baner({
        rodzaj: 'ostrzezenie',
        ikona: 'ostrzezenie',
        glowa: 'uruchomienie.blad.baner.glowa',
        tresc: 'uruchomienie.blad.baner.tresc',
        odliczanieTresci: PONOWIENIE_S,
        postacOdliczania: 'sekundy',
      }),
      frazaNawigacyjna({ klucz: 'uruchomienie.blad.fraza' }),
    ]),
  ];
}

export function pasyUruchomienia(): HTMLElement[] {
  return [
    pasDzialan({ widok: 'w-laczenie', aktywny: true, czynnosci: [ZAMKNIJ] }),
    pasDzialan({ widok: 'w-token', czynnosci: [ZAMKNIJ] }),
    pasDzialan({
      widok: 'w-blad',
      czynnosci: [
        ZAMKNIJ,
        { klucz: 'dzialania.ustawieniaPolaczenia', komunikat: 'ustawienia' },
        {
          klucz: 'dzialania.sprobujPonownie',
          klasa: KLASA_GLOWNA,
          czynnosc: 'ponow-polaczenie',
          ikona: 'ponow',
        },
      ],
    }),
  ];
}

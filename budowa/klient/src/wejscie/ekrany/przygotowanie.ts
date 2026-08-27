/**
 * EKRAN OKNA PRZYGOTOWANIA ŚRODOWISKA PRACY.
 *
 * Trzeci etap wejścia: konto jest już rozpoznane, program odtwarza stan pracy
 * sprzed zamknięcia. Okno nie pyta o nic — wykaz etapów mówi, co się dzieje,
 * a pas działań daje dwa wyjścia: pominąć przywracanie albo przerwać i się
 * wylogować.
 *
 * Kolumna tożsamości niesie tu animację powłok zamiast wykazu zdań: na tym
 * etapie użytkownik nie wybiera już programu, tylko czeka, aż się złoży.
 *
 * Stany i miary pięciu etapów pochodzą z przebiegu, czyli z odpowiedzi rdzenia.
 * Etap, dla którego droga wejścia nie ma komendy, zostaje w stanie oczekiwania
 * — etap oznaczony jako gotowy bez pomiaru mówiłby nieprawdę.
 */

import { el } from '../narzedzia.ts';
import type { PostepPrzygotowania } from '../przebieg.ts';
import { naglowekEkranu } from '../skladniki/naglowek-ekranu.ts';
import { pasDzialan } from '../skladniki/pas-dzialan.ts';
import { pasekPostepu } from '../skladniki/pasek-postepu.ts';
import { listaEtapow } from '../skladniki/lista-etapow.ts';

/** Wykaz etapów przygotowania odbudowywany przy każdej zmianie stanu. */
export function wykazPrzygotowania(postep: PostepPrzygotowania): HTMLElement {
  return listaEtapow({
    nazwy: 'przygotowanie.etapy',
    miary: 'przygotowanie.miary',
    obszar: 'przygotowanie.obszarEtapow',
    stany: postep.stany,
    wartosciMiar: postep.miary,
  });
}

/** Pasek postępu odbudowywany przy każdej zmianie stanu. */
export function paskiPrzygotowania(postep: PostepPrzygotowania): HTMLElement {
  return pasekPostepu({
    etykieta: 'przygotowanie.postep.etykieta',
    opisPaska: 'przygotowanie.postep.opisPaska',
    wartosc: postep.wartosc,
  });
}

export function ekranPrzygotowania(postep: PostepPrzygotowania): HTMLElement {
  return el('div', { klasa: 'we-panel', dane: { widok: 'przygotowanie' } }, [
    naglowekEkranu({
      nadtytul: 'przygotowanie.nadtytul',
      tytul: 'przygotowanie.tytul',
      lid: 'przygotowanie.lid',
    }),
    el('div', { dane: { 'etapy-przygotowania': true } }, [wykazPrzygotowania(postep)]),
    el('div', { dane: { 'postep-przygotowania': true } }, [paskiPrzygotowania(postep)]),
    el('div', {
      klasa: 'we-komunikaty',
      'aria-live': 'polite',
      dane: { komunikaty: 'przygotowanie' },
    }),
  ]);
}

export function pasPrzygotowania(): HTMLElement {
  return pasDzialan({
    czynnosci: [
      { klucz: 'dzialania.pominPrzywracanie', komunikat: 'pominiecie' },
      { klucz: 'dzialania.przerwijIWyloguj', czynnosc: 'przerwij-i-wyloguj' },
    ],
  });
}

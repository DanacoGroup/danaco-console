/**
 * Ekran okna przygotowania środowiska pracy — trzeci etap wejścia: konto
 * jest już rozpoznane, program odtwarza stan pracy sprzed zamknięcia. Okno
 * nie pyta o nic, wykaz etapów mówi, co się dzieje.
 */

import { el } from '../narzedzia.ts';
import type { PostepPrzygotowania } from '../przebieg.ts';
import { naglowekEkranu } from '../skladniki/naglowek-ekranu.ts';
import { pasDzialan } from '../skladniki/pas-dzialan.ts';
import { pasekPostepu } from '../skladniki/pasek-postepu.ts';
import { listaEtapow } from '../skladniki/lista-etapow.ts';

/** Wykaz etapów przygotowania odbudowywany przy każdej zmianie stanu, na podstawie miar odczytanych z odpowiedzi rdzenia. */
export function wykazPrzygotowania(postep: PostepPrzygotowania): HTMLElement {
  return listaEtapow({
    nazwy: 'przygotowanie.etapy',
    miary: 'przygotowanie.miary',
    obszar: 'przygotowanie.obszarEtapow',
    stany: postep.stany,
    wartosciMiar: postep.miary,
  });
}

/** Pasek postępu odbudowywany przy każdej zmianie stanu, dobiegający stu procent dopiero po ostatnim etapie. */
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

/** Znacznik kontrolki ponowienia; montaż nastawia po nim jej widoczność, pokazując ją wyłącznie przy etapie nieudanym. */
export const PONOWIENIE = 'przyg-ponow';

export function pasPrzygotowania(): HTMLElement {
  return pasDzialan({
    czynnosci: [
      {
        klucz: 'dzialania.sprobujPonownie',
        czynnosc: 'ponow-przygotowanie',
        id: PONOWIENIE,
      },
      { klucz: 'dzialania.przerwijIWyloguj', czynnosc: 'przerwij-i-wyloguj' },
    ],
  });
}

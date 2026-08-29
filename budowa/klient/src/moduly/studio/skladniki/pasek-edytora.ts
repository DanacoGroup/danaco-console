/**
 * Pasek narzędziowy edytora — strefa `dn-edytor-pasek` ze źródła kształtu
 * (`design/05-okna/moduly/studio.html`): trzy grupy przycisków rozdzielone
 * kreską, odstęp rozpierający i przycisk znajdowania po prawej.
 *
 * Postać znaku i akapitu zmienia rdzeń: `studio.format.character.set` oraz
 * `studio.format.paragraph.set` biorą zakres znaków, więc pasek podaje im
 * zaznaczenie z pola treści. Bez zaznaczenia komenda nie ma na czym stanąć —
 * pasek mówi o tym wprost zamiast działać na całości.
 *
 * Czynności bez pokrycia w kontrakcie stoją zapowiedziane, nie ukryte: wykaz
 * komend Studia nie zna wstawienia tabeli ani bloku kodu z paska, a zasada zero
 * blokad każe niegotowość nazwać, nie chować przycisk.
 */

import { Command, StudioUnderlineStyle } from '../../../../../shared/contract.ts';
import type { Kanal } from '../../../protokol/kanal.ts';
import { wywolaj } from '../../../protokol/wywolanie.ts';
import { ikony, type NazwaZnaku } from '../ikony.ts';
import { el, tekst, zeZnacznika, type Dziecko } from '../narzedzia.ts';
import { zaloguj } from '../odmowa.ts';

export interface ZaleznosciPaska {
  /** Kanał, którym pasek woła komendy rdzenia. */
  kanal: Kanal;
  /** Dokument, na którym pasek pracuje; pusty znaczy brak otwartego dokumentu. */
  idDokumentu(): string | null;
  /** Zakres zaznaczenia w polu treści; `null` znaczy brak zaznaczenia. */
  zaznaczenie(): { poczatek: number; koniec: number } | null;
  /** Nazywa Operatorowi, dlaczego czynność nie doszła do skutku. */
  powiadom(zdanie: string): void;
}

export interface ZamontowanyPasekEdytora {
  wezel: HTMLElement;
}

function znak(rysunek: NazwaZnaku): SVGElement {
  const wezel = zeZnacznika(ikony[rysunek]);
  wezel.setAttribute('aria-hidden', 'true');
  return wezel;
}

/** Przycisk znakowy paska: sam znak, nazwę czynności niesie etykieta dostępności. */
function przyciskZnaku(etykieta: string, tresc: Dziecko, przyKliknieciu: () => void): HTMLElement {
  const guzik = el('button', {
    klasa: 'dn-btn-ikona dn-btn-ikona--sm',
    type: 'button',
    'aria-label': etykieta,
    title: etykieta,
  }, [tresc]);
  guzik.addEventListener('click', przyKliknieciu);
  return guzik;
}

/** Przycisk z napisem: rozwinięcia postaci akapitu, gdzie znak nie wystarcza. */
function przyciskNapisu(napis: string, przyKliknieciu: () => void): HTMLElement {
  const guzik = el('button', {
    klasa: 'dn-btn dn-btn--duch dn-btn--sm',
    type: 'button',
    tekst: `${napis} ▾`,
  });
  guzik.addEventListener('click', przyKliknieciu);
  return guzik;
}

export function pasekEdytora(zaleznosci: ZaleznosciPaska): ZamontowanyPasekEdytora {
  /*
  postacZnaku wysyła jedną własność postaci na zaznaczony fragment.

  Zakres idzie z pola treści, nie z domysłu: komenda przyjmuje `rangeStart`
  i `rangeEnd` jako opcjonalne, a pominięcie ich objęłoby cały dokument —
  Operator, który zaznaczył zdanie, dostałby pogrubienie całości.
  */
  async function postacZnaku(wlasnosc: Record<string, unknown>): Promise<void> {
    const dokument = zaleznosci.idDokumentu();
    if (dokument === null) return;
    const zakres = zaleznosci.zaznaczenie();
    if (zakres === null) {
      zaleznosci.powiadom(tekst('pasekEdytora.bezZaznaczenia'));
      return;
    }
    const wynik = await wywolaj(zaleznosci.kanal, Command.StudioFormatCharacterSet, {
      documentId: dokument,
      rangeStart: zakres.poczatek,
      rangeEnd: zakres.koniec,
      ...wlasnosc,
    });
    if (!wynik.udany) {
      zaloguj(wynik.blad, 'pasekEdytora');
      zaleznosci.powiadom(tekst('pasekEdytora.zapowiedziane'));
    }
  }

  /** Czynność, której kontrakt nie niesie: nazywa niegotowość zamiast milczeć. */
  function zapowiedziana(): void {
    zaleznosci.powiadom(tekst('pasekEdytora.zapowiedziane'));
  }

  const grupaPostaci = el('span', { klasa: 'dn-edytor-grupa' }, [
    przyciskZnaku(tekst('pasekEdytora.pogrubienie'), el('b', { tekst: 'B' }),
      () => void postacZnaku({ bold: true })),
    przyciskZnaku(tekst('pasekEdytora.kursywa'), el('i', { tekst: 'I' }),
      () => void postacZnaku({ italic: true })),
    przyciskZnaku(tekst('pasekEdytora.podkreslenie'), el('u', { tekst: 'U' }),
      () => void postacZnaku({ underline: StudioUnderlineStyle.Single })),
  ]);

  const grupaAkapitu = el('span', { klasa: 'dn-edytor-grupa' }, [
    przyciskNapisu(tekst('pasekEdytora.naglowek'), zapowiedziana),
    przyciskNapisu(tekst('pasekEdytora.stylAkapitu'), zapowiedziana),
  ]);

  const grupaBlokow = el('span', { klasa: 'dn-edytor-grupa' }, [
    przyciskZnaku(tekst('pasekEdytora.lista'), znak('lista'), zapowiedziana),
    przyciskZnaku(tekst('pasekEdytora.tabela'), znak('tabela'), zapowiedziana),
    przyciskZnaku(tekst('pasekEdytora.kod'), znak('kod'), zapowiedziana),
    przyciskZnaku(tekst('pasekEdytora.cytat'), znak('cytat'), zapowiedziana),
  ]);

  const szukanie = el('button', {
    klasa: 'dn-btn dn-btn--duch dn-btn--sm',
    type: 'button',
    'aria-label': tekst('pasekEdytora.szukaj'),
    title: tekst('pasekEdytora.szukaj'),
  }, [znak('szukaj')]);
  szukanie.addEventListener('click', zapowiedziana);

  return {
    wezel: el('div', {
      klasa: 'dn-edytor-pasek',
      role: 'toolbar',
      'aria-label': tekst('pasekEdytora.etykieta'),
    }, [
      grupaPostaci,
      grupaAkapitu,
      grupaBlokow,
      el('span', { klasa: 'st-edytor-odstep', 'aria-hidden': 'true' }),
      szukanie,
    ]),
  };
}

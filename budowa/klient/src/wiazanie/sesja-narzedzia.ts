// Narzędzia stojące przy sesji i zamknięcie sesji: dwie pozycje menu wiersza
// sesji w Centrum dowodzenia.
import { Command } from '../../../shared/contract.ts';
import type { Kanal, Wynik } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import { zapytajWOknieModalnym } from './pytanie-modalne.ts';

const NAGLOWEK = 'Sesja';

const IKONA_NARZEDZI = 'M14 6 8 12l6 6M4 6h4M4 18h4';
const IKONA_ZAMKNIECIA = 'M9 3h6M6 6h12l-1 14H7L6 6Z';

/*
dolozPozycjeSesji dostawia do menu wiersza sesji dwie pozycje, których prototyp
nie postawił: wgląd w narzędzia sesji i jej zamknięcie. Pozycje powstają z tego
samego kształtu, co pozostałe — przycisk z ikoną i napisem — więc menu zostaje
jednorodne, a mechanika otwierania i klawiatury działa bez zmian.
*/
export function dolozPozycjeSesji(wiersz: HTMLElement): void {
  const menu = wiersz.querySelector('[data-menu-tresc]');
  if (menu === null) return;
  if (menu.querySelector('[data-poz-akcja="narzedzia"]') !== null) return;
  menu.append(
    pozycja(wiersz.ownerDocument, 'narzedzia', 'Narzędzia sesji', IKONA_NARZEDZI),
    pozycja(wiersz.ownerDocument, 'zamknij', 'Zamknij sesję', IKONA_ZAMKNIECIA),
  );
}

function pozycja(
  dokument: Document,
  czynnosc: string,
  napis: string,
  sciezka: string,
): HTMLButtonElement {
  const przycisk = dokument.createElement('button');
  przycisk.type = 'button';
  przycisk.className = 'sta-menu-poz';
  przycisk.setAttribute('role', 'menuitem');
  przycisk.dataset.pozAkcja = czynnosc;
  const ikona = dokument.createElementNS('http://www.w3.org/2000/svg', 'svg');
  ikona.setAttribute('class', 'dn-menu-poz-ikona');
  ikona.setAttribute('viewBox', '0 0 24 24');
  ikona.setAttribute('fill', 'none');
  ikona.setAttribute('stroke', 'currentColor');
  ikona.setAttribute('stroke-width', '1.75');
  ikona.setAttribute('stroke-linecap', 'round');
  ikona.setAttribute('stroke-linejoin', 'round');
  ikona.setAttribute('aria-hidden', 'true');
  const rysunek = dokument.createElementNS('http://www.w3.org/2000/svg', 'path');
  rysunek.setAttribute('d', sciezka);
  ikona.appendChild(rysunek);
  przycisk.append(ikona, dokument.createTextNode(napis));
  return przycisk;
}

/* Jedno okno modalne prowadzi obie strony wykazu: to, co przy sesji stoi,
   nazywa opis, a dwa pola rozstrzygają, co dołożyć i co odjąć. */
export async function prowadzNarzedziaSesji(
  kanal: Kanal,
  idSesji: string,
): Promise<Wynik<unknown> | null> {
  const wykaz = await wywolaj(kanal, Command.SessionToolList, { sessionId: idSesji });
  if (!wykaz.udany || wykaz.wynik === undefined) {
    oglos(NAGLOWEK, wykaz.blad?.message ?? 'Wykaz narzędzi sesji nie doszedł.', 'ostrzezenie');
    return wykaz;
  }
  const stojace = wykaz.wynik.tools;
  const katalog = await wywolaj(kanal, Command.ToolsCatalogList, { sessionId: idSesji });
  const doWyboru = (katalog.wynik?.entries ?? [])
    .filter((pozycjaKatalogu) => pozycjaKatalogu.attachable)
    .map((pozycjaKatalogu) =>
      [pozycjaKatalogu.name, pozycjaKatalogu.shortName] as const);
  const wartosci = await zapytajWOknieModalnym({
    tytul: 'Narzędzia sesji',
    opis: stojace.length === 0
      ? 'Przy sesji nie stoi jeszcze żadne narzędzie.'
      : `Przy sesji stoją: ${stojace.map((narzedzie) => narzedzie.shortName).join(', ')}.`,
    pola: [
      {
        klucz: 'dolacz',
        etykieta: 'Dołóż narzędzie',
        wybor: [['', 'żadnego'] as const, ...doWyboru],
      },
      {
        klucz: 'odlacz',
        etykieta: 'Odejmij narzędzie',
        wybor: [
          ['', 'żadnego'] as const,
          ...stojace.map((narzedzie) => [narzedzie.name, narzedzie.shortName] as const),
        ],
      },
    ],
    wykonanie: 'Zapisz wybór',
  });
  if (wartosci === null) return null;
  return zastosujWybor(kanal, idSesji, wartosci.dolacz ?? '', wartosci.odlacz ?? '');
}

async function zastosujWybor(
  kanal: Kanal,
  idSesji: string,
  dolaczane: string,
  odejmowane: string,
): Promise<Wynik<unknown> | null> {
  let ostatni: Wynik<unknown> | null = null;
  if (odejmowane !== '') {
    ostatni = await wywolaj(kanal, Command.SessionToolDetach, {
      sessionId: idSesji,
      toolName: odejmowane,
    });
    if (!ostatni.udany) {
      oglos(NAGLOWEK, ostatni.blad?.message ?? 'Rdzeń odmówił odjęcia narzędzia.', 'ostrzezenie');
      return ostatni;
    }
    oglos(NAGLOWEK, `Narzędzie „${odejmowane}" nie stoi już przy sesji.`);
  }
  if (dolaczane !== '') {
    ostatni = await wywolaj(kanal, Command.SessionToolAttach, {
      sessionId: idSesji,
      toolName: dolaczane,
    });
    if (!ostatni.udany) {
      oglos(NAGLOWEK, ostatni.blad?.message ?? 'Rdzeń odmówił dołożenia narzędzia.', 'ostrzezenie');
      return ostatni;
    }
    oglos(NAGLOWEK, `Narzędzie „${dolaczane}" stoi teraz przy sesji.`);
  }
  /* Wybór bez żadnej strony nie jest błędem, ale też nic nie robi — okno mówi
     to wprost, żeby zamknięcie modalu nie wyglądało na wykonaną czynność. */
  if (ostatni === null) oglos(NAGLOWEK, 'Wykaz narzędzi sesji został bez zmian.');
  return ostatni;
}

/* Zamknięcie zwija okna sesji, ale sesji nie kasuje: rozmowa zostaje w wykazie
   i daje się podjąć na nowo. */
export async function zamknijSesje(
  kanal: Kanal,
  idSesji: string,
): Promise<Wynik<unknown> | null> {
  const wartosci = await zapytajWOknieModalnym({
    tytul: 'Zamknięcie sesji',
    pola: [],
    wykonanie: 'Zamknij sesję',
    nieodwracalne: 'Okna sesji zostaną zwinięte; biegnące tury przestaną pracować.',
  });
  if (wartosci === null) return null;
  const wynik = await wywolaj(kanal, Command.SessionClose, { sessionId: idSesji });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zamknięcia sesji.', 'ostrzezenie');
    return wynik;
  }
  oglos(NAGLOWEK, 'Sesja zamknięta; rozmowa zostaje w wykazie.');
  return wynik;
}

// Wiązanie karty modułu Translate z rdzeniem: okno tłumaczenia i wykaz glosariusza. Węzły pochodzą ze znacznika Właściciela, treść z rdzenia.

import {
  Command,
  type GlossaryTerm,
} from '../../../shared/contract.ts';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import {
  opiszNaglowek,
  zapewnijOknoModulu,
  zdejmijSterowanieWspolne,
  zdejmijTrescWspolna,
} from './okno-modulu.ts';
import { zwiazKatalogModulu } from './katalog-modulu.ts';

const KOD_MODULU = 'translate';

interface WiazanieKarty {
  korzen: Element;
  odlaczenia: Odsubskrybuj[];
}

const WIAZANIA = new Map<string, WiazanieKarty>();

export function zwiazTlumaczenie(
  kanal: Kanal,
  nazwaSrodowiska: string,
  idOknaStojacego: string,
  wskazanieKorzenia: Element | string,
): boolean {
  const korzen = korzenKarty(wskazanieKorzenia);
  if (korzen === null) return false;
  const idKarty = korzen.getAttribute('data-karta') ?? '';
  if (idKarty === '') return false;
  if (WIAZANIA.get(idKarty)?.korzen === korzen) return false;
  zwolnijTlumaczenie(idKarty);

  const odlaczenia: Odsubskrybuj[] = [];
  const sterowanie = new AbortController();
  const przy = { signal: sterowanie.signal };
  odlaczenia.push(() => {
    sterowanie.abort();
  });
  WIAZANIA.set(idKarty, { korzen, odlaczenia });

  if (korzen instanceof HTMLElement) {
    zdejmijTrescWspolna(korzen);
    zdejmijSterowanieWspolne(korzen);
  }
  zdejmijTrescPrzykladowa(korzen);
  void opiszNaglowek(kanal, nazwaSrodowiska, korzen);

  let idOkna = idOknaStojacego;
  const odswiez = async (): Promise<void> => {
    if (idOkna === '') return;
    await wypelnijGlosariusz(kanal, korzen);
  };

  void (async (): Promise<void> => {
    idOkna = await zapewnijOknoModulu(kanal, idKarty, KOD_MODULU, 'Translate', idOkna);
    if (idOkna === '') {
      niegotowe(panel(korzen, 'panel-glossary'), 'Rdzeń nie dał okna tłumaczenia dla tej karty.');
      return;
    }
    await odswiez();
  })();

  korzen.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    if (cel.closest('#panel-glossary .sta-okno-akcje .dn-btn-ikona') === null) return;
    zdarzenie.stopPropagation();
    void wniesPanel(kanal, korzen, idOkna).then(odswiez);
  }, przy);

  const katalog = zwiazKatalogModulu(kanal, idOkna, korzen, KOD_MODULU, 'Tłumaczenie');
  if (katalog !== null) odlaczenia.push(katalog);
  return true;
}

export function zwolnijTlumaczenie(idKarty: string): void {
  const wiazanie = WIAZANIA.get(idKarty);
  if (wiazanie === undefined) return;
  for (const odlacz of wiazanie.odlaczenia) odlacz();
  WIAZANIA.delete(idKarty);
}

function korzenKarty(wskazanie: Element | string): Element | null {
  if (typeof wskazanie !== 'string') return wskazanie;
  return document.querySelector(`.cd-tresc--modul[data-karta="${wskazanie}"]`);
}

function panel(korzen: Element, identyfikator: string): Element | null {
  return korzen.querySelector(`#${identyfikator} .sta-okno-tresc`);
}

function niegotowe(cialo: Element | null, zdanie: string): void {
  if (cialo === null) return;
  const napis = cialo.ownerDocument.createElement('div');
  napis.className = 'dn-meta';
  napis.textContent = zdanie;
  cialo.replaceChildren(napis);
}

/* Znacznik niesie źródła, ustalenia i miary wpisane wprost. Schodzą przed
   pierwszym pytaniem rdzenia, żeby okno nie pokazywało cudzego badania. */
function zdejmijTrescPrzykladowa(korzen: Element): void {
  niegotowe(panel(korzen, 'panel-glossary'), 'Wykaz czeka na odpowiedź rdzenia.');
  /* Znacznik niesie tekst źródłowy, gotowe tłumaczenia i miary segmentów
     wpisane wprost. Rdzeń ich nie odtworzy, więc schodzą razem z panelami. */
  for (const segment of korzen.querySelectorAll('.tr-seg, .tr-kol-tresc > *')) segment.remove();
  for (const miara of korzen.querySelectorAll('.tr-glowa span:not(.tr-kod)')) miara.remove();
  // Panele języków i wykaz plików prototyp wpisuje wprost; rdzeń poda swoje.
  for (const kolumna of korzen.querySelectorAll('.tr-kol')) kolumna.remove();
  niegotowe(panel(korzen, 'panel-pliki'), 'Rdzeń nie podaje plików tłumaczenia dla tego okna.');
  niegotowe(panel(korzen, 'panel-obszar-tlum'), 'Panele tłumaczenia czekają na wskazanie języka.');
  niegotowe(panel(korzen, 'panel-plan'), 'Plan tłumaczenia czeka na pierwszy panel języka.');
  niegotowe(panel(korzen, 'panel-kolejka'), 'Rdzeń nie podaje zadań w tle dla tego okna.');
  niegotowe(panel(korzen, 'panel-zadania'), 'Rdzeń nie podaje zadań tłumaczenia dla tego okna.');
}

async function wypelnijGlosariusz(kanal: Kanal, korzen: Element): Promise<void> {
  const cialo = panel(korzen, 'panel-glossary');
  if (cialo === null) return;
  const odpowiedz = await wywolaj(kanal, Command.TranslateGlossaryList, {});
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    niegotowe(cialo, odpowiedz.blad?.message ?? 'Wykaz glosariusza nie doszedł.');
    return;
  }
  const wpisy: GlossaryTerm[] = odpowiedz.wynik.terms;
  if (wpisy.length === 0) {
    niegotowe(cialo, 'Glosariusz jest pusty — żadnego terminu jeszcze nie wniesiono.');
    return;
  }
  cialo.replaceChildren(...wpisy.map((w) => pozycja(cialo, w.source, w.target ?? w.language)));
}

function pozycja(cialo: Element, tytul: string, podpis: string): HTMLElement {
  const wiersz = cialo.ownerDocument.createElement('div');
  wiersz.className = 'dn-wykaz-modulu-poz';
  const nazwa = cialo.ownerDocument.createElement('span');
  nazwa.textContent = tytul;
  wiersz.append(nazwa);
  if (podpis !== '') {
    const meta = cialo.ownerDocument.createElement('span');
    meta.className = 'dn-meta';
    meta.textContent = podpis;
    wiersz.append(meta);
  }
  return wiersz;
}

/* Okna zakładania źródła wydanie nie niesie: tytuł wchodzi zapytaniem w oknie,
   tak samo jak nazwa komponentu w Centrum. */
async function wniesPanel(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  if (idOkna === '') {
    oglos('Tłumaczenie', 'Okno tłumaczenia jeszcze nie powstało, więc nie ma do czego dołożyć panelu.');
    return;
  }
  const jezyk = globalThis.prompt?.('Język panelu (kod, np. DE)') ?? '';
  if (jezyk.trim() === '') return;
  const odpowiedz = await wywolaj(kanal, Command.TranslateTargetAdd, {
    windowId: idOkna,
    language: jezyk.trim(),
  });
  if (!odpowiedz.udany) {
    oglos('Tłumaczenie', odpowiedz.blad?.message ?? 'Rdzeń odmówił dołożenia panelu.', 'blad');
    return;
  }
  oglos('Tłumaczenie', `Panel języka ${jezyk.trim()} dołożony.`);
  void korzen;
}

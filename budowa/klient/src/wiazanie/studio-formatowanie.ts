/**
 * Pasek formatowania kanwy Studia: postać znaku zaznaczonego fragmentu oraz
 * wyrównanie jego akapitów. Komendy postaci biorą brak zakresu jako cały
 * dokument, więc czynność bez zaznaczenia wraca odmową zamiast objąć całość.
 */

import {
  Command,
  EventType,
  StudioTextAlign,
  StudioUnderlineStyle,
  type StudioCharacterFormat,
  type StudioDocument,
} from '../../../shared/contract.ts';
import { fragmentZaznaczony, sledzZaznaczenie } from '../model/zaznaczenie-dokumentu.ts';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen.ts';
import { zglosUchwyt } from '../polaczenie/rozdzielacz-zdarzen.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import { nazwijBilans, nazwijOdmowe, wskazDokumentOkna } from './studio-dokument.ts';

const NAZWY_WYROWNANIA: Readonly<Record<StudioTextAlign, string>> = {
  [StudioTextAlign.Left]: 'do lewej',
  [StudioTextAlign.Center]: 'wyśrodkowanie',
  [StudioTextAlign.Right]: 'do prawej',
  [StudioTextAlign.Justify]: 'wyjustowanie',
};

const BEZ_ZAZNACZENIA = 'Wskaż fragment — postać dotyczy zaznaczenia.';

interface WezlyPaska {
  kanwa: HTMLElement;
  pogrubienie: HTMLElement;
  kursywa: HTMLElement;
  podkreslenie: HTMLElement;
  wyrownanie: HTMLElement;
  szukaj: HTMLElement | null;
}

/** Nadaje przyciskom paska stan spoczynku: postać z prototypu opisuje cudzy dokument. */
export function zdejmijTrescPrzykladowaFormatowania(korzen: ParentNode): boolean {
  const wezly = zbierzWezly(korzen);
  if (wezly === null) return false;
  for (const przycisk of [wezly.pogrubienie, wezly.kursywa, wezly.podkreslenie]) {
    przycisk.setAttribute('aria-pressed', 'false');
  }
  wezly.wyrownanie.textContent = 'Wyrównanie ▾';
  nazwijSzukanie(wezly.szukaj);
  return true;
}

// Zamiana treści wymaga dwóch pól, których pasek nie ma; przycisk prowadzi
// do panelu różnic, gdzie stoi pole wzorca.
function nazwijSzukanie(szukaj: HTMLElement | null): void {
  if (szukaj === null) return;
  szukaj.setAttribute('aria-label', 'Znajdź');
  const napis = [...szukaj.childNodes].find((wezel) => wezel.nodeType === Node.TEXT_NODE);
  if (napis !== undefined) napis.textContent = ' Znajdź';
}

/** Wiąże pasek formatowania z dokumentem okna; zwraca odłączenie, a pustkę, gdy paska w karcie nie ma. */
export function zwiazFormatowanie(
  kanal: Kanal,
  idOkna: string,
  korzen: ParentNode,
): Odsubskrybuj | null {
  const znalezione = zbierzWezly(korzen);
  if (znalezione === null) return null;
  const wezly: WezlyPaska = znalezione;
  const sterowanie = new AbortController();
  const przy = { signal: sterowanie.signal };
  sledzZaznaczenie(wezly.kanwa, przy);
  let dokument: StudioDocument | null = null;
  let wyrownanie: StudioTextAlign | null = null;

  const opiszWyrownanie = (): void => {
    const nazwa = wyrownanie === null ? '' : `: ${NAZWY_WYROWNANIA[wyrownanie]}`;
    wezly.wyrownanie.textContent = `Wyrównanie${nazwa} ▾`;
  };

  const odczytajPostac = async (): Promise<void> => {
    const zakres = fragmentZaznaczony(wezly.kanwa);
    if (dokument === null || zakres === null) return;
    const znak = await wywolaj(kanal, Command.StudioFormatCharacterGet, {
      documentId: dokument.id,
      rangeStart: zakres.poczatek,
      rangeEnd: zakres.koniec,
    });
    if (znak.udany && znak.wynik !== undefined) oznaczPostac(wezly, znak.wynik.character);
    const akapit = await wywolaj(kanal, Command.StudioFormatParagraphGet, {
      documentId: dokument.id,
      rangeStart: zakres.poczatek,
      rangeEnd: zakres.koniec,
    });
    if (!akapit.udany || akapit.wynik === undefined) return;
    wyrownanie = akapit.wynik.paragraph.align ?? null;
    opiszWyrownanie();
  };

  const ustawZnak = async (
    przycisk: HTMLElement,
    ceche: (wciety: boolean) => Partial<StudioCharacterFormat>,
  ): Promise<void> => {
    const zakres = fragmentZaznaczony(wezly.kanwa);
    if (dokument === null || zakres === null) {
      oglos('Studio', BEZ_ZAZNACZENIA);
      return;
    }
    const wciety = przycisk.getAttribute('aria-pressed') !== 'true';
    const wynik = await wywolaj(kanal, Command.StudioFormatCharacterSet, {
      documentId: dokument.id,
      rangeStart: zakres.poczatek,
      rangeEnd: zakres.koniec,
      ...ceche(wciety),
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      oglos('Studio', nazwijOdmowe('Postać znaku', wynik.blad), 'ostrzezenie');
      return;
    }
    przycisk.setAttribute('aria-pressed', wciety ? 'true' : 'false');
    if (wynik.wynik.balance.skippedCount > 0) {
      oglos('Studio', nazwijBilans('Postać znaku', wynik.wynik.balance), 'ostrzezenie');
    }
  };

  const przestawWyrownanie = async (): Promise<void> => {
    const zakres = fragmentZaznaczony(wezly.kanwa);
    if (dokument === null || zakres === null) {
      oglos('Studio', BEZ_ZAZNACZENIA);
      return;
    }
    const nastepne = nastepneWyrownanie(wyrownanie);
    const wynik = await wywolaj(kanal, Command.StudioFormatParagraphSet, {
      documentId: dokument.id,
      rangeStart: zakres.poczatek,
      rangeEnd: zakres.koniec,
      align: nastepne,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      oglos('Studio', nazwijOdmowe('Wyrównanie akapitu', wynik.blad), 'ostrzezenie');
      return;
    }
    wyrownanie = nastepne;
    opiszWyrownanie();
    if (wynik.wynik.balance.skippedCount > 0) {
      oglos('Studio', nazwijBilans('Wyrównanie akapitu', wynik.wynik.balance), 'ostrzezenie');
    }
  };

  wezly.pogrubienie.addEventListener('click', () => {
    void ustawZnak(wezly.pogrubienie, (wciety) => ({ bold: wciety }));
  }, przy);
  wezly.kursywa.addEventListener('click', () => {
    void ustawZnak(wezly.kursywa, (wciety) => ({ italic: wciety }));
  }, przy);
  wezly.podkreslenie.addEventListener('click', () => {
    void ustawZnak(wezly.podkreslenie, (wciety) => ({
      underline: wciety ? StudioUnderlineStyle.Single : StudioUnderlineStyle.None,
    }));
  }, przy);
  wezly.wyrownanie.addEventListener('click', () => {
    void przestawWyrownanie();
  }, przy);
  wezly.kanwa.addEventListener('mouseup', () => {
    void odczytajPostac();
  }, przy);
  wezly.szukaj?.addEventListener('click', () => {
    otworzWzorzec(korzen);
  }, przy);

  const odlacz = zglosUchwyt(EventType.StudioDocumentChanged, (tresc) => {
    if (tresc.document.windowId !== idOkna) return;
    dokument = tresc.document;
  });

  void wskazDokumentOkna(kanal, idOkna).then((otwarty) => {
    dokument = otwarty;
  });

  return () => {
    sterowanie.abort();
    odlacz();
  };
}

function otworzWzorzec(korzen: ParentNode): void {
  const karta = korzen.querySelector<HTMLElement>('.st-karty [role="tab"][data-karta="diff"]');
  karta?.click();
  korzen.querySelector<HTMLInputElement>('#panel-diff input[type="search"]')?.focus();
}

function oznaczPostac(wezly: WezlyPaska, postac: StudioCharacterFormat): void {
  wezly.pogrubienie.setAttribute('aria-pressed', postac.bold === true ? 'true' : 'false');
  wezly.kursywa.setAttribute('aria-pressed', postac.italic === true ? 'true' : 'false');
  const podkreslony = postac.underline !== undefined && postac.underline !== StudioUnderlineStyle.None;
  wezly.podkreslenie.setAttribute('aria-pressed', podkreslony ? 'true' : 'false');
}

function nastepneWyrownanie(biezace: StudioTextAlign | null): StudioTextAlign {
  const wykaz = Object.values(StudioTextAlign);
  const miejsce = biezace === null ? -1 : wykaz.indexOf(biezace);
  return wykaz[(miejsce + 1) % wykaz.length] ?? StudioTextAlign.Left;
}

function zbierzWezly(korzen: ParentNode): WezlyPaska | null {
  const pasek = korzen.querySelector('.dn-edytor-pasek');
  const kanwa = korzen.querySelector('.dn-kanwa');
  const grupy = [...(pasek?.querySelectorAll('.dn-edytor-grupa') ?? [])];
  const wyrownanie = [...(grupy[1]?.querySelectorAll('button') ?? [])][1];
  const szukaj = pasek?.querySelector('[aria-label="Znajdź i zamień"], [aria-label="Znajdź"]');
  const pogrubienie = pasek?.querySelector('[aria-label="Pogrubienie"]');
  const kursywa = pasek?.querySelector('[aria-label="Kursywa"]');
  const podkreslenie = pasek?.querySelector('[aria-label="Podkreślenie"]');
  if (
    !(kanwa instanceof HTMLElement) ||
    !(pogrubienie instanceof HTMLElement) ||
    !(kursywa instanceof HTMLElement) ||
    !(podkreslenie instanceof HTMLElement) ||
    !(wyrownanie instanceof HTMLElement)
  ) {
    return null;
  }
  return {
    kanwa,
    pogrubienie,
    kursywa,
    podkreslenie,
    wyrownanie,
    szukaj: szukaj instanceof HTMLElement ? szukaj : null,
  };
}

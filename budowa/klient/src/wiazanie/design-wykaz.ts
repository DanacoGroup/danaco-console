// Dwa wspólne cele wypełniania: panel „Plan" przyjmuje wykaz roboczy wskazany
// paskiem kompozycji, panel „Pliki" przyjmuje wytwory wyniesione z rdzenia.

import { nieGotowe, panel, tresc, wypelnijWykaz, rozmiar, type Kontekst, type Wiersz } from './design-wspolne.ts';

export interface PlikWytworzony {
  fileName: string;
  mediaType: string;
  contentBase64: string;
  sizeBytes?: number;
}

const PLIKI = new Map<string, PlikWytworzony[]>();

export interface PozycjaKorzenia {
  nazwa: string;
  opis: string;
  otworz: (kontekst: Kontekst) => void;
}

// Rejestr zamiast wzajemnych importów między podgrupami wiązania.
const KORZEN: PozycjaKorzenia[] = [];

export function dopiszDoKorzenia(pozycja: PozycjaKorzenia): void {
  if (KORZEN.some((obecna) => obecna.nazwa === pozycja.nazwa)) return;
  KORZEN.push(pozycja);
}

export function pozycjeKorzenia(): readonly PozycjaKorzenia[] {
  return KORZEN;
}

export function pokazWykaz(
  kontekst: Kontekst,
  nazwa: string,
  wiersze: Wiersz[],
  zdanieBraku?: string,
): void {
  const okno = panel(kontekst.korzen, 'panel-plan');
  const cialo = tresc(okno);
  if (okno === null || cialo === null) return;
  okno.removeAttribute('hidden');
  const naglowek = okno.querySelector('.sta-okno-tytul b');
  if (naglowek !== null) naglowek.textContent = nazwa;
  wypelnijWykaz(kontekst, cialo, wiersze, zdanieBraku);
  uzgodnijPrzelacznik(kontekst.korzen, 'panel-plan');
}

export function pustyWykaz(kontekst: Kontekst, zdanie: string): void {
  const cialo = tresc(panel(kontekst.korzen, 'panel-plan'));
  nieGotowe(cialo, zdanie);
}

export function dolozPlik(kontekst: Kontekst, plik: PlikWytworzony): void {
  const zebrane = PLIKI.get(kontekst.idKarty) ?? [];
  zebrane.unshift(plik);
  PLIKI.set(kontekst.idKarty, zebrane.slice(0, 40));
  odswiezPliki(kontekst);
}

export function zapomnijPliki(idKarty: string): void {
  PLIKI.delete(idKarty);
}

export function odswiezPliki(kontekst: Kontekst): void {
  const okno = panel(kontekst.korzen, 'panel-pliki');
  const cialo = tresc(okno);
  if (okno === null || cialo === null) return;
  const zebrane = PLIKI.get(kontekst.idKarty) ?? [];
  if (zebrane.length === 0) {
    nieGotowe(cialo, 'Żaden plik nie został jeszcze wyniesiony w tym oknie.');
    return;
  }
  okno.removeAttribute('hidden');
  const wiersze = zebrane.map((plik) => {
    const opis = plik.sizeBytes === undefined
      ? plik.mediaType
      : `${plik.mediaType} · ${rozmiar(plik.sizeBytes)}`;
    return {
      tekst: plik.fileName,
      meta: opis,
      naKlik: () => {
        zapisz(kontekst.korzen.ownerDocument, plik);
      },
    } satisfies Wiersz;
  });
  wypelnijWykaz(kontekst, cialo, wiersze);
  uzgodnijPrzelacznik(kontekst.korzen, 'panel-pliki');
}

// Biblioteka zakładek czyta stan panelu z atrybutu `hidden`, więc odsłonięcie
// panelu musi przestawić także znacznik pozycji w menu paneli.
function uzgodnijPrzelacznik(korzen: Element, identyfikator: string): void {
  const przelacznik = korzen.querySelector(`[data-panel-toggle="${identyfikator}"]`);
  przelacznik?.setAttribute('aria-checked', 'true');
}

function zapisz(dokument: Document, plik: PlikWytworzony): void {
  const bajty = Uint8Array.from(atob(plik.contentBase64), (znak) => znak.charCodeAt(0));
  const adres = URL.createObjectURL(new Blob([bajty], { type: plik.mediaType }));
  const odsylacz = dokument.createElement('a');
  odsylacz.href = adres;
  odsylacz.download = plik.fileName;
  odsylacz.click();
  URL.revokeObjectURL(adres);
}

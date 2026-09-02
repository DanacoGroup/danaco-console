// Style nazwane w pasku kanwy Studia: przycisk nagłówka krąży po stylach
// arkusza dokumentu, przycisk cytatu stosuje styl o nazwie własnej etykiety.

import {
  Command,
  EventType,
  StudioStyleKind,
  type StudioDocument,
  type StudioNamedStyle,
} from '../../../shared/contract.ts';
import { fragmentZaznaczony } from '../model/zaznaczenie-dokumentu.ts';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen.ts';
import { zglosUchwyt } from '../polaczenie/rozdzielacz-zdarzen.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import { nazwijBilans, nazwijOdmowe, wskazDokumentOkna } from './studio-dokument.ts';

interface WezlyStylow {
  kanwa: HTMLElement;
  naglowek: HTMLElement;
  cytat: HTMLElement;
}

// Przycisk bloku kodu schodzi: arkusz stylów rdzenia nie ma stylu o tej
// nazwie, a innej komendy dla niego kontrakt nie niesie.
export function zdejmijTrescPrzykladowaStylow(korzen: ParentNode): boolean {
  const wezly = zbierzWezly(korzen);
  korzen.querySelector('.dn-edytor-pasek [aria-label="Blok kodu"]')?.remove();
  if (wezly === null) return false;
  wezly.naglowek.textContent = 'Nagłówek ▾';
  return true;
}

export function zwiazStyle(kanal: Kanal, idOkna: string, korzen: ParentNode): Odsubskrybuj | null {
  const znalezione = zbierzWezly(korzen);
  if (znalezione === null) return null;
  const wezly: WezlyStylow = znalezione;
  const sterowanie = new AbortController();
  const przy = { signal: sterowanie.signal };
  let dokument: StudioDocument | null = null;
  let style: StudioNamedStyle[] = [];
  let miejsce = 0;

  const zastosuj = async (styl: StudioNamedStyle | undefined): Promise<boolean> => {
    const zakres = fragmentZaznaczony(wezly.kanwa);
    if (dokument === null || styl === undefined) return false;
    if (zakres === null) {
      oglos('Studio', 'Wskaż fragment — styl nazwany dotyczy zaznaczenia.');
      return false;
    }
    const wynik = await wywolaj(kanal, Command.StudioStyleApply, {
      documentId: dokument.id,
      name: styl.name,
      rangeStart: zakres.poczatek,
      rangeEnd: zakres.koniec,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      oglos('Studio', nazwijOdmowe('Styl nazwany', wynik.blad), 'ostrzezenie');
      return false;
    }
    if (wynik.wynik.balance.skippedCount > 0) {
      oglos('Studio', nazwijBilans('Styl nazwany', wynik.wynik.balance), 'ostrzezenie');
    }
    return true;
  };

  wezly.naglowek.addEventListener('click', () => {
    if (style.length === 0) {
      oglos('Studio', 'Arkusz stylów dokumentu jeszcze nie stoi.');
      return;
    }
    const nastepne = (miejsce + 1) % style.length;
    void zastosuj(style[nastepne]).then((zastosowany) => {
      if (!zastosowany) return;
      miejsce = nastepne;
      wezly.naglowek.textContent = `${nazwaStylu(style[miejsce])} ▾`;
    });
  }, przy);

  wezly.cytat.addEventListener('click', () => {
    const szukana = wezly.cytat.getAttribute('aria-label') ?? '';
    const styl = style.find((pozycja) => nazwaStylu(pozycja) === szukana);
    if (styl === undefined) {
      oglos('Studio', `Arkusz stylów dokumentu nie ma stylu ${szukana}.`, 'ostrzezenie');
      return;
    }
    void zastosuj(styl);
  }, przy);

  const odlacz = zglosUchwyt(EventType.StudioDocumentChanged, (tresc) => {
    if (tresc.document.windowId !== idOkna) return;
    dokument = tresc.document;
  });

  void wskazDokumentOkna(kanal, idOkna).then(async (otwarty) => {
    dokument = otwarty;
    if (otwarty === null) return;
    const wynik = await wywolaj(kanal, Command.StudioStyleList, {
      documentId: otwarty.id,
      kind: StudioStyleKind.Paragraph,
    });
    if (!wynik.udany || wynik.wynik === undefined) return;
    style = uporzadkujStyle(wynik.wynik.styles);
    if (style.length > 0) wezly.naglowek.textContent = `${nazwaStylu(style[0])} ▾`;
  });

  return () => {
    sterowanie.abort();
    odlacz();
  };
}

// Styl podstawowy stoi pierwszy, dalej nagłówki wedle poziomu konspektu.
function uporzadkujStyle(style: StudioNamedStyle[]): StudioNamedStyle[] {
  const poziom = (styl: StudioNamedStyle): number => styl.paragraph?.outlineLevel ?? 0;
  const naglowki = style
    .filter((styl) => poziom(styl) > 0)
    .sort((pierwszy, drugi) => poziom(pierwszy) - poziom(drugi));
  const nazwaPodstawy = naglowki[0]?.basedOn;
  const podstawa = style.find((styl) => styl.name === nazwaPodstawy)
    ?? style.find((styl) => poziom(styl) === 0);
  return podstawa === undefined ? naglowki : [podstawa, ...naglowki];
}

function nazwaStylu(styl: StudioNamedStyle | undefined): string {
  return styl?.displayName ?? styl?.name ?? '';
}

function zbierzWezly(korzen: ParentNode): WezlyStylow | null {
  const pasek = korzen.querySelector('.dn-edytor-pasek');
  const kanwa = korzen.querySelector('.dn-kanwa');
  const grupy = [...(pasek?.querySelectorAll('.dn-edytor-grupa') ?? [])];
  const naglowek = [...(grupy[1]?.querySelectorAll('button') ?? [])][0];
  const cytat = pasek?.querySelector('[aria-label="Cytat"]');
  if (
    !(kanwa instanceof HTMLElement) ||
    !(naglowek instanceof HTMLElement) ||
    !(cytat instanceof HTMLElement)
  ) {
    return null;
  }
  return { kanwa, naglowek, cytat };
}

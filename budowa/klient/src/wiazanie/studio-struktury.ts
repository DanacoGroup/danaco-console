// Struktury wstawiane z paska kanwy Studia: wypunktowanie zaznaczenia i tabela
// w miejscu kursora. Rozmiar tabeli stoi tutaj, bo pasek nie ma go gdzie przyjąć.

import {
  Command,
  EventType,
  StudioListKind,
  type StudioDocument,
} from '../../../shared/contract.ts';
import { fragmentZaznaczony, miejsceKursora } from '../model/zaznaczenie-dokumentu.ts';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen.ts';
import { zglosUchwyt } from '../polaczenie/rozdzielacz-zdarzen.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import { nazwijBilans, nazwijOdmowe, wskazDokumentOkna } from './studio-dokument.ts';

const WIERSZE_TABELI = 3;
const KOLUMNY_TABELI = 3;

interface WezlyStruktur {
  kanwa: HTMLElement;
  lista: HTMLElement;
  tabela: HTMLElement;
}

export function zdejmijTrescPrzykladowaStruktur(korzen: ParentNode): boolean {
  const wezly = zbierzWezly(korzen);
  if (wezly === null) return false;
  wezly.lista.setAttribute('aria-pressed', 'false');
  return true;
}

export function zwiazStruktury(
  kanal: Kanal,
  idOkna: string,
  korzen: ParentNode,
): Odsubskrybuj | null {
  const znalezione = zbierzWezly(korzen);
  if (znalezione === null) return null;
  const wezly: WezlyStruktur = znalezione;
  const sterowanie = new AbortController();
  const przy = { signal: sterowanie.signal };
  let dokument: StudioDocument | null = null;

  const przestawListe = async (): Promise<void> => {
    const zakres = fragmentZaznaczony(wezly.kanwa);
    if (dokument === null) return;
    if (zakres === null) {
      oglos('Studio', 'Wskaż fragment — wypunktowanie obejmuje zaznaczone akapity.');
      return;
    }
    const wciety = wezly.lista.getAttribute('aria-pressed') !== 'true';
    const wynik = await wywolaj(kanal, Command.StudioListApply, {
      documentId: dokument.id,
      rangeStart: zakres.poczatek,
      rangeEnd: zakres.koniec,
      kind: wciety ? StudioListKind.Bullet : StudioListKind.None,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      oglos('Studio', nazwijOdmowe('Wypunktowanie', wynik.blad), 'ostrzezenie');
      return;
    }
    wezly.lista.setAttribute('aria-pressed', wciety ? 'true' : 'false');
    if (wynik.wynik.balance.skippedCount > 0) {
      oglos('Studio', nazwijBilans('Wypunktowanie', wynik.wynik.balance), 'ostrzezenie');
    }
  };

  const wstawTabele = async (): Promise<void> => {
    if (dokument === null) return;
    const wynik = await wywolaj(kanal, Command.StudioTableInsert, {
      documentId: dokument.id,
      offset: miejsceKursora(wezly.kanwa),
      rows: WIERSZE_TABELI,
      columns: KOLUMNY_TABELI,
      headerRows: 1,
      repeatHeader: true,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      oglos('Studio', nazwijOdmowe('Wstawienie tabeli', wynik.blad), 'ostrzezenie');
      return;
    }
    oglos('Studio', nazwijBilans('Wstawienie tabeli', wynik.wynik.balance));
  };

  wezly.lista.addEventListener('click', () => {
    void przestawListe();
  }, przy);
  wezly.tabela.addEventListener('click', () => {
    void wstawTabele();
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

function zbierzWezly(korzen: ParentNode): WezlyStruktur | null {
  const pasek = korzen.querySelector('.dn-edytor-pasek');
  const kanwa = korzen.querySelector('.dn-kanwa');
  const lista = pasek?.querySelector('[aria-label="Lista punktowana"]');
  const tabela = pasek?.querySelector('[aria-label="Tabela"]');
  if (
    !(kanwa instanceof HTMLElement) ||
    !(lista instanceof HTMLElement) ||
    !(tabela instanceof HTMLElement)
  ) {
    return null;
  }
  return { kanwa, lista, tabela };
}

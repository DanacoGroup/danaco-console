import type { ResearchReportExportResponse } from '../../../../shared/contract';
import { pokazKomunikat } from '../../aplikacja/komunikaty';
import { opisOdmowyBledu } from '../../komponenty/odmowa';
import type { WierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import type { AkcjaBadania } from './akcje-okien';
import { pokazPustke, PUSTKA_EKSPORTU } from './pustka-okien';
import type { StanBadania } from './stan-badania';
import { czyKomendaBadania, wykonajKomendeBadania } from './wywolania-komend';
import type { StanOknaBadania } from './stan-okna-badania';
import type { ZlecenieEksportu } from './zlecenia-badania';

/**
 * Czynności operatora w Export Panel: przycisk wydania zostaje klikalny
 * nawet wtedy, gdy raportu jeszcze nie ma.
 */
export interface KontekstEksportu {
  stan: StanBadania;
  okno: StanOknaBadania;
  odpowiedz: WierszOdpowiedzi;
  /** Zlecenie złożone z pól okna; identyfikator raportu dokłada czynność. */
  zlecenie(idRaportu: string): ZlecenieEksportu;
  /** Odsłania kontrolkę ponowienia po nieudanym wydaniu. */
  odslonPonowienie(): void;
  /** Dokłada wpis do historii wydań bieżącej sesji — treścią odpowiedzi rdzenia. */
  odnotujWydanie(zdanie: string): void;
  odswiez(): void;
}

/** Rozdziela akcję panelu na drogę własną okna i drogę generyczną, wspólną dla całej rodziny komend eksportu. */
export async function wykonajAkcjeEksportu(
  kontekst: KontekstEksportu,
  akcja: AkcjaBadania,
): Promise<void> {
  if (akcja.kod === 'research.export.run') {
    await wydajRaport(kontekst, 'Eksport raportu');
    return;
  }
  if (akcja.kod === 'research.export.again') {
    await wydajRaport(kontekst, 'Ponowne pobranie raportu');
    return;
  }
  if (akcja.droga === 'komenda' && czyKomendaBadania(akcja.kod)) {
    const wynik = await wykonajKomendeBadania(
      { stan: kontekst.stan, tekst: tekstDlaKomendy(kontekst) },
      akcja,
    );
    kontekst.odpowiedz.pokaz(wynik.opis, wynik.udany);
    return;
  }
  await przezPanelAkcji(kontekst, akcja);
}

/** Droga generyczna: `window.action` z identyfikatorem raportu w parametrach żądania tego okna panelu eksportu. */
async function przezPanelAkcji(kontekst: KontekstEksportu, akcja: AkcjaBadania): Promise<void> {
  const { stan, odpowiedz } = kontekst;
  if (stan.idOkna() === '') {
    odpowiedz.pokaz(
      `Akcja „${akcja.nazwa}" wymaga okna badania, którego rdzeń jeszcze nie wskazał.`,
      false,
    );
    return;
  }
  odpowiedz.pokaz(`Wysłano window.action ${akcja.kod}…`, true);
  const wynik = await stan.okna.akcja(stan.idOkna(), akcja.kod, {
    reportId: stan.raport()?.id ?? '',
  });
  if (!wynik.udany) {
    odpowiedz.pokaz(opisOdmowyBledu(`Akcja „${akcja.nazwa}"`, wynik.blad, wynik.nieznanyTyp), false);
    return;
  }
  // Rdzeń oddaje sukces window.action tylko wtedy, gdy akcję wykonał, nie samo wywołanie.
  odpowiedz.pokaz(`Rdzeń oddał wynik akcji „${akcja.nazwa}".`, true);
}

/** Znacznik trwającego wydania — idempotencja po stronie logiki, nie kontrolki panelu eksportu raportu. */
const trwajace = new WeakSet<KontekstEksportu>();

/** Wydanie raportu komendą `research.report.export` z pól zlecenia okna panelu eksportu tego badania rdzenia. */
export async function wydajRaport(kontekst: KontekstEksportu, czynnosc: string): Promise<void> {
  const { stan, okno, odpowiedz } = kontekst;
  const raport = stan.raport();
  if (raport === null) {
    odpowiedz.pokaz(
      'Nie ma czego wydać: raport nie został jeszcze złożony w Report Builderze. Przycisk zostaje czynny — to warunek merytoryczny, nie blokada.',
      false,
    );
    return;
  }
  if (trwajace.has(kontekst)) {
    odpowiedz.pokaz('Eksport już trwa — czekam na odpowiedź rdzenia na poprzednie żądanie.', true);
    return;
  }
  trwajace.add(kontekst);
  okno.ladowanie('Wydawanie raportu…');
  // Zlecenie zdejmowane z pól jeden raz, zanim operator zdąży je przestawić przed odpowiedzią.
  const zlecenie = kontekst.zlecenie(raport.id);
  const wynik = await stan.zrodlo.eksportuj(zlecenie);
  trwajace.delete(kontekst);
  kontekst.odslonPonowienie();

  if (!wynik.udany || wynik.wynik === undefined) {
    const opis = opisOdmowyBledu(czynnosc, wynik.blad, wynik.nieznanyTyp);
    // Komunikat zostaje po naciśnięciu ponowienia: ponowienie wyzwala wywołanie, nie sprząta ekranu.
    okno.blad(opis);
    odpowiedz.pokaz(opis, false);
    return;
  }
  kontekst.odswiez();
  const skutek = opisWydania(zlecenie, wynik.wynik);
  // Do historii wchodzi zdanie o skutku, nie o zamówieniu, widoczne w wykazie nawet bez pliku.
  kontekst.odnotujWydanie(skutek.zdanie);
  pokazKomunikat({
    tytul: skutek.udany ? 'Raport wydany' : 'Rdzeń przyjął zlecenie, pliku nie oddał',
    tresc: skutek.zdanie,
    waga: skutek.udany ? 'sukces' : 'ostrz',
  });
  odpowiedz.pokaz(skutek.zdanie, skutek.udany);
}

/**
 * Skutek wydania nazwany odpowiedzią rdzenia, a nie zamówieniem operatora,
 * bo rdzeń oddaje sam format, bez ścieżki pliku.
 */
function opisWydania(
  zlecenie: ZlecenieEksportu,
  wynik: ResearchReportExportResponse,
): { zdanie: string; udany: boolean } {
  const zamowionoMiejsce = zlecenie.sciezka.trim() !== '' || zlecenie.doRepozytorium;
  const braki: string[] = [];
  if (zlecenie.sciezka.trim() !== '' && (wynik.path ?? '') === '') braki.push('ścieżki pliku');
  if (zlecenie.doRepozytorium && (wynik.libraryFileId ?? '') === '') {
    braki.push('dokumentu repozytorium');
  }

  if (braki.length === 0) {
    return {
      zdanie: zamowionoMiejsce
        ? `Rdzeń wydał raport w formacie ${zlecenie.format}${opisMiejsca(wynik)}.`
        : `Rdzeń przyjął zlecenie wydania w formacie ${zlecenie.format}. ` +
          'Miejsca docelowego nie wskazano, więc nie ma czego potwierdzać poza samym zleceniem.',
      udany: true,
    };
  }
  return {
    zdanie:
      `Rdzeń przyjął zlecenie wydania w formacie ${zlecenie.format}, ale NIE ODDAŁ ${braki.join(' ani ')} — ` +
      'zapisał sam ślad zlecenia, pliku pod wskazanym miejscem nie ma. ' +
      'Rdzeń nie ma magazynu plików wyjściowych (zmierzone), więc ponowienie tego nie zmieni.',
    udany: false,
  };
}

/** Miejsce oddane przez rdzeń — wypisywane wyłącznie z pól odpowiedzi, nigdy z zamówienia tego okna panelu. */
function opisMiejsca(wynik: ResearchReportExportResponse): string {
  const czesci: string[] = [];
  if ((wynik.path ?? '') !== '') czesci.push(`plik ${String(wynik.path)}`);
  if ((wynik.libraryFileId ?? '') !== '') {
    czesci.push(`dokument repozytorium ${String(wynik.libraryFileId)}`);
  }
  return czesci.length === 0 ? '' : `, ${czesci.join(', ')}`;
}

/**
 * Trzy stany podglądu eksportu: pytam, mam dokument, nie mam czego wydać,
 * z treścią pustki w osobnym pliku.
 */
export function ustawStanEksportu(kontekst: KontekstEksportu): void {
  const { stan, okno } = kontekst;
  if (stan.faza() === 'odczyt') {
    okno.ladowanie('Odczyt okna badania z rdzenia…');
    return;
  }
  if (stan.raport() !== null) {
    okno.gotowe();
    return;
  }
  if (stan.faza() === 'blad') {
    okno.blad(stan.powod());
    return;
  }
  pokazPustke(okno, stan, PUSTKA_EKSPORTU);
}

/**
 * Tekst swobodny okna przekazywany komendom bez własnego formularza, bez
 * wskazania nazywanego brakiem.
 */
function tekstDlaKomendy(_kontekst: KontekstEksportu): string {
  // Panel eksportu nie ma pola tekstowego; komendy tej rodziny biorą wskazania ze stanu badania.
  return '';
}

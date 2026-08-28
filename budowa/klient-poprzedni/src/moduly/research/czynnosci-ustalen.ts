import { ResearchFindingStatus, type ResearchFinding } from '../../../../shared/contract';
import { opisOdmowyBledu } from '../../komponenty/odmowa';
import type { WierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import type { AkcjaBadania } from './akcje-okien';
import { KODY_OKIEN } from './kody-okien';
import { pokazPustke, PUSTKA_USTALEN } from './pustka-okien';
import type { StanBadania } from './stan-badania';
import { czyKomendaBadania, wykonajKomendeBadania } from './wywolania-komend';
import type { StanOknaBadania } from './stan-okna-badania';

/** Czynności operatora w Findings Panel: zapis, powiązanie ze źródłem i edycja ustalenia w jednej komendzie. */
export interface KontekstUstalen {
  stan: StanBadania;
  okno: StanOknaBadania;
  odpowiedz: WierszOdpowiedzi;
  tresc: HTMLTextAreaElement;
  /** Identyfikator ustalenia w edycji; pusty zakłada ustalenie nowe. */
  edytowane(): string;
  /** Wciąga ustalenie do formularza albo czyści go pod zapis nowy. */
  wciagnij(identyfikator: string, tresc: string): void;
  przejdz(kodOkna: string): void;
}

/** Rozdziela akcję panelu na drogę własną okna i drogę generyczną, wspólną dla całej rodziny komend ustaleń. */
export async function wykonajAkcjeUstalen(
  kontekst: KontekstUstalen,
  akcja: AkcjaBadania,
): Promise<void> {
  if (akcja.kod === 'research.finding.resolve') {
    await zapiszUstalenie(kontekst, ResearchFindingStatus.Resolved);
    return;
  }
  if (akcja.kod === 'research.finding.toReport') {
    kontekst.przejdz(KODY_OKIEN.raport);
    kontekst.odpowiedz.pokaz(
      `Zaznaczono ${kontekst.stan.wybraneUstalenia.wybrane().length} ustaleń — Report Builder weźmie je w pole findingIds.`,
      true,
    );
    return;
  }
  if (akcja.droga === 'okno') {
    if (akcja.kod === 'research.finding.new') kontekst.wciagnij('', '');
    kontekst.tresc.focus();
    kontekst.odpowiedz.pokaz(`„${akcja.nazwa}" — formularz ustalenia jest wyżej.`, true);
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

/** Droga generyczna: `window.action` z zaznaczonymi ustaleniami w parametrach żądania tego okna panelu. */
async function przezPanelAkcji(kontekst: KontekstUstalen, akcja: AkcjaBadania): Promise<void> {
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
    findingIds: stan.wybraneUstalenia.wybrane(),
  });
  if (!wynik.udany) {
    odpowiedz.pokaz(opisOdmowyBledu(`Akcja „${akcja.nazwa}"`, wynik.blad, wynik.nieznanyTyp), false);
    return;
  }
  // Rdzeń oddaje sukces window.action tylko wtedy, gdy akcję wykonał, nie samo wywołanie okna.
  odpowiedz.pokaz(`Rdzeń oddał wynik akcji „${akcja.nazwa}".`, true);
}

/** Zapis ustalenia; obecność `findingId` zamienia zapis w zmianę istniejącego ustalenia, nie zapis nowego. */
export async function zapiszUstalenie(
  kontekst: KontekstUstalen,
  nowyStan: ResearchFindingStatus,
): Promise<void> {
  const { stan, okno, odpowiedz, tresc } = kontekst;
  if (tresc.value.trim() === '') {
    tresc.focus();
    odpowiedz.pokaz('Wpisz treść ustalenia — pole content jest w kontrakcie obowiązkowe.', false);
    return;
  }
  if (stan.idOkna() === '') {
    odpowiedz.pokaz(
      'Rdzeń nie wskazał okna badania; komenda research.finding.add wymaga pola windowId.',
      false,
    );
    return;
  }
  okno.ladowanie('Zapis ustalenia…');
  const zamowioneZrodla = stan.wybraneZrodla.wybrane();
  const wynik = await stan.zrodlo.zapiszUstalenie({
    idOkna: stan.idOkna(),
    tresc: tresc.value,
    idUstalenia: kontekst.edytowane(),
    idZrodel: zamowioneZrodla,
    stan: nowyStan,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    const opis = opisOdmowyBledu('Zapis ustalenia', wynik.blad, wynik.nieznanyTyp);
    okno.blad(opis);
    odpowiedz.pokaz(opis, false);
    return;
  }
  stan.wchlonUstalenie(wynik.wynik.finding);
  kontekst.wciagnij('', '');
  const skutek = opisZapisu(wynik.wynik.finding, zamowioneZrodla, nowyStan);
  odpowiedz.pokaz(skutek.zdanie, skutek.udany);
}

/**
 * Skutek zapisu nazwany ustaleniem, które wróciło, a nie formularzem, który
 * poszedł, bo powiązanie może przepaść.
 */
function opisZapisu(
  ustalenie: ResearchFinding,
  zamowioneZrodla: readonly string[],
  zamowionyStan: ResearchFindingStatus,
): { zdanie: string; udany: boolean } {
  const czesci: string[] = [];
  let udany = true;

  // Stan w ustaleniu, które wróciło, ma prawo różnić się od stanu zamówionego rozstrzygnięciem.
  if (ustalenie.status !== zamowionyStan) {
    czesci.push(
      `Rdzeń NIE nadał ustaleniu stanu ${zamowionyStan} — oddał je w stanie ${ustalenie.status}.`,
    );
    udany = false;
  }

  if (zamowioneZrodla.length > 0) {
    // Rdzeń oddaje sourceIds w innej kolejności niż zamówiona, więc porównanie idzie po przynależności.
    const zwiazane = new Set(ustalenie.sourceIds ?? []);
    const pominiete = zamowioneZrodla.filter((kod) => !zwiazane.has(kod));
    if (pominiete.length > 0) {
      czesci.push(
        'Zaznaczonych źródeł, których rdzeń NIE związał z tym ustaleniem, jest ' +
          `${String(pominiete.length)} z ${String(zamowioneZrodla.length)} — ` +
          'powiązania nie ma w ustaleniu, które oddał.',
      );
      udany = false;
    } else {
      czesci.push(`Rdzeń związał z nim wszystkie zaznaczone źródła (${String(zamowioneZrodla.length)}).`);
    }
  }

  return {
    zdanie: [`Rdzeń zapisał ustalenie ${ustalenie.id}.`, ...czesci].join(' '),
    udany,
  };
}

/**
 * Trzy stany wykazu ustaleń: pytam, mam treść, nie mam czego pokazać, z
 * treścią pustki w osobnym pliku.
 */
export function ustawStanUstalen(kontekst: KontekstUstalen, liczba: number): void {
  const { stan, okno } = kontekst;
  if (stan.faza() === 'odczyt') {
    okno.ladowanie('Odczyt okna badania z rdzenia…');
    return;
  }
  if (liczba > 0) {
    okno.gotowe();
    return;
  }
  if (stan.faza() === 'blad') {
    okno.blad(stan.powod());
    return;
  }
  pokazPustke(okno, stan, PUSTKA_USTALEN);
}

/**
 * Tekst swobodny okna przekazywany komendom bez własnego formularza, bez
 * wskazania nazywanego brakiem.
 */
function tekstDlaKomendy(kontekst: KontekstUstalen): string {
  return kontekst.tresc.value.trim();
}

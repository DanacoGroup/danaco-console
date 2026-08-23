import { opisOdmowyBledu } from '../../komponenty/odmowa';
import type { WierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import type { AkcjaBadania } from './akcje-okien';
import { KODY_OKIEN } from './kody-okien';
import { pokazPustke, PUSTKA_ZAKRESU } from './pustka-okien';
import type { StanBadania } from './stan-badania';
import { czyKomendaBadania, wykonajKomendeBadania } from './wywolania-komend';
import type { StanOknaBadania } from './stan-okna-badania';

/**
 * Czynności okna wiodącego Research Workspace: zapis zakresu badania i rozdział
 * akcji panelu.
 *
 * Dwie akcje kończą się w tym oknie albo komendą `research.workspace.set`.
 * Pozostałe — pytania badawcze, notatka robocza, luki, świeżość, cztery
 * reprezentacje materiału i diagram przesiewu — własnej komendy nie mają
 * i idą generycznym `window.action`, którego odmowę okno wypisuje słowami
 * rdzenia.
 */
export interface KontekstZakresu {
  stan: StanBadania;
  okno: StanOknaBadania;
  odpowiedz: WierszOdpowiedzi;
  zakres: HTMLTextAreaElement;
  etapy: HTMLTextAreaElement;
  przejdz(kodOkna: string): void;
  odswiez(): void;
}

/** Rozdziela akcję panelu na drogę własną okna i drogę generyczną. */
export async function wykonajAkcjeZakresu(
  kontekst: KontekstZakresu,
  akcja: AkcjaBadania,
): Promise<void> {
  if (akcja.kod === 'research.scope.investigate') {
    await ustawZakres(kontekst);
    // Zbieranie materiału zaczyna się od wyszukania, nie od katalogowania:
    // porządek pracy z opracowania (rozdz. 4.1) prowadzi Discovery Panel →
    // Sources Manager, a nie wprost do formularza źródła.
    kontekst.przejdz(KODY_OKIEN.odkrywanie);
    return;
  }
  if (akcja.kod === 'research.scope.edit') {
    kontekst.zakres.focus();
    kontekst.odpowiedz.pokaz(
      'Zakres badania jest w polu poniżej — po zmianie naciśnij „Zapisz zakres badania".',
      true,
    );
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

/** Droga generyczna: `window.action` z zakresem badania w parametrach. */
async function przezPanelAkcji(kontekst: KontekstZakresu, akcja: AkcjaBadania): Promise<void> {
  const { stan, odpowiedz } = kontekst;
  if (stan.idOkna() === '') {
    odpowiedz.pokaz(
      `Akcja „${akcja.nazwa}" wymaga okna badania, którego rdzeń jeszcze nie wskazał.`,
      false,
    );
    return;
  }
  odpowiedz.pokaz(`Wysłano window.action ${akcja.kod}…`, true);
  const wynik = await stan.okna.akcja(stan.idOkna(), akcja.kod, { scope: stan.zakres() });
  if (!wynik.udany) {
    odpowiedz.pokaz(opisOdmowyBledu(`Akcja „${akcja.nazwa}"`, wynik.blad, wynik.nieznanyTyp), false);
    return;
  }
  odpowiedz.pokaz(`Rdzeń oddał wynik akcji „${akcja.nazwa}".`, true);
}

/** Zapis zakresu i etapów komendą `research.workspace.set`. */
export async function ustawZakres(kontekst: KontekstZakresu): Promise<void> {
  const { stan, okno, odpowiedz } = kontekst;
  const zakres = kontekst.zakres.value.trim();
  if (zakres === '') {
    kontekst.zakres.focus();
    odpowiedz.pokaz('Wpisz zakres badania — rdzeń odmówi zapisu bez treści pola scope.', false);
    return;
  }
  const etapy = kontekst.etapy.value
    .split('\n')
    .map((wiersz) => wiersz.trim())
    .filter((wiersz) => wiersz !== '');

  okno.ladowanie('Zapis zakresu badania…');
  odpowiedz.pokaz('Wysłano research.workspace.set — czekam na odpowiedź rdzenia.', true);

  const wynik = await stan.zrodlo.ustawZakres({ zakres, etapy });
  if (!wynik.udany || wynik.wynik === undefined) {
    const opis = opisOdmowyBledu('Zapis zakresu badania', wynik.blad, wynik.nieznanyTyp);
    okno.blad(opis);
    odpowiedz.pokaz(opis, false);
    return;
  }
  stan.wchlonZakres(wynik.wynik.scope, wynik.wynik.stages);
  const skutek = opisZapisuZakresu(wynik.wynik.scope, wynik.wynik.stages, zakres, etapy);
  odpowiedz.pokaz(skutek.zdanie, skutek.udany);
  kontekst.odswiez();
}

/**
 * Skutek zapisu zakresu — porównanie zamówienia z odpowiedzią, nie samo
 * przepisanie odpowiedzi.
 *
 * Rdzeń ma prawo zapisać co innego niż przyszło (przyciąć zakres, odsiać etap),
 * a przy samym zdaniu potwierdzającym podmiana byłaby niewidoczna. Porównanie
 * kosztuje jeden przebieg po wykazie i nazywa rozbieżność wprost.
 */
function opisZapisuZakresu(
  oddanyZakres: string,
  oddaneEtapy: readonly string[],
  zamowionyZakres: string,
  zamowioneEtapy: readonly string[],
): { zdanie: string; udany: boolean } {
  const czesci: string[] = [];
  let udany = true;

  if (oddanyZakres !== zamowionyZakres) {
    czesci.push(
      `Rdzeń zapisał zakres INNY niż wysłany — wysłano „${zamowionyZakres}", oddał „${oddanyZakres}".`,
    );
    udany = false;
  }
  const pominiete = zamowioneEtapy.filter((etap) => !oddaneEtapy.includes(etap));
  if (pominiete.length > 0) {
    czesci.push(
      'Wysłanych etapów, których rdzeń NIE zapisał, jest ' +
        `${String(pominiete.length)} z ${String(zamowioneEtapy.length)}.`,
    );
    udany = false;
  }
  if (czesci.length > 0) return { zdanie: czesci.join(' '), udany };

  return {
    zdanie:
      `Rdzeń potwierdził zakres badania: ${oddanyZakres}` +
      (oddaneEtapy.length === 0 ? '.' : ` — etapów ${String(oddaneEtapy.length)}.`),
    udany: true,
  };
}

/**
 * Stany okna wiodącego: odczyt, błąd, zakres wpisany, zakres jeszcze pusty.
 *
 * Stan pusty opisuje `pustka-okien.ts` wspólnie dla okien modułu, a nie
 * `stan.powod()`, które w tej fazie niesie komunikat rdzenia o niewskazanym
 * oknie badania — zdanie techniczne w miejscu opisu pustego zakresu.
 */
export function ustawStanZakresu(kontekst: KontekstZakresu): void {
  const { stan, okno } = kontekst;
  if (stan.faza() === 'odczyt') {
    okno.ladowanie('Odczyt okna badania z rdzenia…');
    return;
  }
  if (stan.faza() === 'blad') {
    okno.blad(stan.powod());
    return;
  }
  if (stan.zakres() !== '') {
    okno.gotowe();
    return;
  }
  pokazPustke(okno, stan, PUSTKA_ZAKRESU);
}

/**
 * Tekst swobodny okna przekazywany komendom bez własnego formularza.
 *
 * Żądanie składane bez wskazania Operatora wracałoby odmową walidacji, z której
 * nic dla niego nie wynika. Ten jeden krok mówi, skąd okno bierze treść — i gdy
 * jej nie ma, `wywolania-komend.ts` nazywa brak, zamiast wysyłać puste pole.
 */
function tekstDlaKomendy(kontekst: KontekstZakresu): string {
  // Pole zakresu niesie temat badania — z niego biorą treść pytania badawcze
  // i notatka robocza, wypisywane po jednym w wierszu.
  return kontekst.zakres.value.trim();
}

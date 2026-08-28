import { opisOdmowyBledu } from '../../komponenty/odmowa';
import type { WierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import type { AkcjaBadania } from './akcje-okien';
import { KODY_OKIEN } from './kody-okien';
import { pokazPustke, PUSTKA_ZAKRESU } from './pustka-okien';
import type { StanBadania } from './stan-badania';
import { czyKomendaBadania, wykonajKomendeBadania } from './wywolania-komend';
import type { StanOknaBadania } from './stan-okna-badania';

/**
 * Plik niesie czynności okna wiodącego Research Workspace: zapis zakresu
 * badania komendą research.workspace.set oraz rozdział pozostałych akcji panelu.
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

/**
 * Rozdziela akcję panelu zakresu na drogę własną okna, obsługiwaną komendą
 * dedykowaną, oraz drogę generyczną przekazywaną rdzeniowi jako window.action.
 */
export async function wykonajAkcjeZakresu(
  kontekst: KontekstZakresu,
  akcja: AkcjaBadania,
): Promise<void> {
  if (akcja.kod === 'research.scope.investigate') {
    await ustawZakres(kontekst);
    // Zbieranie materiału zaczyna się od wyszukania, nie od katalogowania.
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

/**
 * Przekazuje kod akcji oraz zakres badania do rdzenia komendą window.action,
 * bez własnej obsługi po stronie okna Research Workspace.
 */
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

/**
 * Zapisuje zakres badania wraz z etapami komendą research.workspace.set,
 * przekazując treść pobraną z pola zakresu okna wiodącego.
 */
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
 * Porównuje zamówiony zakres z odpowiedzią rdzenia, zamiast przepisywać samą
 * odpowiedź, i nazywa rozbieżność wprost, gdy rdzeń zapisał co innego.
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
 * Ustawia jeden z czterech stanów okna wiodącego: odczyt trwa, wystąpił błąd,
 * zakres jest wpisany albo zakres pozostaje pusty.
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
 * Podaje tekst swobodny okna przekazywany komendom, które nie mają własnego
 * formularza; brak treści nazywa plik wywolania-komend.ts.
 */
function tekstDlaKomendy(kontekst: KontekstZakresu): string {
  // Pole zakresu niesie temat badania, z którego biorą treść pytania badawcze.
  return kontekst.zakres.value.trim();
}

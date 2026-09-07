// Argumenty i wytwory debaty w oknie Roundtable: oznaczenie argumentu
// kluczowego, wydanie grafu, macierz decyzyjna, katalog błędów i nagranie mowy.
import { Command, RoundtableArgumentFormat } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import {
  odmowaCzynnosci,
  polePasa,
  postawPas,
  potwierdzone,
  przyciskPasa,
  wartoscPola,
} from './roundtable-pas.ts';

const NAGLOWEK = 'Debata';
const SELEKTOR_NAZWY = '[data-argumenty-nazwa]';

export function zwiazArgumenty(
  kanal: Kanal,
  korzen: Element,
  idOkna: () => string,
  odswiez: () => void,
  przy: AddEventListenerOptions,
): void {
  postawPas(korzen, 'panel-debate', [
    polePasa(korzen, 'argumentyNazwa', 'Nazwa macierzy decyzyjnej'),
    przyciskPasa(korzen, 'argumenty', 'oznacz', 'Oznacz argument kluczowy'),
    przyciskPasa(korzen, 'argumenty', 'wydaj-graf', 'Wydaj graf argumentów'),
    przyciskPasa(korzen, 'argumenty', 'macierz', 'Macierz decyzyjna'),
    przyciskPasa(korzen, 'argumenty', 'zapisz-macierz', 'Załóż macierz'),
    przyciskPasa(korzen, 'argumenty', 'bledy', 'Katalog błędów'),
    przyciskPasa(korzen, 'argumenty', 'wlacz-bledy', 'Włącz cały katalog'),
    przyciskPasa(korzen, 'argumenty', 'mowa', 'Nagraj mowę'),
  ]);
  korzen.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const czynnosc = cel.closest<HTMLElement>('[data-argumenty]')?.dataset.argumenty;
    if (czynnosc === undefined) return;
    zdarzenie.stopPropagation();
    void wykonaj(kanal, korzen, czynnosc, idOkna(), odswiez);
  }, przy);
}

async function wykonaj(
  kanal: Kanal,
  korzen: Element,
  czynnosc: string,
  idOkna: string,
  odswiez: () => void,
): Promise<void> {
  if (idOkna === '') {
    oglos(NAGLOWEK, 'Rdzeń nie dał okna debaty dla tej karty.', 'ostrzezenie');
    return;
  }
  if (czynnosc === 'oznacz') return oznaczArgument(kanal, idOkna, odswiez);
  if (czynnosc === 'wydaj-graf') return wydajGraf(kanal, idOkna);
  if (czynnosc === 'macierz') return opiszMacierz(kanal, idOkna);
  if (czynnosc === 'zapisz-macierz') return zalozMacierz(kanal, korzen, idOkna);
  if (czynnosc === 'bledy') return opiszKatalogBledow(kanal, idOkna);
  if (czynnosc === 'wlacz-bledy') return wlaczKatalogBledow(kanal, korzen, idOkna);
  if (czynnosc === 'mowa') return nagrajMowe(kanal, idOkna);
  odmowaCzynnosci(NAGLOWEK, czynnosc);
}

/* Oznaczenie idzie na węzeł pierwszy grafu: okno nie prowadzi wskazania węzła,
   a komenda bez węzła odmawia. */
async function oznaczArgument(kanal: Kanal, idOkna: string, odswiez: () => void): Promise<void> {
  const graf = await wywolaj(kanal, Command.RoundtableArgumentList, { windowId: idOkna });
  const wezel = graf.wynik?.graph.nodes[0]?.id ?? '';
  if (wezel === '') {
    oglos(NAGLOWEK, 'Graf argumentów nie ma węzła do oznaczenia.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.RoundtableArgumentPin, {
    windowId: idOkna,
    nodeId: wezel,
    pinned: true,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił oznaczenia argumentu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Argument oznaczony jako kluczowy.');
  odswiez();
}

/* Graf wraca jako wytwór magazynu, nie jako treść, więc okno nazywa jego
   oznaczenie zamiast udawać zapis do schowka. */
async function wydajGraf(kanal: Kanal, idOkna: string): Promise<void> {
  const wynik = await wywolaj(kanal, Command.RoundtableArgumentExport, {
    windowId: idOkna,
    format: RoundtableArgumentFormat.Argdown,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił wydania grafu argumentów.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Graf argumentów zapisany jako wytwór ${wynik.wynik.artifactId}.`);
}

async function opiszMacierz(kanal: Kanal, idOkna: string): Promise<void> {
  const wynik = await wywolaj(kanal, Command.RoundtableDecisionMatrixGet, { windowId: idOkna });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń nie podał macierzy decyzyjnej.', 'ostrzezenie');
    return;
  }
  const macierz = wynik.wynik.matrix;
  oglos(NAGLOWEK, `Macierz „${macierz.name}”: ${macierz.criteria.length} kryteriów `
    + `wobec ${macierz.options.length} wariantów.`);
}

/* Macierz zakładana z okna jest pusta: kryteria i warianty wpisuje Operator
   w module decyzyjnym, a zmyślona zawartość byłaby decyzją za niego. */
async function zalozMacierz(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const nazwa = wartoscPola(korzen, SELEKTOR_NAZWY);
  if (nazwa === '') {
    oglos(NAGLOWEK, 'Macierz decyzyjna musi mieć nazwę wpisaną w polu.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.RoundtableDecisionMatrixSet, {
    windowId: idOkna,
    name: nazwa,
    criteria: [],
    options: [],
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił założenia macierzy.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Macierz decyzyjna „${nazwa}” założona.`);
}

async function opiszKatalogBledow(kanal: Kanal, idOkna: string): Promise<void> {
  const wynik = await wywolaj(kanal, Command.RoundtableFallacyCatalogGet, { windowId: idOkna });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń nie podał katalogu błędów.', 'ostrzezenie');
    return;
  }
  const wlaczone = wynik.wynik.definitions.filter((d) => d.enabled);
  oglos(NAGLOWEK, `Katalog błędów logicznych: ${wlaczone.length} czynnych `
    + `z ${wynik.wynik.definitions.length} znanych rdzeniowi.`);
}

/* Zapis katalogu nadpisuje wybór czynnych błędów, więc pierwsze naciśnięcie
   uzbraja przycisk zamiast kasować ustawienie Operatora. */
async function wlaczKatalogBledow(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  if (!potwierdzone(korzen, '[data-argumenty="wlacz-bledy"]', 'Potwierdź włączenie',
    'Włącz cały katalog')) return;
  const katalog = await wywolaj(kanal, Command.RoundtableFallacyCatalogGet, { windowId: idOkna });
  const kody = katalog.wynik?.definitions.map((d) => d.code) ?? [];
  if (kody.length === 0) {
    oglos(NAGLOWEK, 'Rdzeń nie zna żadnego błędu logicznego do włączenia.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.RoundtableFallacyCatalogSet, {
    windowId: idOkna,
    codes: kody,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu katalogu błędów.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Katalog błędów włączony w całości: ${kody.length} pozycji.`);
}

/* Nagranie wraca jako wytwór magazynu, więc okno nazywa jego oznaczenie. */
async function nagrajMowe(kanal: Kanal, idOkna: string): Promise<void> {
  const wynik = await wywolaj(kanal, Command.RoundtableSpeechSynthesize, { windowId: idOkna });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił nagrania mowy.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Nagranie debaty zapisane jako wytwór ${wynik.wynik.artifactId}.`);
}

import type { StudioOperation } from '../../../../shared/contract';
import { KATEGORIE_OPERACJI, nazwaOperacji } from './kategorie-operacji';
import type { PozycjaWsadu } from './petla-stan';
import { przycisk } from './zadania-wykaz';

/** Czynności trybu wsadowego sięgające poza ten widok: puszczenie wsadu dokumentów i przerwanie go w połowie. */
export interface CzynnosciWsadu {
  /** Puszcza wsad: dokumenty po kolei, wskazaną czynnością. */
  puscWsad(idDokumentow: readonly string[], idAkcji: string): void;
  /** Przerywa wsad; to, co przetworzone, zostaje przetworzone. */
  przerwij(): void;
}

export interface TrybWsadowy {
  element: HTMLElement;
  odswiez(pozycje: readonly PozycjaWsadu[], przerwany: boolean, wlasne: readonly StudioOperation[]): void;
}

/** Nazwa stanu pozycji wsadu widoczna dla operatora, dla każdego z pięciu stanów niesionych przez kontrakt. */
const NAZWA_STANU_POZYCJI: Record<PozycjaWsadu['stan'], string> = {
  oczekuje: 'oczekuje',
  'w-biegu': 'w realizacji',
  udana: 'gotowe',
  odrzucona: 'odrzucona',
  przerwana: 'przerwana zatrzymaniem',
};

export function utworzTrybWsadowy(czynnosci: CzynnosciWsadu): TrybWsadowy {
  const poleDokumentow = document.createElement('textarea');
  poleDokumentow.className = 'dn-pole petla-wsad__dokumenty';
  poleDokumentow.rows = 3;
  poleDokumentow.placeholder =
    'Identyfikatory dokumentów, po jednym w wierszu albo rozdzielone przecinkami';
  poleDokumentow.setAttribute('aria-label', 'Dokumenty objęte wsadem');

  const wyborCzynnosci = document.createElement('select');
  wyborCzynnosci.className = 'dn-wybor';
  wyborCzynnosci.setAttribute('aria-label', 'Czynność wykonywana na każdym dokumencie');

  const puscGuzik = przycisk('Puść wsad', () => {
    const dokumenty = rozbijWykaz(poleDokumentow.value);
    if (dokumenty.length === 0 || wyborCzynnosci.value === '') return;
    czynnosci.puscWsad(dokumenty, wyborCzynnosci.value);
  });
  const przerwijGuzik = przycisk('Przerwij wsad', () => czynnosci.przerwij());

  const pas = document.createElement('div');
  pas.className = 'petla-wsad__pas';
  pas.append(wyborCzynnosci, puscGuzik, przerwijGuzik);

  const lista = document.createElement('ul');
  lista.className = 'petla-wsad__wykaz';
  lista.setAttribute('aria-label', 'Wynik wsadu, dokument po dokumencie');

  const bilans = document.createElement('p');
  bilans.className = 'dn-tekst-3 petla-wsad__bilans';

  const element = document.createElement('section');
  element.className = 'petla-wsad';
  element.append(poleDokumentow, pas, bilans, lista);

  function wypelnijWybor(wlasne: readonly StudioOperation[]): void {
    const wybrane = wyborCzynnosci.value;
    wyborCzynnosci.replaceChildren();
    const puste = document.createElement('option');
    puste.value = '';
    puste.textContent = '— wybierz czynność —';
    wyborCzynnosci.append(puste);
    for (const kategoria of KATEGORIE_OPERACJI) {
      const grupa = document.createElement('optgroup');
      grupa.label = kategoria.nazwa;
      for (const operacja of kategoria.operacje) {
        const pozycja = document.createElement('option');
        pozycja.value = operacja.id;
        pozycja.textContent = operacja.nazwa;
        grupa.append(pozycja);
      }
      wyborCzynnosci.append(grupa);
    }
    if (wlasne.length > 0) {
      const grupa = document.createElement('optgroup');
      grupa.label = 'Operacje własne Operatora';
      for (const operacja of wlasne) {
        const pozycja = document.createElement('option');
        pozycja.value = operacja.id;
        pozycja.textContent = `${operacja.name} (własna)`;
        grupa.append(pozycja);
      }
      wyborCzynnosci.append(grupa);
    }
    wyborCzynnosci.value = wybrane;
  }

  function wiersz(pozycja: PozycjaWsadu): HTMLElement {
    const wpis = document.createElement('li');
    wpis.className = 'dn-karta petla-wsad__wiersz';
    wpis.dataset['stan'] = pozycja.stan;

    const nazwa = document.createElement('span');
    nazwa.className = 'petla-wsad__dokument';
    nazwa.textContent = pozycja.idDokumentu;

    const stan = document.createElement('span');
    stan.className = 'dn-plakietka';
    stan.textContent = NAZWA_STANU_POZYCJI[pozycja.stan];

    const naglowek = document.createElement('div');
    naglowek.className = 'petla-wsad__naglowek';
    naglowek.append(nazwa, stan);
    if (pozycja.stan === 'w-biegu') {
      const obrot = document.createElement('span');
      obrot.className = 'dn-spinner';
      obrot.setAttribute('role', 'status');
      obrot.setAttribute('aria-label', 'Dokument w realizacji');
      naglowek.append(obrot);
    }
    wpis.append(naglowek);

    // Powód odrzucenia stoi przy dokumencie, bo sama liczba odrzuconych byłaby przemilczeniem straty.
    if (pozycja.powod !== '') {
      const powod = document.createElement('p');
      powod.className = 'dn-tekst-3 petla-wsad__powod';
      powod.textContent = pozycja.powod;
      wpis.append(powod);
    }
    return wpis;
  }

  return {
    element,

    odswiez(pozycje, przerwany, wlasne) {
      wypelnijWybor(wlasne);
      lista.replaceChildren();
      for (const pozycja of pozycje) lista.append(wiersz(pozycja));

      const wBiegu = pozycje.some(
        (pozycja) => pozycja.stan === 'w-biegu' || pozycja.stan === 'oczekuje',
      );
      puscGuzik.disabled = wBiegu;
      przerwijGuzik.disabled = !wBiegu || przerwany;

      if (pozycje.length === 0) {
        bilans.textContent =
          'Wsadu jeszcze nie było. Wskaż dokumenty i czynność — każdy dokument dostanie ' +
          'własny wiersz z własnym wynikiem, a wsad da się przerwać w połowie bez ' +
          'wycofywania tego, co już weszło.';
        return;
      }
      const udane = pozycje.filter((pozycja) => pozycja.stan === 'udana').length;
      const odrzucone = petlaIleOdrzuconych(pozycje);
      const przerwane = pozycje.filter((pozycja) => pozycja.stan === 'przerwana').length;
      bilans.textContent =
        `Bilans wsadu: ${udane} gotowych, ${odrzucone} odrzuconych, ` +
        `${przerwane} przerwanych z ${pozycje.length} dokumentów.` +
        (przerwany ? ' Wsad przerwany decyzją Operatora — to, co weszło, zostało.' : '');
    },
  };
}

/** Liczy, ile pozycji wsadu zostało odrzuconych, na podstawie stanu każdej pozycji w całym wykazie wsadu. */
function petlaIleOdrzuconych(pozycje: readonly PozycjaWsadu[]): number {
  return pozycje.filter((pozycja) => pozycja.stan === 'odrzucona').length;
}

/** Rozbija wpisany wykaz dokumentów na pojedyncze identyfikatory, przyjmując zarówno wiersze, jak i przecinki. */
export function rozbijWykaz(tresc: string): string[] {
  const znalezione = tresc
    .split(/[\n,;]/)
    .map((wpis) => wpis.trim())
    .filter((wpis) => wpis !== '');
  return [...new Set(znalezione)];
}

/** Nazwa czynności dla wiersza wsadu, odczytana z katalogu operacji albo, gdy jej tam nie ma, sam identyfikator. */
export function nazwaCzynnosciWsadu(idAkcji: string): string {
  return nazwaOperacji(idAkcji) ?? idAkcji;
}

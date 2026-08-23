import type { StudioRecognizedWord } from '../../../../shared/contract';
import { opisOdmowyBledu } from '../../komponenty/odmowa';
import { przycisk, utworzWierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import type { Wynik } from '../../protokol/kanal';
import { DROGI_CYFRYZACJI, NASTAWY_CYFRYZACJI, SKLADNIKI_PAKIETU_SERWERA } from './braki-cyfryzacji';
import {
  dolozArchiwum,
  dolozMaterial,
  odczytajKolejke,
  odczytajUrzadzenia,
  poprawSlowo,
  przyjmijWynik,
  rozpoznajKolejke,
  rozpoznajPozycje,
  type ZapleczeCyfryzacji,
} from './czynnosci-cyfryzacji';
import { utworzKolejkeCyfryzacji } from './kolejka-cyfryzacji';
import { utworzOknoStudio } from './okno-studio';
import { utworzPolaCyfryzacji } from './pola-cyfryzacji';
import type { StanStudio } from './stan-studio';
import { utworzWierszCyfryzacji } from './wiersz-cyfryzacji';
import type { ZrodloDokumentuStudio } from './zrodlo-dokumentu-studio';
import {
  POPRAWKI_OBRAZU,
  type WskazanieObrazu,
  type ZrodloMaterialuStudio,
} from './zrodlo-materialu-studio';
import type { ZrodloWstawienStudio } from './zrodlo-wstawien-studio';

/** Nagłówki kolumn kolejki — zostają widoczne także w stanie pustym. */
const KOLUMNY = ['Wskazanie', 'Rodzaj', 'Stan', 'Wynik rdzenia', 'Czynność'];

/**
 * Narzędziownia cyfryzacji — kolejka RDZENIA, rozpoznanie ze sterowaniem,
 * poprawka słów i przyjęcie wyniku jako dokumentu.
 *
 * ── Co się zmieniło ─────────────────────────────────────────────────────────
 * Panel stał na `document.text.extract`: jedno wydobycie na wywołanie, jeden
 * język, bez silnika, bez progu pewności, z kolejką prowadzoną w oknie i ginącą
 * z jego odświeżeniem. Prowadzi teraz rodzinę `studio.ingest.*` — sześć komend:
 * kolejka po stronie rdzenia (`queue.add`, `queue.list`), rozpoznanie z pełnym
 * sterowaniem (`recognize`), poprawka rozpoznanego słowa przed przyjęciem
 * (`correction.set`), przyjęcie wyniku jako dokumentu wraz z pierwszą wersją
 * (`item.accept`) i wykaz urządzeń wejściowych (`device.list`).
 *
 * ── Poprawianie słów idzie dwiema powierzchniami ────────────────────────────
 * Wymaganie wprost: poprawianie wymaga pracy na obrazie obok tekstu, nie ciasnego
 * paska. Panel stawia obok siebie wykaz słów wraz z ich położeniem na stronie
 * i pewnością rozpoznania oraz warstwę tekstową pozycji. Obrazu skanu tu NIE MA
 * i panel mówi to wprost: komendy pobierającej bajty zasobu do przeglądarki
 * kontrakt nie niesie, więc położenie słowa jest podane liczbami — strona,
 * odsunięcie i rozmiar pola — a nie zaznaczone na obrazku.
 *
 * ── Powierzchnia należy do dokumentu ────────────────────────────────────────
 * Narzędziownia nie jest stałą kolumną: wchodzi przyciskiem jako nakładka
 * i schodzi. Kolejka z wieloma pozycjami potrzebuje miejsca na wykaz, więc
 * dostaje nakładkę, a nie pasek, który zabierałby szerokość także wtedy, gdy
 * nikt nic nie cyfryzuje.
 */
export interface OknoIngestOcrPanel {
  element: HTMLElement;
  odswiez(): void;
  /** Otwiera albo zamyka narzędziownię — dla przycisku spoza tego pliku. */
  przestawWidocznosc(): void;
}

const OBJASNIENIE =
  'Cyfryzacja materiału wejściowego rodziną studio.ingest.*: kolejka wczytywania stoi po stronie ' +
  'rdzenia, rozpoznanie tekstu ma sterowanie silnikiem, zestawem języków i progiem pewności, ' +
  'rozpoznane słowa da się poprawić PRZED przyjęciem, a przyjęcie zakłada dokument roboczy wraz ' +
  'z pierwszą wersją w repozytorium sesji.';

const BEZ_POZYCJI =
  'Kolejka wczytywania rdzenia jest pusta. Wskaż materiał — zasób magazynu rdzenia albo ścieżkę ' +
  'widzianą przez rdzeń — i dołóż go do kolejki. Rozpoznanie wydobędzie z niego tekst, a przyjęcie ' +
  'złoży z gotowych pozycji dokument roboczy.';

const BEZ_ZRODLA_INGEST =
  'Źródło rodziny studio.ingest.* nie zostało podane tej narzędziowni przy montażu modułu, więc ' +
  'okno nie ma czym wywołać kolejki rdzenia. To brak MONTAŻU, nie brak kontraktu ani rdzenia: ' +
  'sześć komend czeka gotowych. Narzędziownia dostaje źródło czwartym argumentem swojej wytwórni.';

export function utworzOknoIngestOcrPanel(
  stan: StanStudio,
  dokumenty: ZrodloDokumentuStudio,
  material: ZrodloMaterialuStudio,
  wstawienia?: ZrodloWstawienStudio,
): OknoIngestOcrPanel {
  const rama = utworzOknoStudio({
    kod: 'studio.ingest-ocr-panel',
    tytul: 'Narzędziownia cyfryzacji',
    rola: 'pomocnicze',
    objasnienie: OBJASNIENIE,
  });

  const pola = utworzPolaCyfryzacji();
  const kolejka = utworzKolejkeCyfryzacji();
  const odpowiedz = utworzWierszOdpowiedzi();

  const cyfryzuj = przycisk('Rozpoznaj kolejkę', 'dn-btn dn-btn--sm dn-btn--sygnal');
  cyfryzuj.dataset['czynnosc'] = 'cyfryzuj';

  const cialoKolejki = document.createElement('tbody');
  const glowa = document.createElement('thead');
  const wierszGlowy = document.createElement('tr');
  for (const kolumna of KOLUMNY) {
    const komorka = document.createElement('th');
    komorka.scope = 'col';
    komorka.textContent = kolumna;
    wierszGlowy.append(komorka);
  }
  glowa.append(wierszGlowy);

  const tabela = document.createElement('table');
  tabela.className = 'dn-tabela ms-kolejka';
  tabela.append(glowa, cialoKolejki);

  const bilans = document.createElement('p');
  bilans.className = 'dn-pole-opis ms-kolejka__bilans';

  /* ── Dwie powierzchnie poprawiania ─────────────────────────────────────── */

  const podglad = document.createElement('pre');
  podglad.className = 'ms-cyfryzacja__podglad';

  const slowa = document.createElement('ul');
  slowa.className = 'ms-cyfryzacja__slowa';
  slowa.setAttribute('aria-label', 'Rozpoznane słowa wraz z położeniem i pewnością');

  const oObrazie = document.createElement('p');
  oObrazie.className = 'dn-pole-opis';
  oObrazie.textContent =
    'Poprawianie idzie dwiema powierzchniami: po jednej stronie słowa wraz ze stroną, odsunięciem, ' +
    'rozmiarem pola i pewnością rozpoznania, po drugiej warstwa tekstowa pozycji. Samego OBRAZU ' +
    'skanu tu nie ma — komendy pobierającej bajty zasobu do przeglądarki kontrakt nie niesie, więc ' +
    'położenie słowa jest podane liczbami, a nie zaznaczone na obrazku. To brak nazwany, nie ' +
    'przemilczany.';

  const powierzchnie = document.createElement('div');
  powierzchnie.className = 'ms-cyfryzacja__powierzchnie';
  powierzchnie.append(slowa, podglad);

  const doEdytora = przycisk(
    'Przyjmij wynik do edytora jako dokument',
    'dn-btn dn-btn--sm dn-btn--atrament',
  );
  doEdytora.dataset['czynnosc'] = 'do-edytora';

  /* ── Zaplecze ──────────────────────────────────────────────────────────── */

  /**
   * Zaplecze czynności powstaje wyłącznie wtedy, gdy źródło rodziny
   * `studio.ingest.*` naprawdę jest. Zamiast udawać wywołanie, panel nazywa brak
   * montażu — i nazywa go jako brak montażu, nie jako brak funkcji rdzenia.
   */
  function zaplecze(): ZapleczeCyfryzacji | null {
    if (wstawienia === undefined) {
      odpowiedz.pokaz(BEZ_ZRODLA_INGEST, false);
      return null;
    }
    return {
      stan,
      wstawienia,
      kolejka,
      pas: rama.stan,
      odpowiedz,
      ustawienia: () => pola.ustawienia(),
      tytul: () => pola.tytul(),
      format: () => pola.format(),
      odswiez: () => odswiez(),
    };
  }

  // ── Czynności okna ─────────────────────────────────────────────────────────

  /** Wskazanie obrazu dla poprawki — bierze się z pozycji wskazanej w kolejce. */
  function wskazanieObrazu(): WskazanieObrazu | null {
    const pozycja = kolejka.pozycja(kolejka.wskazana());
    if (pozycja === null) return null;
    return {
      idZasobu: pozycja.assetId ?? '',
      sciezka: pozycja.assetId === undefined ? pozycja.sourcePath ?? '' : '',
    };
  }

  /**
   * Poprawka obrazu wykonana OSOBNO, przed rozpoznaniem.
   *
   * Czyszczenie w ramach samego rozpoznania idzie nastawami (`deskew`, `denoise`,
   * `binarize`, `trimMargins`) i nie zakłada nowego zasobu. Ta droga zostaje, bo
   * oddaje NOWY zasób pod sumą kontrolną, więc Operator może obejrzeć skutek
   * i wskazać go jako materiał kolejnej pozycji — źródło zostaje nietknięte.
   */
  async function przygotuj(
    nazwa: string,
    wykonaj: (wskazanie: WskazanieObrazu) => Promise<Wynik<{ idZasobu: string }>>,
  ): Promise<void> {
    const wskazanie = wskazanieObrazu();
    if (wskazanie === null) {
      odpowiedz.pokaz(BRAK_WSKAZANEJ, false);
      return;
    }
    rama.stan.ladowanie(`${nazwa} w toku…`);
    const wynik = await wykonaj(wskazanie);
    if (!wynik.udany || wynik.wynik === undefined) {
      rama.stan.blad(opisOdmowyBledu(nazwa, wynik.blad));
      return;
    }
    rama.stan.gotowe();
    const zapleczeCzynnosci = zaplecze();
    if (zapleczeCzynnosci === null) return;
    // Obraz poprawiony wchodzi do kolejki jako POZYCJA NOWA, a nie podmienia
    // pozycji istniejącej: pozycja rdzenia niesie swój wynik rozpoznania i jego
    // ciche przestawienie na inny materiał byłoby podmianą dowodu.
    await dolozMaterial(zapleczeCzynnosci, wynik.wynik.idZasobu, true);
    odpowiedz.pokaz(
      `${nazwa}: rdzeń oddał NOWY zasób ${wynik.wynik.idZasobu} i wszedł on do kolejki jako pozycja ` +
        'osobna. Źródło zostało nietknięte, więc próba nic nie kosztuje — rozpoznaj nową pozycję ' +
        'i porównaj wynik z poprzednią.',
      true,
    );
  }

  // ── Układ ──────────────────────────────────────────────────────────────────

  const przygotowanie = document.createElement('div');
  przygotowanie.className = 'ms-cyfryzacja__poprawki';
  przygotowanie.setAttribute('aria-label', 'Przygotowanie obrazu przed rozpoznaniem');
  for (const poprawka of POPRAWKI_OBRAZU) {
    const kontrolka = przycisk(poprawka.nazwa, 'dn-btn dn-btn--sm dn-btn--duch');
    kontrolka.title = poprawka.opis;
    kontrolka.setAttribute('aria-description', poprawka.opis);
    kontrolka.dataset['poprawka'] = poprawka.kod;
    kontrolka.addEventListener('click', () => {
      void przygotuj(poprawka.nazwa, (wskazanie) =>
        material.popraw(stan.idOkna(), wskazanie, poprawka.rodzaj, 0),
      );
    });
    przygotowanie.append(kontrolka);
  }

  const prostuj = przycisk('Prostuj skos obrotem', 'dn-btn dn-btn--sm dn-btn--duch');
  prostuj.dataset['poprawka'] = 'skos';
  prostuj.title =
    'Obrót o zadany kąt oddaje nowy zasób. Prostowanie skosu w ramach samego rozpoznania włącza ' +
    'się nastawą wyżej i nie zakłada osobnego zasobu.';
  prostuj.addEventListener('click', () => {
    void przygotuj('Prostowanie skosu', (wskazanie) =>
      material.obroc(stan.idOkna(), wskazanie, 0),
    );
  });
  przygotowanie.append(prostuj);

  const oPakiecieSerwera = document.createElement('p');
  oPakiecieSerwera.className = 'dn-pole-opis';
  oPakiecieSerwera.textContent = SKLADNIKI_PAKIETU_SERWERA.rozpoznanie;

  const oCzyszczeniu = document.createElement('p');
  oCzyszczeniu.className = 'dn-pole-opis';
  oCzyszczeniu.textContent = DROGI_CYFRYZACJI.czyszczenie;

  const oKorekcie = document.createElement('p');
  oKorekcie.className = 'dn-pole-opis';
  oKorekcie.textContent = NASTAWY_CYFRYZACJI.korekta;

  const oAdresie = document.createElement('p');
  oAdresie.className = 'dn-pole-opis';
  oAdresie.textContent = NASTAWY_CYFRYZACJI.adres;

  rama.pasek.append(cyfryzuj);
  rama.stan.tresc.append(
    ...pola.wsad,
    ...pola.nastawy,
    oCzyszczeniu,
    przygotowanie,
    tabela,
    bilans,
    oKorekcie,
    oObrazie,
    powierzchnie,
    ...pola.przyjecie,
    doEdytora,
    odpowiedz.element,
    oPakiecieSerwera,
    oAdresie,
  );

  /* ── Panel na żądanie, nie stała kolumna ─────────────────────────────────── */

  const nakladka = document.createElement('div');
  nakladka.className = 'ms-cyfryzacja__nakladka';
  nakladka.hidden = true;
  nakladka.append(rama.element);

  const wyzwalacz = przycisk('Narzędziownia cyfryzacji ▾', 'dn-btn dn-btn--sm dn-btn--zarys');
  wyzwalacz.dataset['czynnosc'] = 'narzedziownia';
  wyzwalacz.setAttribute('aria-expanded', 'false');
  wyzwalacz.title =
    'Otwiera narzędziownię cyfryzacji nakładką: kolejka wczytywania rdzenia, rozpoznanie tekstu ze ' +
    'sterowaniem silnikiem, językami i progiem pewności, poprawka rozpoznanych słów przed ' +
    'przyjęciem, wykaz urządzeń wejściowych i przyjęcie wyniku jako dokumentu. Schodzi drugim ' +
    'naciśnięciem albo Escapem.';
  wyzwalacz.addEventListener('click', () => przestawNarzedziownie(nakladka.hidden));

  const powloka = document.createElement('div');
  powloka.className = 'ms-cyfryzacja';
  powloka.append(wyzwalacz, nakladka);
  powloka.addEventListener('keydown', (zdarzenie) => {
    if (zdarzenie.key === 'Escape' && !nakladka.hidden) przestawNarzedziownie(false);
  });

  function przestawNarzedziownie(otwarta: boolean): void {
    nakladka.hidden = !otwarta;
    wyzwalacz.setAttribute('aria-expanded', otwarta ? 'true' : 'false');
    opiszWyzwalacz();
    if (!otwarta) return;
    const zapleczeCzynnosci = zaplecze();
    if (zapleczeCzynnosci === null) return;
    void odczytajKolejke(zapleczeCzynnosci);
  }

  function opiszWyzwalacz(): void {
    const liczby = kolejka.bilans();
    const znak = nakladka.hidden ? '▾' : '▴';
    wyzwalacz.textContent =
      liczby.wszystkie === 0
        ? `Narzędziownia cyfryzacji ${znak}`
        : `Narzędziownia cyfryzacji: ${liczby.wszystkie} w kolejce ${znak}`;
  }

  pola.doloz.addEventListener('click', () => {
    const zapleczeCzynnosci = zaplecze();
    if (zapleczeCzynnosci === null) return;
    void dolozMaterial(zapleczeCzynnosci, pola.wskazanie(), pola.czyZasob()).then(() =>
      pola.wyczyscWskazanie(),
    );
  });

  pola.rozpakuj.addEventListener('click', () => {
    const zapleczeCzynnosci = zaplecze();
    if (zapleczeCzynnosci === null) return;
    void dolozArchiwum(zapleczeCzynnosci, pola.archiwum());
  });

  pola.odczytajUrzadzenia.addEventListener('click', () => {
    const zapleczeCzynnosci = zaplecze();
    if (zapleczeCzynnosci === null) return;
    void odczytajUrzadzenia(zapleczeCzynnosci, (urzadzenia) => pola.ustawUrzadzenia(urzadzenia));
  });

  cyfryzuj.addEventListener('click', () => {
    const zapleczeCzynnosci = zaplecze();
    if (zapleczeCzynnosci === null) return;
    void rozpoznajKolejke(zapleczeCzynnosci);
  });

  doEdytora.addEventListener('click', () => {
    const zapleczeCzynnosci = zaplecze();
    if (zapleczeCzynnosci === null) return;
    const wskazana = kolejka.wskazana();
    void przyjmijWynik(zapleczeCzynnosci, wskazana === '' ? [] : [wskazana]);
  });

  // ── Odświeżenie ────────────────────────────────────────────────────────────

  function przerysujKolejke(): void {
    cialoKolejki.replaceChildren(
      ...kolejka.pozycje().map((pozycja) => {
        const wiersz = utworzWierszCyfryzacji(
          pozycja,
          pozycja.id === kolejka.wskazana(),
          kolejka.wToku(pozycja.id),
        );
        wiersz.rozpoznaj.addEventListener('click', () => {
          const zapleczeCzynnosci = zaplecze();
          if (zapleczeCzynnosci === null) return;
          void rozpoznajPozycje(zapleczeCzynnosci, pozycja.id);
        });
        wiersz.wskaz.addEventListener('click', () => {
          kolejka.ustawWskazana(pozycja.id);
          odswiez();
        });
        wiersz.przyjmij.addEventListener('click', () => {
          const zapleczeCzynnosci = zaplecze();
          if (zapleczeCzynnosci === null) return;
          void przyjmijWynik(zapleczeCzynnosci, [pozycja.id]);
        });
        return wiersz.element;
      }),
    );
  }

  /** Wykaz słów rozpoznanych wraz z polem poprawki przy każdym. */
  function przerysujSlowa(): void {
    const wskazana = kolejka.wskazana();
    const wykaz = kolejka.slowa(wskazana);
    if (wskazana === '') {
      const puste = document.createElement('li');
      puste.className = 'dn-pole-opis';
      puste.textContent =
        'Wskaż pozycję w kolejce, żeby zobaczyć jej rozpoznane słowa i warstwę tekstową.';
      slowa.replaceChildren(puste);
      return;
    }
    if (wykaz.length === 0) {
      const puste = document.createElement('li');
      puste.className = 'dn-pole-opis';
      puste.textContent =
        'Rdzeń nie oddał dla tej pozycji słów wraz z położeniem, więc poprawianie nie ma na czym ' +
        'pracować. Rozpoznaj pozycję ponownie — słowa przychodzą w odpowiedzi rozpoznania, a wykaz ' +
        'kolejki ich nie powtarza.';
      slowa.replaceChildren(puste);
      return;
    }
    slowa.replaceChildren(...wykaz.map((slowo) => wierszSlowa(wskazana, slowo)));
  }

  /** Jedno słowo: brzmienie, położenie, pewność i pole poprawki. */
  function wierszSlowa(idPozycji: string, slowo: StudioRecognizedWord): HTMLElement {
    const opis = document.createElement('p');
    opis.className = 'ms-cyfryzacja__slowo-opis';
    opis.textContent =
      `„${slowo.text}" · strona ${slowo.page} · odsunięcie ${slowo.x}, ${slowo.y} · pole ` +
      `${slowo.width} × ${slowo.height} punktów · pewność ${slowo.confidence}` +
      (slowo.corrected === true ? ' · POPRAWIONE przez Operatora' : '');

    const poprawka = document.createElement('input');
    poprawka.type = 'text';
    poprawka.className = 'dn-pole-kontrolka';
    poprawka.value = slowo.text;
    poprawka.setAttribute('aria-label', `Poprawka słowa o numerze ${slowo.index}`);

    const zapisz = przycisk('Popraw', 'dn-btn dn-btn--sm dn-btn--zarys');
    zapisz.dataset['slowo'] = String(slowo.index);
    zapisz.addEventListener('click', () => {
      const zapleczeCzynnosci = zaplecze();
      if (zapleczeCzynnosci === null) return;
      void poprawSlowo(zapleczeCzynnosci, idPozycji, slowo.index, poprawka.value.trim());
    });

    const pozycja = document.createElement('li');
    pozycja.dataset['slowo'] = String(slowo.index);
    pozycja.dataset['poprawione'] = String(slowo.corrected === true);
    // Pewność niżej niż połowa jest treścią, nie barwą: znacznik idzie do danych,
    // żeby wykaz dał się przejść także bez odczytu koloru.
    pozycja.dataset['niepewne'] = String(slowo.confidence < 0.5);
    pozycja.append(opis, poprawka, zapisz);
    return pozycja;
  }

  function odswiez(): void {
    przerysujKolejke();
    przerysujSlowa();
    const liczby = kolejka.bilans();
    bilans.textContent =
      `Kolejka rdzenia: pozycji ${liczby.wszystkie} · oczekuje ${liczby.oczekujace} · ` +
      `w toku ${liczby.przetwarzane} · gotowych ${liczby.gotowe} · do ponowienia ` +
      `${liczby.ponowienia} · odmów ${liczby.odmowy} · z tekstem gotowym do przyjęcia ` +
      `${liczby.zTekstem}.`;

    const pozycja = kolejka.pozycja(kolejka.wskazana());
    podglad.textContent = pozycja === null ? '' : pozycja.text ?? '';
    podglad.dataset['pozycja'] = pozycja === null ? '' : pozycja.id;
    opiszWyzwalacz();

    // Odczyt w toku i odmowa rdzenia są stanami trwałymi: odświeżenie wywołane
    // zmianą w innym oknie nie może ich zdjąć.
    if (rama.stan.faza() === 'ladowanie' || rama.stan.faza() === 'blad') return;
    if (liczby.wszystkie === 0) {
      rama.stan.puste('Kolejka wczytywania bez pozycji', BEZ_POZYCJI);
      return;
    }
    rama.stan.gotowe();
  }

  // Zamiana formatu dokumentu ma dziś własną drogę w rodzinie `studio.*`, więc
  // źródło obszaru `document` zostaje w umowie wytwórni dla zgodności montażu,
  // ale narzędziownia niczego nim już nie woła. Wskazanie tego wprost jest
  // uczciwsze niż milczące ignorowanie argumentu.
  void dokumenty;

  odswiez();

  return {
    element: powloka,
    odswiez,
    przestawWidocznosc: () => przestawNarzedziownie(nakladka.hidden),
  };
}

const BRAK_WSKAZANEJ =
  'Poprawka obrazu dotyczy pozycji wskazanej w kolejce — naciśnij „Wskaż" przy pozycji, którą ' +
  'chcesz przygotować przed rozpoznaniem.';

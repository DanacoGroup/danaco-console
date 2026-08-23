import { StudioDocumentFormat, StudioExportFormat } from '../../../../shared/contract';
import { opisOdmowyBledu } from '../../komponenty/odmowa';
import {
  poleTekstowe,
  przycisk,
  utworzWierszOdpowiedzi,
  wybor,
  type WierszOdpowiedzi,
} from '../../modele/kontrolki-formularza';
import type { StanStudio } from './stan-studio';
import type { ZrodloDokumentuStudio } from './zrodlo-dokumentu-studio';
import {
  wstawieniaOpiszPominiecia,
  wstawieniaOpiszWydanie,
  type ZrodloWstawienStudio,
} from './zrodlo-wstawien-studio';

/**
 * Wydanie dokumentu w formacie docelowym — jedna czynność dwóch okien.
 *
 * Wydania żąda i pasek narzędzi Studio Editora („Zapisz jako …"), i Preview
 * Window („Eksportuj"). Czynność jest jedna, więc stoi w jednym miejscu: dwie
 * kopie rozjechałyby się przy pierwszej poprawce zdania o wyniku.
 *
 * Przedmiotem zamiany jest treść zaakceptowana modułu, nie plik na dysku.
 * Dokument Studia mieszka w rdzeniu (`studio.document.save`), a nie w systemie
 * plików Operatora, więc `document.convert` dostaje treść wprost polem
 * `content`. Wynik jest zasobem magazynu rdzenia — komenda oddaje jego
 * identyfikator i rozmiar, a nie bajty do pobrania.
 */

/**
 * Formaty docelowe wymienione w opisie komendy `document.convert`.
 *
 * Wykaz pochodzi z kontraktu, nie z wyobrażenia okna: rdzeń odmawia formatu
 * spoza swojego słownika i wylicza w odmowie, co umie.
 */
export const FORMATY_KONWERSJI: readonly string[] = [
  'markdown',
  'html',
  'docx',
  'odt',
  'pdf',
  'epub',
  'rtf',
  'csv',
  'txt',
];

/**
 * Format źródłowy odczytany z dokumentu wskazanego przez rdzeń.
 *
 * `StudioDocument.format` niesie cztery wartości, a zamiana formatu przyjmuje
 * nazwy własnego słownika. Trzy z czterech mają w nim odpowiednik wprost; PDF
 * odpowiednika użytecznego nie ma, bo treść jedzie tu napisem, a nie plikiem —
 * napis nie jest PDF-em, choćby rdzeń tak nazywał dokument, z którego powstał.
 * Zwracamy wtedy tekst czysty i mówimy o tym Operatorowi zdaniem `POWOD_ZRODLA`.
 */
export function formatZrodlowy(format: StudioDocumentFormat): string {
  if (format === StudioDocumentFormat.Pdf) return 'txt';
  if (format === StudioDocumentFormat.Docx) return 'txt';
  return format;
}

/**
 * Zdanie o formacie źródłowym — mówi, czym naprawdę jest treść idąca do zamiany.
 *
 * Bez niego Operator widziałby „format z rdzenia: DOCX" i wydanie zamawiane
 * z formatu, którego rdzeń w tej drodze nie czyta; zdanie nazywa różnicę między
 * formatem DOKUMENTU a formatem TREŚCI, która idzie do zamiany.
 */
export const POWOD_ZRODLA =
  'Zamiana formatu dostaje treść dokumentu napisem, a nie plikiem. Dla dokumentu, który rdzeń ' +
  'prowadzi jako PDF albo DOCX, treścią jest sam tekst — format źródłowy zamiany jest więc ' +
  'tekstem czystym, a układ, style i osadzenia pliku wejściowego nie mają czego przenieść.';

/** Zaplecze wydania — okno podaje własne miejsca na komunikat. */
export interface ZapleczeKonwersji {
  stan: StanStudio;
  dokumenty: ZrodloDokumentuStudio;
  odpowiedz: WierszOdpowiedzi;
}

/**
 * Zamawia wydanie dokumentu w formacie docelowym i nazywa skutek.
 *
 * Zwraca `true`, gdy rdzeń oddał zasób — okno wołające dopisuje wtedy własne
 * zdanie o tym, co się z zasobem dzieje dalej.
 */
export async function wydajDokument(
  zaplecze: ZapleczeKonwersji,
  formatDocelowy: string,
): Promise<boolean> {
  const { stan, dokumenty, odpowiedz } = zaplecze;
  const dokument = stan.dokument();
  if (dokument === null) {
    odpowiedz.pokaz(BRAK_DOKUMENTU, false);
    return false;
  }
  const tresc = stan.trescZaakceptowana();
  if (tresc === '') {
    odpowiedz.pokaz(BRAK_TRESCI, false);
    return false;
  }
  odpowiedz.pokaz(`Zamiana treści dokumentu na format ${formatDocelowy} w toku…`, true);
  const wynik = await dokumenty.zamienFormat(
    stan.idOkna(),
    tresc,
    formatZrodlowy(dokument.format),
    formatDocelowy,
  );
  if (!wynik.udany || wynik.wynik === undefined) {
    odpowiedz.pokaz(opisOdmowyBledu('Wydanie dokumentu', wynik.blad), false);
    return false;
  }
  const zasob = wynik.wynik.asset;
  odpowiedz.pokaz(
    `Dokument wydany w formacie ${formatDocelowy}: zasób ${zasob.id} w magazynie rdzenia ` +
      `(${wynik.wynik.sizeBytes} bajtów). Zasób zostaje po stronie rdzenia — komendy pobierającej ` +
      'jego bajty do przeglądarki kontrakt nie niesie.',
    true,
  );
  return true;
}

const BRAK_DOKUMENTU =
  'Wydanie dotyczy dokumentu — wczytaj go najpierw w Studio Editorze; komenda zamiany formatu ' +
  'bierze jego treść, a nie ścieżkę pliku.';

const BRAK_TRESCI =
  'Dokument nie ma jeszcze treści zaakceptowanej, więc nie ma czego wydać. Zapisz go w Studio ' +
  'Editorze albo przyjmij wynik operacji — dopiero wtedy treść jest ustalona.';

/* ── Wydanie dokumentu rodziną `studio.document.*` ──────────────────────────── */

/**
 * Wydanie dokumentu — zapis pod nazwą, wydanie do formatu, wydanie wsadowe
 * i format dokumentu.
 *
 * ── Dlaczego to stoi obok drogi wyżej, a nie zamiast niej ───────────────────
 * Droga wyżej (`document.convert`) bierze treść NAPISEM i nie ma czym przenieść
 * postaci: dla dokumentu prowadzonego jako PDF albo DOCX treścią jest sam tekst.
 * Cztery komendy niżej pracują na DOKUMENCIE po stronie rdzenia, więc przenoszą
 * arkusz stylów, sekcje, tabele, obrazy i aparat:
 *   — `document.save.as` zapisuje pod nową nazwą albo do wskazanego pliku,
 *   — `document.export.format` wydaje do txt, md, docx, odt, pdf, html i rtf,
 *   — `document.export.batch` wydaje wiele dokumentów naraz,
 *   — `document.format.set` przestawia format dokumentu Studia.
 * Starsza droga zostaje dla treści, która nie jest dokumentem rdzenia.
 *
 * ── Bilans cech pominiętych jest widoczny zawsze ────────────────────────────
 * Format uboższy niż dokument jest normalną sytuacją; przemilczenie straty nie
 * jest. Wydanie do tekstu czystego, które zgubiło tabelę i przypisy, mówi to
 * wprost — i nie ma w tym panelu drogi, którą Operator zobaczyłby samo
 * „zapisano". Wydanie wsadowe wypisuje bilans dokument po dokumencie, a odmowa
 * jednego nie ukrywa się za liczbą wydanych.
 */
export interface WydaniePanel {
  element: HTMLElement;
  przestawWidocznosc(): void;
}

/** Siedem formatów wydania z kontraktu wraz z tym, co każdy niesie. */
const FORMATY_WYDANIA: readonly (readonly [string, string])[] = [
  [StudioExportFormat.Txt, 'tekst czysty — postać schodzi'],
  [StudioExportFormat.Md, 'markdown — style nazwane na znaczniki, tabele na tabele markdown'],
  [StudioExportFormat.Docx, 'dokument Word (docx) — postać zachowana'],
  [StudioExportFormat.Odt, 'dokument OpenDocument (odt) — postać zachowana'],
  [StudioExportFormat.Pdf, 'PDF — przez profil wydania, z paginacją i stopką'],
  [StudioExportFormat.Html, 'HTML wraz z arkuszem stylów'],
  [StudioExportFormat.Rtf, 'RTF'],
];

export function utworzWydaniePanel(
  stan: StanStudio,
  zrodlo: ZrodloWstawienStudio,
): WydaniePanel {
  const odpowiedz = utworzWierszOdpowiedzi();

  /* ── Zapis pod nazwą ───────────────────────────────────────────────────── */

  const nowaNazwa = poleTekstowe({
    etykieta: 'Nowa nazwa dokumentu',
    podpowiedz: 'puste zostawia nazwę bieżącą',
  });
  const sciezkaZapisu = poleTekstowe({
    etykieta: 'Ścieżka pliku, do którego zapisać',
    podpowiedz: 'puste zostawia dokument po stronie rdzenia',
  });
  const zWersja = document.createElement('input');
  zWersja.type = 'checkbox';
  zWersja.className = 'dn-przelacznik';
  zWersja.checked = true;
  zWersja.setAttribute('aria-label', 'Załóż wersję w repozytorium sesji');
  const etykietaWersji = document.createElement('label');
  etykietaWersji.className = 'ms-wydanie__zawezenie';
  etykietaWersji.append(zWersja, document.createTextNode('załóż wersję w repozytorium sesji'));

  const zapiszJako = przycisk('Zapisz jako', 'dn-btn dn-btn--sm dn-btn--sygnal');
  zapiszJako.dataset['czynnosc'] = 'zapisz-jako';
  zapiszJako.title =
    'studio.document.save.as zapisuje dokument WRAZ Z POSTACIĄ: arkuszem stylów, nastawami strony, ' +
    'sekcjami, tabelami, obrazami i aparatem.';
  zapiszJako.addEventListener('click', () => void zapiszPodNazwa());

  /* ── Kopia dokumentu ───────────────────────────────────────────────────── */

  const nazwaKopii = poleTekstowe({ etykieta: 'Nazwa kopii', podpowiedz: '' });
  const zHistoria = document.createElement('input');
  zHistoria.type = 'checkbox';
  zHistoria.className = 'dn-przelacznik';
  zHistoria.setAttribute('aria-label', 'Przenieś historię wersji do kopii');
  const etykietaHistorii = document.createElement('label');
  etykietaHistorii.className = 'ms-wydanie__zawezenie';
  etykietaHistorii.append(
    zHistoria,
    document.createTextNode('przenieś historię wersji — wybór jawny Operatora'),
  );

  const zeZnakowaniem = document.createElement('input');
  zeZnakowaniem.type = 'checkbox';
  zeZnakowaniem.className = 'dn-przelacznik';
  zeZnakowaniem.setAttribute('aria-label', 'Przenieś znakowanie i komentarze do kopii');
  const etykietaZnakowania = document.createElement('label');
  etykietaZnakowania.className = 'ms-wydanie__zawezenie';
  etykietaZnakowania.append(
    zeZnakowaniem,
    document.createTextNode('przenieś znakowanie i komentarze'),
  );

  const zBlokadami = document.createElement('input');
  zBlokadami.type = 'checkbox';
  zBlokadami.className = 'dn-przelacznik';
  zBlokadami.checked = true;
  zBlokadami.setAttribute('aria-label', 'Przenieś blokady fragmentów do kopii');
  const etykietaBlokad = document.createElement('label');
  etykietaBlokad.className = 'ms-wydanie__zawezenie';
  etykietaBlokad.append(zBlokadami, document.createTextNode('przenieś blokady fragmentów'));

  const skopiuj = przycisk('Załóż kopię dokumentu', 'dn-btn dn-btn--sm dn-btn--atrament');
  skopiuj.dataset['czynnosc'] = 'kopia-dokumentu';
  skopiuj.title =
    'Kopia jest OSOBNYM dokumentem, nie drugim odwołaniem do tego samego: zmiana w kopii nie rusza ' +
    'oryginału.';
  skopiuj.addEventListener('click', () => void zalozKopie());

  /* ── Wydanie do formatu ────────────────────────────────────────────────── */

  const formatWydania = wybor('Format wydania', FORMATY_WYDANIA.map((pozycja) => pozycja));
  const sciezkaWydania = poleTekstowe({
    etykieta: 'Ścieżka pliku wydania',
    podpowiedz: 'puste zostawia wynik zasobem magazynu rdzenia',
  });
  const profil = poleTekstowe({
    etykieta: 'Profil wydania',
    podpowiedz: 'obowiązkowy przy PDF — niesie paginację i stopkę',
  });
  const zeSzablonu = poleTekstowe({
    etykieta: 'Szablon, którego postać narzucić wydaniu',
    podpowiedz: 'puste zostawia postać dokumentu',
  });
  const zKomentarzami = document.createElement('input');
  zKomentarzami.type = 'checkbox';
  zKomentarzami.className = 'dn-przelacznik';
  zKomentarzami.setAttribute('aria-label', 'Wydaj komentarze');
  const etykietaKomentarzy = document.createElement('label');
  etykietaKomentarzy.className = 'ms-wydanie__zawezenie';
  etykietaKomentarzy.append(zKomentarzami, document.createTextNode('wydaj komentarze'));

  const zeZmianami = document.createElement('input');
  zeZmianami.type = 'checkbox';
  zeZmianami.className = 'dn-przelacznik';
  zeZmianami.setAttribute('aria-label', 'Wydaj zmiany śledzone');
  const etykietaZmian = document.createElement('label');
  etykietaZmian.className = 'ms-wydanie__zawezenie';
  etykietaZmian.append(
    zeZmianami,
    document.createTextNode('wydaj zmiany śledzone — pismo do wysłania idzie BEZ adiustacji'),
  );

  const wydaj = przycisk('Wydaj dokument do formatu', 'dn-btn dn-btn--sm dn-btn--sygnal');
  wydaj.dataset['czynnosc'] = 'wydaj-format';
  wydaj.addEventListener('click', () => void wydajDoFormatu());

  /* ── Wydanie wsadowe ───────────────────────────────────────────────────── */

  const dokumentyWsadu = poleTekstowe({
    etykieta: 'Dokumenty wydawane',
    podpowiedz: 'identyfikatory rozdzielone przecinkiem',
  });
  const katalogWsadu = poleTekstowe({
    etykieta: 'Katalog, do którego zapisać',
    podpowiedz: 'ścieżka widziana przez rdzeń',
  });
  const wydajWsad = przycisk('Wydaj wsad dokumentów', 'dn-btn dn-btn--sm dn-btn--atrament');
  wydajWsad.dataset['czynnosc'] = 'wydaj-wsad';
  wydajWsad.title =
    'Odmowa jednego dokumentu nie wstrzymuje pozostałych, a każdy wynik niesie swój bilans cech ' +
    'pominiętych.';
  wydajWsad.addEventListener('click', () => void wydajWsadDokumentow());

  /* ── Format dokumentu ──────────────────────────────────────────────────── */

  const formatDokumentu = wybor(
    'Format dokumentu Studia',
    Object.values(StudioDocumentFormat).map((format) => [format, format] as const),
  );
  const przestawFormat = przycisk('Przestaw format dokumentu', 'dn-btn dn-btn--sm dn-btn--zarys');
  przestawFormat.dataset['czynnosc'] = 'format-dokumentu';
  przestawFormat.title =
    'Format nadaje rdzeń przy wczytaniu, a studio.document.save go nie zmienia — ta komenda jest ' +
    'jedyną drogą jego przestawienia.';
  przestawFormat.addEventListener('click', () => void przestawFormatDokumentu());

  /* ── Czynności ─────────────────────────────────────────────────────────── */

  function idDokumentu(): string | null {
    const dokument = stan.dokument();
    if (dokument === null) {
      odpowiedz.pokaz(BRAK_DOKUMENTU_WYDANIA, false);
      return null;
    }
    return dokument.id;
  }

  async function zapiszPodNazwa(): Promise<void> {
    const dokument = idDokumentu();
    if (dokument === null) return;
    const wynik = await zrodlo.zapiszJako({
      documentId: dokument,
      ...(nowaNazwa.kontrolka.value.trim() === ''
        ? {}
        : { title: nowaNazwa.kontrolka.value.trim() }),
      ...(sciezkaZapisu.kontrolka.value.trim() === ''
        ? {}
        : { path: sciezkaZapisu.kontrolka.value.trim() }),
      createVersion: zWersja.checked,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowyBledu('Zapis pod nową nazwą', wynik.blad), false);
      return;
    }
    stan.wchlon(wynik.wynik.document);
    const wydanie = wynik.wynik.export;
    odpowiedz.pokaz(
      `Dokument zapisany jako ${wynik.wynik.document.id} („${
        wynik.wynik.document.title ?? 'bez tytułu'
      }") wraz z całą postacią. ` +
        (wydanie === undefined
          ? 'Do pliku nie zapisywano — dokument został po stronie rdzenia.'
          : wstawieniaOpiszWydanie(wydanie)),
      true,
    );
  }

  async function zalozKopie(): Promise<void> {
    const dokument = idDokumentu();
    if (dokument === null) return;
    const wynik = await zrodlo.skopiuj({
      documentId: dokument,
      ...(nazwaKopii.kontrolka.value.trim() === ''
        ? {}
        : { title: nazwaKopii.kontrolka.value.trim() }),
      ...(stan.idOkna() === '' ? {} : { windowId: stan.idOkna() }),
      includeVersions: zHistoria.checked,
      includeMarkup: zeZnakowaniem.checked,
      includeLocks: zBlokadami.checked,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowyBledu('Założenie kopii dokumentu', wynik.blad), false);
      return;
    }
    odpowiedz.pokaz(
      `Kopia założona jako OSOBNY dokument ${wynik.wynik.document.id} — zmiana w niej nie rusza ` +
        `oryginału ${dokument}. Wersji przeniesionych: ${wynik.wynik.versionsCopied}` +
        (zHistoria.checked
          ? '.'
          : ' (historia nie była przenoszona — kopia zaczyna swój szereg wersji od nowa).'),
      true,
    );
  }

  async function wydajDoFormatu(): Promise<void> {
    const dokument = idDokumentu();
    if (dokument === null) return;
    const format = formatWydania.value as StudioExportFormat;
    if (format === StudioExportFormat.Pdf && profil.kontrolka.value.trim() === '') {
      // Odmowa nazwana PRZED próbą: profil wydania niesie paginację i stopkę, więc
      // PDF bez niego nie jest wydaniem, tylko treścią wsypaną na stronę.
      odpowiedz.pokaz(BRAK_PROFILU, false);
      return;
    }
    const wynik = await zrodlo.wydaj({
      documentId: dokument,
      format,
      ...(sciezkaWydania.kontrolka.value.trim() === ''
        ? {}
        : { path: sciezkaWydania.kontrolka.value.trim() }),
      ...(profil.kontrolka.value.trim() === ''
        ? {}
        : { profileId: profil.kontrolka.value.trim() }),
      ...(zeSzablonu.kontrolka.value.trim() === ''
        ? {}
        : { templateId: zeSzablonu.kontrolka.value.trim() }),
      includeComments: zKomentarzami.checked,
      includeTrackedChanges: zeZmianami.checked,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowyBledu('Wydanie dokumentu do formatu', wynik.blad), false);
      return;
    }
    odpowiedz.pokaz(wstawieniaOpiszWydanie(wynik.wynik.result), true);
  }

  async function wydajWsadDokumentow(): Promise<void> {
    const identyfikatory = dokumentyWsadu.kontrolka.value
      .split(',')
      .map((czesc) => czesc.trim())
      .filter((czesc) => czesc !== '');
    if (identyfikatory.length === 0) {
      odpowiedz.pokaz(BRAK_WSADU, false);
      return;
    }
    const wynik = await zrodlo.wydajWsad({
      documentIds: identyfikatory,
      format: formatWydania.value as StudioExportFormat,
      ...(katalogWsadu.kontrolka.value.trim() === ''
        ? {}
        : { directory: katalogWsadu.kontrolka.value.trim() }),
      ...(profil.kontrolka.value.trim() === ''
        ? {}
        : { profileId: profil.kontrolka.value.trim() }),
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowyBledu('Wydanie wsadowe', wynik.blad), false);
      return;
    }
    const odmowy = wynik.wynik.failures ?? [];
    odpowiedz.pokaz(
      `Wsad wydany: dokumentów ${wynik.wynik.succeeded}, odmów ${wynik.wynik.failed}. ` +
        // Bilans dokument po dokumencie, nie jedna liczba: format uboższy niż
        // dokument gubi cechy w KAŻDYM z nich osobno.
        wynik.wynik.results.map((pozycja) => wstawieniaOpiszWydanie(pozycja)).join(' ') +
        (odmowy.length === 0
          ? ''
          : ` Odmowy: ${wstawieniaOpiszPominiecia(odmowy)}`),
      wynik.wynik.failed === 0,
    );
  }

  async function przestawFormatDokumentu(): Promise<void> {
    const dokument = idDokumentu();
    if (dokument === null) return;
    const wynik = await zrodlo.przestawFormat({
      documentId: dokument,
      format: formatDokumentu.value as StudioDocumentFormat,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowyBledu('Przestawienie formatu dokumentu', wynik.blad), false);
      return;
    }
    stan.wchlon(wynik.wynik.document);
    odpowiedz.pokaz(
      `Format dokumentu ${wynik.wynik.document.id} przestawiony na ${wynik.wynik.document.format}.`,
      true,
    );
  }

  /* ── Nakładka ──────────────────────────────────────────────────────────── */

  const tresc = document.createElement('div');
  tresc.className = 'ms-wydanie__tresc';
  tresc.append(
    czescWydania('Zapis pod nazwą i do pliku', [
      nowaNazwa.element,
      sciezkaZapisu.element,
      etykietaWersji,
      zapiszJako,
    ]),
    czescWydania('Kopia dokumentu', [
      nazwaKopii.element,
      etykietaHistorii,
      etykietaZnakowania,
      etykietaBlokad,
      skopiuj,
    ]),
    czescWydania('Wydanie do formatu', [
      formatWydania,
      sciezkaWydania.element,
      profil.element,
      zeSzablonu.element,
      etykietaKomentarzy,
      etykietaZmian,
      wydaj,
    ]),
    czescWydania('Wydanie wsadowe', [dokumentyWsadu.element, katalogWsadu.element, wydajWsad]),
    czescWydania('Format dokumentu Studia', [formatDokumentu, przestawFormat]),
    odpowiedz.element,
  );

  const nakladka = document.createElement('div');
  nakladka.className = 'ms-wydanie__nakladka';
  nakladka.hidden = true;
  nakladka.append(tresc);

  const wyzwalacz = przycisk('Zapis i wydanie ▾', 'dn-btn dn-btn--sm dn-btn--zarys');
  wyzwalacz.dataset['czynnosc'] = 'warsztat-wydania';
  wyzwalacz.setAttribute('aria-expanded', 'false');
  wyzwalacz.title =
    'Otwiera zapis i wydanie nakładką: zapis pod nową nazwą i do pliku wraz z całą postacią, kopia ' +
    'dokumentu, wydanie do txt, md, docx, odt, pdf, html i rtf wraz z wykazem cech pominiętych, ' +
    'wydanie wsadowe oraz format dokumentu. Schodzi drugim naciśnięciem albo Escapem.';
  wyzwalacz.addEventListener('click', () => przestaw(nakladka.hidden));

  const element = document.createElement('section');
  element.className = 'ms-wydanie';
  element.setAttribute('aria-label', 'Zapis, kopia i wydanie dokumentu');
  element.append(wyzwalacz, nakladka);
  element.addEventListener('keydown', (zdarzenie) => {
    if (zdarzenie.key === 'Escape' && !nakladka.hidden) przestaw(false);
  });

  function przestaw(otwarta: boolean): void {
    nakladka.hidden = !otwarta;
    wyzwalacz.setAttribute('aria-expanded', otwarta ? 'true' : 'false');
    wyzwalacz.textContent = `Zapis i wydanie ${otwarta ? '▴' : '▾'}`;
  }

  return { element, przestawWidocznosc: () => przestaw(nakladka.hidden) };
}

function czescWydania(tytul: string, elementy: readonly HTMLElement[]): HTMLElement {
  const naglowek = document.createElement('p');
  naglowek.className = 'ms-wydanie__tytul';
  naglowek.textContent = tytul;
  const sekcja = document.createElement('section');
  sekcja.className = 'ms-wydanie__czesc';
  sekcja.append(naglowek, ...elementy);
  return sekcja;
}

const BRAK_DOKUMENTU_WYDANIA =
  'Wydanie i zapis dotyczą dokumentu — wczytaj go albo załóż nowy. Komendy studio.document.* ' +
  'przyjmują identyfikator dokumentu jako pole obowiązkowe.';

const BRAK_PROFILU =
  'Wydanie do PDF idzie PRZEZ PROFIL WYDANIA — on niesie paginację, nagłówek i stopkę. Bez profilu ' +
  'PDF byłby treścią wsypaną na stronę bez numeracji. Zapisz profil na zakładce projektowania ' +
  'wstążki i wskaż go tutaj.';

const BRAK_WSADU =
  'Wydanie wsadowe potrzebuje identyfikatorów dokumentów rozdzielonych przecinkiem — bez nich nie ' +
  'ma czego wydać.';

import {
  StudioDocumentFormat,
  StudioImportFormat,
  StudioObjectSource,
  type StudioDocumentOpenRequest,
} from '../../../../shared/contract';
import { utworzDymekObjasnienia } from '../../komponenty/dymek';
import { opisOdmowyBledu } from '../../komponenty/odmowa';
import {
  poleLiczbowe,
  poleTekstowe,
  poleWyboru,
  przycisk,
  utworzWierszOdpowiedzi,
  wybor,
} from '../../modele/kontrolki-formularza';
import type { StanStudio } from './stan-studio';
import {
  wstawieniaOpiszBilans,
  wstawieniaOpiszWniesienie,
  type ZrodloWstawienStudio,
} from './zrodlo-wstawien-studio';

/**
 * Wskazanie dokumentu do wczytania — z Library, z sesji albo z urządzenia.
 *
 * `studio.document.open` przyjmuje trzy wzajemnie wykluczające się pola:
 * `documentId` (dokument otwarty wcześniej w tej sesji), `libraryFileId` (plik
 * repozytorium) i `path` (ścieżka na urządzeniu). Formularz pyta o źródło wprost,
 * zamiast zgadywać, bo pomyłka kończyłaby się odmową rdzenia o powodzie trudnym
 * do odczytania.
 *
 * Wybierak źródła nie jest wykazem plików: wykaz zasobów Library należy do okna
 * Library Explorer i przychodzi komendą `library.file.list`, a to pole przyjmuje
 * sam identyfikator.
 *
 * Wykaz formatów pochodzi z kontraktu (`StudioDocumentFormat`), nie z listy
 * zapisanej w widoku; format wybiera rdzeń przy wczytaniu, okno go tylko pokazuje.
 */
export interface WczytanieDokumentu {
  element: HTMLElement;
  /** Przycisk „Wczytaj dokument" — okno wiodące podpina do niego swoją czynność. */
  wczytaj: HTMLButtonElement;
  /** Treść żądania złożona z pól formularza; `null`, gdy pole wskazania jest puste. */
  zadanie(idOkna: string): StudioDocumentOpenRequest | null;
  /** Zdanie mówiące, czego brakuje — do pokazania po nieudanym złożeniu żądania. */
  brak(): string;
}

/** Trzy źródła dokumentu odpowiadające trzem polom żądania kontraktu. */
const ZRODLA = [
  { wartosc: 'library', etykieta: 'Plik z Library (libraryFileId)' },
  { wartosc: 'sesja', etykieta: 'Dokument otwarty w tej sesji (documentId)' },
  { wartosc: 'urzadzenie', etykieta: 'Ścieżka na urządzeniu (path)' },
];

export function utworzWczytanieDokumentu(): WczytanieDokumentu {
  const zrodlo = poleWyboru(
    {
      etykieta: 'Źródło dokumentu',
      opis: 'Kontrakt przyjmuje dokładnie jedno z trzech wskazań; wybór ustala, które pole poleci do rdzenia.',
    },
    ZRODLA,
  );
  zrodlo.element.append(
    utworzDymekObjasnienia(
      'Wykaz plików repozytorium należy do okna Library Explorer i przychodzi komendą library.file.list. ' +
        'Studio Editor nie buduje jego kopii — przyjmuje wskazanie.',
      { powloka: 'ms-dymek', znak: 'ms-dymek__znak' },
    ),
  );

  const wskazanie = poleTekstowe({
    etykieta: 'Wskazanie',
    podpowiedz: 'identyfikator albo ścieżka',
    opis: `Formaty przyjmowane przez rdzeń: ${Object.values(StudioDocumentFormat).join(' · ')}.`,
  });

  const wczytaj = przycisk('Wczytaj dokument', 'dn-btn dn-btn--sm dn-btn--atrament');
  wczytaj.dataset['czynnosc'] = 'wczytaj';

  const element = document.createElement('div');
  element.className = 'ms-wczytanie';
  element.append(zrodlo.element, wskazanie.element, wczytaj);

  return {
    element,
    wczytaj,

    zadanie(idOkna) {
      const wartosc = wskazanie.kontrolka.value.trim();
      if (wartosc === '') return null;
      const zadanie: StudioDocumentOpenRequest = { windowId: idOkna };
      if (zrodlo.kontrolka.value === 'library') zadanie.libraryFileId = wartosc;
      else if (zrodlo.kontrolka.value === 'sesja') zadanie.documentId = wartosc;
      else zadanie.path = wartosc;
      return zadanie;
    },

    brak: () =>
      'Wskaż dokument — rdzeń nie wczyta niczego bez identyfikatora pliku, dokumentu albo ścieżki.',
  };
}

/* ── Wniesienie pliku WPROST DO EDYTORA ─────────────────────────────────────── */

/**
 * Wniesienie pliku, PDF-u i obrazu wprost do edytora — nakładka na żądanie.
 *
 * ── Czym to się różni od wczytania wyżej ────────────────────────────────────
 * `studio.document.open` otwiera dokument, który rdzeń już prowadzi — pozycję
 * repozytorium sesji albo pliku Library. Nie wnosi POSTACI pliku Operatora:
 * arkusza stylów, sekcji, tabel ani obrazów. Trzy komendy niżej to robią
 * i dlatego stoją osobno:
 *   — `document.import.file` wnosi docx, dotx, odt, ott, tekst czysty, markdown,
 *     RTF i HTML wraz z rozpoznaniem zapisu znaków (także stron kodowych innych
 *     niż UTF-8, bo pliki Operatora bywają starsze);
 *   — `document.import.pdf` odzyskuje z PDF-u tekst, akapity, tabele i obrazy —
 *     na tyle, na ile PDF je niesie;
 *   — `document.image.import` wnosi obraz wprost w miejsce kursora.
 *
 * ── Bilans jest obowiązkowy, nie ozdobny ────────────────────────────────────
 * Odzyskanie z PDF-u jest ODTWORZENIEM, nie odczytem: PDF nie niesie struktury
 * akapitu ani tabeli wprost. Panel wypisuje więc bilans za każdym razem — strony
 * z warstwą tekstową i bez niej, tabele rozpoznane i nierozpoznane, obrazy
 * osadzone i pominięte. PDF ze samych skanów kieruje na rozpoznanie tekstu
 * i panel mówi to wprost, zamiast oddać pustą kartkę jako gotowy dokument.
 */
export interface WniesieniePliku {
  element: HTMLElement;
  przestawWidocznosc(): void;
}

/** Dziewięć formatów wnoszonych z kontraktu; PDF ma własną drogę niżej. */
const FORMATY_WNOSZONE: readonly (readonly [string, string])[] = [
  ['', 'rozpoznanie po zawartości'],
  [StudioImportFormat.Docx, 'dokument Word (docx)'],
  [StudioImportFormat.Dotx, 'szablon Word (dotx)'],
  [StudioImportFormat.Odt, 'dokument OpenDocument (odt)'],
  [StudioImportFormat.Ott, 'szablon OpenDocument (ott)'],
  [StudioImportFormat.Txt, 'tekst czysty'],
  [StudioImportFormat.Md, 'markdown'],
  [StudioImportFormat.Rtf, 'RTF'],
  [StudioImportFormat.Html, 'HTML'],
];

/** Skąd obraz pochodzi — pięć dróg, których obraz naprawdę używa. */
const ZRODLA_OBRAZU: readonly (readonly [string, string])[] = [
  [StudioObjectSource.File, 'plik wskazany przez Operatora'],
  [StudioObjectSource.CoreAsset, 'magazyn zasobów rdzenia'],
  [StudioObjectSource.DesignModule, 'moduł Design'],
  [StudioObjectSource.PhotoBank, 'baza zdjęciowa'],
  [StudioObjectSource.LibraryFile, 'plik Biblioteki Library'],
];

export function utworzWniesieniePliku(
  stan: StanStudio,
  zrodlo: ZrodloWstawienStudio,
): WniesieniePliku {
  const odpowiedz = utworzWierszOdpowiedzi();

  const sciezka = poleTekstowe({
    etykieta: 'Ścieżka pliku widziana przez rdzeń',
    podpowiedz: 'np. /dane/pismo.docx',
    opis:
      'Klient dysku nie czyta i nie zapisuje — podaje wskazanie, a plik wciąga rdzeń. Plik ' +
      'Biblioteki i zasób magazynu wskazuje się polami niżej.',
  });
  const plikBiblioteki = poleTekstowe({
    etykieta: 'Plik Biblioteki Library',
    podpowiedz: 'identyfikator pliku',
  });
  const zasob = poleTekstowe({
    etykieta: 'Zasób magazynu rdzenia',
    podpowiedz: 'identyfikator zasobu',
  });
  const format = wybor('Format pliku wnoszonego', FORMATY_WNOSZONE.map((pozycja) => pozycja));
  const zapisZnakow = poleTekstowe({
    etykieta: 'Zapis znaków, gdy Operator go zna',
    podpowiedz: 'np. windows-1250 — puste zostawia rozpoznanie rdzeniowi',
    opis:
      'Pliki Operatora bywają starsze niż UTF-8. Rdzeń rozpoznaje zapis sam i mówi w bilansie, ' +
      'co rozpoznał; wskazanie jawne jest dla przypadków, w których się pomylił.',
  });
  const doDokumentu = poleTekstowe({
    etykieta: 'Dokument, do którego wnieść treść',
    podpowiedz: 'puste zakłada nowy dokument',
  });
  const wDokumencie = document.createElement('input');
  wDokumencie.type = 'checkbox';
  wDokumencie.className = 'dn-przelacznik';
  wDokumencie.setAttribute('aria-label', 'Wnieś w miejsce kursora, a nie jako cały dokument');
  const etykietaMiejsca = document.createElement('label');
  etykietaMiejsca.className = 'ms-wniesienie__zawezenie';
  etykietaMiejsca.append(
    wDokumencie,
    document.createTextNode('wnieś w miejsce kursora, a nie jako cały dokument'),
  );

  const wniesPlik = przycisk('Wnieś plik do edytora', 'dn-btn dn-btn--sm dn-btn--sygnal');
  wniesPlik.dataset['czynnosc'] = 'wnies-plik';
  wniesPlik.addEventListener('click', () => void wniesPlikDoEdytora());

  /* ── PDF ────────────────────────────────────────────────────────────────── */

  const zakresStron = poleTekstowe({
    etykieta: 'Zakres stron PDF',
    podpowiedz: 'np. 1-5 — puste znaczy cały plik',
  });
  const zTabelami = document.createElement('input');
  zTabelami.type = 'checkbox';
  zTabelami.className = 'dn-przelacznik';
  zTabelami.checked = true;
  zTabelami.setAttribute('aria-label', 'Próbuj odtworzyć tabele');
  const etykietaTabel = document.createElement('label');
  etykietaTabel.className = 'ms-wniesienie__zawezenie';
  etykietaTabel.append(zTabelami, document.createTextNode('próbuj odtworzyć tabele'));

  const zObrazami = document.createElement('input');
  zObrazami.type = 'checkbox';
  zObrazami.className = 'dn-przelacznik';
  zObrazami.checked = true;
  zObrazami.setAttribute('aria-label', 'Osadzaj obrazy');
  const etykietaObrazow = document.createElement('label');
  etykietaObrazow.className = 'ms-wniesienie__zawezenie';
  etykietaObrazow.append(zObrazami, document.createTextNode('osadzaj obrazy'));

  const oOdzyskaniu = document.createElement('p');
  oOdzyskaniu.className = 'dn-pole-opis';
  oOdzyskaniu.textContent =
    'PDF nie niesie struktury akapitu ani tabeli wprost, więc odzyskanie jest ODTWORZENIEM, nie ' +
    'odczytem. Odpowiedź niesie bilans: strony z warstwą tekstową i bez niej, tabele rozpoznane ' +
    'i nierozpoznane, obrazy osadzone i pominięte. PDF ze samych skanów rdzeń kieruje na ' +
    'rozpoznanie tekstu, zamiast oddać pustą kartkę.';

  const wniesPdf = przycisk('Zamień PDF na dokument edytowalny', 'dn-btn dn-btn--sm dn-btn--sygnal');
  wniesPdf.dataset['czynnosc'] = 'wnies-pdf';
  wniesPdf.addEventListener('click', () => void wniesPdfDoEdytora());

  /* ── Obraz ──────────────────────────────────────────────────────────────── */

  const zrodloObrazu = wybor('Skąd obraz pochodzi', ZRODLA_OBRAZU.map((pozycja) => pozycja));
  const wskazanieObrazu = poleTekstowe({
    etykieta: 'Wskazanie obrazu',
    podpowiedz: 'ścieżka, zasób, węzeł Designu, pozycja bazy zdjęciowej albo plik Biblioteki',
  });
  const szerokoscObrazu = poleLiczbowe('Szerokość w milimetrach', '');
  const wysokoscObrazu = poleLiczbowe('Wysokość w milimetrach', '');
  const tekstZastepczy = poleTekstowe({
    etykieta: 'Tekst zastępczy',
    podpowiedz: 'co obraz przedstawia',
  });
  const podpisObrazu = poleTekstowe({ etykieta: 'Podpis pod obrazem', podpowiedz: '' });

  const wniesObraz = przycisk('Wnieś obraz w miejsce kursora', 'dn-btn dn-btn--sm dn-btn--sygnal');
  wniesObraz.dataset['czynnosc'] = 'wnies-obraz';
  wniesObraz.addEventListener('click', () => void wniesObrazDoEdytora());

  /* ── Czynności ─────────────────────────────────────────────────────────── */

  function idOkna(): string | null {
    if (stan.idOkna() === '') {
      odpowiedz.pokaz(BRAK_OKNA, false);
      return null;
    }
    return stan.idOkna();
  }

  function miejsceKursora(): number {
    const zakres = stan.zaznaczenie();
    return zakres === null ? stan.trescRobocza().length : zakres.koniec;
  }

  /** Wskazanie pliku w polu, które rdzeń dla niego czyta; `null` znaczy brak. */
  function wskazaniePliku(): Record<string, string> | null {
    if (sciezka.kontrolka.value.trim() !== '') return { path: sciezka.kontrolka.value.trim() };
    if (plikBiblioteki.kontrolka.value.trim() !== '') {
      return { libraryFileId: plikBiblioteki.kontrolka.value.trim() };
    }
    if (zasob.kontrolka.value.trim() !== '') return { assetId: zasob.kontrolka.value.trim() };
    odpowiedz.pokaz(BRAK_WSKAZANIA, false);
    return null;
  }

  async function wniesPlikDoEdytora(): Promise<void> {
    const okno = idOkna();
    const wskazanie = wskazaniePliku();
    if (okno === null || wskazanie === null) return;
    const dokument = doDokumentu.kontrolka.value.trim();
    const wynik = await zrodlo.wniesPlik({
      windowId: okno,
      ...wskazanie,
      ...(dokument === '' ? {} : { documentId: dokument }),
      ...(format.value === '' ? {} : { format: format.value as StudioImportFormat }),
      ...(zapisZnakow.kontrolka.value.trim() === ''
        ? {}
        : { encoding: zapisZnakow.kontrolka.value.trim() }),
      ...(wDokumencie.checked ? { insertAtOffset: miejsceKursora() } : {}),
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowyBledu('Wniesienie pliku do edytora', wynik.blad), false);
      return;
    }
    stan.wchlon(wynik.wynik.document);
    odpowiedz.pokaz(
      `Plik wniesiony do edytora jako dokument ${wynik.wynik.document.id}. ` +
        wstawieniaOpiszWniesienie(wynik.wynik.balance) +
        (wynik.wynik.provenance === undefined
          ? ''
          : ` Pochodzenie zapisane: ${wynik.wynik.provenance.id}.`),
      true,
    );
  }

  async function wniesPdfDoEdytora(): Promise<void> {
    const okno = idOkna();
    const wskazanie = wskazaniePliku();
    if (okno === null || wskazanie === null) return;
    const dokument = doDokumentu.kontrolka.value.trim();
    const wynik = await zrodlo.wniesPdf({
      windowId: okno,
      ...wskazanie,
      ...(dokument === '' ? {} : { documentId: dokument }),
      ...(zakresStron.kontrolka.value.trim() === ''
        ? {}
        : { pages: zakresStron.kontrolka.value.trim() }),
      recoverTables: zTabelami.checked,
      recoverImages: zObrazami.checked,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowyBledu('Zamiana PDF na dokument edytowalny', wynik.blad), false);
      return;
    }
    stan.wchlon(wynik.wynik.document);
    const doRozpoznania = wynik.wynik.ingestItemId;
    odpowiedz.pokaz(
      `PDF wniesiony do edytora jako dokument ${wynik.wynik.document.id}. ` +
        wstawieniaOpiszWniesienie(wynik.wynik.balance) +
        (doRozpoznania === undefined || doRozpoznania === ''
          ? ''
          : ` Rdzeń założył pozycję kolejki rozpoznania tekstu: ${doRozpoznania} — ten PDF nie ma ` +
            'warstwy tekstowej i konwersji z niego zrobić nie da się. Otwórz narzędziownię ' +
            'cyfryzacji, rozpoznaj tę pozycję i przyjmij wynik.'),
      true,
    );
  }

  async function wniesObrazDoEdytora(): Promise<void> {
    const dokument = stan.dokument();
    if (dokument === null) {
      odpowiedz.pokaz(BRAK_DOKUMENTU, false);
      return;
    }
    const wartosc = wskazanieObrazu.kontrolka.value.trim();
    if (wartosc === '') {
      odpowiedz.pokaz(BRAK_OBRAZU, false);
      return;
    }
    const zrodloWybrane = zrodloObrazu.value as StudioObjectSource;
    const wskazanie: Record<string, string> =
      zrodloWybrane === StudioObjectSource.CoreAsset
        ? { assetId: wartosc }
        : zrodloWybrane === StudioObjectSource.DesignModule
          ? { designNodeId: wartosc }
          : zrodloWybrane === StudioObjectSource.PhotoBank
            ? { photoBankId: wartosc }
            : zrodloWybrane === StudioObjectSource.LibraryFile
              ? { libraryFileId: wartosc }
              : { path: wartosc };

    const wynik = await zrodlo.wniesObraz({
      documentId: dokument.id,
      offset: miejsceKursora(),
      source: zrodloWybrane,
      ...wskazanie,
      ...(liczba(szerokoscObrazu.value) > 0 ? { widthMm: liczba(szerokoscObrazu.value) } : {}),
      ...(liczba(wysokoscObrazu.value) > 0 ? { heightMm: liczba(wysokoscObrazu.value) } : {}),
      ...(tekstZastepczy.kontrolka.value.trim() === ''
        ? {}
        : { altText: tekstZastepczy.kontrolka.value.trim() }),
      ...(podpisObrazu.kontrolka.value.trim() === ''
        ? {}
        : { caption: podpisObrazu.kontrolka.value.trim() }),
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowyBledu('Wniesienie obrazu', wynik.blad), false);
      return;
    }
    odpowiedz.pokaz(
      `Obraz osadzony jako ${wynik.wynik.object.id} na znaku ${miejsceKursora()}. ` +
        wstawieniaOpiszBilans(wynik.wynik.balance) +
        (wynik.wynik.provenance === undefined
          ? ' Rdzeń nie oddał zapisu pochodzenia tego obrazu.'
          : ` Pochodzenie zapisane: ${wynik.wynik.provenance.id}.`),
      wynik.wynik.balance.applied > 0,
    );
  }

  /* ── Nakładka ──────────────────────────────────────────────────────────── */

  const tresc = document.createElement('div');
  tresc.className = 'ms-wniesienie__tresc';
  tresc.append(
    czescWniesienia('Wskazanie pliku', [
      sciezka.element,
      plikBiblioteki.element,
      zasob.element,
      doDokumentu.element,
      etykietaMiejsca,
    ]),
    czescWniesienia('Plik dokumentu i szablonu', [
      format,
      zapisZnakow.element,
      wniesPlik,
    ]),
    czescWniesienia('PDF na dokument edytowalny', [
      zakresStron.element,
      etykietaTabel,
      etykietaObrazow,
      oOdzyskaniu,
      wniesPdf,
    ]),
    czescWniesienia('Obraz w miejsce kursora', [
      zrodloObrazu,
      wskazanieObrazu.element,
      szerokoscObrazu,
      wysokoscObrazu,
      tekstZastepczy.element,
      podpisObrazu.element,
      wniesObraz,
    ]),
    odpowiedz.element,
  );

  const nakladka = document.createElement('div');
  nakladka.className = 'ms-wniesienie__nakladka';
  nakladka.hidden = true;
  nakladka.append(tresc);

  const wyzwalacz = przycisk('Wnieś plik do edytora ▾', 'dn-btn dn-btn--sm dn-btn--zarys');
  wyzwalacz.dataset['czynnosc'] = 'warsztat-wniesienia';
  wyzwalacz.setAttribute('aria-expanded', 'false');
  wyzwalacz.title =
    'Otwiera wniesienie pliku nakładką: dokument i szablon Operatora (docx, dotx, odt, ott, tekst ' +
    'czysty, markdown, RTF, HTML) z rozpoznaniem zapisu znaków, PDF na dokument edytowalny wraz ' +
    'z bilansem odzyskania oraz obraz w miejsce kursora. Schodzi drugim naciśnięciem albo Escapem.';
  wyzwalacz.addEventListener('click', () => przestaw(nakladka.hidden));

  const element = document.createElement('section');
  element.className = 'ms-wniesienie';
  element.setAttribute('aria-label', 'Wniesienie pliku wprost do edytora');
  element.append(wyzwalacz, nakladka);
  element.addEventListener('keydown', (zdarzenie) => {
    if (zdarzenie.key === 'Escape' && !nakladka.hidden) przestaw(false);
  });

  function przestaw(otwarta: boolean): void {
    nakladka.hidden = !otwarta;
    wyzwalacz.setAttribute('aria-expanded', otwarta ? 'true' : 'false');
    wyzwalacz.textContent = `Wnieś plik do edytora ${otwarta ? '▴' : '▾'}`;
  }

  return { element, przestawWidocznosc: () => przestaw(nakladka.hidden) };
}

function liczba(wartosc: string): number {
  const odczytana = Number.parseFloat(wartosc);
  return Number.isFinite(odczytana) ? odczytana : 0;
}

function czescWniesienia(tytul: string, elementy: readonly HTMLElement[]): HTMLElement {
  const naglowek = document.createElement('p');
  naglowek.className = 'ms-wniesienie__tytul';
  naglowek.textContent = tytul;
  const sekcja = document.createElement('section');
  sekcja.className = 'ms-wniesienie__czesc';
  sekcja.append(naglowek, ...elementy);
  return sekcja;
}

const BRAK_OKNA =
  'Wniesienie pliku należy do okna komunikacji sesji — wskaż je w pasie osadzenia modułu. Komendy ' +
  'studio.document.import.* przyjmują identyfikator okna jako pole obowiązkowe.';

const BRAK_WSKAZANIA =
  'Wskaż plik: ścieżkę widzianą przez rdzeń, plik Biblioteki albo zasób magazynu. Klient dysku nie ' +
  'czyta i nie ma czego wysłać.';

const BRAK_DOKUMENTU =
  'Obraz wchodzi w miejsce kursora w dokumencie — wczytaj go albo załóż nowy. Bez dokumentu nie ma ' +
  'gdzie osadzić obrazu.';

const BRAK_OBRAZU =
  'Wskaż obraz — rdzeń nie osadzi niczego bez ścieżki, zasobu, węzła Designu, pozycji bazy ' +
  'zdjęciowej albo pliku Biblioteki.';

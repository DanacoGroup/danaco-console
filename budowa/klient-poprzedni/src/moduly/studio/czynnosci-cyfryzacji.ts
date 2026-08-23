import {
  StudioIngestState,
  type StudioDocumentFormat,
  type StudioIngestQueueAddRequest,
  type StudioInputDevice,
  type StudioRecognitionSettings,
} from '../../../../shared/contract';
import { opisOdmowyBledu } from '../../komponenty/odmowa';
import type { WierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import type { KolejkaCyfryzacji } from './kolejka-cyfryzacji';
import type { StanOknaStudio } from './stan-okna-studio';
import type { StanStudio } from './stan-studio';
import { SKLADNIKI_PAKIETU_SERWERA } from './braki-cyfryzacji';
import type { ZrodloWstawienStudio } from './zrodlo-wstawien-studio';

/**
 * Czynności narzędziowni cyfryzacji — sześć komend rodziny `studio.ingest.*`.
 *
 * ── Co się zmieniło wobec poprzedniej postaci ───────────────────────────────
 * Wsad szedł `document.text.extract`: jedno wydobycie na jedno wywołanie, jeden
 * język, bez silnika, bez progu pewności i bez pojęcia kolejki. Teraz kolejkę
 * prowadzi rdzeń (`ingest.queue.add`, `ingest.queue.list`), rozpoznanie ma pełne
 * sterowanie (`ingest.recognize`), poprawka słowa wchodzi przed przyjęciem
 * (`ingest.correction.set`), a przyjęcie zakłada DOKUMENT wraz z pierwszą wersją
 * (`ingest.item.accept`) — nie sam bufor edytora.
 *
 * ── Przekazanie do edytora przestało być półśrodkiem ────────────────────────
 * Poprzednio wynik szedł do bufora edytora i panel musiał tłumaczyć, że wersji
 * pierwszej nie ma, bo wydobycie tekstu żadnej nie zakłada. `ingest.item.accept`
 * zakłada dokument I wersję, więc panel może wreszcie powiedzieć prawdę bez
 * zastrzeżenia. Dokument wchodzi do stanu modułu, żeby okno pracy zobaczyło go
 * natychmiast.
 *
 * ── Odmowa jednej pozycji nie zatrzymuje pozostałych ────────────────────────
 * Zostaje bez zmian: zatrzymanie całej kolejki na pierwszym nieczytelnym skanie
 * byłoby karą za materiał, a nie obsługą błędu. Bilans na końcu mówi, ile
 * przeszło, ile wróciło do ponowienia i ile odpadło.
 */
export interface ZapleczeCyfryzacji {
  stan: StanStudio;
  /** Źródło rodziny `studio.ingest.*` wraz z pozostałymi wstawieniami. */
  wstawienia: ZrodloWstawienStudio;
  kolejka: KolejkaCyfryzacji;
  pas: StanOknaStudio;
  odpowiedz: WierszOdpowiedzi;
  /** Nastawy rozpoznawania wspólne pozycjom wsadu. */
  ustawienia(): StudioRecognitionSettings;
  /** Tytuł dokumentu zakładanego przy przyjęciu; pusty zostawia wybór rdzeniowi. */
  tytul(): string;
  /** Format dokumentu zakładanego przy przyjęciu; pusty zostawia wybór rdzeniowi. */
  format(): string;
  /** Przerysowuje widok po każdej zmianie stanu pozycji. */
  odswiez(): void;
}

/** Odczytuje kolejkę rdzenia do odbicia w oknie. */
export async function odczytajKolejke(zaplecze: ZapleczeCyfryzacji): Promise<void> {
  const idOkna = zaplecze.stan.idOkna();
  if (idOkna === '') return;
  const wynik = await zaplecze.wstawienia.kolejka({ windowId: idOkna });
  if (!wynik.udany || wynik.wynik === undefined) {
    zaplecze.odpowiedz.pokaz(opisOdmowyBledu('Odczyt kolejki wczytywania', wynik.blad), false);
    return;
  }
  zaplecze.kolejka.ustawPozycje(wynik.wynik.items);
  zaplecze.odswiez();
}

/** Dokłada wskazany materiał do kolejki rdzenia; wiele wskazań jednym żądaniem. */
export async function dolozMaterial(
  zaplecze: ZapleczeCyfryzacji,
  wskazanie: string,
  czyZasob: boolean,
): Promise<void> {
  const idOkna = zaplecze.stan.idOkna();
  if (idOkna === '') {
    zaplecze.odpowiedz.pokaz(BRAK_OKNA, false);
    return;
  }
  const wskazania = rozbij(wskazanie);
  if (wskazania.length === 0) {
    zaplecze.odpowiedz.pokaz(BRAK_WSKAZANIA, false);
    return;
  }
  const zadanie: StudioIngestQueueAddRequest = { windowId: idOkna };
  if (czyZasob) zadanie.assetIds = wskazania;
  else zadanie.sourcePaths = wskazania;
  const nastawy = zaplecze.ustawienia();
  // Nastawy wspólne wsadu wchodzą już przy dołożeniu, bo kontrakt je tu
  // przyjmuje: pozycja niesie wtedy swoje nastawy i rozpoznanie nie musi ich
  // powtarzać przy każdym wywołaniu.
  if (Object.keys(nastawy).length > 0) zadanie.settings = nastawy;

  zaplecze.pas.ladowanie('Dołożenie materiału do kolejki rdzenia w toku…');
  const wynik = await zaplecze.wstawienia.dolozDoKolejki(zadanie);
  if (!wynik.udany || wynik.wynik === undefined) {
    zaplecze.pas.blad(opisOdmowyBledu('Dołożenie do kolejki', wynik.blad));
    return;
  }
  zaplecze.kolejka.dolozPozycje(wynik.wynik.items);
  zaplecze.pas.gotowe();
  zaplecze.odpowiedz.pokaz(
    `Do kolejki rdzenia weszło pozycji: ${wynik.wynik.items.length}. Kolejka stoi po stronie ` +
      'rdzenia, więc przeżywa odświeżenie okna.',
    wynik.wynik.items.length > 0,
  );
  zaplecze.odswiez();
}

/** Rozpakowuje archiwum wsadu wprost do kolejki rdzenia. */
export async function dolozArchiwum(
  zaplecze: ZapleczeCyfryzacji,
  sciezka: string,
): Promise<void> {
  const idOkna = zaplecze.stan.idOkna();
  if (idOkna === '') {
    zaplecze.odpowiedz.pokaz(BRAK_OKNA, false);
    return;
  }
  if (sciezka === '') {
    zaplecze.odpowiedz.pokaz(BRAK_ARCHIWUM, false);
    return;
  }
  zaplecze.pas.ladowanie('Rozpakowanie archiwum wsadu do kolejki rdzenia w toku…');
  const wynik = await zaplecze.wstawienia.dolozDoKolejki({
    windowId: idOkna,
    archivePath: sciezka,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    // Odmowa z powodu rozpakowywacza jest usterką WDROŻENIA, nie brakiem funkcji,
    // i tak musi być nazwana — inaczej Operator uzna, że produkt wsadu nie ma.
    zaplecze.pas.blad(
      `${opisOdmowyBledu('Rozpakowanie archiwum', wynik.blad)} ${SKLADNIKI_PAKIETU_SERWERA.archiwum}`,
    );
    return;
  }
  zaplecze.kolejka.dolozPozycje(wynik.wynik.items);
  zaplecze.pas.gotowe();
  zaplecze.odpowiedz.pokaz(
    `Archiwum rozpakowane po stronie rdzenia: do kolejki weszło pozycji ` +
      `${wynik.wynik.items.length}.`,
    wynik.wynik.items.length > 0,
  );
  zaplecze.odswiez();
}

/** Rozpoznaje tekst jednej pozycji kolejki wraz z pełnym sterowaniem. */
export async function rozpoznajPozycje(
  zaplecze: ZapleczeCyfryzacji,
  idPozycji: string,
): Promise<void> {
  const pozycja = zaplecze.kolejka.pozycja(idPozycji);
  if (pozycja === null) return;
  if (zaplecze.kolejka.wToku(idPozycji)) return;

  const nastawy = zaplecze.ustawienia();
  zaplecze.kolejka.ustawWToku(idPozycji, true);
  zaplecze.odswiez();

  const wynik = await zaplecze.wstawienia.rozpoznaj({
    itemId: idPozycji,
    // Brak nastaw znaczy „weź nastawy kolejki" — to rozstrzygnięcie rdzenia,
    // a nie okna, więc pustego obiektu nie wysyłamy.
    ...(Object.keys(nastawy).length === 0 ? {} : { settings: nastawy }),
  });
  zaplecze.kolejka.ustawWToku(idPozycji, false);

  if (!wynik.udany || wynik.wynik === undefined) {
    zaplecze.odpowiedz.pokaz(
      `${opisOdmowyBledu('Rozpoznanie tekstu', wynik.blad)} ${SKLADNIKI_PAKIETU_SERWERA.rozpoznanie}`,
      false,
    );
    zaplecze.odswiez();
    return;
  }
  zaplecze.kolejka.ustawPozycje1(wynik.wynik.item);
  zaplecze.kolejka.ustawRozpoznanie(idPozycji, wynik.wynik.words, wynik.wynik.layout);
  zaplecze.odpowiedz.pokaz(opiszRozpoznanie(wynik.wynik.item.state, wynik.wynik), true);
  zaplecze.odswiez();
}

/** Rozpoznaje wszystkie pozycje oczekujące i wracające do ponowienia. */
export async function rozpoznajKolejke(zaplecze: ZapleczeCyfryzacji): Promise<void> {
  if (zaplecze.kolejka.nastepnaDoRozpoznania() === null) {
    zaplecze.odpowiedz.pokaz(BRAK_OCZEKUJACYCH, false);
    return;
  }
  zaplecze.pas.ladowanie('Rozpoznanie kolejki wczytywania w toku…');
  // Pozycje przechodzimy po migawce identyfikatorów, a nie pytając kolejki
  // w pętli o „następną": pozycja, której rdzeń nie ruszył ze stanu oczekiwania,
  // wracałaby wtedy w nieskończoność.
  const doRozpoznania = zaplecze.kolejka
    .pozycje()
    .filter(
      (pozycja) =>
        pozycja.state === StudioIngestState.Oczekuje ||
        pozycja.state === StudioIngestState.Ponowienie,
    )
    .map((pozycja) => pozycja.id);
  for (const idPozycji of doRozpoznania) {
    await rozpoznajPozycje(zaplecze, idPozycji);
  }
  const bilans = zaplecze.kolejka.bilans();
  zaplecze.pas.gotowe();
  zaplecze.odpowiedz.pokaz(
    `Kolejka rozpoznana: gotowych ${bilans.gotowe}, do ponowienia ${bilans.ponowienia}, ` +
      `odmów ${bilans.odmowy}, pozycji łącznie ${bilans.wszystkie}. Pozycja w stanie ponowienia ` +
      'wypadła poniżej progu pewności — obniż próg albo popraw obraz i rozpoznaj ją jeszcze raz. ' +
      'Powód odmowy stoi przy pozycji.',
    bilans.odmowy === 0,
  );
}

/** Poprawia rozpoznane słowo na warstwie tekstowej PRZED przyjęciem do edytora. */
export async function poprawSlowo(
  zaplecze: ZapleczeCyfryzacji,
  idPozycji: string,
  numerSlowa: number,
  tresc: string,
): Promise<void> {
  if (tresc === '') {
    zaplecze.odpowiedz.pokaz(PUSTA_POPRAWKA, false);
    return;
  }
  const wynik = await zaplecze.wstawienia.popraw({
    itemId: idPozycji,
    wordIndex: numerSlowa,
    text: tresc,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    zaplecze.odpowiedz.pokaz(opisOdmowyBledu('Poprawka rozpoznanego słowa', wynik.blad), false);
    return;
  }
  zaplecze.kolejka.ustawPozycje1(wynik.wynik.item);
  zaplecze.odpowiedz.pokaz(
    `Słowo o numerze ${numerSlowa} poprawione na „${tresc}" na warstwie tekstowej pozycji. ` +
      'Poprawka weszła do kolejki, a nie do dokumentu — dokument powstanie przy przyjęciu wyniku.',
    true,
  );
  zaplecze.odswiez();
}

/**
 * Przyjmuje wynik pozycji jako dokument roboczy wraz z pierwszą wersją.
 *
 * Wiele pozycji wskazanych naraz składa się w JEDEN dokument w kolejności
 * podania — tak wchodzi skan wielostronicowy rozłożony na osobne pliki.
 */
export async function przyjmijWynik(
  zaplecze: ZapleczeCyfryzacji,
  idPozycji: readonly string[],
): Promise<void> {
  const idOkna = zaplecze.stan.idOkna();
  if (idOkna === '') {
    zaplecze.odpowiedz.pokaz(BRAK_OKNA, false);
    return;
  }
  if (idPozycji.length === 0) {
    zaplecze.odpowiedz.pokaz(BRAK_WYNIKU, false);
    return;
  }
  const bezTekstu = idPozycji.filter((identyfikator) => {
    const pozycja = zaplecze.kolejka.pozycja(identyfikator);
    return pozycja === null || (pozycja.text ?? '') === '';
  });
  if (bezTekstu.length > 0) {
    zaplecze.odpowiedz.pokaz(
      `Pozycji bez tekstu: ${bezTekstu.length} (${bezTekstu.join(', ')}). Rdzeń nie ma z nich czego ` +
        'złożyć w dokument. Rozpoznaj je najpierw albo przeczytaj powód odmowy stojący przy pozycji.',
      false,
    );
    return;
  }

  const format = zaplecze.format();
  const tytul = zaplecze.tytul();
  zaplecze.pas.ladowanie('Przyjęcie wyniku cyfryzacji do edytora w toku…');
  const wynik = await zaplecze.wstawienia.przyjmij({
    windowId: idOkna,
    itemIds: [...idPozycji],
    ...(tytul === '' ? {} : { title: tytul }),
    ...(format === '' ? {} : { format: format as StudioDocumentFormat }),
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    zaplecze.pas.blad(opisOdmowyBledu('Przyjęcie wyniku cyfryzacji', wynik.blad));
    return;
  }
  // Dokument wchodzi do stanu modułu, więc okno pracy widzi go natychmiast —
  // przyjęcie kończy się dokumentem, nie samą treścią w buforze.
  zaplecze.stan.wchlon(wynik.wynik.document);
  zaplecze.pas.gotowe();
  zaplecze.odpowiedz.pokaz(
    `Dokument ${wynik.wynik.document.id} („${wynik.wynik.document.title ?? 'bez tytułu'}") ` +
      `złożony z pozycji: ${idPozycji.length}. Pierwsza wersja w repozytorium sesji: ` +
      `${wynik.wynik.version.id}. Dokument stoi w oknie pracy i jest gotowy do redakcji.`,
    true,
  );
  zaplecze.odswiez();
}

/**
 * Odczytuje urządzenia wejściowe maszyny rdzenia.
 *
 * Wykaz oddaje się wywołaniem, a nie zwracaną wartością, bo wstawia go kontrolka
 * wyboru w polach panelu — a ta wie o urządzeniach wszystko, czego potrzebuje do
 * etykiety, i nie ma po co przechodzić przez widok pośredni.
 */
export async function odczytajUrzadzenia(
  zaplecze: ZapleczeCyfryzacji,
  pokaz: (urzadzenia: readonly StudioInputDevice[]) => void,
): Promise<void> {
  const wynik = await zaplecze.wstawienia.urzadzenia();
  if (!wynik.udany || wynik.wynik === undefined) {
    zaplecze.odpowiedz.pokaz(
      `${opisOdmowyBledu('Odczyt urządzeń wejściowych', wynik.blad)} ` +
        SKLADNIKI_PAKIETU_SERWERA.urzadzenie,
      false,
    );
    return;
  }
  pokaz(wynik.wynik.devices);
  zaplecze.odpowiedz.pokaz(
    wynik.wynik.devices.length === 0
      ? 'Rdzeń nie widzi ani jednego skanera i ani jednej kamery. To odpowiedź, nie odmowa: ' +
          'urządzenie należy do maszyny rdzenia i tam musi być podłączone.'
      : `Rdzeń widzi urządzeń wejściowych: ${wynik.wynik.devices.length}.`,
    wynik.wynik.devices.length > 0,
  );
}

/** Zdanie o wyniku rozpoznania — stan, pewność i to, co wróciło ze strukturą. */
function opiszRozpoznanie(
  stan: StudioIngestState,
  wynik: { words?: readonly unknown[]; layout?: readonly unknown[]; item: { confidence?: number } },
): string {
  const pewnosc =
    wynik.item.confidence === undefined
      ? 'pewności rdzeń nie podał'
      : `średnia pewność rozpoznania ${wynik.item.confidence}`;
  const slowa =
    wynik.words === undefined
      ? 'słów z położeniem rdzeń nie oddał, więc poprawiania na obrazie nie ma na czym oprzeć'
      : `słów z położeniem i pewnością: ${wynik.words.length} — da się je poprawić przed przyjęciem`;
  const uklad =
    wynik.layout === undefined
      ? 'układu rdzeń nie odtwarzał — włącz odtwarzanie układu, jeśli materiał ma kolumny albo tabele'
      : `bloków układu: ${wynik.layout.length}`;
  if (stan === StudioIngestState.Ponowienie) {
    return (
      `Pozycja wróciła do PONOWIENIA: ${pewnosc}, poniżej ustawionego progu. To nie odmowa — ` +
      `obniż próg, popraw obraz albo dołóż język i rozpoznaj ponownie. ${slowa}. ${uklad}.`
    );
  }
  if (stan === StudioIngestState.Odmowa) {
    return `Rdzeń odmówił rozpoznania tej pozycji. Powód stoi przy pozycji w kolejce.`;
  }
  return `Pozycja rozpoznana: ${pewnosc}. ${slowa}. ${uklad}.`;
}

/** Rozbija wskazania rozdzielone przecinkiem; puste odpadają. */
function rozbij(wartosc: string): string[] {
  return wartosc
    .split(',')
    .map((czesc) => czesc.trim())
    .filter((czesc) => czesc !== '');
}

const BRAK_OKNA =
  'Kolejka wczytywania należy do okna komunikacji sesji — wskaż je w pasie osadzenia modułu. ' +
  'Komendy studio.ingest.* przyjmują identyfikator okna jako pole obowiązkowe.';

const BRAK_WSKAZANIA =
  'Wskaż materiał — rdzeń nie założy pozycji kolejki bez identyfikatora zasobu albo ścieżki. ' +
  'Wiele wskazań rozdziel przecinkiem; wejdą jednym żądaniem.';

const BRAK_ARCHIWUM =
  'Wskaż ścieżkę archiwum — rozpakowanie bierze plik po stronie rdzenia, a klient dysku nie czyta.';

const BRAK_OCZEKUJACYCH =
  'Kolejka nie ma pozycji oczekującej ani wracającej do ponowienia. Dołóż wskazanie materiału ' +
  'i uruchom rozpoznanie ponownie.';

const BRAK_WYNIKU =
  'Przyjęcie dotyczy pozycji wskazanej — zaznacz w kolejce tę, której wynik ma wejść do edytora.';

const PUSTA_POPRAWKA =
  'Poprawka pusta zdjęłaby słowo z warstwy tekstowej, a nie poprawiła — wpisz brzmienie właściwe.';

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

/** Czynności narzędziowni cyfryzacji, sześć komend rodziny studio.ingest.*: kolejka, rozpoznanie, poprawka słowa i przyjęcie wyniku jako dokument z pierwszą wersją. */
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

/** Odczytuje kolejkę wczytywania z rdzenia i odbija jej pozycje w stanie okna, żeby widok pokazywał aktualną zawartość kolejki. */
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

/** Dokłada wskazany materiał, jedną albo wieloma ścieżkami lub zasobami naraz, do kolejki wczytywania rdzenia jednym żądaniem. */
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
  // Nastawy wspólne wsadu wchodzą już przy dołożeniu, więc rozpoznanie nie musi ich powtarzać.
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

/** Rozpakowuje wskazane archiwum wsadu po stronie rdzenia i dokłada powstałe z niego pozycje wprost do kolejki wczytywania. */
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
    // Odmowa z powodu rozpakowywacza jest usterką wdrożenia, nie brakiem funkcji, i tak ma być nazwana.
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

/** Rozpoznaje tekst jednej pozycji kolejki, przekazując rdzeniowi pełne nastawy sterujące silnikiem, jeśli okno je ustawiło. */
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
    // Brak nastaw znaczy: weź nastawy kolejki — rozstrzyga rdzeń, więc pustego obiektu nie wysyłamy.
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

/** Rozpoznaje kolejno wszystkie pozycje kolejki będące w stanie oczekiwania albo wracające do ponowienia po niskiej pewności. */
export async function rozpoznajKolejke(zaplecze: ZapleczeCyfryzacji): Promise<void> {
  if (zaplecze.kolejka.nastepnaDoRozpoznania() === null) {
    zaplecze.odpowiedz.pokaz(BRAK_OCZEKUJACYCH, false);
    return;
  }
  zaplecze.pas.ladowanie('Rozpoznanie kolejki wczytywania w toku…');
  // Pozycje przechodzimy po migawce identyfikatorów, inaczej pytanie o następną wracałoby bez końca.
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

/** Poprawia treść jednego rozpoznanego słowa na warstwie tekstowej pozycji kolejki, zanim pozycja zostanie przyjęta do edytora. */
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
  // Dokument wchodzi do stanu modułu, więc okno pracy widzi go od razu, nie dopiero po samej treści.
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

/** Odczytuje urządzenia wejściowe maszyny rdzenia i przekazuje wykaz wywołaniem do kontrolki wyboru pól panelu, która sama zna wszystko potrzebne do etykiety. */
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

/** Zdanie opisujące wynik rozpoznania tekstu: stan pozycji, osiągniętą pewność oraz to, czy rdzeń oddał także strukturę układu. */
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

/** Rozbija ciąg wskazań materiału rozdzielonych przecinkiem na tablicę pojedynczych wskazań, odrzucając wpisy puste. */
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

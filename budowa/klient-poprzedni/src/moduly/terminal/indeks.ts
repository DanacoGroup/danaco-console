import './terminal.css';

import type { Kanal } from '../../protokol/kanal';
import { utworzPokrycieKomend, type PokrycieKomend } from '../pokrycie-komend';
import { widokZOknaSesji, type OpisModulu } from '../rejestracja';
import { WIDOCZNOSCI, type CzynnoscOkna } from './czynnosci-okna';
import { utworzOknoKonsoli } from './okno-output-console';
import { utworzOknoMonitora } from './okno-process-monitor';
import { utworzOknoBibliotekiSkryptow } from './okno-script-library';
import { utworzOknoZarzadcySesji } from './okno-session-manager';
import { utworzOknoZadanIHarmonogramu } from './okno-task-schedule';
import { utworzOknoKart } from './okno-terminal-tabs';
import { utworzPaletePolecen } from './paleta-polecen';
import { WIAZANIA, wykazSkrotow, zwiazSkroty, type CzynnoscSkrotu } from './skroty-klawiszowe';
import { utworzStanTerminala, type StanTerminala } from './stan-terminala';
import { utworzWyborDrzewem } from './wybor-drzewem';
import { KOD_MODULU } from './wynik-czastkowy';
import { utworzZrodloTerminala, type ZrodloTerminala } from './zrodlo-terminala';

/**
 * Moduł Terminal — złożenie sześciu okien operacyjnych wokół jednego okna
 * komunikacji.
 *
 * Układ wynika z ról okien. Terminal Tabs jest oknem wiodącym, więc stoi
 * w obszarze głównym jako punkt wejścia. Output Console i Process Monitor są
 * oknami monitora — stoją pod nim w pasie obserwacyjnym, bo obserwuje się to,
 * co okno wiodące uruchomiło. Session Manager, Task & Schedule i Script Library
 * są oknami zarządcy i kreatora: przygotowują pracę, którą okno wiodące
 * wykonuje, więc stoją w pasie rozszerzeń pod pasem obserwacyjnym. Okno rozmowy
 * modułu (Chat Window) i okno pętli wykonawczej nie należą do tego złożenia:
 * są bytami sesji i składa je warstwa rozmowy.
 *
 * Wszystkie sześć okien patrzy na jeden stan modułu, więc karta otwarta
 * w Session Managerze jest od razu kartą bieżącą okna wiodącego, a proces
 * uruchomiony w Task & Schedule stoi w Process Monitorze i w Output Console
 * bez żadnego przekazywania między oknami.
 *
 * Nad oknami stoi pas modułu z dwiema nastawami wspólnymi: widocznością warstw
 * i paletą poleceń. Warstwy zdejmują z ekranu kontrolki, których bieżące
 * zadanie nie wymaga; paleta pilnuje, żeby zdjęcie z ekranu nie stało się
 * schowaniem — każda czynność każdego okna jest w niej o jedno wskazanie.
 *
 * Sześć okien, pas nastaw, skróty klawiszowe i załącznik skrótów należą do
 * pozycji samodzielnej. Złożenie wchodzące jako okno pomocnicze gospodarza
 * niesie trzy okna rdzenne — dokładnie to, co obiecuje rama tamtego okna.
 *
 * Moduł widać wyłącznie w środowisku CodeStudio. Macierz widoczności jest
 * własnością nawigacji, a nie modułu — złożenie nie sprawdza środowiska samo,
 * bo drugi egzekutor widoczności rozjechałby się z pierwszym.
 *
 * Terminal nie jest samodzielnym modułem, tylko dodatkowym oknem pomocniczym
 * sesji CodeStudio. Złożenie wchodzi więc dwiema drogami:
 *
 *  1. jako okno pomocnicze gospodarza — `moduly/terminal/okno-pomocnicze.ts`,
 *     wołane z modułów Developer, Diagnostics i Apps; karta powstaje w oknie
 *     gospodarza i dziedziczy jego tryb uprawnień oraz katalog roboczy,
 *  2. jako samodzielna pozycja nawigacji. Nie da się jej zdjąć z klienta:
 *     pozycja pochodzi z macierzy widoczności rdzenia
 *     (`store/migracja_007_zaczyn_slownikow.sql`, wiersz `('codestudio',
 *     'terminal', 3)`), a klient bierze wykaz modułów wyłącznie z
 *     `environment.enter` / `module.list`. Dopóki rdzeń pozycję oddaje, moduł
 *     mówi wprost, czym jest — zniknięcie widoku dałoby pustkę czytającą się
 *     jak usterka.
 */
export interface ZamontowanyTerminal {
  /** Element osadzony w dokumencie. */
  element: HTMLElement;
  /** Stan modułu wspólny trzem oknom. */
  stan: StanTerminala;
  /** Odczytuje okna z rdzenia. */
  odswiez(): void;
  /** Zamyka nasłuch zdarzeń modułu. */
  zamknij(): void;
}

/** Zależności złożenia. */
export interface OpcjeTerminala {
  /** Okno komunikacji, w którym pracuje moduł — bez niego nie ma czego otworzyć. */
  okno: string;
  /**
   * Złożenie stoi jako samodzielna pozycja nawigacji, a nie w oknie pomocniczym
   * gospodarza. Wtedy — i tylko wtedy — nad oknami staje nota mówiąca, czym
   * Terminal jest. W pasie okien pomocniczych ta sama nota byłaby powtórzeniem
   * zdania, które rama okna pomocniczego już niesie.
   */
  samodzielny?: boolean;
}

export function zamontujTerminal(
  gospodarz: HTMLElement,
  kanal: Kanal,
  opcje: OpcjeTerminala,
): ZamontowanyTerminal {
  const zrodlo = utworzZrodloTerminala(kanal);
  const stan = utworzStanTerminala(zrodlo, { okno: opcje.okno });
  // Pozycje, których okna jeszcze nie wykonują, biorą swoje zdanie z odczytu
  // wykazu komend rdzenia, a nie z napisu w kodzie. Po scaleniu kontraktu
  // wszystkie komendy tego modułu w nim stoją, więc zdanie „kontrakt tego nie
  // ma" byłoby dziś nieprawdą — brak przeszedł na stronę rdzenia i tylko odczyt
  // potrafi powiedzieć, kiedy przestanie tam być.
  const pokrycie = utworzPokrycieKomend(kanal);

  const obszar = document.createElement('div');
  obszar.className = 'dt-modul';
  obszar.dataset['modul'] = KOD_MODULU;
  // Widoczność podstawowa jest stanem wyjściowym: funkcja niepotrzebna do
  // bieżącego zadania nie stoi na ekranie. To nie jest blokada — kontrolka
  // zdjęta z pola widzenia działa i sięga po nią paleta poleceń.
  obszar.dataset['widocznosc'] = 'podstawowa';

  const karty = utworzOknoKart(zrodlo, stan, pokrycie);
  const konsola = utworzOknoKonsoli(zrodlo, stan, pokrycie);
  const monitor = utworzOknoMonitora(zrodlo, stan, pokrycie);

  const gorny = document.createElement('div');
  gorny.className = 'dt-modul__gora';
  gorny.append(karty.element);

  const dolny = document.createElement('div');
  dolny.className = 'dt-modul__dol';
  dolny.append(konsola.element, monitor.element);

  // Pas rozszerzeń, nastawy wspólne i skróty wchodzą wyłącznie w pozycji
  // samodzielnej. W oknie pomocniczym gospodarza rama obiecuje kartę powłoki,
  // jej procesy i wyjście na żywo — trzy okna. Doklejenie tam zarządcy sesji,
  // zadań i biblioteki rozrosłoby panel wewnątrz cudzego modułu, a przejęcie
  // skrótów klawiszowych odebrałoby je gospodarzowi.
  const rozszerzenia =
    opcje.samodzielny === true ? zlozRozszerzenia(zrodlo, stan, pokrycie, obszar) : null;

  if (opcje.samodzielny === true) obszar.append(notaOknaPomocniczego());
  if (rozszerzenia !== null) obszar.append(rozszerzenia.pas);
  obszar.append(gorny, dolny);
  if (rozszerzenia !== null) obszar.append(rozszerzenia.element, rozszerzenia.zalacznik);
  gospodarz.replaceChildren(obszar);

  // Paleta i skróty biorą ten sam wykaz czynności: pozycja palety i skrót mają
  // wywoływać dokładnie to samo, a dwa wykazy składane osobno rozjechałyby się
  // przy pierwszej dołożonej czynności.
  const wszystkieCzynnosci: readonly CzynnoscOkna[] = [
    ...karty.czynnosci,
    ...konsola.czynnosci,
    ...monitor.czynnosci,
    ...(rozszerzenia?.czynnosci ?? []),
  ];
  rozszerzenia?.ustawCzynnosci(wszystkieCzynnosci);

  const skroty =
    rozszerzenia === null
      ? null
      : zwiazSkroty(
          obszar,
          przypisaniaSkrotow(wszystkieCzynnosci, rozszerzenia.otworzPalete),
        );

  function odswiez(): void {
    karty.odswiez();
    konsola.odswiez();
    monitor.odswiez();
    rozszerzenia?.odswiez();
  }

  odswiez();
  // Odczyt wykazu komend rdzenia idzie raz na połączenie, nie raz na okno:
  // wszystkie moduły kanału czekają na tę samą odpowiedź.
  void pokrycie.odczytaj();

  return {
    element: obszar,
    stan,
    odswiez,
    zamknij() {
      skroty?.zamknij();
      rozszerzenia?.zwin();
      pokrycie.zamknij();
      stan.zamknij();
    },
  };
}

/** Pas rozszerzeń modułu wraz z nastawami wspólnymi — wyłącznie pozycja samodzielna. */
interface RozszerzeniaModulu {
  /** Pas trzech okien: zarządca sesji, zadania i harmonogram, biblioteka skryptów. */
  element: HTMLElement;
  /** Pas nastaw wspólnych: widoczność warstw i paleta poleceń. */
  pas: HTMLElement;
  /** Załącznik skrótów klawiszowych. */
  zalacznik: HTMLElement;
  czynnosci: readonly CzynnoscOkna[];
  /** Podaje palecie komplet czynności wszystkich okien złożenia. */
  ustawCzynnosci(wszystkie: readonly CzynnoscOkna[]): void;
  otworzPalete(): void;
  odswiez(): void;
  zwin(): void;
}

/**
 * Składa trzy okna pasa rozszerzeń wraz z nastawami wspólnymi modułu.
 *
 * Fragment wyjęty z wytwórni, bo pozycja samodzielna i okno pomocnicze różnią
 * się dokładnie o niego — a wytwórnia z dwiema drogami budowy w środku
 * przestałaby się czytać.
 */
function zlozRozszerzenia(
  zrodlo: ZrodloTerminala,
  stan: StanTerminala,
  pokrycie: PokrycieKomend,
  obszar: HTMLElement,
): RozszerzeniaModulu {
  const zarzadcaSesji = utworzOknoZarzadcySesji(zrodlo, stan, pokrycie);
  const zadania = utworzOknoZadanIHarmonogramu(zrodlo, stan, pokrycie);
  const biblioteka = utworzOknoBibliotekiSkryptow(zrodlo, stan, pokrycie);

  const element = document.createElement('div');
  element.className = 'dt-modul__rozszerzenia';
  element.append(zarzadcaSesji.element, zadania.element, biblioteka.element);

  const paleta = utworzPaletePolecen();

  const widocznosc = utworzWyborDrzewem({
    nastawa: 'Widoczność warstw modułu',
    pozycje: WIDOCZNOSCI,
  });
  widocznosc.naZmiane((wartosc) => {
    obszar.dataset['widocznosc'] = wartosc;
  });

  const pas = document.createElement('div');
  pas.className = 'dt-modul__pas';
  pas.setAttribute('aria-label', 'Nastawy wspólne modułu Terminal');
  pas.append(widocznosc.element, paleta.element);

  const zalacznik = wykazSkrotow();
  zalacznik.dataset['warstwa'] = 'kontekstowa';

  return {
    element,
    pas,
    zalacznik,
    czynnosci: [...zarzadcaSesji.czynnosci, ...zadania.czynnosci, ...biblioteka.czynnosci],
    ustawCzynnosci: (wszystkie) => paleta.ustawCzynnosci(wszystkie),
    otworzPalete: () => paleta.otworz(),
    odswiez() {
      zarzadcaSesji.odswiez();
      zadania.odswiez();
      biblioteka.odswiez();
    },
    zwin: () => paleta.zwin(),
  };
}

/**
 * Wiąże czynności okien ze skrótami klawiszowymi po zapisie skrótu.
 *
 * Zapis skrótu jest jedynym łącznikiem: czynność zna swój skrót, bo pokazuje go
 * w palecie, a wiązanie zna kombinację klawiszy. Dobieranie po zapisie zamyka
 * to w jednym miejscu — czynność bez zapisu nie dostaje wiązania, a wiązanie
 * bez czynności nie przechwytuje klawisza, więc skrót bez skutku nie powstaje.
 *
 * Paleta jest jedynym wyjątkiem: należy do modułu, nie do żadnego okna, więc
 * jej wykonanie wchodzi tu wprost.
 */
function przypisaniaSkrotow(
  czynnosci: readonly CzynnoscOkna[],
  otworzPalete: () => void,
): Partial<Record<CzynnoscSkrotu, () => void>> {
  const przypisania: Partial<Record<CzynnoscSkrotu, () => void>> = {
    'paleta-polecen': otworzPalete,
  };
  for (const wiazanie of WIAZANIA) {
    const czynnosc = czynnosci.find((pozycja) => pozycja.skrot === wiazanie.zapis);
    if (czynnosc === undefined) continue;
    przypisania[wiazanie.czynnosc] = () => czynnosc.wykonaj();
  }
  return przypisania;
}

/**
 * Nota pozycji samodzielnej — zdanie o tym, czym Terminal jest.
 *
 * Nie jest ostrzeżeniem o usterce i nic nie odbiera: okna pod nią działają
 * w pełni. Mówi, że to samo złożenie stoi wewnątrz Developera, Diagnostics
 * i Apps, i skąd bierze się ta pozycja w nawigacji.
 */
function notaOknaPomocniczego(): HTMLElement {
  const element = document.createElement('p');
  element.className = 'dn-pole-opis dt-modul__nota';
  element.setAttribute('role', 'note');
  element.dataset['pozycja'] = 'samodzielna';
  element.textContent =
    'Terminal nie jest samodzielnym modułem. Te same trzy okna wchodzą jako OKNO POMOCNICZE ' +
    'wewnątrz modułów Developer, Diagnostics i Apps środowiska CodeStudio i pracują wtedy ' +
    'w oknie tamtego modułu — z jego trybem uprawnień i katalogiem roboczym. Ta pozycja ' +
    'w bocznej nawigacji pochodzi z macierzy widoczności rdzenia; zdejmuje ją zmiana po ' +
    'stronie rdzenia, nie klienta.';
  return element;
}

/**
 * Samoopisujący się moduł dla rejestru powłoki.
 *
 * Kod siedzi w module, nie w mapie po stronie powłoki: dodanie modułu to jeden
 * wpis, a nie dwa, więc nie da się dodać modułu i zapomnieć o wytwórni.
 */
export const MODUL: OpisModulu = {
  kod: KOD_MODULU,
  utworzWidok: (kanal) =>
    widokZOknaSesji(
      kanal,
      (gospodarz: HTMLElement, k: Kanal, okno: string) =>
        zamontujTerminal(gospodarz, k, { okno, samodzielny: true }),
      // Okno własne, nie pierwsze z brzegu. Bez tego kodu `widokZOknaSesji`
      // bierze `okna[0]` — okno, które w wykazie sesji stoi pierwsze, choćby
      // należało do cudzego modułu. Wszystkie cztery komendy tego modułu niosą
      // `windowId`, więc terminal zakładałby wtedy karty powłoki w cudzym oknie
      // komunikacji.
      KOD_MODULU,
    ),
};

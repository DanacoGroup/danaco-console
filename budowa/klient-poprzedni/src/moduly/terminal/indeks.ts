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
 * Moduł Terminal — złożenie sześciu okien operacyjnych wokół jednego okna komunikacji,
 * ułożonych wedle roli: wiodące, obserwacyjne, zarządcy i kreatora.
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

/** Zależności złożenia modułu Terminal: okno komunikacji, w którym pracuje moduł, i tryb pozycji nawigacji. */
export interface OpcjeTerminala {
  /** Okno komunikacji, w którym pracuje moduł — bez niego nie ma czego otworzyć. */
  okno: string;
  /** Złożenie stoi jako samodzielna pozycja nawigacji, a nie w oknie pomocniczym gospodarza. */
  samodzielny?: boolean;
}

export function zamontujTerminal(
  gospodarz: HTMLElement,
  kanal: Kanal,
  opcje: OpcjeTerminala,
): ZamontowanyTerminal {
  const zrodlo = utworzZrodloTerminala(kanal);
  const stan = utworzStanTerminala(zrodlo, { okno: opcje.okno });
  // Pozycje bez wykonania biorą zdanie z odczytu komend rdzenia, nie z napisu w kodzie.
  const pokrycie = utworzPokrycieKomend(kanal);

  const obszar = document.createElement('div');
  obszar.className = 'dt-modul';
  obszar.dataset['modul'] = KOD_MODULU;
  // Widoczność podstawowa jest stanem wyjściowym; to nie jest blokada, tylko zdjęcie z ekranu.
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

  // Pas rozszerzeń i skróty wchodzą tylko w pozycji samodzielnej, nie u gospodarza.
  const rozszerzenia =
    opcje.samodzielny === true ? zlozRozszerzenia(zrodlo, stan, pokrycie, obszar) : null;

  if (opcje.samodzielny === true) obszar.append(notaOknaPomocniczego());
  if (rozszerzenia !== null) obszar.append(rozszerzenia.pas);
  obszar.append(gorny, dolny);
  if (rozszerzenia !== null) obszar.append(rozszerzenia.element, rozszerzenia.zalacznik);
  gospodarz.replaceChildren(obszar);

  // Paleta i skróty biorą ten sam wykaz czynności, żeby oba wywołania nie rozjechały się z czasem.
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
  // Odczyt wykazu komend rdzenia idzie raz na połączenie, nie raz na okno.
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

/** Pas rozszerzeń modułu Terminal wraz z nastawami wspólnymi — wyłącznie pozycja samodzielna nawigacji. */
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
 * Wiąże czynności okien ze skrótami klawiszowymi po zapisie skrótu; zapis skrótu jest jedynym
 * łącznikiem między nimi.
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
 * Nota pozycji samodzielnej — zdanie o tym, czym Terminal jest; nie jest ostrzeżeniem o usterce
 * i nic nie odbiera.
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
      // Okno własne, nie pierwsze z brzegu — inaczej trafiłoby do cudzego modułu.
      KOD_MODULU,
    ),
};

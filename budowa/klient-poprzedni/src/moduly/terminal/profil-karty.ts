import type { TerminalSession } from '../../../../shared/contract';
import { TerminalShell as Powloka } from '../../../../shared/contract';
import type { StanTerminala } from './stan-terminala';
import type { PozycjaWyboru } from './wybor-drzewem';

/**
 * Profil karty powłoki: wykaz powłok kontraktu, schematy barw i opis karty bieżącej.
 */

/**
 * Powłoki karty terminala w kolejności zgodnej z wykazem modułów uruchomieniowych dostępnych w tej budowie.
 */
export const POWLOKI: readonly PozycjaWyboru[] = [
  [Powloka.Powershell, 'PowerShell', 'Powłoka Windows z poleceniami cmdlet.'],
  [Powloka.Cmd, 'CMD', 'Wiersz poleceń Windows — składnia wsadowa.'],
  [Powloka.Bash, 'Bash', 'Powłoka uniksowa; na Windowsie wymaga zainstalowanej.'],
  [Powloka.Node, 'Node.js REPL', 'Pętla wyliczająca JavaScript, nie powłoka systemu.'],
  [Powloka.Python, 'Python REPL', 'Pętla wyliczająca Python, nie powłoka systemu.'],
  [Powloka.Ssh, 'SSH (powłoka zdalna)', 'Powłoka na maszynie zdalnej; katalog roboczy jest jej katalogiem.'],
  [
    Powloka.Container,
    'Kontener',
    'Polecenie wykonuje się w kontenerze, który już biegnie. Kartę wskazuje się identyfikatorem kontenera.',
  ],
  [
    Powloka.Pod,
    'Pod (Kubernetes)',
    'Polecenie wykonuje się w podzie klastra. Kontekst i przestrzeń nazw bierze konfiguracja maszyny rdzenia, o ile karta nie wskaże własnych.',
  ],
  [
    Powloka.Serial,
    'Konsola szeregowa',
    'Konsola urządzenia podłączonego do maszyny RDZENIA; wskazuje się ją ścieżką urządzenia i prędkością transmisji.',
  ],
  [
    Powloka.Telnet,
    'Telnet',
    'Sesja do urządzenia sieciowego. Telnet nie szyfruje ruchu — do maszyn z SSH właściwa jest karta zdalna.',
  ],
];

/**
 * Schematy barw karty pochodzą z żetonów motywu aplikacji, nie z własnych, niezależnych wartości barw.
 */
export const SCHEMATY: readonly PozycjaWyboru[] = [
  ['atrament', 'Atrament', 'Ciemne tło karty. Czynność widoku — rdzeń o niej nie wie.'],
  ['pergamin', 'Pergamin', 'Jasne tło karty. Czynność widoku — rdzeń o niej nie wie.'],
  ['kontrast', 'Wysoki kontrast', 'Największa różnica tła i pisma. Czynność widoku — rdzeń o niej nie wie.'],
];

/**
 * Inicjator uruchomienia rozstrzyga, kto odpowiada za polecenie; odmowę uprawnień wydaje zawsze rdzeń.
 */
export const INICJATORZY_KARTY: readonly PozycjaWyboru[] = [
  ['operator', 'Operator', 'Polecenie idzie na odpowiedzialność Operatora tego okna.'],
  [
    'model',
    'AI (model)',
    'Polecenie zapisuje się w rejestrze jako uruchomione przez model. W oknie o trybie ręcznym rdzeń odmawia (permission_denied) — kontrakt nie ma dla terminala rundy zgody.',
  ],
];

/**
 * Opis karty bieżącej: powłoka, katalog roboczy, stan bieżącej sesji i ostatnie wykonane polecenie terminala.
 */
export function opisKarty(karta: TerminalSession, stan: StanTerminala): HTMLElement {
  const opis = document.createElement('dl');
  opis.className = 'dt-opis-karty';
  const pozycje: ReadonlyArray<[string, string]> = [
    ['Powłoka', karta.shell],
    ['Katalog roboczy', karta.workingDir ?? 'własny katalog okna'],
    ['Stan karty', karta.status],
    ['Przypięta', stan.czyPrzypieta(karta.id) ? 'tak' : 'nie'],
    ['Ostatnie polecenie', stan.ostatniePolecenie(karta.id) === '' ? '—' : stan.ostatniePolecenie(karta.id)],
  ];
  for (const [nazwa, wartosc] of pozycje) {
    const podpis = document.createElement('dt');
    podpis.textContent = nazwa;
    const wpis = document.createElement('dd');
    wpis.textContent = wartosc;
    opis.append(podpis, wpis);
  }
  return opis;
}

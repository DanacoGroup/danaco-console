import type { DeveloperBuild, ProblemSeverity } from '../../../../shared/contract';
import { przyciskAkcji } from '../../modele/kontrolki-formularza-braki';

/** Widok przebiegu budowania — czysta konstrukcja z danych oddanych przez rdzeń, bez własnego stanu. */

/** Zawężenie widoku: fraza w logu i waga zgłoszeń, obie czynności wyłącznie kliencka, nad materiałem gotowym. */
export interface ZawezeniePrzebiegu {
  /** Fraza szukana w wierszach logu; pusta znaczy „bez zawężenia”. */
  szukaj: string;
  /** Waga zgłoszeń wpuszczana do wykazu; pusta znaczy „wszystkie wagi”. */
  waga: ProblemSeverity | '';
}

/** Materiał widoku przebiegu wraz z zawężeniem i przejściem do pliku zgłoszenia we wspólnym stanie modułu. */
export interface OpisWidokuPrzebiegu {
  logi: readonly string[];
  /** Czy okno nie widziało początku przebiegu. */
  uciety: boolean;
  zawezenie: ZawezeniePrzebiegu;
  naWskazanie: (sciezka: string) => void;
}

/** Zdanie o przebiegu: zadanie, stan, kod wyjścia, czas trwania i liczba zgłoszeń podana operatorowi w panelu. */
function opisPrzebiegu(przebieg: DeveloperBuild): string {
  const kod = przebieg.exitCode === undefined ? '' : `; kod wyjścia ${przebieg.exitCode}`;
  const zgloszenia = przebieg.problems?.length ?? 0;
  const czas =
    przebieg.finishedAt === undefined
      ? 'w toku'
      : `zakończono o ${new Date(przebieg.finishedAt).toLocaleTimeString('pl-PL')}`;
  return `${przebieg.task} — ${przebieg.status}${kod}; ${czas}; zgłoszeń: ${zgloszenia}`;
}

/**
 * Rysuje panel przebiegu: nagłówek stanu, log narastający i zgłoszenia.
 *
 * Zdanie o ucięciu mówi wyłącznie o tym oknie — czego ono nie widziało i czym
 * tego nie dogoni. O tym, czym rozporządza rdzeń, mówi stan pusty okna złożony
 * z jego rejestru komend.
 */
export function rysujPrzebieg(przebieg: DeveloperBuild, opis: OpisWidokuPrzebiegu): HTMLElement {
  const panel = document.createElement('div');
  panel.className = 'mdev-wiersz';

  const naglowek = document.createElement('p');
  naglowek.className = 'mdev-przebieg';
  naglowek.dataset['stanBudowania'] = przebieg.status;
  naglowek.textContent = opisPrzebiegu(przebieg);
  panel.append(naglowek);

  if (opis.uciety) {
    const ostrzezenie = document.createElement('p');
    ostrzezenie.className = 'mdev-potwierdzenie';
    ostrzezenie.dataset['udane'] = 'false';
    ostrzezenie.textContent =
      'Log ucięty: to okno otworzyło się w trakcie przebiegu (albo przebieg uruchomiło inne okno ' +
      'tego konta) — początek logu nie przyszedł zdarzeniami, a tego okna nie ma czym dogonić.';
    panel.append(ostrzezenie);
  }

  panel.append(...rysujLog(opis));

  if (przebieg.problems !== undefined && przebieg.problems.length > 0) {
    panel.append(...rysujCzescZgloszen(przebieg.problems, opis));
  }
  return panel;
}

/**
 * Log wraz ze zdaniem o zawężeniu.
 *
 * Zdanie o liczbie wierszy pominiętych stoi obok logu zawsze, gdy fraza jest
 * wpisana: log przycięty bez takiego zdania czytałby się jak log krótki, a to
 * dwie różne rzeczy.
 */
function rysujLog(opis: OpisWidokuPrzebiegu): readonly HTMLElement[] {
  const fraza = opis.zawezenie.szukaj.trim().toLocaleLowerCase('pl-PL');
  const widoczne =
    fraza === ''
      ? opis.logi
      : opis.logi.filter((wiersz) => wiersz.toLocaleLowerCase('pl-PL').includes(fraza));

  const czesci: HTMLElement[] = [];
  if (fraza !== '') {
    const zdanie = document.createElement('p');
    zdanie.className = 'mdev-zawezenie';
    zdanie.textContent =
      `Log zawężony frazą „${opis.zawezenie.szukaj.trim()}”: ${widoczne.length} z ${opis.logi.length} ` +
      'wierszy zebranych w tym oknie. Szukanie obejmuje wyłącznie wiersze, które okno usłyszało — ' +
      'kontrakt nie ma komendy, którą dałoby się dopytać rdzeń o resztę logu.';
    czesci.push(zdanie);
  }

  const log = document.createElement('pre');
  log.className = 'mdev-log';
  log.textContent = tresc(opis.logi.length, widoczne, fraza !== '');
  czesci.push(log);
  return czesci;
}

/** Treść pola logu — pustka bez zawężenia i pustka po zawężeniu to dwa różne zdania panelu przebiegu okna. */
function tresc(zebrane: number, widoczne: readonly string[], zawezone: boolean): string {
  if (widoczne.length > 0) return widoczne.join('\n');
  if (zebrane === 0) return 'Log jeszcze nie przyniósł żadnego wiersza.';
  if (zawezone) return 'Żaden z zebranych wierszy logu nie zawiera szukanej frazy.';
  return 'Log jeszcze nie przyniósł żadnego wiersza.';
}

/** Zgłoszenia wraz ze zdaniem o zawężeniu wagą, wyświetlane operatorowi w panelu przebiegu budowania okna. */
function rysujCzescZgloszen(
  problems: NonNullable<DeveloperBuild['problems']>,
  opis: OpisWidokuPrzebiegu,
): readonly HTMLElement[] {
  const waga = opis.zawezenie.waga;
  const widoczne = waga === '' ? problems : problems.filter((pozycja) => pozycja.severity === waga);

  const czesci: HTMLElement[] = [];
  if (waga !== '') {
    const zdanie = document.createElement('p');
    zdanie.className = 'mdev-zawezenie';
    zdanie.textContent =
      `Zgłoszenia zawężone do wagi „${waga}”: ${widoczne.length} z ${problems.length}. ` +
      'Waga pochodzi z pola zgłoszenia oddanego przez rdzeń, nie z rozpoznawania treści wiersza logu.';
    czesci.push(zdanie);
  }

  if (widoczne.length === 0) {
    const puste = document.createElement('p');
    puste.className = 'mdev-zawezenie';
    puste.textContent = `Przebieg nie ma zgłoszeń o wadze „${waga}”.`;
    czesci.push(puste);
    return czesci;
  }

  czesci.push(rysujZgloszenia(widoczne, opis.naWskazanie));
  return czesci;
}

/**
 * Wykaz zgłoszeń budowania — czysta konstrukcja z danych. Zgłoszenie ze
 * ścieżką jest przejściem do pliku, bez otwierania miejsca w nim.
 */
function rysujZgloszenia(
  problems: NonNullable<DeveloperBuild['problems']>,
  naWskazanie: (sciezka: string) => void,
): HTMLElement {
  const wykaz = document.createElement('ul');
  wykaz.className = 'mdev-wykaz';
  wykaz.setAttribute('aria-label', 'Zgłoszenia przebiegu budowania');
  for (const zgloszenie of problems) {
    const pozycja = document.createElement('li');
    pozycja.className = 'mdev-pozycja';
    pozycja.dataset['waga'] = zgloszenie.severity;
    const tytul = document.createElement('strong');
    tytul.className = 'mdev-pozycja__tytul';
    tytul.textContent =
      zgloszenie.path === undefined
        ? zgloszenie.severity
        : `${zgloszenie.path}${zgloszenie.line === undefined ? '' : `:${zgloszenie.line}`}`;
    const opis = document.createElement('span');
    opis.className = 'mdev-pozycja__opis';
    opis.textContent = zgloszenie.message;
    const akcje = document.createElement('span');
    akcje.className = 'mdev-pozycja__akcje';
    if (zgloszenie.path !== undefined && zgloszenie.path !== '') {
      akcje.append(przejscieDoPliku(zgloszenie.path, naWskazanie));
    }
    pozycja.append(tytul, opis, akcje);
    wykaz.append(pozycja);
  }
  return wykaz;
}

/** Przycisk „Otwórz w edytorze” — wskazuje plik zgłoszenia we wspólnym stanie modułu Developer tego okna. */
function przejscieDoPliku(
  sciezka: string,
  naWskazanie: (sciezka: string) => void,
): HTMLButtonElement {
  const otworz = przyciskAkcji('Otwórz w edytorze', 'dn-btn dn-btn--zarys');
  otworz.dataset['sciezka'] = sciezka;
  otworz.addEventListener('click', () => naWskazanie(sciezka));
  return otworz;
}

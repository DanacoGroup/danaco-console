import { TerminalShell } from '../../../../shared/contract';

/**
 * Programy zewnętrzne, bez których czynności modułu Terminal nie mają czym się
 * wykonać.
 *
 * Instalka Danaco Console nie niesie ani jednego z nich. Terminal nie ma
 * własnego interpretera powłoki, własnego klienta SSH ani własnego lintera —
 * uruchamia programy leżące na maszynie, na której stoi rdzeń serwera, i to one
 * wykonują pracę. Maszyna dewelopera ma je zwykle doinstalowane ręcznie, więc
 * pomiar zdolności robiony na niej zawyża to, co dostanie odbiorca instalki.
 *
 * Milczenie o tym byłoby brakiem funkcji, nie oszczędnością słowa: Operator,
 * który klika „Uruchom" i dostaje odmowę startu procesu bez wskazania programu,
 * nie ma jak odróżnić braku narzędzia na serwerze od usterki platformy. Dlatego
 * każde okno, którego czynność zależy od programu zewnętrznego, nazywa ten
 * program wprost i mówi, co się stanie, gdy go nie będzie.
 *
 * Wykaz powłok jest odwzorowaniem wykazu wykonawczego rdzenia
 * (`adapter_modul_terminal_powloki.go` i `..._powloki_urzadzen.go`), a nie
 * domysłem: rdzeń startuje dokładnie te programy z dokładnie tymi argumentami.
 */

/** Jeden program zewnętrzny wraz z jego rolą w module. */
export interface ProgramZewnetrzny {
  /** Nazwa pliku wykonywalnego szukanego na ścieżce PATH serwera. */
  program: string;
  /** Czynność modułu, która bez tego programu nie ma czym się wykonać. */
  potrzebnyDo: string;
  /** Co Operator zobaczy, gdy programu nie ma na serwerze. */
  skutekBraku: string;
}

/**
 * Program uruchamiany przez rdzeń dla rodzaju powłoki karty.
 *
 * Wykaz jest niepełny wobec słownika powłok kontraktu i taki ma być. Kontrakt
 * zna dziś dziesięć rodzajów powłoki, a wykaz wykonawczy rdzenia
 * (`adapter_modul_terminal_powloki.go`) — sześć. Dopisanie tu programu dla
 * powłoki, której rdzeń nie uruchamia, byłoby domysłem podanym Operatorowi jako
 * zależność zmierzona: nie wiadomo ani jaki program rdzeń wybierze, ani z jakimi
 * argumentami. Powłoka spoza wykazu dostaje więc zdanie o braku po stronie
 * rdzenia, a nie nazwę programu wziętą z wyobrażenia.
 *
 * Argumenty podane w opisie mają znaczenie dla Operatora: powłoka startuje bez
 * czytania profilu użytkownika, a każde polecenie jest odrębnym procesem tej
 * powłoki. Stąd wynika zachowanie, które inaczej wygląda jak usterka — `cd`
 * wydane jednym poleceniem nie przesuwa katalogu następnemu.
 */
const PROGRAM_POWLOKI: Readonly<Partial<Record<TerminalShell, ProgramZewnetrzny>>> = {
  [TerminalShell.Powershell]: {
    program: 'powershell',
    potrzebnyDo: 'karta powłoki PowerShell (uruchamiana z -NoProfile -NonInteractive -Command)',
    skutekBraku: 'Karta otworzy się, ale każde polecenie zakończy się odmową startu procesu.',
  },
  [TerminalShell.Cmd]: {
    program: 'cmd',
    potrzebnyDo: 'karta wiersza poleceń Windows (uruchamiana z /D /C)',
    skutekBraku: 'Poza Windowsem programu nie ma wcale — karta CMD nie wykona ani jednego polecenia.',
  },
  [TerminalShell.Bash]: {
    program: 'bash',
    potrzebnyDo: 'karta powłoki Bash (uruchamiana z -c)',
    skutekBraku: 'Na Windowsie bez zainstalowanego Basha karta nie wykona ani jednego polecenia.',
  },
  [TerminalShell.Node]: {
    program: 'node',
    potrzebnyDo: 'karta pętli wyliczającej Node.js (uruchamiana z -e)',
    skutekBraku: 'Bez środowiska Node.js na serwerze karta nie wykona ani jednego wyrażenia.',
  },
  [TerminalShell.Python]: {
    program: 'python',
    potrzebnyDo: 'karta pętli wyliczającej Python (uruchamiana z -u -c)',
    skutekBraku:
      'Instalacje, w których program nosi nazwę python3, nie są widziane pod nazwą python — karta odmówi startu.',
  },
  [TerminalShell.Ssh]: {
    program: 'ssh',
    potrzebnyDo: 'karta powłoki zdalnej (uruchamiana z -T i adresem celu)',
    skutekBraku:
      'Bez klienta OpenSSH na serwerze żadna karta zdalna nie ruszy, a książka hostów nie ma czym otworzyć połączenia.',
  },
  [TerminalShell.Container]: {
    program: 'docker',
    potrzebnyDo: 'karta powłoki wewnątrz kontenera (uruchamiana z exec -i i powłoką sh -c)',
    skutekBraku:
      'Bez klienta Dockera na serwerze karta kontenera nie wykona ani jednego polecenia; ' +
      'karta nie zarządza kontenerem — wykonuje polecenie w kontenerze, który już biegnie.',
  },
  [TerminalShell.Pod]: {
    program: 'kubectl',
    potrzebnyDo: 'karta powłoki wewnątrz poda (uruchamiana z exec -i i powłoką sh -c)',
    skutekBraku:
      'Bez kubectl na serwerze karta poda nie ruszy. Kontekst klastra i przestrzeń nazw bierze ' +
      'kubectl z konfiguracji maszyny rdzenia, o ile karta nie wskaże własnych.',
  },
  [TerminalShell.Serial]: {
    program: 'picocom',
    potrzebnyDo: 'karta konsoli portu szeregowego (uruchamiana z --baud i --initstring)',
    skutekBraku:
      'Bez programu picocom na serwerze karta portu szeregowego nie ruszy. Urządzenie musi leżeć ' +
      'na maszynie RDZENIA i konto procesu rdzenia musi mieć do niego prawo.',
  },
  [TerminalShell.Telnet]: {
    program: 'telnet',
    potrzebnyDo: 'karta sesji Telnet do urządzenia sieciowego (adres i port argumentem)',
    skutekBraku:
      'Bez klienta Telnet na serwerze karta nie ruszy. Telnet nie szyfruje ruchu — do maszyn, ' +
      'które mają SSH, właściwą kartą jest karta zdalna.',
  },
};

/**
 * Zależność powłoki albo pusto, gdy rdzeń jej nie uruchamia.
 *
 * Typ zostaje CZĘŚCIOWY (`Partial`) mimo że wykaz pokrywa dziś komplet słownika
 * kontraktu: kolejna wartość `TerminalShell` pojawi się w kontrakcie wcześniej
 * niż w wykazie wykonawczym rdzenia, a wtedy pusto ma być odpowiedzią, nie
 * wpisem zgadniętym. Nazwa programu wpisana tu bez pokrycia w rdzeniu byłaby
 * domysłem podanym Operatorowi jako zależność zmierzona.
 */
export function programPowloki(powloka: TerminalShell): ProgramZewnetrzny | null {
  return PROGRAM_POWLOKI[powloka] ?? null;
}

/** Zależności powłok, które rdzeń naprawdę uruchamia — treść noty okna kart. */
export const PROGRAMY_POWLOK: readonly ProgramZewnetrzny[] = Object.values(PROGRAM_POWLOKI).filter(
  (zaleznosc): zaleznosc is ProgramZewnetrzny => zaleznosc !== undefined,
);

/**
 * Programy zewnętrzne, po które sięgają czynności okien poza samym startem
 * powłoki. Wykaz obejmuje wyłącznie to, co moduł naprawdę wywołuje — narzędzia
 * wymienione w dokumentacji, a przez moduł niewywoływane, stoją w katalogu
 * funkcji (`katalog-funkcji.ts`), nie tutaj.
 */
export const PROGRAMY_CZYNNOSCI: readonly ProgramZewnetrzny[] = [
  {
    program: 'npm',
    potrzebnyDo: 'uruchomienie zadania wykrytego w pliku package.json',
    skutekBraku: 'Zadanie zakończy się niezerowym kodem wyjścia z komunikatem powłoki o braku programu.',
  },
  {
    program: 'make',
    potrzebnyDo: 'uruchomienie celu wykrytego w pliku Makefile',
    skutekBraku: 'Zadanie zakończy się niezerowym kodem wyjścia z komunikatem powłoki o braku programu.',
  },
  {
    program: 'task',
    potrzebnyDo: 'uruchomienie zadania wykrytego w pliku Taskfile.yml',
    skutekBraku: 'Zadanie zakończy się niezerowym kodem wyjścia z komunikatem powłoki o braku programu.',
  },
  {
    program: 'just',
    potrzebnyDo: 'uruchomienie przepisu wykrytego w pliku justfile',
    skutekBraku: 'Zadanie zakończy się niezerowym kodem wyjścia z komunikatem powłoki o braku programu.',
  },
];

/**
 * Zdanie o programie powłoki — wchodzi w podpowiedź kontrolki otwierającej
 * kartę i w potwierdzenie czynności.
 */
export function zdanieProgramuPowloki(powloka: TerminalShell): string {
  const zaleznosc = programPowloki(powloka);
  if (zaleznosc === null) {
    return (
      `Powłoka ${powloka} stoi w słowniku kontraktu, ale rdzeń nie ma dla niej programu w swoim wykazie ` +
      'wykonawczym — karta tej powłoki nie ruszy, dopóki rdzeń go nie dostanie. Jaki to będzie program ' +
      'i z jakimi argumentami, rozstrzygnie rdzeń; okno tego nie zgaduje.'
    );
  }
  return (
    `Karta ${powloka} uruchamia na serwerze program ${zaleznosc.program}. ` +
    `Instalka Danaco Console go nie niesie. ${zaleznosc.skutekBraku}`
  );
}

/**
 * Akapit o zależnościach zewnętrznych osadzany w ciele okna.
 *
 * Stoi w widoku na stałe, a nie w dymku pod znakiem zapytania: to jest warunek
 * działania okna, a nie objaśnienie pojedynczej kontrolki.
 */
export function notaZaleznosci(
  wstep: string,
  programy: readonly ProgramZewnetrzny[],
): HTMLElement {
  const blok = document.createElement('section');
  blok.className = 'dt-zaleznosci';

  const podpis = document.createElement('h4');
  podpis.className = 'dt-zaleznosci__podpis';
  podpis.textContent = 'Programy zewnętrzne, których wymaga to okno';
  blok.append(podpis);

  const zdanie = document.createElement('p');
  zdanie.className = 'dn-pole-opis';
  zdanie.textContent = wstep;
  blok.append(zdanie);

  const lista = document.createElement('dl');
  lista.className = 'dt-zaleznosci__wykaz';
  for (const zaleznosc of programy) {
    const nazwa = document.createElement('dt');
    nazwa.textContent = zaleznosc.program;
    const opis = document.createElement('dd');
    opis.textContent = `${zaleznosc.potrzebnyDo} — ${zaleznosc.skutekBraku}`;
    lista.append(nazwa, opis);
  }
  blok.append(lista);
  return blok;
}

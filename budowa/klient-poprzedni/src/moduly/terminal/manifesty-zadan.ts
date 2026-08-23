import { TerminalShell } from '../../../../shared/contract';
import type { PozycjaWyboru } from './wybor-drzewem';

/**
 * Zadania projektu wyczytane z manifestu leżącego w katalogu roboczym karty.
 *
 * Treść manifestu czyta rdzeń komendą `terminal.file.read`, a nie polecenie
 * powłoki: dzięki temu wykrycie zadań nie zależy ani od programu wypisującego
 * plik, ani od składni powłoki karty, i działa jednakowo w każdej z nich.
 * Przeglądania katalogu kontrakt nie ma wcale, więc nazwa manifestu pochodzi
 * z zamkniętego wykazu poniżej — okno nie zgaduje, co leży w katalogu.
 * Gdy manifestu nie ma, rdzeń odmawia kodem `not_found`, a okno pokazuje tę
 * odmowę dosłownie zamiast pustego wykazu zadań.
 *
 * Rozbiór jest czytaniem, nie wykonaniem: treść manifestu nigdy nie trafia do
 * powłoki z powrotem. Do powłoki idzie wyłącznie polecenie złożone z nazwy
 * programu i nazwy zadania odczytanej z manifestu.
 *
 * Nazwy plików pochodzą z zamkniętego wykazu poniżej, nie z pola wpisywania —
 * dlatego wchodzą w treść polecenia bez cudzysłowu obliczanego w czasie
 * działania.
 */

/** Manifest, z którego czytane są zadania projektu. */
export type RodzajManifestu = 'package.json' | 'Makefile' | 'Taskfile.yml' | 'justfile';

/** Wykaz manifestów wraz z programem, który uruchamia ich zadania. */
export const MANIFESTY: readonly PozycjaWyboru[] = [
  ['package.json', 'package.json', 'Pole scripts manifestu Node.js; zadania uruchamia program npm.'],
  ['Makefile', 'Makefile', 'Cele pliku Makefile; zadania uruchamia program make.'],
  ['Taskfile.yml', 'Taskfile.yml', 'Blok tasks pliku Taskfile; zadania uruchamia program task.'],
  ['justfile', 'justfile', 'Przepisy pliku justfile; zadania uruchamia program just.'],
];

/** Jedno zadanie odczytane z manifestu. */
export interface ZadanieManifestu {
  /** Nazwa zadania w manifeście. */
  nazwa: string;
  /** Polecenie uruchamiające zadanie w powłoce karty. */
  polecenie: string;
  zrodlo: RodzajManifestu;
}

/**
 * Polecenie wypisujące treść manifestu w danej powłoce.
 *
 * Każda powłoka wykazu kontraktu ma swoją drogę, także obie pętle wyliczające:
 * Node.js i Python czytają plik własną biblioteką standardową, więc nie
 * potrzebują do tego żadnego programu zewnętrznego. Powłoki systemowe sięgają
 * po program wypisujący plik i to jest ich zależność zewnętrzna.
 */
export function poleceniePodgladuManifestu(powloka: TerminalShell, plik: RodzajManifestu): string {
  switch (powloka) {
    case TerminalShell.Bash:
    case TerminalShell.Ssh:
      return `cat -- ${plik}`;
    case TerminalShell.Powershell:
      return `Get-Content -Raw -LiteralPath '${plik}'`;
    case TerminalShell.Cmd:
      return `type ${plik}`;
    case TerminalShell.Node:
      return `process.stdout.write(require('fs').readFileSync('${plik}','utf8'))`;
    case TerminalShell.Python:
      return `import sys; sys.stdout.write(open('${plik}', encoding='utf-8').read())`;
    default:
      // Powłoka spoza wykazu kontraktu: okno powie o braku drogi zamiast wysłać
      // polecenie złożone dla innej składni.
      return '';
  }
}

/** Rozbiór treści manifestu na zadania; treść nierozpoznana daje pusty wykaz. */
export function czytajZadania(rodzaj: RodzajManifestu, tresc: string): ZadanieManifestu[] {
  switch (rodzaj) {
    case 'package.json':
      return zadaniaManifestuNode(tresc);
    case 'Makefile':
      return celeMakefile(tresc);
    case 'Taskfile.yml':
      return zadaniaTaskfile(tresc);
    default:
      return przepisyJustfile(tresc);
  }
}

/**
 * Pole `scripts` manifestu Node.js.
 *
 * Treść niebędąca poprawnym JSON-em nie jest błędem programu, tylko odpowiedzią
 * powłoki — plik bywa nieobecny, a wtedy w wyjściu stoi komunikat błędu. Rozbiór
 * oddaje wtedy pusty wykaz, a zdanie o powodzie składa okno z treści wyjścia.
 */
function zadaniaManifestuNode(tresc: string): ZadanieManifestu[] {
  let odczytane: unknown;
  try {
    odczytane = JSON.parse(tresc);
  } catch {
    return [];
  }
  if (typeof odczytane !== 'object' || odczytane === null) return [];
  const skrypty = (odczytane as { scripts?: unknown }).scripts;
  if (typeof skrypty !== 'object' || skrypty === null) return [];
  return Object.keys(skrypty).map((nazwa) => ({
    nazwa,
    polecenie: `npm run ${nazwa}`,
    zrodlo: 'package.json' as const,
  }));
}

/**
 * Cele pliku Makefile.
 *
 * Brany jest wyłącznie wiersz zaczynający się od nazwy celu w pierwszej
 * kolumnie. Wiersz wcięty jest przepisem celu, a nie celem; nazwa zaczynająca
 * się kropką (`.PHONY`, `.SUFFIXES`) jest dyrektywą programu make; zapis
 * `nazwa :=` jest przypisaniem zmiennej.
 */
function celeMakefile(tresc: string): ZadanieManifestu[] {
  const zadania: ZadanieManifestu[] = [];
  const znane = new Set<string>();
  for (const wiersz of tresc.split('\n')) {
    const trafienie = /^([A-Za-z0-9][A-Za-z0-9._-]*)\s*:(?![=:])/.exec(wiersz);
    const nazwa = trafienie?.[1];
    if (nazwa === undefined || znane.has(nazwa)) continue;
    znane.add(nazwa);
    zadania.push({ nazwa, polecenie: `make ${nazwa}`, zrodlo: 'Makefile' });
  }
  return zadania;
}

/**
 * Zadania bloku `tasks:` pliku Taskfile.
 *
 * Czytany jest jeden poziom zagnieżdżenia: klucze wcięte bezpośrednio pod
 * `tasks:`. Głębsze klucze są polami zadania (`cmds`, `desc`), a nie zadaniami,
 * więc wejście na nie dawałoby wykaz nazw, których program task nie zna.
 */
function zadaniaTaskfile(tresc: string): ZadanieManifestu[] {
  const zadania: ZadanieManifestu[] = [];
  let wBloku = false;
  let wciecie = 0;
  for (const wiersz of tresc.split('\n')) {
    if (/^tasks\s*:/.test(wiersz)) {
      wBloku = true;
      wciecie = 0;
      continue;
    }
    if (!wBloku) continue;
    if (wiersz.trim() === '' || wiersz.trim().startsWith('#')) continue;
    const trafienie = /^(\s+)([A-Za-z0-9][A-Za-z0-9._:-]*)\s*:/.exec(wiersz);
    if (trafienie === null) {
      // Wiersz bez wcięcia kończy blok zadań — zaczyna się kolejny klucz pliku.
      if (!/^\s/.test(wiersz)) wBloku = false;
      continue;
    }
    const dlugosc = (trafienie[1] ?? '').length;
    if (wciecie === 0) wciecie = dlugosc;
    if (dlugosc !== wciecie) continue;
    const nazwa = trafienie[2] ?? '';
    zadania.push({ nazwa, polecenie: `task ${nazwa}`, zrodlo: 'Taskfile.yml' });
  }
  return zadania;
}

/**
 * Przepisy pliku justfile.
 *
 * Przepis stoi w pierwszej kolumnie i bywa z parametrami (`wdroz srodowisko:`);
 * przypisanie zmiennej (`wersja := "1"`) przepisem nie jest.
 */
function przepisyJustfile(tresc: string): ZadanieManifestu[] {
  const zadania: ZadanieManifestu[] = [];
  const znane = new Set<string>();
  for (const wiersz of tresc.split('\n')) {
    if (/^\s/.test(wiersz) || wiersz.trim().startsWith('#')) continue;
    const trafienie = /^([A-Za-z0-9][A-Za-z0-9._-]*)[^:=]*:(?![=:])/.exec(wiersz);
    const nazwa = trafienie?.[1];
    if (nazwa === undefined || znane.has(nazwa)) continue;
    znane.add(nazwa);
    zadania.push({ nazwa, polecenie: `just ${nazwa}`, zrodlo: 'justfile' });
  }
  return zadania;
}

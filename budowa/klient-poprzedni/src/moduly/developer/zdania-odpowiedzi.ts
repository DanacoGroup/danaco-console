import {
  BuildStatus,
  GitActionKind,
  type DeveloperBuild,
  type DeveloperFile,
  type GitActionResult,
} from '../../../../shared/contract';

/**
 * Zdania potwierdzeń modułu Developer — składane wyłącznie z odpowiedzi rdzenia.
 *
 * Reguła, której ten plik pilnuje: potwierdzenie mówi, co zrobił rdzeń, a nie
 * co wysłało okno. Zdanie sukcesu biorące wartość z żądania jest usterką nawet
 * wtedy, gdy dziś przypadkiem się zgadza — rozjedzie się przy pierwszej zmianie
 * po stronie rdzenia i nikt tego nie zauważy.
 *
 * Odmowa nie jest tu składana: treść odmowy rdzenia niesie stan błędu okna
 * (`komponenty/stan-tresci.ts`), z kodem i wiadomością wprost.
 */

/** Potwierdzenie gotowe do podania oknu: zdanie wraz z jego wydźwiękiem. */
export interface Potwierdzenie {
  zdanie: string;
  udane: boolean;
}

/** Nazwy czynności repozytorium — jedno źródło dla przycisku i dla odpowiedzi. */
export const NAZWY_CZYNNOSCI: Readonly<Record<GitActionKind, string>> = {
  [GitActionKind.Stage]: 'Dodanie do indeksu',
  [GitActionKind.Unstage]: 'Wycofanie z indeksu',
  [GitActionKind.Commit]: 'Zatwierdzenie zmian',
  [GitActionKind.Amend]: 'Poprawienie ostatniego zatwierdzenia',
  [GitActionKind.Revert]: 'Odwrócenie zatwierdzenia',
  [GitActionKind.Checkout]: 'Przełączenie gałęzi',
  [GitActionKind.Merge]: 'Scalenie gałęzi',
  [GitActionKind.Rebase]: 'Przestawienie gałęzi',
  [GitActionKind.Tag]: 'Nadanie etykiety',
  [GitActionKind.Fetch]: 'Pobranie zmian zdalnych',
  [GitActionKind.Pull]: 'Pobranie i scalenie',
  [GitActionKind.Push]: 'Wysłanie zmian',
  [GitActionKind.Stash]: 'Odłożenie zmian',
  [GitActionKind.StashPop]: 'Przywrócenie odłożonych zmian',
};

/** Nazwa czynności oddanej przez rdzeń; nieznanej nie tłumaczymy na siłę. */
function nazwaCzynnosci(czynnosc: string): string {
  return NAZWY_CZYNNOSCI[czynnosc as GitActionKind] ?? `czynność ${czynnosc}`;
}

/**
 * Zdanie o czynności repozytorium — nazwa bierze się z pola `action` wyniku.
 *
 * Rozjazd żądanej czynności z wykonaną nie jest przemilczany: gdy rdzeń wykona
 * inną niż zamówiona, mówi o tym zdanie, a nie etykieta naciśniętego przycisku.
 */
export function zdanieCzynnosciRepozytorium(
  zadana: GitActionKind,
  wynik: GitActionResult,
): Potwierdzenie {
  const wykonana = nazwaCzynnosci(wynik.action);
  const rozjazd =
    wynik.action === zadana
      ? ''
      : ` (żądano innej czynności: ${nazwaCzynnosci(zadana)} — rdzeń odesłał ${wynik.action})`;
  const skutek = wynik.succeeded
    ? 'rdzeń melduje powodzenie'
    : 'repozytorium odmówiło — powód w wyniku poniżej';
  return { zdanie: `${wykonana}${rozjazd}: ${skutek}.`, udane: wynik.succeeded };
}

/**
 * Zdanie o zapisie pliku — wersja rozstrzygana zmianą `versionId`, nie
 * przełącznikiem okna.
 *
 * Dlaczego zmianą, a nie obecnością: zapis z `createVersion:false` też oddaje
 * `versionId`, bo rdzeń dokłada wtedy identyfikator wersji najnowszej,
 * założonej wcześniej. Okno pytające „czy pole jest obecne” meldowałoby nową
 * wersję przy każdym zapisie pliku, który kiedykolwiek jakąś miał. Pytanie
 * „czy identyfikator się zmienił” odróżnia wersję świeżo założoną od zastanej.
 */
export function zdanieZapisuPliku(
  plik: DeveloperFile,
  wersjaPrzedZapisem: string,
  zadanoWersji: boolean,
): Potwierdzenie {
  const wersjaPo = plik.versionId ?? '';
  const zalozona = wersjaPo !== '' && wersjaPo !== wersjaPrzedZapisem;
  const czolo = `Rdzeń zapisał plik „${plik.path}”`;
  const otresci =
    plik.content === undefined
      ? ' Rdzeń nie odesłał zapisanej treści, więc okno nie potwierdza, że na dysku leży dokładnie to, co w polu.'
      : '';
  if (zadanoWersji && !zalozona) {
    return {
      zdanie:
        `${czolo}, ale NIE oddał nowej wersji — punktu powrotu do treści sprzed zapisu nie ma.` +
        otresci,
      udane: false,
    };
  }
  if (zalozona) {
    const dopisek = zadanoWersji ? '' : ' — o wersję nie proszono, rdzeń założył ją sam';
    return { zdanie: `${czolo} i założył wersję ${wersjaPo}${dopisek}.${otresci}`, udane: true };
  }
  return { zdanie: `${czolo}. Wersji nie zakładano.${otresci}`, udane: true };
}

/** Zdanie o stanie przebiegu — stan, kod wyjścia i czas zakończenia z odpowiedzi. */
function opisStanuPrzebiegu(przebieg: DeveloperBuild): string {
  const kod = przebieg.exitCode === undefined ? '' : `, kod wyjścia ${przebieg.exitCode}`;
  if (przebieg.finishedAt === undefined) return `stan ${przebieg.status}${kod}`;
  const czas = new Date(przebieg.finishedAt).toLocaleTimeString('pl-PL');
  return `stan ${przebieg.status}${kod}, zakończony o ${czas}`;
}

/**
 * Zdanie o uruchomieniu budowania — zadanie i stan z odpowiedzi rdzenia.
 *
 * „Uruchomione” nie jest wpisane na stałe: rdzeń może oddać przebieg już
 * domknięty (wyścig odpowiedzi ze zdarzeniem, zadanie kończące się natychmiast),
 * a wtedy zdanie o uruchomieniu byłoby zdaniem o czymś, co się skończyło.
 */
export function zdanieUruchomienia(przebieg: DeveloperBuild): Potwierdzenie {
  const czolo = `Rdzeń przyjął zadanie „${przebieg.task}” (przebieg ${przebieg.id})`;
  if (przebieg.finishedAt === undefined) {
    return { zdanie: `${czolo}: ${opisStanuPrzebiegu(przebieg)}.`, udane: true };
  }
  return {
    zdanie: `${czolo}, ale oddaje go już domkniętego: ${opisStanuPrzebiegu(przebieg)}.`,
    udane: przebieg.status === BuildStatus.Succeeded,
  };
}

/**
 * Zdanie o przerwaniu budowania.
 *
 * `bylZakonczony` mówi, czy przebieg trzymany przez okno był domknięty już
 * przed naciśnięciem — tego rdzeń nie powie, bo na oba przypadki odpowiada tą
 * samą migawką. Bez tego rozróżnienia okno meldowałoby przyjęcie żądania także
 * wtedy, gdy przerywać nie było czego.
 */
export function zdaniePrzerwania(przebieg: DeveloperBuild, bylZakonczony: boolean): Potwierdzenie {
  const opis = `przebieg „${przebieg.task}” (${przebieg.id}): ${opisStanuPrzebiegu(przebieg)}`;
  if (bylZakonczony) {
    return {
      zdanie: `Przerywać nie było czego — ${opis}. Ten przebieg był domknięty już przed żądaniem.`,
      udane: false,
    };
  }
  if (przebieg.finishedAt === undefined) {
    return { zdanie: `Rdzeń przyjął żądanie, ale ${opis} — przebieg wciąż trwa.`, udane: false };
  }
  return { zdanie: `Rdzeń domknął ${opis}.`, udane: true };
}

import {
  Command,
  EventType,
  SessionToolSource,
  SlashEntryKind,
  type ToolCatalogEntry,
} from '../../../shared/contract';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import type { Kanal } from '../protokol/kanal';

/**
 * Wejście danych wykazu po ukośniku, wzorowane na
 * `okno-komunikacji/katalog-akcji.ts`.
 *
 * Typ funkcyjny oddający `(pozycje, powod)` jest w tym repozytorium wzorcem
 * zastanym: panel akcji modułu bierze pozycje tak samo i z tego samego powodu —
 * wykaz pusty bez powodu jest atrapą, wykaz pusty z powodem jest stanem
 * opisanym.
 *
 * Pozycją jest `ToolCatalogEntry` z kontraktu, a nie własny kształt: drugi model
 * dawałby dwie prawdy o tym, czym jest pozycja wykazu, i pierwszy rozjazd nazw
 * w rdzeniu przeszedłby tu niezauważony. Rozróżnienie powołania narzędzia od
 * komendy akcji niesie pole `kind` tej struktury, więc stoi w danych, a nie
 * w domyśle widoku.
 */
export type ZrodloWykazuUkosnika = (
  oddaj: (pozycje: readonly ToolCatalogEntry[], powod: string) => void,
) => void;

/**
 * Wykonanie wyboru pozycji rodzaju `tool` — dołożenie na czas sesji.
 *
 * Oddaje zdanie dla Operatora i to, czy dołożenie weszło. Zdanie idzie do wątku
 * rozmowy: model dostał narzędzie, którego nie miał, więc musi to być widoczne
 * — tak samo jak niepowodzenie, bo cisza po wyborze byłaby nieodróżnialna od
 * powodzenia.
 */
export type DolozenieNarzedzia = (
  pozycja: ToolCatalogEntry,
  oddaj: (zdanie: string, udane: boolean) => void,
) => void;

/**
 * Źródło wykazu czytane z rdzenia komendą `tools.catalog.list`.
 *
 * Wykaz idzie w komplecie, bez pola `query`. Kontrakt umie zawęzić wykaz po
 * stronie rdzenia, ale filtrowanie ma działać od pierwszego znaku, a runda do
 * rdzenia na każde uderzenie w klawisz tego nie daje. Zawężanie, porządek
 * według trafności i wytłuszczenie trafień niesie mechanizm z biblioteki
 * (`komponenty/menu-drzewo.ts`), któremu wykaz jest podawany w całości.
 *
 * `sessionId` idzie, gdy sesja stoi — wtedy rdzeń oznacza pozycje już dołożone
 * polem `attached` i Operator nie dokłada po raz drugi tego, co ma.
 */
export function zrodloWykazuKanalu(kanal: Kanal): ZrodloWykazuUkosnika {
  return (oddaj) => {
    const idSesji = kanal.sesja().id();
    kanal.wyslij(
      Command.ToolsCatalogList,
      idSesji === '' ? {} : { sessionId: idSesji },
      (wynik) => {
        if (!wynik.udany) {
          oddaj([], `Wykaz po ukośniku nieodczytany — ${opisBledu(wynik.blad)}`);
          return;
        }
        const pozycje = wynik.wynik?.entries ?? [];
        oddaj(
          pozycje,
          pozycje.length > 0
            ? ''
            : 'Rdzeń oddał wykaz po ukośniku bez ani jednej pozycji do wyboru.',
        );
      },
    );
  };
}

/**
 * Dołożenie narzędzia do sesji komendą `session.tool.attach`.
 *
 * Powtórzenie nie jest błędem: kontrakt oddaje `alreadyAttached` i to pole jest
 * tu czytane wprost, więc sięgnięcie po narzędzie już dołożone daje zdanie
 * o tym, a nie drugi wpis „dołożono".
 *
 * `source` idzie jawnie jako `slashCommand`, choć kontrakt przyjmuje tę wartość
 * także milczeniem. Komenda po ukośniku jest jedyną drogą poszerzenia zestawu
 * w trakcie pracy, więc zapis sprawcy ma stać w danych wprost.
 */
export function dolozenieKanalu(kanal: Kanal): DolozenieNarzedzia {
  return (pozycja, oddaj) => {
    const idSesji = kanal.sesja().id();
    if (idSesji === '') {
      oddaj(
        'Narzędzia nie dołożono: rdzeń nie założył jeszcze sesji, a dołożenie żyje ' +
          'w stanie sesji i nie ma dokąd wejść.',
        false,
      );
      return;
    }
    kanal.wyslij(
      Command.SessionToolAttach,
      {
        sessionId: idSesji,
        toolName: pozycja.name,
        source: SessionToolSource.SlashCommand,
      },
      (wynik) => {
        if (!wynik.udany) {
          oddaj(`Narzędzia „${pozycja.name}" nie dołożono — ${opisBledu(wynik.blad)}`, false);
          return;
        }
        // Zameldowane TU nie ma być zameldowane drugi raz przez rozgłoszone
        // zdarzenie `session.tool.attached`, które wróci także do tego okna.
        wlasneDolozenia.add(pozycja.name);
        oddaj(
          wynik.wynik?.alreadyAttached === true
            ? `Narzędzie „${pozycja.name}" było już dołożone do tej sesji — zestaw bez zmian.`
            : zdanieDolozenia(pozycja),
          true,
        );
      },
    );
  };
}

/**
 * Nasłuch zdarzenia `session.tool.attached` — widoczność dołożenia spoza tego
 * okna.
 *
 * Wybór z tego okna niesie własne zdanie z odpowiedzi komendy, ale zestaw
 * narzędzi sesji poszerza także asystent działający za Operatora
 * (`SessionToolSource.assistant`). Bez tego nasłuchu takie dołożenie byłoby
 * niewidoczne.
 */
export function nasluchDolozen(
  kanal: Kanal,
  zglos: (zdanie: string) => void,
): Odsubskrybuj {
  return kanal.naZdarzenie(EventType.SessionToolAttached, (tresc) => {
    const nazwa = tresc.tool.name;
    if (nazwa === '' || czyWlasneDolozenie(nazwa)) return;
    zglos(
      `Zestaw narzędzi sesji poszerzony poza tym oknem: model dostał „${nazwa}"` +
        `${opisemDo(tresc.tool.description)}`,
    );
  });
}

/* ------------------------------------------------------------------------- */

/**
 * Nazwy dołożeń zameldowanych już przez odpowiedź komendy z tego okna.
 *
 * Zbiór jest maleńki i sesyjny. Stoi tutaj, bo tylko ten plik widzi obie drogi
 * wejścia tej samej wiadomości: odpowiedź komendy i rozgłoszone zdarzenie.
 * Bez niego Operator, który dołożył narzędzie sam, przeczytałby o tym w wątku
 * dwa razy — a to jest ten sam rodzaj podwojenia, co dwa wpisy jednej tury.
 */
const wlasneDolozenia = new Set<string>();

function czyWlasneDolozenie(nazwa: string): boolean {
  if (!wlasneDolozenia.has(nazwa)) return false;
  wlasneDolozenia.delete(nazwa);
  return true;
}

/** Zdanie stawiane w wątku po udanym dołożeniu. */
function zdanieDolozenia(pozycja: ToolCatalogEntry): string {
  return (
    `Model dostał narzędzie, którego nie miał: „${pozycja.name}"` +
    `${opisemDo(pozycja.description)} Dołożenie trwa do końca tej sesji — definicja ` +
    'eksperta pozostaje nietknięta.'
  );
}

/** Dokleja opis pozycji jako zdanie poboczne; pusty opis nie dokleja nic. */
function opisemDo(opis: string): string {
  return opis === '' ? '.' : ` — ${opis}`;
}

/** Czy pozycja poszerza zestaw narzędzi modelu, czy wykonuje czynność. */
export function czyPowolanieNarzedzia(pozycja: ToolCatalogEntry): boolean {
  return pozycja.kind === SlashEntryKind.Tool;
}

/** Opis błędu kontraktu dla Operatora — jak w `katalog-akcji.ts`. */
function opisBledu(blad: { message: string; code: string } | undefined): string {
  if (blad === undefined) return 'rdzeń nie podał przyczyny';
  return `${blad.message} (${blad.code})`;
}

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
 * Wejście danych wykazu po ukośniku, wzorowane na module katalogu akcji i korzystające z
 * tego samego kształtu pozycji.
 */
export type ZrodloWykazuUkosnika = (
  oddaj: (pozycje: readonly ToolCatalogEntry[], powod: string) => void,
) => void;

/**
 * Wykonanie wyboru pozycji rodzaju narzędzia — dołożenie go na czas sesji, ze zdaniem dla
 * Operatora o skutku.
 */
export type DolozenieNarzedzia = (
  pozycja: ToolCatalogEntry,
  oddaj: (zdanie: string, udane: boolean) => void,
) => void;

/**
 * Źródło wykazu czytane z rdzenia komendą tools.catalog.list, w komplecie, bez zawężania
 * po stronie rdzenia.
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
 * Dołożenie narzędzia do sesji komendą session.tool.attach; powtórzenie nie jest błędem,
 * tylko odrębnym zdaniem.
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
        // Zameldowane tu dołożenie nie melduje się drugi raz przez rozgłoszone zdarzenie tej samej
        // sesji.
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
 * Nasłuch zdarzenia dołożenia narzędzia — widoczność dołożenia wykonanego spoza tego okna
 * przez asystenta.
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
 * Nazwy dołożeń zameldowanych już przez odpowiedź komendy z tego okna, trzymane w małym
 * zbiorze sesyjnym.
 */
const wlasneDolozenia = new Set<string>();

function czyWlasneDolozenie(nazwa: string): boolean {
  if (!wlasneDolozenia.has(nazwa)) return false;
  wlasneDolozenia.delete(nazwa);
  return true;
}

/**
 * Zdanie stawiane w wątku bieżącej rozmowy zaraz po udanym dołożeniu narzędzia do
 * trwającej sesji Operatora.
 */
function zdanieDolozenia(pozycja: ToolCatalogEntry): string {
  return (
    `Model dostał narzędzie, którego nie miał: „${pozycja.name}"` +
    `${opisemDo(pozycja.description)} Dołożenie trwa do końca tej sesji — definicja ` +
    'eksperta pozostaje nietknięta.'
  );
}

/**
 * Dokleja opis wybranej pozycji wykazu jako zdanie poboczne w wątku rozmowy; pusty opis
 * nie dokleja niczego.
 */
function opisemDo(opis: string): string {
  return opis === '' ? '.' : ` — ${opis}`;
}

/**
 * Czy wybrana pozycja wykazu poszerza zestaw narzędzi modelu, czy raczej wykonuje czynność
 * aplikacji.
 */
export function czyPowolanieNarzedzia(pozycja: ToolCatalogEntry): boolean {
  return pozycja.kind === SlashEntryKind.Tool;
}

/**
 * Opis błędu kontraktu przeznaczony dla Operatora, złożony z treści komunikatu i kodu
 * błędu rdzenia.
 */
function opisBledu(blad: { message: string; code: string } | undefined): string {
  if (blad === undefined) return 'rdzeń nie podał przyczyny';
  return `${blad.message} (${blad.code})`;
}

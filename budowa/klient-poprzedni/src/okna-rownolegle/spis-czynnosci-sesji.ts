import type { NazwaIkony } from '../ikony/ikony';
import { pokazKomunikat } from '../aplikacja/komunikaty';
import { otworzUsuniecieSesji } from '../powloka/potwierdzenie-usuniecia';
import { czyRozliczeniePuste, trescRozliczenia } from '../powloka/rozliczenie-usuniecia';
import { zadajUsuniecieSesji } from '../powloka/usuniecie-sesji';
import type { Kanal } from '../protokol/kanal';
import { czynnosciWiersza } from '../strona-glowna/czynnosci-sesji';
import { nazwa as nazwaSesji } from '../strona-glowna/meldunki-sesji';
import { utworzCzynnosciHistorii } from '../strona-glowna/wykonanie-czynnosci';
import type { WpisSesji } from '../strona-glowna/zrodlo-sesji';
import type { PozycjaCzynnosciMenu } from './wiersz-czynnosci';

/**
 * Które czynności sesji stają w menu `⋮` okna rozmowy — i czym są.
 *
 * O tym, czy czynność ma sens dla danego stanu sesji, rozstrzyga
 * `strona-glowna/czynnosci-sesji.ts` (`czynnosciWiersza`), nie ten moduł.
 * Powtórzenie tamtych warunków tutaj dałoby dwie odpowiedzi na jedno pytanie.
 * Tak samo z wykonaniem: komendy `session.*` stoją złożone
 * w `wykonanie-czynnosci.ts`; ten moduł nie pisze ani jednej nowej.
 *
 * Moduł robi trzy rzeczy, których tamten wykaz nie zna, bo nie zna menu:
 * wybiera z niego pozycje należące do sekcji sesji, dokłada „Usuń" (usuwanie
 * żyje w `powloka/`) oraz nadaje ikonę i literę skrótu.
 *
 * Rozgałęzienia rozmowy menu nie oferuje: nie ma komendy ani bytu, który by je
 * wykonał. Wiersz wygaszony byłby bramką, a wiersz czynny — atrapą.
 *
 * Każda czynność oddaje zdanie (`MeldunekCzynnosci`) i to zdanie trafia do
 * komunikatu: menu zwija się w chwili naciśnięcia, więc bez komunikatu nie
 * byłoby widać ani powodzenia, ani odmowy rdzenia.
 */

/** Zależności spisu: droga do rdzenia i ponowne odpytanie po zmianie. */
export interface ZaleznosciSpisu {
  kanal: Kanal;
  /** Ponowne odpytanie rdzenia o wpis sesji po udanej zmianie. */
  odswiez(): void;
}

/** Ikona i litera skrótu pozycji menu; klucz jest kluczem `czynnosciWiersza`. */
interface ZnakiPozycji {
  ikona: NazwaIkony;
  skrot: string;
}

/**
 * Pozycje sekcji sesji, w kolejności wyświetlania.
 *
 * Wykaz jest jednocześnie filtrem: czynność historii bez wpisu w tej mapie nie
 * wchodzi do menu okna rozmowy, choćby `czynnosciWiersza` uznał ją za sensowną.
 */
const ZNAKI_WZORCA: ReadonlyMap<string, ZnakiPozycji> = new Map([
  ['nazwa', { ikona: 'olowek', skrot: 'R' }],
  ['archiwum', { ikona: 'archiwum', skrot: 'A' }],
]);

/**
 * Składa spis czynności dla wpisu sesji tego okna.
 *
 * Pusty wynik jest odpowiedzią poprawną: bez wpisu sesji nie wiemy, w jakim
 * ona stanie, a czynność nazwana na ślepo obiecywałaby skutek, którego nikt nie
 * sprawdził.
 */
export function spisCzynnosciSesji(
  wpis: WpisSesji | null,
  zaleznosci: ZaleznosciSpisu,
): readonly PozycjaCzynnosciMenu[] {
  if (wpis === null) return [];

  const czynnosci = utworzCzynnosciHistorii(zaleznosci);
  const wykaz: PozycjaCzynnosciMenu[] = [];

  for (const pozycja of czynnosciWiersza(wpis, czynnosci)) {
    const znaki = ZNAKI_WZORCA.get(pozycja.klucz);
    if (znaki === undefined) continue;
    wykaz.push({
      klucz: pozycja.klucz,
      nazwa: pozycja.napis,
      przeznaczenie: pozycja.wyjasnienie,
      ikona: znaki.ikona,
      skrot: znaki.skrot,
      wykonaj: () => {
        void pozycja.wykonaj(wpis).then(zamelduj);
      },
    });
  }

  wykaz.push(pozycjaUsuniecia(wpis, zaleznosci));
  return wykaz;
}

/**
 * „Usuń" — jedyna droga utraty zapisu sesji w produkcie.
 *
 * Potwierdzenie z wykazem tego, co zginie, stoi
 * w `powloka/potwierdzenie-usuniecia.ts`: kontrakt żąda pola `confirm`, a rdzeń
 * bez niego odmawia (`adapter_sesje_usuwanie.go`). Menu tego potwierdzenia nie
 * dubluje własnym pytaniem — wołanie idzie prosto w tamto okno.
 *
 * Wiersz jest wyróżniony barwą (`grozna`), ale klikalny jak każdy inny.
 */
function pozycjaUsuniecia(
  wpis: WpisSesji,
  zaleznosci: ZaleznosciSpisu,
): PozycjaCzynnosciMenu {
  const tytul = nazwaSesji(wpis);
  const nazwaPo = (id: string): string => (id === wpis.sesja.id ? tytul : '');

  return {
    klucz: 'usun',
    nazwa: 'Usuń',
    przeznaczenie: 'Kasuje zapis sesji wraz z wiadomościami i oknami. Nie ma odwrotu.',
    ikona: 'kosz',
    skrot: 'D',
    grozna: true,
    wykonaj: () => {
      void otworzUsuniecieSesji(
        [{ id: wpis.sesja.id, tytul }],
        (idSesji) => zadajUsuniecieSesji(zaleznosci.kanal, idSesji, true),
        nazwaPo,
      ).then((rozliczenie) => {
        // `null` znaczy „brak potwierdzenia albo odmowa rdzenia" — powód został
        // już nazwany w oknie potwierdzenia, więc drugi komunikat o tej samej
        // odmowie byłby powtórzeniem.
        if (rozliczenie === null) return;
        zaleznosci.odswiez();
        pokazKomunikat({
          tytul: 'Trwałe usunięcie sesji',
          tresc: trescRozliczenia(rozliczenie, nazwaPo),
          waga: czyRozliczeniePuste(rozliczenie) ? 'blad' : 'info',
        });
      });
    },
  };
}

/** Pokazuje zdanie czynności jako komunikat; `null` znaczy „bez meldunku". */
function zamelduj(meldunek: string | null): void {
  if (meldunek === null) return;
  pokazKomunikat({ tytul: 'Sesja', tresc: meldunek, waga: 'info' });
}

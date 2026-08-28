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
 * Spis rozstrzyga, które czynności sesji stają w menu okna rozmowy i czym są, wybierając pozycje należące do sekcji sesji, dokładając usunięcie oraz nadając ikonę i literę skrótu, bez powtarzania warunków ani komend zapisanych gdzie indziej.
 */
export interface ZaleznosciSpisu {
  kanal: Kanal;
  /** Ponowne odpytanie rdzenia o wpis sesji po udanej zmianie. */
  odswiez(): void;
}

/** Ikona i litera skrótu każdej pozycji menu sesji, przypisane do tego samego klucza, którym woła się wiersz czynności. */
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
 * Usunięcie jest jedyną drogą utraty zapisu sesji w produkcie; wiersz jest wyróżniony barwą ostrzegawczą, ale klikalny jak każdy inny, a potwierdzenie z wykazem strat pokazuje osobne okno.
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
        // Wartość pusta znaczy brak potwierdzenia albo odmowę rdzenia, nazwaną już w oknie potwierdzenia.
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

/** Pokazuje zdanie czynności jako komunikat dla Operatora; wartość pusta znaczy brak meldunku do pokazania. */
function zamelduj(meldunek: string | null): void {
  if (meldunek === null) return;
  pokazKomunikat({ tytul: 'Sesja', tresc: meldunek, waga: 'info' });
}

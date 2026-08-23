import { otworzPowierzchnieAod } from '../aod/indeks';
import { otworzOknoDostepow } from '../dostepy/indeks';
import { otworzOknoKonfiguracji } from '../konfiguracja/indeks';
import { otworzOknoMobile } from '../mobile/indeks';
import { otworzOknoModeli } from '../modele/indeks';
import { otworzOknoPunktowIzolacji } from '../punkty-izolacji/indeks';
import { otworzOknoUstawien } from '../ustawienia/indeks';
import type { Kanal } from '../protokol/kanal';
import type { PozycjaUstawienia } from '../strona-glowna/indeks';
import { pokazKomunikat } from './komunikaty';

/**
 * Skutek wyboru pozycji listwy ustawień Centrum dowodzenia.
 *
 * Listwa sama niczego nie otwiera — zgłasza wybór, a rozdział wyboru na
 * czynności należy do tej warstwy. Dzięki temu dopisanie kolejnego ekranu
 * ustawień zmienia wyłącznie ten plik.
 *
 * Pozycja bez własnego widoku pozostaje klikalna, a naciśnięcie mówi wprost,
 * że ekran jeszcze nie powstał, zamiast otwierać atrapę.
 */
export function wykonajZamiarUstawien(pozycja: PozycjaUstawienia, kanal: Kanal): void {
  if (pozycja.kod === 'konfiguracja') {
    otworzOknoKonfiguracji(kanal);
    return;
  }

  if (pozycja.kod === 'dostepy') {
    otworzOknoDostepow(kanal);
    return;
  }

  if (pozycja.kod === 'modele') {
    otworzOknoModeli(kanal);
    return;
  }

  if (pozycja.kod === 'ustawienia') {
    otworzOknoUstawien(kanal);
    return;
  }

  if (pozycja.kod === 'punkty-izolacji') {
    otworzOknoPunktowIzolacji(kanal);
    return;
  }

  if (pozycja.kod === 'mobile') {
    otworzOknoMobile(kanal);
    return;
  }

  if (pozycja.kod === 'aod') {
    // Always On Display nie ma własnego okna: pozycja listwy otwiera
    // powierzchnię interakcji jako rozszerzenie boczne (opracowanie funkcji
    // globalnej, rozdz. 2.1 i 8.2).
    otworzPowierzchnieAod(kanal);
    return;
  }

  pokazKomunikat({
    tytul: pozycja.nazwa,
    tresc: `${pozycja.wyjasnienie}. Pozycja czeka na własny ekran — jeszcze go nie ma.`,
    waga: 'info',
  });
}

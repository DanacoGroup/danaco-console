/**
 * Zaczep pływającego Asystenta — jedyne miejsce, w którym favikon, dymek i stan
 * schodzą się w jedną warstwę. Warstwa osadza się w widoku środowiska, bo tylko
 * tam znana jest jednocześnie powłoka i droga do rdzenia.
 */

import './asystent-plywajacy.css';

import type { ZrodloPosuniec } from '../aplikacja/zrodlo-posuniec';
import type { Kanal } from '../protokol/kanal';
import { czyGlosDziala } from './dostepnosc-mowy';
import { utworzFavikonPlywajacy } from './favikon-plywajacy';
import { utworzOknoDymkowe } from './okno-dymkowe';
import { utworzStanDymka } from './stan-dymka';


export interface AsystentPlywajacy {
  /** Warstwa do osadzenia; leży poza obszarem roboczym powłoki. */
  element: HTMLElement;
  /** Zdejmuje subskrypcje zdarzeń rdzenia i wspólnego źródła posunięć. */
  rozlacz(): void;
}
/**
 * Składa warstwę asystenta z favikonu, dymka i stanu.
 * @param posuniecia wspólne źródło posunięć, to samo co dla pasa dolnego;
 *   pominięte znaczy, że montaż go nie podał.
 * @param idKlienta identyfikator połączenia, po którym poznaje się urządzenie
 *   obce.
 */
export function zaczepAsystenta(
  kanal: Kanal,
  posuniecia?: ZrodloPosuniec,
  idKlienta?: string,
): AsystentPlywajacy {
  const element = document.createElement('div');
  element.className = 'ap-warstwa';

  const stan = utworzStanDymka(kanal, { posuniecia, idKlienta });
  const dymek = utworzOknoDymkowe(stan, () => przelacz(false));

  const favikon = utworzFavikonPlywajacy({
    glosDziala: czyGlosDziala(),
    naNacisniecie: () => przelacz(!dymek.widoczny()),
  });

  function przelacz(otwarty: boolean): void {
    if (otwarty) dymek.pokaz();
    else dymek.schowaj();
    favikon.ustawOtwarty(otwarty);
    if (otwarty) stan.oznaczPrzeczytane();
  }

  /** Jedyna droga ze stanu na favikon. */
  function odswiezFavikon(): void {
    if (dymek.widoczny()) stan.oznaczPrzeczytane();
    favikon.ustawPrace(stan.praca());
    favikon.ustawNieprzeczytane(stan.nieprzeczytane());
  }

  const przestanObserwowac = stan.obserwuj(odswiezFavikon);
  odswiezFavikon();

  element.append(dymek.element, favikon.element);

  return {
    element,
    rozlacz() {
      przestanObserwowac();
      dymek.rozlacz();
      stan.rozlacz();
    },
  };
}

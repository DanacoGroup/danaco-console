import './asystent-plywajacy.css';

import type { ZrodloPosuniec } from '../aplikacja/zrodlo-posuniec';
import type { Kanal } from '../protokol/kanal';
import { czyGlosDziala } from './dostepnosc-mowy';
import { utworzFavikonPlywajacy } from './favikon-plywajacy';
import { utworzOknoDymkowe } from './okno-dymkowe';
import { utworzStanDymka } from './stan-dymka';

/**
 * Zaczep pływającego Asystenta — jedyne miejsce, w którym favikon, dymek i stan
 * schodzą się w jedną warstwę.
 *
 * Warstwa osadza się w `aplikacja/widok-srodowiska.ts`, bo tylko tam znana jest
 * jednocześnie powłoka i droga do rdzenia. Favikon ląduje w korzeniu powłoki,
 * obok obszaru roboczego: moduł podmienia zawartość obszaru
 * (`aplikacja/przestrzen-modulu.ts`), więc favikon leżący poza nim zostaje
 * widoczny przy każdej zmianie modułu.
 *
 * Dwa inne zaczepy nie nadają się. `aplikacja/scena-sesji.ts` trzyma okna
 * równoległe liczone do sufitu `LICZBA_MAX` (`okna-rownolegle/identyfikatory.ts`);
 * asystent nie jest oknem sceny — profil daje mu `granicaOkien: 1` i postać
 * `dymek-glosowy` — więc zabierałby gniazdo oknu komunikacji i znikał razem ze
 * sceną. `powloka/powloka.ts` buduje także stanowisko podglądu
 * (`powloka/podglad.ts`), które o rdzeniu nie wie; zaczep tam wymusiłby albo
 * wersję niemą, albo przeciek drogi do rdzenia w dół.
 *
 * Dymek jest jeden na widok środowiska. Liczby `granicaOkien` nie przepisano
 * tutaj — `okno-dymkowe.ts` bierze ją z rejestru profilów.
 *
 * Profil daje `pamiecSesyjna: true`, więc zwinięcie dymka nie czyści historii:
 * `schowaj()` chowa element. Rozmowa ginąca po zwinięciu wyglądałaby na awarię.
 *
 * Favikon i stan spina jedna subskrypcja: `stan.obserwuj` przenosi na favikon
 * pracę asystenta i licznik nieprzeczytanych posunięć. Otwarcie dymka zeruje
 * licznik, zamknięcie znów go zbiera — dopóki dymek stoi odsłonięty, posunięcia
 * są czytane na bieżąco.
 */

export interface AsystentPlywajacy {
  /** Warstwa do osadzenia; leży poza obszarem roboczym powłoki. */
  element: HTMLElement;
  /** Zdejmuje subskrypcje zdarzeń rdzenia i wspólnego źródła posunięć. */
  rozlacz(): void;
}

/**
 * @param posuniecia wspólne źródło posunięć — ten sam egzemplarz, którym
 *   karmi się pas dolny. Pominięte znaczy, że miejsce montażu źródła nie
 *   podało; dymek mówi o tym wprost, zamiast milczeć.
 * @param idKlienta identyfikator tego połączenia (`uzgodnienie.klient.id`) —
 *   rozstrzyga „Operator z innego urządzenia" wobec `actorClientId`.
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

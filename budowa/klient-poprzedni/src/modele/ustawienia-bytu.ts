import '../konfiguracja/konfiguracja.css';

import { ConfigAxis, ConfigScope } from '../../../shared/contract';
import { utworzNawigacjeKategorii } from '../konfiguracja/nawigacja-kategorii';
import { utworzPanelKategorii } from '../konfiguracja/panel-kategorii';
import { utworzStanKonfiguracji } from '../konfiguracja/stan-konfiguracji';
import { NAZWY_OSI } from '../konfiguracja/zasiegi';
import type { Kanal } from '../protokol/kanal';
import { otworzOknoPunktowIzolacji } from '../punkty-izolacji/indeks';
import { wskazanieKompletne, type WskazanieOsi } from './wybor-osi';

/**
 * Ustawienia per model i per konto — te same pola co w oknie konfiguracji,
 * oglądane i zapisywane z punktu widzenia jednej osi.
 *
 * Katalog kategorii, katalog definicji, budowa kontrolek, wskaźnik dziedziczenia
 * i zapis `config.set` pochodzą w całości z modułów okna konfiguracji; sekcja
 * modeli dokłada wyłącznie punkt widzenia zawężony do wskazanej osi. Drugi
 * generator pól oznaczałby dwie prawdy o tej samej wartości i dwa miejsca do
 * zmiany po dopisaniu rodzaju wartości do kontraktu.
 *
 * Wskaźnik zasięgu nanosi punkt widzenia na wybór adresu zapisu raz — przy
 * pierwszym odświeżeniu pola. Formularz zbudowany dla jednego modelu, a
 * pozostawiony po przejściu na inny, zapisywałby dalej pod adres poprzedniego,
 * więc po zmianie osi panel dostaje `pokaz`, a nie samo `odswiez`.
 *
 * Pozycja katalogu, która nie dopuszcza wskazanej osi, nie znika z formularza:
 * jej wskaźnik zasięgu poda wtedy poziomy i osie dopuszczone przez katalog,
 * a rozstrzygnięcie o dopuszczalności zapisu zostaje przy rdzeniu.
 */
export interface UstawieniaBytu {
  /** Panel osadzany w sekcji modeli. */
  element: HTMLElement;
  /** Ustawia oś i byt, względem których liczone jest dziedziczenie. */
  ustawWskazanie(wskazanie: WskazanieOsi): void;
  /** Wczytuje katalog i wpisy z rdzenia, po czym przebudowuje formularz. */
  odswiez(): void;
  /** Odłącza subskrypcje kanału. */
  rozlacz(): void;
}

export function utworzUstawieniaBytu(kanal: Kanal): UstawieniaBytu {
  const stan = utworzStanKonfiguracji(kanal);
  const kategorie = utworzNawigacjeKategorii();
  // Przejście do okna punktów izolacji idzie tą samą drogą, co z okna
  // Konfiguracji: tamto okno jest jedno na klienta i pamięta swój stan, więc
  // otwarcie stąd trafia w ten sam egzemplarz. Pominięcie przejścia zostawiłoby
  // w tym oknie pozycję „Izolacja", która nie prowadzi donikąd.
  const panel = utworzPanelKategorii(stan, () => void otworzOknoPunktowIzolacji(kanal));

  const podpis = document.createElement('p');
  podpis.className = 'dn-pole-opis dm-ustawienia__podpis';

  const cialo = document.createElement('div');
  cialo.className = 'dk-okno__cialo dm-ustawienia__cialo';
  cialo.append(kategorie.element, panel.element);

  const element = document.createElement('section');
  element.className = 'dm-ustawienia';
  element.append(podpis, cialo);

  /** Oś czynna; zaczynamy od platformy, bo tak zaczyna okno konfiguracji. */
  let wskazanie: WskazanieOsi = { os: ConfigAxis.Platform, bytOsi: '' };

  /** Przebudowa formularza — jedyna droga po zmianie osi albo katalogu. */
  function przebuduj(): void {
    kategorie.odswiez(stan.kategorie());
    panel.pokaz(kategorie.wybrana());
  }

  function opisz(): void {
    podpis.textContent = wskazanieKompletne(wskazanie)
      ? `Wartości oglądane i zapisywane dla osi: ${NAZWY_OSI[wskazanie.os]}${
          wskazanie.bytOsi === '' ? '' : ` · ${wskazanie.bytOsi}`
        }. Wskaźnik przy polu mówi, z którego poziomu i której osi pochodzi wartość obowiązująca.`
      : `Oś ${NAZWY_OSI[wskazanie.os]} wymaga wskazania bytu. Formularz pokazuje wartości platformy, dopóki byt nie zostanie podany.`;
  }

  kategorie.naWybor((kategoria) => panel.pokaz(kategoria));

  // Zmiana stanu nanosi wartości na pola już zbudowane; składu katalogu nie
  // rusza, więc zapis dokonany gdzie indziej nie przerywa pracy przy polu.
  stan.naZmiane(() => panel.odswiez());

  opisz();

  return {
    element,

    ustawWskazanie(nowe) {
      wskazanie = nowe;
      // Byt niepodany znaczy oś platformy: adresowanie osi bez bytu nie ma
      // w kontrakcie znaczenia.
      const skuteczne = wskazanieKompletne(nowe)
        ? nowe
        : { os: ConfigAxis.Platform, bytOsi: '' };
      stan.ustawPunkt({
        zasieg: ConfigScope.Global,
        bytZasiegu: '',
        os: skuteczne.os,
        bytOsi: skuteczne.bytOsi,
      });
      opisz();
      przebuduj();
    },

    odswiez() {
      void stan.odswiez().then(przebuduj);
    },

    rozlacz: () => stan.rozlacz(),
  };
}

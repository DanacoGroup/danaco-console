import { nazwaNakladu } from './etykiety-sterowania';
import {
  KluczUstawieniaOkna,
  NAKLAD_NAJMNIEJSZY,
  NAKLAD_NAJWIEKSZY,
  nakladZPolozenia,
  polozenieNakladu,
} from './klucze-ustawien';
import { utworzNaglowekSterowania } from './naglowek-sterowania';
import type { MigawkaSterowania, StanSterowania } from './stan-sterowania';
import { zdanieSteru } from './wartosc-obowiazujaca';
import type { ZmianaUstawienia } from './zmiana-ustawienia';

const NAZWA = 'Nakład rozumowania';

/**
 * Sterowanie nakładem rozumowania: suwak od odpowiedzi szybkiej po pełny namysł.
 *
 * Stopień jest wartością wyliczenia katalogu rdzenia (`naklad_rozumowania`), nie
 * nazwą modelu — zmiana kanału nie unieważnia ustawienia, a kanał przekłada
 * stopień na własny parametr. Suwak pokazuje położenie w wykazie stopni, wysyła
 * zaś sam stopień, bo rdzeń nie przyjmuje numeru położenia. Wartość idzie
 * ustawieniem poziomu okna (`config.set`, zasięg `window`), bo treść
 * `window.update` nie ma dla niej pola.
 *
 * Suwak nie ma stanu wyłączonego: przesunięcie jest czynne zawsze, także gdy
 * rdzeń nie odpowiada.
 */
export function utworzSterowanieNakladu(
  stan: StanSterowania,
  ustawienia: ZmianaUstawienia,
): HTMLElement {
  const identyfikator = `dc-ster-naklad-${stan.idOkna()}`;

  const element = document.createElement('div');
  element.className = 'dc-ster-pole';

  const naglowek = utworzNaglowekSterowania(NAZWA, identyfikator);

  const suwak = document.createElement('input');
  suwak.id = identyfikator;
  suwak.type = 'range';
  suwak.className = 'dc-ster-suwak';
  suwak.min = String(NAKLAD_NAJMNIEJSZY);
  suwak.max = String(NAKLAD_NAJWIEKSZY);
  suwak.step = '1';

  const opis = document.createElement('span');
  opis.className = 'dc-ster-suwak__opis';

  const skala = document.createElement('span');
  skala.className = 'dc-ster-suwak__skala';
  skala.textContent = 'szybciej ↔ mądrzej';

  suwak.addEventListener('input', () => {
    opis.textContent = nazwaNakladu(nakladZPolozenia(Number(suwak.value)));
  });
  suwak.addEventListener('change', () => {
    ustawienia.zapisz(
      NAZWA,
      KluczUstawieniaOkna.NakladRozumowania,
      nakladZPolozenia(Number(suwak.value)),
    );
  });

  /**
   * Suwak i zdanie pod nim niosą dwie różne wartości.
   *
   * Położenie suwaka to zapis poziomu okna — wartość ustawiona tutaj i zmieniana
   * przesunięciem. Zdanie pod suwakiem to wartość obowiązująca po rozstrzygnięciu
   * poziomów zasięgu. Gdy wartość przychodzi z poziomu szerszego niż okno, zdanie
   * podaje ten poziom; bez tego suwak stojący na „Bez wskazania" przeczyłby
   * nakładowi narzuconemu z poziomu globalnego.
   */
  function odrysuj(migawka: MigawkaSterowania): void {
    suwak.value = String(polozenieNakladu(migawka.ustawienia.nakladRozumowania));
    opis.textContent = zdanieSteru(
      migawka.obowiazujace.nakladRozumowania,
      migawka.obowiazujace.znane,
      migawka.ustawienia.nakladRozumowania,
      stan.idOkna(),
      nazwaNakladu,
    );
  }

  element.append(naglowek, suwak, skala, opis);

  stan.naZmiane(odrysuj);
  odrysuj(stan.migawka());

  return element;
}

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

/** Sterowanie nakładem rozumowania: suwak od odpowiedzi szybkiej po pełny namysł, zapisywany ustawieniem poziomu okna. */
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

  // Suwak niesie zapis poziomu okna, zdanie pod nim wartość obowiązującą po rozstrzygnięciu.
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

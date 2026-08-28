import type { Component } from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import { poleWyboru, ustawPozycje, type PozycjaWyboru } from '../../modele/kontrolki-formularza';
import type { StanAssistant } from './stan-assistant';

/**
 * Selektor Profilu asystenta — warstwa 2 nagłówka Voice Console. Profil jest
 * komponentem własnym strefy 2 strony głównej, więc jego wykazem jest odczyt
 * `component.list` zawężony rodzajem `assistant`.
 */
export interface SelektorProfilu {
  /** Kontrolka osadzana w rzędzie kontrolek paska promptu. */
  element: HTMLElement;
  /** Odczyt katalogu profili. */
  wczytaj(): Promise<void>;
}

/**
 * Pozycja pusta — profil domyślny rdzenia; zawsze pierwsza w wykazie. Wybranie jej
 * zostawia pole puste, więc zlecenie idzie bez wskazania profilu, a rdzeń dobiera
 * własny domyślny zamiast odmawiać.
 */
const DOMYSLNY: PozycjaWyboru = {
  wartosc: '',
  etykieta: 'Profil asystenta: domyślny rdzenia',
};

export function utworzSelektorProfilu(
  stan: StanAssistant,
  naWybor: (kodProfilu: string) => void,
): SelektorProfilu {
  const kontrolka = poleWyboru(
    {
      etykieta: 'Katalog profili asystenta',
      opis:
        'Wykaz komponentów własnych rodzaju „Profil asystenta" (component.list). ' +
        'Wybór wpisuje kod profilu w pole obok — to ono jedzie do rdzenia polem profileId.',
    },
    [DOMYSLNY],
  );
  kontrolka.kontrolka.addEventListener('change', () => naWybor(kontrolka.kontrolka.value));

  const powod = document.createElement('p');
  powod.className = 'dn-pole-opis ma-brak';
  powod.hidden = true;

  const element = document.createElement('div');
  element.className = 'ma-profile';
  element.append(kontrolka.element, powod);

  function nazwij(zdanie: string): void {
    powod.textContent = zdanie;
    powod.hidden = zdanie === '';
  }

  return {
    element,

    async wczytaj() {
      const wynik = await stan.zaplecze.profile();
      if (!wynik.udany || wynik.wynik === undefined) {
        nazwij(opisOdmowy('Odczyt katalogu profili', wynik.blad?.code, wynik.blad?.message));
        return;
      }
      const profile = wynik.wynik.components;
      ustawPozycje(kontrolka.kontrolka, [DOMYSLNY, ...profile.map(pozycja)]);
      nazwij(
        profile.length === 0
          ? 'Rdzeń nie zna ani jednego komponentu własnego rodzaju „Profil asystenta". ' +
              'Zlecenie pójdzie profilem domyślnym rdzenia albo kodem wpisanym w pole obok.'
          : '',
      );
    },
  };
}

/**
 * Jeden wiersz katalogu; nazwa mówi także to, czego wiersz nie ma. Dopiski
 * o wyłączeniu w strefie 2 i o braku bytu docelowego stoją w nazwie, bo wiersz
 * czyta się w zwiniętej liście, gdzie nie ma miejsca na drugi wiersz opisu.
 */
function pozycja(komponent: Component): PozycjaWyboru {
  const byt = komponent.targetId ?? '';
  const dopiski: string[] = [];
  if (!komponent.enabled) dopiski.push('wyłączony w strefie 2');
  if (byt === '') dopiski.push('bez bytu docelowego — rdzeń nie ma czym wskazać profilu');
  return {
    wartosc: byt,
    etykieta:
      dopiski.length === 0
        ? `Profil asystenta: ${komponent.name}`
        : `Profil asystenta: ${komponent.name} (${dopiski.join('; ')})`,
  };
}

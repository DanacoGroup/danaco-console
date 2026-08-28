import type { Action } from '../../../../shared/contract';
import { utworzDymekObjasnienia } from '../../komponenty/dymek';
import { opisOdmowy } from '../../komponenty/odmowa';
import { przycisk, utworzWierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import type { StanBiblioteki } from './stan-biblioteki';
import type { ZrodloOtoczenia } from './zrodlo-otoczenia';

/**
 * Panel akcji okna operacyjnego bierze pozycje z katalogu rdzenia komendą
 * `action.list` i wykonuje je komendą `window.action`. Czynności, którym
 * kontrakt komendy nie przypisał, stoją obok jako przyciski nazywające brak.
 */
export interface PozycjaBezKomendy {
  etykieta: string;
  powod: string;
}

export interface PanelAkcji {
  element: HTMLElement;
  /** Odczytuje katalog akcji rdzenia; wolno wołać wielokrotnie. */
  wczytaj(): Promise<void>;
}

export function utworzPanelAkcji(
  otoczenie: ZrodloOtoczenia,
  stan: StanBiblioteki,
  bezKomendy: readonly PozycjaBezKomendy[],
): PanelAkcji {
  const odpowiedz = utworzWierszOdpowiedzi();

  const katalog = document.createElement('div');
  katalog.className = 'ml-panel__pozycje';

  const wykazBezKomendy = document.createElement('div');
  wykazBezKomendy.className = 'ml-panel__pozycje ml-panel__pozycje--bez-komendy';
  for (const pozycja of bezKomendy) {
    wykazBezKomendy.append(przyciskBezKomendy(pozycja, odpowiedz.pokaz));
  }

  const tytul = document.createElement('p');
  tytul.className = 'ml-panel__tytul';
  tytul.append(
    'Panel akcji',
    utworzDymekObjasnienia(
      'Pozycje panelu pochodzą z katalogu akcji rdzenia (action.list) i wykonuje ' +
        'je window.action. Czynności wymienione niżej nie mają dziś komendy kontraktu — ' +
        'przycisk mówi to wprost, zamiast milczeć.',
      { powloka: 'ml-dymek', znak: 'ml-dymek__znak' },
    ),
  );

  const element = document.createElement('div');
  element.className = 'ml-panel';
  element.append(tytul, katalog, wykazBezKomendy, odpowiedz.element);

  async function wykonaj(akcja: Action): Promise<void> {
    const idOkna = stan.idOkna();
    if (idOkna === '') {
      odpowiedz.pokaz(
        `Akcja „${akcja.name}" nie ma okna wykonawczego: rdzeń nie oddał ani jednego okna ` +
          'komunikacji tej sesji (window.list).',
        false,
      );
      return;
    }
    odpowiedz.pokaz(`Zgłaszanie akcji „${akcja.name}" do rdzenia…`, true);
    const wynik = await otoczenie.wykonajAkcje(idOkna, akcja.id);
    // Powodzenie znaczy wynik akcji, a nie samo przyjęcie zgłoszenia.
    odpowiedz.pokaz(
      wynik.udany
        ? `Rdzeń oddał wynik akcji „${akcja.name}".`
        : opisOdmowy(`Akcja „${akcja.name}"`, wynik.blad?.code, wynik.blad?.message),
      wynik.udany,
    );
  }

  return {
    element,

    async wczytaj() {
      const wynik = await otoczenie.akcje();
      if (!wynik.udany || wynik.wynik === undefined) {
        katalog.replaceChildren(
          zdanie(opisOdmowy('Katalog akcji', wynik.blad?.code, wynik.blad?.message)),
        );
        return;
      }
      const pozycje = [...wynik.wynik.actions].sort((pierwsza, druga) => pierwsza.order - druga.order);
      if (pozycje.length === 0) {
        katalog.replaceChildren(
          zdanie(
            'Rdzeń nie ma ani jednej akcji w katalogu modułu Library — panel zostaje pusty ' +
              'zamiast pokazywać pozycje zmyślone przez klienta.',
          ),
        );
        return;
      }
      katalog.replaceChildren(
        ...pozycje.map((akcja) => {
          const przyciskAkcji = przycisk(akcja.name, 'dn-btn dn-btn--sm dn-btn--zarys');
          przyciskAkcji.dataset['akcja'] = akcja.id;
          if (akcja.description !== undefined) przyciskAkcji.title = akcja.description;
          przyciskAkcji.addEventListener('click', () => void wykonaj(akcja));
          return przyciskAkcji;
        }),
      );
    },
  };
}

/**
 * Przycisk czynności, której kontrakt nie niesie, pozostaje klikalny i po
 * naciśnięciu nazywa powód braku. Kontrolka wygaszona albo milcząca zostawiłaby
 * Operatora bez odpowiedzi na pytanie, dlaczego czynność się nie odbywa.
 */
function przyciskBezKomendy(
  pozycja: PozycjaBezKomendy,
  pokaz: (tresc: string, powodzenie: boolean) => void,
): HTMLElement {
  const element = przycisk(pozycja.etykieta, 'dn-btn dn-btn--sm dn-btn--duch');
  element.dataset['bezKomendy'] = 'tak';
  element.addEventListener('click', () => {
    pokaz(`„${pozycja.etykieta}" — ${pozycja.powod}`, false);
  });
  return element;
}

/**
 * Zdanie wyjaśniające zajmuje miejsce wykazu pozycji wtedy, gdy katalog rdzenia
 * jest pusty albo jego odczyt się nie powiódł. Puste miejsce bez zdania nie
 * odróżnia katalogu pustego od katalogu nieodczytanego.
 */
function zdanie(tresc: string): HTMLElement {
  const element = document.createElement('p');
  element.className = 'dn-pole-opis';
  element.textContent = tresc;
  return element;
}

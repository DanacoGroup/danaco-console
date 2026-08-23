import type { Action } from '../../../../shared/contract';
import { utworzDymekObjasnienia } from '../../komponenty/dymek';
import { opisOdmowy } from '../../komponenty/odmowa';
import { przycisk, utworzWierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import type { StanBiblioteki } from './stan-biblioteki';
import type { ZrodloOtoczenia } from './zrodlo-otoczenia';

/**
 * Panel akcji okna operacyjnego.
 *
 * Panel nie jest zaszytym wykazem: pozycje przychodzą z katalogu rdzenia
 * (`action.list`, zasięg modułu) i każda niesie własny kod, którym wykonuje ją
 * `window.action`. Zaszycie ich po stronie klienta byłoby drugą kopią katalogu.
 *
 * Czynności, którym kontrakt nie przypisał komendy (eksport, archiwizacja, kosz,
 * wykrywanie duplikatów, udostępnienie odnośnikiem, porównanie), stoją obok jako
 * przyciski klikalne: naciśnięcie nie wysyła nic i nazywa brak — zamiast kontrolki
 * wygaszonej albo milczącej.
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
    // Powodzenie znaczy tu wynik akcji, nie samo przyjęcie zgłoszenia: akcję,
    // której rdzeń nie wykonuje, rdzeń odrzuca wprost (`conflict`, z kodem
    // komendy do wywołania), a panel pokazuje tę odmowę.
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

/** Przycisk czynności, której kontrakt nie niesie: klikalny, odpowiada powodem. */
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

/** Zdanie wyjaśniające w miejscu wykazu pozycji. */
function zdanie(tresc: string): HTMLElement {
  const element = document.createElement('p');
  element.className = 'dn-pole-opis';
  element.textContent = tresc;
  return element;
}

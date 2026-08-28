import { Command, IsolationLayer } from '../../../shared/contract';
import { opisOdmowyBledu } from '../komponenty/odmowa';
import { wiersz, wybor } from '../modele/kontrolki-formularza-braki';
import type { Kanal } from '../protokol/kanal';
import { idSesji, posijKomende } from './komenda';
import { NAZWY_WARSTW, OPISY_WARSTW, type StanWarstwy } from './stan-warstwy';

/** Wybór warstwy izolacji — kontrolka pasa narzędzi okna punktów izolacji, a nie osobna zakładka widoku. */
export interface SterowanieWarstwa {
  /** Element montowany w pasie narzędzi ramy okna. */
  element: HTMLElement;
  /** Ustawia listę na warstwę czynną — po powrocie do okna albo zmianie z zewnątrz. */
  odswiez(): void;
}

const OPCJE_WARSTW: ReadonlyArray<readonly [string, string]> = [
  [IsolationLayer.Default, NAZWY_WARSTW[IsolationLayer.Default]],
  [IsolationLayer.Session, NAZWY_WARSTW[IsolationLayer.Session]],
];

export function utworzSterowanieWarstwa(kanal: Kanal, stan: StanWarstwy): SterowanieWarstwa {
  const lista = wybor('Warstwa izolacji', OPCJE_WARSTW);
  lista.value = stan.warstwa();

  const opis = document.createElement('p');
  opis.className = 'pi-warstwa__opis';

  const meldunek = document.createElement('p');
  meldunek.className = 'pi-warstwa__meldunek';
  meldunek.setAttribute('role', 'status');
  meldunek.hidden = true;

  const pole = wiersz('Warstwa izolacji', lista, {
    klasa: 'dn-pole pi-warstwa__pole',
    objasnienie: 'Dotyczy wszystkich obszarów tego okna — odczytu i zapisu.',
  });

  const element = document.createElement('div');
  element.className = 'pi-warstwa';
  element.append(pole, opis, meldunek);

  function nanies(): void {
    const warstwa = stan.warstwa();
    lista.value = warstwa;
    opis.textContent = OPISY_WARSTW[warstwa];
  }

  function melduj(zdanie: string, udane: boolean): void {
    meldunek.textContent = zdanie;
    meldunek.dataset['udane'] = String(udane);
    meldunek.hidden = zdanie === '';
  }

  async function przelacz(wybrana: IsolationLayer): Promise<void> {
    const poprzednia = stan.warstwa();
    melduj(`Proszę rdzeń o przełączenie warstwy na „${NAZWY_WARSTW[wybrana]}” (isolation.layer.set)…`, true);

    const wynik = await posijKomende(kanal, Command.IsolationLayerSet, {
      layer: wybrana,
      sessionId: idSesji(kanal),
    });

    if (!wynik.udany || wynik.wynik === undefined) {
      nanies();
      melduj(
        `${opisOdmowyBledu(`Przełączenie warstwy na „${NAZWY_WARSTW[wybrana]}”`, wynik.blad)}. ` +
          `Okno pracuje dalej na warstwie „${NAZWY_WARSTW[poprzednia]}”. ` +
          'Wybierz warstwę z listy ponownie — lista pozostaje czynna i nic nie jest zablokowane.',
        false,
      );
      return;
    }

    stan.ustaw(wynik.wynik.layer);
    nanies();
    melduj(`Warstwa czynna: „${NAZWY_WARSTW[wynik.wynik.layer]}”. Obszary czytają teraz z tej warstwy.`, true);
  }

  lista.addEventListener('change', () => {
    const wybrana = lista.value === IsolationLayer.Session ? IsolationLayer.Session : IsolationLayer.Default;
    void przelacz(wybrana);
  });

  nanies();

  return { element, odswiez: nanies };
}

import { ConfigScope, type ConfigEntry } from '../../../shared/contract';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import type { Kanal } from '../protokol/kanal';
import { naZmianeKlucza, odczytajNastawe } from './zrodlo-nastaw';

/** Wiersz wymogu logowania pokazuje wartość bieżącą nastawy i jej nie zmienia, bo okno konfiguracji rysuje dla tego klucza własną kontrolkę na właściwym poziomie zasięgu. */
export interface WierszWymogu {
  /** Element montowany w sekcji. */
  element: HTMLElement;
  /** Odczytuje wartość z rdzenia i przerysowuje wiersz. */
  odswiez(): void;
  /** Zdejmuje subskrypcję kanału. */
  rozlacz(): void;
}

/** Klucz nastawy w katalogu rdzenia — jedyne jego wystąpienie w kliencie, używane do odczytu i nasłuchu zmian. */
export const KLUCZ_WYMOGU_LOGOWANIA = 'gateway.requireLogin';

export function utworzWierszWymoguLogowania(kanal: Kanal): WierszWymogu {
  const element = document.createElement('div');
  element.className = 'du-wiersz du-wiersz--odczyt';
  element.dataset['klucz'] = KLUCZ_WYMOGU_LOGOWANIA;

  const nazwa = document.createElement('span');
  nazwa.className = 'du-wiersz__nazwa';
  nazwa.textContent = 'Wymóg logowania';

  const wartosc = document.createElement('span');
  wartosc.className = 'dn-plakietka du-wiersz__wartosc';
  wartosc.textContent = 'odczyt…';

  const opis = document.createElement('p');
  opis.className = 'du-wiersz__opis';

  const ostrzezenie = document.createElement('p');
  ostrzezenie.className = 'du-wiersz__zdanie';
  ostrzezenie.setAttribute('role', 'alert');
  ostrzezenie.hidden = true;

  element.append(nazwa, wartosc, opis, ostrzezenie);

  // Trzy stany, nie dwa: brak wartości znaczy bez wskazania, rozstrzyga adres nasłuchu.
  function napis(surowa: unknown): string {
    if (surowa === true || surowa === 'true') return 'wymagane';
    if (surowa === false || surowa === 'false') return 'niewymagane';
    return 'bez wskazania — rozstrzyga adres nasłuchu';
  }

  function pokazWartosc(surowa: unknown, opisKatalogu: string): void {
    wartosc.textContent = napis(surowa);
    opis.textContent = opisKatalogu;
  }

  function ostrzez(tresc: string): void {
    ostrzezenie.textContent = tresc;
    ostrzezenie.hidden = tresc === '';
  }

  function odczytaj(): void {
    void odczytajNastawe(kanal, KLUCZ_WYMOGU_LOGOWANIA).then((odczyt) => {
      if (odczyt.definicja === undefined) {
        wartosc.textContent = 'nie do odczytania';
        opis.textContent = odczyt.odmowa;
        return;
      }
      pokazWartosc(
        odczyt.wpis?.value ?? odczyt.definicja.defaultValue,
        `${odczyt.definicja.description ?? ''} Zmiana obowiązuje od następnego ` +
          'połączenia — rdzenia nie trzeba zatrzymywać.',
      );
      ostrzez(odczyt.odmowa);
    });
  }

  const odsubskrybuj: Odsubskrybuj = naZmianeKlucza(
    kanal,
    KLUCZ_WYMOGU_LOGOWANIA,
    (nowa, wpis: ConfigEntry) => {
      if (wpis.scope !== ConfigScope.Application) {
        // Zapis pod adresem, którego katalog nie dopuszcza, a który mimo to wygrywa rozstrzyganie.
        ostrzez(
          `Ktoś zapisał tę nastawę na poziomie „${wpis.scope}", którego katalog ` +
            'nie dopuszcza (dozwolony jest wyłącznie „application"). Rdzeń takiego ' +
            'zapisu dziś nie odrzuca, a poziom węższy przesłania właściwy — ' +
            'wartość obowiązująca może być inna niż pokazana wyżej.',
        );
        return;
      }
      ostrzez('');
      pokazWartosc(nowa, opis.textContent ?? '');
    },
  );

  odczytaj();

  return {
    element,
    odswiez: odczytaj,
    rozlacz: () => odsubskrybuj(),
  };
}

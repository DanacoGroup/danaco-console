import {
  Command,
  NARZEDZIA_MODELU,
  type ToolDeclaration,
} from '../../../../shared/contract';
import { opisOdmowyBledu } from '../../komponenty/odmowa';
import { przycisk } from '../../modele/kontrolki-formularza';
import type { Kanal } from '../../protokol/kanal';
import { czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';
import { utworzRozwiniecie } from './warstwy-translate';

/**
 * Wystawienie operacji modułu na zewnątrz rozdziela dwie rzeczy: operacje modułu są zadeklarowane
 * jako narzędzia modelu, a webhooka wywoływanego z zewnątrz kontrakt nie wystawia.
 */
export interface WystawienieOperacji {
  element: HTMLElement;
}

/** Przedrostek nazw narzędzi obszaru pozwala rozpoznać, które deklaracje kontraktu należą do modułu tłumaczeń. */
const PRZEDROSTEK_NARZEDZIA = 'danaco_translate_';

/** Deklaracje narzędzi obszaru translate stanowią komplet operacji modułu wystawionych modelowi jako narzędzia wywołania. */
export function deklaracjeModulu(): readonly ToolDeclaration[] {
  return NARZEDZIA_MODELU.filter((narzedzie) => narzedzie.name.startsWith(PRZEDROSTEK_NARZEDZIA));
}

export function utworzWystawienieOperacji(kanal: Kanal): WystawienieOperacji {
  const deklaracje = deklaracjeModulu();

  const rozwiniecie = utworzRozwiniecie({
    warstwa: 4,
    nazwa: 'Wystawienie operacji na zewnątrz',
    wyjasnienie:
      `Operacji modułu zadeklarowanych jako narzędzia modelu: ${String(deklaracje.length)}. ` +
      'Webhooka wywoływanego z zewnątrz kontrakt nie wystawia.',
    znacznik: '☰',
  });

  const wykaz = document.createElement('ul');
  wykaz.className = 'mt-wystawienie';
  wykaz.replaceChildren(...deklaracje.map(wierszDeklaracji));

  const zdanie = document.createElement('p');
  zdanie.className = 'mt-wystawienie__zdanie';
  zdanie.textContent =
    'Katalogu pozycji po ukośniku jeszcze nie odpytano — sprawdzenie mówi, czy operacje modułu ' +
    'widać w wykazie, którym Operator dokłada narzędzia do sesji.';

  const katalog = document.createElement('ul');
  katalog.className = 'mt-wystawienie__katalog';

  const sprawdz = przycisk('Sprawdź w katalogu pozycji po ukośniku', 'dn-btn dn-btn--sm dn-btn--zarys');
  sprawdz.addEventListener('click', () => {
    void sprawdzKatalog(kanal, zdanie, katalog);
  });

  const pasek = document.createElement('div');
  pasek.className = 'mt-pasek';
  pasek.append(sprawdz);

  rozwiniecie.tresc.append(wykaz, pasek, zdanie, katalog);
  return { element: rozwiniecie.element };
}

/**
 * `tools.catalog.list` — czy operacje modułu widać w wykazie pozycji sesji.
 *
 * Zawężenie idzie po przedrostku nazwy obszaru, a dopasowanie kontraktu działa
 * także środkiem nazwy, więc pozycja znajdzie się niezależnie od źródła, które
 * ją dostarcza.
 */
async function sprawdzKatalog(
  kanal: Kanal,
  zdanie: HTMLElement,
  katalog: HTMLElement,
): Promise<void> {
  zdanie.textContent = 'Odczyt katalogu pozycji po ukośniku…';
  const wynik = sprawdzKsztalt(
    await wywolaj(kanal, Command.ToolsCatalogList, { query: 'translate' }),
    Command.ToolsCatalogList,
    (tresc) => czyTablica(tresc.entries),
  );

  if (!wynik.udany || wynik.wynik === undefined) {
    katalog.replaceChildren();
    zdanie.textContent = opisOdmowyBledu('Odczyt katalogu pozycji po ukośniku', wynik.blad);
    return;
  }

  const pozycje = wynik.wynik.entries;
  zdanie.textContent =
    `Katalog oddał ${String(pozycje.length)} pozycji przy zawężeniu „translate" ` +
    `(spełniających je łącznie: ${String(wynik.wynik.total)}).`;
  katalog.replaceChildren(
    ...pozycje.map((pozycja) => {
      const element = document.createElement('li');
      element.className = 'mt-wystawienie__pozycja';
      element.textContent =
        `${pozycja.name} — ${pozycja.group} · ` +
        (pozycja.attachable ? 'da się dołożyć do sesji' : 'bez dokładania do sesji');
      return element;
    }),
  );
}

/** Jedna deklaracja narzędzia niesie nazwę, komendę oraz pola żądania, którymi model wywołuje operację modułu. */
function wierszDeklaracji(narzedzie: ToolDeclaration): HTMLElement {
  const nazwa = document.createElement('span');
  nazwa.className = 'mt-wystawienie__nazwa';
  nazwa.textContent = narzedzie.name;

  const komenda = document.createElement('span');
  komenda.className = 'mt-wystawienie__komenda';
  komenda.textContent = narzedzie.command;

  const pola = document.createElement('p');
  pola.className = 'mt-wystawienie__pola';
  pola.textContent =
    narzedzie.parameters.length === 0
      ? 'Żądanie bez pól.'
      : `Pola żądania: ${narzedzie.parameters
          .map((pole) => `${pole.name}${pole.required ? '' : ' (nieobowiązkowe)'}`)
          .join(' · ')}`;

  const element = document.createElement('li');
  element.className = 'mt-wystawienie__wiersz';
  element.append(nazwa, komenda, pola);
  return element;
}

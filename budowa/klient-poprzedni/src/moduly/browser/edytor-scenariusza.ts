/**
 * Widok tekstowy scenariusza — warstwa czwarta Automation Studio, dla treści,
 * których nie da się wyklikać kartami kroków.
 */

import { AutomationStepKind, type AutomationStep } from '../../../../shared/contract';
import { utworzDymekObjasnienia } from '../../komponenty/dymek';
import { poleWielowierszowe } from '../../modele/kontrolki-formularza';
import { KLASY_DYMKA, OBJASNIENIA } from './etykiety-browser';
import { KLASA_PRZYCISKU, przyciskCzynnosci } from './przyciski-browser';

/** Kroki odczytane z widoku edytora wraz z powodem odrzucenia; pusty powód znaczy „przyjęte bez zastrzeżeń”. */
export interface OdczytKrokow {
  kroki: AutomationStep[];
  blad: string;
}

export interface EdytorScenariusza {
  element: HTMLElement;
  /** Wypełnia oba widoki krokami — z rdzenia albo z kart kroków okna. */
  wczytaj(kroki: readonly AutomationStep[]): void;
  /** Kroki wpisane w widoku JSON wraz z powodem odrzucenia. */
  odczytaj(): OdczytKrokow;
}

/** Rodzaje kroku scenariusza, które niesie kontrakt — sprawdzian treści przyjmuje wyłącznie te wartości. */
const RODZAJE_KROKU: readonly string[] = Object.values(AutomationStepKind);

export function utworzEdytorScenariusza(): EdytorScenariusza {
  const zapis = poleWielowierszowe(
    {
      etykieta: 'Kroki scenariusza (JSON)',
      opis: 'Tablica kroków. Pola: id, name, kind, command, params, dependsOn, condition, order.',
    },
    10,
  );

  const podglad = document.createElement('pre');
  podglad.className = 'mb-scenariusz__yaml';
  podglad.hidden = true;

  const notacja = przyciskCzynnosci('Pokaż w notacji YAML', KLASA_PRZYCISKU.zarys, () => {
    const yaml = !podglad.hidden;
    podglad.hidden = yaml;
    zapis.element.hidden = !yaml;
    notacja.textContent = yaml ? 'Pokaż w notacji YAML' : 'Wróć do zapisu w JSON';
  });

  const uwaga = document.createElement('p');
  uwaga.className = 'dn-pole-opis mb-uwaga';
  uwaga.textContent =
    'Widok YAML jest odczytem kroków; zapisywana jest treść pola JSON. ' +
    'Klient nie niesie czytnika YAML.';

  const pasek = document.createElement('div');
  pasek.className = 'mb-panel__pasek';
  pasek.append(notacja, utworzDymekObjasnienia(OBJASNIENIA.edytorScenariusza, KLASY_DYMKA));

  const element = document.createElement('section');
  element.className = 'mb-panel mb-scenariusz';
  element.setAttribute('aria-label', 'Widok tekstowy scenariusza — warstwa czwarta');
  element.append(zapis.element, podglad, pasek, uwaga);

  return {
    element,

    wczytaj(kroki) {
      zapis.kontrolka.value = JSON.stringify(kroki, null, 2);
      podglad.textContent = naYaml(kroki);
    },

    odczytaj() {
      const tresc = zapis.kontrolka.value.trim();
      if (tresc === '') return { kroki: [], blad: '' };
      const wynik = rozbierz(tresc);
      // Widok YAML nadąża za treścią przyjętą, nie za każdą wpisaną literą.
      if (wynik.blad === '') podglad.textContent = naYaml(wynik.kroki);
      return wynik;
    },
  };
}

/** Treść pola JSON przełożona na kroki scenariusza wraz ze sprawdzianem poprawności każdej jego pozycji. */
function rozbierz(tresc: string): OdczytKrokow {
  let odczytane: unknown;
  try {
    odczytane = JSON.parse(tresc);
  } catch (blad) {
    const powod = blad instanceof Error ? blad.message : 'treść nie jest poprawnym zapisem JSON';
    return { kroki: [], blad: `Kroków nie odczytano: ${powod}.` };
  }
  if (!Array.isArray(odczytane)) {
    return { kroki: [], blad: 'Kroki mają być tablicą — w polu stoi pojedyncza wartość.' };
  }

  const kroki: AutomationStep[] = [];
  for (const [pozycja, wpis] of odczytane.entries()) {
    const numer = pozycja + 1;
    if (typeof wpis !== 'object' || wpis === null || Array.isArray(wpis)) {
      return { kroki: [], blad: `Krok ${numer} nie jest obiektem.` };
    }
    const pola = wpis as Record<string, unknown>;
    const identyfikator = typeof pola['id'] === 'string' ? pola['id'].trim() : '';
    if (identyfikator === '') {
      return { kroki: [], blad: `Krok ${numer} nie ma pola id — rdzeń wymaga go w każdym kroku.` };
    }
    const rodzaj = typeof pola['kind'] === 'string' ? pola['kind'] : '';
    if (!RODZAJE_KROKU.includes(rodzaj)) {
      return {
        kroki: [],
        blad: `Krok ${numer} ma rodzaj „${rodzaj}"; kontrakt zna: ${RODZAJE_KROKU.join(', ')}.`,
      };
    }
    kroki.push(wpis as AutomationStep);
  }
  return { kroki, blad: '' };
}

/**
 * Kroki w notacji YAML — blok sekwencji z odwzorowaniami, gdzie treść żądania
 * kroku zostaje zapisem JSON w jednym wierszu notacji przepływowej YAML.
 */
function naYaml(kroki: readonly AutomationStep[]): string {
  if (kroki.length === 0) return '# scenariusz nie ma ani jednego kroku';
  return kroki.map(krokNaYaml).join('\n');
}

function krokNaYaml(krok: AutomationStep): string {
  const wiersze: string[] = [`- id: ${JSON.stringify(krok.id)}`];
  if (krok.name !== undefined) wiersze.push(`  name: ${JSON.stringify(krok.name)}`);
  wiersze.push(`  kind: ${JSON.stringify(krok.kind)}`);
  if (krok.command !== undefined) wiersze.push(`  command: ${JSON.stringify(krok.command)}`);
  if (krok.params !== undefined) wiersze.push(`  params: ${JSON.stringify(krok.params)}`);
  if (krok.dependsOn !== undefined) wiersze.push(`  dependsOn: ${JSON.stringify(krok.dependsOn)}`);
  if (krok.condition !== undefined) wiersze.push(`  condition: ${JSON.stringify(krok.condition)}`);
  if (krok.order !== undefined) wiersze.push(`  order: ${String(krok.order)}`);
  return wiersze.join('\n');
}

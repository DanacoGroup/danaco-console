import { AutomationStepKind, Command, type AutomationStep } from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import {
  pobierzPlik,
  pole,
  poleTresci,
  przyciskAkcji,
  utworzWierszOdpowiedzi,
  wiersz,
  type WierszOdpowiedzi,
} from '../../modele/kontrolki-formularza';
import { BRAKI, zglosBrak } from './braki-kontraktu';
import { ODCZYTY } from './etykiety-assistant';
import type { ZrodloNarzedzi } from './zrodlo-narzedzi';

/**
 * Edytor makra w Command & Tools Hub definiuje kroki i przekazuje je do modułu Automations
 * komendą zapisu przebiegu, bo makro nie ma własnego magazynu w kontrakcie.
 */
export interface EdytorMakra {
  element: HTMLElement;
}

/**
 * Kroki podpowiedziane w pustym edytorze pokazują kształt, nie treść do wysłania, z nazwą
 * komendy pochodzącą wprost z kontraktu.
 */
const WZOR_KROKOW = `[
  {
    "id": "krok-1",
    "name": "przykład: odczyt dziennika asystenta",
    "kind": "${AutomationStepKind.Command}",
    "command": "${Command.AssistantActivityList}",
    "order": 1
  }
]`;

export function utworzEdytorMakra(zrodlo: ZrodloNarzedzi): EdytorMakra {
  const odpowiedz: WierszOdpowiedzi = utworzWierszOdpowiedzi();

  const nazwa = pole('Nazwa makra', 'nazwa, pod którą makro stanie w module Automations');
  const opis = pole('Opis makra', 'po co makro jest i kiedy je uruchamiać');
  const kroki = poleTresci('Kroki makra w JSON', 10, WZOR_KROKOW, 'ma-kod');

  const wyslij = przyciskAkcji('→ Wyślij do Automations', 'dn-btn dn-btn--sm dn-btn--atrament');
  wyslij.addEventListener('click', () => void wyslijMakro());

  const sprawdz = przyciskAkcji('Sprawdź składnię kroków', 'dn-btn dn-btn--sm dn-btn--zarys');
  sprawdz.addEventListener('click', () => {
    const rozbior = rozbierzKroki(kroki.value);
    odpowiedz.pokaz(
      rozbior.powod !== ''
        ? rozbior.powod
        : `Składnia poprawna: ${String(rozbior.kroki.length)} kroków gotowych do zapisu.`,
      rozbior.powod === '',
    );
  });

  const pobierz = przyciskAkcji('Pobierz definicję makra', 'dn-btn dn-btn--sm dn-btn--zarys');
  pobierz.addEventListener('click', () => {
    // Pobranie nie potrzebuje komendy: plik oddaje dokładnie to, co Operator widzi w polu.
    pobierzPlik(
      `makro-asystenta-${nazwaPliku(nazwa.value)}.json`,
      JSON.stringify({ name: nazwa.value, description: opis.value, steps: kroki.value }, null, 2),
      'application/json',
    );
    odpowiedz.pokaz('Definicja makra pobrana jako plik JSON.', true);
  });

  const zapisWlasny = przyciskAkcji('Zapisz makro w profilu', 'dn-btn dn-btn--sm dn-btn--duch');
  zapisWlasny.dataset['brak'] = 'makro-wlasne';
  zapisWlasny.addEventListener('click', () => zglosBrak('Zapis makra w profilu', BRAKI.makroWlasne));

  const skladnia = przyciskAkcji('Definicja w YAML', 'dn-btn dn-btn--sm dn-btn--duch');
  skladnia.dataset['brak'] = 'makro-yaml';
  skladnia.addEventListener('click', () => zglosBrak('Definicja makra w YAML', BRAKI.makroYaml));

  const przyciski = document.createElement('div');
  przyciski.className = 'ma-formularz__przyciski';
  przyciski.append(wyslij, sprawdz, pobierz, zapisWlasny, skladnia);

  const element = document.createElement('div');
  element.className = 'ma-obszar';
  element.dataset['obszar'] = 'makra';
  element.append(
    wiersz('Nazwa', nazwa, {
      klasa: 'ma-wiersz',
      objasnienie: 'Kontrakt wymaga nazwy automatyki; bez niej rdzeń zapisu odmówi.',
    }),
    wiersz('Opis', opis, { klasa: 'ma-wiersz', objasnienie: 'Pole opcjonalne definicji.' }),
    wiersz('Kroki', kroki, {
      klasa: 'ma-wiersz',
      objasnienie:
        'Tablica kroków w kształcie AutomationStep: wymagane są identyfikator (id) ' +
        'i rodzaj (kind). Krok rodzaju „command" wywołuje komendę kontraktu polem ' +
        'command, a jej ładunek idzie polem params.',
    }),
    przyciski,
    odpowiedz.element,
  );

  async function wyslijMakro(): Promise<void> {
    if (nazwa.value.trim() === '') {
      odpowiedz.pokaz('Nadaj makru nazwę — rdzeń wymaga jej w polu name.', false);
      return;
    }
    const rozbior = rozbierzKroki(kroki.value);
    if (rozbior.powod !== '') {
      odpowiedz.pokaz(rozbior.powod, false);
      return;
    }
    odpowiedz.pokaz(ODCZYTY.wysylkaMakra, true);
    const wynik = await zrodlo.zapiszAutomatyke({
      name: nazwa.value.trim(),
      steps: rozbior.kroki,
      enabled: true,
      ...(opis.value.trim() === '' ? {} : { description: opis.value.trim() }),
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(
        opisOdmowy('Wysłanie makra do Automations', wynik.blad?.code, wynik.blad?.message),
        false,
      );
      return;
    }
    // Potwierdzenie mówi to, co zapisał rdzeń — liczba kroków bywa inna niż w zamówieniu.
    const zapisana = wynik.wynik.workflow;
    odpowiedz.pokaz(
      `Rdzeń zapisał automatykę ${zapisana.id} („${zapisana.name}") z ` +
        `${String(zapisana.steps?.length ?? 0)} krokami. Harmonogram nadasz jej w zakładce ` +
        'rutyn.',
      true,
    );
  }

  return { element };
}

/** Wynik rozbioru pola kroków niesie kroki albo powód, dla którego ich nie ma, gdy zapis Operatora nie jest poprawnym JSON. */
interface RozbiorKrokow {
  kroki: AutomationStep[];
  /** Pusty napis znaczy „rozbiór się udał". */
  powod: string;
}

/**
 * Rozbiór pola kroków wraz ze sprawdzeniem kształtu wymaganego przez kontrakt: identyfikator
 * i rodzaj każdego kroku.
 */
function rozbierzKroki(zapis: string): RozbiorKrokow {
  const tresc = zapis.trim();
  if (tresc === '') {
    return { kroki: [], powod: 'Wpisz kroki makra — pusta definicja nie ma czego wykonać.' };
  }
  let rozebrane: unknown;
  try {
    rozebrane = JSON.parse(tresc);
  } catch (blad) {
    // Powód niesie komunikat parsera, bo mówi, w którym miejscu składnia siada.
    const zdanie = blad instanceof Error ? blad.message : String(blad);
    return { kroki: [], powod: `Kroki nie są poprawnym JSON-em: ${zdanie}` };
  }
  if (!Array.isArray(rozebrane)) {
    return {
      kroki: [],
      powod: 'Kroki muszą być tablicą — kontrakt przyjmuje w polu steps tablicę AutomationStep.',
    };
  }
  const rodzaje: readonly string[] = Object.values(AutomationStepKind);
  const kroki: AutomationStep[] = [];
  for (const [numer, pozycja] of rozebrane.entries()) {
    if (typeof pozycja !== 'object' || pozycja === null) {
      return { kroki: [], powod: `Krok ${String(numer + 1)} nie jest obiektem.` };
    }
    const krok = pozycja as Partial<AutomationStep>;
    if (typeof krok.id !== 'string' || krok.id === '') {
      return { kroki: [], powod: `Krok ${String(numer + 1)} nie ma identyfikatora (pole id).` };
    }
    if (typeof krok.kind !== 'string' || !rodzaje.includes(krok.kind)) {
      return {
        kroki: [],
        powod:
          `Krok ${String(numer + 1)} ma rodzaj spoza kontraktu. Dopuszczalne rodzaje: ` +
          `${rodzaje.join(', ')}.`,
      };
    }
    kroki.push(krok as AutomationStep);
  }
  return { kroki, powod: '' };
}

/** Nazwa makra sprowadzona do bezpiecznego członu nazwy pliku, bez znaków, których system plików nie akceptuje. */
function nazwaPliku(nazwa: string): string {
  const oczyszczona = nazwa
    .trim()
    .toLocaleLowerCase('pl')
    .replace(/[^\p{Letter}\p{Number}]+/gu, '-')
    .replace(/^-+|-+$/g, '');
  return oczyszczona === '' ? 'bez-nazwy' : oczyszczona;
}

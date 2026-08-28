import type { GlossaryTerm, TranslateGlossarySetRequest } from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import {
  poleLogiczne,
  poleTekstowe,
  przycisk,
  type PoleFormularza,
  type WierszOdpowiedzi,
} from '../../modele/kontrolki-formularza';
import { OBJASNIENIA, PODPOWIEDZ_JEZYKOW } from './etykiety-translate';
import { dopnijDymek, podepnijPodpowiedz } from './kontrolki-translate';
import type { StanOkna } from './stan-okna-translate';
import { rozbieznoscOdpowiedzi, type PoleOdpowiedzi } from './zgodnosc-odpowiedzi';
import type { ZrodloGlosariusza } from './zrodlo-glosariusza';

/**
 * Formularz definicji terminu obsługuje w module Translate dwie czynności Glossary Managera jednym kodem.
 */
export interface FormularzTerminu {
  element: HTMLElement;
  /** Wczytuje termin z wykazu do edycji. */
  wczytaj(termin: GlossaryTerm): void;
}

interface PolaTerminu {
  zrodlowy: PoleFormularza<HTMLInputElement>;
  jezyk: PoleFormularza<HTMLInputElement>;
  podpowiedzi: HTMLDataListElement;
  odpowiednik: PoleFormularza<HTMLInputElement>;
  uwaga: PoleFormularza<HTMLInputElement>;
  nieTlumacz: PoleFormularza<HTMLInputElement>;
}

/**
 * Termin będący przedmiotem edycji w tym formularzu glosariusza; pusty identyfikator znaczy termin nowy.
 */
interface Edycja {
  id: string;
}

export function utworzFormularzTerminu(
  zrodlo: ZrodloGlosariusza,
  okno: StanOkna,
  odpowiedz: WierszOdpowiedzi,
  naZapis: (termin: GlossaryTerm) => void,
): FormularzTerminu {
  const pola = zbudujPolaTerminu();
  const edycja: Edycja = { id: '' };

  const nowy = przycisk('+ Nowy termin', 'dn-btn dn-btn--sm dn-btn--zarys');
  nowy.addEventListener('click', () => {
    wyczysc(pola, edycja);
    odpowiedz.pokaz('Formularz pusty — wpisz nowy termin.', true);
  });

  const zapisz = przycisk('Zapisz termin', 'dn-btn dn-btn--sm dn-btn--atrament');
  zapisz.addEventListener('click', () => {
    void zapiszTermin(zrodlo, pola, edycja, { okno, odpowiedz, naZapis });
  });

  const pasek = document.createElement('div');
  pasek.className = 'mt-pasek';
  pasek.append(nowy, zapisz);

  const element = document.createElement('div');
  element.className = 'mt-formularz-terminu';
  element.append(
    pola.zrodlowy.element,
    pola.jezyk.element,
    pola.podpowiedzi,
    pola.odpowiednik.element,
    pola.uwaga.element,
    pola.nieTlumacz.element,
    pasek,
  );

  return {
    element,
    wczytaj(termin) {
      edycja.id = termin.id;
      pola.zrodlowy.kontrolka.value = termin.source;
      pola.jezyk.kontrolka.value = termin.language;
      pola.odpowiednik.kontrolka.value = termin.target ?? '';
      pola.uwaga.kontrolka.value = termin.note ?? '';
      pola.nieTlumacz.kontrolka.checked = termin.doNotTranslate === true;
      odpowiedz.pokaz(`Termin „${termin.source}" wczytany do edycji.`, true);
    },
  };
}

/**
 * Odbiorcy wyniku zapisu terminu glosariusza: okno, wiersz odpowiedzi formularza i cały wykaz terminów.
 */
interface OdbiorcyZapisu {
  okno: StanOkna;
  odpowiedz: WierszOdpowiedzi;
  naZapis: (termin: GlossaryTerm) => void;
}

async function zapiszTermin(
  zrodlo: ZrodloGlosariusza,
  pola: PolaTerminu,
  edycja: Edycja,
  odbiorcy: OdbiorcyZapisu,
): Promise<void> {
  const termin = pola.zrodlowy.kontrolka.value.trim();
  const kod = pola.jezyk.kontrolka.value.trim();
  if (termin === '' || kod === '') {
    odbiorcy.odpowiedz.pokaz(
      'Termin źródłowy i język odpowiednika są wymagane — bez nich rdzeń odmówi zapisu.',
      false,
    );
    return;
  }
  const zadanie: TranslateGlossarySetRequest = { source: termin, language: kod };
  if (edycja.id !== '') zadanie.termId = edycja.id;
  const cel = pola.odpowiednik.kontrolka.value.trim();
  if (cel !== '') zadanie.target = cel;
  const uwaga = pola.uwaga.kontrolka.value.trim();
  if (uwaga !== '') zadanie.note = uwaga;
  if (pola.nieTlumacz.kontrolka.checked) zadanie.doNotTranslate = true;

  // Stan ładowania obejmuje całe okno, bo zapis terminu jest jedynym wywołaniem rdzenia z formularza.
  odbiorcy.okno.ladowanie(
    edycja.id === ''
      ? `Rdzeń zapisuje nowy termin „${termin}" w glosariuszu.`
      : `Rdzeń zapisuje zmianę terminu „${termin}" w glosariuszu.`,
  );
  odbiorcy.odpowiedz.pokaz('Zapisywanie terminu…', true);
  const wynik = await zrodlo.zapisz(zadanie);
  if (!wynik.udany || wynik.wynik === undefined) {
    const zdanie = opisOdmowy('Zapis terminu glosariusza', wynik.blad?.code, wynik.blad?.message);
    odbiorcy.odpowiedz.pokaz(zdanie, false);
    odbiorcy.okno.blad(zdanie);
    return;
  }
  // Sprawdzane są cztery pola zamówienia przeciw terminowi, który wrócił z rdzenia po zapisie.
  const zapisany = wynik.wynik.term;
  const sprawdzane: PoleOdpowiedzi[] = [
    { nazwa: 'termin źródłowy', zamowione: termin, oddane: zapisany.source },
    { nazwa: 'język odpowiednika', zamowione: kod, oddane: zapisany.language },
  ];
  if (cel !== '') {
    sprawdzane.push({ nazwa: 'odpowiednik docelowy', zamowione: cel, oddane: zapisany.target });
  }
  if (edycja.id !== '') {
    sprawdzane.push({ nazwa: 'identyfikator terminu', zamowione: edycja.id, oddane: zapisany.id });
  }
  const rozbiezne = rozbieznoscOdpowiedzi('Zapis terminu glosariusza', sprawdzane);

  // Wykaz okna napełnia się terminem z rdzenia także przy rozbieżności odpowiedzi z zamówieniem.
  edycja.id = '';
  odbiorcy.naZapis(zapisany);
  if (rozbiezne !== null) {
    odbiorcy.odpowiedz.pokaz(rozbiezne, false);
    odbiorcy.okno.blad(rozbiezne);
    return;
  }
  // Fazę ładowania i fazę błędu zdejmuje wyłącznie udany zapis terminu, nie samo przerysowanie okna.
  odbiorcy.okno.gotowe();
  odbiorcy.odpowiedz.pokaz(`Termin „${zapisany.source}" zapisany w rdzeniu.`, true);
}

function wyczysc(pola: PolaTerminu, edycja: Edycja): void {
  edycja.id = '';
  pola.zrodlowy.kontrolka.value = '';
  pola.odpowiednik.kontrolka.value = '';
  pola.uwaga.kontrolka.value = '';
  pola.nieTlumacz.kontrolka.checked = false;
}

function zbudujPolaTerminu(): PolaTerminu {
  const zrodlowy = poleTekstowe({ etykieta: 'Termin źródłowy', podpowiedz: 'np. przetarg' });
  const jezyk = poleTekstowe({ etykieta: 'Język odpowiednika', podpowiedz: 'np. en' });
  const podpowiedzi = podepnijPodpowiedz(jezyk.kontrolka, 'mt-jezyki-terminu', PODPOWIEDZ_JEZYKOW);
  const odpowiednik = poleTekstowe({ etykieta: 'Odpowiednik docelowy', podpowiedz: 'np. tender' });
  const uwaga = poleTekstowe({ etykieta: 'Uwaga do terminu', podpowiedz: 'kontekst użycia' });
  const nieTlumacz = poleLogiczne({ etykieta: 'Nie tłumacz tego terminu' });
  dopnijDymek(nieTlumacz.element, OBJASNIENIA.nieTlumacz);
  return { zrodlowy, jezyk, podpowiedzi, odpowiednik, uwaga, nieTlumacz };
}

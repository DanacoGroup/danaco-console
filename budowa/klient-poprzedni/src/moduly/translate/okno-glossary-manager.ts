import type { GlossaryTerm } from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import {
  poleTekstowe,
  przyciskBezKomendy,
  utworzWierszOdpowiedzi,
  type WierszOdpowiedzi,
} from '../../modele/kontrolki-formularza';
import { PUSTE } from './etykiety-translate';
import { utworzFormularzTerminu } from './formularz-terminu';
import { naglowekOkna } from './kontrolki-translate';
import { utworzStanOkna, type StanOkna } from './stan-okna-translate';
import type { StanTranslate } from './stan-translate';
import { utworzRozwiniecie } from './warstwy-translate';
import { utworzWierszTerminu } from './wiersz-terminu';
import { utworzWymianeGlosariusza } from './wymiana-glosariusza';
import { rozbieznoscOdpowiedzi } from './zgodnosc-odpowiedzi';
import type { ZrodloGlosariusza } from './zrodlo-glosariusza';

/**
 * Glossary Manager — okno **zarządca** modułu Translate
 * (`translate.glossary-manager`).
 *
 * Kontrakt nie ma komendy odczytu glosariusza — w obszarze `translate` nie
 * występuje ani `translate.glossary.list`, ani `.get`. Okno pokazuje więc
 * wyłącznie terminy zapisane w tej sesji i mówi o tym wprost w stanie pustym;
 * wykaz zmyślony albo pusta lista podana jako „glosariusz jest pusty" byłyby
 * atrapą.
 */
export interface OknoGlossaryManager {
  element: HTMLElement;
  odswiez(): void;
}

const BRAKI = {
  dziedzina:
    'Kategorii dziedziny termin w kontrakcie nie ma: niesie źródło, język, odpowiednik, znacznik ' +
    '„nie tłumacz" i uwagę. Bazy nie da się zawęzić do dziedziny, bo nie ma czym jej oznaczyć.',
  ekstrakcja:
    'Kontrakt nie ma komendy wyprowadzającej kandydatów na terminy z tekstu źródłowego.',
  listy:
    'Termin ma w kontrakcie jeden znacznik stanu — „nie tłumacz". Statusu zatwierdzony, kandydat ' +
    'ani zabroniony nie ma czym zapisać, więc list preferowanych i zabronionych nie ma z czego złożyć.',
  historia:
    'Historii zmian wpisu kontrakt nie prowadzi: odpowiedź zapisu niesie termin po zmianie, ' +
    'a wersji poprzedniej nie oddaje żadna komenda obszaru.',
  wymuszanie:
    'Sposobu wstrzyknięcia odpowiednika do zapytania silnika kontrakt nie wystawia. Ujednolicenie ' +
    'terminologii wykonuje rdzeń po swojemu, a reguł tego zabiegu żądanie nie niesie.',
} as const;

export function utworzOknoGlossaryManager(stan: StanTranslate): OknoGlossaryManager {
  const okno: StanOkna = utworzStanOkna(PUSTE.glosariusz);
  const zapisane = new Map<string, GlossaryTerm>();
  const odpowiedz = utworzWierszOdpowiedzi();

  const wykaz = document.createElement('ul');
  wykaz.className = 'mt-terminy';

  const wystapienia = document.createElement('ul');
  wystapienia.className = 'mt-wystapienia';

  const formularz = utworzFormularzTerminu(stan.glosariusz, okno, odpowiedz, (termin) => {
    zapisane.set(termin.id, termin);
    odswiez();
  });

  const wymiana = utworzWymianeGlosariusza(stan.glosariusz, odpowiedz);

  // Wyszukiwanie działa na tym, co okno ma: na terminach zapisanych w tej sesji.
  // Kontrakt nie ma komendy odczytu glosariusza, więc pole nie pyta rdzenia
  // i nie udaje, że przeszukuje bazę — przesiewa wykaz widoczny obok.
  const szukaj = poleTekstowe({
    etykieta: 'Szukaj w terminach zapisanych w tej sesji',
    podpowiedz: 'termin źródłowy, odpowiednik albo język',
  });
  szukaj.kontrolka.addEventListener('input', () => odswiez());

  okno.tresc.append(
    formularz.element,
    szukaj.element,
    odpowiedz.element,
    wykaz,
    wystapienia,
    wymiana.element,
    filtrDziedziny(),
    menuBazy(),
    regulyWymuszania(),
  );

  const element = document.createElement('section');
  element.className = 'mt-okno mt-okno--zarzadca';
  element.dataset['okno'] = 'glossary-manager';
  element.append(naglowekOkna('Glossary Manager', 'zarządca'), okno.element);

  /**
   * Przerysowanie odbudowuje wykaz terminów, ale nie rusza fazy trwającej ani
   * fazy błędu — te zdejmuje czynność, która je postawiła (zapis terminu,
   * odczyt wystąpień). Inaczej zapis, który napełnia wykaz w środku własnego
   * wywołania, kasowałby sobie zapowiedź tego wywołania.
   */
  function odswiez(): void {
    const szukane = szukaj.kontrolka.value.trim().toLowerCase();
    const terminy = [...zapisane.values()].filter((termin) => pasujeTermin(termin, szukane));
    wykaz.replaceChildren(
      ...terminy.map((termin) =>
        utworzWierszTerminu(termin, {
          naEdycje: formularz.wczytaj,
          naWystapienia: (wskazany) => {
            void pokazWystapienia(stan.glosariusz, wskazany, wystapienia, okno, odpowiedz);
          },
        }),
      ),
    );
    if (okno.faza() === 'ladowanie' || okno.faza() === 'blad') return;
    // Dwie pustki są tu różne i nie wolno ich zlać: baza bez ani jednego wpisu
    // to stan pierwszego użycia, a zawężenie bez trafienia to wynik szukania.
    if (zapisane.size === 0) {
      okno.puste(PUSTE.glosariusz);
      return;
    }
    if (terminy.length === 0) {
      okno.puste({
        tytul: 'Zawężenie bez trafienia',
        opis:
          `Żaden z ${String(zapisane.size)} terminów zapisanych w tej sesji nie pasuje do ` +
          `„${szukane}". Wyczyść pole szukania, aby zobaczyć wykaz w całości.`,
      });
      return;
    }
    okno.gotowe();
  }

  odswiez();
  return { element, odswiez };
}

/** Dopasowanie terminu do zawężenia — środkiem źródła, odpowiednika i języka. */
function pasujeTermin(termin: GlossaryTerm, szukane: string): boolean {
  if (szukane === '') return true;
  const stog = `${termin.source} ${termin.target ?? ''} ${termin.language} ${termin.note ?? ''}`;
  return stog.toLowerCase().includes(szukane);
}

/** Warstwa druga: filtr dziedziny — kategoria, której termin w kontrakcie nie ma. */
function filtrDziedziny(): HTMLElement {
  const rozwiniecie = utworzRozwiniecie({
    warstwa: 2,
    nazwa: 'Filtr dziedziny',
    wyjasnienie: BRAKI.dziedzina,
    znacznik: '▼',
  });
  rozwiniecie.tresc.append(przyciskBezKomendy('Zawężenie do dziedziny', BRAKI.dziedzina));
  return rozwiniecie.element;
}

/** Warstwa trzecia: menu bazy — czynności zbiorcze poza wymianą plikową. */
function menuBazy(): HTMLElement {
  const rozwiniecie = utworzRozwiniecie({
    warstwa: 3,
    nazwa: 'Menu bazy terminologicznej',
    wyjasnienie:
      'Ekstrakcja kandydatów na terminy, listy preferowane i zabronione oraz historia wpisu. ' +
      'Import i eksport TBX/CSV stoją wyżej, bo mają komendy kontraktu.',
    znacznik: '☰',
  });
  rozwiniecie.tresc.append(
    przyciskBezKomendy('Ekstrakcja kandydatów na terminy', BRAKI.ekstrakcja),
    przyciskBezKomendy('Listy preferowane i zabronione', BRAKI.listy),
    przyciskBezKomendy('Historia zmian wpisu', BRAKI.historia),
  );
  return rozwiniecie.element;
}

/** Warstwa czwarta: reguły wymuszania terminologii w silnikach przekładu. */
function regulyWymuszania(): HTMLElement {
  const rozwiniecie = utworzRozwiniecie({
    warstwa: 4,
    nazwa: 'Reguły wymuszania terminologii',
    wyjasnienie: BRAKI.wymuszanie,
    znacznik: '☰',
  });
  rozwiniecie.tresc.append(
    przyciskBezKomendy('Sposób wstrzyknięcia odpowiednika', BRAKI.wymuszanie),
  );
  return rozwiniecie.element;
}

/**
 * „Pokaż wystąpienia" — `translate.glossary.occurrences`.
 *
 * Wykaz pusty jest tu wynikiem, nie pustką okna: rdzeń odpowiedział i nie
 * znalazł terminu. Odmowa czyści wykaz i mówi wprost, że wystąpień nie
 * sprawdzono — te dwa stany nie mogą wyglądać tak samo.
 *
 * Zdanie nazywa termin, którego rdzeń szukał: bierze go z pola `term`
 * odpowiedzi (`TranslateGlossaryOccurrencesResponse`), a nie z żądania. Zdanie
 * zbudowane z żądania przypisywałoby rdzeniowi przeszukanie o zakresie, którego
 * okno nie widziało; rozbieżność echa jest odmową, bo znaczy, że szukano czegoś
 * innego.
 */
async function pokazWystapienia(
  zrodlo: ZrodloGlosariusza,
  termin: GlossaryTerm,
  wykaz: HTMLElement,
  okno: StanOkna,
  odpowiedz: WierszOdpowiedzi,
): Promise<void> {
  okno.ladowanie(`Rdzeń szuka wystąpień terminu „${termin.source}" w tekście źródłowym.`);
  odpowiedz.pokaz(`Szukanie wystąpień terminu „${termin.source}"…`, true);
  const wynik = await zrodlo.wystapienia(termin.source);
  if (!wynik.udany || wynik.wynik === undefined) {
    wykaz.replaceChildren();
    const zdanie = opisOdmowy('Wystąpienia terminu', wynik.blad?.code, wynik.blad?.message);
    odpowiedz.pokaz(zdanie, false);
    okno.blad(zdanie);
    return;
  }
  const rozbiezne = rozbieznoscOdpowiedzi('Wystąpienia terminu', [
    { nazwa: 'szukany termin', zamowione: termin.source, oddane: wynik.wynik.term },
  ]);
  if (rozbiezne !== null) {
    wykaz.replaceChildren();
    odpowiedz.pokaz(rozbiezne, false);
    okno.blad(rozbiezne);
    return;
  }
  const szukany = wynik.wynik.term;
  wykaz.replaceChildren(...wynik.wynik.occurrences.map(wierszWystapienia));
  okno.gotowe();
  odpowiedz.pokaz(
    wynik.wynik.occurrences.length === 0
      ? `Rdzeń nie znalazł wystąpień terminu „${szukany}".`
      : `Wystąpienia terminu „${szukany}": ${wynik.wynik.occurrences.length}.`,
    true,
  );
}

function wierszWystapienia(tresc: string): HTMLElement {
  const element = document.createElement('li');
  element.className = 'mt-wystapienia__wiersz';
  element.textContent = tresc;
  return element;
}

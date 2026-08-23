import { utworzPanelRodzin } from './panel-rodzin';
import { sekcjePorzadku } from './sekcje-rodzin';
import type { BrowserNote } from '../../../../shared/contract';
import { pokazKomunikat } from '../../aplikacja/komunikaty';
import { utworzDymekObjasnienia } from '../../komponenty/dymek';
import { utworzNaglowekOkna } from '../../komponenty/naglowek-okna';
import { poleTekstowe, utworzWierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import { KLASA_PRZYCISKU, przyciskCzynnosci } from './przyciski-browser';
import { utworzCzynnosciNotatek, type CzynnosciNotatek } from './czynnosci-notatek';
import {
  KLASY_DYMKA,
  KODY_OKIEN,
  OBJASNIENIA,
  POZYCJE_BEZ_OBSLUGI,
  STANY_PUSTE,
  WYKAZY,
  nazwaKlasyfikacji,
} from './etykiety-browser';
import { utworzFormularzNotatki, type FormularzNotatki } from './formularz-notatki';
import { utworzStanOkna, type StanOkna } from './stan-okna';
import type { StanPrzegladania } from './stan-przegladania';
import { utworzSterWyboru } from './ster-wyboru';
import { utworzWierszNotatki, type AkcjeNotatki } from './wiersz-notatki';

/**
 * Notes Panel — okno **pomocnicze** modułu Browser. Notatka wiąże się
 * z konkretnym źródłem z Sources Panel i zasila Library Explorer.
 *
 * Jedna odpowiedzialność: złożenie okna i wykaz notatek. Formularz jest
 * w `formularz-notatki.ts`, rozmowa z rdzeniem w `czynnosci-notatek.ts`,
 * jeden wiersz w `wiersz-notatki.ts`.
 *
 * Edycja nie udaje zmiany: `browser.note.update` kontrakt niesie, ale to okno
 * jeszcze jej nie wywołuje. Pozycja „Edytuj" wypełnia formularz treścią notatki
 * i mówi wprost — powodem wziętym z odczytu wykazu komend rdzenia — że zapis
 * utworzy notatkę nową.
 */
export interface OknoNotesPanel {
  element: HTMLElement;
  odswiez(): void;
  /** Przenosi zaznaczony fragment strony do formularza notatki. */
  przygotujZFragmentu(fragment: string): void;
}

export function utworzOknoNotesPanel(stan: StanPrzegladania): OknoNotesPanel {
  const okno = utworzStanOkna(STANY_PUSTE.notatki);
  const odpowiedz = utworzWierszOdpowiedzi();
  const powiedz = (tresc: string, ok: boolean): void => odpowiedz.pokaz(tresc, ok);
  const formularz = utworzFormularzNotatki();

  const lista = document.createElement('ul');
  lista.className = 'mb-notatki__lista';

  const szukane = poleTekstowe({
    etykieta: 'Szukaj w treści notatek',
    podpowiedz: 'treść albo cytat',
  });
  const widok = utworzSterWyboru({
    nastawa: 'Widok wykazu',
    pozycje: WIDOKI_NOTATEK,
    klasa: 'mb-notatki__ster',
    podpis: false,
    naZmiane: () => odswiez(),
  });

  const uwaga = document.createElement('p');
  uwaga.className = 'dn-pole-opis mb-uwaga';

  const akcje = utworzAkcjeNotatek({
    stan,
    formularz,
    czynnosci: utworzCzynnosciNotatek(stan, powiedz),
    powiedz,
    odswiez: () => odswiez(),
  });

  // Zestawy źródeł, wątki notatek i zmiana notatki — oznaczenia utrwalane
  // w rdzeniu, więc przeżywają przeładowanie karty.
  const rodziny = utworzPanelRodzin(sekcjePorzadku(stan));

  okno.tresc.append(lista, rodziny.element, odpowiedz.element);

  const filtry = document.createElement('div');
  filtry.className = 'mb-panel__pasek';
  filtry.append(
    szukane.element,
    widok.element,
    utworzDymekObjasnienia(OBJASNIENIA.klasyfikacjaNotatki, KLASY_DYMKA),
  );
  szukane.kontrolka.addEventListener('input', () => odswiez());

  const element = document.createElement('section');
  element.className = 'mb-okno mb-okno--pomocnicze';
  element.dataset['okno'] = KODY_OKIEN.notatki;
  element.setAttribute('aria-label', 'Notes Panel — notatki powiązane ze źródłami');
  element.append(
    utworzNaglowekOkna({ tytul: 'Notes Panel', klasa: 'mb-okno__naglowek' }),
    formularz.element,
    utworzPasekNotatek(akcje),
    filtry,
    uwaga,
    okno.element,
  );

  function odswiez(): void {
    formularz.ustawZrodla(stan.zebrane.zrodla());
    const wszystkie = stan.zebrane.notatki();
    const widoczne = przefiltruj(wszystkie, szukane.kontrolka.value);
    uwaga.textContent = zdanieUwagi(stan, wszystkie.length, widoczne.length);
    lista.replaceChildren(...wykazNotatek(stan, widoczne, widok.wartosc(), akcje));
    nanieStan(okno, stan, widoczne.length);
  }

  okno.ustawPonowienie(odswiez);

  return {
    element,
    odswiez,
    przygotujZFragmentu: (fragment) => formularz.ustawCytat(fragment),
  };
}

/** Czynności Operatora w oknie notatek: panel akcji okna i pozycje wykazu. */
interface AkcjeOknaNotatek extends AkcjeNotatki {
  /** Zapis formularza w rdzeniu; pusta treść nie wychodzi z okna. */
  zapisz(): Promise<void>;
  /** Ponowny odczyt wykazów z rdzenia na żądanie Operatora. */
  zaciagnij(): Promise<void>;
  /** Przekazanie wykazu notatek do modułu wskazanego w formularzu. */
  przekaz(): Promise<void>;
}

/** Wszystko, czego czynności potrzebują — bez sięgania do wnętrza okna. */
interface ZapleczeNotatek {
  stan: StanPrzegladania;
  formularz: FormularzNotatki;
  czynnosci: CzynnosciNotatek;
  powiedz(tresc: string, powodzenie: boolean): void;
  odswiez(): void;
}

/**
 * Co dzieje się po naciśnięciu: zapis, edycja, przypięcie, odczyt i przekazanie.
 *
 * Stoi osobno od składania okna: tam mieszka układ elementów, tutaj reguły
 * rozmowy z rdzeniem i treść odpowiedzi.
 */
function utworzAkcjeNotatek(zaplecze: ZapleczeNotatek): AkcjeOknaNotatek {
  const { stan, formularz, czynnosci, powiedz, odswiez } = zaplecze;
  return {
    async zapisz() {
      const zapis = formularz.odczytaj();
      if (zapis.tresc === '') {
        formularz.pokazBlad('Treść notatki jest wymagana — rdzeń odmówi zapisu pustej notatki.');
        powiedz('Notatka nie została wysłana: pole treści jest puste.', false);
        return;
      }
      formularz.pokazBlad('');
      if (await czynnosci.dodaj(zapis)) odswiez();
    },

    edytuj(notatka) {
      const powiazanie = formularz.wypelnij(notatka);
      const oZmianie = zdanieOZmianie(stan);
      pokazKomunikat({ tytul: 'Edycja notatki', tresc: oZmianie, waga: 'ostrz' });
      // Źródła notatki może już nie być w wykazie okna — ster wraca wtedy do
      // „bez powiązania" i okno mówi to wprost, zamiast pozwolić zapisać
      // notatkę o cicho zmienionym powiązaniu.
      const oPowiazaniu = powiazanie
        ? ''
        : ' Źródła tej notatki nie ma w wykazie okna — ster stanął na „bez powiązania".';
      powiedz(`Treść notatki wczytana do formularza. ${oZmianie}${oPowiazaniu}`, false);
    },

    przypnij(notatka) {
      const teraz = stan.zebrane.przelaczPrzypiecie(notatka.id);
      powiedz(teraz ? 'Notatka przypięta na górze wykazu.' : 'Notatka odpięta.', true);
      odswiez();
    },

    otworzZrodlo: (notatka) => void czynnosci.otworzZrodlo(notatka),

    sklasyfikuj(notatka, kod) {
      stan.zebrane.ustawKlasyfikacje(notatka.id, kod);
      powiedz(
        kod === ''
          ? `Klasyfikacja notatki ${notatka.id} zdjęta (oznaczenie żyje w tej karcie).`
          : `Notatka ${notatka.id} oznaczona jako ${nazwaKlasyfikacji(kod)} — oznaczenie żyje w tej karcie, rdzeń go nie zna.`,
        true,
      );
      odswiez();
    },

    przekaz: () => czynnosci.przekaz(formularz.modulDocelowy()),

    async zaciagnij() {
      powiedz('Odczyt wykazu notatek z rdzenia…', true);
      await stan.zaciagnijZebrane();
      const uwagi = stan.powodZebranego();
      powiedz(
        uwagi === '' ? `Wykaz zaciągnięty: ${stan.zebrane.notatki().length} notatek.` : uwagi,
        uwagi === '',
      );
    },
  };
}

/** Panel akcji okna: trzy pozycje wykazu narzędzi Notes Panel. */
function utworzPasekNotatek(akcje: AkcjeOknaNotatek): HTMLElement {
  const element = document.createElement('div');
  element.className = 'mb-panel__pasek';
  element.append(
    przyciskCzynnosci('+ Nowa notatka', KLASA_PRZYCISKU.glowny, () => void akcje.zapisz()),
    przyciskCzynnosci('Odśwież z rdzenia', KLASA_PRZYCISKU.zarys, () => void akcje.zaciagnij()),
    przyciskCzynnosci(
      '→ Wyślij do Research/Library',
      KLASA_PRZYCISKU.zarys,
      () => void akcje.przekaz(),
    ),
    utworzDymekObjasnienia(OBJASNIENIA.przekazanie, KLASY_DYMKA),
  );
  return element;
}

/**
 * Powód, dla którego zapis zakłada notatkę nową zamiast zmieniać zastaną.
 *
 * Zdanie bierze się z odczytu wykazu komend rdzenia, nie z napisu w module:
 * `browser.note.update` stoi w kontrakcie, więc od dnia, w którym rdzeń dostanie
 * jej uchwyt, powód ma brzmieć inaczej — i zabrzmi, bez wchodzenia w ten plik.
 */
function zdanieOZmianie(stan: StanPrzegladania): string {
  return stan.pokrycie.zdanie(
    POZYCJE_BEZ_OBSLUGI.zmianaNotatki.komenda,
    POZYCJE_BEZ_OBSLUGI.zmianaNotatki.czynnosc,
  );
}

/**
 * Trzy porządki wykazu z opracowania modułu.
 *
 * Grupowanie „według wątku" nie stoi w wykazie: wątków tematycznych własnych
 * nie niesie ani kontrakt, ani ten moduł, więc pozycja bez treści byłaby
 * obietnicą porządku, którego nie ma czym zbudować.
 */
const WIDOKI_NOTATEK = [
  {
    wartosc: 'chronologia',
    etykieta: 'Chronologicznie',
    opis: 'Najświeższe na górze; przypięte zawsze przed pozostałymi.',
  },
  {
    wartosc: 'zrodlo',
    etykieta: 'Według źródła',
    opis: 'Notatki zebrane pod stroną, której dotyczą.',
  },
  {
    wartosc: 'klasyfikacja',
    etykieta: 'Według klasyfikacji',
    opis: 'Obserwacje, cytaty, pytania otwarte i wnioski osobno.',
  },
];

/** Zawężenie wykazu frazą szukaną w treści notatki i w cytacie. */
function przefiltruj(notatki: readonly BrowserNote[], fraza: string): BrowserNote[] {
  const szukana = fraza.trim().toLowerCase();
  if (szukana === '') return [...notatki];
  return notatki.filter((notatka) =>
    `${notatka.content} ${notatka.quote ?? ''}`.toLowerCase().includes(szukana),
  );
}

/**
 * Wykaz notatek w wybranym porządku. Przy grupowaniu przed każdą grupą staje
 * jej nagłówek — bez niego kolejność zmieniałaby się bez widocznej przyczyny.
 */
function wykazNotatek(
  stan: StanPrzegladania,
  notatki: readonly BrowserNote[],
  porzadek: string,
  akcje: AkcjeNotatki,
): HTMLElement[] {
  if (porzadek === 'chronologia') {
    return notatki.map((notatka) => wierszNotatki(stan, notatka, akcje));
  }
  const grupy = new Map<string, BrowserNote[]>();
  for (const notatka of notatki) {
    const nazwa =
      porzadek === 'zrodlo' ? nazwaZrodla(stan, notatka) : nazwaKlasyfikacji(stan.zebrane.klasyfikacja(notatka.id));
    const zastane = grupy.get(nazwa);
    if (zastane === undefined) grupy.set(nazwa, [notatka]);
    else zastane.push(notatka);
  }
  return [...grupy].flatMap(([nazwa, wpisy]) => [
    naglowekGrupy(`${nazwa} · ${wpisy.length}`),
    ...wpisy.map((notatka) => wierszNotatki(stan, notatka, akcje)),
  ]);
}

/** Nazwa źródła notatki albo zdanie o jego braku — klucz grupowania. */
function nazwaZrodla(stan: StanPrzegladania, notatka: BrowserNote): string {
  const zrodlo = stan.zebrane.zrodlo(notatka.sourceId ?? '');
  if (zrodlo === null) return 'bez powiązania ze źródłem';
  const tytul = (zrodlo.title ?? '').trim();
  return tytul === '' ? zrodlo.url : tytul;
}

function naglowekGrupy(napis: string): HTMLElement {
  const element = document.createElement('li');
  element.className = 'mb-notatki__grupa';
  element.textContent = napis;
  return element;
}

/** Jedna notatka wraz z powiązanym źródłem, przypięciem i klasyfikacją. */
function wierszNotatki(
  stan: StanPrzegladania,
  notatka: BrowserNote,
  akcje: AkcjeNotatki,
): HTMLElement {
  return utworzWierszNotatki(
    notatka,
    stan.zebrane.zrodlo(notatka.sourceId ?? ''),
    stan.zebrane.czyPrzypieta(notatka.id),
    stan.zebrane.klasyfikacja(notatka.id),
    akcje,
  );
}

/**
 * Uwaga panelu mówi o pochodzeniu wykazu i o braku komendy zmiany; gdy
 * ostatni odczyt miał powód do zgłoszenia, mówi najpierw jego. Zawężony wykaz
 * mówi także, ile pozycji filtr ukrył.
 */
function zdanieUwagi(stan: StanPrzegladania, wszystkich: number, widocznych: number): string {
  const uwagi = stan.powodZebranego();
  const podstawa = `${uwagi === '' ? WYKAZY.notatki : uwagi} ${zdanieOZmianie(stan)}`;
  if (widocznych === wszystkich) return podstawa;
  return `${podstawa} Filtr pokazuje ${widocznych} z ${wszystkich} notatek.`;
}

/**
 * Trzy stany obowiązkowe wykazu notatek — ten sam zestaw reguł, co w panelu
 * źródeł, bo oba wykazy przychodzą jednym zaciągnięciem.
 *
 * Kolejność pytań: czekanie, odmowa, pustka. Odmowa wykazu jest błędem panelu
 * tylko przy pustym wykazie — z notatkami na ekranie wpisy zostają, a powód
 * idzie zdaniem przy wykazie.
 */
function nanieStan(okno: StanOkna, stan: StanPrzegladania, pozycji: number): void {
  if (stan.faza() === 'odczyt') {
    okno.ladowanie('Rdzeń ustala okno przeglądarki tej sesji…');
    return;
  }
  if (stan.stanZebranego() === 'odczyt') {
    okno.ladowanie(stan.powodZebranego());
    return;
  }
  if (stan.faza() === 'blad') {
    okno.blad(stan.powod());
    return;
  }
  if (stan.stanZebranego() === 'blad' && pozycji === 0) {
    okno.blad(stan.powodZebranego());
    return;
  }
  if (pozycji === 0) {
    // Opis, z którym okno powstało — jedno brzmienie pustki na okno.
    okno.pusteZOpisu();
    return;
  }
  okno.gotowe();
}

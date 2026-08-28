import type { BrowserSource } from '../../../../shared/contract';
import { utworzDymekObjasnienia } from '../../komponenty/dymek';
import { utworzNaglowekOkna } from '../../komponenty/naglowek-okna';
import {
  pobierzPlik,
  poleLogiczne,
  poleTekstowe,
  utworzWierszOdpowiedzi,
} from '../../modele/kontrolki-formularza';
import { KLASA_PRZYCISKU, przyciskCzynnosci } from './przyciski-browser';
import { utworzCzynnosciZrodel, type CzynnosciZrodel } from './czynnosci-zrodel';
import { KLASY_DYMKA, KODY_OKIEN, OBJASNIENIA, STANY_PUSTE, WYKAZY } from './etykiety-browser';
import { utworzStanOkna, type StanOkna } from './stan-okna';
import type { StanPrzegladania } from './stan-przegladania';
import { utworzSterWyboru } from './ster-wyboru';
import { utworzWierszZrodla } from './wiersz-zrodla';
import { notacja, wypisz, NOTACJE_WYPISU } from './wypis-zrodel';

/**
 * Sources Panel — okno pomocnicze modułu Browser, narastające w miarę
 * przeglądania kolejnych stron, zasilające moduł Research.
 */
export interface OknoSourcesPanel {
  element: HTMLElement;
  odswiez(): void;
}

/** Zdanie okna o skutku wykonanej czynności — wiersz odpowiedzi wyświetlany zawsze pod wykazem źródeł panelu. */
type Powiedz = (tresc: string, powodzenie: boolean) => void;

export function utworzOknoSourcesPanel(stan: StanPrzegladania): OknoSourcesPanel {
  const okno = utworzStanOkna(STANY_PUSTE.zrodla);
  const odpowiedz = utworzWierszOdpowiedzi();
  const powiedz: Powiedz = (tresc, ok) => odpowiedz.pokaz(tresc, ok);
  const czynnosci = utworzCzynnosciZrodel(stan, powiedz);
  const wybrane = new Set<string>();

  const adres = poleTekstowe({ etykieta: 'Adres źródła', podpowiedz: 'https://' });
  const nazwa = poleTekstowe({ etykieta: 'Tytuł źródła' });
  const kluczowe = poleLogiczne({ etykieta: 'Źródło kluczowe' });
  const szukane = poleTekstowe({
    etykieta: 'Szukaj w wykazie',
    podpowiedz: 'tytuł, domena albo data',
  });
  const tylkoKluczowe = poleLogiczne({ etykieta: 'Pokaż wyłącznie źródła kluczowe' });
  const samoczynnie = poleLogiczne({ etykieta: 'Dopisuj źródła samoczynnie' });
  samoczynnie.kontrolka.checked = true;
  const format = utworzSterWyboru({
    nastawa: 'Notacja wypisu',
    pozycje: NOTACJE_WYPISU,
    klasa: 'mb-zrodla__ster',
    podpis: false,
  });
  const lista = document.createElement('ul');
  lista.className = 'mb-zrodla__lista';
  const wypis = utworzWypisZrodel(stan, powiedz, () => format.wartosc());

  /** Migawka, z której źródło już dopisano — próg samoczynnego dopisywania. */
  let ostatniaMigawka = '';

  async function dodaj(url: string, tytul: string, klucz: boolean): Promise<void> {
    if (await czynnosci.dodaj(url, tytul, klucz)) odswiez();
  }

  function dodajZBiezacej(): void {
    const wpis = zrodloBiezacejStrony(stan, powiedz);
    if (wpis !== null) void dodaj(wpis.url, wpis.tytul, kluczowe.kontrolka.checked);
  }

  // Samoczynne dopisanie źródła po przejściu na nową stronę, według progu identyfikatora migawki.
  function dopiszSamoczynnie(): void {
    if (!samoczynnie.kontrolka.checked) return;
    const migawka = stan.migawka();
    if (migawka === null || migawka.id === ostatniaMigawka) return;
    ostatniaMigawka = migawka.id;
    if (stan.zebrane.zrodla().some((wpis) => wpis.url === migawka.url)) return;
    void dodaj(migawka.url, migawka.title ?? '', false);
  }

  const pasek = utworzPasek({
    dodaj: () => void dodaj(adres.kontrolka.value, nazwa.kontrolka.value, kluczowe.kontrolka.checked),
    zBiezacej: dodajZBiezacej,
    odswiezZRdzenia: () => void zaciagnijWykaz(stan, powiedz),
    eksportuj: wypis.wykonaj,
    doResearch: () => void czynnosci.przekaz([...wybrane]),
    notacja: format.element,
  });

  const uwaga = document.createElement('p');
  uwaga.className = 'dn-pole-opis mb-uwaga';

  okno.tresc.append(lista, wypis.element, odpowiedz.element);

  const element = document.createElement('section');
  element.className = 'mb-okno mb-okno--pomocnicze';
  element.dataset['okno'] = KODY_OKIEN.zrodla;
  element.setAttribute('aria-label', 'Sources Panel — wykaz źródeł przeglądania');
  const filtry = document.createElement('div');
  filtry.className = 'mb-panel__pasek';
  filtry.append(
    szukane.element,
    tylkoKluczowe.element,
    samoczynnie.element,
    utworzDymekObjasnienia(OBJASNIENIA.samoczynneZrodla, KLASY_DYMKA),
  );
  szukane.kontrolka.addEventListener('input', () => odswiez());
  tylkoKluczowe.kontrolka.addEventListener('change', () => odswiez());

  element.append(
    utworzNaglowekOkna({ tytul: 'Sources Panel', klasa: 'mb-okno__naglowek' }),
    adres.element,
    nazwa.element,
    kluczowe.element,
    utworzDymekObjasnienia(OBJASNIENIA.kluczowe, KLASY_DYMKA),
    pasek,
    filtry,
    uwaga,
    okno.element,
  );

  function odswiez(): void {
    dopiszSamoczynnie();
    const widoczne = przefiltruj(stan.zebrane.zrodla(), {
      fraza: szukane.kontrolka.value,
      tylkoKluczowe: tylkoKluczowe.kontrolka.checked,
    });
    uwaga.textContent = zdanieUwagi(stan, stan.zebrane.zrodla().length, widoczne.length);
    const oznaczKluczowe = (wpis: BrowserSource): void => {
      void dodaj(wpis.url, wpis.title ?? '', true);
    };
    lista.replaceChildren(...wierszeZrodel(widoczne, wybrane, czynnosci, oznaczKluczowe));
    nanieStan(okno, stan, widoczne.length);
  }

  okno.ustawPonowienie(odswiez);
  return { element, odswiez };
}

/**
 * Adres i tytuł strony widocznej w Browser Window.
 *
 * Bez migawki nie ma czego zapisać — okno mówi o tym wprost zamiast wysyłać
 * do rdzenia żądanie z pustym adresem.
 */
function zrodloBiezacejStrony(
  stan: StanPrzegladania,
  powiedz: Powiedz,
): { url: string; tytul: string } | null {
  const migawka = stan.migawka();
  if (migawka === null) {
    powiedz('Brak migawki — najpierw przejdź do strony w Browser Window.', false);
    return null;
  }
  return { url: migawka.url, tytul: migawka.title ?? '' };
}

/** Ponowny odczyt wykazów źródeł wprost z rdzenia na żądanie operatora, nie z pamięci przechowanej w karcie. */
async function zaciagnijWykaz(stan: StanPrzegladania, powiedz: Powiedz): Promise<void> {
  powiedz('Odczyt wykazu źródeł z rdzenia…', true);
  await stan.zaciagnijZebrane();
  const uwagi = stan.powodZebranego();
  powiedz(
    uwagi === '' ? `Wykaz zaciągnięty: ${stan.zebrane.zrodla().length} źródeł.` : uwagi,
    uwagi === '',
  );
}

/**
 * Uwaga panelu mówi o pochodzeniu wykazu i o ostatnim odczycie z rdzenia.
 * Przy zawężonym wykazie mówi też, ile pozycji filtr ukrył przed operatorem.
 */
function zdanieUwagi(stan: StanPrzegladania, wszystkich: number, widocznych: number): string {
  const uwagi = stan.powodZebranego();
  const podstawa = uwagi === '' ? WYKAZY.zrodla : uwagi;
  if (widocznych === wszystkich) return podstawa;
  return `${podstawa} Filtr pokazuje ${widocznych} z ${wszystkich} pozycji wykazu.`;
}

/** Zawężenie wykazu źródeł: fraza szukana w tytule, adresie i dacie strony, wraz z samą istotnością źródła. */
function przefiltruj(
  zrodla: readonly BrowserSource[],
  filtr: { fraza: string; tylkoKluczowe: boolean },
): BrowserSource[] {
  const fraza = filtr.fraza.trim().toLowerCase();
  return zrodla.filter((zrodlo) => {
    if (filtr.tylkoKluczowe && zrodlo.key !== true) return false;
    if (fraza === '') return true;
    const data = new Date(zrodlo.createdAt).toISOString().slice(0, 10);
    return `${zrodlo.title ?? ''} ${zrodlo.url} ${data}`.toLowerCase().includes(fraza);
  });
}

/** Wykaz źródeł okna wraz ze znacznikiem zaznaczenia, do przekazania wybranych źródeł do innego modułu. */
function wierszeZrodel(
  zrodla: readonly BrowserSource[],
  wybrane: Set<string>,
  czynnosci: CzynnosciZrodel,
  oznaczKluczowe: (zrodlo: BrowserSource) => void,
): HTMLElement[] {
  return zrodla.map((zrodlo) =>
    utworzWierszZrodla(zrodlo, wybrane.has(zrodlo.id), {
      otworz: (wpis) => void czynnosci.otworz(wpis),
      migawka: czynnosci.pokazMigawke,
      oznaczKluczowe,
      usun: czynnosci.usun,
      przelaczWybor: (wpis) => {
        if (wybrane.has(wpis.id)) wybrane.delete(wpis.id);
        else wybrane.add(wpis.id);
      },
    }),
  );
}

/** Czynności panelu akcji Sources Panel wraz ze sterem notacji wypisu bibliografii z wykazu źródeł okna. */
interface AkcjePaska {
  dodaj(): void;
  zBiezacej(): void;
  odswiezZRdzenia(): void;
  eksportuj(): void;
  doResearch(): void;
  /** Ster notacji wypisu — buduje go okno, bo to ono zna wybraną wartość. */
  notacja: HTMLElement;
}

/** Panel akcji okna: pięć pozycji wykazu narzędzi Sources Panel dostępnych bezpośrednio dla operatora sesji. */
function utworzPasek(akcje: AkcjePaska): HTMLElement {
  const element = document.createElement('div');
  element.className = 'mb-panel__pasek';
  element.append(
    przyciskCzynnosci('+ Dodaj ręcznie', KLASA_PRZYCISKU.glowny, akcje.dodaj),
    przyciskCzynnosci('Dodaj z bieżącej strony', KLASA_PRZYCISKU.zarys, akcje.zBiezacej),
    przyciskCzynnosci('Odśwież z rdzenia', KLASA_PRZYCISKU.zarys, akcje.odswiezZRdzenia),
    akcje.notacja,
    przyciskCzynnosci('Eksportuj listę', KLASA_PRZYCISKU.zarys, akcje.eksportuj),
    utworzDymekObjasnienia(OBJASNIENIA.eksportZrodel, KLASY_DYMKA),
    przyciskCzynnosci('→ Wyślij do Research', KLASA_PRZYCISKU.zarys, akcje.doResearch),
    utworzDymekObjasnienia(OBJASNIENIA.przekazanie, KLASY_DYMKA),
  );
  return element;
}

/**
 * Trzy stany obowiązkowe wykazu źródeł: czekanie, odmowa, pustka, w tej
 * kolejności. Odmowa jest błędem okna tylko przy pustym wykazie, bo pozycje
 * na ekranie zostają, a powód idzie zdaniem przy wykazie.
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
    // Ten sam opis, z którym okno powstało — jedno brzmienie pustki na okno.
    okno.pusteZOpisu();
    return;
  }
  okno.gotowe();
}

/** Wypis wykazu źródeł — jedyne miejsce, gdzie okno oddaje zebrane źródła tekstem oraz plikiem na dysku. */
interface WypisZrodel {
  element: HTMLElement;
  /** Odtwarza wypis z wykazu bieżącego, pobiera plik i mówi, ile pozycji poszło. */
  wykonaj(): void;
}

/**
 * Wypis staje na widoku i schodzi na dysk jednym naciśnięciem — dwa różne
 * zastosowania tej samej treści bibliografii, złożone w jedną czynność.
 */
function utworzWypisZrodel(
  stan: StanPrzegladania,
  powiedz: Powiedz,
  wybranaNotacja: () => string,
): WypisZrodel {
  const element = document.createElement('pre');
  element.className = 'mb-zrodla__eksport';
  element.hidden = true;
  return {
    element,
    wykonaj() {
      const zrodla = stan.zebrane.zrodla();
      if (zrodla.length === 0) {
        element.hidden = true;
        powiedz('Nie ma czego eksportować — wykaz jest pusty.', false);
        return;
      }
      const postac = notacja(wybranaNotacja());
      const tresc = wypisz(zrodla, postac.wartosc);
      element.hidden = false;
      element.textContent = tresc;
      pobierzPlik(`zrodla-sesji.${postac.rozszerzenie}`, tresc, postac.rodzaj);
      powiedz(
        `Wypisano ${zrodla.length} źródeł w notacji ${postac.etykieta} ` +
          `i pobrano plik zrodla-sesji.${postac.rozszerzenie}.`,
        true,
      );
    },
  };
}

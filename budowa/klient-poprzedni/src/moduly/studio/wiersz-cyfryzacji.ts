import { StudioIngestState, type StudioIngestItem } from '../../../../shared/contract';

/**
 * Jeden wiersz kolejki wczytywania narzędziowni cyfryzacji.
 *
 * Wiersz nie wywołuje niczego sam: oddaje przyciski, a czynności podpina okno.
 * Dzięki temu rozmowa z rdzeniem zostaje w jednym miejscu, a wiersz odpowiada
 * wyłącznie za to, co widać.
 *
 * ── Stan jest treścią, nie barwą ────────────────────────────────────────────
 * Nazwa stanu idzie do `data-stan` i do plakietki słownej naraz, więc pozycja
 * odmówiona jest rozpoznawalna także bez odczytu barwy. Stan przychodzi z rdzenia
 * jako `StudioIngestState` i ma PIĘĆ wartości, nie cztery: „ponowienie" jest
 * stanem osobnym i nazywa rzecz, której ani „gotowa", ani „odmowa" nie opisuje —
 * rozpoznanie wypadło poniżej progu pewności i pozycja wraca do rozpoznania.
 */
export interface WierszCyfryzacji {
  element: HTMLElement;
  rozpoznaj: HTMLButtonElement;
  wskaz: HTMLButtonElement;
  przyjmij: HTMLButtonElement;
}

/** Plakietka stanu: napis dla Operatora i odmiana biblioteczna. */
const PLAKIETKA_STANU: Record<StudioIngestState, { napis: string; odmiana: string }> = {
  [StudioIngestState.Oczekuje]: { napis: 'oczekuje', odmiana: '' },
  [StudioIngestState.Przetwarzanie]: {
    napis: 'rozpoznawanie w toku',
    odmiana: 'dn-plakietka--informacja',
  },
  [StudioIngestState.Gotowa]: { napis: 'gotowa', odmiana: 'dn-plakietka--sukces' },
  [StudioIngestState.Ponowienie]: {
    napis: 'ponowienie — poniżej progu pewności',
    odmiana: 'dn-plakietka--ostrz',
  },
  [StudioIngestState.Odmowa]: { napis: 'odmowa rdzenia', odmiana: 'dn-plakietka--blad' },
};

/** Zdanie o wyniku pozycji — skąd wziął się tekst, ile go jest i jak pewny. */
function opiszWynik(pozycja: StudioIngestItem): string {
  if (pozycja.state === StudioIngestState.Odmowa) {
    return pozycja.failureReason ?? 'Rdzeń odmówił, ale powodu nie podał.';
  }
  const tekst = pozycja.text ?? '';
  if (tekst === '' && pozycja.state !== StudioIngestState.Gotowa) return '—';
  if (tekst === '') {
    return (
      'Rdzeń oddał tekst PUSTY — to odpowiedź, nie odmowa. Materiał nie ma warstwy tekstowej, ' +
      'a rozpoznanie pisma nic z niego nie odczytało.'
    );
  }
  const zrodloTekstu =
    pozycja.usedOcr === true ? 'tekst z rozpoznania pisma' : 'tekst z warstwy tekstowej';
  const strony = pozycja.pages === undefined ? '' : ` · stron ${pozycja.pages}`;
  const pewnosc =
    pozycja.confidence === undefined ? '' : ` · pewność ${pozycja.confidence}`;
  return `${zrodloTekstu} · ${tekst.length} znaków${strony}${pewnosc}`;
}

/** Wskazanie materiału widziane przez rdzeń — ścieżka albo zasób magazynu. */
function opiszWskazanie(pozycja: StudioIngestItem): { tresc: string; rodzaj: string; opis: string } {
  if (pozycja.assetId !== undefined && pozycja.assetId !== '') {
    return {
      tresc: pozycja.assetId,
      rodzaj: 'zasób magazynu',
      opis: 'Zasób magazynu rdzenia — treść stoi pod sumą kontrolną.',
    };
  }
  if (pozycja.sourcePath !== undefined && pozycja.sourcePath !== '') {
    return {
      tresc: pozycja.sourcePath,
      rodzaj: 'ścieżka rdzenia',
      opis: 'Ścieżka widziana przez rdzeń; klient dysku nie czyta.',
    };
  }
  return {
    tresc: pozycja.id,
    rodzaj: 'pozycja bez wskazania',
    opis:
      'Rdzeń nie podał ani ścieżki, ani zasobu — pozycja powstała inaczej, na przykład z pobrania ' +
      'strony sieciowej albo z urządzenia wejściowego.',
  };
}

export function utworzWierszCyfryzacji(
  pozycja: StudioIngestItem,
  wskazana: boolean,
  wToku: boolean,
): WierszCyfryzacji {
  const opisany = opiszWskazanie(pozycja);

  const wskazanie = document.createElement('td');
  wskazanie.className = 'ms-kolejka__wskazanie';
  wskazanie.textContent = opisany.tresc;
  wskazanie.title = opisany.opis;

  const rodzaj = document.createElement('td');
  rodzaj.textContent = opisany.rodzaj;

  const stan = document.createElement('td');
  const plakietka = document.createElement('span');
  const opisStanu = PLAKIETKA_STANU[pozycja.state];
  plakietka.className =
    opisStanu.odmiana === '' ? 'dn-plakietka' : `dn-plakietka ${opisStanu.odmiana}`;
  plakietka.textContent = wToku ? 'wywołanie w toku' : opisStanu.napis;
  stan.append(plakietka);

  const wynik = document.createElement('td');
  wynik.className = 'ms-kolejka__wynik';
  wynik.textContent = opiszWynik(pozycja);

  const rozpoznaj = przyciskWiersza(
    pozycja.state === StudioIngestState.Ponowienie ? 'Rozpoznaj ponownie' : 'Rozpoznaj',
    'dn-btn dn-btn--sm dn-btn--atrament',
    'rozpoznaj',
  );
  const wskaz = przyciskWiersza(
    wskazana ? 'Wskazana' : 'Wskaż',
    'dn-btn dn-btn--sm dn-btn--zarys',
    'wskaz',
  );
  wskaz.setAttribute('aria-pressed', String(wskazana));
  const przyjmij = przyciskWiersza(
    'Przyjmij do edytora',
    'dn-btn dn-btn--sm dn-btn--sygnal',
    'przyjmij',
  );
  przyjmij.title =
    'studio.ingest.item.accept zakłada dokument roboczy WRAZ z pierwszą wersją w repozytorium ' +
    'sesji — nie sam bufor edytora.';

  const czynnosci = document.createElement('td');
  czynnosci.className = 'ms-kolejka__czynnosci';
  czynnosci.append(rozpoznaj, wskaz, przyjmij);

  const element = document.createElement('tr');
  element.dataset['pozycja'] = pozycja.id;
  element.dataset['stan'] = pozycja.state;
  element.dataset['wskazana'] = String(wskazana);
  element.dataset['wToku'] = String(wToku);
  element.append(wskazanie, rodzaj, stan, wynik, czynnosci);

  return { element, rozpoznaj, wskaz, przyjmij };
}

function przyciskWiersza(etykieta: string, klasa: string, czynnosc: string): HTMLButtonElement {
  const kontrolka = document.createElement('button');
  kontrolka.type = 'button';
  kontrolka.className = klasa;
  kontrolka.textContent = etykieta;
  kontrolka.dataset['czynnosc'] = czynnosc;
  return kontrolka;
}

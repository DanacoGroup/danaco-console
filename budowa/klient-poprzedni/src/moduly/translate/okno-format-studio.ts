import type { ExportFormat } from '../../../../shared/contract';
import { opisOdmowyBledu } from '../../komponenty/odmowa';
import {
  poleLogiczne,
  poleTekstowe,
  poleWyboru,
  przycisk,
  przyciskBezKomendy,
  ustawPozycje,
  utworzWierszOdpowiedzi,
  type WierszOdpowiedzi,
} from '../../modele/kontrolki-formularza';
import { eksportuj } from './czynnosci-panelu';
import { FORMATY_EKSPORTU, OBJASNIENIA, PUSTE } from './etykiety-translate';
import { dopnijDymek, naglowekOkna } from './kontrolki-translate';
import {
  WYDLUZENIE_DOMYSLNE,
  pseudolokalizuj,
  zdanieOPseudolokalizacji,
} from './pseudolokalizacja';
import { utworzStanOkna, type StanOkna } from './stan-okna-translate';
import type { StanTranslate } from './stan-translate';
import { oznaczWarstwe, utworzRozwiniecie } from './warstwy-translate';
import type { ZrodloDokumentuTranslate } from './zrodlo-dokumentu-translate';

/**
 * Format Studio — okno robocze modułu Translate, otwierane przy pracy
 * z dokumentem.
 *
 * Okno stoi na styku dwóch obszarów kontraktu i tylko dlatego ma czym pracować:
 * dokument wchodzi komendami obszaru `document` (wydobycie tekstu wraz
 * z rozpoznaniem pisma oraz zamiana formatu), a wychodzi komendą eksportu
 * panelu z obszaru `translate`. Obieg jest więc realny na obu końcach, ale nie
 * jest obiegiem zamkniętym: wydobycie oddaje sam tekst, więc styl, tabela
 * i osadzenie dokumentu wejściowego nie mają jak przetrwać przekładu. Okno mówi
 * to przy wyniku, zamiast obiecywać wierność formatu.
 *
 * Podgląd jest warstwą tekstową, nie renderem strony: kontrakt nie ma komendy
 * rysującej dokument ani porównującej jego układ, więc porównania układów okno
 * nie pokazuje i nie udaje.
 *
 * Tekst wydobyty z dokumentu wchodzi do pola Source Panel, ale nie zapisuje się
 * sam. Zapis źródła jest czynnością Operatora i uruchamia aktualizację
 * wszystkich paneli — wykonanie go bez naciśnięcia byłoby przekładem
 * zamówionym przez okno, nie przez człowieka.
 */
export interface OknoFormatStudio {
  element: HTMLElement;
  odswiez(): void;
}

/** Formaty zamiany dokumentu wymienione w opisie komendy kontraktu. */
const FORMATY_DOKUMENTU: readonly string[] = [
  'markdown',
  'html',
  'docx',
  'odt',
  'pdf',
  'epub',
  'rtf',
  'csv',
];

const BRAKI = {
  ukladIPorownanie:
    'Kontrakt nie ma komendy rysującej dokument ani porównującej jego układ, więc przepełnień ' +
    'i złamanego układu nie ma z czego wyliczyć.',
  zasobyLokalizacyjne:
    'Kontrakt nie zna formatów zasobów lokalizacyjnych ani pojęcia klucza zasobu — nie ma czego ' +
    'wczytać ani czego chronić.',
  formyGramatyczne:
    'Reguł liczby mnogiej i rodzaju kontrakt nie stosuje: nie ma komendy przyjmującej zasób ' +
    'ani reguł języka docelowego.',
  napisy:
    'Formatów napisowych nie ma ani w wykazie eksportu panelu, ani w wykazie zamiany formatu ' +
    'dokumentu, więc taktowania nie ma czym wczytać.',
  dwujezyczny:
    'Eksport panelu wydaje sam przekład w jednym z pięciu formatów kontraktu. Formatu XLIFF ani ' +
    'wariantu dwujęzycznego w wykazie nie ma.',
} as const;

export function utworzOknoFormatStudio(
  stan: StanTranslate,
  dokument: ZrodloDokumentuTranslate,
  wstawDoZrodla: (tekst: string) => void,
): OknoFormatStudio {
  const okno: StanOkna = utworzStanOkna(PUSTE.formaty);
  const odpowiedz = utworzWierszOdpowiedzi();

  const sciezka = poleTekstowe({
    etykieta: 'Ścieżka dokumentu po stronie rdzenia',
    podpowiedz: 'ścieżka pliku widziana przez rdzeń',
    opis: 'Klient plików nie czyta ani nie zapisuje — robi to rdzeń pod wskazaną ścieżką.',
  });

  const rozpoznajPismo = poleLogiczne({
    etykieta: 'Wymuś rozpoznanie pisma, także gdy dokument ma warstwę tekstową',
  });

  const podglad = document.createElement('pre');
  podglad.className = 'mt-formaty__podglad';

  const wydobadz = przycisk('Wydobądź tekst dokumentu', 'dn-btn dn-btn--sm dn-btn--atrament');
  wydobadz.addEventListener('click', () => {
    void wydobadzTekst(
      dokument,
      { sciezka: sciezka.kontrolka.value, rozpoznajPismo: rozpoznajPismo.kontrolka.checked },
      { podglad, okno, odpowiedz, wstawDoZrodla },
    );
  });

  const formatDokumentu = poleWyboru(
    { etykieta: 'Format docelowy dokumentu' },
    FORMATY_DOKUMENTU.map((nazwa) => ({ wartosc: nazwa, etykieta: nazwa })),
  );

  const zamien = przycisk('Zamień format dokumentu', 'dn-btn dn-btn--sm dn-btn--zarys');
  zamien.addEventListener('click', () => {
    void zamienFormat(
      dokument,
      sciezka.kontrolka.value,
      formatDokumentu.kontrolka.value,
      okno,
      odpowiedz,
    );
  });

  const wejscie = document.createElement('div');
  wejscie.className = 'mt-formaty__wejscie';
  oznaczWarstwe(wejscie, 1);
  wejscie.append(
    sciezka.element,
    rozpoznajPismo.element,
    formatDokumentu.element,
    pasek(wydobadz, zamien),
    odpowiedz.element,
    podglad,
  );

  const panel = poleWyboru({ etykieta: 'Panel do wydania' }, []);
  const formatPanelu = poleWyboru(
    { etykieta: 'Format pliku wyniku' },
    FORMATY_EKSPORTU.map((pozycja) => ({ wartosc: pozycja.wartosc, etykieta: pozycja.etykieta })),
  );
  dopnijDymek(formatPanelu.element, OBJASNIENIA.formatEksportu);

  const wydaj = przycisk('Wydaj panel do pliku', 'dn-btn dn-btn--sm dn-btn--atrament');
  wydaj.addEventListener('click', () => {
    void wydajPanel(
      stan,
      panel.kontrolka.value,
      formatPanelu.kontrolka.value as ExportFormat,
      okno,
      odpowiedz,
    );
  });

  const wydanie = document.createElement('div');
  wydanie.className = 'mt-formaty__wydanie';
  oznaczWarstwe(wydanie, 2);
  wydanie.append(
    panel.element,
    formatPanelu.element,
    pasek(wydaj, przyciskBezKomendy('Wydanie dwujęzyczne', BRAKI.dwujezyczny)),
  );

  okno.tresc.append(
    wejscie,
    wydanie,
    menuDokumentu(),
    pseudolokalizacja(stan, okno),
    regulyGramatyczne(),
  );

  const element = document.createElement('section');
  element.className = 'mt-okno mt-okno--robocze';
  element.dataset['okno'] = 'format-studio';
  element.append(naglowekOkna('Format Studio', 'robocze'), okno.element);

  function odswiez(): void {
    ustawPozycje(
      panel.kontrolka,
      stan
        .panelJezykow()
        .map((wpis) => ({ wartosc: wpis.id, etykieta: `${wpis.language} — panel ${wpis.id}` })),
    );
    if (okno.faza() === 'ladowanie' || okno.faza() === 'blad') return;
    if (podglad.textContent === '' && stan.panelJezykow().length === 0) {
      okno.puste(PUSTE.formaty);
      return;
    }
    okno.gotowe();
  }

  odswiez();
  return { element, odswiez };
}

function pasek(...przyciski: readonly HTMLElement[]): HTMLElement {
  const element = document.createElement('div');
  element.className = 'mt-pasek';
  element.append(...przyciski);
  return element;
}

/**
 * `document.text.extract` — jedyne wejście modułu od strony pliku.
 *
 * Zdanie o wyniku rozróżnia dwie drogi, którymi tekst mógł powstać, bo różnią
 * się pewnością: warstwa tekstowa dokumentu jest odczytem, rozpoznanie pisma
 * jest odgadnięciem z pikseli. Odpowiedź mówi, która droga zaszła, i okno tego
 * nie zaciera.
 */
async function wydobadzTekst(
  dokument: ZrodloDokumentuTranslate,
  zadanie: { sciezka: string; rozpoznajPismo: boolean },
  widok: {
    podglad: HTMLElement;
    okno: StanOkna;
    odpowiedz: WierszOdpowiedzi;
    wstawDoZrodla: (tekst: string) => void;
  },
): Promise<void> {
  const { okno, odpowiedz } = widok;
  const sciezka = zadanie.sciezka.trim();
  if (sciezka === '') {
    const zdanie = 'Wskaż ścieżkę dokumentu — żądanie wydobycia tekstu niesie ścieżkę pliku.';
    odpowiedz.pokaz(zdanie, false);
    okno.blad(zdanie);
    return;
  }

  okno.ladowanie('Rdzeń wydobywa tekst ze wskazanego dokumentu.');
  odpowiedz.pokaz('Wydobywanie tekstu dokumentu…', true);
  const wynik = await dokument.wydobadzTekst(sciezka, zadanie.rozpoznajPismo);
  if (!wynik.udany || wynik.wynik === undefined) {
    const zdanie = opisOdmowyBledu('Wydobycie tekstu dokumentu', wynik.blad);
    odpowiedz.pokaz(zdanie, false);
    okno.blad(zdanie);
    return;
  }

  const tekst = wynik.wynik.text;
  widok.podglad.textContent = tekst;
  widok.wstawDoZrodla(tekst);
  okno.gotowe();

  const strony = wynik.wynik.pages;
  const droga = wynik.wynik.usedOcr
    ? 'tekst powstał z rozpoznania pisma, więc jest odczytem z pikseli, nie warstwą tekstową dokumentu'
    : 'tekst pochodzi z warstwy tekstowej dokumentu';
  odpowiedz.pokaz(
    `Wydobyto ${String(tekst.length)} znaków` +
      `${strony === undefined ? '' : ` z ${String(strony)} stron`} — ${droga}. ` +
      'Treść wpisano do pola Source Panel; zapis źródła zostaje czynnością Operatora. ' +
      'Stylów, tabel ani osadzeń odpowiedź nie niesie.',
    true,
  );
}

/** `document.convert` — zamiana formatu dokumentu po stronie rdzenia. */
async function zamienFormat(
  dokument: ZrodloDokumentuTranslate,
  sciezkaPliku: string,
  formatDocelowy: string,
  okno: StanOkna,
  odpowiedz: WierszOdpowiedzi,
): Promise<void> {
  const sciezka = sciezkaPliku.trim();
  if (sciezka === '') {
    const zdanie = 'Wskaż ścieżkę dokumentu — żądanie zamiany formatu niesie ścieżkę pliku.';
    odpowiedz.pokaz(zdanie, false);
    okno.blad(zdanie);
    return;
  }

  okno.ladowanie(`Rdzeń zamienia dokument na format ${formatDocelowy}.`);
  odpowiedz.pokaz('Zamiana formatu dokumentu…', true);
  const wynik = await dokument.zamienFormat(sciezka, formatDocelowy);
  if (!wynik.udany || wynik.wynik === undefined) {
    const zdanie = opisOdmowyBledu('Zamiana formatu dokumentu', wynik.blad);
    odpowiedz.pokaz(zdanie, false);
    okno.blad(zdanie);
    return;
  }
  okno.gotowe();
  odpowiedz.pokaz(
    `Rdzeń odłożył wynik w magazynie jako zasób ${wynik.wynik.asset.id} ` +
      `(${String(wynik.wynik.sizeBytes)} bajtów). Okno magazynu nie czyta — potwierdza ` +
      'odpowiedź rdzenia, nie istnienie pliku.',
    true,
  );
}

/** `translate.panel.export` — wydanie panelu w formacie pliku. */
async function wydajPanel(
  stan: StanTranslate,
  idPanelu: string,
  format: ExportFormat,
  okno: StanOkna,
  odpowiedz: WierszOdpowiedzi,
): Promise<void> {
  if (idPanelu === '') {
    const zdanie = 'Wskaż panel — żądanie eksportu niesie panel i format pliku.';
    odpowiedz.pokaz(zdanie, false);
    okno.blad(zdanie);
    return;
  }
  okno.ladowanie('Rdzeń wydaje panel do pliku.');
  odpowiedz.pokaz('Wydawanie panelu…', true);
  const sprawozdanie = await eksportuj(stan.panele, idPanelu, format);
  if (sprawozdanie.powodzenie) okno.gotowe();
  else okno.blad(sprawozdanie.tresc);
  odpowiedz.pokaz(sprawozdanie.tresc, sprawozdanie.powodzenie);
}

/** Warstwa trzecia: operacje dokumentu, których kontrakt nie niesie. */
function menuDokumentu(): HTMLElement {
  const rozwiniecie = utworzRozwiniecie({
    warstwa: 3,
    nazwa: 'Menu operacji dokumentu',
    wyjasnienie: 'Detekcja układu, porównanie układów i obieg zasobów lokalizacyjnych.',
    znacznik: '⋮',
  });
  rozwiniecie.tresc.append(
    przyciskBezKomendy('Ponowna detekcja układu', BRAKI.ukladIPorownanie),
    przyciskBezKomendy('Porównanie układów przed i po', BRAKI.ukladIPorownanie),
    przyciskBezKomendy('Ochrona kluczy zasobów lokalizacyjnych', BRAKI.zasobyLokalizacyjne),
    przyciskBezKomendy('Wczytanie napisów z taktowaniem', BRAKI.napisy),
  );
  return rozwiniecie.element;
}

/**
 * Warstwa czwarta: pseudolokalizacja.
 *
 * Jedyna funkcja warstwy czwartej tego okna, którą da się wykonać bez rdzenia,
 * bo jest przekształceniem znaków, a nie zapytaniem o cokolwiek. Wynik stoi
 * w oknie do przeniesienia ręcznego i mówi o sobie, że nigdzie się nie zapisał.
 */
function pseudolokalizacja(stan: StanTranslate, okno: StanOkna): HTMLElement {
  const rozwiniecie = utworzRozwiniecie({
    warstwa: 4,
    nazwa: 'Pseudolokalizacja',
    wyjasnienie:
      'Test wydłużenia tekstu i znaków diakrytycznych w interfejsie przed właściwym przekładem.',
    znacznik: '☰',
  });

  const wydluzenie = poleTekstowe({
    etykieta: 'Wydłużenie w procentach',
    podpowiedz: String(WYDLUZENIE_DOMYSLNE),
  });
  wydluzenie.kontrolka.value = String(WYDLUZENIE_DOMYSLNE);

  const zdanie = document.createElement('p');
  zdanie.className = 'mt-formaty__zdanie';

  const wynik = document.createElement('pre');
  wynik.className = 'mt-formaty__pseudo';

  const uruchom = przycisk('Przekształć tekst źródłowy', 'dn-btn dn-btn--sm dn-btn--zarys');
  uruchom.addEventListener('click', () => {
    const przeksztalcony = pseudolokalizuj(stan.tekstZrodlowy(), Number(wydluzenie.kontrolka.value));
    wynik.textContent = przeksztalcony.tekst;
    zdanie.textContent = zdanieOPseudolokalizacji(przeksztalcony);
    // Czynność jest miejscowa i nie pyta rdzenia, więc nie stawia okna
    // w ładowaniu; zdejmuje natomiast komunikat poprzedniej odmowy, bo od tej
    // chwili okno pokazuje wynik, a nie powód niewykonania.
    if (okno.faza() === 'blad') okno.gotowe();
  });

  rozwiniecie.tresc.append(wydluzenie.element, pasek(uruchom), zdanie, wynik);
  return rozwiniecie.element;
}

/** Warstwa czwarta: reguły liczby mnogiej i rodzaju — bez komendy w kontrakcie. */
function regulyGramatyczne(): HTMLElement {
  const rozwiniecie = utworzRozwiniecie({
    warstwa: 4,
    nazwa: 'Reguły liczby mnogiej i rodzaju',
    wyjasnienie: BRAKI.formyGramatyczne,
    znacznik: '☰',
  });
  rozwiniecie.tresc.append(
    przyciskBezKomendy('Zastosuj reguły języka docelowego', BRAKI.formyGramatyczne),
  );
  return rozwiniecie.element;
}

import {
  TelemetryFormat,
  UsageDimension,
  type UsageAggregate,
} from '../../../../shared/contract';
import {
  pobierzPlik,
  pole,
  poleLiczbowe,
  przelacznik,
  przyciskAkcji,
  wiersz,
  wybor,
} from '../../modele/kontrolki-formularza-braki';
import type { StanDiagnostyki } from './stan-diagnostyki';
import { utworzStanTresci } from './stany-okna';
import {
  cialoNarzedzia,
  objasnienieNarzedzia,
  pasekNarzedzia,
  type PozycjaZakladkiNarzedzi,
} from './zakladki-narzedzi';
import {
  nazwaWymiaru,
  POSTACI_RAPORTU,
  rozszerzenieRaportu,
  WYMIARY_ZUZYCIA,
  type ZrodloZuzycia,
} from './zuzycie-zrodlo';

/**
 * Usage & Cost — zużycie tokenów, żądań i kosztu wraz z raportem
 * rozliczeniowym.
 *
 * Zakładka nazywa się tak, jak nazywa ją opracowanie modułu (rozdz. 2.5,
 * „Zużycie kont, kosztów i budżetów (Usage & Cost)"), i stoi tam, gdzie
 * opracowanie ją stawia: w kontenerze Observability Tools, nie w Diagnostics
 * Center.
 *
 * Do dziś zakładka mówiła, że zużycia nie ma z czego policzyć, bo kontrakt nie
 * niesie rodziny komend rozliczenia. To zdanie przestało być prawdziwe:
 * kontrakt niesie `usage.summary.get` i `usage.report.build`, a rdzeń ma dla
 * nich uchwyty. Zdanie zostało więc zastąpione czynnością, a nie przepisane.
 *
 * Trzy rozstrzygnięcia zakładki:
 *
 *  1. Jeden wymiar naraz — tak stanowi kontrakt. Wykaz mówi, po którym wymiarze
 *     jest zebrany; tabela z dwoma wymiarami pod jednym nagłówkiem byłaby
 *     wykazem, którego rdzeń nigdy nie oddał.
 *  2. Puste zestawienie nie jest brakiem. Okres bez ani jednego wywołania jest
 *     poprawną odpowiedzią i zakładka nazywa to zdaniem, osobno od odmowy
 *     rdzenia i osobno od „jeszcze nie pytałem".
 *  3. Koszt niepełny mówi o sobie. Wywołanie kanału bez cennika nie wchodzi do
 *     kosztu i kontrakt niesie na to dwa pola (`priceCoverage`,
 *     `costWithoutPrice`). Suma podana bez tego zastrzeżenia wyglądałaby na
 *     rachunek, a jest częścią rachunku.
 *
 * Raport wytwarza plik. `usage.report.build` oddaje treść, nie zasób
 * w magazynie, więc wytworem jest plik pobrany na dysk Operatora — nazwany
 * okresem i postacią, żeby dwa raporty nie nadpisały się wzajemnie.
 */
export interface ZakladkaZuzycia {
  /** Kod obszaru i nazwa zakładki wraz z jej ciałem — gotowa pozycja kontenera. */
  pozycja: PozycjaZakladkiNarzedzi;
  /** Odczyt zestawienia z rdzenia. */
  odswiez(): void;
  /** Czy zakładka pytała już rdzeń — po tym kontener wie, czy odświeżać. */
  czytano(): boolean;
}

export function utworzZakladkeZuzycia(
  zrodlo: ZrodloZuzycia,
  stan: StanDiagnostyki,
): ZakladkaZuzycia {
  const tresc = utworzStanTresci();
  let czytano = false;

  const wymiar = wybor(
    'Wymiar zestawienia',
    WYMIARY_ZUZYCIA.map((pozycja) => [pozycja.kod, pozycja.nazwa] as const),
  );
  const granica = poleLiczbowe('Górna granica liczby pozycji', 'domyślnie wg rdzenia');
  const byty = pole(
    'Byty wymiaru zawężające wynik',
    'identyfikatory rozdzielone przecinkiem; puste znaczy komplet',
  );
  const prognoza = przelacznik('Dolicz prognozę do końca okresu rozliczeniowego');
  const postac = wybor(
    'Postać raportu',
    POSTACI_RAPORTU.map((pozycja) => [pozycja.kod, pozycja.nazwa] as const),
  );

  const odczyt = przyciskAkcji('Odczytaj zestawienie zużycia', 'dn-btn dn-btn--atrament');
  odczyt.addEventListener('click', () => void odczytaj());

  const raport = przyciskAkcji('Zbuduj raport rozliczeniowy', 'dn-btn dn-btn--zarys');
  raport.addEventListener('click', () => void zbudujRaport());

  /** Okres bierze się z zakresu wspólnego modułowi; końce nieustawione nie jadą. */
  function okres(): { fromTime?: number; toTime?: number } {
    const zakres = stan.zakres();
    return {
      ...(zakres.od === undefined ? {} : { fromTime: zakres.od }),
      ...(zakres.do === undefined ? {} : { toTime: zakres.do }),
    };
  }

  /** Wymiar wskazany; wartość spoza wykazu bierze kanał modelu. */
  function wskazanyWymiar(): UsageDimension {
    const pozycja = WYMIARY_ZUZYCIA.find((wpis) => wpis.kod === wymiar.value);
    return pozycja === undefined ? UsageDimension.Channel : pozycja.kod;
  }

  function wskazanaPostac(): TelemetryFormat {
    const pozycja = POSTACI_RAPORTU.find((wpis) => wpis.kod === postac.value);
    return pozycja === undefined ? TelemetryFormat.Csv : pozycja.kod;
  }

  /** Wykaz bytów wymiaru z pola; puste pozycje odpadają, bo nie są bytem. */
  function wskazaneByty(): string[] {
    return byty.value
      .split(',')
      .map((pozycja) => pozycja.trim())
      .filter((pozycja) => pozycja !== '');
  }

  async function odczytaj(): Promise<void> {
    czytano = true;
    const kod = wskazanyWymiar();
    tresc.ladowanie(`Rdzeń zbiera zużycie po wymiarze „${nazwaWymiaru(kod)}"…`);

    const liczba = Number(granica.value);
    const wybraneByty = wskazaneByty();
    const wynik = await zrodlo.zestawienie({
      dimension: kod,
      ...okres(),
      ...(wybraneByty.length === 0 ? {} : { dimensionIds: wybraneByty }),
      ...(prognoza.checked ? { includeForecast: true } : {}),
      ...(granica.value.trim() === '' || !Number.isFinite(liczba) ? {} : { limit: liczba }),
    });

    if (!wynik.udany || wynik.wynik === undefined) {
      tresc.blad(
        `Zestawienia zużycia po wymiarze „${nazwaWymiaru(kod)}" rdzeń nie oddał.`,
        wynik.blad,
      );
      return;
    }

    const odpowiedz = wynik.wynik;
    if (odpowiedz.aggregates.length === 0) {
      // Pusty wykaz jest odpowiedzią, nie brakiem funkcji, i zdanie musi to
      // rozstrzygnąć: „nie było ani jednego wywołania" znaczy co innego niż
      // „zużycia nie ma z czego policzyć".
      tresc.pusto(
        `Zużycia w tym okresie nie było: rdzeń nie zapisał ani jednego wywołania kanału ` +
          `modelu po wymiarze „${nazwaWymiaru(kod)}" między ${czas(odpowiedz.fromTime)} ` +
          `a ${czas(odpowiedz.toTime)}. To pusty okres, nie brak pomiaru — pomiar odpowiedział.`,
      );
      return;
    }

    const miejsce = tresc.tresc();
    miejsce.replaceChildren(
      naglowekZestawienia(odpowiedz.fromTime, odpowiedz.toTime, kod, odpowiedz.priceCoverage),
      ...(odpowiedz.overall === undefined
        ? []
        : [zdanieLaczne(odpowiedz.overall, odpowiedz.currency)]),
      wykazPozycji(odpowiedz.aggregates, odpowiedz.currency),
    );
  }

  async function zbudujRaport(): Promise<void> {
    const wybrana = wskazanaPostac();
    tresc.potwierdzenie('Rdzeń składa raport rozliczeniowy — trwa…', true);
    raport.disabled = true;
    const wynik = await zrodlo.raport({ ...okres(), format: wybrana });
    raport.disabled = false;

    if (!wynik.udany || wynik.wynik === undefined) {
      tresc.blad('Raportu rozliczeniowego rdzeń nie złożył.', wynik.blad);
      return;
    }

    const odpowiedz = wynik.wynik;
    // Raport pusty też jest raportem i też schodzi na dysk: plik z nagłówkiem
    // bez wierszy jest dowodem, że w okresie nic nie było. Zdanie mówi to
    // wprost, żeby nikt nie wziął pustego pliku za nieudane wywołanie.
    const nazwa = `zuzycie-raport-${odpowiedz.generatedAt}.${rozszerzenieRaportu(odpowiedz.format)}`;
    pobierzPlik(nazwa, odpowiedz.content, rodzajTresci(odpowiedz.format));
    tresc.potwierdzenie(
      `Raport rozliczeniowy złożony i pobrany jako ${nazwa}: wywołań objętych ` +
        `${odpowiedz.requests}, postać ${odpowiedz.format}, chwila złożenia ` +
        `${czas(odpowiedz.generatedAt)}.` +
        (odpowiedz.requests === 0
          ? ' Wywołań w okresie nie było — plik niesie sam nagłówek i to jest jego treść.'
          : ''),
      true,
    );
  }

  const element = cialoNarzedzia(
    objasnienieNarzedzia(
      'Zużycie tokenów, żądań i kosztu zebrane w JEDNYM wymiarze — tak stanowi kontrakt: ' +
        'zestawienie po dwóch wymiarach naraz jest dwoma żądaniami, nie jednym. Okres bierze ' +
        'się z zakresu czasu wspólnego modułowi, ustawianego w Diagnostics Center.',
    ),
    objasnienieNarzedzia(
      'Koszt pusty przy pozycji znaczy brak cennika kanału, a nie koszt zerowy — i zakładka ' +
        'mówi osobno, jaki udział wywołań cennik obejmuje. Suma podana bez tego zastrzeżenia ' +
        'wyglądałaby na rachunek, a jest częścią rachunku.',
    ),
    wiersz('Wymiar zestawienia', wymiar, {
      klasa: 'dg-narzedzie__pole',
      objasnienie: 'Siedem wymiarów kontraktu; jeden na odczyt.',
    }),
    wiersz('Byty wymiaru', byty, {
      klasa: 'dg-narzedzie__pole',
      objasnienie: 'Puste znaczy komplet bytów wymiaru, nie brak zawężenia do pokazania.',
    }),
    wiersz('Górna granica liczby pozycji', granica, {
      klasa: 'dg-narzedzie__pole',
      objasnienie: 'Puste zostawia granicę rdzeniowi.',
    }),
    wiersz('Prognoza okresu rozliczeniowego', prognoza, {
      klasa: 'dg-narzedzie__pole',
      objasnienie: 'Doliczenie prognozy do końca okresu rozliczeniowego wykonuje rdzeń.',
    }),
    wiersz('Postać raportu', postac, {
      klasa: 'dg-narzedzie__pole',
      objasnienie: 'Cztery postaci kontraktu. Raport schodzi na dysk Operatora jako plik.',
    }),
    pasekNarzedzia(odczyt, raport),
    tresc.element,
  );

  return {
    pozycja: { kod: 'usage-cost', nazwa: 'Usage & Cost', element },
    odswiez: () => void odczytaj(),
    czytano: () => czytano,
  };
}

/** Nagłówek zestawienia: okres, wymiar i pokrycie cennikiem. */
function naglowekZestawienia(
  od: number,
  doo: number,
  kod: UsageDimension,
  pokrycie: number | undefined,
): HTMLElement {
  const element = document.createElement('p');
  element.className = 'dn-pole-opis';
  // Pokrycie nieoddane znaczy „rdzeń o nim nie powiedział", a nie „pełne".
  // Dopisanie stu procent z własnej głowy byłoby zapewnieniem, którego nikt
  // nie wydał.
  const zdanie =
    pokrycie === undefined
      ? 'Udziału wywołań objętych cennikiem rdzeń nie podał — o pełności kosztu ta odpowiedź nie orzeka.'
      : pokrycie >= 100
        ? 'Cennik obejmuje wszystkie wywołania okresu, więc koszt jest pełny.'
        : `Cennik obejmuje ${pokrycie}% wywołań okresu — koszt poniżej jest NIEPEŁNY.`;
  element.textContent =
    `Zużycie po wymiarze „${nazwaWymiaru(kod)}" za okres ${czas(od)} – ${czas(doo)}. ${zdanie}`;
  return element;
}

/** Zdanie zestawienia łącznego dla całego zakresu. */
function zdanieLaczne(laczne: UsageAggregate, waluta: string | undefined): HTMLElement {
  const element = document.createElement('p');
  element.className = 'dn-pole-opis';
  element.textContent = `Łącznie w zakresie: ${opisPozycji(laczne, waluta)}`;
  return element;
}

/** Wykaz pozycji zestawienia — jedna pozycja na byt wymiaru. */
function wykazPozycji(
  pozycje: readonly UsageAggregate[],
  waluta: string | undefined,
): HTMLElement {
  const element = document.createElement('ul');
  element.className = 'dg-wykaz';
  element.setAttribute('aria-label', 'Pozycje zestawienia zużycia');
  for (const pozycja of pozycje) {
    const wpis = document.createElement('li');
    wpis.className = 'dg-pozycja';
    const tytul = document.createElement('strong');
    tytul.className = 'dg-pozycja__tytul';
    // Nazwa bytu bywa nieoddana; identyfikator niesie tożsamość zawsze, więc
    // wtedy mówi się identyfikatorem, a nie zwrotem „bez nazwy".
    tytul.textContent =
      (pozycja.dimensionLabel ?? '').trim() === ''
        ? pozycja.dimensionId
        : `${pozycja.dimensionLabel} (${pozycja.dimensionId})`;
    const opis = document.createElement('span');
    opis.className = 'dg-pozycja__opis';
    opis.textContent = opisPozycji(pozycja, waluta);
    wpis.append(tytul, opis);
    element.append(wpis);
  }
  return element;
}

/**
 * Opis jednej pozycji zestawienia.
 *
 * Koszt pusty nie zamienia się w zero, a wywołania pominięte w koszcie stają
 * w zdaniu obok niego: liczba z gwiazdką bez wyjaśnienia byłaby liczbą, której
 * nikt nie umie użyć.
 */
function opisPozycji(pozycja: UsageAggregate, waluta: string | undefined): string {
  const jednostka = (pozycja.currency ?? waluta ?? '').trim();
  const koszt =
    pozycja.cost === undefined
      ? 'kosztu nie ma, bo kanał nie ma cennika'
      : `koszt ${pozycja.cost}${jednostka === '' ? '' : ` ${jednostka}`}`;
  const nieujete =
    pozycja.costWithoutPrice === undefined || pozycja.costWithoutPrice === 0
      ? ''
      : ` Poza kosztem zostało wywołań: ${pozycja.costWithoutPrice} — brak cennika kanału.`;
  const nieudane =
    pozycja.failedRequests === undefined || pozycja.failedRequests === 0
      ? ''
      : ` Zakończonych błędem: ${pozycja.failedRequests}.`;
  const zPamieci =
    pozycja.cachedTokens === undefined
      ? ''
      : ` Tokenów z pamięci podręcznej: ${pozycja.cachedTokens}.`;
  return (
    `żądań ${pozycja.requests}, tokenów łącznie ${pozycja.totalTokens} ` +
    `(prompt ${pozycja.promptTokens}, odpowiedź ${pozycja.completionTokens}), ${koszt}.` +
    zPamieci +
    nieudane +
    nieujete
  );
}

/** Rodzaj treści pliku raportu — postać rozstrzyga, nie nazwa komendy. */
function rodzajTresci(postac: TelemetryFormat): string {
  return postac === TelemetryFormat.Csv ? 'text/csv' : 'application/json';
}

/**
 * Chwila w zapisie miejscowym Operatora.
 *
 * Rdzeń podaje milisekundy epoki i nie zna strefy Operatora, więc formatowanie
 * należy do widoku. Zero znaczy „bez granicy", a nie „rok tysiąc dziewięćset
 * siedemdziesiąty".
 */
function czas(znacznik: number): string {
  if (znacznik === 0) return 'bez granicy';
  return new Date(znacznik).toLocaleString('pl-PL');
}

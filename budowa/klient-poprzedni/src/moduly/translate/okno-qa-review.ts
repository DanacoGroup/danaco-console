import { TranslationStatus, type ExportFormat, type TranslationPanel } from '../../../../shared/contract';
import {
  poleWyboru,
  przycisk,
  przyciskBezKomendy,
  utworzWierszOdpowiedzi,
  type WierszOdpowiedzi,
} from '../../modele/kontrolki-formularza';
import { eksportuj, kontrolaJakosci, type Sprawozdanie } from './czynnosci-panelu';
import { FORMATY_EKSPORTU, PUSTE, nazwaStanuPanelu } from './etykiety-translate';
import { naglowekOkna } from './kontrolki-translate';
import { liczSlowa, utworzKarteLqa, type KartaLqa } from './karta-lqa';
import { utworzStanOkna, type StanOkna } from './stan-okna-translate';
import type { StanTranslate } from './stan-translate';
import { oznaczWarstwe, utworzRozwiniecie } from './warstwy-translate';

/**
 * QA & Review Center — okno **zarządca** modułu Translate.
 *
 * Okno robi to, czego pojedynczy panel zrobić nie może: prowadzi kontrolę
 * jakości wszystkich paneli naraz i zbiera ich zastrzeżenia w jednym wykazie.
 * Kontrakt ma kontrolę jakości jednego panelu, więc kontrola zbiorcza jest
 * powtórzeniem tej komendy dla każdego panelu — i okno mówi to wprost, zamiast
 * sugerować zdolność wsadową, której rdzeń nie ma.
 *
 * Wywołania idą równolegle. Kontrola jednego panelu nie zależy od kontroli
 * drugiego, więc szeregowanie ich tylko wydłużałoby czekanie; odmowa jednego
 * panelu zostaje przy nim i nie przerywa pozostałym (`Promise.allSettled` nie
 * jest tu potrzebny, bo źródło modułu oddaje odmowę wynikiem, nie wyjątkiem).
 *
 * Przebieg akceptacji tłumaczenie → korekta → zatwierdzenie pokazuje stany,
 * które kontrakt zna: oczekuje, tłumaczenie w toku, gotowe, błąd. Etapu
 * zatwierdzenia ani autora zmiany w kontrakcie nie ma, więc wskaźnik mówi
 * o stanie wykonania i nazywa tę różnicę, zamiast malować przebieg, którego
 * rdzeń nie prowadzi.
 */
export interface OknoQaReview {
  element: HTMLElement;
  odswiez(): void;
}

const BRAKI = {
  profil:
    'Żądanie kontroli jakości niesie sam panel — pola na profil reguł w nim nie ma, więc zakresu ' +
    'kontroli nie da się zmienić z okna.',
  korekta:
    'Korekty językowej — gramatyki, ortografii, stylu, czytelności, typografii rynku — kontrakt ' +
    'nie prowadzi. Kontrola jakości zna sześć rodzajów niezgodności i korekty wśród nich nie ma.',
  akceptacja:
    'Panel ma w kontrakcie cztery stany wykonania. Etapu akceptacji, śladu autora zmiany ani ' +
    'zatwierdzenia nie ma w nim czym zapisać.',
  pakiet:
    'Pakiet dla wykonawcy zewnętrznego wymaga materiału w formacie XLIFF, pamięci tłumaczeń ' +
    'i bazy terminów w jednym pliku. Kontrakt wydaje bazę terminów i panele osobno, pamięci nie ' +
    'wydaje wcale, a spakowania całości nie przewiduje.',
  wycena:
    'Wyceny nie ma z czego złożyć: kontrakt nie niesie ani siatki stawek, ani analizy materiału ' +
    'względem pamięci tłumaczeń.',
  poprawkiMasowe:
    'Poprawki masowej kontrakt nie ma — treść panelu zmienia się zapisem korekty jednego panelu.',
} as const;

export function utworzOknoQaReview(stan: StanTranslate): OknoQaReview {
  const okno: StanOkna = utworzStanOkna(PUSTE.jakosc);
  const odpowiedz = utworzWierszOdpowiedzi();

  const uruchom = przycisk('Uruchom kontrolę wszystkich paneli', 'dn-btn dn-btn--sm dn-btn--atrament');
  uruchom.addEventListener('click', () => {
    void sprawdzWszystkie(stan, okno, odpowiedz);
  });

  const format = poleWyboru(
    { etykieta: 'Format eksportu zbiorczego' },
    FORMATY_EKSPORTU.map((pozycja) => ({ wartosc: pozycja.wartosc, etykieta: pozycja.etykieta })),
  );

  const eksportZbiorczy = przycisk('Eksport zbiorczy paneli', 'dn-btn dn-btn--sm dn-btn--zarys');
  eksportZbiorczy.addEventListener('click', () => {
    void wydajWszystkie(stan, format.kontrolka.value as ExportFormat, okno, odpowiedz);
  });

  const pasek = document.createElement('div');
  pasek.className = 'mt-pasek';
  pasek.append(uruchom, eksportZbiorczy);

  const zastrzezenia = document.createElement('div');
  zastrzezenia.className = 'mt-jakosc__wykaz';

  const przebieg = document.createElement('p');
  przebieg.className = 'mt-jakosc__przebieg';

  const kontrola = document.createElement('div');
  kontrola.className = 'mt-jakosc__kontrola';
  oznaczWarstwe(kontrola, 1);
  kontrola.append(format.element, pasek, odpowiedz.element, przebieg, zastrzezenia);

  const karta: KartaLqa = utworzKarteLqa(() => liczSlowa(stan.tekstZrodlowy()));

  okno.tresc.append(kontrola, profilKontroli(), menuWydania(), korektaJezykowa(), karta.element);

  const element = document.createElement('section');
  element.className = 'mt-okno mt-okno--zarzadca';
  element.dataset['okno'] = 'qa-review-center';
  element.append(naglowekOkna('QA & Review Center', 'zarządca'), okno.element);

  /**
   * Przerysowanie odbudowuje wykaz zastrzeżeń z paneli stanu modułu.
   *
   * Zastrzeżenia nie są tu wynikiem ostatniego kliknięcia, tylko odbiciem tego,
   * co rdzeń trzyma przy panelach: kontrola zapisuje je w panelu, a panel
   * wraca do modułu zdarzeniem zmiany. Dzięki temu wykaz jest prawdziwy także
   * wtedy, gdy kontrolę uruchomiono z paska narzędzi pojedynczego panelu.
   */
  function odswiez(): void {
    const panele = stan.panelJezykow();
    zastrzezenia.replaceChildren(...panele.map(grupaPanelu));
    przebieg.textContent = zdaniePrzebiegu(panele);
    if (okno.faza() === 'ladowanie' || okno.faza() === 'blad') return;
    if (panele.length === 0) {
      okno.puste(PUSTE.jakosc);
      return;
    }
    okno.gotowe();
  }

  odswiez();
  return { element, odswiez };
}

/**
 * `translate.quality.check` powtórzona dla każdego panelu okna.
 *
 * Bilans jest obowiązkowy: przy wielu wywołaniach część potrafi się nie udać,
 * a zdanie mówiące wyłącznie o powodzeniu ukryłoby panele, których nie
 * sprawdzono. Powody odmów idą po nazwie panelu, żeby wiadomo było który.
 */
async function sprawdzWszystkie(
  stan: StanTranslate,
  okno: StanOkna,
  odpowiedz: WierszOdpowiedzi,
): Promise<void> {
  const panele = stan.panelJezykow();
  if (panele.length === 0) {
    const zdanie =
      'Nie ma czego sprawdzać: okno nie zna ani jednego panelu języka. Dodaj język docelowy ' +
      'w Translation Panels.';
    odpowiedz.pokaz(zdanie, false);
    okno.blad(zdanie);
    return;
  }

  okno.ladowanie(
    `Rdzeń sprawdza ${String(panele.length)} paneli — po jednym wywołaniu kontroli na panel.`,
  );
  odpowiedz.pokaz('Kontrola jakości wszystkich paneli…', true);

  const sprawozdania = await Promise.all(
    panele.map(async (panel) => ({
      jezyk: panel.language,
      sprawozdanie: await kontrolaJakosci(stan.panele, panel.id),
    })),
  );

  okno.gotowe();
  odpowiedz.pokaz(bilans('Kontrola jakości', sprawozdania), czyWszystkieUdane(sprawozdania));
}

/** `translate.panel.export` powtórzony dla każdego panelu okna. */
async function wydajWszystkie(
  stan: StanTranslate,
  format: ExportFormat,
  okno: StanOkna,
  odpowiedz: WierszOdpowiedzi,
): Promise<void> {
  const panele = stan.panelJezykow();
  if (panele.length === 0) {
    const zdanie = 'Nie ma czego wydać: okno nie zna ani jednego panelu języka.';
    odpowiedz.pokaz(zdanie, false);
    okno.blad(zdanie);
    return;
  }

  okno.ladowanie(`Rdzeń wydaje ${String(panele.length)} paneli — po jednym wywołaniu na panel.`);
  odpowiedz.pokaz('Eksport zbiorczy paneli…', true);

  const sprawozdania = await Promise.all(
    panele.map(async (panel) => ({
      jezyk: panel.language,
      sprawozdanie: await eksportuj(stan.panele, panel.id, format),
    })),
  );

  okno.gotowe();
  odpowiedz.pokaz(bilans('Eksport zbiorczy', sprawozdania), czyWszystkieUdane(sprawozdania));
}

/** Jedno sprawozdanie wraz z językiem panelu, którego dotyczy. */
interface SprawozdaniePanelu {
  jezyk: string;
  sprawozdanie: Sprawozdanie;
}

function czyWszystkieUdane(sprawozdania: readonly SprawozdaniePanelu[]): boolean {
  return sprawozdania.every((wpis) => wpis.sprawozdanie.powodzenie);
}

function bilans(czynnosc: string, sprawozdania: readonly SprawozdaniePanelu[]): string {
  const nieudane = sprawozdania.filter((wpis) => !wpis.sprawozdanie.powodzenie);
  const naglowek =
    `${czynnosc}: paneli ${String(sprawozdania.length)}, udanych ` +
    `${String(sprawozdania.length - nieudane.length)}, nieudanych ${String(nieudane.length)}.`;
  if (nieudane.length === 0) return naglowek;
  return `${naglowek} ${nieudane
    .map((wpis) => `panel ${wpis.jezyk}: ${wpis.sprawozdanie.tresc}`)
    .join(' · ')}`;
}

/** Grupa zastrzeżeń jednego panelu — nagłówek z językiem i stanem, pod nim wykaz. */
function grupaPanelu(panel: TranslationPanel): HTMLElement {
  const tytul = document.createElement('h5');
  tytul.className = 'mt-jakosc__tytul';
  const zastrzezenia = panel.issues ?? [];
  tytul.textContent =
    `Panel ${panel.language} — ${nazwaStanuPanelu(panel.status)}, ` +
    `zastrzeżeń ${String(zastrzezenia.length)}`;

  const element = document.createElement('section');
  element.className = 'mt-jakosc__grupa';
  element.dataset['panel'] = panel.id;
  element.append(tytul);

  if (zastrzezenia.length === 0) {
    const puste = document.createElement('p');
    puste.className = 'mt-jakosc__puste';
    puste.textContent =
      'Rdzeń nie trzyma przy tym panelu żadnego zastrzeżenia. Wykaz pusty znaczy tyle, ' +
      'ile powiedziała ostatnia kontrola — nie znaczy, że kontrolę uruchomiono.';
    element.append(puste);
    return element;
  }

  const wykaz = document.createElement('ul');
  wykaz.className = 'mt-zastrzezenia';
  wykaz.replaceChildren(
    ...zastrzezenia.map((zastrzezenie) => {
      const wiersz = document.createElement('li');
      wiersz.className = 'mt-zastrzezenia__wiersz';
      wiersz.dataset['rodzaj'] = zastrzezenie.kind;
      wiersz.textContent = [zastrzezenie.kind, zastrzezenie.segment, zastrzezenie.detail]
        .filter((czesc): czesc is string => czesc !== undefined && czesc !== '')
        .join(' — ');
      return wiersz;
    }),
  );
  element.append(wykaz);
  return element;
}

/** Wskaźnik przebiegu — liczony ze stanów paneli, jedynego, co kontrakt niesie. */
function zdaniePrzebiegu(panele: readonly TranslationPanel[]): string {
  if (panele.length === 0) return 'Przebieg: okno nie zna ani jednego panelu języka.';
  const gotowe = panele.filter((panel) => panel.status === TranslationStatus.Ready).length;
  const bledne = panele.filter((panel) => panel.status === TranslationStatus.Error).length;
  return (
    `Przebieg wykonania: gotowych ${String(gotowe)} z ${String(panele.length)}` +
    `${bledne === 0 ? '' : `, w błędzie ${String(bledne)}`}. ` +
    'To stan wykonania panelu, nie etap akceptacji — etapu korekty ani zatwierdzenia ' +
    'kontrakt nie prowadzi.'
  );
}

/** Warstwa druga: profil kontroli jakości. */
function profilKontroli(): HTMLElement {
  const rozwiniecie = utworzRozwiniecie({
    warstwa: 2,
    nazwa: 'Profil kontroli jakości',
    wyjasnienie: 'Zestaw reguł, wedle którego biegnie kolejny przebieg kontroli.',
    znacznik: '▼',
  });
  rozwiniecie.tresc.append(przyciskBezKomendy('Wybór profilu kontroli', BRAKI.profil));
  return rozwiniecie.element;
}

/** Warstwa trzecia: menu wydania i przebiegu akceptacji. */
function menuWydania(): HTMLElement {
  const rozwiniecie = utworzRozwiniecie({
    warstwa: 3,
    nazwa: 'Menu wydania',
    wyjasnienie:
      'Pakiet dla wykonawcy zewnętrznego, statystyka objętości z wyceną i przebieg akceptacji.',
    znacznik: '☰',
  });
  rozwiniecie.tresc.append(
    przyciskBezKomendy('Pakiet dla wykonawcy zewnętrznego', BRAKI.pakiet),
    przyciskBezKomendy('Odbiór zwrotu od wykonawcy', BRAKI.pakiet),
    przyciskBezKomendy('Statystyka objętości i wycena', BRAKI.wycena),
    przyciskBezKomendy('Zatwierdzenie materiału', BRAKI.akceptacja),
  );
  return rozwiniecie.element;
}

/** Warstwa trzecia: korekta językowa panelu. */
function korektaJezykowa(): HTMLElement {
  const rozwiniecie = utworzRozwiniecie({
    warstwa: 3,
    nazwa: 'Korekta językowa panelu',
    wyjasnienie: BRAKI.korekta,
    znacznik: '⋮',
  });
  rozwiniecie.tresc.append(
    przyciskBezKomendy('Gramatyka, ortografia i interpunkcja', BRAKI.korekta),
    przyciskBezKomendy('Styl, rejestr i czytelność', BRAKI.korekta),
    przyciskBezKomendy('Typografia rynku docelowego', BRAKI.korekta),
    przyciskBezKomendy('Poprawki masowe wskazanych segmentów', BRAKI.poprawkiMasowe),
  );
  return rozwiniecie.element;
}

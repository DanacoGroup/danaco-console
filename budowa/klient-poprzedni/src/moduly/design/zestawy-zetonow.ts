import {
  Command,
  DesignTokenKind,
  DesignTokenTarget,
  type DesignToken,
  type DesignTokenSet,
} from '../../../../shared/contract';
import { opisOdmowyBledu } from '../../komponenty/odmowa';
import {
  poleTekstowe,
  poleWielowierszowe,
  poleWyboru,
  przycisk,
  ustawPozycje,
  utworzWierszOdpowiedzi,
} from '../../modele/kontrolki-formularza';
import type { StanDesignu } from './stan-designu';

/**
 * Zestawy żetonów w Tokens & System Panel — `design.tokenset.save`,
 * `design.tokenset.list`, `design.tokenset.export`, `design.tokenset.import`
 * oraz `design.styleguide.publish`.
 *
 * ── Rdzeń nie nadpisuje motywu produktu ─────────────────────────────────────
 * Motyw jest własnością powłoki. Zestaw żetonów jest bytem OBOK niego: panel
 * odczytuje żetony motywu obowiązującego (`zetony-systemu.ts`), a Operator może
 * je odłożyć jako zestaw, wczytać cudzy system i wydać go do kodu. Zapis
 * zestawu nie zmienia wyglądu produktu ani o jeden piksel i panel mówi to
 * wprost, zamiast zostawiać Operatora z domysłem.
 *
 * ── Import wnosi role nieznane i JE WYMIENIA ────────────────────────────────
 * System projektowy klienta ma własne nazwy ról. Odrzucenie ich byłoby
 * zgubieniem pracy, przemilczenie — obietnicą, że wszystko pasuje. Rola wchodzi
 * do zestawu, a jej nazwa wraca w `unknownNames` i panel wypisuje ją w całości.
 *
 * ── Postacie wydania składa rdzeń ───────────────────────────────────────────
 * CSS, SCSS, Tailwind, moduł JavaScript, Swift i Kotlin powstają w rdzeniu przez
 * sklejenie napisów — bez jednego programu spoza instalki. Panel pokazuje treść
 * wydania, a nie samo „gotowe": plik bez treści jest kopertą udaną i pustą.
 */
export interface ZestawyZetonow {
  element: HTMLElement;
  /** Zleca odczyt zestawów żetonów okna. */
  wczytaj(): Promise<void>;
}

/** Czym panel steruje w oknie żetonów. */
export interface SterowanieZestawami {
  /** Żetony odczytane z motywu obowiązującego — materiał na zestaw. */
  zetonyMotywu(): readonly DesignToken[];
}

const POSTACIE_WYDANIA = [
  { wartosc: DesignTokenTarget.CssVariables, etykieta: 'Zmienne CSS' },
  { wartosc: DesignTokenTarget.Scss, etykieta: 'Zmienne SCSS' },
  { wartosc: DesignTokenTarget.Tailwind, etykieta: 'Konfiguracja Tailwind' },
  { wartosc: DesignTokenTarget.Javascript, etykieta: 'Moduł JavaScript' },
  { wartosc: DesignTokenTarget.Ios, etykieta: 'Zasoby iOS (Swift)' },
  { wartosc: DesignTokenTarget.Android, etykieta: 'Zasoby Android (Kotlin)' },
] as const;

export function utworzZestawyZetonow(
  stan: StanDesignu,
  sterowanie: SterowanieZestawami,
): ZestawyZetonow {
  let zbior: readonly DesignTokenSet[] = [];

  const nazwa = poleTekstowe({
    etykieta: 'Nazwa zestawu',
    podpowiedz: 'np. System marki — motyw ciemny',
    opis: 'Pole name żądania design.tokenset.save.',
  });
  const motyw = poleTekstowe({
    etykieta: 'Motyw zestawu',
    podpowiedz: 'np. ciemny',
    opis:
      'Pole theme. Oba motywy są równoprawne i niosą własne wartości tych samych ról, ' +
      'więc motyw jasny i ciemny to dwa zestawy, a nie jeden.',
  });
  const wybor = poleWyboru(
    {
      etykieta: 'Zestaw',
      opis:
        'Wskazany zestaw jest celem nadpisania (pole tokenSetId), wydania do kodu ' +
        'i przewodnika stylu. Pozycja „nowy zestaw" zakłada kolejny obok.',
    },
    [{ wartosc: '', etykieta: 'nowy zestaw' }],
  );
  const postac = poleWyboru(
    { etykieta: 'Postać wydania', opis: 'Pole target żądania design.tokenset.export.' },
    [...POSTACIE_WYDANIA],
  );
  const modul = poleTekstowe({
    etykieta: 'Moduł docelowy',
    podpowiedz: 'np. library',
    opis:
      'Pole targetModuleId. Przy wydaniu żetonów puste oddaje treść samemu oknu; ' +
      'przy przewodniku stylu jest wymagane — przewodnik wydany donikąd jest plikiem, ' +
      'o którym nikt się nie dowie.',
  });
  const zapis = poleWielowierszowe(
    {
      etykieta: 'Zapis do wczytania',
      podpowiedz: '{"sygnal":"#c8a24a"} albo :root { --dn-sygnal: #c8a24a; }',
      opis:
        'Pole contentBase64 żądania design.tokenset.import — okno koduje treść samo. ' +
        'Rdzeń czyta zapis JSON (płaski albo zagnieżdżony) oraz zmienne CSS.',
    },
    4,
  );

  const zapisz = przycisk('Odłóż żetony motywu jako zestaw', 'dn-btn dn-btn--sm dn-btn--atrament');
  const odczytaj = przycisk('Odczytaj zestawy', 'dn-btn dn-btn--sm dn-btn--zarys');
  const wydaj = przycisk('Wydaj zestaw do kodu', 'dn-btn dn-btn--sm dn-btn--zarys');
  const wczytajZapis = przycisk('Wczytaj zapis zewnętrzny', 'dn-btn dn-btn--sm dn-btn--zarys');
  const przewodnik = przycisk('Wydaj przewodnik stylu', 'dn-btn dn-btn--sm dn-btn--zarys');
  const odpowiedz = utworzWierszOdpowiedzi();

  const pasek = document.createElement('div');
  pasek.className = 'md-czynnosci__pasek';
  pasek.append(zapisz, odczytaj, wydaj, wczytajZapis, przewodnik);

  const wydanie = document.createElement('pre');
  wydanie.className = 'md-zetony__wydanie';

  const element = document.createElement('div');
  element.className = 'md-zestawy';
  element.append(
    nazwa.element,
    motyw.element,
    wybor.element,
    postac.element,
    modul.element,
    zapis.element,
    pasek,
    wydanie,
    odpowiedz.element,
  );

  zapisz.addEventListener('click', () => void zapiszZestaw());
  odczytaj.addEventListener('click', () => void wczytaj());
  wydaj.addEventListener('click', () => void wydajZestaw());
  wczytajZapis.addEventListener('click', () => void wczytajZewnetrzny());
  przewodnik.addEventListener('click', () => void wydajPrzewodnik());
  wybor.kontrolka.addEventListener('change', () => {
    const zestaw = zbior.find((pozycja) => pozycja.id === wybor.kontrolka.value);
    nazwa.kontrolka.value = zestaw?.name ?? '';
    motyw.kontrolka.value = zestaw?.theme ?? '';
  });

  function bezOkna(komenda: string): boolean {
    if (stan.idOkna() !== '') return false;
    odpowiedz.pokaz(`Komenda ${komenda} wymaga okna modułu. ${stan.opisOkna()}`, false);
    return true;
  }

  function pokazZestawy(): void {
    ustawPozycje(wybor.kontrolka, [
      { wartosc: '', etykieta: 'nowy zestaw' },
      ...zbior.map((zestaw) => ({
        wartosc: zestaw.id,
        etykieta: `${zestaw.name} — żetonów: ${zestaw.tokenCount}`,
      })),
    ]);
  }

  async function zapiszZestaw(): Promise<void> {
    if (bezOkna(Command.DesignTokensetSave)) return;
    // Żetony biorą się z motywu obowiązującego, czytane w chwili naciśnięcia,
    // a nie z kopii zrobionej przy otwarciu panelu: między jednym a drugim
    // Operator mógł przełączyć motyw.
    const zetony = sterowanie.zetonyMotywu();
    if (zetony.length === 0) {
      odpowiedz.pokaz(
        'Motyw obowiązujący nie oddał ani jednego żetonu — zestaw pusty nie jest systemem ' +
          'projektowym i rdzeń go odmówi.',
        false,
      );
      return;
    }
    odpowiedz.pokaz('Zapis zestawu żetonów…', true);
    const wynik = await stan.czuwanie.prowadz(
      'zapis zestawu żetonów',
      stan.zrodlo.zapiszZestawZetonow({
        idOkna: stan.idOkna(),
        nazwa: nazwa.kontrolka.value === '' ? 'Zestaw z motywu obowiązującego' : nazwa.kontrolka.value,
        zetony,
        idZestawu: wybor.kontrolka.value,
        motyw: motyw.kontrolka.value,
      }),
      {
        wToku: (zdanie) => odpowiedz.pokaz(zdanie, true),
        cisza: (zdanie) => odpowiedz.pokaz(zdanie, false),
        powrot: (zdanie) => odpowiedz.pokaz(zdanie, false),
        spozniona: (zdanie) => odpowiedz.pokaz(zdanie, false),
      },
    );
    if (wynik === null) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowyBledu('Zapis zestawu żetonów', wynik.blad), false);
      return;
    }
    const zestaw = wynik.wynik.tokenSet;
    zbior = [zestaw, ...zbior.filter((pozycja) => pozycja.id !== zestaw.id)];
    pokazZestawy();
    wybor.kontrolka.value = zestaw.id;
    odpowiedz.pokaz(
      `Rdzeń utrwalił zestaw „${zestaw.name}" (${zestaw.id}) o ${zestaw.tokenCount} żetonach. ` +
        'Motyw produktu został nietknięty — zestaw jest bytem obok niego.',
      true,
    );
  }

  async function wczytaj(): Promise<void> {
    if (bezOkna(Command.DesignTokensetList)) return;
    odpowiedz.pokaz('Odczyt zestawów żetonów…', true);
    const wynik = await stan.czuwanie.prowadz(
      'odczyt zestawów żetonów',
      stan.zrodlo.zestawyZetonow(stan.idOkna(), ''),
      {
        wToku: (zdanie) => odpowiedz.pokaz(zdanie, true),
        cisza: (zdanie) => odpowiedz.pokaz(zdanie, false),
        powrot: (zdanie) => odpowiedz.pokaz(zdanie, false),
        spozniona: (zdanie) => odpowiedz.pokaz(zdanie, false),
      },
    );
    if (wynik === null) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowyBledu('Odczyt zestawów żetonów', wynik.blad), false);
      return;
    }
    zbior = wynik.wynik.tokenSets;
    pokazZestawy();
    odpowiedz.pokaz(
      zbior.length === 0
        ? 'Rdzeń nie zna ani jednego zestawu żetonów tego okna.'
        : `Zestawów w oknie: ${wynik.wynik.total}.`,
      true,
    );
  }

  async function wydajZestaw(): Promise<void> {
    if (wybor.kontrolka.value === '') {
      odpowiedz.pokaz('Wskaż zestaw — pozycja „nowy zestaw" nie ma czego wydać.', false);
      return;
    }
    odpowiedz.pokaz('Wydanie zestawu do kodu…', true);
    const wynik = await stan.czuwanie.prowadz(
      'wydanie zestawu żetonów',
      stan.zrodlo.wydajZetony(
        wybor.kontrolka.value,
        postac.kontrolka.value as DesignTokenTarget,
        modul.kontrolka.value.trim(),
      ),
      {
        wToku: (zdanie) => odpowiedz.pokaz(zdanie, true),
        cisza: (zdanie) => odpowiedz.pokaz(zdanie, false),
        powrot: (zdanie) => odpowiedz.pokaz(zdanie, false),
        spozniona: (zdanie) => odpowiedz.pokaz(zdanie, false),
      },
    );
    if (wynik === null) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowyBledu('Wydanie zestawu żetonów', wynik.blad), false);
      return;
    }
    // Treść wydania idzie na ekran w całości: „gotowe" bez pliku nie odróżnia
    // wydania od pustki.
    wydanie.textContent = wynik.wynik.content;
    odpowiedz.pokaz(
      `Rdzeń złożył plik „${wynik.wynik.fileName}".` +
        (wynik.wynik.delivered === true
          ? ` Wydanie trafiło także do modułu ${modul.kontrolka.value.trim()}.`
          : ''),
      true,
    );
  }

  async function wczytajZewnetrzny(): Promise<void> {
    if (bezOkna(Command.DesignTokensetImport)) return;
    const tresc = zapis.kontrolka.value.trim();
    if (tresc === '') {
      odpowiedz.pokaz('Wklej zapis do wczytania — puste pole nie niesie ani jednej roli.', false);
      return;
    }
    odpowiedz.pokaz('Wczytywanie zapisu żetonów…', true);
    const wynik = await stan.czuwanie.prowadz(
      'wczytanie zestawu żetonów',
      stan.zrodlo.wczytajZetony({
        idOkna: stan.idOkna(),
        nazwa: nazwa.kontrolka.value === '' ? 'Zestaw wczytany' : nazwa.kontrolka.value,
        trescBase64: wBaza64(tresc),
        // Postaci nie narzucamy: rdzeń rozpoznaje ją po treści, a wskazanie
        // z okna byłoby zgadywaniem za Operatora.
        postac: '',
      }),
      {
        wToku: (zdanie) => odpowiedz.pokaz(zdanie, true),
        cisza: (zdanie) => odpowiedz.pokaz(zdanie, false),
        powrot: (zdanie) => odpowiedz.pokaz(zdanie, false),
        spozniona: (zdanie) => odpowiedz.pokaz(zdanie, false),
      },
    );
    if (wynik === null) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowyBledu('Wczytanie zestawu żetonów', wynik.blad), false);
      return;
    }
    const zestaw = wynik.wynik.tokenSet;
    zbior = [zestaw, ...zbior];
    pokazZestawy();
    wybor.kontrolka.value = zestaw.id;
    const nieznane = wynik.wynik.unknownNames ?? [];
    odpowiedz.pokaz(
      `Rdzeń założył zestaw „${zestaw.name}" o ${zestaw.tokenCount} żetonach.` +
        (nieznane.length === 0
          ? ' Wszystkie role zapisu system produktu zna.'
          : ` Ról, których system produktu nie zna: ${nieznane.length} — ${nieznane.join(', ')}. ` +
            'Weszły do zestawu; wymieniamy je, żeby rozjazd nazw był widoczny.'),
      true,
    );
  }

  async function wydajPrzewodnik(): Promise<void> {
    if (wybor.kontrolka.value === '') {
      odpowiedz.pokaz('Wskaż zestaw — przewodnik powstaje z zestawu żetonów.', false);
      return;
    }
    if (modul.kontrolka.value.trim() === '') {
      odpowiedz.pokaz(
        'Wskaż moduł docelowy — przewodnik wydany donikąd jest plikiem, o którym nikt się ' +
          'nie dowie.',
        false,
      );
      return;
    }
    odpowiedz.pokaz('Wydanie przewodnika stylu…', true);
    const wynik = await stan.czuwanie.prowadz(
      'wydanie przewodnika stylu',
      // Kolekcji docelowej okno nie narzuca: kolekcje zakłada się w Assets
      // Panelu i tam się je wskazuje, a wymyślenie jej tutaj wydałoby
      // przewodnik do zbioru, którego Operator nie wybrał.
      stan.zrodlo.wydajPrzewodnik(wybor.kontrolka.value, modul.kontrolka.value.trim(), ''),
      {
        wToku: (zdanie) => odpowiedz.pokaz(zdanie, true),
        cisza: (zdanie) => odpowiedz.pokaz(zdanie, false),
        powrot: (zdanie) => odpowiedz.pokaz(zdanie, false),
        spozniona: (zdanie) => odpowiedz.pokaz(zdanie, false),
      },
    );
    if (wynik === null) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(opisOdmowyBledu('Wydanie przewodnika stylu', wynik.blad), false);
      return;
    }
    odpowiedz.pokaz(
      `Przewodnik leży w magazynie rdzenia jako zasób ${wynik.wynik.assetId}. ` +
        'Assets Panel pokaże go po odczycie zasobów; jego treść czyta się przyciskiem ' +
        '„Sprawdź treść w magazynie".',
      true,
    );
  }

  return { element, wczytaj };
}

/**
 * Koduje treść w base64 z zachowaniem znaków spoza ASCII.
 *
 * `btoa` przyjmuje wyłącznie bajty do 255, a zapis żetonów bywa po polsku
 * (opisy ról) — kodowanie wprost urywałoby się na pierwszym „ż". Droga przez
 * `TextEncoder` daje bajty UTF-8, czyli to, co rdzeń rozkłada po drugiej
 * stronie.
 */
function wBaza64(tekst: string): string {
  const bajty = new TextEncoder().encode(tekst);
  let zapis = '';
  for (const bajt of bajty) zapis += String.fromCharCode(bajt);
  return btoa(zapis);
}

/** Rodzaje żetonów kontraktu — pomocnik dla okna składającego materiał. */
export const RODZAJE_ZETONOW = DesignTokenKind;

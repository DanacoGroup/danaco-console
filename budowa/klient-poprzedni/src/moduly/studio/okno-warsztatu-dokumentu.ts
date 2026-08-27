import { Command, type DesignAsset } from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import {
  poleLogiczne,
  poleTajne,
  poleTekstowe,
  poleWielowierszowe,
  poleWyboru,
  przycisk,
  ustawPozycje,
  utworzWierszOdpowiedzi,
  type PozycjaWyboru,
} from '../../modele/kontrolki-formularza';
import {
  CZYNNOSCI_WARSZTATU,
  opiszSkutek,
  type CzynnoscWarsztatu,
  type PoleCzynnosci,
} from './czynnosci-warsztatu';
import { utworzOknoStudio } from './okno-studio';
import type { StanStudio } from './stan-studio';
import type { ZrodloWarsztatuDokumentu } from './zrodlo-warsztatu-dokumentu';

/** Wstążka PDF — zakładka kontekstowa rodzin studio.pdf.* i studio.security.*, domyślnie ukryta i wchodząca na żądanie albo samoczynnie przy dokumencie PDF. */
export interface OknoWarsztatuDokumentu {
  element: HTMLElement;
  /** Wczytuje wykaz materiału z magazynu okna. */
  wczytaj(): Promise<void>;
  // Przenosi ognisko do okna — woła to grupa Ochrona wstążki okna pracy.
  przenieOgnisko(): void;
  odswiez(): void;
}

/** Kontrolka pola wraz z odczytem jej wartości — element widoczny w formularzu i funkcja zwracająca wpisaną albo wybraną wartość jako tekst. */
export interface KontrolkaPola {
  element: HTMLElement;
  odczytaj(): string;
}

/** Buduje kontrolkę pola czynności i sposób odczytu jej wartości — wyeksportowana, bo korzysta z niej także okno redakcji dokumentu. */
export function utworzKontrolke(pole: PoleCzynnosci, zasoby: readonly DesignAsset[]): KontrolkaPola {
  const opis = {
    etykieta: pole.etykieta,
    ...(pole.opis === undefined ? {} : { opis: pole.opis }),
    ...(pole.podpowiedz === undefined ? {} : { podpowiedz: pole.podpowiedz }),
  };

  switch (pole.rodzaj) {
    case 'tajne': {
      const kontrolka = poleTajne(opis);
      return { element: kontrolka.element, odczytaj: () => kontrolka.kontrolka.value };
    }
    case 'liczba': {
      const kontrolka = poleTekstowe(opis);
      kontrolka.kontrolka.type = 'number';
      return { element: kontrolka.element, odczytaj: () => kontrolka.kontrolka.value };
    }
    case 'logiczne': {
      const kontrolka = poleLogiczne(opis);
      return {
        element: kontrolka.element,
        odczytaj: () => (kontrolka.kontrolka.checked ? 'tak' : 'nie'),
      };
    }
    case 'wielowiersz': {
      const kontrolka = poleWielowierszowe(opis, 4);
      return { element: kontrolka.element, odczytaj: () => kontrolka.kontrolka.value };
    }
    case 'wybor': {
      const kontrolka = poleWyboru(opis, pole.pozycje ?? []);
      return { element: kontrolka.element, odczytaj: () => kontrolka.kontrolka.value };
    }
    case 'zasob': {
      const kontrolka = poleWyboru(opis, pozycjeZasobow(zasoby));
      return { element: kontrolka.element, odczytaj: () => kontrolka.kontrolka.value };
    }
    case 'zasoby': {
      // Wielokrotny wybór zamiast listy pojedynczej — scalanie bierze materiały w kolejności zaznaczenia.
      const kontrolka = poleWyboru(opis, pozycjeZasobow(zasoby));
      kontrolka.kontrolka.multiple = true;
      kontrolka.kontrolka.size = 5;
      return {
        element: kontrolka.element,
        odczytaj: () =>
          Array.from(kontrolka.kontrolka.selectedOptions)
            .map((pozycja) => pozycja.value)
            .join(','),
      };
    }
    default: {
      const kontrolka = poleTekstowe(opis);
      return { element: kontrolka.element, odczytaj: () => kontrolka.kontrolka.value };
    }
  }
}

/** Pozycje listy materiału: nazwa zasobu wraz z jego formatem, gotowe do pokazania w polu wyboru czynności. */
function pozycjeZasobow(zasoby: readonly DesignAsset[]): PozycjaWyboru[] {
  return zasoby.map((zasob) => ({
    wartosc: zasob.id,
    etykieta: `${zasob.name ?? zasob.id}${zasob.format === undefined ? '' : ` (${zasob.format})`}`,
  }));
}

/** Grupa wstążki PDF wraz z czynnościami, które do niej należą: kod, nazwa widoczna, opis oraz komendy rdzenia zawężające wykaz czynności. */
export interface GrupaWstazkiPdf {
  kod: string;
  nazwa: string;
  opis: string;
  komendy: readonly Command[];
}

/** Grupy wstążki PDF w kolejności pracy nad plikiem: strony, nakładanie, treść i bezpieczeństwo przed wysyłką. */
export const GRUPY_WSTAZKI_PDF: readonly GrupaWstazkiPdf[] = [
  {
    kod: 'strony',
    nazwa: 'Strony',
    opis:
      'Łączenie wielu plików w jeden, dzielenie po stronach, zakresach albo zakładkach, ' +
      'porządkowanie stron (kolejność, obrót, usunięcie, wstawienie, wyodrębnienie), drzewo ' +
      'zakładek i zmniejszenie rozmiaru rekompresją.',
    komendy: [
      Command.StudioPdfMerge,
      Command.StudioPdfSplit,
      Command.StudioPdfPagesReorder,
      Command.StudioPdfBookmarksSet,
      Command.StudioPdfOptimize,
    ],
  },
  {
    kod: 'nakladanie',
    nazwa: 'Nakładanie',
    opis: 'Znak wodny, pieczęć tekstowa i graficzna oraz numeracja prawna Bates wraz z nagłówkami i stopkami.',
    komendy: [Command.StudioPdfStamp, Command.StudioPdfBates],
  },
  {
    kod: 'tresc',
    nazwa: 'Treść',
    opis: 'Formularze AcroForm — odczyt i wypełnienie — oraz wyciągnięcie osadzonych obrazów i załączników.',
    komendy: [Command.StudioPdfFormFill, Command.StudioPdfExtract],
  },
  {
    kod: 'bezpieczenstwo',
    nazwa: 'Bezpieczeństwo',
    opis:
      'Szyfrowanie i jego zdjęcie wraz z uprawnieniami, podpis certyfikatem i weryfikacja ' +
      'podpisów, trwałe zamazanie, czyszczenie autora i danych ukrytych przed wysyłką oraz ' +
      'oznaczenie numerów identyfikacyjnych, kart płatniczych i danych osobowych.',
    komendy: [
      Command.StudioSecurityEncrypt,
      Command.StudioSecuritySign,
      Command.StudioSecuritySignVerify,
      Command.StudioSecurityRedact,
      Command.StudioSecurityMetadataStrip,
      Command.StudioSecuritySensitiveDetect,
    ],
  },
];

/** Kod grupy narzędziowni cyfryzacji — jej treść przychodzi gniazdem przekazanym z zewnątrz, nie wykazem czynności katalogu. */
const GRUPA_CYFRYZACJI = 'cyfryzacja';

export function utworzOknoWarsztatuDokumentu(
  stan: StanStudio,
  zrodlo: ZrodloWarsztatuDokumentu,
  gniazdoCyfryzacji?: HTMLElement,
): OknoWarsztatuDokumentu {
  const rama = utworzOknoStudio({
    kod: 'studio.document-workshop',
    tytul: 'Wstążka PDF',
    rola: 'pomocnicze',
    objasnienie:
      'Zakładka kontekstowa: czynności na dokumencie PDF i jego bezpieczeństwie — scalanie, ' +
      'podział, strony, pieczęcie, numeracja Bates, zakładki, formularze, szyfrowanie, redakcja ' +
      'i podpis — oraz narzędziownia cyfryzacji. Wchodzi po przycisku i samoczynnie przy ' +
      'dokumencie PDF; nie zajmuje miejsca, gdy nie jest używana. PDF i krypto liczy biblioteka ' +
      'wkompilowana w rdzeń, bez ani jednego programu zewnętrznego. Materiał zostaje nietknięty; ' +
      'wynik jest nowym zasobem magazynu.',
  });

  /* ── Grupy wstążki ───────────────────────────────────────────────────────── */

  let grupaBiezaca = GRUPY_WSTAZKI_PDF[0]?.kod ?? 'strony';

  const zakladkiGrup = document.createElement('div');
  zakladkiGrup.className = 'ms-pdf__grupy';
  zakladkiGrup.setAttribute('role', 'tablist');
  zakladkiGrup.setAttribute('aria-label', 'Grupy wstążki PDF');

  const opisGrupy = document.createElement('p');
  opisGrupy.className = 'dn-pole-opis ms-pdf__opis-grupy';

  const przyciskiGrup = new Map<string, HTMLButtonElement>();
  const grupyDoZlozenia: readonly { kod: string; nazwa: string; opis: string }[] = [
    ...GRUPY_WSTAZKI_PDF.map((grupa) => ({ kod: grupa.kod, nazwa: grupa.nazwa, opis: grupa.opis })),
    {
      kod: GRUPA_CYFRYZACJI,
      nazwa: 'Narzędziownia cyfryzacji',
      opis:
        'Kolejka wczytywania (plik, katalog, archiwum), rozpoznanie tekstu ze sterowaniem ' +
        'silnikiem, językami i progiem pewności, poprawianie rozpoznanych słów przed przyjęciem ' +
        'do edytora, pobranie strony sieciowej z oczyszczeniem oraz wykaz skanerów i kamer.',
    },
  ];

  for (const grupa of grupyDoZlozenia) {
    const przyciskGrupy = document.createElement('button');
    przyciskGrupy.type = 'button';
    przyciskGrupy.className = 'ms-pdf__grupa';
    przyciskGrupy.textContent = grupa.nazwa;
    przyciskGrupy.dataset['grupa'] = grupa.kod;
    przyciskGrupy.setAttribute('role', 'tab');
    przyciskGrupy.title = grupa.opis;
    przyciskGrupy.addEventListener('click', () => ustawGrupe(grupa.kod));
    zakladkiGrup.append(przyciskGrupy);
    przyciskiGrup.set(grupa.kod, przyciskGrupy);
  }

  /** Gniazdo narzędziowni cyfryzacji — panel buduje inny wołacz, wstążka bierze jego element gniazdem. */
  const gniazdoGrupy = document.createElement('div');
  gniazdoGrupy.className = 'ms-pdf__gniazdo';
  if (gniazdoCyfryzacji !== undefined) {
    gniazdoGrupy.append(gniazdoCyfryzacji);
  } else {
    const wskazanie = document.createElement('p');
    wskazanie.className = 'dn-pole-opis';
    wskazanie.textContent =
      'Narzędziownia cyfryzacji stoi jako panel otwierany przyciskiem „Narzędziownia cyfryzacji" ' +
      'w module — wraz z kolejką wczytywania, rozpoznaniem tekstu, poprawianiem słów i przyjęciem ' +
      'wyniku do edytora. Osadzenie jej w tej grupie jest jednym wywołaniem w pliku składającym ' +
      'moduł (modul-studio.ts), który nie należy do tego odcinka prac — brak wpisany ' +
      'w sprawozdanie.';
    gniazdoGrupy.append(wskazanie);
  }

  const wybor = poleWyboru(
    {
      etykieta: 'Czynność',
      opis: 'Wybór przestawia pola poniżej — każda czynność bierze inne wskazania.',
    },
    CZYNNOSCI_WARSZTATU.map((czynnosc) => ({
      wartosc: czynnosc.komenda,
      etykieta: czynnosc.nazwa,
    })),
  );

  const objasnienieCzynnosci = document.createElement('p');
  objasnienieCzynnosci.className = 'dn-pole-opis ms-wskaznik';

  const formularz = document.createElement('div');
  formularz.className = 'ms-warsztat__pola';

  const wykonanie = przycisk('Wykonaj czynność', 'dn-btn dn-btn--sygnal');
  wykonanie.dataset['czynnosc'] = 'wykonaj';
  const odpowiedz = utworzWierszOdpowiedzi();

  const formularzCzynnosci = document.createElement('div');
  formularzCzynnosci.className = 'ms-pdf__czynnosc';
  formularzCzynnosci.append(wybor.element, objasnienieCzynnosci, formularz);

  rama.pasek.append(wykonanie);
  rama.stan.tresc.append(
    zakladkiGrup,
    opisGrupy,
    formularzCzynnosci,
    gniazdoGrupy,
    odpowiedz.element,
  );

  /* ── Ukrycie domyślne i wejście na żądanie ───────────────────────────────── */

  const nakladka = document.createElement('div');
  nakladka.className = 'ms-pdf__nakladka';
  nakladka.hidden = true;
  nakladka.append(rama.element);

  const wyzwalacz = document.createElement('button');
  wyzwalacz.type = 'button';
  wyzwalacz.className = 'dn-btn dn-btn--sm dn-btn--zarys ms-pdf__wyzwalacz';
  wyzwalacz.dataset['czynnosc'] = 'wstazka-pdf';
  wyzwalacz.textContent = 'Wstążka PDF ▾';
  wyzwalacz.setAttribute('aria-expanded', 'false');
  wyzwalacz.title =
    'Otwiera wstążkę PDF: strony, nakładanie, treść, bezpieczeństwo i narzędziownia cyfryzacji. ' +
    'Wstążka jest kontekstowa — wchodzi na żądanie, a przy dokumencie PDF wchodzi sama.';
  wyzwalacz.addEventListener('click', () => przestawWstazke(nakladka.hidden));

  const powloka = document.createElement('div');
  powloka.className = 'ms-pdf';
  powloka.append(wyzwalacz, nakladka);

  /** Czy Operator zamknął wstążkę ręcznie — samoczynne wejście tego nie łamie. */
  let zamknietaRecznie = false;

  function przestawWstazke(otwarta: boolean): void {
    nakladka.hidden = !otwarta;
    wyzwalacz.setAttribute('aria-expanded', otwarta ? 'true' : 'false');
    wyzwalacz.textContent = otwarta ? 'Wstążka PDF ▴' : 'Wstążka PDF ▾';
    zamknietaRecznie = !otwarta;
  }

  /** Czy dokument w pracy jest PDF-em — wtedy wstążka wchodzi sama. */
  function dokumentJestPdf(): boolean {
    return (stan.dokument()?.format ?? '').toLowerCase() === 'pdf';
  }

  /** Przestawia grupę czynną i zawęża wykaz czynności do jej obszaru. */
  function ustawGrupe(kod: string): void {
    grupaBiezaca = kod;
    for (const [nazwa, przyciskGrupy] of przyciskiGrup) {
      przyciskGrupy.dataset['czynna'] = nazwa === kod ? 'tak' : 'nie';
      przyciskGrupy.setAttribute('aria-selected', nazwa === kod ? 'true' : 'false');
    }
    opisGrupy.textContent =
      grupyDoZlozenia.find((grupa) => grupa.kod === kod)?.opis ?? '';
    const cyfryzacja = kod === GRUPA_CYFRYZACJI;
    formularzCzynnosci.hidden = cyfryzacja;
    wykonanie.hidden = cyfryzacja;
    gniazdoGrupy.hidden = !cyfryzacja;
    if (cyfryzacja) return;
    przestawPozycjeGrupy();
    przestawFormularz();
  }

  /** Wykaz czynności zawężony do grupy czynnej. */
  function czynnosciGrupy(): readonly CzynnoscWarsztatu[] {
    const grupa = GRUPY_WSTAZKI_PDF.find((pozycja) => pozycja.kod === grupaBiezaca);
    if (grupa === undefined) return CZYNNOSCI_WARSZTATU;
    return CZYNNOSCI_WARSZTATU.filter((czynnosc) => grupa.komendy.includes(czynnosc.komenda));
  }

  function przestawPozycjeGrupy(): void {
    const wykaz = czynnosciGrupy();
    ustawPozycje(
      wybor.kontrolka,
      wykaz.map((czynnosc) => ({ wartosc: czynnosc.komenda, etykieta: czynnosc.nazwa })),
    );
    const pierwsza = wykaz[0];
    if (pierwsza !== undefined) wybor.kontrolka.value = pierwsza.komenda;
  }

  let dokumenty: DesignAsset[] = [];
  let wszystkie: DesignAsset[] = [];
  let kontrolki = new Map<string, KontrolkaPola>();

  /** Czynność wskazana na liście; pierwsza z grupy, gdy lista dopiero się zakłada. */
  function czynnoscBiezaca(): CzynnoscWarsztatu {
    const wybrana = CZYNNOSCI_WARSZTATU.find((czynnosc) => czynnosc.komenda === wybor.kontrolka.value);
    return wybrana ?? (czynnosciGrupy()[0] ?? CZYNNOSCI_WARSZTATU[0]) as CzynnoscWarsztatu;
  }

  /** Przebudowuje formularz pod czynność wskazaną. */
  function przestawFormularz(): void {
    const czynnosc = czynnoscBiezaca();
    objasnienieCzynnosci.textContent = czynnosc.opis;
    kontrolki = new Map();
    const wiersze: HTMLElement[] = [];
    for (const pole of czynnosc.pola) {
      // Materiał czynności to dokument; pieczęć i certyfikat sięgają po wszystkie zasoby okna.
      const zrodloPozycji = pole.kod === 'assetId' || pole.kod === 'assetIds' ? dokumenty : wszystkie;
      const kontrolka = utworzKontrolke(pole, zrodloPozycji);
      kontrolki.set(pole.kod, kontrolka);
      wiersze.push(kontrolka.element);
    }
    formularz.replaceChildren(...wiersze);
  }

  /** Zbiera wartości formularza w postaci, którą czyta katalog czynności. */
  function zbierzWartosci(): Record<string, string> {
    const wartosci: Record<string, string> = {};
    for (const [kod, kontrolka] of kontrolki) wartosci[kod] = kontrolka.odczytaj();
    return wartosci;
  }

  async function wykonaj(): Promise<void> {
    const czynnosc = czynnoscBiezaca();
    const zlozenie = czynnosc.zloz(zbierzWartosci(), {
      idOkna: stan.idOkna(),
      idDokumentu: stan.dokument()?.id ?? null,
    });
    if ('odmowa' in zlozenie) {
      odpowiedz.pokaz(zlozenie.odmowa, false);
      return;
    }

    rama.stan.ladowanie(`${czynnosc.nazwa} — czynność w toku…`);
    const wynik = await zrodlo.wykonaj(czynnosc.komenda, zlozenie.zadanie);
    if (!wynik.udany) {
      const powod = opisOdmowy(czynnosc.nazwa, wynik.blad?.code, wynik.blad?.message);
      odpowiedz.pokaz(powod, false);
      rama.stan.gotowe();
      return;
    }
    odpowiedz.pokaz(`${czynnosc.nazwa}: ${opiszSkutek(wynik.wynik)}`, true);
    rama.stan.gotowe();
    // Wynik jest nowym zasobem magazynu — wykaz materiału zestarzał się, więc odczyt idzie od razu.
    await wczytaj();
  }

  async function wczytaj(): Promise<void> {
    const idOkna = stan.idOkna();
    if (idOkna === '') return;
    const [wynikDokumentow, wynikZasobow] = await Promise.all([
      zrodlo.dokumenty(idOkna),
      zrodlo.zasoby(idOkna),
    ]);
    if (!wynikDokumentow.udany || wynikDokumentow.wynik === undefined) {
      rama.stan.blad(
        opisOdmowy(
          'Odczyt materiału warsztatu',
          wynikDokumentow.blad?.code,
          wynikDokumentow.blad?.message,
        ),
      );
      return;
    }
    dokumenty = wynikDokumentow.wynik;
    wszystkie = wynikZasobow.udany && wynikZasobow.wynik !== undefined
      ? wynikZasobow.wynik
      : dokumenty;
    przestawFormularz();
    odswiez();
  }

  function odswiez(): void {
    // Wstążka wchodzi sama przy PDF-ie, ale nie wraca po ręcznym zamknięciu przez Operatora.
    if (dokumentJestPdf() && nakladka.hidden && !zamknietaRecznie) przestawWstazke(true);
    if (rama.stan.faza() === 'ladowanie') return;
    if (dokumenty.length === 0) {
      // Nazwa stanu mówi, czego brakuje: czynności są czynne, brakuje materiału do pracy.
      rama.stan.puste(
        'Wstążka PDF bez materiału',
        'Czynności pracują na dokumentach PDF wniesionych do okna. Wnieś dokument — ' +
          'narzędziownią cyfryzacji albo Assets Panelem modułu Design — a czynności zobaczą go ' +
          'na liście materiału.',
      );
      return;
    }
    przestawPozycjeGrupy();
    rama.stan.gotowe();
  }

  wybor.kontrolka.addEventListener('change', () => {
    odpowiedz.wyczysc();
    przestawFormularz();
  });
  wykonanie.addEventListener('click', () => void wykonaj());

  ustawGrupe(grupaBiezaca);

  return {
    element: powloka,
    wczytaj,

    // Otwiera wstążkę i przenosi do niej ognisko — otwarcie jest częścią tej czynności.
    przenieOgnisko() {
      przestawWstazke(true);
      rama.element.dataset['ognisko'] = 'tak';
      rama.element.scrollIntoView({ block: 'nearest' });
    },

    odswiez,
  };
}

import { pustaMigawka, type MigawkaDokumentu } from './pola-stanu';
import type { StanStudio } from './stan-studio';
import {
  utworzPamiecNastawWidoku,
  type TrybDwochDokumentow,
} from './widok-nastawy-operatora';
import type { ZrodloWstawienStudio } from './zrodlo-wstawien-studio';

/**
 * Zakładki dokumentów — jeden z DWÓCH równorzędnych trybów dwóch dokumentów.
 *
 * ── Dwa tryby, żaden zapasowy ───────────────────────────────────────────────
 * Rozstrzygnięcie Właściciela: zakładki i podział powierzchni są równorzędne,
 * a wybór należy do Operatora i jest pamiętany. Zakładki dają większe pole pracy
 * nad jednym pismem — dwie kartki A4 obok siebie schodzą do rozmiaru, w którym
 * pisma się nie czyta — a podział pokazuje oba naraz. Przełączenie trybu niczego
 * nie gubi i nie zamyka: stan każdego dokumentu siedzi w jego migawce niezależnie
 * od trybu, a dokument niewidoczny zostaje OTWARTY i dostępny modelowi.
 *
 * Sam podział powierzchni buduje `widok-podzialu-powierzchni.ts`; ten plik
 * prowadzi zakładki i pamięta, który tryb obowiązuje. Widok wielu stron obok
 * siebie jest rzeczą trzecią i dotyczy stron JEDNEGO dokumentu.
 *
 * ── Skąd niezależność stanu ─────────────────────────────────────────────────
 * Z migawki pól stanu modułu (`pola-stanu.ts`). Zakładka czynna JEST dokumentem
 * czynnym modułu — więc Tools Panel, Session Repository i Ingest/OCR Panel
 * pracują na tym, co Operator widzi. Zakładka odłożona trzyma swoją migawkę:
 * treść, zaznaczenie, propozycję i parę porównania. Przełączenie to podmiana
 * migawek, a nie druga kopia stanu modułu.
 */

/** Jedna zakładka: nazwa widoczna i odłożony stan dokumentu. */
export interface KartaDokumentu {
  /** Kod zakładki — trafia do `data-karta`. */
  kod: string;
  /** Nazwa widoczna: tytuł dokumentu albo zdanie o zakładce pustej. */
  nazwa(): string;
  /** Migawka odłożona; zakładka czynna ma migawkę nieaktualną z założenia. */
  migawka: MigawkaDokumentu;
}

/** Zakładki wraz z ich sterowaniem. */
export interface KartyDokumentow {
  element: HTMLElement;
  /** Zakłada zakładkę nową i czyni ją czynną. */
  dodaj(): void;
  /** Zamyka zakładkę wskazaną; ostatniej zamknąć nie można. */
  zamknij(kod: string): void;
  /** Czyni zakładkę czynną, odkładając stan poprzedniej. */
  przelacz(kod: string): void;
  /** Kod zakładki czynnej. */
  czynna(): string;
  /** Liczba zakładek otwartych. */
  ile(): number;
  /** Przerysowuje pasek zakładek — nazwy idą za tytułami dokumentów. */
  odswiez(): void;
  /** Tryb dwóch dokumentów obowiązujący teraz — nastawa pamiętana Operatora. */
  tryb(): TrybDwochDokumentow;
  /** Przestawia tryb i zapamiętuje wybór; nie zamyka i nie gubi dokumentów. */
  ustawTryb(tryb: TrybDwochDokumentow): void;
  /** Migawka zakładki wskazanej — czyta ją podział powierzchni dla pola drugiego. */
  migawka(kod: string): MigawkaDokumentu | null;
  /** Kody zakładek otwartych, w kolejności paska. */
  kody(): readonly string[];
}

/** Numer nadawany zakładkom kolejno w tej sesji okna. */
let licznikKart = 0;

/**
 * Zakładki wraz z drogą do PUSTEGO DOKUMENTU rdzenia.
 *
 * Źródło wstawień jest nieobowiązkowe, bo zakładki działają też bez niego —
 * i tak działały dotąd, zakładając zakładkę BEZ dokumentu. Gdy źródło jest,
 * przycisk „+ nowy dokument" woła `studio.document.create` i zakładka dostaje
 * pustą stronę gotową do pisania, z arkuszem stylów i nastawami strony
 * domyślnymi. To była jedyna rzecz, której rdzeń nie miał czym zrobić przed tą
 * turą: `document.open` wczytuje istniejący, `template.apply` zakłada z szablonu,
 * a pustego dokumentu nie zakładało nic.
 */
export function utworzKartyDokumentow(
  stan: StanStudio,
  poPrzelaczeniu: () => void,
  wstawienia?: ZrodloWstawienStudio,
): KartyDokumentow {
  const element = document.createElement('div');
  element.className = 'ms-karty';
  element.setAttribute('role', 'tablist');
  element.setAttribute('aria-label', 'Dokumenty otwarte w oknie pracy');

  function nowaKarta(): KartaDokumentu {
    licznikKart += 1;
    const kod = `karta-${licznikKart}`;
    const migawka = pustaMigawka();
    return {
      kod,
      migawka,
      nazwa() {
        const dokument = karty.find((karta) => karta.kod === kod);
        if (dokument === undefined) return kod;
        const wskazany = kod === czynna ? stan.dokument() : dokument.migawka.dokument;
        if (wskazany === null) return 'zakładka bez dokumentu';
        return wskazany.title ?? wskazany.id;
      },
    };
  }

  const karty: KartaDokumentu[] = [];
  const pierwsza = nowaKarta();
  karty.push(pierwsza);
  let czynna = pierwsza.kod;

  const pamiec = utworzPamiecNastawWidoku();
  let trybBiezacy: TrybDwochDokumentow = pamiec.biezace().trybDokumentow;

  /**
   * Plakietka trybu — mówi, w którym z dwóch trybów Operator pracuje.
   *
   * Stoi w pasku zakładek, a nie jako drugi przełącznik: przestawia się go na
   * pasku widoku powierzchni i dwa miejsca tej samej nastawy rozjechałyby się.
   */
  const plakietkaTrybu = document.createElement('span');
  plakietkaTrybu.className = 'dn-plakietka ms-karty__tryb';

  function opiszTryb(): void {
    plakietkaTrybu.textContent =
      trybBiezacy === 'podzial' ? 'podział powierzchni' : 'zakładki';
    plakietkaTrybu.title =
      trybBiezacy === 'podzial'
        ? 'Dokumenty stoją obok siebie, każdy z własnym zaznaczeniem i paskiem statusu. ' +
          'Tryb przestawia się na pasku widoku powierzchni; przełączenie niczego nie zamyka.'
        : 'Jeden dokument na całej powierzchni — większe pole pracy. Dokumenty pozostałych ' +
          'zakładek są OTWARTE, tylko niewidoczne, i model ma do nich dostęp. Tryb przestawia ' +
          'się na pasku widoku powierzchni.';
  }

  function odswiez(): void {
    opiszTryb();
    element.replaceChildren(
      plakietkaTrybu,
      ...karty.map((karta) => {
        const przycisk = document.createElement('button');
        przycisk.type = 'button';
        przycisk.className = 'ms-karty__karta';
        przycisk.dataset['karta'] = karta.kod;
        przycisk.dataset['czynna'] = karta.kod === czynna ? 'tak' : 'nie';
        przycisk.setAttribute('role', 'tab');
        przycisk.textContent = karta.nazwa();
        przycisk.addEventListener('click', () => przelacz(karta.kod));

        const zamkniecie = document.createElement('button');
        zamkniecie.type = 'button';
        zamkniecie.className = 'ms-karty__zamkniecie';
        zamkniecie.textContent = '✕';
        zamkniecie.title =
          karty.length === 1
            ? 'Ostatniej zakładki nie zamyka się — okno pracy prowadzi co najmniej jeden dokument.'
            : 'Zamyka zakładkę. Dokument zostaje w rdzeniu; zamyka się widok, nie pismo.';
        zamkniecie.addEventListener('click', (zdarzenie) => {
          zdarzenie.stopPropagation();
          zamknij(karta.kod);
        });

        const powloka = document.createElement('span');
        powloka.className = 'ms-karty__powloka';
        powloka.append(przycisk, zamkniecie);
        return powloka;
      }),
      przyciskDodania(),
      ...(wstawienia === undefined ? [] : [przyciskNowegoDokumentu()]),
      zdanieZalozenia,
    );
  }

  function przyciskDodania(): HTMLElement {
    const dodaj = document.createElement('button');
    dodaj.type = 'button';
    dodaj.className = 'ms-karty__dodaj';
    dodaj.textContent = '+ nowa zakładka';
    dodaj.dataset['czynnosc'] = 'dodaj-karte';
    dodaj.title =
      'Zakłada zakładkę bez dokumentu. Wskaż w niej dokument („Wczytaj dokument"), wnieś plik albo ' +
      'załóż pusty dokument obok — dwa pisma pracują wtedy obok siebie, każde z własnym zaznaczeniem.';
    dodaj.addEventListener('click', () => dodaj_karte());
    return dodaj;
  }

  /**
   * Zdanie o skutku założenia — stoi w pasku, bo tu Operator naciska.
   *
   * Odmowa rdzenia nie może zniknąć w ciszy: zakładka powstałaby wtedy pusta
   * i wyglądałaby jak dokument, którego nie ma.
   */
  const zdanieZalozenia = document.createElement('span');
  zdanieZalozenia.className = 'dn-pole-opis ms-karty__zdanie';
  zdanieZalozenia.hidden = true;

  function przyciskNowegoDokumentu(): HTMLElement {
    const zaloz = document.createElement('button');
    zaloz.type = 'button';
    zaloz.className = 'ms-karty__dodaj';
    zaloz.textContent = '+ nowy dokument';
    zaloz.dataset['czynnosc'] = 'nowy-dokument';
    zaloz.title =
      'Zakłada PUSTY dokument w rdzeniu komendą studio.document.create — nową stronę gotową do ' +
      'pisania, z arkuszem stylów i nastawami strony domyślnymi. Zakładka dostaje go od razu.';
    zaloz.addEventListener('click', () => void zalozNowyDokument());
    return zaloz;
  }

  async function zalozNowyDokument(): Promise<void> {
    if (wstawienia === undefined) return;
    if (stan.idOkna() === '') {
      pokazZdanie(
        'Założenie dokumentu należy do okna komunikacji sesji — wskaż je w pasie osadzenia modułu.',
        false,
      );
      return;
    }
    // Zakładka powstaje PRZED wywołaniem, żeby dokument nowy nie wyparł dokumentu,
    // nad którym Operator pracuje: nowa strona ma stanąć obok, a nie zamiast.
    dodaj_karte();
    const wynik = await wstawienia.zalozDokument({ windowId: stan.idOkna() });
    if (!wynik.udany || wynik.wynik === undefined) {
      pokazZdanie(
        'Rdzeń odmówił założenia pustego dokumentu' +
          (wynik.blad?.message === undefined ? '.' : `: ${wynik.blad.message}`) +
          ' Zakładka została założona i jest pusta — wskaż w niej dokument albo wnieś plik.',
        false,
      );
      return;
    }
    stan.wchlon(wynik.wynik.document);
    odswiez();
    pokazZdanie(
      `Pusty dokument ${wynik.wynik.document.id} założony i stoi w tej zakładce — strona gotowa do ` +
        'pisania, z arkuszem stylów i nastawami strony domyślnymi.',
      true,
    );
  }

  function pokazZdanie(tresc: string, udane: boolean): void {
    zdanieZalozenia.textContent = tresc;
    zdanieZalozenia.hidden = false;
    zdanieZalozenia.dataset['udane'] = udane ? 'tak' : 'nie';
  }

  function odlozCzynna(): void {
    const karta = karty.find((pozycja) => pozycja.kod === czynna);
    if (karta !== undefined) karta.migawka = stan.migawka();
  }

  function przelacz(kod: string): void {
    if (kod === czynna) return;
    const docelowa = karty.find((karta) => karta.kod === kod);
    if (docelowa === undefined) return;
    odlozCzynna();
    czynna = kod;
    stan.przywrocMigawke(docelowa.migawka);
    odswiez();
    poPrzelaczeniu();
  }

  function dodaj_karte(): void {
    odlozCzynna();
    const karta = nowaKarta();
    karty.push(karta);
    czynna = karta.kod;
    stan.przywrocMigawke(karta.migawka);
    odswiez();
    poPrzelaczeniu();
  }

  function zamknij(kod: string): void {
    if (karty.length === 1) return;
    const numer = karty.findIndex((karta) => karta.kod === kod);
    if (numer < 0) return;
    karty.splice(numer, 1);
    if (kod === czynna) {
      const nastepna = karty[Math.max(0, numer - 1)];
      if (nastepna !== undefined) {
        czynna = nastepna.kod;
        stan.przywrocMigawke(nastepna.migawka);
        poPrzelaczeniu();
      }
    }
    odswiez();
  }

  odswiez();

  return {
    element,
    dodaj: dodaj_karte,
    zamknij,
    przelacz,
    czynna: () => czynna,
    ile: () => karty.length,
    odswiez,
    tryb: () => trybBiezacy,

    ustawTryb(tryb) {
      if (tryb === trybBiezacy) return;
      trybBiezacy = tryb;
      pamiec.przestaw({ trybDokumentow: tryb });
      // Odłożenie stanu zakładki czynnej przed przełączeniem trybu: w podziale
      // pracują dwa pola naraz i pole drugie czyta migawkę, która musi być świeża.
      odlozCzynna();
      odswiez();
    },

    migawka(kod) {
      const karta = karty.find((pozycja) => pozycja.kod === kod);
      if (karta === undefined) return null;
      return kod === czynna ? stan.migawka() : karta.migawka;
    },

    kody: () => karty.map((karta) => karta.kod),
  };
}

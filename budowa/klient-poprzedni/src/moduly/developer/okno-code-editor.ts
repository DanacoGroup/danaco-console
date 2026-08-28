import type { DeveloperFile } from '../../../../shared/contract';
import { utworzRameOkna } from '../../komponenty/rama-okna';
import {
  poleTresci,
  przelacznik,
  przyciskAkcji as przycisk,
  wiersz,
} from '../../modele/kontrolki-formularza-braki';
import { osadzAkcjeKanoniczne, type AkcjaKanoniczna } from './akcje-kanoniczne';
import { opisOdmowyBledu } from '../../komponenty/odmowa';
import { rysujZaleznosci, zaleznosci } from './zaleznosci-zewnetrzne';
import type { StanDevelopera } from './stan-developer';
import { utworzStanTresci } from './stany-okna';
import { zdanieZapisuPliku } from './zdania-odpowiedzi';
import type { ZrodloDeveloper } from './zrodlo-developer';
import type { ZrodloWarsztatu } from './zrodlo-warsztatu';

/**
 * Code Editor — okno wiodące modułu Developer: edycja treści pliku wskazanego w Project
 * Tree i zapis z opcjonalnym założeniem wersji; stan pliku należy do modułu, nie do okna.
 */
export interface OknoCodeEditora {
  element: HTMLElement;
  odswiez(): void;
  /** Zamyka nasłuch `stan.naZmiane` założony przy konstrukcji okna. */
  zamknij(): void;
}

export function utworzOknoCodeEditora(
  zrodlo: ZrodloDeveloper,
  warsztat: ZrodloWarsztatu,
  stan: StanDevelopera,
): OknoCodeEditora {
  const rama = utworzRameOkna({
    tytul: 'Code Editor',
    rola: 'wiodące',
    kod: 'code-editor',
    przeznaczenie:
      'Edycja treści pliku wskazanego w Project Tree; zapis z opcjonalnym założeniem wersji.',
    modul: 'Developer',
    przedrostek: 'mdev',
  });
  const tresc = utworzStanTresci();
  // Wykaz programów spoza instalki stoi w pasku narzędzi — trzeba o nim wiedzieć przy planowaniu pracy.
  rama.narzedzia.append(
    rysujZaleznosci(zaleznosci(['serwer-jezyka', 'formater-linter', 'ripgrep'])),
  );
  const powierzchnia = zlozPowierzchnieEdytora(rama, tresc.element);
  const { otworz, zapisz, porzuc } = zlozCzynnosciEdytora(zrodlo, stan, tresc, powierzchnia);

  podepnijAkcjeEdytora(powierzchnia, { zapisz, porzuc });
  podepnijOperacjeKontekstowe(powierzchnia, warsztat, stan, tresc);

  const odsubskrybuj = stan.naZmiane(otworz);

  return {
    element: rama.element,
    odswiez: otworz,
    zamknij: odsubskrybuj,
  };
}

/** Czynności okna: odczyt pliku wraz z osadzeniem jego treści w polu edycji, zapis i porzucenie zmiany. */
interface CzynnosciEdytora {
  otworz(): void;
  zapisz(): void;
  porzuc(): void;
}

/**
 * Składa czynności okna, wydzielone z wytwórni, żeby ta zmieściła się pod
 * progiem długości funkcji. W domknięciu zostaje jedna zmienna — znacznik
 * czystości pola; ścieżka, treść i wersja pliku należą do stanu modułu.
 */
function zlozCzynnosciEdytora(
  zrodlo: ZrodloDeveloper,
  stan: StanDevelopera,
  tresc: ReturnType<typeof utworzStanTresci>,
  powierzchnia: PowierzchniaEdytora,
): CzynnosciEdytora {
  /** Jedyny stan, którego stan modułu nie zna: treść, względem której pole jest czyste. */
  let odniesienieCzystosci: string | null = null;

  function niezapisana(): boolean {
    return odniesienieCzystosci !== null && powierzchnia.kod.value !== odniesienieCzystosci;
  }

  function wczytajDoPola(plik: DeveloperFile): void {
    odniesienieCzystosci = plik.content ?? '';
    powierzchnia.sciezka.textContent = plik.path;
    powierzchnia.sciezka.dataset['wskazana'] = 'tak';
    opiszStatusPliku(powierzchnia.status, plik);
    const miejsce = tresc.tresc();
    powierzchnia.kod.value = odniesienieCzystosci;
    miejsce.append(powierzchnia.wiersz);
  }

  function otworz(): void {
    const sciezka = stan.sciezka();
    // Plik oddany przez rdzeń niesie własną ścieżkę — to on mówi, co leży w polu edycji.
    const wPolu = stan.plik();
    if (sciezka === '') {
      tresc.pusto('Nie wskazano pliku — wybierz go w Project Tree.');
      return;
    }
    if (wPolu !== null && wPolu.path === sciezka) return;
    if (niezapisana()) {
      // potwierdzenie, nie blad — blad czyściłby miejsce treści razem z polem .mdev-kod i pracą Operatora.
      tresc.potwierdzenie(
        `Plik „${wPolu?.path ?? ''}” ma niezapisaną zmianę w edytorze — zapisz ją albo porzuć, zanim otworzysz „${sciezka}”.`,
        false,
      );
      return;
    }
    tresc.ladowanie('Odczyt pliku…');
    void zrodlo.otworzPlik({ windowId: stan.okno(), path: sciezka }).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad('Rdzeń odmówił odczytu pliku.', wynik.blad);
        return;
      }
      // Kolejność ma znaczenie: wczytajDoPola osadza treść, zanim ustawPlik rozgłosi zmianę.
      wczytajDoPola(wynik.wynik);
      stan.ustawPlik(wynik.wynik);
    });
  }

  // zapisz i porzuc używają wyłącznie tresc.potwierdzenie — inne metody skasują pole .mdev-kod z pracą.
  function zapisz(): void {
    const wPolu = stan.plik();
    if (wPolu === null) {
      tresc.potwierdzenie(
        'Nie wiadomo, który plik zapisać — otwórz go najpierw z Project Tree.',
        false,
      );
      return;
    }
    // Wersja pliku sprzed żądania — porównanie z wersją po zapisie rozstrzyga, czy punkt powrotu powstał.
    const wersjaPrzedZapisem = wPolu.versionId ?? '';
    tresc.potwierdzenie('Zapis w toku…', true);
    void zrodlo
      .zapiszPlik({
        windowId: stan.okno(),
        path: wPolu.path,
        content: powierzchnia.kod.value,
        createVersion: powierzchnia.zaloz.checked,
      })
      .then((wynik) => {
        if (!wynik.udany || wynik.wynik === undefined) {
          tresc.potwierdzenie(`Rdzeń odmówił zapisu pliku. ${opisOdmowyZapisu(wynik.blad)}`, false);
          return;
        }
        // Zdanie powstaje przed stan.ustawPlik, zanim wersja w stanie modułu zostanie nadpisana.
        const potwierdzenie = zdanieZapisuPliku(
          wynik.wynik,
          wersjaPrzedZapisem,
          powierzchnia.zaloz.checked,
        );
        // Treść zapisana bierze się z odpowiedzi; gdy rdzeń jej nie odesłał, znacznikiem zostaje treść pola.
        odniesienieCzystosci = wynik.wynik.content ?? powierzchnia.kod.value;
        powierzchnia.sciezka.textContent = wynik.wynik.path;
        opiszStatusPliku(powierzchnia.status, wynik.wynik);
        stan.ustawPlik(wynik.wynik);
        tresc.potwierdzenie(potwierdzenie.zdanie, potwierdzenie.udane);
      });
  }

  /** Porzuca zmianę w polu i wraca do treści ostatnio wczytanej ze ścieżki bieżącej. */
  function porzuc(): void {
    if (odniesienieCzystosci === null) return;
    powierzchnia.kod.value = odniesienieCzystosci;
    tresc.potwierdzenie('', false);
    otworz();
  }

  return { otworz, zapisz, porzuc };
}

/** Zdanie o odmowie zapisu pliku do rdzenia, z kodem i treścią przyczyny kontraktu, gdy rdzeń je podał. */
function opisOdmowyZapisu(powod: { code: string; message: string } | undefined): string {
  if (powod === undefined) return '';
  const tresc = powod.message === '' ? 'rdzeń nie podał przyczyny' : powod.message;
  return `Powód: ${tresc} (kod ${powod.code}).`;
}

/** Kontrolki paska akcji Code Editora: przyciski zapisu, porzucenia zmiany oraz operacji kontekstowych. */
interface AkcjeEdytora {
  /** Uchwyt wykonania operacji kontekstowych panelu akcji. */
  operacje: UchwytOperacji;
  zapiszPrzycisk: HTMLButtonElement;
  porzucPrzycisk: HTMLButtonElement;
}

/**
 * Składa pasek akcji: zapis, porzucenie zmiany i dziewięć operacji kontekstowych bez
 * pokrycia w kontrakcie; wydzielone z wytwórni okna pod progiem długości funkcji.
 */
function zlozAkcjeEdytora(gospodarz: HTMLElement): AkcjeEdytora {
  const zapiszPrzycisk = przycisk('Zapisz plik', 'dn-btn dn-btn--atrament');
  const porzucPrzycisk = przycisk('Porzuć zmianę i otwórz wskazany plik', 'dn-btn dn-btn--zarys');
  gospodarz.append(zapiszPrzycisk, porzucPrzycisk);
  // Wykonanie podpina się później przez uchwyt zmienny — stan treści jeszcze nie istnieje.
  const uchwyt: UchwytOperacji = { wykonaj: () => undefined };
  osadzAkcjeKanoniczne(gospodarz, (akcja) => uchwyt.wykonaj(akcja));
  return { zapiszPrzycisk, porzucPrzycisk, operacje: uchwyt };
}

/** Uchwyt wykonania operacji kontekstowej panelu akcji — wypełniany dopiero po złożeniu całego okna edytora. */
interface UchwytOperacji {
  wykonaj(akcja: AkcjaKanoniczna): void;
}

/** Kontrolki okna: pasek akcji, pole edycji treści pliku oraz przełącznik wersji zapisu obok niego samego. */
interface PowierzchniaEdytora extends AkcjeEdytora {
  sciezka: HTMLElement;
  /** Pasek statusu — wyłącznie pola, które rdzeń o pliku naprawdę powiedział. */
  status: HTMLElement;
  kod: HTMLTextAreaElement;
  zaloz: HTMLInputElement;
  wiersz: HTMLElement;
}

/**
 * Pasek statusu edytora złożony z pól `DeveloperFile`: język, rozmiar, wersja i czas
 * ostatniej zmiany, jedyne pola, które kontrakt niesie.
 */
function opiszStatusPliku(miejsce: HTMLElement, plik: DeveloperFile): void {
  const czesci: string[] = [];
  czesci.push(plik.language === undefined || plik.language === ''
    ? 'język: rdzeń nie podał'
    : `język: ${plik.language}`);
  czesci.push(plik.sizeBytes === undefined
    ? 'rozmiar: rdzeń nie podał'
    : `rozmiar: ${plik.sizeBytes} B`);
  czesci.push(plik.versionId === undefined || plik.versionId === ''
    ? 'wersja: brak punktu powrotu'
    : `wersja: ${plik.versionId}`);
  czesci.push(`zmieniono: ${new Date(plik.updatedAt).toLocaleString('pl-PL')}`);
  // Zgłoszenia lintera i kodowanie nie mają w kontrakcie pola — pasek mówi to wprost, nie zero domyślne.
  czesci.push('zgłoszenia lintera i kodowanie: niezmierzone, kontrakt nie ma tych pól');
  miejsce.textContent = czesci.join(' · ');
  miejsce.dataset['wskazana'] = 'tak';
}

/**
 * Składa kontrolki okna i osadza pasek ścieżki oraz przełącznik wersji w ciele
 * ramy. Pole kodu trafia do treści dopiero po udanym odczycie — czysta
 * konstrukcja, nie domyka się na stanie modułu ani na źródle.
 */
function zlozPowierzchnieEdytora(
  rama: { akcje: HTMLElement; cialo: HTMLElement },
  stanTresci: HTMLElement,
): PowierzchniaEdytora {
  const akcje = zlozAkcjeEdytora(rama.akcje);
  const sciezka = document.createElement('p');
  sciezka.className = 'mdev-sciezka';
  sciezka.dataset['wskazana'] = 'nie';
  sciezka.textContent = 'Brak wskazanego pliku.';

  const status = document.createElement('p');
  status.className = 'mdev-sciezka mdev-status';
  status.dataset['wskazana'] = 'nie';
  status.textContent = 'Pasek statusu wypełni się po odczycie pliku.';

  const kod = poleTresci('Treść pliku', 20, '', 'mdev-kod');
  const zaloz = przelacznik('Załóż nową wersję przy zapisie');
  const wierszWersji = wiersz('Wersja', zaloz, {
    klasa: 'mdev-wiersz',
    objasnienie: 'Decyzja Operatora — zapis bez tego pola nie zakłada wersji pliku.',
  });

  const edytor = document.createElement('div');
  edytor.className = 'mdev-edytor';
  edytor.append(kod, wierszWersji);

  rama.cialo.append(sciezka, status, stanTresci);
  return { ...akcje, sciezka, status, kod, zaloz, wiersz: edytor };
}

/** Podpina pasek akcji do czynności okna edytora: zapis, porzucenie zmiany i operacje kontekstowe rdzenia. */
function podepnijAkcjeEdytora(
  powierzchnia: PowierzchniaEdytora,
  obsluga: { zapisz: () => void; porzuc: () => void },
): void {
  powierzchnia.zapiszPrzycisk.addEventListener('click', obsluga.zapisz);
  powierzchnia.porzucPrzycisk.addEventListener('click', obsluga.porzuc);
}

/**
 * Podpina dziewięć operacji panelu akcji do komendy `developer.contextual.op`: wynik
 * wraca propozycją, okno pokazuje go w miejscu treści, nie wstawia do pola edycji.
 */
function podepnijOperacjeKontekstowe(
  powierzchnia: PowierzchniaEdytora,
  warsztat: ZrodloWarsztatu,
  stan: StanDevelopera,
  tresc: ReturnType<typeof utworzStanTresci>,
): void {
  powierzchnia.operacje.wykonaj = (akcja) => {
    void wykonajOperacjeKontekstowa(akcja, powierzchnia, warsztat, stan, tresc);
  };
}

/** Wykonuje jedną operację kontekstową panelu akcji i pokazuje w polu treści jej wynik albo odmowę rdzenia. */
async function wykonajOperacjeKontekstowa(
  akcja: AkcjaKanoniczna,
  powierzchnia: PowierzchniaEdytora,
  warsztat: ZrodloWarsztatu,
  stan: StanDevelopera,
  tresc: ReturnType<typeof utworzStanTresci>,
): Promise<void> {
  const sciezka = stan.sciezka();
  if (sciezka === '') {
    tresc.blad('Operacja dotyczy pliku — wskaż go w Project Tree albo otwórz w edytorze.');
    return;
  }
  const zaznaczenie = powierzchnia.kod.value.slice(
    powierzchnia.kod.selectionStart,
    powierzchnia.kod.selectionEnd,
  );

  tresc.ladowanie(`${akcja.nazwa} — operacja w toku…`);
  const zadanie: Parameters<ZrodloWarsztatu['operacjaKontekstowa']>[0] = {
    windowId: stan.okno(),
    operation: akcja.rodzaj,
    path: sciezka,
  };
  if (zaznaczenie.trim() !== '') zadanie.selection = zaznaczenie;

  const wynik = await warsztat.operacjaKontekstowa(zadanie);
  if (!wynik.udany || wynik.wynik === undefined) {
    tresc.blad(opisOdmowyBledu(akcja.nazwa, wynik.blad), wynik.blad);
    return;
  }
  const miejsce = tresc.tresc();
  const zdanie = document.createElement('p');
  zdanie.className = 'mdev-wiersz';
  zdanie.textContent = `${akcja.nazwa} — wynik jest propozycją; okno niczego nie zapisało.`;
  const blok = document.createElement('pre');
  blok.className = 'mdev-blok';
  blok.textContent = wynik.wynik.result;
  miejsce.replaceChildren(zdanie, blok);
}

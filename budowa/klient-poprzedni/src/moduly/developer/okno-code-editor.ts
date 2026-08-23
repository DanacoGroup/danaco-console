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
 * Code Editor — okno wiodące modułu Developer: edycja treści pliku wskazanego
 * w Project Tree i zapis z opcjonalnym założeniem wersji.
 *
 * Dziewięć operacji kontekstowych panelu akcji jedzie komendą
 * `developer.contextual.op`: rdzeń składa polecenie z treści pliku, zaznaczenia
 * i wskazania Operatora, a wynik wraca PROPOZYCJĄ — okno pokazuje ją w miejscu
 * treści i niczego nie zapisuje. Nazwy kanoniczne i rodzaje operacji stoją
 * w `akcje-kanoniczne.ts`, jednym wykazem dla całego modułu.
 *
 * Stan należy do modułu, nie do tego okna: ścieżkę pliku i jego wersję okno
 * czyta z `stan.plik()`. Własną trzyma wyłącznie treść, względem której pole
 * jest czyste — bywa nią treść pola, a nie treść z rdzenia.
 *
 * Wskazanie pliku przychodzi ze stanu. Project Tree woła `stan.wskazPlik(...)`;
 * to okno nasłuchuje `stan.naZmiane(...)` i na nową ścieżkę otwiera plik przez
 * `zrodlo.otworzPlik`. Po udanym odczycie i po zapisie woła `stan.ustawPlik`,
 * żeby Project Tree i Git Panel widziały to samo.
 *
 * Jedyny punkt wejścia do odczytu to `odswiez()`. Wytwórnia okna nie czyta
 * sama — złożenie modułu woła `odswiez()` raz po montażu, a odczyt w wytwórni
 * dałby podwójne żądanie. `odswiez()` nie zdejmuje subskrypcji `stan.naZmiane`:
 * ta żyje od konstrukcji do `zamknij()`, inaczej wskazanie pliku w Project Tree
 * przestałoby cokolwiek robić po pierwszym odświeżeniu.
 *
 * Okno nie gubi pracy Operatora cicho. Gdy treść w polu edycji różni się od
 * treści ostatnio wczytanego pliku, a stan wskazuje inną ścieżkę, okno woła
 * `tresc.potwierdzenie(...)`, nie `tresc.blad(...)` — ten drugi czyściłby
 * miejsce treści razem z polem edycji. Pole `.mdev-kod` zostaje widoczne
 * i edytowalne, a Operator ma dwa jawne wyjścia: „Zapisz plik” albo „Porzuć
 * zmianę i otwórz wskazany plik”.
 *
 * Potwierdzenie zapisu mówi, co zrobił rdzeń: zdanie składa
 * `zdania-odpowiedzi.ts` z pól odpowiedzi, nie z przełącznika wersji. Rdzeń
 * robi wersję z treści sprzed zapisu, więc gdy takiej nie było, `versionId`
 * nie wraca wcale, a potwierdzenie nie ma prawa obiecywać punktu powrotu.
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
  // Wykaz programów spoza instalki stoi w pasku narzędzi, nie w podpowiedzi
  // przycisku: o wymaganym serwerze języka trzeba wiedzieć przy planowaniu
  // pracy, a nie dopiero z odmowy czynności, która go potrzebowała.
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

/** Czynności okna: odczyt (wraz z osadzeniem w polu), zapis i porzucenie zmiany. */
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
  /**
   * Jedyny stan, którego stan modułu nie zna, i dlatego jedyny trzymany tutaj:
   * treść, względem której pole edycji jest czyste.
   *
   * Nie jest to treść z rdzenia. Po zapisie, którego rdzeń nie potwierdził
   * treścią, znacznikiem czystości zostaje treść pola — inaczej edytor
   * uznawałby pracę za niezapisaną na zawsze — a stan modułu nie ma prawa
   * takiej treści nieść, bo rdzeń jej nie odesłał. `null` znaczy „nic jeszcze
   * nie wczytano”.
   */
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
    // Plik oddany przez rdzeń niesie własną ścieżkę — to on, a nie osobna kopia
    // w domknięciu, mówi, co leży w polu edycji.
    const wPolu = stan.plik();
    if (sciezka === '') {
      tresc.pusto('Nie wskazano pliku — wybierz go w Project Tree.');
      return;
    }
    if (wPolu !== null && wPolu.path === sciezka) return;
    if (niezapisana()) {
      // `potwierdzenie`, nie `blad`: `blad` czyściłby miejsce treści i
      // skasowałby razem z nim pole `.mdev-kod` — czyli dokładnie tę pracę
      // Operatora, przed której utratą to ostrzeżenie ma bronić.
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
      // Kolejność ma znaczenie: `wczytajDoPola` musi osadzić treść w polu zanim
      // `stan.ustawPlik` rozgłosi zmianę, bo nasłuch wywołuje `otworz` ponownie.
      // Reentrantne wywołanie kończy się wtedy natychmiast na porównaniu
      // `stan.plik()?.path === stan.sciezka()` — obie wartości ustawia
      // `ustawPlik` przed powiadomieniem — więc drugie żądanie odczytu nie leci.
      wczytajDoPola(wynik.wynik);
      stan.ustawPlik(wynik.wynik);
    });
  }

  // `zapisz` i `porzuc` używają wyłącznie `tresc.potwierdzenie(...)`, nigdy
  // `tresc.ladowanie`/`tresc.blad`: obie czyszczą miejsce treści i skasowałyby
  // `.mdev-kod` razem z pracą Operatora — tą samą pracą, którą zapis ma utrwalić.
  function zapisz(): void {
    const wPolu = stan.plik();
    if (wPolu === null) {
      tresc.potwierdzenie(
        'Nie wiadomo, który plik zapisać — otwórz go najpierw z Project Tree.',
        false,
      );
      return;
    }
    // Wersja pliku wedle ostatniej odpowiedzi rdzenia, zdjęta ze stanu przed
    // wysłaniem żądania. Porównanie jej z wersją oddaną po zapisie rozstrzyga,
    // czy punkt powrotu naprawdę powstał — obecność pola tego nie rozstrzyga,
    // bo rdzeń dokłada tam wersję zastaną (`zdania-odpowiedzi.ts`).
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
        // Zdanie powstaje przed `stan.ustawPlik`, czyli zanim wersja w stanie
        // modułu zostanie nadpisana: rozstrzyga je porównanie wersji sprzed
        // zapisu z wersją oddaną przez rdzeń.
        const potwierdzenie = zdanieZapisuPliku(
          wynik.wynik,
          wersjaPrzedZapisem,
          powierzchnia.zaloz.checked,
        );
        // Treść zapisana bierze się z odpowiedzi. Gdy rdzeń jej nie odesłał,
        // znacznikiem czystości zostaje treść pola — inaczej edytor uznawałby
        // pracę za niezapisaną na zawsze — a zdanie potwierdzenia mówi wprost,
        // że okno nie ma czym potwierdzić zgodności z dyskiem. Znacznik
        // ustawiamy przed `ustawPlik`, bo ten rozgłasza zmianę i budzi `otworz`.
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

/** Zdanie o odmowie zapisu, z kodem kontraktu, gdy rdzeń go podał. */
function opisOdmowyZapisu(powod: { code: string; message: string } | undefined): string {
  if (powod === undefined) return '';
  const tresc = powod.message === '' ? 'rdzeń nie podał przyczyny' : powod.message;
  return `Powód: ${tresc} (kod ${powod.code}).`;
}

/** Kontrolki paska akcji Code Editora. */
interface AkcjeEdytora {
  /** Uchwyt wykonania operacji kontekstowych panelu akcji. */
  operacje: UchwytOperacji;
  zapiszPrzycisk: HTMLButtonElement;
  porzucPrzycisk: HTMLButtonElement;
}

/**
 * Składa pasek akcji: zapis, porzucenie zmiany i dziewięć operacji
 * kontekstowych bez pokrycia w kontrakcie.
 *
 * Wydzielone z wytwórni okna, żeby ta zmieściła się pod progiem długości
 * funkcji — fragment jest czystą konstrukcją bez domknięcia na stanie.
 *
 * „Porzuć zmianę” jest jawnym drugim wyjściem z ostrzeżenia o niezapisanej
 * pracy (obok „Zapisz plik”) — bez niego Operator widziałby ostrzeżenie bez
 * żadnej drogi naprzód poza ręczną edycją treści z powrotem do stanu wyjściowego.
 *
 * Nazwy operacji nie stoją tutaj — mają jedno źródło w `akcje-kanoniczne.ts`,
 * żeby przycisk i komenda dołożona później mówiły o operacji tym samym słowem.
 */
function zlozAkcjeEdytora(gospodarz: HTMLElement): AkcjeEdytora {
  const zapiszPrzycisk = przycisk('Zapisz plik', 'dn-btn dn-btn--atrament');
  const porzucPrzycisk = przycisk('Porzuć zmianę i otwórz wskazany plik', 'dn-btn dn-btn--zarys');
  gospodarz.append(zapiszPrzycisk, porzucPrzycisk);
  // Wykonanie podpina się później — pasek powstaje razem z powierzchnią okna,
  // a droga do rdzenia potrzebuje stanu treści, którego wtedy jeszcze nie ma.
  // Dlatego wykonanie idzie przez uchwyt zmienny, a nie przez drugi zestaw
  // nasłuchów doklejany do tych samych przycisków.
  const uchwyt: UchwytOperacji = { wykonaj: () => undefined };
  osadzAkcjeKanoniczne(gospodarz, (akcja) => uchwyt.wykonaj(akcja));
  return { zapiszPrzycisk, porzucPrzycisk, operacje: uchwyt };
}

/** Uchwyt wykonania operacji kontekstowej — wypełniany po złożeniu okna. */
interface UchwytOperacji {
  wykonaj(akcja: AkcjaKanoniczna): void;
}

/** Kontrolki okna: pasek akcji, pole edycji i przełącznik wersji zapisu. */
interface PowierzchniaEdytora extends AkcjeEdytora {
  sciezka: HTMLElement;
  /** Pasek statusu — wyłącznie pola, które rdzeń o pliku naprawdę powiedział. */
  status: HTMLElement;
  kod: HTMLTextAreaElement;
  zaloz: HTMLInputElement;
  wiersz: HTMLElement;
}

/**
 * Pasek statusu edytora złożony z pól `DeveloperFile`.
 *
 * Opracowanie wymienia w tym pasku język, kodowanie, znaki końca linii,
 * wcięcia, liczbę zgłoszeń lintera i gałąź. Kontrakt niesie z tego wyłącznie
 * język, rozmiar, wersję i czas ostatniej zmiany — i tylko te pola pasek
 * pokazuje. Reszta jest tu NAZWANA jako niezmierzona, zamiast być pokazana
 * wartością domyślną: „UTF-8” wypisane bez odczytu byłoby zgadywaniem, a „0
 * błędów” bez lintera mówiłoby o pliku sprawdzonym, choć nikt go nie sprawdzał.
 *
 * Pole nieobecne w odpowiedzi znaczy „rdzeń tego nie podał” i tak jest opisane;
 * milczenie rdzenia nie zamienia się tu w wartość.
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
  // Zgłoszenia lintera i kodowanie nie mają w kontrakcie ani jednego pola —
  // pasek mówi to wprost zamiast wypisywać zero albo wartość domyślną.
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

/** Podpina pasek akcji do czynności okna. */
function podepnijAkcjeEdytora(
  powierzchnia: PowierzchniaEdytora,
  obsluga: { zapisz: () => void; porzuc: () => void },
): void {
  powierzchnia.zapiszPrzycisk.addEventListener('click', obsluga.zapisz);
  powierzchnia.porzucPrzycisk.addEventListener('click', obsluga.porzuc);
}

/**
 * Podpina dziewięć operacji panelu akcji do komendy `developer.contextual.op`.
 *
 * Zaznaczeniem jest to, co Operator zaznaczył w polu edycji; brak zaznaczenia
 * znaczy operację nad całym plikiem. Wynik wraca propozycją i okno pokazuje go
 * w miejscu treści — nie wstawia go do pola edycji, bo to odebrałoby Operatorowi
 * tę jedną chwilę, w której da się pracę modelu odrzucić.
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

/** Wykonuje jedną operację kontekstową i pokazuje jej wynik albo odmowę. */
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
